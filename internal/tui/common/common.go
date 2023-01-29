package common

import (
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const Ellipsis = "…"

type ErrMsg struct {
	Description string
}

func (e ErrMsg) Error() string {
	return e.Description
}

func NewErrMsg(description string) ErrMsg {
	return ErrMsg{Description: description}
}

func NewErrCmd(description string) tea.Cmd {
	return func() tea.Msg {
		return NewErrMsg(description)
	}
}

func NewErrCmds(descriptions ...string) []tea.Cmd {
	var cmds []tea.Cmd
	for _, description := range descriptions {
		if description == "" {
			continue
		}
		cmds = append(cmds, NewErrCmd(description))
	}
	return cmds
}

func Title(s string) string {
	return cases.Title(language.English).String(s)
}
