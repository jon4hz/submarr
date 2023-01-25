package core

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ClientsItem interface {
	fmt.Stringer

	Render(itemWidth int, isSelected bool) string
	FilterValue() string
}

type FetchClientsSuccessMsg struct {
	Items []list.Item
}

type FetchClientsErrorMsg struct {
	Description string
	Err         error
}

func (c *Client) FetchClients() tea.Cmd {
	return func() tea.Msg {
		var items []list.Item
		if c.Sonarr != nil {
			if err := c.Sonarr.Init(); err != nil {
				return FetchClientsErrorMsg{Description: "Failed to initialize sonarr", Err: err}
			}
			items = append(items, c.Sonarr.ListItem(), c.Sonarr.ListItem(), c.Sonarr.ListItem())
		}

		return FetchClientsSuccessMsg{items}
	}
}
