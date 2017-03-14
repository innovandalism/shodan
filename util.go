package shodan

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
)

func MentionsMe(userID string, content string) bool {
	mention := fmt.Sprintf("<@%s>", userID)
	mentionNicknamed := fmt.Sprintf("<@!%s>", userID)
	if !strings.HasPrefix(content, mention) && !strings.HasPrefix(content, mentionNicknamed) {
		return false
	}
	return true
}

func IsMention(needle string) bool {
	return strings.HasPrefix(needle, "<@") && strings.HasSuffix(needle, ">")
}

func makeSignalChannel() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	return c
}
