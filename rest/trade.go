package rest

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/buger/jsonparser"
)

type Trade struct {
	client *Client
}

type ResponseForPagination struct {
	Limit         int    `json:"limit"`
	Order         string `json:"order"`
	StartingAfter string `json:"starting_after,omitempty"`
	EndingBefore  string `json:"ending_before,omitempty"`
}

type ResponseForTrade struct {
	ID        int       `json:"id"`
	Amount    string    `json:"amount"`
	Rate      string    `json:"rate"`
	Pair      string    `json:"pair"`
	OrderType string    `json:"order_type"`
	CreatedAt time.Time `json:"created_at"`
}

// 最新の取引履歴を取得できます。
func (a Trade) Get(symbol string) ([]ResponseForTrade, error) {
	s := a.client.Request("GET", fmt.Sprintf("api/trades?pair=%s", symbol), "")
	isSuccess, err := jsonparser.GetBoolean([]byte(s), "success")
	if err != nil {
		return nil, err
	}
	if !isSuccess {
		return nil, nil
	}
	var trades []ResponseForTrade
	data, _, _, err := jsonparser.Get([]byte(s), "data")
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &trades); err != nil {
		return nil, err
	}

	return trades, nil
}
