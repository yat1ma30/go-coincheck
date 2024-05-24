# go-coincheck

Coincheck API, connect websocket

## Description

go-coincheck is a go client library for [Coincheck API](https://coincheck.com/languages/en).

Connecting websocket, read data[trads, orderbook]. as of 2024/05.  
Use [official package](https://github.com/coincheckjp/coincheck-go) for REST.

## Installation

```
$ go get -u github.com/go-numb/go-coincheck
```

## Usage
``` 
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