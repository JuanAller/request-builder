package builder

import (
	"testing"
	"net/http"
	"github.com/JuanAller/request-builder/src/api/mock"
)

func TestGet(t *testing.T) {
	responseMap := make(map[string]string)
	response := Get(&mock.HttpClientMock{
		MakeResponseFunction: func(request *http.Request) (*http.Response, error) {
			return mock.NewJsonResponse(http.StatusOK, map[string]string{"name": "aName"})
		},
	}, "http://test/get_ok").Execute(&responseMap)
	if response.Error != nil {
		t.Errorf("not expected error")
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

	if response.StatusCode != http.StatusNotFound {
		t.Errorf("expected not found")
	}
}

func TestGetServerError(t *testing.T) {
	responseMap := make(map[string]string)
	response := Get(&mock.HttpClientMock{
		MakeResponseFunction: func(request *http.Request) (*http.Response, error) {
			return mock.NewJsonResponse(http.StatusInternalServerError, map[string]string{"status": "internal_server_error"})
		},
	}, "http://test/get_server_error").Execute(&responseMap)

	if response.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected server error")
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
	if response.StatusCode != http.StatusOK {
		t.Errorf("expected 200 status code")
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
	if response.StatusCode != http.StatusOK {
		t.Errorf("expected 200 status code")
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

	if response.StatusCode != http.StatusOK {
		t.Errorf("expecte 200 status code")
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

	if response.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201")
	}
	if response.Error != nil {
		t.Errorf("not error expected")
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
	if response.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 status code")
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
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected 200")
	}
}
