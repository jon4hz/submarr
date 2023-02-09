package serie

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	CursorUp   key.Binding
	CursorDown key.Binding
	NextPage   key.Binding
	PrevPage   key.Binding
	Quit       key.Binding
	Back       key.Binding
	Help       key.Binding
	Select     key.Binding
	Tab        key.Binding
	ShiftTab   key.Binding
}

var DefaultKeyMap = KeyMap{
	CursorUp:   key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	CursorDown: key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Quit:       key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit")),
	Back:       key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Help:       key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "close help")),
	Select:     key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	Tab:        key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next focus")),
	ShiftTab:   key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev focus")),
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.CursorUp, k.CursorDown, k.NextPage, k.PrevPage},
		{k.Select, k.Tab},
		{k.Help, k.Back, k.Quit},
	}
}
