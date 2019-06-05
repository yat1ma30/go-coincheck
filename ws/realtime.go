package ws

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Realtime struct {
	conn   *websocket.Conn
	Result chan []byte
}

const (
	WSENDPOINT = "wss://ws-api.coincheck.com/"
)

/*
	WebSocket APIはHTTP通信では実装が難しかった、リアルタイムでのやり取りが容易になり、取引所で公開されている取引の履歴、板情報を取得することができます。

	WebSocket APIはβ版であるため、仕様に変更があったり、動作が不安定になる可能性があります。また、利用可能な通貨ペアは"btc_jpy"のみになります。(2017/03/03時点)
*/
func NewRealtime() *Realtime {
	conn, _, err := websocket.DefaultDialer.Dial(WSENDPOINT, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &Realtime{
		conn:   conn,
		Result: make(chan []byte, 512),
	}
}

func (p *Realtime) Close() {
	p.conn.Close()
}

type JsonRPC2 struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
}

const (
	// 現在btc_jpyのみ
	WsBTCJPY = "btc_jpy"

	WsExecution = "trades"
	WsOrderbook = "orderbook"
)

func (p *Realtime) Subscribe(products, channels []string) {
	for _, channel := range channels {
		for _, product := range products {
			readCh := fmt.Sprintf("%s-%s", product, channel)
			fmt.Printf("%+v\n", readCh)
			if err := p.conn.WriteJSON(
				&JsonRPC2{
					Type:    "subscribe",
					Channel: readCh,
				}); err != nil {
				log.Fatal("subscribe:", err)
			}

			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (p *Realtime) Connect() {

	p.read()
}

type WsTrade []interface{}

func (p *Realtime) read() {
	for {
		_, msg, err := p.conn.ReadMessage()
		if err != nil {
			continue
		}
		p.Result <- msg
	}
}
