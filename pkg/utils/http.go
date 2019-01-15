package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

// HTTPPost send json data to external api service
func HTTPPost(apiURL string, contentType string, body io.Reader) ([]byte, error) {
	resp, err := http.Post(apiURL, contentType, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func HTTPGet(apiURL string) ([]byte, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func UploadFile(url string, fieldName string, fileName string, src io.Reader) ([]byte, error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	var fw io.Writer
	var err error
	if fw, err = w.CreateFormFile(fieldName, fileName); err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, src); err != nil {
		return nil, err
	}

	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return nil, err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		Log.Errorf("upload file to [%v] fail, http status = [%v]", url, res.StatusCode)
	}

	return ioutil.ReadAll(res.Body)
}
