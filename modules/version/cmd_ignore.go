package version

import (
	"fmt"
	"github.com/innovandalism/shodan"
)

// IgnoreCommand is an empty struct that holds this commands methods
type IgnoreCommand struct{}

// GetNames returns the command aliases for this command
func (*IgnoreCommand) GetNames() []string {
	return []string{"ignore"}
}

// Invoke runs the command. This command will block a user. Currently inoperable.
func (command *IgnoreCommand) Invoke(ci *shodan.CommandInvocation) error {
	if len(ci.Arguments) != 1 {
		err := ci.Helpers.Reply("This command requires one argument")
		if err != nil {
			return shodan.WrapError(err)
		}
		return nil
	}
	if !shodan.IsMention(ci.Arguments[0]) {
		err := ci.Helpers.Reply("Argument must be a mention")
		if err != nil {
			return shodan.WrapError(err)
		}
		return nil
	}
	// FIXME: This needs to go into redis or something
	err := ci.Shodan.GetDiscord().RelationshipUserBlock(ci.Event.Mentions[0].ID)
	shodan.ReportThreadError(false, err)
	if err == nil {
		err := ci.Helpers.Reply("Command successful.")
		if err != nil {
			return shodan.WrapError(err)
		}
		return nil
	} else {
		err := ci.Helpers.Reply(fmt.Sprintf("Command failed: %s", err))
		if err != nil {
			return shodan.WrapError(err)
		}
		return nil
	}
}
