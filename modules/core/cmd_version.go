package version

import (
	"fmt"
	"runtime"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/innovandalism/shodan"
	"github.com/innovandalism/shodan/config"
)

// VersionCommand holds information for this command, including bot startup time and a MemStats struct pointer
type VersionCommand struct {
	startupTime time.Time
	memStats    *runtime.MemStats
}

// GetNames returns the command aliases for this command
func (vc *VersionCommand) GetNames() []string {
	return []string{"version", "info"}
}

// Invoke runs the command. This commands returns a rich embed with version and runtime information.
func (vc *VersionCommand) Invoke(ci *shodan.CommandInvocation) error {
	v := fmt.Sprintf("%d.%d.%d", config.VersionMajor, config.VersionMinor, config.VersionRevision)
	me := &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       config.OEMName,
		Description: config.OEMDescription,
		Fields: []*discordgo.MessageEmbedField{
			{"Version", v, true},
			{"Uptime", time.Now().Sub(vc.startupTime).String(), true},
			{"Git Hash", config.VersionGitHash, false},
			{"Code", config.OEMBuildfile, false},
		},
		URL:    config.OEMVendorURI,
		Author: &discordgo.MessageEmbedAuthor{URL: config.OEMVendorURI, Name: config.OEMVendor},
		Footer: &discordgo.MessageEmbedFooter{
			Text: config.OEMAnecdote,
		},
	}
	err := ci.Helpers.ReplyEmbed(me)
	if err != nil {
		return shodan.WrapError(err)
	}
	return nil
}
