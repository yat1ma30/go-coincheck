# go-coincheck

Coincheck API, connect websocket.  
add part of READ API.

## Description

go-coincheck is a go client library for [Coincheck API](https://coincheck.com/ja/documents/exchange/api).

Connecting websocket, read data[trads, orderbook]. as of 2024/05.  
Use [official package](https://github.com/coincheckjp/coincheck-go) for REST.

## Installation

```
$ go get -u github.com/go-numb/go-coincheck
```

## Usage
```go


func main() {
	apiKey := ""
	apiSecret := ""
	client := new(Client).NewClient(apiKey, apiSecret)
	res, err := client.Ticker.Get()
	assert.NoError(t, err)
	// &{Last:1.0796181e+07 Bid:1.0796108e+07 Ask:1.0798757e+07 High:1.0862513e+07 Low:1.0492623e+07 Volume:863.84966125 Timestamp:1716589070}

	res, err := client.Order.Create(RequestForOrder{
		Pair:        "btc_jpy",
		OrderType:   "buy",
		Price:       "1000000",
		Amount:      "0.01",
		TimeInForce: "post_only",
	})
	assert.NoError(t, err)

	// 注文のキャンセル
	res, err := client.Order.Cancel("12345")
	assert.NoError(t, err)
	// id:12345
}
```


```go
package main

import (
 "fmt"
 "github.com/go-numb/go-coincheck/ws"
)


func main() {
	c := ws.NewRealtime()
	defer c.Close()

	products := []string{ws.WsBTCJPY}
	channels := []string{ws.WsExecution, ws.WsOrderbook}

	c.Subscribe(products, channels)

	ctx := context.Background()
	go c.Connect(ctx)

	for {
		select {
		case v := <-c.Result:
			fmt.Printf("%+v\n", string(v))

			// OR
			var s Orderbook
			if err := json.Unmarshal(v, &s); err != nil {
				t.Fatal(err)
			}

			fmt.Printf("%+v\n", s)
		}
	}
}
```

## Author

[@_numbP](https://twitter.com/_numbP)

## License

[MIT](https://github.com/go-numb/go-coincheck/blob/master/LICENSE)