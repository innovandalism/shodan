package oauth

import (
	"fmt"
	"github.com/innovandalism/shodan/api"
	"net/http"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/innovandalism/shodan/util"
	"github.com/bwmarrin/discordgo"
)

type ExchangeTokenRequest struct {
	Code string `json:"code"`
}

func handleAuthenticate(w http.ResponseWriter, _ *http.Request) {
	uri := fmt.Sprintf("https://discordapp.com/oauth2/authorize?client_id=%s&scope=identify guilds&response_type=code&redirect_uri=%s", *mod.clientid, *mod.returnuri)
	api.Forward(w, uri)
}

func handleExchangeToken(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		api.Error(w, 400, util.WrapError(errors.New("code parameter missing")))
		return
	}
	err, tokenInfo := exchangeDiscordToken(code)
	if err != nil {
		api.Error(w, 500, err)
		return
	}
	user, err := api.FetchProfile(tokenInfo.AccessToken)
	if err != nil {
		api.Error(w, 500, err)
		return
	}
	id, err := DBAddToken(mod.shodan.GetDatabase(), tokenInfo.AccessToken, user.ID)
	if err != nil {
		api.Error(w, 500, err)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": fmt.Sprint(id),
		"discord_uid": user.ID,
	})
	strToken, err := token.SignedString([]byte(*mod.jwtSecret))
	if err != nil {
		api.Error(w, 500, err)
		return
	}
	api.Forward(w, fmt.Sprintf("/#!?jwt=%s", strToken))
}

func handleGetUserProfile(w http.ResponseWriter, r *http.Request) {
	var res api.ResponseEnvelope
	req, err := api.ReadRequest(r)
	if err != nil {
		api.ErrorInfer(w, err)
		return
	}
	if len(req.Token) < 1 {
		api.Error(w, 401, util.WrapError(errors.New("Authorization required")))
		return
	}

	u, err := api.FetchProfile(req.Token)

	if err != nil {
		api.ErrorInfer(w, err)
		return
	}

	g, err := api.GetUserGuilds(req.Token)

	if err != nil {
		api.ErrorInfer(w, err)
		return
	}

	res = api.ResponseEnvelope{
		Status: 200,
		Data: struct{
			User *discordgo.User `json:"user"`
			Guilds []*discordgo.UserGuild `json:"guilds"`
		}{u, g},
	}

	api.SendResponse(w, &res)
}
