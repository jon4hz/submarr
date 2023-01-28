package statusbar

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
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

func NewHelpMsg(h [][]key.Binding) SetHelpMsg {
	return SetHelpMsg(h)
}

type msgQueueMsg struct {
	Message string
	isError bool
	timeout time.Duration
}

type dispatchMsgQueueMsg struct{}

func newMsgQueueMsg(msg string, timeout int, isError bool) msgQueueMsg {
	return msgQueueMsg{
		Message: msg,
		isError: isError,
		timeout: time.Second * time.Duration(timeout),
	}
}
