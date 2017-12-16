package addrstate

import (
	"encoding/json"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/distributed_lab/logan/v3"
)

type RequesterMock struct {
	mock.Mock
}

func (r *RequesterMock) Requester(method, endpoint string, target interface{}) error {
	args := r.Called(method, endpoint, target)
	return args.Error(0)
}

func TestNewLedgersProvider(t *testing.T) {
	requesterMock := RequesterMock{}
	t.Run("ledger unmarshal", func(t *testing.T) {
		ch := NewLedgersProvider(logan.New(), requesterMock.Requester)()
		requesterMock.
			On("Requester", "GET", mock.Anything, mock.Anything).
			Run(func(arguments mock.Arguments) {
				json.Unmarshal(ledgersResponse, arguments.Get(2))
			}).
			Return(nil).Once()
		defer requesterMock.AssertExpectations(t)

		ledger := <-ch
		assert.Equal(t, "103079215104", ledger.ID)
		closedAt, _ := time.Parse(time.RFC3339, "2017-12-13T10:23:33Z")
		assert.Equal(t, closedAt, ledger.ClosedAt)
		assert.Equal(t, int64(2), ledger.TXCount)

		ledger = <-ch
		assert.Equal(t, "107374182400", ledger.ID)
		closedAt, _ = time.Parse(time.RFC3339, "2017-12-13T10:23:38Z")
		assert.Equal(t, closedAt, ledger.ClosedAt)
		assert.Equal(t, int64(0), ledger.TXCount)
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
				if !assertFailed && !assert.Equal(t, "107374182400", u.Query().Get("cursor")) {
					assertFailed = true
				}
			}).
			Return(nil).Once()
		<-ch
		assert.False(t, assertFailed)
	})
}

var (
	ledgersResponse = []byte(`{
	"_embedded": {
        "records": [
            {
                "_links": {
                    "operations": {
                        "href": "http://localhost:8000/ledgers/24/operations{?cursor,limit,order}",
                        "templated": true
                    },
                    "payments": {
                        "href": "http://localhost:8000/ledgers/24/payments{?cursor,limit,order}",
                        "templated": true
                    },
                    "self": {
                        "href": "http://localhost:8000/ledgers/24"
                    },
                    "transactions": {
                        "href": "http://localhost:8000/ledgers/24/transactions{?cursor,limit,order}",
                        "templated": true
                    }
                },
                "base_fee": 0,
                "base_reserve": "0.0000",
                "closed_at": "2017-12-13T10:23:33Z",
                "fee_pool": "0.0000",
                "hash": "58ef1f6eca2e50c3b334a13184d200884ac99fb6c27e6019f0e918d01b36c029",
                "id": "58ef1f6eca2e50c3b334a13184d200884ac99fb6c27e6019f0e918d01b36c029",
                "max_tx_set_size": 500,
                "operation_count": 2,
                "paging_token": "103079215104",
                "prev_hash": "6e00c3c6ce5e90d7355f99e79ee2b9270a0287936bddc49a5323754b330f36f4",
                "sequence": 24,
                "total_coins": "0.0000",
                "transaction_count": 2
            },
            {
                "_links": {
                    "operations": {
                        "href": "http://localhost:8000/ledgers/25/operations{?cursor,limit,order}",
                        "templated": true
                    },
                    "payments": {
                        "href": "http://localhost:8000/ledgers/25/payments{?cursor,limit,order}",
                        "templated": true
                    },
                    "self": {
                        "href": "http://localhost:8000/ledgers/25"
                    },
                    "transactions": {
                        "href": "http://localhost:8000/ledgers/25/transactions{?cursor,limit,order}",
                        "templated": true
                    }
                },
                "base_fee": 0,
                "base_reserve": "0.0000",
                "closed_at": "2017-12-13T10:23:38Z",
                "fee_pool": "0.0000",
                "hash": "33edfd0c6c04f497d5be1bab64ffce596d8f6b77f9e00e16ca9aa080af775b77",
                "id": "33edfd0c6c04f497d5be1bab64ffce596d8f6b77f9e00e16ca9aa080af775b77",
                "max_tx_set_size": 500,
                "operation_count": 0,
                "paging_token": "107374182400",
                "prev_hash": "58ef1f6eca2e50c3b334a13184d200884ac99fb6c27e6019f0e918d01b36c029",
                "sequence": 25,
                "total_coins": "0.0000",
                "transaction_count": 0
            }
        ]
    },
    "_links": {
        "next": {
            "href": "http://localhost:8000/ledgers?order=asc&limit=2&cursor=107374182400"
        },
        "prev": {
            "href": "http://localhost:8000/ledgers?order=desc&limit=2&cursor=103079215104"
        },
        "self": {
            "href": "http://localhost:8000/ledgers?order=asc&limit=2&cursor=98784247808"
        }
    }}`)
)
