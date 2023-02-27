package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)


type HTTPError struct {
	StatusCode int
	Message    string
}


func (e *HTTPError) Error() string {
	return e.Message
}

func NewHTTPError(message string, statusCode int) error {
	return &HTTPError{
		StatusCode: statusCode,
		Message:    message,
	}
}


func DoProxyRequest(proxyRequest *ProxyRequest) (*ProxyResponse, error) {
	client := &http.Client{}
	request, err := http.NewRequest(proxyRequest.Method, proxyRequest.URL, nil)
	if err != nil {
		return nil, errors.New("bad request")
	}

	for k, v := range proxyRequest.Headers {
		request.Header.Set(k, v)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, errors.New("error connecting to service")
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("error reading response body")
	}
	headers := make(map[string]string)
	for k, v := range response.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	proxyResponse := &ProxyResponse{
		ID:      len(requests) + 1,
		Status:  response.StatusCode,
		Headers: headers,
		Length:  int64(len(responseBody)),
	}

	requests[proxyResponse.ID] = *proxyRequest
	responses[proxyResponse.ID] = *proxyResponse

	return proxyResponse, nil
}

func writeJSONResponse(w http.ResponseWriter, response interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	if err := encoder.Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(buffer.Bytes()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
