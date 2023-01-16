package sonarr

import (
	"context"
	"strconv"

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

// Ping pings the sonarr server
func (c *Client) Ping(ctx context.Context) (*Ping, error) {
	var res Ping
	_, err := c.http.Get(ctx, c.cfg.Host, "/ping", &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetSeries returns a list of all series
func (c *Client) GetSeries(ctx context.Context) ([]SeriesResource, error) {
	var res []SeriesResource
	_, err := c.http.Get(ctx, c.cfg.Host, "/api/v3/series", &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

// GetSerie returns a serie by its TVDB ID
func (c *Client) GetSerie(ctx context.Context, tvdbID int) (*SeriesResource, error) {
	var res []SeriesResource
	_, err := c.http.Get(ctx, c.cfg.Host, "/api/v3/series", &res, map[string]string{"tvdbId": strconv.Itoa(tvdbID)})
	if err != nil {
		return nil, err
	}
	return &res[0], nil
}

// GetQueue returns the current download queue
func (c *Client) GetQueue(ctx context.Context) (QueueResourcePagingResource, error) {
	var res QueueResourcePagingResource
	_, err := c.http.Get(ctx, c.cfg.Host, "/api/v3/queue", &res)
	if err != nil {
		return res, err
	}
	return res, nil
}
