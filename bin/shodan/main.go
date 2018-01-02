package main

import (
	"github.com/innovandalism/shodan"

	// module imports
	_ "github.com/innovandalism/shodan/modules/core"
	_ "github.com/innovandalism/shodan/modules/gmtools"
	_ "github.com/innovandalism/shodan/modules/greet"
	_ "github.com/innovandalism/shodan/modules/log"
	_ "github.com/innovandalism/shodan/modules/oauth"
	_ "github.com/innovandalism/shodan/modules/webui"
)

func main() {
	shodan.Run()
}
