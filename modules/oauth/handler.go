package oauth

import (
	"fmt"
	"github.com/innovandalism/shodan"
	"net/http"
)

func handleAuthenticate(w http.ResponseWriter, _ *http.Request) {
	uri := fmt.Sprintf("https://discordapp.com/oauth2/authorize?client_id=%s&scope=identify guilds&response_type=code&redirect_uri=%s", mod.clientid, mod.returnuri)
	shodan.HttpForward(w, uri)
}

func handleExchangeToken(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		shodan.HttpSendError(w, shodan.ErrorHttp("code parameter missing", 400))
		return
	}
	err, tokenInfo := exchangeDiscordToken(code)
	if err != nil {
		shodan.HttpSendError(w, err)
		return
	}
	user, err := shodan.DUFetchProfile(tokenInfo.AccessToken)
	if err != nil {
		shodan.HttpSendError(w, err)
		return
	}
	id, err := DBAddToken(mod.shodan.GetDatabase(), tokenInfo.AccessToken, user.ID)
	if err != nil {
		shodan.HttpSendError(w, err)
		return
	}
	strToken, err := makeJwt(id)
	if err != nil {
		shodan.HttpSendError(w, err)
		return
	}
	shodan.HttpForward(w, fmt.Sprintf("/#!?jwt=%s", strToken))
}
