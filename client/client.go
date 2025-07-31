package client

import (
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	HeaderContentType = "Content-Type"
)

type Client interface {
	DoRequestWithoutBody(method, url string) (respBody []byte, statusCode int, elapsed time.Duration, err error)
}

type PureHttpClient struct {
	httpClient *http.Client
	bodyReader io.Reader
}

func NewPureHttpClient() Client {
	httpClient := http.DefaultClient
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &PureHttpClient{httpClient: httpClient}
}

func (c *PureHttpClient) DoRequestWithoutBody(method, url string) (
	respBody []byte,
	statusCode int,
	elapsed time.Duration,
	err error,
) {
	start := time.Now()
	resp, err := c.doRequest(method, url, "")
	if err != nil {
		return nil, 0, 0, err
	}
	stop := time.Now()
	elapsed = stop.Sub(start)

	defer func() {
		closeErr := resp.Body.Close()
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()
	body, err := io.ReadAll(resp.Body)
	return body, resp.StatusCode, elapsed, err
}

func (c *PureHttpClient) doRequest(method, url, requestBody string) (*http.Response, error) {
	if requestBody != "" {
		c.bodyReader = strings.NewReader(requestBody)
	}

	req, err := http.NewRequest(method, url, c.bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add(HeaderContentType, "application/json")
	resp, err := c.httpClient.Do(req)
	return resp, err
}
