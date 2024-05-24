package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Books struct {
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
	LastUpdateAt time.Time  `json:"last_update_at"`
}

type Orderbook struct {
	Pair  string `json:"pair"`
	Books Books  `json:"orderBook"`
}

func (o *Orderbook) UnmarshalJSON(data []byte) error {
	var rawOrderbook []interface{}
	if err := json.Unmarshal(data, &rawOrderbook); err != nil {
		return err
	}

	if len(rawOrderbook) != 2 {
		return errors.New("expected an array with 2 elements")
	}

	pair, ok := rawOrderbook[0].(string)
	if !ok {
		return errors.New("first element must be a string")
	}

	bookData, ok := rawOrderbook[1].(map[string]interface{})
	if !ok {
		return errors.New("second element must be a JSON object")
	}

	var books Books
	bids, ok := bookData["bids"].([]interface{})
	if !ok {
		return errors.New("bids must be an array of arrays")
	}
	for _, bid := range bids {
		bidArray, ok := bid.([]interface{})
		if !ok || len(bidArray) != 2 {
			return errors.New("each bid must be an array with 2 elements")
		}
		bidStr1, ok1 := bidArray[0].(string)
		bidStr2, ok2 := bidArray[1].(string)
		if !ok1 || !ok2 {
			return errors.New("both elements of a bid must be strings")
		}
		books.Bids = append(books.Bids, []string{bidStr1, bidStr2})
	}

	asks, ok := bookData["asks"].([]interface{})
	if !ok {
		return errors.New("asks must be an array of arrays")
	}
	for _, ask := range asks {
		askArray, ok := ask.([]interface{})
		if !ok || len(askArray) != 2 {
			return errors.New("each ask must be an array with 2 elements")
		}
		askStr1, ok1 := askArray[0].(string)
		askStr2, ok2 := askArray[1].(string)
		if !ok1 || !ok2 {
			return errors.New("both elements of an ask must be strings")
		}
		books.Asks = append(books.Asks, []string{askStr1, askStr2})
	}

	lastUpdateAt, ok := bookData["last_update_at"].(string)
	if !ok {
		return errors.New("last_update_at must be a string")
	}
	// unix string time to time.Time
	lastUpdateAtTime, err := UnixTimeStringToTime(lastUpdateAt)
	if err != nil {
		return err
	}
	books.LastUpdateAt = lastUpdateAtTime

	o.Pair = pair
	o.Books = books
	return nil
}

type Trades struct {
	Trades []Trade `json:"trades"`
}

type Trade struct {
	TradeID      string    `json:"trade_id"`
	Pair         string    `json:"pair"`
	Rate         string    `json:"rate"`
	Amount       string    `json:"amount"`
	OrderType    string    `json:"order_type"`
	TakerOrderID string    `json:"taker_order_id"`
	MakerOrderID string    `json:"maker_order_id"`
	Timestamp    time.Time `json:"timestamp"`
}

func (p *Trades) UnmarshalJSON(data []byte) error {
	// JSONデータをデコードする
	var rawData [][]string
	err := json.Unmarshal([]byte(data), &rawData)
	if err != nil {
		return err
	}

	// 二次元スライスを TradeData 構造体のスライスに変換
	trades := make([]Trade, len(rawData))
	for i, trade := range rawData {
		// unix string time to time.Time
		at, err := UnixTimeStringToTime(trade[0])
		if err != nil {
			return err
		}

		if len(trade) == 8 {
			trades[i] = Trade{
				Timestamp:    at,
				TradeID:      trade[1],
				Pair:         trade[2],
				Rate:         trade[3],
				Amount:       trade[4],
				OrderType:    trade[5],
				TakerOrderID: trade[6],
				MakerOrderID: trade[7],
			}
		}
	}

	p.Trades = trades

	return nil
}

func UnixTimeStringToTime(unixTimeStr string) (time.Time, error) {
	// Convert the string to an int64.
	unixTimeInt, err := strconv.ParseInt(unixTimeStr, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid Unix time string: %v", err)
	}

	// Convert the int64 to a time.Time object.
	timestamp := time.Unix(unixTimeInt, 0)
	return timestamp, nil
}
