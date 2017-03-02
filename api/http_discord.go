package api

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/innovandalism/shodan/util"
)

func DoAsUser(token string, cb func(session *discordgo.Session) error) error {
	t := fmt.Sprintf("Bearer %s", token)
	s, err := discordgo.New(t)
	defer s.Close()
	if err != nil {
		err = util.WrapError(err)
		return err
	}
	err = cb(s)
	return err
}

func FetchProfile(token string) (*discordgo.User, error) {
	var u *discordgo.User
	var err error
	err = DoAsUser(token, func(s *discordgo.Session) error {

		u, err = s.User("@me")
		if err != nil {
			err = util.WrapError(err)
			return err
		}
		return nil
	})
	return u, err
}


func GetUserGuilds(token string) ([]*discordgo.UserGuild, error) {
	var u []*discordgo.UserGuild
	var err error
	err = DoAsUser(token, func(s *discordgo.Session) error {

		u, err = s.UserGuilds()
		if err != nil {
			err = util.WrapError(err)
			return err
		}
		return nil
	})
	return u, err
}