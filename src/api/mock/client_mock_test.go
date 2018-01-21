package mock

import (
	"testing"
	"net/http"
	"github.com/JuanAller/request-builder/src/api/builder"
)

func TestHttpClientMock_DoWithJsonBody(t *testing.T) {
	mockClient := &HttpClientMock{
		makeResponseFunction: func(request *http.Request) (*http.Response, error) {
			return NewJsonResponse(http.StatusOK, map[string]string{"hola": "mundo"})
		},
	}

	responseMap := make(map[string]string)
	response := builder.Get(mockClient, "http://mock.test/id").
		LogRequestBody().
		LogResponseBody().
		Execute(&responseMap)

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected 200")
	}
	if responseMap["hola"] != "mundo" {
		t.Errorf("Expected body {hola : mundo}")
	}
}

func TestHttpClientMock_DoWithXMLBody(t *testing.T) {
	mockClient := &HttpClientMock{
		makeResponseFunction: func(request *http.Request) (*http.Response, error) {
			return NewXmlResponse(http.StatusOK, &builder.Response{
				StatusCode: http.StatusOK,
				Error:      nil,
			})
		},
	}

	responseStruct := &builder.Response{}
	response := builder.Get(mockClient, "http://mock.test/id").
		WithContentType("application/xml").
		LogRequestBody().
		LogResponseBody().
		Execute(responseStruct)

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expecte 200")
	}
	if responseStruct.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 in body response")
	}
}
