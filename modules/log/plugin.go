package log

import (
	"flag"
	"github.com/innovandalism/shodan"
)

// The logging module
type Module struct {
	mongodb *string
}

var mod = Module{}

func init() {
	shodan.Loader.LoadModule(&mod)
}

// Identify this module
func (_ *Module) GetIdentifier() string {
	return "log"
}

// Signal this module to register flags
func (l *Module) FlagHook() {
	l.mongodb = flag.String("mongo", "", "MongoDB database URI")
}

// Gives this module the Shodan object to mutate
func (l *Module) Attach(session *shodan.Shodan) {
	dataChannel = make(chan *LogMessage)
	session.GetDiscord().AddHandler(onMessageCreate)
	session.GetDiscord().AddHandler(onMessageUpdate)
	go MongoLogger(l.mongodb, dataChannel)
}
