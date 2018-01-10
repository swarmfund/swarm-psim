package account

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/mocks"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/resources"
)

func TestQ_Balances(t *testing.T) {
	client := mocks.Client{}
	q := NewQ(&client)

	t.Run("account not found", func(t *testing.T) {
		client.On("Get", mock.Anything).Return(nil, nil).Once()
		defer client.AssertExpectations(t)

		got, err := q.Balances("foobar")
		assert.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("existing account", func(t *testing.T) {
		data := []byte(`[{
        	"account_id": "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
        	"asset": "SUN",
        	"balance_id": "BBXG7L5SE6SHBG2UEOUUCMB3WK5PLQ2NFMDUVCLJNSXUN3FKVMKJCXRU"
    	}]`)
		expected := []resources.Balance{
			{
				Asset:     "SUN",
				AccountID: "GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636",
				BalanceID: "BBXG7L5SE6SHBG2UEOUUCMB3WK5PLQ2NFMDUVCLJNSXUN3FKVMKJCXRU",
			},
		}
		client.On("Get", mock.Anything).Return(data, nil).Once()
		defer client.AssertExpectations(t)

		got, err := q.Balances("foobar")
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})
}
