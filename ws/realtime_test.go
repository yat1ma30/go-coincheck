package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

func TestConnect(t *testing.T) {
	c := NewRealtime()
	c.SetPing(30)
	defer c.Close()

	products := []string{WsBTCJPY}
	channels := []string{WsExecution}

	c.Subscribe(products, channels)

	ctx := context.Background()
	go c.Connect(ctx)

	for v := range c.Result {
		var s Trades
		if err := json.Unmarshal(v, &s); err != nil {
			t.Fatal(err)
		}

		fmt.Printf("%+v\n", s)
	}
}

func TestConnectForBoard(t *testing.T) {
	c := NewRealtime()
	c.SetPing(30)
	defer c.Close()

	products := []string{WsBTCJPY}
	channels := []string{WsOrderbook}

	c.Subscribe(products, channels)

	ctx := context.Background()
	go c.Connect(ctx)

	for v := range c.Result {
		var s Orderbook
		if err := json.Unmarshal(v, &s); err != nil {
			t.Fatal(err)
		}

		fmt.Printf("%+v\n", s)
	}
}
