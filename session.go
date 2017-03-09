package shodan

import (
	discordgo "github.com/bwmarrin/discordgo"
	"net/http"
	"database/sql"
	"fmt"
)

// Shodan is the central object holding all services
type Shodan struct {
	discord      *discordgo.Session
	moduleLoader *ModuleLoader
	mux          *http.ServeMux
	cmdStack     *CommandStack
	kvs          KVS
	database     *sql.DB
}

// GetDiscord returns the discordgo session
func (session *Shodan) GetDiscord() *discordgo.Session {
	return session.discord
}

// GetDatabase returns the DAL (Database Access Layer)
func (session *Shodan) GetDatabase() *sql.DB {
	return session.database
}

// GetMux returns the HTTP Muxer
func (session *Shodan) GetMux() *http.ServeMux {
	return session.mux
}

// GetCommandStack returns the command stack
func (session *Shodan) GetCommandStack() *CommandStack {
	return session.cmdStack
}

// GetRedis returns the attached redis key-value-store
func (session *Shodan) GetRedis() KVS {
	return session.kvs
}

// Bootstrap wires all services up with each other
func (session *Shodan) Bootstrap() error {
	session.cmdStack.Attach(session)
	session.moduleLoader.Attach(session)
	session.discord.Debug = true
	fmt.Printf("HELLO WORLD")
	session.discord.Open()

	return nil
}