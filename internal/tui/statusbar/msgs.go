package statusbar

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	SetTitleMsg struct {
		title      string
		foreground lipgloss.TerminalColor
		background lipgloss.TerminalColor
	}

	SetMessageMsg struct {
		message string
		timeout int
	}

	SetErrMsg struct {
		Description string
	}

	SetHelpMsg [][]key.Binding
)

type SetTitleOpts func(*SetTitleMsg)

func WithTitleForeground(fg lipgloss.TerminalColor) SetTitleOpts {
	return func(m *SetTitleMsg) {
		m.foreground = fg
	}
}

func WithTitleBackground(bg lipgloss.TerminalColor) SetTitleOpts {
	return func(m *SetTitleMsg) {
		m.background = bg
	}
}

func NewTitleMsg(title string, opts ...SetTitleOpts) SetTitleMsg {
	m := SetTitleMsg{
		title:      title,
		foreground: nil,
		background: nil,
	}
	for _, opt := range opts {
		opt(&m)
	}
	return m
}

func NewTitleCmd(title string, opts ...SetTitleOpts) tea.Cmd {
	return func() tea.Msg {
		return NewTitleMsg(title, opts...)
	}
}

type SetMessageOpts func(*SetMessageMsg)

func WithMessageTimeout(timeout int) SetMessageOpts {
	return func(m *SetMessageMsg) {
		m.timeout = timeout
	}
}

func NewMessageMsg(message string, opts ...SetMessageOpts) SetMessageMsg {
	m := SetMessageMsg{
		message: message,
		timeout: 4,
	}
	for _, opt := range opts {
		opt(&m)
	}
	return m
}

func NewMessageCmd(message string, opts ...SetMessageOpts) tea.Cmd {
	return func() tea.Msg {
		return NewMessageMsg(message, opts...)
	}
}

func NewHelpMsg(h [][]key.Binding) SetHelpMsg {
	return SetHelpMsg(h)
}

func NewHelpCmd(h [][]key.Binding) tea.Cmd {
	return func() tea.Msg {
		return NewHelpMsg(h)
	}
}

func (e SetErrMsg) Error() string {
	return e.Description
}

func NewErrMsg(description string) SetErrMsg {
	return SetErrMsg{Description: description}
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

type msgQueueMsg struct {
	Message string
	isError bool
	timeout time.Duration
}

func newMsgQueueMsg(msg string, timeout int, isError bool) *msgQueueMsg {
	return &msgQueueMsg{
		Message: msg,
		isError: isError,
		timeout: time.Second * time.Duration(timeout),
	}
}

type messageTimeoutMsg struct{}
