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

func (vc *VersionCommand) Invoke(ci *shodan.CommandInvocation) bool {
	v := fmt.Sprintf("%d.%d.%d", config.VersionMajor, config.VersionMinor, config.VersionRevision)
	me := &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       "SHODAN",
		Description: "A platform for building Discord Bots in Go",
		Fields: []*discordgo.MessageEmbedField{
			{"Version", v, true},
			{"Uptime", time.Now().Sub(vc.startupTime).String(), true},
			{"Git Hash", config.VersionGitHash, false},
			{"Modules", getPluginList(), false},
			{"Code", "https://github.com/innovandalism/shodan", false},
		},
		URL:    "https://shodan-bot.neocities.org/",
		Author: &discordgo.MessageEmbedAuthor{URL: "https://innovandalism.eu", Name: "Innovandalism"},
	}
	_, err := ci.Session.ChannelMessageSendEmbed(ci.Event.ChannelID, me)
	if err != nil {
		util.ReportThreadError(false, err)
	}
	return true
}

func getPluginList() string {
	plugins := ""
	for _, moduleInstance := range shodan.Loader.Modules {
		if !*moduleInstance.Enabled {
			continue
		}
		module := *moduleInstance.Module
		plugins = plugins + fmt.Sprintf("\nm_"+module.GetIdentifier())
	}
	return plugins
}
