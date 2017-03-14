package test

import (
	"github.com/innovandalism/shodan"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/mux"
	"database/sql"
)

type MockShodan struct{
	FuncGetDiscord func() *discordgo.Session
	FuncGetMux func() *mux.Router
	FuncGetDatabase func() *sql.DB
	FuncGetCommandStack func() *shodan.CommandStack
	FuncGetRedis func() shodan.KVS
}

func (ms *MockShodan) GetDiscord() *discordgo.Session {
	return ms.FuncGetDiscord()
}

func (ms *MockShodan) GetDatabase() *sql.DB {
	return ms.FuncGetDatabase()
}

func (ms *MockShodan) GetMux() *mux.Router {
	return ms.FuncGetMux()
}

func (ms *MockShodan) GetCommandStack() *shodan.CommandStack {
	return ms.FuncGetCommandStack()
}

func (ms *MockShodan) GetRedis() shodan.KVS {
	return ms.FuncGetRedis()
}



