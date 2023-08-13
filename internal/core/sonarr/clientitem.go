package sonarr

import (
	"fmt"

	"github.com/jon4hz/subrr/internal/tui/common"
)

type ClientItem struct {
	c *Client
}

func (i ClientItem) String() string { return "sonarr" }

func (i ClientItem) FilterValue() string { return "" }

func (i ClientItem) Title() string { return common.Title(i.String()) }

func (i ClientItem) Available() bool { return i.c.available }

func (i ClientItem) Stats() []string {
	return []string{
		fmt.Sprintf("%d queued", i.c.totalQueued),
		fmt.Sprintf("%d missing", i.c.totalMissing),
	}
}
