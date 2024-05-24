package rest

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/buger/jsonparser"
)

type Order struct {
	client *Client
}

type RequestForOrder struct {
	Pair            string `json:"pair"`                        // 取引ペア（例: "btc_jpy"）
	OrderType       string `json:"order_type"`                  // 注文方法（"buy", "sell", "market_buy", "market_sell"）
	Price           string `json:"rate,omitempty"`              // 注文のレート（例: "28000"）、指値注文の場合に必要
	Amount          string `json:"amount,omitempty"`            // 注文での量（例: "0.1"）、指値注文および成行売り注文の場合に必要
	MarketBuyAmount string `json:"market_buy_amount,omitempty"` // 成行買で利用する日本円の金額（例: "10000"）、成行買い注文の場合に必要
	StopLossRate    string `json:"stop_loss_rate,omitempty"`    // 逆指値レート、オプショナル
	TimeInForce     string `json:"time_in_force,omitempty"`     // 注文有効期間（"good_til_cancelled" あるいは "post_only"）、オプショナル
}

type ResponseForOrder struct {
	ID           int       `json:"id"`
	Rate         string    `json:"rate"`
	Amount       string    `json:"amount"`
	OrderType    string    `json:"order_type"`
	TimeInForce  string    `json:"time_in_force"`
	StopLossRate string    `json:"stop_loss_rate,omitempty"`
	Pair         string    `json:"pair"`
	CreatedAt    time.Time `json:"created_at"`
}

// Create a order object with given parameters.In live mode, this issues a transaction.
func (a Order) Create(op RequestForOrder) (*ResponseForOrder, error) {
	param, err := json.Marshal(op)
	if err != nil {
		return nil, err
	}

	s := a.client.Request("POST", "api/exchange/orders", string(param))
	isSuccess, err := jsonparser.GetBoolean([]byte(s), "success")
	if err != nil {
		return nil, err
	}
	if !isSuccess {
		errs, _ := jsonparser.GetString([]byte(s), "error")
		return nil, fmt.Errorf("failed to create order, error: %s", errs)
	}

	var res *ResponseForOrder
	if err := json.Unmarshal([]byte(s), res); err != nil {
		return nil, err
	}

	return res, nil

}

// cancel a created order specified by order id. Optional argument amount is to refund partially.
func (a Order) Cancel(id string) string {
	s := a.client.Request("DELETE", "api/exchange/orders/"+id, "")
	isSuccess, err := jsonparser.GetBoolean([]byte(s), "success")
	if err != nil {
		return err.Error()
	}
	if !isSuccess {
		return "failed to cancel order"
	}

	result, err := jsonparser.GetInt([]byte(s), "id")
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%d", result)
}

// List charges filtered by params
func (a Order) Opens() string {
	return a.client.Request("GET", "api/exchange/orders/opens", "")
}

// Get Order Transactions
func (a Order) Transactions() string {
	return a.client.Request("GET", "api/exchange/orders/transactions", "")
}
