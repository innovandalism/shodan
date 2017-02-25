package greet

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/innovandalism/shodan"
	"github.com/innovandalism/shodan/util"
	"strings"
)

type GreetChannelCmd struct{}

func (_ *GreetChannelCmd) GetNames() []string {
	return []string{"greet"}
}

func (_ *GreetChannelCmd) Invoke(ci *shodan.CommandInvocation) error {
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
	guildKeyChannelId := fmt.Sprintf(greetChannelIdKeyFormat, channel.GuildID)
	guildKeyGreetMessage := fmt.Sprintf(greetMessageFormat, channel.GuildID)
	switch ci.Arguments[0] {
	case "status":
		hasKey, err := ci.Shodan.GetRedis().HasKey(guildKeyChannelId)
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
		err := ci.Shodan.GetRedis().SetPerm(guildKeyChannelId, ci.Event.ChannelID)
		if err != nil {
			return util.WrapError(err)
		}
		err = ci.Helpers.Reply("This channel is now the greet channel")
		if err != nil {
			return util.WrapError(err)
		}
	case "clear":
		_, err := ci.Shodan.GetRedis().Clear(guildKeyChannelId)
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
		err := ci.Shodan.GetRedis().SetPerm(guildKeyGreetMessage, msg)
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

func (_ *GreetChannelCmd) GetRequiredPermission() int {
	return discordgo.PermissionManageServer
}
