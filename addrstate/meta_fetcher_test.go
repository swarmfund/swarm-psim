package addrstate

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/distributed_lab/logan/v3"
)

func TestNewChangesProvider(t *testing.T) {
	requesterMock := RequesterMock{}
	t.Run("ledger unmarshal", func(t *testing.T) {
		provider := NewChangesProvider(logan.New(), requesterMock.Requester)
		requesterMock.
			On("Requester", "GET", mock.Anything, mock.Anything).
			Run(func(arguments mock.Arguments) {
				err := json.Unmarshal(ledgerTxResponse, arguments.Get(2))
				fmt.Println(err)
			}).
			Return(nil).Once()
		defer requesterMock.AssertExpectations(t)

		ch := provider("4242")
		change := <-ch
		fmt.Println(change)
		//change := <-ch
		//change = change
		//assert.Equal(t, "4294967296", ledger.ID)
		//closedAt, _ := time.Parse(time.RFC3339, "1970-01-01T00:00:00Z")
		//assert.Equal(t, closedAt, ledger.ClosedAt)
		//assert.Equal(t, int64(1), ledger.TXCount)
		//
		//ledger = <-ch
		//assert.Equal(t, "8589934592", ledger.ID)
		//closedAt, _ = time.Parse(time.RFC3339, "2017-11-29T18:36:58Z")
		//assert.Equal(t, closedAt, ledger.ClosedAt)
		//assert.Equal(t, int64(2), ledger.TXCount)
	})

	t.Run("cursor handling", func(t *testing.T) {
		ch := NewLedgersProvider(logan.New(), requesterMock.Requester)()
		requesterMock.
			On("Requester", "GET", mock.Anything, mock.Anything).
			Run(func(arguments mock.Arguments) {
				json.Unmarshal(ledgersResponse, arguments.Get(2))
			}).
			Return(nil).Once()
		<-ch
		<-ch
		assertFailed := false
		requesterMock.
			On("Requester", "GET", mock.Anything, mock.Anything).
			Run(func(arguments mock.Arguments) {
				json.Unmarshal(ledgersResponse, arguments.Get(2))
				u, err := url.Parse(arguments.String(1))
				if !assert.NoError(t, err) {
					assertFailed = true
				}
				if !assertFailed && !assert.Equal(t, "8589934592", u.Query().Get("cursor")) {
					assertFailed = true
				}
			}).
			Return(nil).Once()
		<-ch
		assert.False(t, assertFailed)
	})
}

var (
	ledgerTxResponse = []byte(`{
    "_embedded": {
        "records": [
            {
                "_links": {
                    "account": {
                        "href": "http://localhost:8000/accounts/GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636"
                    },
                    "ledger": {
                        "href": "http://localhost:8000/ledgers/36875"
                    },
                    "operations": {
                        "href": "http://localhost:8000/transactions/4301906343c3f2e9cdaa77bddac638413a891362de2163129ea33a3e058bb1fe/operations{?cursor,limit,order}",
                        "templated": true
                    },
                    "precedes": {
                        "href": "http://localhost:8000/transactions?order=asc&cursor=158376919044096"
                    },
                    "self": {
                        "href": "http://localhost:8000/transactions/4301906343c3f2e9cdaa77bddac638413a891362de2163129ea33a3e058bb1fe"
                    },
                    "succeeds": {
                        "href": "http://localhost:8000/transactions?order=desc&cursor=158376919044096"
                    }
                },
                "created_at": "2017-12-06T13:07:06Z",
                "envelope_xdr": "AAAAAP4DpOIcoI8urCJITRZtEDS0wzyPuGojb7AbKpHcMR1gAAAAAAAAAAAAAAAAAAAAAAAAAABaMPuMAAAAAAAAAAEAAAAAAAAAAAAAAACEl4fwgtpGCT/QoJ44YMlKa5yT7R7jSHaetaNeD2Td7AAAAAAAAAAFAAAAAAAAAAAAAAAAAAAAAdwxHWAAAABAzJ20dSKwq8AvTQNqiyvlh8fdqZA56I97tzs1cX9OX0yeBQq4HmKrUt0+wI+2aaclnd0uvAnjNS7JZWL6RcVnCg==",
                "fee_meta_xdr": "AAAAAA==",
                "fee_paid": 0,
                "hash": "4301906343c3f2e9cdaa77bddac638413a891362de2163129ea33a3e058bb1fe",
                "id": "4301906343c3f2e9cdaa77bddac638413a891362de2163129ea33a3e058bb1fe",
                "ledger": 36875,
                "memo_type": "none",
                "operation_count": 1,
                "paging_token": "158376919044096",
                "result_meta_xdr": "AAAAAAAAAAEAAAACAAAAAwAAi8sAAAAAAAAAAISXh/CC2kYJP9CgnjhgyUprnJPtHuNIdp61o14PZN3sAQAAAAAAAAAAAAAAAAAAAAAAAAUAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAACQCwAAAAAAAAAAhJeH8ILaRgk/0KCeOGDJSmuck+0e40h2nrWjXg9k3ewBAAAAAAAAAAAAAAAAAAAAAAAABQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==",
                "result_xdr": "AAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
                "signatures": [
                    "zJ20dSKwq8AvTQNqiyvlh8fdqZA56I97tzs1cX9OX0yeBQq4HmKrUt0+wI+2aaclnd0uvAnjNS7JZWL6RcVnCg=="
                ],
                "source_account": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
                "valid_after": "1970-01-01T00:00:00Z",
                "valid_before": "2017-12-13T10:06:05Z"
            }
        ]
    },
    "_links": {
        "next": {
            "href": "http://localhost:8000/ledgers/36875/transactions?order=asc&limit=10&cursor=158376919044096"
        },
        "prev": {
            "href": "http://localhost:8000/ledgers/36875/transactions?order=desc&limit=10&cursor=158376919044096"
        },
        "self": {
            "href": "http://localhost:8000/ledgers/36875/transactions?order=asc&limit=10&cursor="
        }
    }}`)
)
