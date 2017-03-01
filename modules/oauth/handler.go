package oauth

import (
	"fmt"
	"github.com/innovandalism/shodan/api"
	"net/http"
	"errors"
	"github.com/dgrijalva/jwt-go"
)

type ExchangeTokenRequest struct {
	Code string `json:"code"`
}

func handleAuthenticate(w http.ResponseWriter, _ *http.Request) {
	uri := fmt.Sprintf("https://discordapp.com/oauth2/authorize?client_id=%s&scope=identify&response_type=code&redirect_uri=%s", *mod.clientid, *mod.returnuri)
	api.Forward(w, uri)
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
	user, err := api.FetchProfile(tokenInfo.AccessToken)
	if err != nil {
		api.ErrorInternalServerError(w, err)
		return
	}
	id, err := mod.shodan.GetPostgres().AddToken(tokenInfo.AccessToken, user.ID)
	if err != nil {
		api.ErrorInternalServerError(w, err)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": fmt.Sprint(id),
		"discord_uid": user.ID,
	})
	strToken, err := token.SignedString([]byte(*mod.jwtSecret))
	if err != nil {
		api.ErrorInternalServerError(w, err)
		return
	}
	api.Forward(w, fmt.Sprintf("/#!?jwt=%s", strToken))
}

func handleGetUserProfile(w http.ResponseWriter, r *http.Request) {
	var res api.ResponseEnvelope
	req, err := api.ReadRequest(r)
	if err != nil {
		api.ErrorInternalServerError(w, err)
		return
	}
	if len(req.Token) < 1 {
		api.ErrorBadRequest(w, errors.New("Token is missing"))
		return
	}
	u, err := api.FetchProfile(req.Token)

	if err != nil {
		api.ErrorUnauthorized(w, err)
		return
	} else {
		res = api.ResponseEnvelope{
			Status: 200,
			Data: u,
		}
	}
	api.SendResponse(w, &res)
}
