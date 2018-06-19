package btcfunnel

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
)

func (s *Service) monitorLowBalance(ctx context.Context) {
	app.RunOverIncrementalTimer(ctx, s.log, "low_balance_monitor", s.checkBalance, 10*time.Second, 10*time.Second)
}

func (s *Service) checkBalance(_ context.Context) error {
	balance, err := s.getHotBalance()
	if err != nil {
		return errors.Wrap(err, "Failed to get Hot wallet balance")
	}

	if balance <= s.config.MinBalanceAlarmThreshold {
		_, err = s.sendLowBalanceNotification(balance)
		if err != nil {
			return errors.Wrap(err, "Failed to send Notification", logan.F{
				"balance": balance,
			})
		}
	}

	return nil
}

func (s *Service) sendLowBalanceNotification(currentBalance float64) (bool, error) {
	if time.Now().Sub(s.lastMinBalanceAlarmAt) < s.config.MinBalanceAlarmPeriod {
		return false, nil
	}

	message := fmt.Sprintf("Hey! There is not so much *BTC* left on the *%s* Hot wallet: *%.8f*.\n"+
		"I was asked to notify once Hot wallet balance is not grater than %.8f BTC.\n"+
		"I was asked to notify every *%s*.",
		s.config.OffchainBlockchain, currentBalance, s.config.MinBalanceAlarmThreshold, s.config.MinBalanceAlarmPeriod.String())

	err := s.notificationSender.Send(message)
	if err != nil {
		return false, errors.From(err, logan.F{
			"message_trying_to_send": message,
		})
	}

	s.log.WithFields(logan.F{
		"current_balance":             currentBalance,
		"min_balance_alarm_threshold": s.config.MinBalanceAlarmThreshold,
		"runner":                      "low_balance_monitor",
	}).Info("Notification about small balance on Hot wallet was sent.")

	s.lastMinBalanceAlarmAt = time.Now()

	return true, nil
}
