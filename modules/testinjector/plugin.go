// Package testinjector allows attaching to a shodan session for running unit tests
package testinjector

import "github.com/innovandalism/shodan"

// Module holds AttachHandler, a function to invoke when the module is attached
type Module struct {
	AttachHandler func(session *shodan.Shodan)
}

func (_ *Module) GetIdentifier() string {
	return "testinjector"
}

func (m *Module) FlagHook() {}

func (m *Module) Attach(session *shodan.Shodan) {
	m.AttachHandler(session)
}