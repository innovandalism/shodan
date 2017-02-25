package api

import "net/http"

func Forward(w http.ResponseWriter, uri string) {
	w.Header().Add("Location", uri)
	w.WriteHeader(301)
}