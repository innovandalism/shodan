package version

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/innovandalism/shodan"
	"github.com/innovandalism/shodan/config"
	"github.com/innovandalism/shodan/util"
	"runtime"
	"time"
)

type VersionCommand struct {
	startupTime time.Time
	memStats    *runtime.MemStats
}

func (vc *VersionCommand) GetNames() []string {
	return []string{"version", "info"}
}

func (vc *VersionCommand) Invoke(ci *shodan.CommandInvocation) error {
	v := fmt.Sprintf("%d.%d.%d", config.VersionMajor, config.VersionMinor, config.VersionRevision)
	me := &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       "SHODAN",
		Description: "A platform for building discordgo Bots in Go",
		Fields: []*discordgo.MessageEmbedField{
			{"Version", v, true},
			{"Uptime", time.Now().Sub(vc.startupTime).String(), true},
			{"Git Hash", config.VersionGitHash, false},
			{"Code", "https://github.com/innovandalism/shodan", false},
		},
		URL:    "https://github.com/innovandalism/shodan",
		Author: &discordgo.MessageEmbedAuthor{URL: "https://innovandalism.eu", Name: "Innovandalism"},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "This build of Shodan has super-kapow-powers.",
		},
	}
	err := ci.Helpers.ReplyEmbed(me)
	if err != nil {
		return util.WrapError(err)
	}
	return nil
}