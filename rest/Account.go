package rest

import "encoding/json"

type Account struct {
	client *Client
}

type AccountBalance struct {
	Success      bool   `json:"success"`
	JPY          string `json:"jpy"`
	BTC          string `json:"btc"`
	JPYReserved  string `json:"jpy_reserved"`
	BTCReserved  string `json:"btc_reserved"`
	JPYLendInUse string `json:"jpy_lend_in_use"`
	BTCLendInUse string `json:"btc_lend_in_use"`
	JPYLent      string `json:"jpy_lent"`
	BTCLent      string `json:"btc_lent"`
	JPYDebt      string `json:"jpy_debt"`
	BTCDebt      string `json:"btc_debt"`
	JPYTsumitate string `json:"jpy_tsumitate"`
	BTCTsumitate string `json:"btc_tsumitate"`
}

// Make sure a balance.
func (a Account) Balance() (*AccountBalance, error) {
	s := a.client.Request("GET", "api/accounts/balance", "")
	var balance *AccountBalance
	if err := json.Unmarshal([]byte(s), balance); err != nil {
		return nil, err
	}

	return balance, nil
}

type ExchangeFees struct {
	MakerFeeRate string `json:"maker_fee_rate"`
	TakerFeeRate string `json:"taker_fee_rate"`
}

type AccountInfo struct {
	Success        bool                    `json:"success"`
	ID             int                     `json:"id"`
	Email          string                  `json:"email"`
	IdentityStatus string                  `json:"identity_status"`
	BitcoinAddress string                  `json:"bitcoin_address"`
	TakerFee       string                  `json:"taker_fee"`
	MakerFee       string                  `json:"maker_fee"`
	ExchangeFees   map[string]ExchangeFees `json:"exchange_fees"`
}

// Get account information.
func (a Account) Info() (*AccountInfo, error) {
	s := a.client.Request("GET", "api/accounts", "")
	var v *AccountInfo
	if err := json.Unmarshal([]byte(s), v); err != nil {
		return nil, err
	}

	return v, nil
}
