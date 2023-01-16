package sonarr

import (
	"context"

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

type PingRes struct {
	Status string `json:"status"`
}

func (c *Client) Ping(ctx context.Context) (PingRes, error) {
	var res PingRes
	_, err := c.http.Get(ctx, c.cfg.Host, "/ping", &res)
	if err != nil {
		return res, err
	}
	return res, nil
}
