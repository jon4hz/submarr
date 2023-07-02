package series

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	CursorUp      key.Binding
	CursorDown    key.Binding
	Quit          key.Binding
	Back          key.Binding
	Help          key.Binding
	Select        key.Binding
	Tab           key.Binding
	ShiftTab      key.Binding
	Reload        key.Binding
	ToggleMonitor key.Binding
}

var DefaultKeyMap = KeyMap{
	CursorUp:      key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	CursorDown:    key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Quit:          key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit")),
	Back:          key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Help:          key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "close help")),
	Select:        key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	Tab:           key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next focus")),
	ShiftTab:      key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev focus")),
	Reload:        key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "reload")),
	ToggleMonitor: key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "toggle monitor")),
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Tab, k.CursorUp, k.CursorDown},
		{k.Select, k.Reload, k.ToggleMonitor},
		{k.Help, k.Back, k.Quit},
	}
}
