package builder

import (
	"compress/gzip"
	"bytes"
	"io/ioutil"
	"net/http"
)

type compressionAlgorithm func(data []byte) ([]byte, error)

func gzipAlgorithm(data []byte) ([]byte, error) {
	gzipReader, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()
	return ioutil.ReadAll(gzipReader)
}

func notAlgorithm(data []byte) ([]byte, error) {
	return data, nil
}

func compressionType(response *http.Response) string {
	return response.Header.Get("Content-Encoding")
}

func compressionFunctionsMap() map[string]compressionAlgorithm {
	return map[string]compressionAlgorithm{
		"gzip": gzipAlgorithm,
		"":     notAlgorithm,
	}
}
