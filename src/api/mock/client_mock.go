package mock

import (
	"net/http"
)

type HttpClientMock struct {
	makeResponseFunction func(request *http.Request) (*http.Response, error)
}

func (mock *HttpClientMock) Do(request *http.Request) (*http.Response, error) {
	if mock.makeResponseFunction != nil {
		return mock.makeResponseFunction(request)
	}
	return &http.Response{}, nil
}
