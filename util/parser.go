package util

import (
	"fmt"
	"strings"
)

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

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
