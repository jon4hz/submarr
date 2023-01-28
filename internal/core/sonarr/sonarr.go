package sonarr

import (
	"context"
	"strings"

	"github.com/jon4hz/subrr/pkg/sonarr"
)

type Client struct {
	sonarr *sonarr.Client

	// is the client available?
	available bool

	// some client stats
	missing int
	queued  int
}

func New(sonarr *sonarr.Client) *Client {
	if sonarr == nil {
		return nil
	}
	return &Client{
		sonarr: sonarr,
	}
}

func (c *Client) Init() error {
	ping, err := c.sonarr.Ping(context.Background())
	if err != nil {
		return err
	}
	if strings.ToLower(ping.Status) == "ok" {
		c.available = true
	} else {
		return nil
	}

	queue, err := c.sonarr.GetQueue(context.Background())
	if err != nil {
		return err
	}
	c.queued = int(queue.TotalRecords)

	return nil
}

func (c *Client) ListItem() ClientItem {
	return ClientItem{c}
}
