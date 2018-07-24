package builder

import (
	"testing"
	"net/http"
	"github.com/JuanAller/request-builder/src/api/mock"
	"fmt"
	"errors"
)

type checkFunc func(response *Response) error

func checkStatusCode(statusCode int) checkFunc {
	return func(response *Response) error {
		if statusCode != response.StatusCode {
			return fmt.Errorf("Expected : %v , but got : %v ", statusCode, response.StatusCode)
		}
		return nil
	}
}

func checkNotError() checkFunc {
	return func(response *Response) error {
		if response.Error != nil {
			return fmt.Errorf("Not expected error : %v ", response.Error)
		}
		return nil
	}
}

func checkErrorMessage(errMessage string) checkFunc {
	return func(response *Response) error {
		if errMessage != response.Error.Error() {
			fmt.Errorf("Expected error message : %v, but got : %v ", errMessage, response.Error.Error())
		}
		return nil
	}
}

func checkFuncs(checks ...checkFunc) checkFunc {
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
			return mock.NewJsonResponse(http.StatusOK, map[string]string{"name": "aName"})
		},
	}, "http://test/get_ok").Execute(&responseMap)

	if err := checkFuncs(checkNotError())(response); err != nil {
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
			return mock.NewJsonResponse(http.StatusNotFound, map[string]string{"status": "not_found"})
		},
	}, "http://test/get_not_found").Execute(&responseMap)

	if err := checkFuncs(checkStatusCode(http.StatusNotFound), checkNotError())(response); err != nil {
		t.Error(err)
	}
}

func TestGetWithError(t *testing.T) {
	responseMap := make(map[string]string)
	response := Get(&mock.HttpClientMock{
		MakeResponseFunction: func(request *http.Request) (*http.Response, error) {
			return nil, errors.New("an error")
		},
	}, "http://test/get_not_found").Execute(&responseMap)

	if err := checkFuncs(checkErrorMessage("an error"))(response); err != nil {
		t.Error(err)
	}
}

func TestGetServerError(t *testing.T) {
	responseMap := make(map[string]string)
	response := Get(&mock.HttpClientMock{
		MakeResponseFunction: func(request *http.Request) (*http.Response, error) {
			return mock.NewJsonResponse(http.StatusInternalServerError, map[string]string{"status": "internal_server_error"})
		},
	}, "http://test/get_server_error").Execute(&responseMap)

	if err := checkFuncs(checkStatusCode(http.StatusInternalServerError), checkNotError())(response); err != nil {
		t.Error(err)
	}
}

func TestGetWithContentType(t *testing.T) {
	responseMap := make(map[string]interface{})
	response := Get(&mock.HttpClientMock{
		MakeResponseFunction: func(request *http.Request) (*http.Response, error) {
			if request.Header.Get("Content-Type") != APPLICATIONJSON {
				return mock.NewJsonResponse(http.StatusInternalServerError, map[string]string{"status": "fail"})
			}
			return mock.NewJsonResponse(http.StatusOK, map[string]string{"status": "status_ok"})
		},
	}, "http://test/get_with_content_type").
		WithJSONContentType().
		Execute(&responseMap)

	if err := checkFuncs(checkStatusCode(http.StatusOK), checkNotError())(response); err != nil {
		t.Error(err)
	}
}

func TestGetXML(t *testing.T) {
	responseMap := &Response{}
	response := Get(&mock.HttpClientMock{
		MakeResponseFunction: func(request *http.Request) (*http.Response, error) {
			if request.Header.Get("Content-Type") != APPLICATIONXML {
				return mock.NewXmlResponse(http.StatusInternalServerError, &Response{
					StatusCode: http.StatusInternalServerError,
					Error:      nil,
				})
			}
			return mock.NewXmlResponse(http.StatusOK, &Response{
				StatusCode: http.StatusOK,
				Error:      nil,
			})
		},
	}, "http://test/get_with_xml").
		WithXMLContentType().
		Execute(responseMap)

	if err := checkFuncs(checkStatusCode(http.StatusOK), checkNotError())(response); err != nil {
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
			if request.URL.Query().Get("query_param") != "my_param" {
				mock.NewJsonResponse(http.StatusInternalServerError, map[string]string{"status": "error"})
			}
			return mock.NewJsonResponse(http.StatusOK, map[string]string{"status": "status_ok"})
		},
	}, "http://test/get_with_query_param").
		WithQueryParam("query_param", "my_param").
		Execute(&responseMap)

	if err := checkFuncs(checkStatusCode(http.StatusOK), checkNotError())(response); err != nil {
		t.Error(err)
	}
}

func TestPostOk(t *testing.T) {
	responseMap := make(map[string]string)
	response := Post(&mock.HttpClientMock{
		MakeResponseFunction: func(request *http.Request) (*http.Response, error) {
			return mock.NewJsonResponse(http.StatusCreated, map[string]string{"my_field": "my_value"})
		},
	}, "http://test/post_ok").
		WithBody(map[string]string{"my_field": "my_value"}).
		Execute(&responseMap)

	if err := checkFuncs(checkStatusCode(http.StatusCreated), checkNotError())(response); err != nil {
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
			return mock.NewJsonResponse(http.StatusInternalServerError, map[string]string{"status": "error"})
		},
	}, "http://test/post_server_error").Execute(&responseMap)

	if err := checkFuncs(checkStatusCode(http.StatusInternalServerError), checkNotError())(response); err != nil {
		t.Error(err)
	}
}

func TestGetWithBasicAuthentication(t *testing.T) {
	responseMap := make(map[string]interface{})
	response := Get(&mock.HttpClientMock{
		MakeResponseFunction: func(request *http.Request) (*http.Response, error) {
			if request.Header.Get("Authorization") != "Basic YWRtaW46YWRtaW4=" {
				return mock.NewJsonResponse(http.StatusInternalServerError, map[string]string{"status": "fail"})
			}
			return mock.NewJsonResponse(http.StatusOK, map[string]string{"status": "ok"})
		},
	}, "http://test/get_basic_auth").
		WithJSONContentType().
		WithBasicAuthorization("admin", "admin").
		Execute(&responseMap)

	if err := checkFuncs(checkStatusCode(http.StatusOK), checkNotError())(response); err != nil {
		t.Error(err)
	}
}

func TestGetWithGZIPCompression(t *testing.T) {
	responseMap := make(map[string]interface{})
	response := Get(&mock.HttpClientMock{
		MakeResponseFunction: func(request *http.Request) (*http.Response, error) {
			if request.Header.Get("Accept-Encoding") != "gzip" {
				return mock.NewJsonGzipResponse(http.StatusInternalServerError, map[string]string{"status": "fail"})
			}
			return mock.NewJsonGzipResponse(http.StatusOK, map[string]string{"status": "ok"})
		},
	}, "http://test/get_with_gzip").
		AcceptGzipEncoding().
		WithJSONContentType().
		LogRequestBody().
		LogResponseBody().
		Execute(&responseMap)

	if err := checkFuncs(checkStatusCode(http.StatusOK), checkNotError())(response); err != nil {
		t.Error(err)
	}

	if responseMap["status"] != "ok" {
		t.Errorf("Expected ok")
	}
}
