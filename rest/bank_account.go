package rest

import (
	"encoding/json"
	"fmt"

	"github.com/buger/jsonparser"
)

type BankAccount struct {
	client *Client
}

type Bank struct {
	ID              int    `json:"id"`
	BankName        string `json:"bank_name"`
	BranchName      string `json:"branch_name"`
	BankAccountType string `json:"bank_account_type"`
	Number          string `json:"number"`
	Name            string `json:"name"`
}

// Create a new BankAccount.
func (a BankAccount) Create(param string) (*Bank, error) {
	s := a.client.Request("POST", "api/bank_accounts", param)
	isSuccess, err := jsonparser.GetBoolean([]byte(s), "success")
	if err != nil {
		return nil, err
	}
	if !isSuccess {
		return nil, fmt.Errorf("failed to create bank account")
	}

	data, _, _, err := jsonparser.Get([]byte(s), "data")
	if err != nil {
		return nil, err
	}

	var banks *Bank
	if err := json.Unmarshal(data, banks); err != nil {
		return nil, err
	}

	return banks, nil
}

// Get account information.
func (a BankAccount) Get() ([]Bank, error) {
	s := a.client.Request("GET", "api/bank_accounts", "")
	isSuccess, err := jsonparser.GetBoolean([]byte(s), "success")
	if err != nil {
		return nil, err
	}
	if !isSuccess {
		return nil, fmt.Errorf("failed to get bank account")
	}

	data, _, _, err := jsonparser.Get([]byte(s), "data")
	if err != nil {
		return nil, err
	}

	var banks []Bank
	if err := json.Unmarshal(data, &banks); err != nil {
		return nil, err
	}

	return banks, nil
}

// Delete a BankAccount.
func (a BankAccount) Delete(id string) error {
	s := a.client.Request("DELETE", "api/bank_accounts/"+id, "")
	isSuccess, err := jsonparser.GetBoolean([]byte(s), "success")
	if err != nil {
		return err
	}
	if !isSuccess {
		return fmt.Errorf("failed to delete bank account")
	}

	return nil
}
