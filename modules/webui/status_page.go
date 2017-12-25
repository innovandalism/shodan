package webui

import (
	"fmt"
	"net/http"
	"strings"
)

func handleStatusPage(w http.ResponseWriter, r *http.Request) {
	discord := mod.shodan.GetDiscord()
	ug, err := discord.UserGuilds(100, "", "")
	if err != nil {
		panic(err)
	}
	guildsWrapperFmt := `Guilds connected to:
%s
`
	guildFmt := "%s\t%b\t%s"
	guilds := []string{}
	for _, g := range ug {
		guilds = append(guilds, fmt.Sprintf(guildFmt, g.ID, g.Permissions, g.Name))
	}
	fmt.Fprintf(w, guildsWrapperFmt, strings.Join(guilds, "\n"))
}
