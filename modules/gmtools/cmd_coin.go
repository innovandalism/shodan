package gmtools

import (
	"math/rand"

	"github.com/innovandalism/shodan"
)

// CoinCommand is an empty struct that holds this commands methods
type CoinCommand struct{}

// GetNames returns the command aliases for this command
func (*CoinCommand) GetNames() []string {
	return []string{"coin"}
}

// Invoke runs the command. This command will block a user. Currently inoperable.
func (command *CoinCommand) Invoke(ci *shodan.CommandInvocation) error {
	message := "Heads"
	toss := rand.Intn(2)
	if toss == 1 {
		message = "Tails"
	}
	return ci.Helpers.Reply(message)
}
