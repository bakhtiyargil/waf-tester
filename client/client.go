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
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &Client{httpClient: httpClient}
}

func (c *Client) DoRequestWithoutBody(method, url string) (respBody []byte, statusCode int, err error) {
	resp, err := c.doRequest(method, url, "")
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

func (c *Client) doRequest(method, url, requestBody string) (*http.Response, error) {
	var bodyReader io.Reader
	if requestBody != "" {
		bodyReader = strings.NewReader(requestBody)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add(HeaderContentType, "application/json")
	resp, err := c.httpClient.Do(req)

	return resp, err
}
