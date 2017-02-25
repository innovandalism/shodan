package log

import "time"

// A message to be logged to the backend database
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
