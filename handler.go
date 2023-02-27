package main

import (
	"log"
	"net/http"
)

var requests = make(map[int]ProxyRequest)
var responses = make(map[int]ProxyResponse)

func handler(w http.ResponseWriter, r *http.Request) {
	proxyRequest, err := NewProxyRequest(r)
	if err != nil {
		log.Println(err)
		if httpError, ok := err.(*HTTPError); ok {
			writeJSONResponse(w, map[string]string{"error": httpError.Message}, httpError.StatusCode)
		} else {
			writeJSONResponse(w, map[string]string{"error": "Internal server error"}, http.StatusInternalServerError)
		}
		return
	}



	proxyResponse, err := DoProxyRequest(proxyRequest)
	if err != nil {
		log.Println(err)
		writeJSONResponse(w, map[string]string{"error": "Internal server error"}, http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, proxyResponse, http.StatusOK)
}
