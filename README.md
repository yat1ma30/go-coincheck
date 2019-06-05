# go-coincheck

Coincheck API

## Description

go-coincheck is a go client library for [Coincheck API](https://coincheck.com/languages/en).

Connecting websocket, read data[trads, orderbook]. as of 2019/06.

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

	go c.Connect()

	for {
		select {
		case v := <-c.Result:
			fmt.Printf("%+v\n", string(v))
		}
	}
}
```

## Author

[@_numbP](https://twitter.com/_numbP)

## License

[MIT](https://github.com/go-numb/go-coincheck/blob/master/LICENSE)