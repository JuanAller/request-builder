package mock

import (
	"strconv"
	"bytes"
	"net/http"
	"encoding/json"
	"encoding/xml"
	"compress/gzip"
)

func NewJsonResponse(status int, body interface{}) (*http.Response, error) {
	encoded, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	response := newBytesResponse(status, encoded)
	response.Header.Set("Content-Type", "application/json")
	return response, nil
}

func NewXmlResponse(status int, body interface{}) (*http.Response, error) {
	encoded, err := xml.Marshal(body)
	if err != nil {
		return nil, err
	}
	response := newBytesResponse(status, encoded)
	response.Header.Set("Content-Type", "application/xml")
	return response, nil
}

func NewJsonGzipResponse(status int, body interface{}) (*http.Response, error) {
	encoded, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write(encoded)
	w.Close()
	response := newBytesResponse(status, b.Bytes())
	response.Header.Set("Content-Encoding", "gzip")
	return response, nil
}

func newBytesResponse(status int, body []byte) *http.Response {
	return &http.Response{
		Status:     strconv.Itoa(status),
		StatusCode: status,
		Body:       &dummyReadCloser{bytes.NewReader(body)},
		Header:     http.Header{},
	}
}
