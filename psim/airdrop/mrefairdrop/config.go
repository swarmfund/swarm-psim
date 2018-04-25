package mrefairdrop

import (
	"time"

	"gitlab.com/swarmfund/psim/psim/airdrop"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	IssuanceAsset         string    `fig:"issuance_asset,required"`
	SnapshotTime          time.Time `fig:"snapshot_time,required"`
	ApproveWaitFinishTime time.Time `fig:"approve_wait_finish_time,required"`

	Source keypair.Address `fig:"source,required"`
	Signer keypair.Full    `fig:"signer,required" mapstructure:"signer"`

	airdrop.EmailsConfig `fig:"emails,required"`

	BlackList []string `fig:"black_list"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"issuance_asset":           c.IssuanceAsset,
		"snapshot_time":            c.SnapshotTime,
		"approve_wait_finish_time": c.ApproveWaitFinishTime,
		"emails":                   c.EmailsConfig,
		"black_list_length":        len(c.BlackList),
	}
}
