package version

import (
	"fmt"
	"github.com/innovandalism/shodan"
	"github.com/innovandalism/shodan/util"
)

type IgnoreCommand struct{}

func (_ *IgnoreCommand) GetNames() []string {
	return []string{"ignore"}
}

func (command *IgnoreCommand) Invoke(ci *shodan.CommandInvocation) bool {
	if len(ci.Arguments) != 1 {
		err := ci.Helpers.Reply("This command requires one argument")
		util.ReportThreadError(false, err)
		return false
	}
	if !util.IsMention(ci.Arguments[0]) {
		err := ci.Helpers.Reply("Argument must be a mention")
		util.ReportThreadError(false, err)
		return false
	}
	// FIXME: This needs to go into redis or something
	err := ci.Session.RelationshipUserBlock(ci.Event.Mentions[0].ID)
	util.ReportThreadError(false, err)
	if err == nil {
		err := ci.Helpers.Reply("Command successful.")
		util.ReportThreadError(false, err)
		return true
	} else {
		err := ci.Helpers.Reply(fmt.Sprintf("Command failed: %s", err))
		util.ReportThreadError(false, err)
		return false
	}
}
