package abx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	PricesEndpoint = "https://api.abx.com/prices"
)

type Connector struct {
}

func NewConnector() *Connector {
	return &Connector{}
}

func (c *Connector) Get() (*PricesResponse, error) {
	response, err := http.Get(PricesEndpoint)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get prices: %d", response.StatusCode)
	}
	var result *PricesResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	return result, err
}

type PricesResponse map[string]Hub

func (p PricesResponse) GetHub(id string) *Hub {
	value, ok := p[id]
	if !ok {
		return nil
	}
	return &value
}

type Hub map[string]AssetPrice

func (h Hub) GetTicker(ticker string) *AssetPrice {
	for _, v := range h {
		if v.Ticker == ticker {
			return &v
		}
	}
	return nil
}

type AssetPrice struct {
	Ticker string  `json:"ticker"`
	Price  float64 `json:"mid"`
}
