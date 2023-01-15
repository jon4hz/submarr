package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is a http client.
type Client interface {
	DoRequest(ctx context.Context, baseURL, method string, expRes, reqData any, params ...map[string]string) (int, error)
}

// client is the actual http client.
type client struct {
	http *http.Client
}

// New creates a new http client.
func New(opts ...ClientOpts) Client {
	c := &client{
		http: &http.Client{},
	}
	c.http.Timeout = 30 * time.Second
	for _, o := range opts {
		o(c)
	}
	return c
}

// DoRequest performs the request to the API.
func (c *client) DoRequest(ctx context.Context, baseURL, method string, expRes, reqData any, params ...map[string]string) (int, error) {
	callURL, err := buildRequestURL(baseURL, params...)
	if err != nil {
		return 0, err
	}

	var dataReq []byte
	if reqData != nil {
		dataReq, err = json.Marshal(reqData)
		if err != nil {
			return 0, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, callURL, bytes.NewBuffer(dataReq))
	if err != nil {
		return 0, err
	}
	if dataReq != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	switch resp.StatusCode {
	case 200:
		if expRes != nil {
			err = json.Unmarshal(body, expRes)
			if err != nil {
				return 0, err
			}
		}
		return resp.StatusCode, nil

	default:
		return resp.StatusCode, fmt.Errorf("%s", body)
	}
}

func buildRequestURL(baseURL string, params ...map[string]string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	if len(params) == 0 {
		return u.String(), nil
	}
	p := url.Values{}
	for k, v := range params[0] {
		p.Set(k, v)
	}
	u.RawQuery = p.Encode()
	return u.String(), nil
}
