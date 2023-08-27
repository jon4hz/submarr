package lidarr

import (
	"github.com/jon4hz/submarr/internal/config"
	"github.com/jon4hz/submarr/internal/httpclient"
)

// Client represents a lidarr client
type Client struct {
	http httpclient.Client
	cfg  *config.LidarrConfig
}

// New creates a new lidarr client
func New(httpClient httpclient.Client, cfg *config.LidarrConfig) *Client {
	return &Client{
		http: httpClient,
		cfg:  cfg,
	}
}
