package sbankenSDK

import (
	"encoding/json"

	"github.com/larwef/sbankenSDK/common"
	"github.com/larwef/sbankenSDK/authentication"
	"github.com/larwef/sbankenSDK/client"
	"errors"
)

type accountsRepository struct {
	url    string
	client *client.SbankenClient
}

type accountsResponse struct {
	AvailableItems int       `json:"availableItems"`
	Items          []Account `json:"items"`
	common.Error
}

type accountResponse struct {
	Item Account `json:"item"`
	common.Error
}

type Account struct {
	AccountNumber   string  `json:"accountNumber"`
	CustomerId      string  `json:"customerId"`
	OwnerCustomerId string  `json:"ownerCustomerId"`
	Name            string  `json:"name"`
	AccountType     string  `json:"accountType"`
	Available       float64 `json:"available"`
	Balance         float64 `json:"balance"`
	CreditLimit     float64 `json:"creditLimit"`
	DefaultAccount  bool    `json:"defaultAccount"`
}

func NewAccountRepository(config Config) (*accountsRepository) {
	token := authentication.NewSbankenToken(config.IdentityServer, config.ClientId, config.ClientSecret)
	return &accountsRepository{url: config.AccountsEndpoint, client: client.NewSbankenClient(&token)}
}

func (ar accountsRepository) GetAccounts(customerId string) ([]Account, error) {
	var accountsRsp accountsResponse
	response, err := ar.client.Get(ar.url + customerId, nil)
	if err != nil {
		return accountsRsp.Items, err
	}

	err = json.Unmarshal(response, &accountsRsp)
	if accountsRsp.IsError == true {
		return accountsRsp.Items, errors.New(accountsRsp.ErrorMessage)
	}

	return accountsRsp.Items, err
}

func (ar accountsRepository) GetAccount(customerId string, accountNumber string) (Account, error) {
	var accountRsp accountResponse
	response, err := ar.client.Get(ar.url + customerId + "/" + accountNumber, nil)
	if err != nil {
		return accountRsp.Item, err
	}

	err = json.Unmarshal(response, &accountRsp)
	if accountRsp.IsError == true {
		return accountRsp.Item, errors.New(accountRsp.ErrorMessage)
	}

	return accountRsp.Item, err
}
