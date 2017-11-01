package builder

import (
	"net/http"
	"encoding/json"
	"bytes"
	"net/http/httputil"
	"log"
)

type Request struct {
	Method      string
	Path        string
	Headers     map[string]string
	QueryParams map[string]string
	Body        interface{}
}

func NewRequest(method string, path string) *Request {
	return &Request{
		Method:      method,
		Path:        path,
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
	}
}

func (simpleRequest *Request) Execute(client HttpClient) (*http.Response, error) {
	byteSlice, marshallErr := json.Marshal(simpleRequest.Body)
	if marshallErr != nil {
		return nil, marshallErr
	}
	request, _ := http.NewRequest(simpleRequest.Method, simpleRequest.Path, bytes.NewBuffer(byteSlice))
	query := request.URL.Query()
	for key, value := range simpleRequest.QueryParams {
		query.Add(key, value)
	}
	request.URL.RawQuery = query.Encode()
	for key, value := range simpleRequest.Headers {
		request.Header.Set(key, value)
	}
	rawRequest, _ := httputil.DumpRequestOut(request, true)
	log.Println(string(rawRequest))
	return client.Do(request)
}
