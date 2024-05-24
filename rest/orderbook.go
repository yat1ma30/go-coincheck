package rest

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type OrderBook struct {
	client *Client
}

type Book struct {
	Price  float64
	Amount float64
}

type ResponseForOrderBook struct {
	Asks []Book
	Bids []Book
}

// 板情報を取得できます。
func (a OrderBook) Get(symbol string) (*ResponseForOrderBook, error) {
	s := a.client.Request("GET", fmt.Sprintf("api/order_books?pair=%s", symbol), "")
	var res ResponseForOrderBook
	if err := json.Unmarshal([]byte(s), &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (o *ResponseForOrderBook) UnmarshalJSON(data []byte) error {
	var raw map[string][][]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for _, rawAsk := range raw["asks"] {
		price, ok := rawAsk[0].(string)
		if !ok {
			price = "0"
		}
		fprice, err := strconv.ParseFloat(price, 64)
		if err != nil {
			fprice = 0
		}
		amount, ok := rawAsk[1].(string)
		if !ok {
			amount = "0"
		}
		fsize, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			fsize = 0
		}

		o.Asks = append(o.Asks, Book{Price: fprice, Amount: fsize})
	}

	for _, rawBid := range raw["bids"] {
		price, ok := rawBid[0].(string)
		if !ok {
			price = "0"
		}
		fprice, err := strconv.ParseFloat(price, 64)
		if err != nil {
			fprice = 0
		}
		amount, ok := rawBid[1].(string)
		if !ok {
			amount = "0"
		}
		fsize, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			fsize = 0
		}

		o.Bids = append(o.Bids, Book{Price: fprice, Amount: fsize})
	}

	return nil
}
