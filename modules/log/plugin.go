package log

import (
	"github.com/bwmarrin/discordgo"
	"github.com/innovandalism/shodan"
	"log"
)

// Module holds data for this module and implements the shodan.Module interface
type Module struct {
	shodan shodan.Shodan
}

var mod = Module{}

func init() {
	shodan.Loader.LoadModule(&mod)
}

// GetIdentifier returns the identifier for this module
func (m *Module) GetIdentifier() string {
	return "log"
}

// Attach attaches this module to a Shodan session
func (m *Module) Attach(session shodan.Shodan) {
	m.shodan = session
	session.GetDiscord().AddHandler(onMessageCreate)
}

func onMessageCreate(s *discordgo.Session, message *discordgo.MessageCreate) {
	err := shodan.DBUpsertUser(mod.shodan.GetDatabase(), message.Author)
	if err != nil {
		log.Printf("logger: err user: %s\n", err.Error())
		return
	}
	channel, err := s.State.Channel(message.ChannelID)
	if err != nil {
		log.Printf("logger: err channel: %s\n", err.Error())
		return
	}
	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		log.Printf("logger: err guild: %s\n", err.Error())
		return
	}
	err = shodan.DBUpsertGuild(mod.shodan.GetDatabase(), guild)
	if err != nil {
		log.Printf("logger: err guild: %s\n", err.Error())
		return
	}
	err = shodan.DBUpsertChannel(mod.shodan.GetDatabase(), channel)
	if err != nil {
		log.Printf("logger: err channel: %s\n", err.Error())
		return
	}
	err = shodan.DBUpsertMessage(mod.shodan.GetDatabase(), message.Message)
	if err != nil {
		log.Printf("logger: err message: %s\n", err.Error())
		return
	}
	log.Printf("logged: %s - %s - %s\n", message.ChannelID, message.Author.ID, message.Content)
}
