package utils

import (
	"io"
	"io/ioutil"
	"net/http"
)

// HTTPPost send json data to external api service
func HTTPPost(apiURL string, contentType string, body io.Reader) (rs []byte, err error) {
	resp, err := http.Post(apiURL, contentType, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func HTTPGet(apiURL string) (rs []byte, err error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
