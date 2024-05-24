package rest

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	accessKey string
	secretKey string
	Account
	BankAccount
	// borrow       Borrow
	Deposit
	// leverage   Leverage
	Order
	OrderBook
	Send
	Ticker
	Trade
	// transfer     Transfer
	Withdraw
}

func (c *Client) NewClient(accessKey string, secretKey string) *Client {
	c.accessKey = accessKey
	c.secretKey = secretKey
	c.Account = Account{c}
	c.BankAccount = BankAccount{c}
	// c.borrow = Borrow{c}
	c.Deposit = Deposit{c}
	// c.leverage = Leverage{c}
	c.Order = Order{c}
	c.OrderBook = OrderBook{c}
	c.Send = Send{c}
	c.Ticker = Ticker{c}
	c.Trade = Trade{c}
	// c.transfer = Transfer{c}
	c.Withdraw = Withdraw{c}
	return c
}

type Response struct {
	Success bool  `json:"success"`
	Data    []any `json:"data"`
}

func (c *Client) Request(method string, path string, param string) string {
	if param != "" && method == "GET" {
		path = path + "?" + param
		param = ""
	}
	url := "https://coincheck.com/" + path
	nonce := strconv.FormatInt(CreateNonce(), 10)
	message := nonce + url + param
	req := &http.Request{}
	if method == "POST" {
		payload := strings.NewReader(param)
		req, _ = http.NewRequest(method, url, payload)
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	signature := ComputeHmac256(message, c.secretKey)
	req.Header.Add("ACCESS-KEY", c.accessKey)
	req.Header.Add("ACCESS-NONCE", nonce)
	req.Header.Add("ACCESS-SIGNATURE", signature)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Sprintf("Error: %s", err.Error())
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Sprintf("Error: %s", err.Error())
	}

	return string(body)
}

// create nonce by milliseconds
func CreateNonce() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

// create signature
func ComputeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
