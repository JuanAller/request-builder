package builder

import (
	"net/http"
	"encoding/xml"
	"encoding/json"
)

func Get(client HttpClient, path string) *requestBuilder {
	return &requestBuilder{
		client:               client,
		request:              newRequest(http.MethodGet, path),
		contentType:          APPLICATIONJSON,
		unmarshalFunctions:   unmarshalFunctionsMap(),
		compressionFunctions: compressionFunctionsMap(),
	}
}

func Post(client HttpClient, path string) *requestBuilder {
	return &requestBuilder{
		client:               client,
		request:              newRequest(http.MethodPost, path),
		contentType:          APPLICATIONJSON,
		unmarshalFunctions:   unmarshalFunctionsMap(),
		compressionFunctions: compressionFunctionsMap(),
	}
}

func Put(client HttpClient, path string) *requestBuilder {
	return &requestBuilder{
		client:               client,
		request:              newRequest(http.MethodPut, path),
		contentType:          APPLICATIONJSON,
		unmarshalFunctions:   unmarshalFunctionsMap(),
		compressionFunctions: compressionFunctionsMap(),
	}
}

func Delete(client HttpClient, path string) *requestBuilder {
	return &requestBuilder{
		client:               client,
		request:              newRequest(http.MethodDelete, path),
		contentType:          APPLICATIONJSON,
		unmarshalFunctions:   unmarshalFunctionsMap(),
		compressionFunctions: compressionFunctionsMap(),
	}
}

func unmarshalFunctionsMap() map[string]func([]byte, interface{}) error {
	return map[string]func([]byte, interface{}) error{
		APPLICATIONJSON: json.Unmarshal,
		APPLICATIONXML:  xml.Unmarshal,
	}
}
