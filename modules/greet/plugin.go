package greet

import (
	"github.com/innovandalism/shodan"
)

type Module struct {}

var mod = Module{}

func init() {
	shodan.Loader.LoadModule(&mod)
}

func (_ *Module) GetIdentifier() string {
	return "greet"
}

func (m *Module) FlagHook() {

}

func (m *Module) Attach(session *shodan.Shodan) {
	session.GetCommandStack().RegisterCommand(&GreetChannelCmd{})
}