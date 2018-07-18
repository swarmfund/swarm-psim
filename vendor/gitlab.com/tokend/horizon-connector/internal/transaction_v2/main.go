package transaction_v2

import (
	"gitlab.com/tokend/horizon-connector/internal"
	"gitlab.com/tokend/horizon-connector/internal/responses"
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/regources"
	"gitlab.com/tokend/horizon-connector/internal/resources"
	"net/url"
	"strconv"
	"fmt"
)

type Q struct {
	client internal.Client
}

func NewQ(client internal.Client) *Q {
	return &Q{
		client,
	}
}

func (q *Q) TransactionsByEffectsAndEntryTypes(effects, entryTypes []int,
) ([]regources.TransactionV2, *resources.PageMeta, error) {
	u := url.Values{}
	u.Add("limit", "1000")
	response, err := q.client.Get(fmt.Sprintf("/v2/transactions?%s%s%s", u.Encode(),
		getStringFromIntSlice("effect", effects), getStringFromIntSlice("entry_type", entryTypes)))
	if err != nil {
		return nil, nil, errors.Wrap(err, "transactions_v2 request failed")
	}

	var result responses.TransactionV2Index
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, nil, errors.Wrap(err, "failed to unmarshal transactions_v2")
	}

	return result.Embedded.Records, &result.Embedded.Meta, nil
}

func getStringFromIntSlice(fieldName string, input []int) string {
	u := url.Values{}
	for _, value := range input {
		u.Add(fieldName, strconv.Itoa(value))
	}

	return u.Encode()
}