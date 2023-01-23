package core

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type ClientsItem interface {
	fmt.Stringer

	Render(itemWidth int) string
	FilterValue() string
}

type FetchClientsSuccessMsg struct {
	Items []ClientsItem
}

type FetchClientsErrorMsg struct {
	Description string
	Err         error
}

func (c *Client) FetchClients() tea.Cmd {
	return func() tea.Msg {
		var items []ClientsItem
		if c.Sonarr != nil {
			if err := c.Sonarr.Init(); err != nil {
				return FetchClientsErrorMsg{Description: "Failed to initialize sonarr", Err: err}
			}
			items = append(items, c.Sonarr.ListItem())
		}

		return FetchClientsSuccessMsg{items}
	}
}
