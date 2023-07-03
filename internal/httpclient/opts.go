package httpclient

import (
	"crypto/tls"
	"net/http"
	"time"
)

// ClientOpts are options for the client.
type ClientOpts func(*client)

// WithoutTLSVerify disables TLS verification.
func WithoutTLSVerfiy(disableTLS ...bool) ClientOpts {
	if len(disableTLS) > 0 && !disableTLS[0] {
		return func(c *client) {}
	}
	// #nosec G402
	return func(c *client) {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		c.http.Transport = tr
	}
}

// WithTimout sets the default request timeout
func WithTimeout(t time.Duration) ClientOpts {
	return func(c *client) {
		c.http.Timeout = t
	}
}

// WithAPIKey sets the API key for the client.
func WithAPIKey(apiKey string) ClientOpts {
	return func(c *client) {
		c.apiKey = apiKey
	}
}

// WithBasicAuth adds basic authentication headers to the http requests
func WithBasicAuth(username, password string) ClientOpts {
	return func(c *client) {
		c.basicAuth = &basicAuth{
			username: username,
			password: password,
		}
	}
}

// WithHeader adds a http header to all requests
func WithHeader(key, value string) ClientOpts {
	return func(c *client) {
		if c.headers == nil {
			c.headers = make([]header, 0, 1)
		}
		c.headers = append(c.headers, header{key: key, value: value})
	}
}
