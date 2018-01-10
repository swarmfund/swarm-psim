package btcwithdraw

import (
	horizonV2 "gitlab.com/swarmfund/horizon-connector/v2"
)

func getRequestLoganFields(key string, request horizonV2.Request) map[string]interface{} {
	result := map[string]interface{}{
		key + "_id": request.ID,
		key + "_state": request.State,
	}

	if request.Details.Withdraw != nil {
		detKey := key + "_withdraw_details"

		result[detKey + "_amount"] = request.Details.Withdraw.Amount
		result[detKey + "_destination_amount"] = request.Details.Withdraw.DestinationAmount
		result[detKey + "_balance_id"] = request.Details.Withdraw.BalanceID
		result[detKey + "_external_details"] = request.Details.Withdraw.ExternalDetails
	}

	return result
}
