package addseries

import "github.com/charmbracelet/bubbles/key"

type defaultKeyMap struct {
	Quit   key.Binding
	Back   key.Binding
	Help   key.Binding
	Select key.Binding
	Toggle key.Binding
}

var DefaultKeyMap = defaultKeyMap{
	Quit:   key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
	Back:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Help:   key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "close help")),
	Select: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	Toggle: key.NewBinding(key.WithKeys("space"), key.WithHelp("space", "toggle")),
}

func (k defaultKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Select, k.Toggle},
		{k.Back, k.Quit},
		{k.Help},
	}
}
