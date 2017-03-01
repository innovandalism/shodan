package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/innovandalism/shodan/util"
	"io/ioutil"
	"net/http"
	"strings"
)

var VerifyJWTFunc func(string) (string, error) = func(t string) (string, error) {
	return "", errors.New("VerifyJWTFunc not overwritten, did you load mod_oauth?")
}

type RequestEnvelope struct {
	Data  []byte
	Token string
}

type ResponseEnvelope struct {
	Status int32       `json:"status"`
	Error  string      `json:"err"`
	Data   interface{} `json:"data"`
}

func ReadRequest(r *http.Request) (*RequestEnvelope, error) {
	req := RequestEnvelope{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = util.WrapError(err)
		return nil, err
	}
	req.Data = body
	if err != nil {
		err = util.WrapError(err)
		return nil, err
	}
	if authHeader:=r.Header.Get("Authorization"); authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return nil, errors.New("ReadRequest: Authentication header invalid")
		}
		jwt := parts[1]
		token, err := VerifyJWTFunc(jwt)
		if err != nil {
			return nil, util.WrapError(err)
		}
		req.Token = token
	}
	return &req, nil
}

func SendResponse(w http.ResponseWriter, res *ResponseEnvelope) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(res.Status))
	resBytes, err := json.Marshal(res)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "%s", resBytes)
	if err != nil {
		return err
	}
	return nil
}
