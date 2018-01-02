package shodan

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// The Command interface enables the implementer to be loaded as a command
type Command interface {
	GetNames() []string
	Invoke(invocation *CommandInvocation) error
}

// A CommandInvocation contains all data (and references) the command stack parses from an incoming command
type CommandInvocation struct {
	Name      string
	Arguments []string
	Empty     bool
	Event     *discordgo.Message
	Shodan    Shodan
	Helpers   *CommandInvocationHelpers
}

// CommandInvocationHelpers are attached to the CommandInvocation and provide easy access to frequently used functionality
type CommandInvocationHelpers struct {
	Reply      func(string) error
	ReplyEmbed func(*discordgo.MessageEmbed) error
}

// The CommandStack holds an array of commands and an optional fallback command
type CommandStack struct {
	commands        []Command
	FallbackCommand Command
}

// PermissionEnabledCommand defines an interface Commands can implement to indicate the
// permission bits corresponding to discordgo ACLs required to run this command
type PermissionEnabledCommand interface {
	GetRequiredPermission() int
}

// Attach the command stack to a shodan instance
func (commandStack *CommandStack) Attach(shodan Shodan) {
	callback := func(session *discordgo.Session, event *discordgo.MessageCreate) {
		if !MentionsMe(session.State.User.ID, event.Content) {
			return
		}

		commandInvocation := prepareCommand(event.Content, event)
		commandInvocation.Shodan = shodan
		err := commandStack.dispatchCommand(commandInvocation)
		if err != nil {
			ReportThreadError(false, err)
		}
	}
	shodan.GetDiscord().AddHandler(callback)
}

// RegisterCommand adds a command to the stack
func (commandStack *CommandStack) RegisterCommand(c Command) {
	log.Printf("CommandStack: Registering command %s\n", strings.Join(c.GetNames(), ","))
	commandStack.commands = append(commandStack.commands, c)
}

// GetCommandList returns a map of all commands.
// The includeAliases variable determines if a command can have multiple entries, one for each alias.
// Otherwise only the first alias is returned
func (commandStack *CommandStack) GetCommandList(includeAliases bool) map[string]Command {
	res := make(map[string]Command)
	for _, c := range commandStack.commands {
		names := c.GetNames()
		if includeAliases {
			for _, name := range names {
				res[name] = c
			}
		} else {
			// if someone is ever so smart to register a command with no names...
			if len(names) < 1 {
				continue
			}
			res[names[0]] = c
		}
	}
	return res
}

// Checks for permissions, matches the invocation to a command and invokes the command handler
func (commandStack *CommandStack) dispatchCommand(ci *CommandInvocation) error {
	var command Command
	for _, c := range commandStack.commands {
		for _, name := range c.GetNames() {
			if ci.Name != name {
				continue
			}
			ok, err := checkDiscordPermissions(ci, c)
			if err != nil {
				return err
			}
			if !ok {
				err := ci.Helpers.Reply("Permission denied.")
				if err != nil {
					return err
				}
				return nil
			}
			command = c
			break
		}
	}
	if commandStack.FallbackCommand != nil {
		command = commandStack.FallbackCommand
	}
	if command != nil {
		return command.Invoke(ci)
	}
	return nil
}

// shorthand for checking if a CommandInvocation can invoke a command
func checkDiscordPermissions(ci *CommandInvocation, c Command) (bool, error) {
	pec, ok := c.(PermissionEnabledCommand)
	if ok {
		perms, err := ci.Shodan.GetDiscord().State.UserChannelPermissions(ci.Event.Author.ID, ci.Event.ChannelID)
		if err != nil {
			return false, WrapError(err)
		}
		return perms&pec.GetRequiredPermission() > 0, nil
	}
	return true, nil
}

// parser logic
func prepareCommand(message string, event *discordgo.MessageCreate) *CommandInvocation {
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
	commandInvocation.Event = event.Message
	attachHelpers(commandInvocation)
	return commandInvocation
}

// attach various helpers
func attachHelpers(ci *CommandInvocation) {
	ci.Helpers = &CommandInvocationHelpers{}
	ci.Helpers.Reply = func(message string) error {
		_, err := ci.Shodan.GetDiscord().ChannelMessageSend(ci.Event.ChannelID, message)
		if err != nil {
			err = WrapError(err)
		}
		return err
	}
	ci.Helpers.ReplyEmbed = func(embed *discordgo.MessageEmbed) error {
		_, err := ci.Shodan.GetDiscord().ChannelMessageSendEmbed(ci.Event.ChannelID, embed)
		if err != nil {
			err = WrapError(err)
		}
		return err
	}
}
