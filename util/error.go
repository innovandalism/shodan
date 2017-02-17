package util

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/innovandalism/shodan/config"
	"os"
	"runtime"
)

var errorChannel chan *ThreadError

type ThreadError struct {
	IsFatal bool
	Error   error
}

func errorMessage(err error) {
	fmt.Printf("Shodan %d.%d.%d (%s) has crashed\n", config.VersionMajor, config.VersionMinor, config.VersionRevision, config.VersionGitHash)
	fmt.Println("")

	// check if we're given a proper error
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		fmt.Println("NO ERROR MESSAGE HAS BEEN PROVIDED")
		fmt.Println("This is a bug, please report the stack trace below")
		fmt.Println("https://github.com/innovandalism/shodan")
	}

	fmt.Println("")

	// try to get a stack trace if this is a debug build
	if config.Debug {
		buf := make([]byte, 1024)
		runtime.Stack(buf, false)
		fmt.Printf("%s\n", buf)
	}
}

// Error handler. Only call from main thread.
func ErrorHandler() {
	if r := recover(); r != nil {
		err, ok := r.(error)
		if ok {
			errorMessage(err)
		} else {
			errorMessage(nil)
		}
		os.Exit(0xDEAD)
	}
}

func ReportThreadError(isFatal bool, error error) {
	if error == nil {
		return
	}

	stack := raven.NewStacktrace(1, 2, nil)
	packet := raven.NewPacket(error.Error(), raven.NewException(error, stack))
	raven.DefaultClient.Capture(packet, nil)

	if isFatal {
		errorMessage(error)
	}

	GetThreadErrorChannel() <- &ThreadError{
		IsFatal: isFatal,
		Error:   error,
	}
	if isFatal {
		runtime.Goexit()
	}
}

func GetThreadErrorChannel() chan *ThreadError {
	if errorChannel == nil {
		errorChannel = make(chan *ThreadError, 10)
	}
	return errorChannel
}