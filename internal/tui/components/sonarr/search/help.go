package search

import "github.com/charmbracelet/bubbles/key"

type inputKeyMap struct {
	Quit   key.Binding
	Back   key.Binding
	Help   key.Binding
	Select key.Binding
}

var InputKeyMap = inputKeyMap{
	Quit:   key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
	Back:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Help:   key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "close help")),
	Select: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "start search")),
}

func (k inputKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Select, k.Back},
		{k.Help, k.Quit},
	}
}

type resultKeyMap struct {
	Quit   key.Binding
	Back   key.Binding
	Help   key.Binding
	Select key.Binding
	Filter key.Binding
}

var ResultKeyMap = resultKeyMap{
	Quit:   key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
	Back:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Help:   key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "close help")),
	Select: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	Filter: key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
}

func (k resultKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Select, k.Filter},
		{k.Help, k.Back, k.Quit},
	}
}
