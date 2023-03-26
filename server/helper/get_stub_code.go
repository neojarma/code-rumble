package helper

import (
	"io"
	"net/http"
)

func GetStubCode(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	content := make([]byte, resp.ContentLength)
	_, err = resp.Body.Read(content)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return content, nil
}
