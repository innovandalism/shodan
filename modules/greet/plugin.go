package greet

import (
	"github.com/innovandalism/shodan"
)

// Module holds data for this module and implements the shodan.Module interface
type Module struct{}

var mod = Module{}

func init() {
	shodan.Loader.LoadModule(&mod)
}

// GetIdentifier returns the identifier for this module
func (m *Module) GetIdentifier() string {
	return "greet"
}

// Attach attaches this module to a Shodan session
func (m *Module) Attach(session shodan.Shodan) error {
	if session.GetRedis() == nil {
		return shodan.Error("mod_greet: redis is nil; make sure the redis driver loads or disable this module")
	}
	session.GetCommandStack().RegisterCommand(&ChannelCmd{})
	session.GetDiscord().AddHandler(getHandleGuildMemberAdd(session))
	return nil
}
