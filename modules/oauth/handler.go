package oauth

import (
	"fmt"
	"github.com/innovandalism/shodan/api"
	"github.com/innovandalism/shodan/bindata"
	"net/http"
	"errors"
)

type TokenTestResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

type ExchangeTokenRequest struct {
	Code string `json:"code"`
}

func handleAuthenticate(w http.ResponseWriter, _ *http.Request) {
	uri := fmt.Sprintf("https://discordapp.com/oauth2/authorize?client_id=%s&scope=identify&response_type=code&redirect_uri=%s", *mod.clientid, *mod.returnuri)
	api.Forward(w, uri)
}

func handleCallback(w http.ResponseWriter, _ *http.Request) {
	page, _ := bindata.Asset("assets/oauth/callback.html")
	fmt.Fprintf(w, "%s", page)
}

func handleExchangeToken(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		api.ErrorBadRequest(w, errors.New("Code not found"))
		return
	}
	err, tokenInfo := exchangeDiscordToken(code)
	if err != nil {
		api.ErrorInternalServerError(w, err)
		return
	}
	fmt.Fprint(w, tokenInfo.AccessToken)
}

func handleGetUserProfile(w http.ResponseWriter, r *http.Request) {
	var res api.ResponseEnvelope
	req, err := api.ReadRequest(r)
	if err != nil {
		api.ErrorInternalServerError(w, err)
	}
	if len(req.Token) < 1 {
		api.ErrorBadRequest(w, errors.New("Token is missing"))
	}
	u, err := api.FetchProfile(req.Token)

	if err != nil {
		api.ErrorUnauthorized(w, err)
	} else {
		res = api.ResponseEnvelope{
			Status: 200,
			Data: TokenTestResponse{
				Token: req.Token,
				User:  u,
			},
		}
	}

	api.SendResponse(w, &res)
}
