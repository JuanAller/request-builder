package builder

import (
	"io/ioutil"
	"encoding/base64"
	"net/http/httputil"
	"log"
)

const (
	APPLICATIONJSON = "application/json"
	APPLICATIONXML  = "application/xml"
)

type requestBuilder struct {
	client             HttpClient
	request            *request
	contentType        string
	unmarshalFunctions map[string]func([]byte, interface{}) error
	logResponseBody    bool
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
	requestBuilder.contentType = contentType
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

func (requestBuilder *requestBuilder) LogResponseBody() *requestBuilder {
	requestBuilder.logResponseBody = true
	return requestBuilder
}

func (requestBuilder *requestBuilder) LogRequestBody() *requestBuilder {
	requestBuilder.request.logRequestBody = true
	return requestBuilder
}

func (requestBuilder *requestBuilder) Execute(entityResponse interface{}) *Response {
	response, err := requestBuilder.request.execute(requestBuilder.client)
	if err != nil {
		return &Response{
			Error: err,
		}
	}
	rawResp, _ := httputil.DumpResponse(response, requestBuilder.logResponseBody)
	log.Println(string(rawResp))
	defer response.Body.Close()
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		body, _ := ioutil.ReadAll(response.Body)
		return &Response{
			StatusCode: response.StatusCode,
			Error:      requestBuilder.unmarshalFunctions[requestBuilder.contentType](body, entityResponse),
		}
	}
	return &Response{
		StatusCode: response.StatusCode,
		Error:      nil,
	}
}
