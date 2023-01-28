package core

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jon4hz/subrr/internal/logging"
)

type ClientsItem interface {
	fmt.Stringer

	Render(itemWidth int, isSelected bool) string
	FilterValue() string
}

type FetchClientsMsg struct {
	Items  []list.Item
	Errors []string
}

type FetchClientsErrorMsg struct {
	Description string
}

func (c *Client) FetchClients() tea.Cmd {
	return func() tea.Msg {
		var (
			items  []list.Item
			errors []string
		)
		if c.Sonarr != nil {
			if err := c.Sonarr.Init(); err != nil {
				logging.Log.Error().Err(err).Msg("Failed to initialize sonarr")
				errors = append(errors, "Failed to initialize sonarr")
			}
			items = append(items, c.Sonarr.ListItem())
		}
		if c.Radarr != nil {
			if err := c.Radarr.Init(); err != nil {
				logging.Log.Error().Err(err).Msg("Failed to initialize radarr")
				errors = append(errors, "Failed to initialize radarr")
			}
			items = append(items, c.Radarr.ListItem())
		}

		return FetchClientsMsg{
			Items:  items,
			Errors: errors,
		}
	}
}
