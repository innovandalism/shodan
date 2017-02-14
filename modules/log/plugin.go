// +build linux
package log

import (
	"flag"
	"github.com/innovandalism/shodan"
)

type Module struct {
	mongodb *string
}

var mod = Module{}

func init() {
	shodan.Loader.LoadModule(&mod)
}

func (_ *Module) GetIdentifier() string {
	return "log"
}

func (l *Module) FlagHook() {
	l.mongodb = flag.String("mongo", "", "MongoDB database URI")
}

func (l *Module) Attach(session *shodan.Shodan) {
	dataChannel = make(chan *LogMessage)
	session.GetDiscord().AddHandler(onMessageCreate)
	session.GetDiscord().AddHandler(onMessageUpdate)
	go MongoLogger(l.mongodb, dataChannel)
}
