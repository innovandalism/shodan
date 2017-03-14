// Package testinjector allows attaching to a shodan session for running unit tests
package testinjector

import "github.com/innovandalism/shodan"

// Module holds data for this module and implements the shodan.Module interface
type Module struct {
	AttachHandler func(session *shodan.Shodan)
}

// GetIdentifier returns the identifier for this module
func (_ *Module) GetIdentifier() string {
	return "testinjector"
}

// FlagHook triggers before flags are parsed to allow this module to add options
func (m *Module) FlagHook() {}

// Attach attaches this module to a Shodan session
func (m *Module) Attach(session *shodan.Shodan) {
	m.AttachHandler(session)
}
