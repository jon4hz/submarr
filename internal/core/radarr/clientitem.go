package radarr

import (
	"fmt"
	"strings"
)

type ClientItem struct {
	c *Client
}

func (i ClientItem) String() string { return "radarr" }

func (i ClientItem) FilterValue() string { return "" }

func (i ClientItem) Title() string { return strings.Title(i.String()) }

func (i ClientItem) Available() bool { return i.c.available }

func (i ClientItem) Stats() []string {
	return []string{
		fmt.Sprintf("%d queued", i.c.queued),
		fmt.Sprintf("%d missing", i.c.missing),
	}
}