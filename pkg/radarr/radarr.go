package radarr

import (
	"github.com/jon4hz/subrr/internal/config"
	"github.com/jon4hz/subrr/internal/httpclient"
)

// Client represents a radarr client
type Client struct {
	http httpclient.Client
	cfg  *config.RadarrConfig
}

// New creates a new radarr client
func New(httpClient httpclient.Client, cfg *config.RadarrConfig) *Client {
	return &Client{
		http: httpClient,
		cfg:  cfg,
	}
}
