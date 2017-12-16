package addrstate

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/go/xdr"
)

func TestNewChangesProvider(t *testing.T) {
	requesterMock := RequesterMock{}
	t.Run("ledger unmarshal", func(t *testing.T) {
		provider := NewChangesProvider(logan.New(), requesterMock.Requester)

		requesterMock.
			On("Requester", "GET", mock.Anything, mock.Anything).
			Run(func(arguments mock.Arguments) {
				// if test fails check error here
				json.Unmarshal(ledgerTxResponse, arguments.Get(2))
			}).
			Return(nil).Once()
		defer requesterMock.AssertExpectations(t)

		ch := provider("4242")
		change := <-ch
		assert.Equal(t, xdr.LedgerEntryChangeTypeCreated, change.Type)
		assert.NotNil(t, change.Created)

		change = <-ch
		assert.Equal(t, xdr.LedgerEntryChangeTypeCreated, change.Type)
		assert.NotNil(t, change.Created)
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
                        "href": "http://localhost:8000/ledgers/24"
                    },
                    "operations": {
                        "href": "http://localhost:8000/transactions/5b675e7d0940a1d3915042cd431670d10e7918d602565da72df110e94f6349d9/operations{?cursor,limit,order}",
                        "templated": true
                    },
                    "precedes": {
                        "href": "http://localhost:8000/transactions?order=asc&cursor=103079219200"
                    },
                    "self": {
                        "href": "http://localhost:8000/transactions/5b675e7d0940a1d3915042cd431670d10e7918d602565da72df110e94f6349d9"
                    },
                    "succeeds": {
                        "href": "http://localhost:8000/transactions?order=desc&cursor=103079219200"
                    }
                },
                "created_at": "2017-12-13T12:23:33Z",
                "envelope_xdr": "AAAAAP4DpOIcoI8urCJITRZtEDS0wzyPuGojb7AbKpHcMR1gAAAAAAAAAAAAAAAAAAAAAAAAAABaOiwRAAAAAAAAAAEAAAAAAAAACwAAAAAAAAAAAAAAAAAAAANTVU4AAAAACFNVTiBuYW1lAAAAAP4DpOIcoI8urCJITRZtEDS0wzyPuGojb7AbKpHcMR1gAAAAC0Rlc2NyaXB0aW9uAAAAABNodHRwczovL215YXNzZXQuY29tAAAAAOjUpRAAAAAAAwAAAAZsb2dvSUQAAAAAAAAAAAAAAAAAAAAAAAHcMR1gAAAAQAlfR2ypEvVujn99pBmz8T6luPnHYRgspvdIWs3empv/uvcH11FjeOi+0zjoVV13CAdG6/RMslFlRRcWua4u2QI=",
                "fee_meta_xdr": "AAAAAA==",
                "fee_paid": 0,
                "hash": "5b675e7d0940a1d3915042cd431670d10e7918d602565da72df110e94f6349d9",
                "id": "5b675e7d0940a1d3915042cd431670d10e7918d602565da72df110e94f6349d9",
                "ledger": 24,
                "memo_type": "none",
                "operation_count": 1,
                "paging_token": "103079219200",
                "result_meta_xdr": "AAAAAAAAAAEAAAAEAAAAAAAAABgAAAAEAAAAAEelFRKMt4xwiQLENINoyxSGBr7mpulf8kWqPYECjzQXAAAAA1NVTgAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGAAAAAQAAAAAbm+vsiekcJtUI6lBMDuyuvXDTSsHSolpbK9G7KqrFJEAAAADU1VOAAAAAAD+A6TiHKCPLqwiSE0WbRA0tMM8j7hqI2+wGyqR3DEdYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAYAAAABAAAAAB+/F+cknrCW45Cu9Sk1XVoprQ2BqH6fdmIQuVgOdnOjwAAAANTVU4AAAAAAAIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABgAAAAGAAAAA1NVTgAAAAAA/gOk4hygjy6sIkhNFm0QNLTDPI+4aiNvsBsqkdwxHWAAAAAIU1VOIG5hbWUAAAAA/gOk4hygjy6sIkhNFm0QNLTDPI+4aiNvsBsqkdwxHWAAAAALRGVzY3JpcHRpb24AAAAAE2h0dHBzOi8vbXlhc3NldC5jb20AAAAA6NSlEAAAAAAAAAAAAAAAAAAAAAAAAAAAAwAAAAZsb2dvSUQAAAAAAAAAAAAA",
                "result_xdr": "AAAAAAAAAAAAAAAAAAAAAQAAAAAAAAALAAAAAAAAAAAAAAABAAAAAQAAAAAAAAAA",
                "signatures": [
                    "CV9HbKkS9W6Of32kGbPxPqW4+cdhGCym90hazd6am/+69wfXUWN46L7TOOhVXXcIB0br9EyyUWVFFxa5ri7ZAg=="
                ],
                "source_account": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
                "valid_after": "1970-01-01T00:00:00Z",
                "valid_before": "2017-12-20T09:23:30Z"
            },
            {
                "_links": {
                    "account": {
                        "href": "http://localhost:8000/accounts/GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636"
                    },
                    "ledger": {
                        "href": "http://localhost:8000/ledgers/24"
                    },
                    "operations": {
                        "href": "http://localhost:8000/transactions/1b8837d153886ecacbe0f5b53e3a8c1e93fb904fcd768de30d955107d7bef60f/operations{?cursor,limit,order}",
                        "templated": true
                    },
                    "precedes": {
                        "href": "http://localhost:8000/transactions?order=asc&cursor=103079223296"
                    },
                    "self": {
                        "href": "http://localhost:8000/transactions/1b8837d153886ecacbe0f5b53e3a8c1e93fb904fcd768de30d955107d7bef60f"
                    },
                    "succeeds": {
                        "href": "http://localhost:8000/transactions?order=desc&cursor=103079223296"
                    }
                },
                "created_at": "2017-12-13T12:23:33Z",
                "envelope_xdr": "AAAAAP4DpOIcoI8urCJITRZtEDS0wzyPuGojb7AbKpHcMR1gAAAAAAAAAAAAAAAAAAAAAAAAAABaOiwRAAAAAAAAAAEAAAAAAAAACwAAAAAAAAAAAAAAAAAAAANVU0QAAAAACFVTRCBuYW1lAAAAAP4DpOIcoI8urCJITRZtEDS0wzyPuGojb7AbKpHcMR1gAAAAC0Rlc2NyaXB0aW9uAAAAABNodHRwczovL215YXNzZXQuY29tAAAAAOjUpRAAAAAAAwAAAAZsb2dvSUQAAAAAAAAAAAAAAAAAAAAAAAHcMR1gAAAAQLSTWqDLLb1wpxtTHBra2+jVcUiQ5X8kpLNVhanpofvuaXMF8mvRGl9f9HT3ikcY/w0yGbC/MixaMSDGkn8lywA=",
                "fee_meta_xdr": "AAAAAA==",
                "fee_paid": 0,
                "hash": "1b8837d153886ecacbe0f5b53e3a8c1e93fb904fcd768de30d955107d7bef60f",
                "id": "1b8837d153886ecacbe0f5b53e3a8c1e93fb904fcd768de30d955107d7bef60f",
                "ledger": 24,
                "memo_type": "none",
                "operation_count": 1,
                "paging_token": "103079223296",
                "result_meta_xdr": "AAAAAAAAAAEAAAAEAAAAAAAAABgAAAAEAAAAAFV68CjIzDtqDwXg1NVIrb7Ut97/Fw8MkOPNGD/iNOA1AAAAA1VTRAAAAAAA/gOk4hygjy6sIkhNFm0QNLTDPI+4aiNvsBsqkdwxHWAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGAAAAAQAAAAAWFEKEJtOs3/HDCQ2cXBTaALXfEwOfLzv4g1aXbYM/XYAAAADVVNEAAAAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAYAAAABAAAAABxZa1+wM/nOtjZPflhKUbLtiWLqHjIQRg7ISXjj0MHXAAAAANVU0QAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABgAAAAGAAAAA1VTRAAAAAAA/gOk4hygjy6sIkhNFm0QNLTDPI+4aiNvsBsqkdwxHWAAAAAIVVNEIG5hbWUAAAAA/gOk4hygjy6sIkhNFm0QNLTDPI+4aiNvsBsqkdwxHWAAAAALRGVzY3JpcHRpb24AAAAAE2h0dHBzOi8vbXlhc3NldC5jb20AAAAA6NSlEAAAAAAAAAAAAAAAAAAAAAAAAAAAAwAAAAZsb2dvSUQAAAAAAAAAAAAA",
                "result_xdr": "AAAAAAAAAAAAAAAAAAAAAQAAAAAAAAALAAAAAAAAAAAAAAACAAAAAQAAAAAAAAAA",
                "signatures": [
                    "tJNaoMstvXCnG1McGtrb6NVxSJDlfySks1WFqemh++5pcwXya9EaX1/0dPeKRxj/DTIZsL8yLFoxIMaSfyXLAA=="
                ],
                "source_account": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
                "valid_after": "1970-01-01T00:00:00Z",
                "valid_before": "2017-12-20T09:23:30Z"
            }
        ]
    },
    "_links": {
        "next": {
            "href": "http://localhost:8000/ledgers/24/transactions?order=asc&limit=10&cursor=103079223296"
        },
        "prev": {
            "href": "http://localhost:8000/ledgers/24/transactions?order=desc&limit=10&cursor=103079219200"
        },
        "self": {
            "href": "http://localhost:8000/ledgers/24/transactions?order=asc&limit=10&cursor="
        }
    }
}`)
)
