package webui

import (
	"fmt"
	"github.com/innovandalism/shodan"
	"net/http"
	"strings"
)

func registerStatusPage(s *shodan.Shodan) {
	s.GetMux().HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		discord := s.GetDiscord()
		ug, err := discord.UserGuilds()
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
	})
}
