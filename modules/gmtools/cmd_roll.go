package gmtools

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"math/rand"

	"github.com/innovandalism/shodan"
)

// RollCommand is an empty struct that holds this commands methods
type RollCommand struct{}

// GetNames returns the command aliases for this command
func (*RollCommand) GetNames() []string {
	return []string{"roll"}
}

// Invoke runs the command. This command will block a user. Currently inoperable.
func (command *RollCommand) Invoke(ci *shodan.CommandInvocation) error {
	if len(ci.Arguments) != 1 {
		err := ci.Helpers.Reply("This command requires one argument")
		if err != nil {
			return shodan.WrapError(err)
		}
		return nil
	}
	rollexp, err := regexp.Compile("([0-9]*)d([0-9]+)")
	if err != nil {
		return shodan.WrapError(err)
	}
	strgs := rollexp.FindStringSubmatch(ci.Arguments[0])
	var iterations, sides float64

	if strgs[1] != "" {
		iterations, err = strconv.ParseFloat(strgs[1], 64)
		if err != nil {
			return shodan.WrapError(err)
		}
		iterations = math.Max(1, math.Min(iterations, 64))
	} else {
		iterations = 1
	}

	sides, err = strconv.ParseFloat(strgs[2], 64)
	if err != nil {
		return shodan.WrapError(err)
	}
	sides = math.Max(1, math.Min(sides, 1024))
	err = ci.Helpers.Reply(fmt.Sprintf("Rolling %dd%d...", int64(iterations), int64(sides)))

	numbers := make([]string, int64(iterations))
	var sum int64

	for i := 0; i < int(iterations); i++ {
		num := rand.Int63n(int64(sides)) + 1
		numbers[i] = strconv.Itoa(int(num))
		sum += num
	}

	err = ci.Helpers.Reply(fmt.Sprintf("```\n\n%s\n=======\n%d```", strings.Join(numbers, "\n"), sum))

	return err
}
