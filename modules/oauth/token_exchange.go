package oauth

import (
	"encoding/json"
	"fmt"
	"github.com/innovandalism/shodan"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// TokenInfo represents the data returned by the OAuth 2 token exchange
type TokenInfo struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	Error        string `json:"error"`
}

func exchangeDiscordToken(code string) (error, *TokenInfo) {
	uri := "https://discordapp.com/api/oauth2/token"

	postForm := url.Values{}
	postForm.Add("grant_type", "authorization_code")
	postForm.Add("client_id", mod.clientid)
	postForm.Add("client_secret", mod.clientsecret)
	postForm.Add("code", code)
	postForm.Add("redirect_uri", mod.returnuri)

	txReq, err := http.NewRequest("POST", uri, strings.NewReader(postForm.Encode()))
	if err != nil {
		err = shodan.WrapError(err)
		return err, nil
	}

	txReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	txRes, err := client.Do(txReq)
	if err != nil {
		err = shodan.WrapError(err)
		return err, nil
	}
	txResContent, err := ioutil.ReadAll(txRes.Body)
	if err != nil {
		err = shodan.WrapError(err)
		return err, nil
	}
	tokenInfo := TokenInfo{}
	fmt.Printf("%s", txResContent)
	err = json.Unmarshal(txResContent, &tokenInfo)
	if err != nil {
		err = shodan.WrapError(err)
		return err, nil
	}
	if tokenInfo.Error != "" {
		return shodan.ErrorHttp(tokenInfo.Error, 403), nil
	}
	return nil, &tokenInfo
}
