package greet

import (
	"github.com/innovandalism/shodan"
	"github.com/bwmarrin/discordgo"
	"fmt"
	"github.com/innovandalism/shodan/util"
)

var greetChannelIdKeyFormat string = "%s_greetChannelId"
var greetMessageFormat string = "%s_greetMessage"

func Greeter(shodan *shodan.Shodan) {
	shodan.GetDiscord().AddHandler(getHandleGuildMemberAdd(shodan))
}

// I'm feeling javascripty
func getHandleGuildMemberAdd(shodan *shodan.Shodan) interface{} {
	return func(session *discordgo.Session, event *discordgo.GuildMemberAdd) {
		// check if we have a channel ID set for this guild
		guildKeyChannelId := fmt.Sprintf(greetChannelIdKeyFormat, event.GuildID)
		hasKey, err := shodan.GetRedis().HasKey(guildKeyChannelId)
		util.ReportThreadError(false, err)

		// bail if not configured
		if !hasKey {
			return
		}

		// grab the channel ID from redis
		channelId, err := shodan.GetRedis().GetString(guildKeyChannelId)
		util.ReportThreadError(false, err)

		// check if the channel still exists
		// TODO: is it better to just crash later if it breaks? It's just one extra query to redis after all
		_, err = session.Channel(channelId)
		util.ReportThreadError(false, err)

		// default message
		greetMessageFormat := "Hi %s! Welcome to the server!"

		// check if we have a custom message in redis
		guildKeyGreetMessage := fmt.Sprintf(greetMessageFormat, event.GuildID)
		hasKey, err = shodan.GetRedis().HasKey(guildKeyGreetMessage)
		util.ReportThreadError(false, err)

		// replace default message if we do
		if hasKey {
			greetMessageFormat, err = shodan.GetRedis().GetString(guildKeyGreetMessage)
			util.ReportThreadError(false, err)
		}

		//
		greetMessage := fmt.Sprintf(greetMessageFormat, "<@!" + event.User.ID + ">")

		session.ChannelMessageSend(channelId, greetMessage)
		util.ReportThreadError(false, err)
	}
}
