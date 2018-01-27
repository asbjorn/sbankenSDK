package authentication

import (
	"net/http"
	"strings"
	"log"
	"os"
	"encoding/base64"
	"io/ioutil"
	"encoding/json"
	"time"
)

type sbankenToken struct {
	tokenBase
	identityServer string
	clientId       string
	clientSecret   string
}

type sbankenTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func NewSbankenToken(identityServer string, clientId string, clientSecret string) (sbankenToken) {
	response := callIdentityServer(identityServer, clientId, clientSecret)

	token := sbankenToken{
		tokenBase: tokenBase{
			tokenString: response.AccessToken,
			validTo:     time.Now().Add(time.Duration(response.ExpiresIn) * time.Second),
			tokenType:   response.TokenType,
		},
		identityServer: identityServer,
		clientId:       clientId,
		clientSecret:   clientSecret,
	}

	return token
}

func (sbt *sbankenToken) GetTokenString() (string) {
	return sbt.tokenString
}

func (sbt *sbankenToken) GetExpirationTime() (time.Time) {
	return sbt.validTo
}

func (sbt *sbankenToken) GetTokenType() (string) {
	return sbt.tokenType
}

func (sbt *sbankenToken) RefreshToken() {
	response := callIdentityServer(sbt.identityServer, sbt.clientId, sbt.clientSecret)

	sbt.tokenString = response.AccessToken
	sbt.validTo = time.Now().Add(time.Duration(response.ExpiresIn) * time.Second)
	sbt.tokenType = response.TokenType
}

func callIdentityServer(identityServer string, clientId string, clientSecret string) (sbankenTokenResponse) {
	client := &http.Client{}
	request, err := http.NewRequest("POST", identityServer, strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(clientId+":"+clientSecret)))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	response, err := client.Do(request)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	var bankToken sbankenTokenResponse
	json.Unmarshal(body, &bankToken)

	return bankToken
}
