package rest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTicker(t *testing.T) {
	client := new(Client).NewClient("", "")
	res, err := client.Ticker.Get()
	assert.NoError(t, err)

	fmt.Printf("%+v", res)
}
func TestGetTrade(t *testing.T) {
	client := new(Client).NewClient("", "")
	res, err := client.Trade.Get("btc_jpy")
	assert.NoError(t, err)

	fmt.Printf("%+v", res)
}
func TestGetOrderbook(t *testing.T) {
	client := new(Client).NewClient("", "")
	res, err := client.OrderBook.Get("btc_jpy")
	assert.NoError(t, err)

	fmt.Printf("%+v", res)
}

func TestCreateOrder(t *testing.T) {
	client := new(Client).NewClient("", "")
	res, err := client.Order.Create(RequestForOrder{
		Pair:        "btc_jpy",
		OrderType:   "buy",
		Price:       "1000000",
		Amount:      "0.01",
		TimeInForce: "post_only",
	})
	assert.NoError(t, err)

	fmt.Printf("%+v", res)
}

func Test(t *testing.T) {
	client := new(Client).NewClient("", "")

	/** Private API */
	// 新規注文
	// "buy" 指値注文 現物取引 買い
	// "sell" 指値注文 現物取引 売り
	// "market_buy" 成行注文 現物取引 買い
	// "market_sell" 成行注文 現物取引 売り
	// "leverage_buy" 指値注文 レバレッジ取引新規 買い
	// "leverage_sell" 指値注文 レバレッジ取引新規 売り
	// "close_long" 指値注文 レバレッジ取引決済 売り
	// "close_short" 指値注文 レバレッジ取引決済 買い

	// 未決済の注文一覧
	client.Order.Opens()
	// 注文のキャンセル
	client.Order.Cancel("12345")
	// 取引履歴
	client.Order.Transactions()
	// ポジション一覧
	// client.leverage.Positions()
	// 残高
	client.Account.Balance()
	// // レバレッジアカウントの残高
	// client.Account.leverage_balance()
	// アカウント情報
	client.Account.Info()
	// ビットコインの送金
	client.Send.Create(`{"address":"","amount":"0.0002"`)
	// ビットコインの送金履歴
	client.Send.Get("currency=BTC")
	// ビットコインの受け取り履歴
	client.Deposit.Get("currency=BTC")
	// // ビットコインの高速入金
	// client.deposit.fast("12345")
	// 銀行口座一覧
	client.BankAccount.Get()
	// 銀行口座の登録
	client.BankAccount.Create(`{"bank_name":"MUFG","branch_name":"tokyo", "bank_account_type":"toza", "number":"1234567", "name":"Danny"}`)
	// 銀行口座の削除
	client.BankAccount.Delete("25621")
	// 出金履歴
	client.Withdraw.Get()
	// 出金申請の作成
	client.Withdraw.Create(`{"bank_account_id":"2222","amount":"50000", "currency":"JPY", "is_fast":"false"}`)
	// 出金申請のキャンセル
	client.Withdraw.Cancel("12345")
	// // 借入申請
	// client.Borrow.create(`{"amount":"100","currency":"JPY"}`)
	// // 借入中一覧
	// client.Borrow.matches()
	// // 返済
	// client.Borrow.repay("1135")
	// // レバレッジアカウントへの振替
	// client.Transfer.to_leverage(`{"amount":"100","currency":"JPY"}`)
	// // レバレッジアカウントからの振替
	// client.Transfer.from_leverage(`{"amount":"100","currency":"JPY"}`)
}
