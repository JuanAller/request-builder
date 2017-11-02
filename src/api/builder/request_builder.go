package builder

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"encoding/base64"
	"net/http/httputil"
	"log"
	"encoding/xml"
)

type requestBuilder struct {
	client             HttpClient
	request            *Request
	responseType       string
	unmarshalFunctions map[string]func([]byte, interface{}) error
}

func (requestBuilder *requestBuilder) WithQueryParam(key string, value string) *requestBuilder {
	requestBuilder.request.QueryParams[key] = value
	return requestBuilder
}

func (requestBuilder *requestBuilder) WithHeader(key string, value string) *requestBuilder {
	requestBuilder.request.Headers[key] = value
	return requestBuilder
}

func (requestBuilder *requestBuilder) WithContentType(contentType string) *requestBuilder {
	requestBuilder.responseType = contentType
	return requestBuilder.WithHeader("Content-Type", contentType)
}

func (requestBuilder *requestBuilder) WithBody(body interface{}) *requestBuilder {
	requestBuilder.request.Body = body
	return requestBuilder
}

func (requestBuilder *requestBuilder) Accept(accept string) *requestBuilder {
	return requestBuilder.WithHeader("Accept", accept)
}

func (requestBuilder *requestBuilder) WithBasicAuthorization(username string, password string) *requestBuilder {
	return requestBuilder.WithHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(username+":"+password)))
}

func Get(client HttpClient, path string) *requestBuilder {
	return &requestBuilder{
		client:             client,
		request:            NewRequest(http.MethodGet, path),
		responseType:       "application/json",
		unmarshalFunctions: unmarshalFunctionsMap(),
	}
}

func Post(client HttpClient, path string) *requestBuilder {
	return &requestBuilder{
		client:             client,
		request:            NewRequest(http.MethodPost, path),
		responseType:       "application/json",
		unmarshalFunctions: unmarshalFunctionsMap(),
	}
}

func Put(client HttpClient, path string) *requestBuilder {
	return &requestBuilder{
		client:             client,
		request:            NewRequest(http.MethodPut, path),
		responseType:       "application/json",
		unmarshalFunctions: unmarshalFunctionsMap(),
	}
}

func Delete(client HttpClient, path string) *requestBuilder {
	return &requestBuilder{
		client:             client,
		request:            NewRequest(http.MethodDelete, path),
		responseType:       "application/json",
		unmarshalFunctions: unmarshalFunctionsMap(),
	}
}

func unmarshalFunctionsMap() map[string]func([]byte, interface{}) error {
	return map[string]func([]byte, interface{}) error{
		"application/json": json.Unmarshal,
		"application/xml":  xml.Unmarshal,
	}
}

func (requestBuilder *requestBuilder) Execute(entityResponse interface{}) *Response {
	response, err := requestBuilder.request.Execute(requestBuilder.client)
	if err != nil {
		return &Response{
			Error: err,
		}
	}
	rawResp, _ := httputil.DumpResponse(response, true)
	log.Println(string(rawResp))
	defer response.Body.Close()
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		body, _ := ioutil.ReadAll(response.Body)
		return &Response{
			StatusCode: response.StatusCode,
			Error:      requestBuilder.unmarshalFunctions[requestBuilder.responseType](body, entityResponse),
		}
	}
	return &Response{
		StatusCode: response.StatusCode,
		Error:      nil,
	}
}
