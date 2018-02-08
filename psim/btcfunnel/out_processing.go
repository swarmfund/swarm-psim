package btcfunnel

import (
	"context"

	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
)

const (
	txTemplateSize = 20
	inSize         = 148
	outSize        = 21
)

var errNoScriptAddresses = errors.New("No Addresses in the ScriptPubKey of the UTXO.")

func (s Service) listenOutStream(ctx context.Context) {
	s.log.Info("Started listening to stream of our Outputs.")

	for {
		select {
		case <-ctx.Done():
			return
		case out, ok := <-s.outCh:
			if !ok {
				// No more Outs will come
				s.log.Info("Stopped listening to stream of our Outputs - channel has been closed.")

				app.RunUntilSuccess(ctx, s.log, "utxo_funnelling", s.funnelUTXOs, 5*time.Second)
				s.utxos = nil
				return
			}

			s.log.WithField("out", out).Debug("Processing our TX Output.")

			app.RunUntilSuccess(ctx, s.log, "checking_our_out", func(ctx context.Context) error {
				utxo, err := s.btcClient.GetTxUTXO(out.TXHash, out.Vout)
				if err != nil {
					return errors.Wrap(err, "Failed to Get TX UTXO")
				}

				if utxo != nil {
					// This Out is unspent.
					s.log.WithField("out", out).WithField("utxo", utxo).Info("Found our UTXO.")
					s.utxos = append(s.utxos, UTXO{
						UTXO: *utxo,
						Out:  out,
					})
				}

				return nil
			}, 5*time.Second)
		}
	}
}

func (s Service) funnelUTXOs(ctx context.Context) error {
	if len(s.utxos) == 0 {
		return nil
	}

	s.log.WithField("utxo_length", len(s.utxos)).Info("Started funnelling batch of our UTXOs.")

	var utxoOuts []bitcoin.Out
	var inputUTXOs []bitcoin.InputUTXO
	var totalInAmount float64
	var privateKeys []string

	for _, utxo := range s.utxos {
		utxoOuts = append(utxoOuts, utxo.Out)

		inputUTXOs = append(inputUTXOs, bitcoin.InputUTXO{
			Out:          utxo.Out,
			ScriptPubKey: utxo.ScriptPubKey.Hex,
			RedeemScript: nil,
		})

		totalInAmount += utxo.Value

		if len(utxo.ScriptPubKey.Addresses) == 0 {
			return errors.From(errNoScriptAddresses, logan.F{"utxo": utxo})
		}
		privateKeys = append(privateKeys, s.addrToPriv[utxo.ScriptPubKey.Addresses[0]])
	}

	fields := logan.F{
		"total_in_amount": totalInAmount,
	}

	hotBalance, err := s.getHotBalance()
	if err != nil {
		return errors.Wrap(err, "Failed to get hot balance", fields)
	}
	fields["hot_balance"] = hotBalance

	txFee, err := s.estimateTXFee(len(inputUTXOs), 2)
	if err != nil {
		return errors.Wrap(err, "Failed to estimate TX Fee", fields)
	}
	fields["tx_fee"] = txFee

	funnelOuts := s.countFunnelOuts(totalInAmount, hotBalance, txFee)
	fields["funnel_outs"] = funnelOuts

	txHash, err := s.craftAndSendTX(utxoOuts, funnelOuts, inputUTXOs, privateKeys)
	if err != nil {
		return errors.Wrap(err, "Failed to craft and send TX")
	}

	s.log.WithFields(fields).WithField("tx_hash", txHash).Info("Funneled BTC successfully.")
	return nil
}

// TODO
func (s Service) getHotBalance() (float64, error) {
	return 0, nil
}

// TODO Add maxPossibleFee to config and compare estimated fee with it
func (s *Service) estimateTXFee(inputs, outputs int) (float64, error) {
	feePerKB, err := s.btcClient.EstimateFee()
	if err != nil {
		return 0, errors.Wrap(err, "Failed to EstimateFee")
	}

	// TODO Add maxPossibleFee to config and compare estimated fee with it

	txSize := txTemplateSize + inSize*inputs + outSize*outputs

	return feePerKB * float64(txSize), nil
}

func (s *Service) countFunnelOuts(totalInAmount, hotBalance, txFee float64) map[string]float64 {
	amountToHot, amountToCold := s.countAmounts(totalInAmount-txFee, hotBalance)

	addrToAmount := make(map[string]float64)
	if amountToHot > 0 {
		addrToAmount[s.config.HotAddress] = amountToHot
	}
	if amountToCold > 0 {
		addrToAmount[s.config.ColdAddress] = amountToCold
	}

	return addrToAmount
}

// TODO Cover with tests.
func (s *Service) countAmounts(totalOutAmount, hotBalance float64) (amountToHot, amountToCold float64) {
	availableToHot := s.config.MaxHotStock - hotBalance

	if availableToHot < s.config.DustOutputLimit {
		// Already enough BTC on the Hot Address - sending everything to the Cold.
		return 0, totalOutAmount
	}

	if totalOutAmount-availableToHot > s.config.DustOutputLimit {
		// After fulfilling of the Hot, still have more than Dust to be sent to the Cold.
		return availableToHot, totalOutAmount - availableToHot
	}

	// Whole money to the Hot.
	return totalOutAmount, 0
}

func (s *Service) craftAndSendTX(utxoOuts []bitcoin.Out, funnelOuts map[string]float64, inputUTXOs []bitcoin.InputUTXO,
	privateKeys []string) (txHash string, err error) {

	unsignedTX, err := s.btcClient.CreateRawTX(utxoOuts, funnelOuts)
	if err != nil {
		return "", errors.Wrap(err, "Failed to CreateRawTX")
	}

	signedTX, err := s.btcClient.SignRawTX(unsignedTX, inputUTXOs, privateKeys)
	if err != nil {
		return "", errors.Wrap(err, "Failed to SignRawTX", logan.F{"unsigned_tx": unsignedTX})
	}

	txHash, err = s.btcClient.SendRawTX(signedTX)
	if err != nil {
		return "", errors.Wrap(err, "Failed to SendRawTX", logan.F{"signed_tx": signedTX})
	}

	return txHash, nil
}
