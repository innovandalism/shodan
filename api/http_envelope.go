package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func ReadRequest(r *http.Request) *RequestEnvelope {
	req := RequestEnvelope{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}
	return &req
}

func SendResponse(w http.ResponseWriter, res *ResponseEnvelope) {
	w.WriteHeader(int(res.Status))
	resBytes, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "%s", resBytes)
}
