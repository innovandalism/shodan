package greet

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/innovandalism/shodan"
	"github.com/innovandalism/shodan/util"
	"strings"
)

// A ChannelCmd holds methods for this command
type ChannelCmd struct{}

// GetNames returns the command aliases for this command
func (c *ChannelCmd) GetNames() []string {
	return []string{"greet"}
}

// Invoke runs the command
func (c *ChannelCmd) Invoke(ci *shodan.CommandInvocation) error {
	if len(ci.Arguments) < 1 {
		err := ci.Helpers.Reply("greet [status|set|clear|msg]")
		if err != nil {
			return util.WrapError(err)
		}
		return nil
	}
	channel, err := ci.Shodan.GetDiscord().State.Channel(ci.Event.ChannelID)
	if err != nil {
		return util.WrapError(err)
	}
	guildKeyChannelID := fmt.Sprintf(channelIDKeyFormat, channel.GuildID)
	guildKeyGreetMessage := fmt.Sprintf(messageKeyFormat, channel.GuildID)
	switch ci.Arguments[0] {
	case "status":
		hasKey, err := ci.Shodan.GetRedis().HasKey(guildKeyChannelID)
		if err != nil {
			return util.WrapError(err)
		}
		if hasKey {
			if err != nil {
				return util.WrapError(err)
			}
			util.ReportThreadError(false, err)
		} else {
			if err != nil {
				return util.WrapError(err)
			}
			util.ReportThreadError(false, err)
		}
	case "set":
		err := ci.Shodan.GetRedis().Set(guildKeyChannelID, ci.Event.ChannelID)
		if err != nil {
			return util.WrapError(err)
		}
		err = ci.Helpers.Reply("This channel is now the greet channel")
		if err != nil {
			return util.WrapError(err)
		}
	case "clear":
		err := ci.Shodan.GetRedis().Clear(guildKeyChannelID)
		if err != nil {
			return util.WrapError(err)
		}
		err = ci.Helpers.Reply("Key deleted.")
		if err != nil {
			return util.WrapError(err)
		}
	case "msg":
		if len(ci.Arguments) < 2 {
			err = ci.Helpers.Reply("You need to specify a message. To disable greetings, try `greet clear`")
			if err != nil {
				return util.WrapError(err)
			}
			return nil
		}
		msg := strings.Join(ci.Arguments[1:], " ")
		err := ci.Shodan.GetRedis().Set(guildKeyGreetMessage, msg)
		if err != nil {
			return util.WrapError(err)
		}
		err = ci.Helpers.Reply("Message set.")
		if err != nil {
			return util.WrapError(err)
		}
	}
	return nil
}

// GetRequiredPermission returns permission bits for discord ACL permission system
func (c *ChannelCmd) GetRequiredPermission() int {
	return discordgo.PermissionManageServer
}
