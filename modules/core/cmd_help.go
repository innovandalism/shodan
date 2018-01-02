package version

import (
	"github.com/innovandalism/shodan"
)

// HelpCommand is an empty struct that holds this commands methods
type HelpCommand struct{}

// GetNames returns the command aliases for this command
func (*HelpCommand) GetNames() []string {
	return []string{"help"}
}

// Invoke runs the command. This command will block a user. Currently inoperable.
func (command *HelpCommand) Invoke(ci *shodan.CommandInvocation) error {
	return nil
}
