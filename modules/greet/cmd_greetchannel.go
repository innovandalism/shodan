package greet

import (
	"github.com/innovandalism/shodan"
	"github.com/innovandalism/shodan/util"
	"fmt"
	"strings"
	"github.com/bwmarrin/discordgo"
)

type GreetChannelCmd struct{}

func (_ *GreetChannelCmd) GetNames() []string {
	return []string{"greet"}
}

func (_ *GreetChannelCmd) Invoke(ci *shodan.CommandInvocation) bool {
	if len(ci.Arguments) < 1 {
		err := ci.Helpers.Reply("greet [status|set|clear|msg]")
		util.ReportThreadError(false, err)
		return false
	}
	channel, err := ci.Session.State.Channel(ci.Event.ChannelID)
	util.ReportThreadError(false, err)
	guildKeyChannelId := fmt.Sprintf(greetChannelIdKeyFormat, channel.GuildID)
	guildKeyGreetMessage := fmt.Sprintf(greetMessageFormat, channel.GuildID)
	switch ci.Arguments[0] {
	case "status":
		hasKey, err := ci.Shodan.GetRedis().HasKey(guildKeyChannelId)
		util.ReportThreadError(false, err)
		if hasKey {
			err = ci.Helpers.Reply("Key is set.")
			util.ReportThreadError(false, err)
		} else {
			err = ci.Helpers.Reply("Key is not set.")
			util.ReportThreadError(false, err)
		}
	case "set":
		err := ci.Shodan.GetRedis().SetPerm(guildKeyChannelId, ci.Event.ChannelID)
		util.ReportThreadError(false, err)
		err = ci.Helpers.Reply("This channel is now the greet channel")
		util.ReportThreadError(false, err)
	case "clear":
		_, err := ci.Shodan.GetRedis().Clear(guildKeyChannelId)
		util.ReportThreadError(false, err)
		_, err = ci.Session.ChannelMessageSend(ci.Event.ChannelID, "Key deleted.")
		util.ReportThreadError(false, err)
	case "msg":
		if len(ci.Arguments) < 2 {
			err = ci.Helpers.Reply("You need to specify a message. To disable greetings, try `greet clear`")
			util.ReportThreadError(false, err)
			return false
		}
		msg := strings.Join(ci.Arguments[1:], " ")
		err := ci.Shodan.GetRedis().SetPerm(guildKeyGreetMessage, msg)
		util.ReportThreadError(false, err)
		err = ci.Helpers.Reply("Message set.")
		util.ReportThreadError(false, err)
	}
	return true
}

func (_ *GreetChannelCmd) GetRequiredPermission() int {
	return discordgo.PermissionManageServer
}