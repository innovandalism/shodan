package log

import (
	"database/sql"
	"github.com/bwmarrin/discordgo"
	"github.com/innovandalism/shodan/test"
	_ "github.com/lib/pq"
	"os"
	"testing"
)

// TestCheckInserts is broken. Don't call it.
func TestCheckInserts(t *testing.T) {
	// goodbye for now!
	t.FailNow()
	ms := test.MockShodan{
		FuncGetDatabase: func() *sql.DB {
			db, err := sql.Open("postgres", os.Getenv("SHODAN_MOCK_POSTGRES"))
			if err != nil {
				panic(err)
			}
			// Connectivity/Auth test
			err = db.Ping()
			if err != nil {
				panic(err)
			}
			return db
		},
	}
	mod.shodan = &ms
	mockSession := discordgo.Session{}
	mockSession.State = discordgo.NewState()

	mockSession.State.GuildAdd(&discordgo.Guild{
		ID: "1",
	})

	mockSession.State.ChannelAdd(&discordgo.Channel{
		ID:      "1",
		Name:    "general",
		GuildID: "1",
	})

	mockMessage := discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author: &discordgo.User{
				ID:            "1",
				Username:      "john doe",
				Discriminator: "0001",
			},
			ChannelID: "1",
			ID:        "1",
			Content:   "Hello World",
		},
	}

	onMessageCreate(&mockSession, &mockMessage)
}
