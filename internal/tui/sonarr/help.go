package sonarr

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	CursorUp   key.Binding
	CursorDown key.Binding
	NextPage   key.Binding
	PrevPage   key.Binding
	Quit       key.Binding
	Help       key.Binding
	Select     key.Binding
	Reload     key.Binding
}

var DefaultKeyMap = KeyMap{
	CursorUp:   key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	CursorDown: key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	NextPage:   key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "next page")),
	PrevPage:   key.NewBinding(key.WithKeys("left", "←"), key.WithHelp("←/h", "prev page")),
	Quit:       key.NewBinding(key.WithKeys("q", "esc"), key.WithHelp("q/esc", "quit")),
	Help:       key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "close help")),
	Select:     key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select client")),
	Reload:     key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "reload list")),
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.CursorUp, k.CursorDown, k.NextPage, k.PrevPage},
		{k.Select, k.Reload},
		{k.Help, k.Quit},
	}
}
