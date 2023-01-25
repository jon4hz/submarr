package sonarr

import (
	"fmt"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ClientItem struct {
	c *Client
}

func (i ClientItem) String() string { return "sonarr" }

func (i ClientItem) FilterValue() string { return "" }

func (i ClientItem) Title() string { return cases.Title(language.AmericanEnglish).String(i.String()) }

func (i ClientItem) Available() bool { return i.c.available }

func (i ClientItem) Stats() []string {
	return []string{
		fmt.Sprintf("%d queued", i.c.queued),
		fmt.Sprintf("%d missing", i.c.missing),
	}
}
