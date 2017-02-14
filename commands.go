package shodan

import (
	"github.com/bwmarrin/discordgo"
	"github.com/innovandalism/shodan/util"
	"strings"
)

type Command interface {
	GetNames() []string
	Invoke(invocation *CommandInvocation) bool
}

type CommandInvocation struct {
	Name      string
	Arguments []string
	Empty     bool
	Session   *discordgo.Session
	Event     *discordgo.MessageCreate
}

type CommandStack struct {
	commands        []Command
	FallbackCommand Command
}

func (commandStack *CommandStack) Attach(session *Shodan) {
	callback := func(session *discordgo.Session, event *discordgo.MessageCreate) {
		if !util.MentionsMe(session.State.User.ID, event.Content) {
			return
		}
		commandInvocation := prepareCommand(event.Content, session, event)
		commandStack.DispatchCommand(commandInvocation)
	}
	session.GetDiscord().AddHandler(callback)
}

func (commandStack *CommandStack) RegisterCommand(c Command) {
	commandStack.commands = append(commandStack.commands, c)
}

func (commandStack *CommandStack) DispatchCommand(ci *CommandInvocation) {
	for _, c := range commandStack.commands {
		for _, name := range c.GetNames() {
			if ci.Name == name {
				c.Invoke(ci)
				return
			}
		}

	}
	if commandStack.FallbackCommand != nil {
		commandStack.FallbackCommand.Invoke(ci)
	}
}

func prepareCommand(message string, session *discordgo.Session, event *discordgo.MessageCreate) *CommandInvocation {
	parts := strings.Split(message, " ")

	// discord the mention
	parts = parts[1:]

	commandInvocation := func() *CommandInvocation {
		// bail if there's only one part
		if len(parts) < 1 {
			return &CommandInvocation{
				Empty: true,
			}
		}

		commandName := parts[0]

		return &CommandInvocation{
			Empty: false,
			Name:  commandName,
		}
	}()

	if len(parts) > 1 {
		commandInvocation.Arguments = parts[1:]
	}

	commandInvocation.Session = session
	commandInvocation.Event = event
	return commandInvocation
}
