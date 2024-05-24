package rest

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/buger/jsonparser"
)

type Deposit struct {
	client *Client
}

type ResponseForDeposit struct {
	ID          int       `json:"id"`
	Amount      string    `json:"amount"`
	Currency    string    `json:"currency"`
	Address     string    `json:"address"`
	Status      string    `json:"status"`
	ConfirmedAt time.Time `json:"confirmed_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// You Get Deposit history
func (a Deposit) Get(param string) ([]ResponseForDeposit, error) {
	s := a.client.Request("GET", "api/deposit_money", param)
	isSuccess, err := jsonparser.GetBoolean([]byte(s), "success")
	if err != nil {
		return nil, err
	}
	if !isSuccess {
		return nil, fmt.Errorf("failed to get deposit")
	}

	deposits, _, _, err := jsonparser.Get([]byte(s), "deposits")
	if err != nil {
		return nil, err
	}

	var res []ResponseForDeposit
	if err := json.Unmarshal(deposits, &res); err != nil {
		return nil, err
	}

	return res, nil
}
