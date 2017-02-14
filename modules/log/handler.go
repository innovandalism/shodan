package log

import (
	"github.com/bwmarrin/discordgo"
	"github.com/innovandalism/shodan/util"
)

var dataChannel chan *LogMessage

func onMessageCreate(s *discordgo.Session, event *discordgo.MessageCreate) {
	channel, err := s.State.Channel(event.ChannelID)
	if err != nil {
		util.ReportThreadError(false, err)
	}
	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		util.ReportThreadError(false, err)
	}
	time, err := event.Timestamp.Parse()
	if err != nil {
		util.ReportThreadError(false, err)
	}
	m := LogMessage{
		ID:                  event.ID,
		Channel:             channel.ID,
		ChannelName:         channel.Name,
		Guild:               guild.ID,
		GuildName:           guild.Name,
		Author:              event.Author.ID,
		AuthorName:          event.Author.Username,
		AuthorDiscriminator: event.Author.Discriminator,
		IsEdit:              false,
		Time:                time,
		Message:             event.Message.ContentWithMentionsReplaced(),
	}
	dataChannel <- &m
}

func onMessageUpdate(_ *discordgo.Session, event *discordgo.MessageUpdate) {
	time, _ := event.Timestamp.Parse()
	m := LogMessage{
		ID:      event.ID,
		IsEdit:  true,
		Time:    time,
		Message: event.Message.ContentWithMentionsReplaced(),
	}
	dataChannel <- &m
}
