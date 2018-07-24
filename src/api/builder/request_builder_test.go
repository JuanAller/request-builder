package builder

import (
	"testing"
	"net/http"
	"github.com/JuanAller/request-builder/src/api/mock"
	"fmt"
	"errors"
)

type checkRequestFunc func(request *http.Request) error

func checkReqHeader(key string, value string) checkRequestFunc {
	return func(request *http.Request) error {
		if request.Header.Get(key) != value {
			return fmt.Errorf("Expected value : %v, in header %v, but got : %v ", value, key, request.Header.Get(key))
		}
		return nil
	}
}

func checkReqQueryParam(key string, value string) checkRequestFunc {
	return func(request *http.Request) error {
		if request.URL.Query().Get(key) != value {
			return fmt.Errorf("Expected value : %v, in query param : %v, but got : %v ", value, key, request.URL.Query().Get(key))
		}
		return nil
	}
}

func checkReqMethod(method string) checkRequestFunc {
	return func(request *http.Request) error {
		if request.Method != method {
			return fmt.Errorf("Expected method : %v, but got : %v ", method, request.Method)
		}
		return nil
	}
}

func checkReqFuncs(checks ...checkRequestFunc) checkRequestFunc {
	return func(request *http.Request) error {
		for _, check := range checks {
			if err := check(request); err != nil {
				return err
			}
		}
		return nil
	}
}

type checkRespFunc func(response *Response) error

func checkStatusCode(statusCode int) checkRespFunc {
	return func(response *Response) error {
		if statusCode != response.StatusCode {
			return fmt.Errorf("Expected : %v , but got : %v ", statusCode, response.StatusCode)
		}
		return nil
	}
}

func checkNotError() checkRespFunc {
	return func(response *Response) error {
		if response.Error != nil {
			return fmt.Errorf("Not expected error : %v ", response.Error)
		}
		return nil
	}
}

func checkErrorMessage(errMessage string) checkRespFunc {
	return func(response *Response) error {
		if errMessage != response.Error.Error() {
			fmt.Errorf("Expected error message : %v, but got : %v ", errMessage, response.Error.Error())
		}
		return nil
	}
}

func checkRespFuncs(checks ...checkRespFunc) checkRespFunc {
	return func(response *Response) error {
		for _, check := range checks {
			if err := check(response); err != nil {
				return err
			}
		}
		return nil
	}
}

func TestGet(t *testing.T) {
	responseMap := make(map[string]string)
	response := Get(&mock.HttpClientMock{
		MakeResponseFunction: func(request *http.Request) (*http.Response, error) {
			if err := checkReqFuncs(checkReqMethod("GET"))(request); err != nil {
				return nil, err
			}
			return mock.NewJsonResponse(http.StatusOK, map[string]string{"name": "aName"})
		},
	}, "http://test/get_ok").
		LogRequestBody().
		LogResponseBody().
		Execute(&responseMap)

	if err := checkRespFuncs(checkStatusCode(http.StatusOK), checkNotError())(response); err != nil {
		t.Error(err)
	}

	if responseMap["name"] != "aName" {
		t.Errorf("expected aName")
	}
}

func TestGetNotFound(t *testing.T) {
	responseMap := make(map[string]string)
	response := Get(&mock.HttpClientMock{
		MakeResponseFunction: func(request *http.Request) (*http.Response, error) {
			if err := checkReqFuncs(checkReqMethod("GET"))(request); err != nil {
				return nil, err
			}
			return mock.NewJsonResponse(http.StatusNotFound, map[string]string{"status": "not_found"})
		},
	}, "http://test/get_not_found").
		LogRequestBody().
		LogResponseBody().
		Execute(&responseMap)

	if err := checkRespFuncs(checkStatusCode(http.StatusNotFound), checkNotError())(response); err != nil {
		t.Error(err)
	}
}

func TestGetWithError(t *testing.T) {
	responseMap := make(map[string]string)
	response := Get(&mock.HttpClientMock{
		MakeResponseFunction: func(request *http.Request) (*http.Response, error) {
			if err := checkReqFuncs(checkReqMethod("GET"))(request); err != nil {
				return nil, err
			}
			return nil, errors.New("an error")
		},
	}, "http://test/get_not_found").
		LogRequestBody().
		LogResponseBody().
		Execute(&responseMap)

	if err := checkRespFuncs(checkErrorMessage("an error"))(response); err != nil {
		t.Error(err)
	}
}

func TestGetServerError(t *testing.T) {
	responseMap := make(map[string]string)
	response := Get(&mock.HttpClientMock{
		MakeResponseFunction: func(request *http.Request) (*http.Response, error) {
			if err := checkReqFuncs(checkReqMethod("GET"))(request); err != nil {
				return nil, err
			}
			return mock.NewJsonResponse(http.StatusInternalServerError, map[string]string{"status": "internal_server_error"})
		},
	}, "http://test/get_server_error").
		LogRequestBody().
		LogResponseBody().
		Execute(&responseMap)

	if err := checkRespFuncs(checkStatusCode(http.StatusInternalServerError), checkNotError())(response); err != nil {
		t.Error(err)
	}
}

func TestGetWithContentType(t *testing.T) {
	responseMap := make(map[string]interface{})
	response := Get(&mock.HttpClientMock{
		MakeResponseFunction: func(request *http.Request) (*http.Response, error) {
			if err := checkReqFuncs(checkReqMethod("GET"),
				checkReqHeader("Content-Type", "application/json"))(request); err != nil {
				return nil, err
			}
			return mock.NewJsonResponse(http.StatusOK, map[string]string{"status": "status_ok"})
		},
	}, "http://test/get_with_content_type").
		WithJSONContentType().
		LogRequestBody().
		LogResponseBody().
		Execute(&responseMap)

	if err := checkRespFuncs(checkStatusCode(http.StatusOK), checkNotError())(response); err != nil {
		t.Error(err)
	}
}

func TestGetXML(t *testing.T) {
	responseMap := &Response{}
	response := Get(&mock.HttpClientMock{
		MakeResponseFunction: func(request *http.Request) (*http.Response, error) {
			if err := checkReqFuncs(checkReqMethod("GET"),
				checkReqHeader("Content-Type", "application/xml"))(request); err != nil {
				return nil, err
			}
			return mock.NewXmlResponse(http.StatusOK, &Response{
				StatusCode: http.StatusOK,
				Error:      nil,
			})
		},
	}, "http://test/get_with_xml").
		WithXMLContentType().
		LogRequestBody().
		LogResponseBody().
		Execute(responseMap)

	if err := checkRespFuncs(checkStatusCode(http.StatusOK), checkNotError())(response); err != nil {
		t.Error(err)
	}

	if responseMap.StatusCode != 200 {
		t.Errorf("expected 200")
	}
}

func TestRequestBuilder_WithQueryParam(t *testing.T) {
	responseMap := make(map[string]interface{})
	response := Get(&mock.HttpClientMock{
		MakeResponseFunction: func(request *http.Request) (*http.Response, error) {
			if err := checkReqFuncs(checkReqMethod("GET"),
				checkReqQueryParam("query_param", "my_param"))(request); err != nil {
				return nil, err
			}
			return mock.NewJsonResponse(http.StatusOK, map[string]string{"status": "status_ok"})
		},
	}, "http://test/get_with_query_param").
		WithQueryParam("query_param", "my_param").
		LogRequestBody().
		LogResponseBody().
		Execute(&responseMap)

	if err := checkRespFuncs(checkStatusCode(http.StatusOK), checkNotError())(response); err != nil {
		t.Error(err)
	}
}

func TestPostOk(t *testing.T) {
	responseMap := make(map[string]string)
	response := Post(&mock.HttpClientMock{
		MakeResponseFunction: func(request *http.Request) (*http.Response, error) {
			if err := checkReqFuncs(checkReqMethod("POST"))(request); err != nil {
				return nil, err
			}
			return mock.NewJsonResponse(http.StatusCreated, map[string]string{"my_field": "my_value"})
		},
	}, "http://test/post_ok").
		WithBody(map[string]string{"my_field": "my_value"}).
		LogRequestBody().
		LogResponseBody().
		Execute(&responseMap)

	if err := checkRespFuncs(checkStatusCode(http.StatusCreated), checkNotError())(response); err != nil {
		t.Error(err)
	}

	if responseMap["my_field"] != "my_value" {
		t.Errorf("Expected my_value in response")
	}
}

func TestPostWithServerError(t *testing.T) {
	responseMap := make(map[string]interface{})
	response := Post(&mock.HttpClientMock{
		MakeResponseFunction: func(request *http.Request) (*http.Response, error) {
			if err := checkReqFuncs(checkReqMethod("POST"))(request); err != nil {
				return nil, err
			}
			return mock.NewJsonResponse(http.StatusInternalServerError, map[string]string{"status": "error"})
		},
	}, "http://test/post_server_error").
		LogRequestBody().
		LogResponseBody().
		Execute(&responseMap)

	if err := checkRespFuncs(checkStatusCode(http.StatusInternalServerError), checkNotError())(response); err != nil {
		t.Error(err)
	}
}

func TestGetWithBasicAuthentication(t *testing.T) {
	responseMap := make(map[string]interface{})
	response := Get(&mock.HttpClientMock{
		MakeResponseFunction: func(request *http.Request) (*http.Response, error) {
			if err := checkReqFuncs(checkReqMethod("GET"),
				checkReqHeader("Authorization", "Basic YWRtaW46YWRtaW4="))(request); err != nil {
				return nil, err
			}
			return mock.NewJsonResponse(http.StatusOK, map[string]string{"status": "ok"})
		},
	}, "http://test/get_basic_auth").
		WithJSONContentType().
		LogRequestBody().
		LogResponseBody().
		WithBasicAuthorization("admin", "admin").
		Execute(&responseMap)

	if err := checkRespFuncs(checkStatusCode(http.StatusOK), checkNotError())(response); err != nil {
		t.Error(err)
	}
}

func TestGetWithGZIPCompression(t *testing.T) {
	responseMap := make(map[string]interface{})
	response := Get(&mock.HttpClientMock{
		MakeResponseFunction: func(request *http.Request) (*http.Response, error) {
			if err := checkReqFuncs(checkReqMethod("GET"),
				checkReqHeader("Accept-Encoding", "gzip"))(request); err != nil {
				return nil, err
			}
			return mock.NewJsonGzipResponse(http.StatusOK, map[string]string{"status": "ok"})
		},
	}, "http://test/get_with_gzip").
		AcceptGzipEncoding().
		WithJSONContentType().
		LogRequestBody().
		LogResponseBody().
		Execute(&responseMap)

	if err := checkRespFuncs(checkStatusCode(http.StatusOK), checkNotError())(response); err != nil {
		t.Error(err)
	}

	if responseMap["status"] != "ok" {
		t.Errorf("Expected ok")
	}
}
