package search

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit   key.Binding
	Back   key.Binding
	Help   key.Binding
	Select key.Binding
}

var DefaultKeyMap = KeyMap{
	Quit:   key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
	Back:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Help:   key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "close help")),
	Select: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "start search")),
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Select},
		{k.Help, k.Back, k.Quit},
	}
}
