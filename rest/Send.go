package rest

type Send struct {
	client *Client
}

// Sending Bitcoin to specified Bitcoin addres.
func (a Send) Create(param string) string {
	return a.client.Request("POST", "api/send_money", param)
}

// You Get Send history
func (a Send) Get(param string) string {
	return a.client.Request("GET", "api/send_money", param)
}
