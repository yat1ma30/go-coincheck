package rest

type Withdraw struct {
	client *Client
}

// Transfer Balance to Leverage.
func (a Withdraw) Create(param string) string {
	return a.client.Request("POST", "api/withdraws", param)
}

// Transfer Balance from Leverage.
func (a Withdraw) Get() string {
	return a.client.Request("GET", "api/withdraws", "")
}

// Transfer Balance from Leverage.
func (a Withdraw) Cancel(id string) string {
	return a.client.Request("DELETE", "api/withdraws/"+id, "")
}
