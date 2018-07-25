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

type unmarshalFunc func([]byte, interface{}) error

type requestBuilder struct {
	client               HttpClient
	request              *request
	contentType          string
	unmarshalFunctions   map[string]func([]byte, interface{}) error
	compressionFunctions map[string]compressionAlgorithm
	logResponseBody      bool
}

func (requestBuilder *requestBuilder) WithQueryParam(key string, value string) *requestBuilder {
	requestBuilder.request.QueryParams[key] = value
	return requestBuilder
}

func (requestBuilder *requestBuilder) WithHeader(key string, value string) *requestBuilder {
	requestBuilder.request.Headers[key] = value
	return requestBuilder
}

func (requestBuilder *requestBuilder) AcceptGzipEncoding() *requestBuilder {
	return requestBuilder.WithHeader("Accept-Encoding", "gzip")
}

func (requestBuilder *requestBuilder) WithJSONContentType() *requestBuilder {
	requestBuilder.contentType = APPLICATIONJSON
	return requestBuilder.WithHeader("Content-Type", APPLICATIONJSON)
}

func (requestBuilder *requestBuilder) WithCustomJSONUnmarshal(custom unmarshalFunc) *requestBuilder {
	requestBuilder.unmarshalFunctions[APPLICATIONJSON] = custom
	return requestBuilder
}

func (requestBuilder *requestBuilder) WithXMLContentType() *requestBuilder {
	requestBuilder.contentType = APPLICATIONXML
	return requestBuilder.WithHeader("Content-Type", APPLICATIONXML)
}

func (requestBuilder *requestBuilder) WithCustomXMLUnmarshal(custom unmarshalFunc) *requestBuilder {
	requestBuilder.unmarshalFunctions[APPLICATIONXML] = custom
	return requestBuilder
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
	if requestBuilder.logResponseBody {
		rawResp, _ := httputil.DumpResponse(response, requestBuilder.logResponseBody)
		log.Println(string(rawResp))
	}
	defer response.Body.Close()
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		body, _ := ioutil.ReadAll(response.Body)
		if body, err = requestBuilder.compressionFunctions[compressionType(response)](body); err != nil {
			return &Response{
				StatusCode: response.StatusCode,
				Error:      err,
			}
		}
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
