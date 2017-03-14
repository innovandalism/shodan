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

// FlagHook triggers before flags are parsed to allow this module to add options
func (m *Module) FlagHook() {

}

// Attach attaches this module to a Shodan session
func (m *Module) Attach(shodan shodan.Shodan) {
	shodan.GetCommandStack().RegisterCommand(&ChannelCmd{})
	shodan.GetDiscord().AddHandler(getHandleGuildMemberAdd(shodan))
}
