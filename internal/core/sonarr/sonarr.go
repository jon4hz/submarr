package sonarr

import (
	"context"
	"fmt"
	"strings"

	"github.com/jon4hz/subrr/pkg/sonarr"
)

type Client struct {
	sonarr *sonarr.Client

	// is the client available?
	available bool

	// some client stats
	missing int
	queued  int32

	// quality profiles by id
	qualityProfiles map[int32]sonarr.QualityProfileResource
}

func New(sonarr *sonarr.Client) *Client {
	if sonarr == nil {
		return nil
	}
	return &Client{
		sonarr: sonarr,
	}
}

// Init initializes the client and fetches some stats
func (c *Client) Init() error {
	ping, err := c.sonarr.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("failed to ping sonarr: %w", err)
	}
	if strings.ToLower(ping.Status) == "ok" {
		c.available = true
	} else {
		return nil
	}

	queue, err := c.sonarr.GetQueue(context.Background())
	if err != nil {
		c.available = false
		return fmt.Errorf("failed to get queue: %w", err)
	}
	c.queued = queue.TotalRecords

	return nil
}

func (c *Client) ClientListItem() ClientItem {
	return ClientItem{c}
}
