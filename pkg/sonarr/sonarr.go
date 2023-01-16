package sonarr

import (
	"context"

	"github.com/jon4hz/subrr/internal/config"
	"github.com/jon4hz/subrr/internal/httpclient"
)

// Client represents a sonarr client
type Client struct {
	http httpclient.Client
	cfg  *config.SonarrConfig
}

// New creates a new sonarr client
func New(httpClient httpclient.Client, cfg *config.SonarrConfig) *Client {
	return &Client{
		http: httpClient,
		cfg:  cfg,
	}
}

// PingRes is the response from the ping endpoint
type PingRes struct {
	Status string `json:"status"`
}

// Ping pings the sonarr server
func (c *Client) Ping(ctx context.Context) (PingRes, error) {
	var res PingRes
	_, err := c.http.Get(ctx, c.cfg.Host, "/ping", &res)
	if err != nil {
		return res, err
	}
	return res, nil
}
