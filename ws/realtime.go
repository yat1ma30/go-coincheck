package ws

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Realtime struct {
	conn    *websocket.Conn
	pingSec int
	Result  chan []byte
}

const (
	WSENDPOINT = "wss://ws-api.coincheck.com/"
)

/*
WebSocket APIはHTTP通信では実装が難しかった、リアルタイムでのやり取りが容易になり、取引所で公開されている取引の履歴、板情報を取得することができます。

WebSocket APIはβ版であるため、仕様に変更があったり、動作が不安定になる可能性があります。また、利用可能な通貨ペアは"btc_jpy"のみになります。(2017/03/03時点)

Update:
利用可能な通貨ペアはbtc_jpy, etc_jpy, lsk_jpy, mona_jpy, plt_jpy, fnct_jpy, dai_jpy, wbtc_jpyになります。(2020/09時点)
*/
func NewRealtime() *Realtime {
	conn, _, err := websocket.DefaultDialer.Dial(WSENDPOINT, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &Realtime{
		conn:    conn,
		pingSec: 0,
		Result:  make(chan []byte, 512),
	}
}

func (p *Realtime) SetPing(sec int) {
	p.pingSec = sec
}

func (p *Realtime) Close() {
	p.conn.Close()
}

type JsonRPC2 struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
}

const (
	WsBTCJPY  = "btc_jpy"
	WsETCJPY  = "etc_jpy"
	WsLSKJPY  = "lsk_jpy"
	WsMONAJPY = "mona_jpy"
	WsPLTJPY  = "plt_jpy"
	WsFNCTJPY = "fnct_jpy"
	WsDAIJPY  = "dai_jpy"
	WsWBTCJPY = "wbtc_jpy"

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

func (p *Realtime) Connect(ctx context.Context) error {
	child, cancel := context.WithCancel(ctx)
	defer cancel()

	if p.pingSec > 0 {
		go p.ping(child)
	}

	if err := p.read(child); err != nil {
		return err
	}

	return nil
}

func (p *Realtime) ping(ctx context.Context) error {
	defer p.conn.Close()

	var (
		writeWait = 10 * time.Second
		ticker    = time.NewTicker(time.Duration(p.pingSec) * time.Second)
	)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return err
			}

		case <-ctx.Done():
			return fmt.Errorf("write context close: %w", ctx.Err())
		}
	}
}

func (p *Realtime) read(ctx context.Context) error {
	defer p.conn.Close()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("read context close: %w", ctx.Err())

		default:
			p.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			_, msg, err := p.conn.ReadMessage()
			if err != nil {
				continue
			}
			p.Result <- msg
		}

	}
}
