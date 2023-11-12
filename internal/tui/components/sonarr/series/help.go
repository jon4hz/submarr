package series

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	CursorUp            key.Binding
	CursorDown          key.Binding
	Quit                key.Binding
	Back                key.Binding
	Help                key.Binding
	Select              key.Binding
	Tab                 key.Binding
	ShiftTab            key.Binding
	Reload              key.Binding
	ToggleMonitor       key.Binding
	ToggleMonitorSeries key.Binding
	Refresh             key.Binding
	AutomaticSearch     key.Binding
	AutomaticSearchAll  key.Binding
	InteractiveSearch   key.Binding
	Delete              key.Binding
}

var DefaultKeyMap = KeyMap{
	CursorUp:            key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	CursorDown:          key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Quit:                key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit")),
	Back:                key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Help:                key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "close help")),
	Select:              key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	Tab:                 key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next focus")),
	ShiftTab:            key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev focus")),
	Reload:              key.NewBinding(key.WithKeys("r", "f5"), key.WithHelp("r", "reload")),
	ToggleMonitor:       key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "toggle monitor")),
	ToggleMonitorSeries: key.NewBinding(key.WithKeys("M"), key.WithHelp("shift+m", "toggle monitor for series")),
	Refresh:             key.NewBinding(key.WithKeys("ctrl+r"), key.WithHelp("ctrl+r", "refresh & scan")),
	AutomaticSearch:     key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "search season")),
	AutomaticSearchAll:  key.NewBinding(key.WithKeys("ctrl+a"), key.WithHelp("ctrl+a", "search all seasons")),
	InteractiveSearch:   key.NewBinding(key.WithKeys("ctrl+w"), key.WithHelp("ctrl+w", "interactive search season")),
	Delete:              key.NewBinding(key.WithKeys("ctrl+d"), key.WithHelp("ctrl+d", "delete series")),
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Tab, k.CursorUp, k.CursorDown, k.Select},
		{k.Reload, k.ToggleMonitorSeries, k.ToggleMonitor, k.Refresh},
		{k.AutomaticSearchAll, k.AutomaticSearch, k.InteractiveSearch, k.Delete},
		{k.Help, k.Back, k.Quit},
	}
}
