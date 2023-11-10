package addseries

import "github.com/charmbracelet/bubbles/key"

type defaultKeyMap struct {
	Quit   key.Binding
	Back   key.Binding
	Help   key.Binding
	Select key.Binding
	Toggle key.Binding
	Add    key.Binding
	Up     key.Binding
	Down   key.Binding
}

var DefaultKeyMap = defaultKeyMap{
	Quit:   key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
	Back:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Help:   key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "close help")),
	Select: key.NewBinding(key.WithKeys("enter", " "), key.WithHelp("enter", "select")),
	Toggle: key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "toggle")),
	Add:    key.NewBinding(key.WithKeys("ctrl+a"), key.WithHelp("ctrl+a", "add")),
	Up:     key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k", "up")),
	Down:   key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j", "down")),
}

func (k defaultKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Select, k.Toggle},
		{k.Add, k.Back},
		{k.Help, k.Quit},
	}
}
