// Package shodan provides core functionality for writing a Discord bot
package shodan

import (
	"errors"
	"log"
	"math/rand"
	"os"
	"time"

	raven "github.com/getsentry/raven-go"
	"github.com/innovandalism/shodan/config"
)

// Run invokes the main loop, core services like HTTP and discordgo.
//
// Meant to be called from an external package after importing additional modules.
func Run() {
	// catch-all error handler
	defer errorHandler()

	// seed the RNG with a somewhat random number. don't do crypto with this kids
	rand.Seed(time.Now().UTC().UnixNano())

	sgm := os.Getenv("SGM")
	token := os.Getenv("DISCORD_TOKEN")
	addr := os.Getenv("WEB_ADDR")
	dsn := os.Getenv("DSN")
	redisURI := os.Getenv("REDIS_URL")
	pgURI := os.Getenv("POSTGRES_URL")

	// Single-Guild-Mode hides all multiguild UI and pre-populates the selected guild
	config.SingleGuildMode = sgm

	// check our required set of flags
	if len(token) == 0 {
		panic(errors.New("discordgo token must be provided"))
	}

	if len(addr) == 0 {
		panic(errors.New("no idea where I should bind, please set WEB_ADDR"))
	}

	if len(dsn) != 0 {
		raven.SetDSN(dsn)
	}

	session := shodanSession{}
	session.moduleLoader = Loader

	var err error

	// set up external services
	session.cmdStack, err = InitCommandStack()
	panicOnError(err)

	session.kvs, err = InitRedis(redisURI)
	if err != nil {
		log.Printf("redis not initialized: %s - you can ignore this if you're running without a KVS\n", err)
	}

	session.database, err = InitPostgres(pgURI)
	if err != nil {
		log.Printf("pgsql not initialized: %s - you can ignore this if you're running without a DB\n", err)
	}

	session.discord, err = InitDiscord(token)
	panicOnError(err)

	session.mux, err = InitHTTP(addr)
	panicOnError(err)

	// raven should receive the git hash when crashing
	raven.SetRelease(config.VersionGitHash)

	err = session.Bootstrap()
	panicOnError(err)

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
				log.Printf("fatal thread error: %s\n", threadError.Error)
				os.Exit(0xDEAD)
			} else {
				log.Printf("non-fatal thread error: %s\n", threadError.Error)
			}
		}
		if quit {
			break
		}
	}
}
