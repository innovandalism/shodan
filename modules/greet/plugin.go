package greet

import (
	"github.com/innovandalism/shodan"
)

// Module holds this modules data and methods
type Module struct {}

var mod = Module{}

func init() {
	shodan.Loader.LoadModule(&mod)
}

// GetIdentifier returns the name of the module. Purely used for statistics and flag registration
func (m *Module) GetIdentifier() string {
	return "greet"
}

// FlagHook registers any flags needed for this module
func (m *Module) FlagHook() {

}

// Attach attaches functionality in this module to Shodan
func (m *Module) Attach(shodan *shodan.Shodan) {
	shodan.GetCommandStack().RegisterCommand(&ChannelCmd{})
	shodan.GetDiscord().AddHandler(getHandleGuildMemberAdd(shodan))
}