// Package shodan provides core functionality for writing a Discord bot
package shodan

import (
	"errors"
	"flag"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/innovandalism/shodan/config"
	shodanRedis "github.com/innovandalism/shodan/services/redis"
	"github.com/innovandalism/shodan/util"
	"github.com/vharitonsky/iniflags"
	"math/rand"
	"os"
	"time"
	"gopkg.in/redis.v5"
)

// Invokes the main loop, starting HTTP listeners and Discordgo.
// Meant to be called from an external package after importing additional modules.
func Run() {
	// catch-all error handler
	defer util.ErrorHandler()

	// seed the RNG with a somewhat random number. don't do crypto with this kids
	rand.Seed(time.Now().UTC().UnixNano())

	session := Init()

	// essential flags
	var (
		token     = flag.String("token", config.DefaultToken, "Discord Bot Authentication Token")
		addr      = flag.String("web_addr", "", "Address to bind HTTP listener to")
		noweb     = flag.Bool("web_disable", false, "Disable WebUI")
		dsn       = flag.String("dsn", "", "Sentry DSN")
		redis_uri = flag.String("redis", "redis://127.0.0.1/", "Redis URI")
	)

	// notify all modules that this is the last chance to ask for flags
	session.FlagHook()

	// use iniflags instead of flag, this allows for loading ini files
	iniflags.Parse()

	// check our required set of flags
	if len(*token) == 0 {
		panic(errors.New("Discord token must be provided"))
	}

	if len(*addr) == 0 && !*noweb {
		panic(errors.New("either web_addr or web_disable must be provided"))
	}

	if len(*dsn) != 0 {
		raven.SetDSN(*dsn)
	}

	// set up redis
	session.redis = &shodanRedis.ShodanRedis{}
	options, err := redis.ParseURL(*redis_uri)
	util.PanicOnError(err)
	session.redis.Init(options)

	// raven should receive the git hash when crashing
	raven.SetRelease(config.VersionGitHash)

	err = session.InitDiscord(token)
	util.PanicOnError(err)

	if !*noweb {
		go session.InitHttp(*addr)
	}

	session.Connect()

	// because discordgo returns immediately and starts its goroutines in the background, we defer closing the session here
	defer session.GetDiscord().Close()

	// Setup channels to react to in the main goroutine
	signalChannel := util.MakeSignalChannel()
	threadErrorChannel := util.GetThreadErrorChannel()

	for {
		quit := false
		select {
		case <-signalChannel:
			quit = true
		case threadError := <-threadErrorChannel:
			if threadError.IsFatal {
				os.Exit(0xDEAD)
			} else {
				fmt.Fprintf(os.Stderr, "non-fatal thread error: %s\n", threadError.Error)
			}
		}
		if quit {
			break
		}
	}
}
