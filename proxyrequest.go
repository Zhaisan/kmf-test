package main

import (
	"encoding/json"
	"net/http"
)

type ProxyRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

func NewProxyRequest(r *http.Request) (*ProxyRequest, error) {
	var proxyRequest ProxyRequest

	if r.Method != "POST" {
		return nil, NewHTTPError("method not allowed", http.StatusMethodNotAllowed)
	}

	if err := json.NewDecoder(r.Body).Decode(&proxyRequest); err != nil {
		return nil, NewHTTPError("bad request", http.StatusBadRequest)
	}

	return &proxyRequest, nil
}
