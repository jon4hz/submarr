package radarr

import (
	"github.com/jon4hz/subrr/pkg/radarr"
)

type Client struct {
	radarr *radarr.Client

	// is the client available?
	available bool

	// some client stats
	missing int
	queued  int
}

func New(radarr *radarr.Client) *Client {
	return &Client{
		radarr: radarr,
	}
}

func (c *Client) Init() error {
	/* ping, err := c.radarr.Ping(context.Background())
	if err != nil {
		return err
	}
	if strings.ToLower(ping.Status) == "ok" {
		c.available = true
	} else {
		return nil
	}

	queue, err := c.radarr.GetQueue(context.Background())
	if err != nil {
		return err
	}
	c.queued = int(queue.TotalRecords) */

	return nil
}

func (c *Client) ListItem() ClientItem {
	return ClientItem{c}
}
