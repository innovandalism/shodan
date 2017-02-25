package oauth

import (
	"net/url"
	"net/http"
	"strings"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"errors"
	"github.com/innovandalism/shodan/util"
)

type TokenInfo struct{
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope string `json:"scope"`
	Error string `json:"error"`
}

func exchangeDiscordToken(code string) (error, *TokenInfo) {
	uri := "https://discordapp.com/api/oauth2/token"

	postForm := url.Values{}
	postForm.Add("grant_type", "authorization_code")
	postForm.Add("client_id", *mod.clientid)
	postForm.Add("client_secret", *mod.clientsecret)
	postForm.Add("code", code)
	postForm.Add("redirect_uri", *mod.returnuri)

	txReq, err := http.NewRequest("POST", uri, strings.NewReader(postForm.Encode()))
	if err != nil {
		err = util.WrapError(err)
		return err, nil
	}

	txReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	txRes, err := client.Do(txReq)
	if err != nil {
		err = util.WrapError(err)
		return err, nil
	}
	txResContent, err := ioutil.ReadAll(txRes.Body)
	if err != nil {
		err = util.WrapError(err)
		return err, nil
	}
	tokenInfo := TokenInfo{}
	fmt.Printf("%s", txResContent)
	err = json.Unmarshal(txResContent, &tokenInfo)
	if err != nil {
		err = util.WrapError(err)
		return err, nil
	}
	if tokenInfo.Error != "" {
		return util.WrapError(errors.New(tokenInfo.Error)), nil
	}
	return nil, &tokenInfo
}