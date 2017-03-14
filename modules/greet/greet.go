package greet

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/innovandalism/shodan"
)

var (
	channelIDKeyFormat = "%s_greetChannelId"
	messageKeyFormat   = "%s_greetMessage"
)

func getHandleGuildMemberAdd(s shodan.Shodan) interface{} {
	return func(session *discordgo.Session, event *discordgo.GuildMemberAdd) {
		// check if we have a channel ID set for this guild
		guildKeyChannelID := fmt.Sprintf(channelIDKeyFormat, event.GuildID)
		hasKey, err := s.GetRedis().HasKey(guildKeyChannelID)
		shodan.ReportThreadError(false, err)

		// bail if not configured
		if !hasKey {
			return
		}

		// grab the channel ID from redis
		channelID, err := s.GetRedis().Get(guildKeyChannelID)
		shodan.ReportThreadError(false, err)

		// check if the channel still exists
		// TODO: is it better to just crash later if it breaks? It's just one extra query to redis after all
		_, err = session.Channel(channelID)
		shodan.ReportThreadError(false, err)

		// check if we have a custom message in redis
		guildKeyGreetMessage := fmt.Sprintf(messageKeyFormat, event.GuildID)
		hasKey, err = s.GetRedis().HasKey(guildKeyGreetMessage)
		shodan.ReportThreadError(false, err)

		messageFormat := "Hi %s!"

		// replace default message if we do
		if hasKey {
			messageFormat, err = s.GetRedis().Get(guildKeyGreetMessage)
			shodan.ReportThreadError(false, err)
		}

		message := fmt.Sprintf(messageFormat, "<@!"+event.User.ID+">")

		session.ChannelMessageSend(channelID, message)
		shodan.ReportThreadError(false, err)
	}
}
