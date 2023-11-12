package overview

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
	Reload     key.Binding
	Filter     key.Binding
	AddNew     key.Binding
}

var DefaultKeyMap = KeyMap{
	CursorUp:   key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	CursorDown: key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	NextPage:   key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "next page")),
	PrevPage:   key.NewBinding(key.WithKeys("left", "←"), key.WithHelp("←/h", "prev page")),
	Quit:       key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit")),
	Back:       key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Help:       key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "close help")),
	Select:     key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select series")),
	Reload:     key.NewBinding(key.WithKeys("r", "f5"), key.WithHelp("r", "reload list")),
	Filter:     key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
	AddNew:     key.NewBinding(key.WithKeys("ctrl+a"), key.WithHelp("ctrl+a", "add new series")),
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.CursorUp, k.CursorDown, k.NextPage, k.PrevPage},
		{k.Filter, k.Select, k.Reload},
		{k.Help, k.Back, k.Quit},
		{k.AddNew},
	}
}
