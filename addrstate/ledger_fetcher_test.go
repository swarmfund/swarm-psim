package addrstate

import (
	"testing"

	"encoding/json"

	"time"

	"net/url"

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
		assert.Equal(t, "4294967296", ledger.ID)
		closedAt, _ := time.Parse(time.RFC3339, "1970-01-01T00:00:00Z")
		assert.Equal(t, closedAt, ledger.ClosedAt)
		assert.Equal(t, int64(1), ledger.TXCount)

		ledger = <-ch
		assert.Equal(t, "8589934592", ledger.ID)
		closedAt, _ = time.Parse(time.RFC3339, "2017-11-29T18:36:58Z")
		assert.Equal(t, closedAt, ledger.ClosedAt)
		assert.Equal(t, int64(2), ledger.TXCount)
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
	ledgersResponse = []byte(`{
    "_embedded": {
        "records": [
            {
                "_links": {
                    "operations": {
                        "href": "http://localhost:8000/ledgers/1/operations{?cursor,limit,order}",
                        "templated": true
                    },
                    "payments": {
                        "href": "http://localhost:8000/ledgers/1/payments{?cursor,limit,order}",
                        "templated": true
                    },
                    "self": {
                        "href": "http://localhost:8000/ledgers/1"
                    },
                    "transactions": {
                        "href": "http://localhost:8000/ledgers/1/transactions{?cursor,limit,order}",
                        "templated": true
                    }
                },
                "base_fee": 0,
                "base_reserve": "0.0000",
                "closed_at": "1970-01-01T00:00:00Z",
                "fee_pool": "0.0000",
                "hash": "97477fe2b40c4064f11389b7655fea12a791105b746dbc11356a6d31d96fadcb",
                "id": "97477fe2b40c4064f11389b7655fea12a791105b746dbc11356a6d31d96fadcb",
                "max_tx_set_size": 100,
                "operation_count": 0,
                "paging_token": "4294967296",
                "sequence": 1,
                "total_coins": "0.0000",
                "transaction_count": 1
            },
            {
                "_links": {
                    "operations": {
                        "href": "http://localhost:8000/ledgers/2/operations{?cursor,limit,order}",
                        "templated": true
                    },
                    "payments": {
                        "href": "http://localhost:8000/ledgers/2/payments{?cursor,limit,order}",
                        "templated": true
                    },
                    "self": {
                        "href": "http://localhost:8000/ledgers/2"
                    },
                    "transactions": {
                        "href": "http://localhost:8000/ledgers/2/transactions{?cursor,limit,order}",
                        "templated": true
                    }
                },
                "base_fee": 0,
                "base_reserve": "0.0000",
                "closed_at": "2017-11-29T18:36:58Z",
                "fee_pool": "0.0000",
                "hash": "ac43309efc3ecfe545226e44f9f1705b1e5c953751d00eaf650afbe2abba5635",
                "id": "ac43309efc3ecfe545226e44f9f1705b1e5c953751d00eaf650afbe2abba5635",
                "max_tx_set_size": 500,
                "operation_count": 0,
                "paging_token": "8589934592",
                "prev_hash": "97477fe2b40c4064f11389b7655fea12a791105b746dbc11356a6d31d96fadcb",
                "sequence": 2,
                "total_coins": "0.0000",
                "transaction_count": 2
            }
        ]
    },
    "_links": {
        "next": {
            "href": "http://localhost:8000/ledgers?order=asc&limit=2&cursor=8589934592"
        },
        "prev": {
            "href": "http://localhost:8000/ledgers?order=desc&limit=2&cursor=4294967296"
        },
        "self": {
            "href": "http://localhost:8000/ledgers?order=asc&limit=2&cursor="
        }
    }}`)
)
