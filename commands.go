package shodan

import (
	"github.com/bwmarrin/discordgo"
	"github.com/innovandalism/shodan/util"
	"strings"
)

type Command interface {
	GetNames() []string
	Invoke(invocation *CommandInvocation) error
}

type CommandInvocation struct {
	Name      string
	Arguments []string
	Empty     bool
	Session   *discordgo.Session
	Event     *discordgo.MessageCreate
	Shodan    *Shodan
	Helpers   *CommandInvocationHelpers
}

type CommandInvocationHelpers struct {
	Reply func(string) (error)
	ReplyEmbed func(*discordgo.MessageEmbed) (error)
}

type CommandStack struct {
	commands        []Command
	FallbackCommand Command
}

type PermissionEnabledCommand interface{
	GetRequiredPermission() int
}

func (commandStack *CommandStack) Attach(shodan *Shodan) {
	callback := func(session *discordgo.Session, event *discordgo.MessageCreate) {
		if !util.MentionsMe(session.State.User.ID, event.Content) {
			return
		}

		commandInvocation := prepareCommand(event.Content, session, event)
		commandInvocation.Shodan = shodan
		commandStack.DispatchCommand(commandInvocation)
	}
	shodan.GetDiscord().AddHandler(callback)
}

func (commandStack *CommandStack) RegisterCommand(c Command) {
	commandStack.commands = append(commandStack.commands, c)
}

func (commandStack *CommandStack) DispatchCommand(ci *CommandInvocation) {
	for _, c := range commandStack.commands {
		for _, name := range c.GetNames() {
			if ci.Name != name {
				continue
			}
			ok, err := checkDiscordPermissions(ci, c)
			util.ReportThreadError(false, err)
			if !ok {
				_,err := ci.Session.ChannelMessageSend(ci.Event.ChannelID, "Permission denied.")
				util.ReportThreadError(false, err)
				return
			}
			c.Invoke(ci)
			return
		}
	}
	if commandStack.FallbackCommand != nil {
		commandStack.FallbackCommand.Invoke(ci)
	}
}

func checkDiscordPermissions(ci *CommandInvocation, c Command) (bool, error) {
	pec, ok := c.(PermissionEnabledCommand)
	if ok {
		perms, err := ci.Session.State.UserChannelPermissions(ci.Event.Author.ID, ci.Event.ChannelID)
		if err != nil {
			return false, err
		}
		return perms & pec.GetRequiredPermission() > 0, nil
	}
	return true, nil
}

func prepareCommand(message string, session *discordgo.Session, event *discordgo.MessageCreate) *CommandInvocation {
	parts := strings.Split(message, " ")

	// discard the mention
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
	attachHelpers(commandInvocation)
	return commandInvocation
}

func attachHelpers(ci *CommandInvocation) {
	ci.Helpers = &CommandInvocationHelpers{}
	ci.Helpers.Reply = func(message string) (error) {
		_, err := ci.Shodan.GetDiscord().ChannelMessageSend(ci.Event.ChannelID, message)
		return err
	}
	ci.Helpers.ReplyEmbed = func(embed *discordgo.MessageEmbed) (error) {
		_, err := ci.Shodan.GetDiscord().ChannelMessageSendEmbed(ci.Event.ChannelID, embed)
		return err
	}
}