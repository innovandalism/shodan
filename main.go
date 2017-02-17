package shodan

import (
	"errors"
	"flag"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/innovandalism/shodan/config"
	"github.com/innovandalism/shodan/util"
	"github.com/vharitonsky/iniflags"
	"math/rand"
	"os"
	"time"
)

// Invokes the main SHODAN loop, starting HTTP listeners and Discordgo.
// Meant to be called from an external package that imports your modules
func Run() {
	// catch-all error handler
	defer util.ErrorHandler()

	// seed the RNG with a somewhat random number. don't do crypto with this kids
	rand.Seed(time.Now().UTC().UnixNano())

	session := Init()

	// essential flags
	var (
		token = flag.String("token", config.DefaultToken, "Discord Bot Authentication Token")
		addr  = flag.String("web_addr", "", "Address to bind HTTP listener to")
		noweb = flag.Bool("web_disable", false, "Disable WebUI")
		dsn   = flag.String("dsn", "", "Sentry DSN")
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

	raven.SetRelease(config.VersionGitHash)

	err := session.InitDiscord(token)
	util.PanicOnError(err)

	if !*noweb {
		go session.InitHttp(*addr)
	}

	session.Connect()
	defer session.GetDiscord().Close()

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
