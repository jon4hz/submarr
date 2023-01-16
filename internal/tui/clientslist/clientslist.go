package clientslist

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jon4hz/subrr/internal/core"
)

type Model struct {
	client *core.Client
}

type ItemsMsg struct {
	Items []list.Item
}

func FetchClientsListItems(c *core.Client) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(time.Second * 3)
		return ItemsMsg{}
	}
}
