package builder

import (
	"net/http"
	"encoding/json"
	"bytes"
	"net/http/httputil"
	"log"
	"encoding/xml"
)

type request struct {
	Method         string
	Path           string
	Headers        map[string]string
	QueryParams    map[string]string
	Body           interface{}
	MarshalFuncs   map[string]func(v interface{}) ([]byte, error)
	ContentType    string
	logRequestBody bool
}

func newRequest(method string, path string) *request {
	return &request{
		Method:      method,
		Path:        path,
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
		MarshalFuncs: map[string]func(v interface{}) ([]byte, error){
			APPLICATIONJSON: json.Marshal,
			APPLICATIONXML:  xml.Marshal,
		},
		ContentType: APPLICATIONJSON,
	}
}

func (request *request) build() (*http.Request, error) {
	byteSlice, marshallErr := request.MarshalFuncs[request.ContentType](request.Body)
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
	if request.logRequestBody {
		rawRequest, _ := httputil.DumpRequestOut(newRequest, request.logRequestBody)
		log.Println(string(rawRequest))
	}
	return newRequest, nil
}
