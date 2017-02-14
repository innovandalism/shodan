package version

import (
	"github.com/innovandalism/shodan"
	"runtime"
	"time"
)

type Module struct {
	shodan   *shodan.Shodan
	memStats *runtime.MemStats
}

var mod = Module{}

func init() {
	shodan.Loader.LoadModule(&mod)
}

func (_ *Module) GetIdentifier() string {
	return "version"
}

func (m *Module) FlagHook() {

}

func (m *Module) Attach(session *shodan.Shodan) {
	m.shodan = session

	memStats := runtime.MemStats{}
	mod.memStats = &memStats

	session.GetCommandStack().RegisterCommand(&MemoryCommand{})
	session.GetCommandStack().RegisterCommand(&VersionCommand{startupTime: time.Now()})
	session.GetCommandStack().RegisterCommand(&IgnoreCommand{})
}
