package sbankenSDK

import (
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
)

type Config struct {
	ClientId             string `json:"clientId"`
	ClientSecret         string `json:"clientSecret"`
	IdentityServer       string `json:"identityServer"`
	AccountsEndpoint     string `json:"accountsEndpoint"`
	TransactionsEndpoint string `json:"transactionsEndpoint"`
}

func ConfigFromFile(filepath string) (Config) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	var config Config
	json.Unmarshal(file, &config)

	return config
}

