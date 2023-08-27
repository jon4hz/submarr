package radarr

import (
	"github.com/jon4hz/submarr/internal/config"
	"github.com/jon4hz/submarr/pkg/radarr"
)

type Client struct {
	Config *config.RadarrConfig
	radarr *radarr.Client

	// is the client available?
	available bool

	// some client stats
	missing int
	queued  int
}

func New(cfg *config.RadarrConfig, radarr *radarr.Client) *Client {
	if radarr == nil || cfg == nil {
		return nil
	}
	return &Client{
		Config: cfg,
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
	//return errors.New("test")
}

func (c *Client) ClientListItem() ClientItem {
	return ClientItem{c}
}
