package client

import (
	"io"
	"net/http"
	"strings"
)

const (
	HeaderContentType = "Content-Type"
)

type Client struct {
	ResponseBody string
}

func (cc *Client) DoRequestWithoutBody(method, url string) (*http.Response, error) {
	return doRequest(method, url, "")
}

func doRequest(method, url, requestBody string) (*http.Response, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	var bodyReader io.Reader
	if requestBody != "" {
		bodyReader = strings.NewReader(requestBody)
	}

	req, _ := http.NewRequest(method, url, bodyReader)
	req.Header.Add(HeaderContentType, "application/json")
	resp, err := client.Do(req)

	return resp, err
}
