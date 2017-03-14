package shodan

import (
	discordgo "github.com/bwmarrin/discordgo"
	"database/sql"
	"github.com/gorilla/mux"
)

type Shodan interface {
	GetDiscord() *discordgo.Session
	GetDatabase() *sql.DB
	GetMux() *mux.Router
	GetCommandStack() *CommandStack
	GetRedis() KVS
}

// Shodan is the central object holding all services
type shodanSession struct {
	discord      *discordgo.Session
	moduleLoader *ModuleLoader
	mux          *mux.Router
	cmdStack     *CommandStack
	kvs          KVS
	database     *sql.DB
}

// GetDiscord returns the discordgo session
func (session *shodanSession) GetDiscord() *discordgo.Session {
	return session.discord
}

// GetDatabase returns the DAL (Database Access Layer)
func (session *shodanSession) GetDatabase() *sql.DB {
	return session.database
}

// GetMux returns the HTTP Muxer
func (session *shodanSession) GetMux() *mux.Router {
	return session.mux
}

// GetCommandStack returns the command stack
func (session *shodanSession) GetCommandStack() *CommandStack {
	return session.cmdStack
}

// GetRedis returns the attached redis key-value-store
func (session *shodanSession) GetRedis() KVS {
	return session.kvs
}

// Bootstrap wires all services up with each other
func (session *shodanSession) Bootstrap() error {
	session.cmdStack.Attach(session)
	session.moduleLoader.Attach(session)
	session.discord.Open()

	return nil
}