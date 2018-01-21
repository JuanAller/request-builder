package builder

import (
	"net/http"
	"encoding/json"
	"bytes"
	"net/http/httputil"
	"log"
)

type request struct {
	Method         string
	Path           string
	Headers        map[string]string
	QueryParams    map[string]string
	Body           interface{}
	logRequestBody bool
}

func newRequest(method string, path string) *request {
	return &request{
		Method:      method,
		Path:        path,
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
	}
}

func (request *request) execute(client HttpClient) (*http.Response, error) {
	byteSlice, marshallErr := json.Marshal(request.Body)
	if marshallErr != nil {
		return nil, marshallErr
	}
	newRequest, _ := http.NewRequest(request.Method, request.Path, bytes.NewBuffer(byteSlice))
	query := newRequest.URL.Query()
	for key, value := range request.QueryParams {
		query.Add(key, value)
	}
	newRequest.URL.RawQuery = query.Encode()
	for key, value := range request.Headers {
		newRequest.Header.Set(key, value)
	}
	rawRequest, _ := httputil.DumpRequestOut(newRequest, request.logRequestBody)
	log.Println(string(rawRequest))
	return client.Do(newRequest)
}
