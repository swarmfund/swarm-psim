package internal

import (
	"context"

	"sync"

	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/create_account_streamer"
)

// CreatedAccountsStreamer describes the interface an CreatedAccountOps streamer
// must implement to be used in this AddressQ.
type CreatedAccountsStreamer interface {
	Run(ctx context.Context)
	ReadinessWaiter() <-chan struct{}
	GetStream() <-chan create_account_streamer.CreateAccountOp
}

// AddressQ listens to CreateAccountOps from CreatedAccountsStreamer
// and saves BTC addresses of Accounts into watchAddresses map.
//
// Listening and watchAddresses filling happens in the blocking method Run().
//
// Allows to GetAccountID by BTC Address.
type AddressQ struct {
	ctx                     context.Context
	createdAccountsStreamer CreatedAccountsStreamer
	watchAddresses          map[string]string
}

// NewAddressQ constructs AddressQ from CreatedAccountsStreamer.
func NewAddressQ(ctx context.Context, createdAccountsStreamer CreatedAccountsStreamer) *AddressQ {
	q := &AddressQ{
		ctx: ctx,
		createdAccountsStreamer: createdAccountsStreamer,
	}

	return q
}

// FIXME
// FIXME
// FIXME
// If this btc Address - "" must be returned.
func (q *AddressQ) GetAccountID(btcAddress string) string {
	// FIXME
	// FIXME
	// FIXME
	//return q.watchAddresses[btcAddress]

	// FIXME
	// FIXME
	// FIXME
	if btcAddress == "mv6xAvLk88ZyoqxfEXfNbwqdcMPXnJuBww" {
		return "some_random_stellar_public_key"
	}

	return ""
}

// ReadinessWaiter returns channel, if reading from it doesn't block i.e. channel is closed -
// then AddressQ is ready.
//
// Don't use AddressQ if it's not ready.
// NOTE: AddressQ can become *not* ready after it has already become ready.
//
// Will never return nil channel.
func (q *AddressQ) ReadinessWaiter() <-chan struct{} {
	return q.createdAccountsStreamer.ReadinessWaiter()
}

// Run is processing CreateAccountOps streamed by createdAccountsStreamer.
func (q *AddressQ) Run() {
	streamer := q.createdAccountsStreamer.GetStream()

	wg := &sync.WaitGroup{}
	app.StartRunners(q.ctx, wg, q.createdAccountsStreamer)

	for {
		select {
		case <-q.ctx.Done():
			return
		case createAccountOp := <-streamer:
			// TODO Save BTC address from CreateAccountOp into map
			createAccountOp = createAccountOp
			//q.watchAddresses[createAccountOp.BTCAddress] = createAccountOp.AccountID
		}
	}

	wg.Wait()
}
