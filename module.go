package shodan

import (
	"flag"
)

// We use a singleton here because the modules need to be registered without context
var Loader = &ModuleLoader{
	Modules: []ModuleInstance{},
}

type ModuleInstance struct {
	Module  *Module
	Enabled *bool
}

type Module interface {
	GetIdentifier() string
	Attach(*Shodan)
	FlagHook()
}

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

// Notifies all modules that it is now time to register your flags if any are required.
// Also adds m_modulename boolean flag to enable the module. All modules are disabled by default.
func (loader *ModuleLoader) FlagHook() {
	for _, moduleInstance := range Loader.Modules {
		m := *moduleInstance.Module
		flag.BoolVar(moduleInstance.Enabled, "m_"+m.GetIdentifier(), false, "enable "+m.GetIdentifier()+" module")
		m.FlagHook()
	}
}

// Attaches enabled modules by handing them the Shodan object.
func (loader *ModuleLoader) Attach(session *Shodan) {
	for _, moduleInstance := range Loader.Modules {
		m := *moduleInstance.Module
		if *moduleInstance.Enabled {
			m.Attach(session)
		}
	}
}
