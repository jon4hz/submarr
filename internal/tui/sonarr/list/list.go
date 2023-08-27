package sonarr_list

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// New returns an opinionated list model for sonarr submodels.
func New(title string, items []list.Item, delegate list.ItemDelegate, width, height int) list.Model {
	l := list.New(items, delegate, width, height)

	l.Title = title
	l.DisableQuitKeybindings()
	l.Styles.Title = l.Styles.Title.Copy().
		Background(lipgloss.Color("#7B61FF"))
	l.SetShowHelp(false)

	l.FilterInput.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#00CCFF"))
	return l
}
