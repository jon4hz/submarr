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
	// Get performs a GET request to the API.
	Get(ctx context.Context, base, endpoint string, expRes any, params ...map[string]string) (int, error)
	// Post performs a POST request to the API.
	Post(ctx context.Context, base, endpoint string, expRes, reqData any, params ...map[string]string) (int, error)
	// Put performs a PUT request to the API.
	Put(ctx context.Context, base, endpoint string, expRes, reqData any, params ...map[string]string) (int, error)
	// Delete performs a DELETE request to the API.
	Delete(ctx context.Context, base, endpoint string, expRes, reqData any, params ...map[string]string) (int, error)
}

// client is the actual http client.
type client struct {
	http   *http.Client
	apiKey string
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

// Get performs a GET request to the API.
func (c *client) Get(ctx context.Context, base, endpoint string, expRes any, params ...map[string]string) (int, error) {
	return c.doRequest(ctx, base, endpoint, http.MethodGet, expRes, nil, params...)
}

// Post performs a POST request to the API.
func (c *client) Post(ctx context.Context, base, endpoint string, expRes, reqData any, params ...map[string]string) (int, error) {
	return c.doRequest(ctx, base, endpoint, http.MethodPost, expRes, reqData, params...)
}

// Put performs a PUT request to the API.
func (c *client) Put(ctx context.Context, base, endpoint string, expRes, reqData any, params ...map[string]string) (int, error) {
	return c.doRequest(ctx, base, endpoint, http.MethodPut, expRes, reqData, params...)
}

// Delete performs a DELETE request to the API.
func (c *client) Delete(ctx context.Context, base, endpoint string, expRes, reqData any, params ...map[string]string) (int, error) {
	return c.doRequest(ctx, base, endpoint, http.MethodDelete, expRes, reqData, params...)
}

// DoRequest performs the request to the API.
func (c *client) doRequest(ctx context.Context, base, endpoint, method string, expRes, reqData any, params ...map[string]string) (int, error) {
	callURL, err := buildRequestURL(base, endpoint, params...)
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

	// set the API key
	if c.apiKey != "" {
		req.Header.Add("X-Api-Key", c.apiKey)
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

	case 401:
		return resp.StatusCode, fmt.Errorf("unauthorized")

	case 404:
		return resp.StatusCode, fmt.Errorf("not found")

	default:
		return resp.StatusCode, fmt.Errorf("%s", body)
	}
}

func buildRequestURL(base, endpoint string, params ...map[string]string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	u.Path = endpoint
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
