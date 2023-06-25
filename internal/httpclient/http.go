package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrNotFound     = errors.New("not found")
)

// Client is a http client.
type Client interface {
	// Get performs a GET request to the API.
	Get(ctx context.Context, base, endpoint string, expRes any, opts ...RequestOpts) (int, error)
	// Post performs a POST request to the API.
	Post(ctx context.Context, base, endpoint string, expRes, reqData any, opts ...RequestOpts) (int, error)
	// Put performs a PUT request to the API.
	Put(ctx context.Context, base, endpoint string, expRes, reqData any, opts ...RequestOpts) (int, error)
	// Delete performs a DELETE request to the API.
	Delete(ctx context.Context, base, endpoint string, expRes, reqData any, opts ...RequestOpts) (int, error)
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
func (c *client) Get(ctx context.Context, base, endpoint string, expRes any, opts ...RequestOpts) (int, error) {
	return c.doRequest(ctx, base, endpoint, http.MethodGet, expRes, nil, opts...)
}

// Post performs a POST request to the API.
func (c *client) Post(ctx context.Context, base, endpoint string, expRes, reqData any, opts ...RequestOpts) (int, error) {
	return c.doRequest(ctx, base, endpoint, http.MethodPost, expRes, reqData, opts...)
}

// Put performs a PUT request to the API.
func (c *client) Put(ctx context.Context, base, endpoint string, expRes, reqData any, opts ...RequestOpts) (int, error) {
	return c.doRequest(ctx, base, endpoint, http.MethodPut, expRes, reqData, opts...)
}

// Delete performs a DELETE request to the API.
func (c *client) Delete(ctx context.Context, base, endpoint string, expRes, reqData any, opts ...RequestOpts) (int, error) {
	return c.doRequest(ctx, base, endpoint, http.MethodDelete, expRes, reqData, opts...)
}

// DoRequest performs the request to the API.
func (c *client) doRequest(ctx context.Context, base, endpoint, method string, expRes, reqData any, opts ...RequestOpts) (int, error) {
	r := newRequest(opts...)
	callURL, err := buildRequestURL(base, endpoint, r)
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

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if expRes != nil {
			err = json.Unmarshal(body, expRes)
			if err != nil {
				return 0, err
			}
		}
		return resp.StatusCode, nil
	}

	switch resp.StatusCode {
	case 401:
		return resp.StatusCode, ErrUnauthorized
	case 404:
		return resp.StatusCode, ErrNotFound
	default:
		return resp.StatusCode, fmt.Errorf("%s", body)
	}
}

func buildRequestURL(base, endpoint string, r *Request) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	u.Path = endpoint
	p := url.Values{}

	// set paging options
	if r.page != 0 {
		p.Set("page", strconv.Itoa(r.page))
	}
	if r.pageSize != 0 {
		p.Set("pageSize", strconv.Itoa(r.pageSize))
	}
	// set sorting options
	if r.sortKey != "" {
		p.Set("sortKey", r.sortKey)
	}
	if r.sortDirection != "" {
		p.Set("sortDirection", string(r.sortDirection))
	}
	// add additional query params if there are any
	if r.params != nil && len(r.params) > 0 {
		for k, v := range r.params {
			p.Set(k, v)
		}
	}

	// encode the query parameters
	u.RawQuery = p.Encode()
	return u.String(), nil
}
