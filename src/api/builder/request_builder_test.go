package builder

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"os"
	"time"
	"encoding/xml"
)

var tmux = http.NewServeMux()
var server = httptest.NewServer(tmux)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	tmux.HandleFunc("/get_ok", func(writer http.ResponseWriter, request *http.Request) {
		body, _ := json.Marshal(map[string]string{"name": "aName"})
		writer.Write(body)
		return
	})
	tmux.HandleFunc("/get_not_found", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNotFound)
	})
	tmux.HandleFunc("/get_server_error", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
	})
	tmux.HandleFunc("/get_timeout", func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(time.Millisecond * 20)
		writer.WriteHeader(http.StatusOK)
	})
	tmux.HandleFunc("/get_with_xml", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		body, _ := xml.Marshal(&Response{
			StatusCode: 200,
			Error: nil,
		})
		writer.Write(body)
	})
	tmux.HandleFunc("/get_with_query_param", func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Query().Get("query_param") == "my_param" {
			writer.WriteHeader(http.StatusOK)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
		}
	})
	tmux.HandleFunc("/get_with_content_type", func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Content-Type") == "application/json" {
			writer.WriteHeader(http.StatusOK)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
		}
	})
	tmux.HandleFunc("/post_ok", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusCreated)
		body, _ := json.Marshal(map[string]string{"my_field": "my_value"})
		writer.Write(body)
	})
	tmux.HandleFunc("/post_server_error", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
	})
}

func TestGetTimeout(t *testing.T) {
	responseMap := make(map[string]string)
	response := Get(&http.Client{Timeout: time.Millisecond * 10}, server.URL+"/get_timeout").Execute(&responseMap)
	if response.Error == nil {
		t.Errorf("expected timeout error")
	}
}

func TestGet(t *testing.T) {
	responseMap := make(map[string]string)
	response := Get(&http.Client{}, server.URL+"/get_ok").Execute(&responseMap)
	if response.Error != nil {
		t.Errorf("not expected error")
	}
	if responseMap["name"] != "aName" {
		t.Errorf("expected aName")
	}
}

func TestGetNotFound(t *testing.T) {
	responseMap := make(map[string]string)
	response := Get(&http.Client{}, server.URL+"/get_not_found").Execute(&responseMap)
	if response.StatusCode != http.StatusNotFound {
		t.Errorf("expected not found")
	}
}

func TestGetServerError(t *testing.T) {
	responseMap := make(map[string]string)
	response := Get(&http.Client{}, server.URL+"/get_server_error").Execute(&responseMap)
	if response.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected server error")
	}
}

func TestGetWithContentType(t *testing.T) {
	responseMap := make(map[string]interface{})
	response := Get(&http.Client{}, server.URL+"/get_with_content_type").
		WithJSONContentType().
		Execute(&responseMap)
	if response.StatusCode != http.StatusOK {
		t.Errorf("expected 200 status code")
	}
}

func TestGetXML(t *testing.T)  {
	responseMap := &Response{}
	response := Get(&http.Client{}, server.URL+"/get_with_xml").
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
	response := Get(&http.Client{}, server.URL+"/get_with_query_param").
		WithQueryParam("query_param", "my_param").
		Execute(&responseMap)
	if response.StatusCode != http.StatusOK {
		t.Errorf("expecte 200 status code")
	}
}

func TestPostOk(t *testing.T) {
	responseMap := make(map[string]string)
	body := map[string]string{"my_field": "my_value"}
	response := Post(&http.Client{}, server.URL+"/post_ok").
		WithBody(body).
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
	response := Post(&http.Client{}, server.URL+"/post_server_error").Execute(&responseMap)
	if response.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 status code")
	}
}
