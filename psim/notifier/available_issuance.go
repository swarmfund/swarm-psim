package notifier

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/psim/psim/notifier/internal/emails"
	"gitlab.com/tokend/psim/psim/notifier/internal/types"
	"gitlab.com/tokend/psim/psim/utils"
)

func (s *Service) checkAssetsIssuanceAmount(ctx context.Context) {
	if s.Assets == nil {
		s.logger.Warn("assets issuance checker is not enabled")
		return
	}

	d, err := time.ParseDuration(s.Assets.CheckPeriod)
	if err != nil {
		s.errors <- errors.Wrap(err, "can't start asset loader")
		return
	}

	ticker := time.NewTicker(d)
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("finished background ticker")
			ticker.Stop()
			return
		case <-ticker.C:
			err = s.loadAssets()
			if err != nil {
				s.errors <- errors.Wrap(err, "load assets runner failed")
			}
		}
	}
}

func (s *Service) loadAssets() (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = errors.FromPanic(rec)
			s.logger.WithStack(err).WithError(err).Error("load assets runner recovered")
		}
	}()

	assets, err := s.getAssetsList()
	if err != nil {
		return errors.Wrap(err, "unable to get assets list")
	}

	for _, asset := range assets {
		err = s.processAsset(asset)
		if err != nil {
			return errors.Wrap(err, "asset processing failed", logan.F{"asset": asset.Code})
		}
	}
	return nil
}

func (s *Service) processAsset(asset types.Asset) error {
	if !contains(s.Assets.Codes, asset.Code) {
		return nil
	}
	if asset.AvailableForIssuance > s.Assets.EmissionThreshold {
		return nil
	}

	err := s.notifyOwner(asset)
	if err != nil {
		return errors.Wrap(err, "Failed to send asset notice letter")
	}
	return nil
}

func (s *Service) notifyOwner(asset types.Asset) error {
	letter := &emails.NoticeLetter{
		ID:       utils.GenerateToken(),
		Header:   fmt.Sprintf("%s Admin Notification", s.ProjectName),
		Email:    s.Assets.NotificationReceiver,
		Template: emails.NoticeTemplateLowIssuance,
		Message: fmt.Sprintf(
			"Asset %s has low emission. Upload more presigned emissions.",
			asset.Code),
	}

	return s.sendEmail(letter, s.Assets.PayloadID)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
