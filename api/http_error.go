package api

import "fmt"
import (
	"net/http"
	"github.com/innovandalism/shodan/util"
)

func Error(w http.ResponseWriter, code int, err error) error {
	res := ResponseEnvelope{
		Status: int32(code),
		Error:  fmt.Sprintf("%s", err),
	}
	err = SendResponse(w, &res)
	if err != nil {
		err = err
		return err
	}
	return nil
}

func ErrorInfer(w http.ResponseWriter, err error) error {
	code := 500
	de, isDebuggable := err.(*util.DebuggableError)
	if isDebuggable {
		code = de.Status
	}
	return Error(w, code, err)
}