package marketmaker

import (
	"context"

	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/regources"
)

type AssetsGetter interface {
	Pairs() ([]regources.AssetPair, error)
}

type AccountInfoProvider interface {
	Offers(address, baseAsset, quoteAsset string, isBuy *bool, offerID string, orderBookID *uint64) ([]regources.Offer, error)
	Balances(address string) ([]horizon.Balance, error)
}

type Submitter interface {
	SubmitE(txEnvelope string) (horizon.SubmitResponseDetails, error)
}

type Service struct {
	log    *logan.Entry
	config Config

	assetsGetter        AssetsGetter
	accountInfoProvider AccountInfoProvider
	submitter           Submitter
	builder             *xdrbuild.Builder

	assetToBalanceID map[string]string
}

func NewService(
	log *logan.Entry,
	config Config,
	assetsGetter AssetsGetter,
	offersGetter AccountInfoProvider,
	submitter Submitter,
	builder *xdrbuild.Builder,
) *Service {
	return &Service{
		log:                 log.WithField("service", conf.ServiceMarketMaker),
		config:              config,
		assetsGetter:        assetsGetter,
		accountInfoProvider: offersGetter,
		submitter:           submitter,
		builder:             builder,

		assetToBalanceID: make(map[string]string),
	}
}

func (s *Service) Run(ctx context.Context) {
	s.log.WithField("c", s.config).Info("Starting.")

	running.UntilSuccess(ctx, s.log, "balances_obtainer", s.obtainBalances, 5*time.Second, 5*time.Minute)
	if running.IsCancelled(ctx) {
		return
	}

	running.WithBackOff(ctx, s.log, "offers_refresher_iteration", func(ctx context.Context) error {
		for _, assetPair := range s.config.AssetPairs {
			err := s.refreshOffers(ctx, assetPair)
			if err != nil {
				return errors.Wrap(err, "Failed to refresh Offer for AssetPair", logan.F{
					"asset_pair": assetPair,
				})
			}
		}

		return nil
	}, s.config.CheckPeriod, 5*time.Second, time.Minute)
}

func (s *Service) obtainBalances(ctx context.Context) (bool, error) {
	balances, err := s.accountInfoProvider.Balances(s.config.Source.Address())
	if err != nil {
		return false, errors.Wrap(err, "failed to get Account Balances")
	}

	managedAssets := make(map[string]struct{})
	for _, assetPair := range s.config.AssetPairs {
		managedAssets[assetPair.BaseAsset] = struct{}{}
		managedAssets[assetPair.QuoteAsset] = struct{}{}
	}

	for _, balance := range balances {
		if _, ok := managedAssets[balance.Asset]; ok {
			// This is asset exists in config
			s.assetToBalanceID[balance.Asset] = balance.BalanceID
		}
	}

	// Check all assets from config were found
	for asset, _ := range managedAssets {
		if _, ok := s.assetToBalanceID[asset]; !ok {
			return false, errors.Errorf("Balance for the Asset (%s) was not found for my Account (%s).", asset, s.config.Source.Address())
		}
	}

	return true, nil
}

func (s *Service) refreshOffers(ctx context.Context, assetPairConfig AssetPairConfig) error {
	currentPrice, err := s.getCurrentPrice(assetPairConfig.BaseAsset, assetPairConfig.QuoteAsset)
	if err != nil {
		return errors.Wrap(err, "failed to obtain current price of the AssetPair")
	}

	buyPriceToOffer, overflow := amount.BigDivide(int64(*currentPrice), amount.One-int64(assetPairConfig.PriceMargin), amount.One, amount.ROUND_DOWN)
	if overflow {
		return errors.New("Overflow on counting buyPriceToOffer.")
	}
	sellPriceToOffer, overflow := amount.BigDivide(int64(*currentPrice), amount.One+int64(assetPairConfig.PriceMargin), amount.One, amount.ROUND_UP)
	if overflow {
		return errors.New("Overflow on counting sellPriceToOffer.")
	}

	tx := s.builder.Transaction(s.config.Source)

	needNewBuyOffer, err := s.removeBuyOffersIfNecessary(ctx, assetPairConfig, buyPriceToOffer, tx)
	if err != nil {
		return errors.Wrap(err, "failed to refresh buy Offer")
	}

	needNewSellOffer, err := s.removeSellOffersIfNecessary(ctx, assetPairConfig, sellPriceToOffer, tx)
	if err != nil {
		return errors.Wrap(err, "failed to refresh sell Offer")
	}

	if needNewBuyOffer {
		baseAmount, overflow := amount.BigDivide(amount.One, int64(assetPairConfig.QuoteAssetVolume), buyPriceToOffer, amount.ROUND_UP)
		if overflow {
			return errors.From(errors.New("Conversion to BaseAmount caught overflow."), logan.F{
				"quote_asset_volume": assetPairConfig.QuoteAssetVolume,
				"buy_price_to_offer": buyPriceToOffer,
			})
		}

		op := xdrbuild.CreateOffer(s.assetToBalanceID[assetPairConfig.BaseAsset], s.assetToBalanceID[assetPairConfig.QuoteAsset],
			true, baseAmount, buyPriceToOffer, 0)
		tx.Op(op)

		s.log.WithFields(logan.F{
			"manage_offer_op": op,
		}).Info("Creating new buy Offer.")
	}

	if needNewSellOffer {
		op := xdrbuild.CreateOffer(s.assetToBalanceID[assetPairConfig.BaseAsset], s.assetToBalanceID[assetPairConfig.QuoteAsset],
			false, int64(assetPairConfig.BaseAssetVolume), sellPriceToOffer, 0)
		tx.Op(op)

		s.log.WithFields(logan.F{
			"manage_offer_op": op,
		}).Info("Creating new sell Offer.")
	}

	envelope, err := tx.Sign(s.config.Signer).Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to marshal TX")
	}

	responseDetails, err := s.submitter.SubmitE(envelope)
	if err != nil {
		return errors.Wrap(err, "failed to submit tx", logan.F{
			"details": responseDetails,
		})
	}

	return nil
}

func (s *Service) removeBuyOffersIfNecessary(ctx context.Context, assetPairConfig AssetPairConfig, priceToOffer int64, tx *xdrbuild.Transaction) (needNewOffer bool, err error) {
	if assetPairConfig.QuoteAssetVolume == 0 {
		// Don't manage buy Offers
		return false, nil
	}

	t := true
	z := uint64(0)
	offers, err := s.accountInfoProvider.Offers(s.config.Source.Address(), assetPairConfig.BaseAsset, assetPairConfig.QuoteAsset, &t, "", &z)
	if err != nil {
		return false, errors.Wrap(err, "failed to obtain my Offers")
	}

	if len(offers) == 1 && int64(offers[0].Price) == priceToOffer && offers[0].QuoteAmount >= assetPairConfig.QuoteAssetVolume {
		// No need to refresh the Offer - it is full and price is actual.
		s.log.WithFields(logan.F{
			"offer":      offers[0],
			"asset_pair": assetPairConfig,
		}).Debug("Buy Offer is not outdated, leaving it as is.")
		return false, nil
	}

	for _, o := range offers {
		tx.Op(xdrbuild.DeleteOffer(o.OfferID))
		s.log.WithFields(logan.F{
			"offer":          o,
			"price_to_offer": priceToOffer,
		}).Info("Removing buy Offer.")
	}

	return true, nil
}

func (s *Service) removeSellOffersIfNecessary(ctx context.Context, assetPairConfig AssetPairConfig, priceToOffer int64, tx *xdrbuild.Transaction) (needNewOffer bool, err error) {

	f := false
	z := uint64(0)
	offers, err := s.accountInfoProvider.Offers(s.config.Source.Address(), assetPairConfig.BaseAsset, assetPairConfig.QuoteAsset, &f, "", &z)
	if err != nil {
		return false, errors.Wrap(err, "failed to obtain my Offers")
	}

	if len(offers) == 1 && int64(offers[0].Price) == priceToOffer && offers[0].BaseAmount >= assetPairConfig.BaseAssetVolume {
		// No need to refresh the Offer - it is full with actual price.
		s.log.WithFields(logan.F{
			"offer":          offers[0],
			"asset_pair":     assetPairConfig,
			"price_to_offer": priceToOffer,
		}).Debug("Sell Offer is not outdated, leaving it as is.")
		return false, nil
	}

	for _, o := range offers {
		tx.Op(xdrbuild.DeleteOffer(o.OfferID))
		s.log.WithField("offer", o).Info("Removing sell Offer.")
	}

	return true, nil
}

func (s *Service) getCurrentPrice(base, quote string) (*regources.Amount, error) {
	assetPairs, err := s.assetsGetter.Pairs()
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain AssetPairs from Horizon")
	}

	var assetPair *regources.AssetPair
	for _, aPair := range assetPairs {
		if aPair.Base == base && aPair.Quote == quote {
			assetPair = &aPair
			break
		}
	}

	if assetPair == nil {
		return nil, errors.From(errors.Errorf("No AssetPair in Horizon with BaseAsset (%s) and QuoteAsset (%s).", base, quote), logan.F{
			"horizon_asset_pairs": assetPairs,
		})
	}

	return &assetPair.CurrentPrice, nil
}
