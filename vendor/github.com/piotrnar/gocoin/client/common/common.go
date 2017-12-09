package common

import (
	"fmt"
	"sync"
	"time"
	"errors"
	"github.com/piotrnar/gocoin/lib/btc"
	"github.com/piotrnar/gocoin/lib/chain"
	"github.com/piotrnar/gocoin/lib/others/sys"
	"github.com/piotrnar/gocoin/lib/others/utils"
)

const (
	ConfigFile = "gocoin.conf"
	Version = uint32(70015)
	Services = uint64(0x00000009)
)

var (
	BlockChain *chain.Chain
	GenesisBlock *btc.Uint256
	Magic [4]byte
	Testnet bool

	Last TheLastBlock

	GocoinHomeDir string
	StartTime time.Time
	MaxPeersNeeded int

	DefaultTcpPort uint16

	MaxExpireTime time.Duration
	ExpirePerKB time.Duration

	DebugLevel int64

	CounterMutex sync.Mutex
	Counter map[string] uint64 = make(map[string]uint64)

	BusyWith string
	Busy_mutex sync.Mutex

	NetworkClosed sys.SyncBool

	AverageBlockSize sys.SyncInt

	AllBalMinVal uint64

	DropSlowestEvery time.Duration
	BlockExpireEvery time.Duration
	PingPeerEvery time.Duration

	UserAgent string
)


type TheLastBlock struct {
	sync.Mutex // use it for writing and reading from non-chain thread
	Block *chain.BlockTreeNode
	time.Time
}


func (b *TheLastBlock) BlockHeight() (res uint32) {
	b.Mutex.Lock()
	res = b.Block.Height
	b.Mutex.Unlock()
	return
}


func CountSafe(k string) {
	CounterMutex.Lock()
	Counter[k]++
	CounterMutex.Unlock()
}

func CountSafeAdd(k string, val uint64) {
	CounterMutex.Lock()
	Counter[k] += val
	CounterMutex.Unlock()
}


func Busy(b string) {
	Busy_mutex.Lock()
	BusyWith = b
	Busy_mutex.Unlock()
}


func BytesToString(val uint64) string {
	if val < 1e6 {
		return fmt.Sprintf("%.1f KB", float64(val)/1e3)
	} else if val < 1e9 {
		return fmt.Sprintf("%.2f MB", float64(val)/1e6)
	}
	return fmt.Sprintf("%.2f GB", float64(val)/1e9)
}


func NumberToString(num float64) string {
	if num>1e24 {
		return fmt.Sprintf("%.2f Y", num/1e24)
	}
	if num>1e21 {
		return fmt.Sprintf("%.2f Z", num/1e21)
	}
	if num>1e18 {
		return fmt.Sprintf("%.2f E", num/1e18)
	}
	if num>1e15 {
		return fmt.Sprintf("%.2f P", num/1e15)
	}
	if num>1e12 {
		return fmt.Sprintf("%.2f T", num/1e12)
	}
	if num>1e9 {
		return fmt.Sprintf("%.2f G", num/1e9)
	}
	if num>1e6 {
		return fmt.Sprintf("%.2f M", num/1e6)
	}
	if num>1e3 {
		return fmt.Sprintf("%.2f K", num/1e3)
	}
	return fmt.Sprintf("%.2f", num)
}


func HashrateToString(hr float64) string {
	return NumberToString(hr)+"H/s"
}

// Calculates average blocks size over the last "CFG.Stat.BSizeBlks" blocks
// Only call from blockchain thread.
func RecalcAverageBlockSize() {
	n := BlockChain.BlockTreeEnd
	var sum, cnt uint
	for maxcnt := CFG.Stat.BSizeBlks; maxcnt>0 && n!=nil; maxcnt-- {
		sum += uint(n.BlockSize)
		cnt++
		n = n.Parent
	}
	if sum>0 && cnt>0 {
		//println("The average block size is", sum/cnt, "at block height", BlockChain.BlockTreeEnd.Height)
		AverageBlockSize.Store(int(sum/cnt))
	} else {
		AverageBlockSize.Store(204)
	}
}

func GetRawTx(BlockHeight uint32, txid *btc.Uint256) (data []byte, er error) {
	data, er = BlockChain.GetRawTx(BlockHeight, txid)
	if er!=nil {
		data = utils.GetTxFromWeb(txid)
		if data!=nil {
			er = nil
		} else {
			er = errors.New("GetRawTx and GetTxFromWeb failed for " + txid.String())
		}
	}
	return
}
