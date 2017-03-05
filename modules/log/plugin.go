package log

import (
	"github.com/innovandalism/shodan"
	"github.com/innovandalism/shodan/util"
	"github.com/bwmarrin/discordgo"
	"time"
	"log"
)

type Module struct {}

var mod = Module{}

func init() {
	shodan.Loader.LoadModule(&mod)
}

func (_ *Module) GetIdentifier() string {
	return "log"
}

func (m *Module) FlagHook() {

}

func (m *Module) Attach(session *shodan.Shodan) {
	stmtInsertGuild, err := session.GetDatabase().DB.Prepare("INSERT INTO discord_guild (id, name) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET name = $2")
	if err != nil {
		util.ReportThreadError(true, util.WrapError(err))
	}
	stmtInsertUser, err := session.GetDatabase().DB.Prepare("INSERT INTO discord_user (id, name, discriminator) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET name = $2,discriminator = $3")
	if err != nil {
		util.ReportThreadError(true, util.WrapError(err))
	}
	stmtInsertChannel, err := session.GetDatabase().DB.Prepare("INSERT INTO discord_channel (id, guild, name) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET guild = $2,name = $3")
	if err != nil {
		util.ReportThreadError(true, util.WrapError(err))
	}
	stmtInsertMessage, err := session.GetDatabase().DB.Prepare("INSERT INTO discord_message (channel, \"user\", message, time) VALUES ($1, $2, $3, $4)")
	if err != nil {
		util.ReportThreadError(true, util.WrapError(err))
	}
	session.GetDiscord().AddHandler(func(s *discordgo.Session, event *discordgo.MessageCreate) {
		_, err := stmtInsertUser.Exec(event.Author.ID, event.Author.Username, event.Author.Discriminator)
		if err != nil {
			log.Printf("logger: err: %s\n", err.Error())
			return
		}
		channel, err := s.State.Channel(event.ChannelID)
		if err != nil {
			log.Printf("logger: err: %s\n", err.Error())
			return
		}
		guild, err := s.State.Guild(channel.GuildID)
		if err != nil {
			log.Printf("logger: err: %s\n", err.Error())
			return
		}
		_, err = stmtInsertGuild.Exec(guild.ID, guild.Name)
		if err != nil {
			log.Printf("logger: err: %s\n", err.Error())
			return
		}
		_, err = stmtInsertChannel.Exec(channel.ID, channel.GuildID, channel.Name)
		if err != nil {
			log.Printf("logger: err: %s\n", err.Error())
			return
		}
		_, err = stmtInsertMessage.Exec(event.ChannelID, event.Author.ID, event.Message.Content, time.Now())
		if err != nil {
			log.Printf("logger: err: %s\n", err.Error())
			return
		}
		log.Printf("logged: %s - %s - %s\n", event.ChannelID, event.Author.ID, event.Content)
	})
}