package sonarr

import (
	"github.com/jon4hz/subrr/internal/config"
	"github.com/jon4hz/subrr/internal/httpclient"
)

type Client struct {
	http httpclient.Client
	cfg  *config.SonarrConfig
}

func New(httpClient httpclient.Client, cfg *config.SonarrConfig) *Client {
	return &Client{
		http: httpClient,
		cfg:  cfg,
	}
}
