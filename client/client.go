package client

import (
	"io"
	"net/http"
	"strings"
)

const (
	HeaderContentType = "Content-Type"
)

type Client struct{}

func (cc *Client) DoRequestWithoutBody(method, url string) (respBody []byte, statusCode int, err error) {
	resp, err := doRequest(method, url, "")
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		closeErr := resp.Body.Close()
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	body, err := io.ReadAll(resp.Body)
	return body, resp.StatusCode, err
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

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add(HeaderContentType, "application/json")
	resp, err := client.Do(req)

	return resp, err
}
