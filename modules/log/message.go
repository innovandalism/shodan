package log

import "time"

type LogMessage struct {
	ID                  string
	Guild               string
	GuildName           string
	Channel             string
	ChannelName         string
	Author              string
	AuthorName          string
	AuthorDiscriminator string
	AuthorNick          string
	IsEdit              bool
	Time                time.Time
	Message             string
}
