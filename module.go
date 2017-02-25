package shodan

import (
	"flag"
)

// We use a singleton here because the modules need to be registered without context
var Loader = &ModuleLoader{
	Modules: []ModuleInstance{},
}

// Type ModuleInstance holds a module and its current state.
// This glue is needed because the Modules are global to the process.
type ModuleInstance struct {
	Module  *Module
	Enabled *bool
}

// A module represents a loadable piece of code that has been added at compile-time
// Modules are, by nature of them being loaded though init(), global.
type Module interface {
	GetIdentifier() string
	Attach(*Shodan)
	FlagHook()
}

// A ModuleLoader, to hold ModuleInstances.
type ModuleLoader struct {
	Modules []ModuleInstance
}

// Adds a module to the loader. Modules are disabled by default.
func (loader *ModuleLoader) LoadModule(m Module) {
	enabled := false
	instance := ModuleInstance{
		Module:  &m,
		Enabled: &enabled,
	}
	loader.Modules = append(loader.Modules, instance)
}

// Notifies all modules that it is now time to register flags if any are required.
// Also adds m_modulename boolean flag to enable the module. All modules are disabled by default.
func (loader *ModuleLoader) FlagHook() {
	for _, moduleInstance := range Loader.Modules {
		m := *moduleInstance.Module
		flag.BoolVar(moduleInstance.Enabled, "m_"+m.GetIdentifier(), false, "enable "+m.GetIdentifier()+" module")
		m.FlagHook()
	}
}

// Attaches enabled modules to the session.
func (loader *ModuleLoader) Attach(session *Shodan) {
	for _, moduleInstance := range Loader.Modules {
		m := *moduleInstance.Module
		if *moduleInstance.Enabled {
			m.Attach(session)
		}
	}
}
