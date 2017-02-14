package util

import (
	"os"
	"os/signal"
)

func MakeSignalChannel() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	return c
}
