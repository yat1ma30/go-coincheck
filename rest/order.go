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
	Pair         string    `json:"pair"`
	Rate         string    `json:"rate"`
	Amount       string    `json:"amount"`
	OrderType    string    `json:"order_type"`
	TimeInForce  string    `json:"time_in_force"`
	StopLossRate string    `json:"stop_loss_rate,omitempty"`
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

type ResponseForOpenOrder struct {
	ID                     int       `json:"id"`
	Pair                   string    `json:"pair"`
	Rate                   float64   `json:"rate"`           // nullの場合があるためポインタ型を使用
	PendingAmount          string    `json:"pending_amount"` // nullの場合があるためポインタ型を使用
	OrderType              string    `json:"order_type"`
	PendingMarketBuyAmount string    `json:"pending_market_buy_amount"` // nullの場合があるためポインタ型を使用
	StopLossRate           string    `json:"stop_loss_rate"`            // nullの場合があるためポインタ型を使用
	CreatedAt              time.Time `json:"created_at"`
}

// List charges filtered by params
func (a Order) Opens() ([]ResponseForOrder, error) {
	s := a.client.Request("GET", "api/exchange/orders/opens", "")
	isSuccess, err := jsonparser.GetBoolean([]byte(s), "success")
	if err != nil {
		return nil, err
	}
	if !isSuccess {
		errs, _ := jsonparser.GetString([]byte(s), "error")
		return nil, fmt.Errorf("failed to get open orders, error: %s", errs)
	}

	data, _, _, err := jsonparser.Get([]byte(s), "orders")
	if err != nil {
		return nil, err
	}

	var res []ResponseForOrder
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res, nil
}

type ResponseForCancelStatus struct {
	ID        int       `json:"id"`
	Cancel    bool      `json:"cancel"`
	CreatedAt time.Time `json:"created_at"`
}

// CancelStatus returns the status of the order cancellation.
func (a Order) CancelStatus(id string) (*ResponseForCancelStatus, error) {
	s := a.client.Request("GET", fmt.Sprintf("api/exchange/orders/cancel_status?id=%s", id), "")
	isSuccess, err := jsonparser.GetBoolean([]byte(s), "success")
	if err != nil {
		return nil, err
	}
	if !isSuccess {
		errs, _ := jsonparser.GetString([]byte(s), "error")
		return nil, fmt.Errorf("failed to get open orders, error: %s", errs)
	}

	var res *ResponseForCancelStatus
	if err := json.Unmarshal([]byte(s), res); err != nil {
		return nil, err
	}

	return res, nil
}

type ResponseForTransaction struct {
	ID      int `json:"id"`
	OrderID int `json:"order_id"`

	Pair string `json:"pair"`
	Side string `json:"side"`
	Rate string `json:"rate"`

	Fee         string `json:"fee"`
	FeeCurrency string `json:"fee_currency"`

	Liquidity string `json:"liquidity"`

	Funds     map[string]interface{} `json:"funds"`
	CreatedAt time.Time              `json:"created_at"`
}

// Get Order Transactions
func (a Order) Transactions() ([]ResponseForTransaction, error) {
	s := a.client.Request("GET", "api/exchange/orders/transactions", "")
	isSuccess, err := jsonparser.GetBoolean([]byte(s), "success")
	if err != nil {
		return nil, err
	}
	if !isSuccess {
		errs, _ := jsonparser.GetString([]byte(s), "error")
		return nil, fmt.Errorf("failed to get open orders, error: %s", errs)
	}

	data, _, _, err := jsonparser.Get([]byte(s), "transactions")
	if err != nil {
		return nil, err
	}

	var res []ResponseForTransaction
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res, nil
}
