package sonarr_list

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/submarr/internal/tui/styles"
)

// New returns an opinionated list model for sonarr submodels.
func New(title string, items []list.Item, delegate list.ItemDelegate, width, height int) list.Model {
	l := list.New(items, delegate, width, height)

	l.Title = title
	l.DisableQuitKeybindings()
	l.Styles.Title = l.Styles.Title.Copy().
		Background(styles.PurpleColor)
	l.SetShowHelp(false)

	l.FilterInput.Cursor.Style = lipgloss.NewStyle().Foreground(styles.SonarrBlue)
	return l
}
