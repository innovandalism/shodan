package oauth

import (
	"fmt"
	Discord "github.com/bwmarrin/discordgo"
	"github.com/innovandalism/shodan/api"
	"github.com/innovandalism/shodan/bindata"
	"net/http"
)

type TokenTestResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

func handleAuthenticate(w http.ResponseWriter, _ *http.Request) {
	uri := fmt.Sprintf("https://discordapp.com/oauth2/authorize?client_id=%s&scope=identify&response_type=token&redirect_uri=%s", *mod.clientid, *mod.returnuri)
	w.Header().Add("Location", uri)
	w.WriteHeader(301)
}

func handleCallback(w http.ResponseWriter, _ *http.Request) {
	page, _ := bindata.Asset("assets/oauth/callback.html")
	fmt.Fprintf(w, "%s", page)
}

func fetchProfile(token string) (error, *Discord.User) {
	t := fmt.Sprintf("Bearer %s", token)
	s, err := Discord.New(t)
	defer s.Close()
	if err != nil {
		return err, nil
	}
	u, err := s.User("@me")
	if err != nil {
		return err, nil
	}
	return nil, u
}

func handleTokenTest(w http.ResponseWriter, r *http.Request) {
	req := api.ReadRequest(r)
	if len(req.Token) < 1 {
		res := api.ResponseEnvelope{
			Status: 500,
			Error:  "Token not found",
		}
		api.SendResponse(w, &res)
		return
	}
	err, u := fetchProfile(req.Token)

	var res api.ResponseEnvelope
	if err != nil {
		res = api.ResponseEnvelope{
			Status: 403,
			Error:  fmt.Sprintf("%s", err),
			Data: TokenTestResponse{
				Token: req.Token,
			},
		}
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
