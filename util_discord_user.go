package shodan

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func dUDo(token string, cb func(session *discordgo.Session) error) error {
	t := fmt.Sprintf("Bearer %s", token)
	s, err := discordgo.New(t)
	defer s.Close()
	if err != nil {
		err = WrapError(err)
		return err
	}
	err = cb(s)
	return err
}

// DUFetchProfile fetches a users own profile using the provided token (DU: Do as User)
func DUFetchProfile(token string) (*discordgo.User, error) {
	var u *discordgo.User
	var err error
	err = dUDo(token, func(s *discordgo.Session) error {

		u, err = s.User("@me")
		if err != nil {
			err = WrapError(err)
			return err
		}
		return nil
	})
	return u, err
}

// DUGetUserGuild gets a users joined guilds using the provided token (DU: Do as User)
func DUGetUserGuilds(token string) ([]*discordgo.UserGuild, error) {
	var u []*discordgo.UserGuild
	var err error
	err = dUDo(token, func(s *discordgo.Session) error {

		u, err = s.UserGuilds(100, "", "")
		if err != nil {
			err = WrapError(err)
			return err
		}
		return nil
	})
	return u, err
}
