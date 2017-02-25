package api

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/innovandalism/shodan/util"
)

func FetchProfile(token string) (*discordgo.User, error) {
	t := fmt.Sprintf("Bearer %s", token)
	s, err := discordgo.New(t)
	defer s.Close()
	if err != nil {
		err = util.WrapError(err)
		return nil, err
	}
	u, err := s.User("@me")
	if err != nil {
		err = util.WrapError(err)
		return nil, err
	}
	return u, nil
}
