package mock

import (
	"net/http"
)

type HttpClientMock struct {
	MakeResponseFunction func(request *http.Request) (*http.Response, error)
}

func (mock *HttpClientMock) Do(request *http.Request) (*http.Response, error) {
	if mock.MakeResponseFunction != nil {
		return mock.MakeResponseFunction(request)
	}
	return &http.Response{}, nil
}
