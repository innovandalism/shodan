package webui

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/mux"
	"github.com/innovandalism/shodan"
	"net/http"
)

func handleGetUserProfile(w http.ResponseWriter, r *http.Request) {
	var res shodan.ResponseEnvelope
	req, err := shodan.ReadRequest(r)
	if err != nil {
		shodan.HttpSendError(w, err)
		return
	}
	if !req.Authenticated() {
		shodan.HttpSendError(w, shodan.ErrorHttp("authorization required", 401))
		return
	}

	u, err := shodan.DUFetchProfile(req.Token)

	if err != nil {
		shodan.HttpSendError(w, err)
		return
	}

	// use Session.UserGuilds if this acts up
	botGuilds := copyGuildState(mod.shodan.GetDiscord().State)
	userGuilds, err := shodan.DUGetUserGuilds(req.Token)
	if err != nil {
		shodan.HttpSendError(w, err)
		return
	}
	// build two arrays:
	// commonGuilds are the guilds both the bot and the user are joined to
	// ownedGuilds are the guilds the user has owner flags for
	var commonGuilds, ownedGuilds []*discordgo.UserGuild
	for _, userGuild := range userGuilds {
		if userGuild.Owner {
			ownedGuilds = append(ownedGuilds, userGuild)
		}
		for _, botGuild := range botGuilds {
			if botGuild.ID == userGuild.ID {
				commonGuilds = append(commonGuilds, userGuild)
				break
			}
		}
	}
	res = shodan.ResponseEnvelope{
		Status: 200,
		Data: struct {
			User        *discordgo.User        `json:"user"`
			Guilds      []*discordgo.UserGuild `json:"commonGuilds"`
			OwnedGuilds []*discordgo.UserGuild `json:"ownedGuilds"`
		}{u, commonGuilds, ownedGuilds},
	}
	shodan.SendResponse(w, &res)
}

func handleGetRolesForGuild(w http.ResponseWriter, r *http.Request) {
	var res shodan.ResponseEnvelope
	req, err := shodan.ReadRequest(r)
	if err != nil {
		shodan.HttpSendError(w, err)
		return
	}
	if !req.Authenticated() {
		shodan.HttpSendError(w, shodan.ErrorHttp("Authorization required", 401))
		return
	}
	// get the guild argument from the request
	requestVars := mux.Vars(r)
	guilds, err := shodan.DUGetUserGuilds(req.Token)
	if err != nil {
		shodan.HttpSendError(w, err)
		return
	}
	isJoinedGuild := false
	for _, guild := range guilds {
		if guild.ID == requestVars["guild"] {
			isJoinedGuild = true
		}
	}
	if !isJoinedGuild {
		shodan.HttpSendError(w, shodan.ErrorHttp("Guild access denied", 403))
		return
	}
	data, err := DBGetAvailableRoles(mod.shodan.GetDatabase(), requestVars["guild"])
	if err != nil {
		shodan.HttpSendError(w, err)
		return
	}
	res = shodan.ResponseEnvelope{
		Status: 200,
		Data:   data,
	}
	shodan.SendResponse(w, &res)
}

// safely copy guilds from state
func copyGuildState(state *discordgo.State) []discordgo.Guild {
	state.RLock()
	defer state.RUnlock()
	retVal := make([]discordgo.Guild, len(state.Guilds))
	for i, g := range state.Guilds {
		retVal[i] = *g
	}
	return retVal
}
