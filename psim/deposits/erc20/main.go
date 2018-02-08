package erc20

import (
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/go/hash"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/psim/addrstate"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/deposit"
	"gitlab.com/swarmfund/psim/psim/deposits/erc20/internal"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	Source        keypair.Address
	Signer        keypair.Full
	Cursor        uint64
	Confirmations uint64
	BaseAsset     string
	DepositAsset  string
}

func init() {
	app.RegisterService(conf.ServiceERC20Deposit, func(ctx context.Context) (app.Service, error) {
		config := Config{}

		err := figure.
			Out(&config).
			With(figure.BaseHooks, utils.ETHHooks).
			From(app.Config(ctx).Get(conf.ServiceERC20Deposit)).
			Please()
		if err != nil {
			return nil, errors.Wrap(err, "failed to figure out")
		}

		horizon := app.Config(ctx).Horizon().WithSigner(config.Signer)

		addrProvider := addrstate.New(
			ctx,
			app.Log(ctx),
			internal.StateMutator(config.BaseAsset, config.DepositAsset),
			horizon.Listener(),
		)

		info, err := horizon.Info()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get horizon info")
		}
		builder := xdrbuild.NewBuilder(info.Passphrase, info.TXExpirationPeriod)

		ethclient := app.Config(ctx).Ethereum()

		return deposit.New(
			app.Log(ctx),
			config.Source,
			config.Signer,
			conf.ServiceERC20Deposit,
			// FIXME implement verifier
			"",
			config.Cursor,
			config.Confirmations,
			app.Config(ctx).Horizon().WithSigner(config.Signer),
			addrProvider,
			app.Config(ctx).Discovery(),
			builder,
			NewERC20Helper(ethclient, config),
		), nil
	})
}

func NewERC20Helper(eth *ethclient.Client, config Config) *ERC20Helper {
	return &ERC20Helper{
		NewConfigHelper(config.DepositAsset, 0, 0),
		NewConverter(),
		NewReferenceBuilder(),
		NewETHHelper(eth),
	}
}

type ERC20Helper struct {
	*ConfigHelper
	*Converter
	*ReferenceBuilder
	*ETHHelper
}

type ConfigHelper struct {
	depositAsset string
	minDeposit   uint64
	depositFee   uint64
}

func NewConfigHelper(depositAsset string, minDeposit, depositFee uint64) *ConfigHelper {
	return &ConfigHelper{
		depositAsset,
		minDeposit,
		depositFee,
	}
}

func (h *ConfigHelper) GetAsset() string {
	return h.depositAsset
}

func (h *ConfigHelper) GetMinDepositAmount() uint64 {
	return h.minDeposit
}

func (h *ConfigHelper) GetFixedDepositFee() uint64 {
	return h.depositFee
}

type Converter struct {
}

func NewConverter() *Converter {
	return &Converter{}
}

func (h *Converter) ConvertToSystem(offchain uint64) uint64 {
	// expected offchain to be in gwei precision (10^9)
	var gwei int64 = 1000000000
	if offchain > math.MaxInt64 {
		panic("overflow")
	}
	result, overflow := amount.BigDivide(amount.One, int64(offchain), gwei, amount.ROUND_DOWN)
	if overflow {
		panic("overflow")
	}
	return uint64(result)
}

type ReferenceBuilder struct {
}

func NewReferenceBuilder() *ReferenceBuilder {
	return &ReferenceBuilder{}
}

func (h *ReferenceBuilder) BuildReference(blockNumber uint64, txHash, offchainAddress string, outIndex uint, maxLen int) string {
	base := fmt.Sprintf("%d:%s:%s:%d", blockNumber, txHash, offchainAddress, outIndex)
	hash := hash.Hash([]byte(base))
	return hex.EncodeToString(hash[:])
}

type ETHHelper struct {
	eth *ethclient.Client
}

func NewETHHelper(eth *ethclient.Client) *ETHHelper {
	return &ETHHelper{eth}
}

func (h *ETHHelper) GetLastKnownBlockNumber() (uint64, error) {
	head, err := h.eth.HeaderByNumber(context.TODO(), nil)
	if err != nil {
		return 0, err
	}
	return head.Number.Uint64(), nil
}

func (h *ETHHelper) GetBlock(number uint64) (*deposit.Block, error) {
	logs, err := h.eth.FilterLogs(context.TODO(), ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(number),
		ToBlock:   new(big.Int).SetUint64(number),
		Addresses: []common.Address{
			common.HexToAddress("0xd159bded8d0c82ab6a6610d7e26235a535d3f64e"),
		},
		Topics: [][]common.Hash{
			{common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")},
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get block logs")
	}

	if len(logs) == 0 {
		// no interesting outputs
		return nil, nil
	}

	ethBlock, err := h.eth.BlockByNumber(context.TODO(), new(big.Int).SetUint64(number))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get block")
	}

	block := deposit.Block{
		Hash:      ethBlock.Hash().Hex(),
		Timestamp: time.Unix(ethBlock.Time().Int64(), 0),
	}

	transactions := map[common.Hash][]deposit.Out{}

	for _, log := range logs {
		if len(log.Topics) != 3 {
			// TODO log invalid log
			continue
		}
		// third indexed topic is 20 bytes receiver address packed in 40 bytes, big-endian layout
		receiver := common.BytesToAddress(log.Topics[2][len(log.Topics[2])-20:])
		amount := new(big.Int).SetBytes(log.Data)
		gweiAmount := toGwei(amount)
		if !gweiAmount.IsUint64() {
			panic("overflow")
		}
		transactions[log.TxHash] = append(transactions[log.TxHash], deposit.Out{
			receiver.Hex(),
			gweiAmount.Uint64(),
		})
		fmt.Println("FOUND", receiver.Hex())
	}

	for hash, outputs := range transactions {
		block.TXs = append(block.TXs, deposit.Tx{
			Hash: hash.Hex(),
			Outs: outputs,
		})
	}

	return &block, nil
}

func toGwei(amount *big.Int) *big.Int {
	return new(big.Int).Div(amount, new(big.Int).SetInt64(1000000000))
}
