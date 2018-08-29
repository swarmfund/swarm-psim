package masternode

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"gitlab.com/swarmfund/psim/psim/app"
)

func init() {
	app.RegisterService("masternode_airdrop", func(ctx context.Context) (app.Service, error) {
		ts, err := time.Parse(time.RFC3339, "2018-08-29T15:40:33Z")
		if err != nil {
			panic(err)
		}
		fmt.Println(ts.String())

		ns := NodesState{
			eth:   app.Config(ctx).Ethereum(),
			token: common.HexToAddress("0x3bd98ca5189e51034abb48b48c2ccb55c43da23f"),
			rq: RewardQueue{
				promoteAfter: 1,
			},
			furnace:       common.HexToAddress("0x5eD3821880c01c2E63931eF6eF7d5d4Fe215610F"),
			firstBlock:    ts,
			blockDuration: 10 * time.Minute,
		}

		for {
			fmt.Println("trying next payout")
			start := time.Now()
			ns.Payout(ctx)
			end := time.Now()
			fmt.Println("TOOK", end.Sub(start).String())
		}
	})
}
