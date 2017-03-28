// Package shodan provides core functionality for writing a Discord bot
package shodan

import (
	"errors"
	"flag"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/innovandalism/shodan/config"
	"github.com/vharitonsky/iniflags"
	"math/rand"
	"os"
	"time"
)

// Run invokes the main loop, core services like HTTP and discordgo.
//
// Meant to be called from an external package after importing additional modules.
func Run() {
	// catch-all error handler
	defer errorHandler()

	// seed the RNG with a somewhat random number. don't do crypto with this kids
	rand.Seed(time.Now().UTC().UnixNano())

	// essential flags
	var (
		err      error
		sgm	 = flag.String("sgm", "", "Single Guild Mode GuildID")
		token    = flag.String("token", config.DefaultToken, "discordgo Bot Authentication Token")
		addr     = flag.String("web_addr", "", "Address to bind HTTP listener to")
		noweb    = flag.Bool("web_disable", false, "Disable WebUI")
		dsn      = flag.String("dsn", "", "Sentry DSN")
		redisURI = flag.String("redis", "redis://127.0.0.1/", "Redis URI")
		pgURI    = flag.String("postgres", "postgres://127.0.0.1/shodan", "Postgres URI")
	)

	// notify all modules that this is the last chance to ask for flags
	Loader.FlagHook()

	// use iniflags instead of flag, this allows for loading ini files
	iniflags.Parse()

	// Single-Guild-Mode hides all multiguild UI and pre-populates the selected guild
	config.SingleGuildMode = *sgm

	// check our required set of flags
	if len(*token) == 0 {
		panic(errors.New("discordgo token must be provided"))
	}

	if len(*addr) == 0 && !*noweb {
		panic(errors.New("either web_addr or web_disable must be provided"))
	}

	if len(*dsn) != 0 {
		raven.SetDSN(*dsn)
	}

	session := shodanSession{}
	session.moduleLoader = Loader

	// set up external services
	session.cmdStack, err = InitCommandStack()
	panicOnError(err)

	session.kvs, err = InitRedis(*redisURI)
	panicOnError(err)

	session.database, err = InitPostgres(*pgURI)
	panicOnError(err)

	session.discord, err = InitDiscord(*token)
	panicOnError(err)

	session.mux, err = InitHTTP(*addr)
	panicOnError(err)

	// raven should receive the git hash when crashing
	raven.SetRelease(config.VersionGitHash)

	session.Bootstrap()

	// because discordgo returns immediately and starts its goroutines in the background, we defer closing the session here
	defer session.GetDiscord().Close()

	// Setup channels to react to in the main goroutine
	signalChannel := makeSignalChannel()
	threadErrorChannel := getThreadErrorChannel()

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
