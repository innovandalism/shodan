package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/innovandalism/shodan/util"
)

type RequestEnvelope struct {
	Token string      `json:"token"`
	Data  interface{} `json:"data"`
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
	err = json.Unmarshal(body, &req)
	if err != nil {
		err = util.WrapError(err)
		return nil, err
	}
	return &req, nil
}

func SendResponse(w http.ResponseWriter, res *ResponseEnvelope) (error) {
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
