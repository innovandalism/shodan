package shodan

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var VerifyJWTFunc func(string) (string, error) = func(t string) (string, error) {
	return "", WrapErrorHttp(errors.New("VerifyJWTFunc not overwritten, did you load mod_oauth?"), 500)
}

// RequestEnvelope wraps a HTTP request to the SHODAN-API
type RequestEnvelope struct {
	Data  []byte
	Token string
}

// ResponseEnvelope wraps a HTTP respone from the SHODAN-API
type ResponseEnvelope struct {
	Status int32       `json:"status"`
	Error  string      `json:"err"`
	Data   interface{} `json:"data"`
}

// Checks if a token is provided. Because ReadRequest will fail if a token is provided, but doesn't validate, no further checking is required; this is a cheap call.
func (req *RequestEnvelope) Authenticated() bool {
	return !(len(req.Token) < 1)
}

// ReadRequest reads the given Request and converts it into a standard API request
func ReadRequest(r *http.Request) (*RequestEnvelope, error) {
	req := RequestEnvelope{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = WrapError(err)
		return nil, err
	}
	req.Data = body
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return nil, WrapErrorHttp(errors.New("ReadRequest: Authentication header invalid"), 400)
		}
		jwt := parts[1]
		token, err := VerifyJWTFunc(jwt)
		if err != nil {
			return nil, err
		}
		req.Token = token
	}
	return &req, nil
}

// SendResponse sends the given ResponseEnvelope to the ResponseWriter
func SendResponse(w http.ResponseWriter, res *ResponseEnvelope) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(res.Status))
	resBytes, err := json.Marshal(res)
	if err != nil {
		return WrapError(err)
	}
	_, err = fmt.Fprintf(w, "%s", resBytes)
	if err != nil {
		return WrapError(err)
	}
	return nil
}

// HttpForward writes a 301 forward to the supplied uri via a ResponseWriter
func HttpForward(w http.ResponseWriter, uri string) {
	w.Header().Add("Location", uri)
	w.WriteHeader(301)
}
