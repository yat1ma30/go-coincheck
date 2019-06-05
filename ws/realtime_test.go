package ws

import (
	"fmt"
	"testing"
)

func TestConnect(t *testing.T) {
	c := NewRealtime()
	defer c.Close()

	products := []string{WsBTCJPY}
	channels := []string{WsExecution, WsOrderbook}

	c.Subscribe(products, channels)

	go c.Connect()

	for {
		select {
		case v := <-c.Result:
			fmt.Printf("%+v\n", string(v))
		}
	}
}
