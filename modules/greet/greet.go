package greet

import (
	"github.com/innovandalism/shodan"
	"github.com/bwmarrin/discordgo"
	"fmt"
	"github.com/innovandalism/shodan/util"
)

var (
	channelIDKeyFormat = "%s_greetChannelId"
	messageKeyFormat = "%s_greetMessage"
)


func getHandleGuildMemberAdd(shodan *shodan.Shodan) interface{} {
	return func(session *discordgo.Session, event *discordgo.GuildMemberAdd) {
		// check if we have a channel ID set for this guild
		guildKeyChannelID := fmt.Sprintf(channelIDKeyFormat, event.GuildID)
		hasKey, err := shodan.GetRedis().HasKey(guildKeyChannelID)
		util.ReportThreadError(false, err)

		// bail if not configured
		if !hasKey {
			return
		}

		// grab the channel ID from redis
		channelID, err := shodan.GetRedis().Get(guildKeyChannelID)
		util.ReportThreadError(false, err)

		// check if the channel still exists
		// TODO: is it better to just crash later if it breaks? It's just one extra query to redis after all
		_, err = session.Channel(channelID)
		util.ReportThreadError(false, err)

		// check if we have a custom message in redis
		guildKeyGreetMessage := fmt.Sprintf(messageKeyFormat, event.GuildID)
		hasKey, err = shodan.GetRedis().HasKey(guildKeyGreetMessage)
		util.ReportThreadError(false, err)

		messageFormat := "Hi %s!"

		// replace default message if we do
		if hasKey {
			messageFormat, err = shodan.GetRedis().Get(guildKeyGreetMessage)
			util.ReportThreadError(false, err)
		}

		message := fmt.Sprintf(messageFormat, "<@!" + event.User.ID + ">")

		session.ChannelMessageSend(channelID, message)
		util.ReportThreadError(false, err)
	}
}
