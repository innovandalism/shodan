package shodan

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"

	raven "github.com/getsentry/raven-go"

	"github.com/innovandalism/shodan/config"
)

var errorChannel chan *ThreadError

// A ThreadError wraps an error that occured in another thread, indicating if the error was fatal and the application should be terminated
type ThreadError struct {
	IsFatal bool
	Error   error
}

// DebuggableError implements Error and stores a HTTP status code and stack traces in addition to an error message.
type DebuggableError struct {
	error  string
	Status int
	stack  *raven.Stacktrace
	packet *raven.Packet
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// Error returns the error message
func (e *DebuggableError) Error() string {
	return e.error
}

// Capture sends this error to Sentry for further analysis
func (e *DebuggableError) Capture() {
	raven.DefaultClient.Capture(e.packet, nil)
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

	// TODO: Check if debuggable
	buf := make([]byte, 1024)
	runtime.Stack(buf, false)
	fmt.Printf("%s\n", buf)
}

// Error handler. Only call from main thread.
func errorHandler() {
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

// ReportThreadError reports the current thread as failed, captures the error if it is debuggable and kills the goroutine
func ReportThreadError(isFatal bool, error error) {
	if error == nil {
		return
	}

	de, isDebuggable := error.(*DebuggableError)
	if isDebuggable {
		de.Capture()
	}

	if isFatal {
		errorMessage(error)
	}

	getThreadErrorChannel() <- &ThreadError{
		IsFatal: isFatal,
		Error:   error,
	}

	// goodbye cruel world
	runtime.Goexit()
}

// GetThreadErrorChannel returns the channel errors should be sent to
func getThreadErrorChannel() chan *ThreadError {
	if errorChannel == nil {
		errorChannel = make(chan *ThreadError, 10)
	}
	return errorChannel
}

// Error is a shorthand function to create a DebuggableError from a string
func Error(msg string) error {
	return WrapError(errors.New(msg))
}

// ErrorHttp is a shorthand function to create a DebuggableError from a string and HTTP status code
func ErrorHttp(msg string, code int) error {
	return WrapErrorHttp(Error(msg), code)
}

// WrapError wraps any error into a DebuggableError. Never double-wraps.
// This function should only be used at module boundaries for performance reasons.
// Sets the HTTP status to 500 by default.
func WrapError(e error) error {
	var de = &DebuggableError{}
	if e == nil {
		panic("Bug: WrapError called with nil. This should never happen.")
	}

	// make sure we don't double-wrap the error
	_, isDebuggableError := e.(*DebuggableError)
	if isDebuggableError {
		return e
	}

	// attach the actual debug information
	de.Status = 500
	de.error = e.Error()
	de.stack = raven.NewStacktrace(1, 2, nil)
	de.packet = raven.NewPacket(e.Error(), raven.NewException(e, de.stack))
	return de
}

// WrapErrorHttp wraps an existing error into a DebuggableError and attaches a status code to it
func WrapErrorHttp(e error, code int) error {
	e = WrapError(e)
	de, _ := e.(*DebuggableError)
	de.Status = code
	return de
}

// HttpSendError sends a response to the given ResponseWriter and captures the error if given a DebuggableError
//
// This function sends 500 Internal Server Error by default
func HttpSendError(w http.ResponseWriter, err error) error {
	status := 500
	de, isDe := err.(*DebuggableError)
	if isDe {
		status = de.Status
		de.Capture()
	}
	res := ResponseEnvelope{
		Status: int32(status),
		Error:  fmt.Sprintf("%s", err),
	}
	err = SendResponse(w, &res)
	return err
}
