package deployer

import (
	"context"
	"time"

	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/keypair"
)

type EntityCountGetter func(systemType string) (uint64, error)

type Deriver interface {
	ChildAddress(uint64) (string, error)
}

type Opts struct {
	Log           *logan.Entry
	ExternalTypes []string
	EntityCount   EntityCountGetter
	TargetCount   uint64
	Deriver       Deriver
	TXBuilder     *xdrbuild.Builder
	Source        keypair.Address
	Signer        keypair.Full
	Horizon       *horizon.Connector
}

func ExternalAccountDeployer(opts Opts) func(context.Context) error {
	// TODO validate opts
	return func(ctx context.Context) (err error) {
		ctx, cancel := context.WithCancel(ctx)
		defer func() {
			if rvr := recover(); rvr != nil {
				// we might spend actual money here,
				// so in case of emergency abandon the operations completely
				cancel()
				err = errors.Wrap(errors.FromPanic(rvr), "service panicked")
			}
		}()
		for _, systemType := range opts.ExternalTypes {
			current, err := opts.EntityCount(systemType)
			if err != nil {
				return errors.Wrap(err, "failed to get current entity count")
			}

			for current <= opts.TargetCount {
				if running.IsCancelled(ctx) {
					return nil
				}
				fields := logan.F{}
				address, err := opts.Deriver.ChildAddress(current)
				if err != nil {
					return errors.Wrap(err, "failed to derive external address")
				}
				fields["external_address"] = address
				opts.Log.WithFields(fields).Info("external address derived")
				// critical section. external address has been derived, we need to create entity at any cost
				running.UntilSuccess(context.Background(), opts.Log, "create-pool-entity", func(i context.Context) (bool, error) {
					tx := opts.TXBuilder.Transaction(opts.Source)
					for _, systemType := range opts.ExternalTypes {
						tx = tx.Op(xdrbuild.CreateExternalPoolEntry(cast.ToInt32(systemType), data, s.deployerID))
					}
					tx = tx.Sign(opts.Signer)
					envelope, err := tx.Marshal()
					if err != nil {
						return false, errors.Wrap(err, "failed to marshal tx")
					}

					result := opts.Horizon.Submitter().Submit(context.TODO(), envelope)
					if result.Err != nil {
						return false, errors.Wrap(result.Err, "failed to submit tx", logan.F{
							"tx_result": result.GetLoganFields(),
						})
					}
					return true, nil
				}, 1*time.Second, 1*time.Minute)

				opts.Log.WithFields(fields).Info("entities created")

				current += 1
			}
		}
		opts.Log.Info("all good")
		return nil
	}
}
