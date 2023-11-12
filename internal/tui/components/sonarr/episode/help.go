package episode

import "github.com/charmbracelet/bubbles/key"

type defaultKeyMap struct {
	Quit key.Binding
	Back key.Binding
	Help key.Binding
	Up   key.Binding
	Down key.Binding
}

var DefaultKeyMap = defaultKeyMap{
	Quit: key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
	Back: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Help: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "close help")),
	Up:   key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k", "up")),
	Down: key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j", "down")),
}

func (k defaultKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back},
		{k.Help, k.Quit},
	}
}
