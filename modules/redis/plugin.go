package redis

import (
	"flag"
	"github.com/innovandalism/shodan"
	"github.com/innovandalism/shodan/util"
	"gopkg.in/redis.v5"
)

type Module struct {
	session *shodan.Shodan
	addr    *string
	pass    *string
	db      *int
	Client  *redis.Client
}

var mod = Module{}

func init() {
	shodan.Loader.LoadModule(&mod)
}

func (_ *Module) GetIdentifier() string {
	return "redis"
}

func (m *Module) FlagHook() {
	m.addr = flag.String("redis_addr", "", "Redis Address")
	m.pass = flag.String("redis_password", "", "Redis Password")
	m.db = flag.Int("redis_db", 1, "Redis DB")
}

func (m *Module) Attach(session *shodan.Shodan) {
	m.session = session
	mod.Client = redis.NewClient(&redis.Options{
		Addr:     *mod.addr,
		Password: *mod.pass,
		DB:       *mod.db,
	})
	// bail if redis crashed - since it has to be explicitly enabled failing is fatal
	util.ReportThreadError(true, mod.Client.Ping().Err())
}
