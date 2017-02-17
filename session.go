package shodan

import (
	"fmt"
	Discord "github.com/bwmarrin/discordgo"
	"github.com/innovandalism/shodan/util"
	"net/http"
)

type Shodan struct {
	discord      *Discord.Session
	moduleLoader *ModuleLoader
	mux          *http.ServeMux
	cmdStack     *CommandStack
}

func (session *Shodan) GetDiscord() *Discord.Session {
	return session.discord
}

func (session *Shodan) GetMux() *http.ServeMux {
	return session.mux
}

func (session *Shodan) GetCommandStack() *CommandStack {
	return session.cmdStack
}

func Init() *Shodan {
	session := Shodan{}
	session.moduleLoader = Loader
	session.mux = http.NewServeMux()
	session.cmdStack = &CommandStack{}
	return &session
}

func (session *Shodan) FlagHook() {
	session.moduleLoader.FlagHook()
}

func (session *Shodan) InitDiscord(token *string) error {
	t := fmt.Sprintf("Bot %s", *token)
	s, err := Discord.New(t)

	if err != nil {
		return err
	}

	s.State.User, err = s.User("@me")

	if err != nil {
		return err
	}

	session.discord = s

	session.cmdStack.Attach(session)
	session.moduleLoader.Attach(session)

	return nil
}

func (session *Shodan) InitHttp(addr string) {
	err := http.ListenAndServe(addr, session.mux)
	if err != nil {
		util.ReportThreadError(true, err)
	}
}

func (session *Shodan) Connect() {
	session.discord.Open()
}