package builder

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"encoding/base64"
	"net/http/httputil"
	"log"
)

type requestBuilder struct {
	Client  HttpClient
	Request *Request
}

func (requestBuilder *requestBuilder) WithQueryParam(key string, value string) *requestBuilder {
	requestBuilder.Request.QueryParams[key] = value
	return requestBuilder
}

func (requestBuilder *requestBuilder) WithHeader(key string, value string) *requestBuilder {
	requestBuilder.Request.Headers[key] = value
	return requestBuilder
}

func (requestBuilder *requestBuilder) WithContentType(contentType string) *requestBuilder {
	return requestBuilder.WithHeader("Content-Type", contentType)
}

func (requestBuilder *requestBuilder) WithBody(body interface{}) *requestBuilder {
	requestBuilder.Request.Body = body
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
		Client:  client,
		Request: NewRequest(http.MethodGet, path),
	}
}

func Post(client HttpClient, path string) *requestBuilder {
	return &requestBuilder{
		Client:  client,
		Request: NewRequest(http.MethodPost, path),
	}
}

func Put(client HttpClient, path string) *requestBuilder {
	return &requestBuilder{
		Client:  client,
		Request: NewRequest(http.MethodPut, path),
	}
}

func Delete(client HttpClient, path string) *requestBuilder {
	return &requestBuilder{
		Client:  client,
		Request: NewRequest(http.MethodDelete, path),
	}
}

func (requestBuilder *requestBuilder) Execute(entityResponse interface{}) *Response {
	response, err := requestBuilder.Request.Execute(requestBuilder.Client)
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
			Error:      json.Unmarshal(body, entityResponse),
		}
	}
	return &Response{
		StatusCode: response.StatusCode,
		Error:      nil,
	}
}
