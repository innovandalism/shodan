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
		err = util.WrapError(err)
		return err
	}
	return nil
}

// 400
func ErrorBadRequest(w http.ResponseWriter, err error) error {
	return Error(w, 400, err)
}

// 403
func ErrorUnauthorized(w http.ResponseWriter, err error) error {
	return Error(w, 403, err)
}

// 500
func ErrorInternalServerError(w http.ResponseWriter, err error) error {
	return Error(w, 500, err)
}