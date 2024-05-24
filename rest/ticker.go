package rest

import (
	"encoding/json"
)

type Ticker struct {
	client *Client
}

type ResponseForTicker struct {
	Last      float64 `json:"last"`
	Bid       float64 `json:"bid"`
	Ask       float64 `json:"ask"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Volume    float64 `json:"volume"`
	Timestamp int64   `json:"timestamp"`
}

// 各種最新情報を簡易に取得することができます。
func (a Ticker) Get() (*ResponseForTicker, error) {
	s := a.client.Request("GET", "api/ticker", "")
	t := new(ResponseForTicker)
	if err := json.Unmarshal([]byte(s), t); err != nil {
		return nil, err
	}

	return t, nil
}
