package shodan

import (
	"os"
	"strings"
)

// Loader is a singleton available from init functions of plugin packages and exposes a module loader
var Loader = &ModuleLoader{
	Modules: []ModuleInstance{},
}

// ModuleInstance holds a module and its current state.
//
// This glue is needed because the Modules are global to the process.
type ModuleInstance struct {
	Module  Module
	Enabled bool
}

// A Module represents a loadable piece of code that has been added at compile-time
// Modules are, by nature of them being loaded though init(), global.
type Module interface {
	GetIdentifier() string
	Attach(Shodan)
}

// A ModuleLoader holds ModuleInstances
type ModuleLoader struct {
	Modules []ModuleInstance
}

// LoadModule adds a module to the loader. Modules are disabled by default.
func (loader *ModuleLoader) LoadModule(m Module) {
	_, enabled := os.LookupEnv("M_" + strings.ToUpper(m.GetIdentifier()))
	instance := ModuleInstance{
		Module:  m,
		Enabled: enabled,
	}
	loader.Modules = append(loader.Modules, instance)
}

// Attach attaches enabled modules to the session.
//
// After this point it does not matter if a module is marked as enabled or not in the ModuleInstance
func (loader *ModuleLoader) Attach(session Shodan) {
	for _, moduleInstance := range Loader.Modules {
		m := moduleInstance.Module
		if moduleInstance.Enabled {
			m.Attach(session)
		}
	}
}
