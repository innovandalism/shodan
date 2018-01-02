package log

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/innovandalism/shodan"
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
func (m *Module) Attach(session shodan.Shodan) error {
	m.shodan = session
	if session.GetDatabase() == nil {
		return shodan.Error("mod_log: database is nil; please make sure you're connected to PG or disable the module")
	}
	session.GetDiscord().AddHandler(onMessageCreate)
	return nil
}

func onMessageCreate(s *discordgo.Session, message *discordgo.MessageCreate) {
	err := DBUpsertUser(mod.shodan.GetDatabase(), message.Author)
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
	err = DBUpsertGuild(mod.shodan.GetDatabase(), guild)
	if err != nil {
		log.Printf("logger: err guild: %s\n", err.Error())
		return
	}
	err = DBUpsertChannel(mod.shodan.GetDatabase(), channel)
	if err != nil {
		log.Printf("logger: err channel: %s\n", err.Error())
		return
	}
	err = DBUpsertMessage(mod.shodan.GetDatabase(), message.Message)
	if err != nil {
		log.Printf("logger: err message: %s\n", err.Error())
		return
	}
	log.Printf("logged: %s - %s - %s\n", message.ChannelID, message.Author.ID, message.Content)
}
