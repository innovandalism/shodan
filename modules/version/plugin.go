package version

import (
	"github.com/innovandalism/shodan"
	"runtime"
	"time"
)

// Module holds data for this module and implements the shodan.Module interface
type Module struct {
	shodan   shodan.Shodan
	memStats *runtime.MemStats
}

var mod = Module{}

func init() {
	shodan.Loader.LoadModule(&mod)
}

// GetIdentifier returns the identifier for this module
func (_ *Module) GetIdentifier() string {
	return "version"
}

// Attach attaches this module to a Shodan session
func (m *Module) Attach(session shodan.Shodan) {
	m.shodan = session

	memStats := runtime.MemStats{}
	mod.memStats = &memStats

	session.GetCommandStack().RegisterCommand(&MemoryCommand{})
	session.GetCommandStack().RegisterCommand(&VersionCommand{startupTime: time.Now()})
	session.GetCommandStack().RegisterCommand(&IgnoreCommand{})
}
