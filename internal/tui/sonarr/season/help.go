package season

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	CursorUp   key.Binding
	CursorDown key.Binding
	Quit       key.Binding
	Back       key.Binding
	Help       key.Binding
	Select     key.Binding
	Reload     key.Binding
}

var DefaultKeyMap = KeyMap{
	CursorUp:   key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	CursorDown: key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Quit:       key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit")),
	Back:       key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Help:       key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "close help")),
	Select:     key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	Reload:     key.NewBinding(key.WithKeys("r", "f5"), key.WithHelp("r", "reload")),
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.CursorUp, k.CursorDown, k.Select},
		{k.Reload},
		{k.Help, k.Back, k.Quit},
	}
}
