package statusbar

import (
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/subrr/internal/tui/common"
	zone "github.com/lrstanley/bubblezone"
	"github.com/muesli/reflow/truncate"
)

var (
	statusBarNoteFg = lipgloss.AdaptiveColor{Light: "#656565", Dark: "#7D7D7D"}
	statusBarBg     = lipgloss.AdaptiveColor{Light: "#E6E6E6", Dark: "#242424"}

	titleStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Height(1).
			Bold(true)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(statusBarNoteFg).
			Background(statusBarBg).
			Height(1).
			Padding(0, 1)

	messageStyle = statusBarStyle.Copy().
			Foreground(lipgloss.AdaptiveColor{Light: "#89F0CB", Dark: "#89F0CB"}).
			Background(statusBarBg)

	errorStyle = statusBarStyle.Copy().
			Foreground(lipgloss.AdaptiveColor{Light: "#F08F89", Dark: "#F08F89"}).
			Background(statusBarBg)

	helpStyle = lipgloss.NewStyle().
			Height(1).
			Padding(0, 1).
			Foreground(statusBarNoteFg).
			Background(lipgloss.AdaptiveColor{Light: "#DCDCDC", Dark: "#323232"})

	helpViewStyle = lipgloss.NewStyle().
			Foreground(statusBarNoteFg).
			Background(lipgloss.AdaptiveColor{Light: "#f2f2f2", Dark: "#1B1B1B"}).
			Padding(1, 2)
)

type Model struct {
	isInitialized bool

	width  int
	height int

	showHelp bool

	Title           string
	TitleForeground lipgloss.TerminalColor
	TitleBackground lipgloss.TerminalColor

	msgQueue     []msgQueueMsg
	msgQueueMu   *sync.Mutex
	message      *msgQueueMsg
	messageTimer *time.Timer
	placeholder  string

	help     help.Model
	helpKeys [][]key.Binding
}

func New(title string) Model {
	m := Model{
		Title:       title,
		placeholder: "Bliblablub Placeholder",
		msgQueue:    make([]msgQueueMsg, 0),
		msgQueueMu:  &sync.Mutex{},
		help:        help.New(),
	}
	m.help.ShowAll = true
	m.help.FullSeparator = "      "

	helpStyle := help.Styles{
		ShortKey:       lipgloss.NewStyle(),
		ShortDesc:      lipgloss.NewStyle(),
		ShortSeparator: lipgloss.NewStyle(),
		Ellipsis:       lipgloss.NewStyle(),
		FullKey:        lipgloss.NewStyle().PaddingRight(1),
		FullDesc:       lipgloss.NewStyle(),
		FullSeparator:  lipgloss.NewStyle(),
	}
	m.help.Styles = helpStyle

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "?":
			m.showHelp = !m.showHelp
		}
		return m, nil

	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseLeft:
			if zone.Get("toggle-help").InBounds(msg) {
				m.showHelp = !m.showHelp
			}
			return m, nil
		}

	case common.ErrMsg:
		m.msgQueueMu.Lock()
		m.msgQueue = append(m.msgQueue, newMsgQueueMsg(msg.Error()+" :(", 5, true))
		m.msgQueueMu.Unlock()

	case SetMessageMsg:
		m.msgQueueMu.Lock()
		m.msgQueue = append(m.msgQueue, newMsgQueueMsg(msg.message, msg.timeout, false))
		m.msgQueueMu.Unlock()

	case dispatchMsgQueueMsg:
		m.msgQueueMu.Lock()
		defer m.msgQueueMu.Unlock()

		// if the queue is empty, we set the message to nil and return
		if len(m.msgQueue) == 0 {
			m.message = nil
			return m, nil
		}

		// get the next message from the queue and set it as the current message
		im := m.msgQueue[0]
		m.message = &im
		// remove the message from the queue
		m.msgQueue = m.msgQueue[1:]
		// start a timer to remove the message after the timeout and dispatch a new message
		return m, tea.Tick(im.timeout, func(t time.Time) tea.Msg { return dispatchMsgQueueMsg{} })

	case SetTitleMsg:
		m.Title = msg.title
		if msg.foreground != nil {
			m.TitleForeground = msg.foreground
		}
		if msg.background != nil {
			m.TitleBackground = msg.background
		}
		return m, nil

	case SetHelpMsg:
		m.helpKeys = msg
		return m, nil
	}

	m.msgQueueMu.Lock()
	defer m.msgQueueMu.Unlock()
	// if the queue has only one message and no message is currently displayed, we trigger a new dispatcher
	if len(m.msgQueue) == 1 && m.message == nil {
		return m, func() tea.Msg {
			return dispatchMsgQueueMsg{}
		}
	}

	return m, nil
}

func (m *Model) SetWidth(width int) {
	m.width = width
	m.help.Width = width - helpViewStyle.GetHorizontalPadding()
}

func (m Model) GetHeight() int {
	return lipgloss.Height(m.View())
}

func (m Model) IsInitialized() bool { return m.isInitialized }

func (m Model) View() string {
	title := titleStyle.Foreground(m.TitleForeground).Background(m.TitleBackground).Render(m.Title)

	help := zone.Mark("toggle-help", helpStyle.Render("? Help"))

	statusWidth := m.width - lipgloss.Width(title) - lipgloss.Width(help)

	var status string
	if m.message == nil {
		status = statusBarStyle.Width(statusWidth).Render(
			truncate.StringWithTail(m.placeholder, uint(statusWidth-statusBarStyle.GetHorizontalPadding()), common.Ellipsis),
		)
	} else {
		if m.message.isError {
			status = errorStyle.Width(statusWidth).Render(
				truncate.StringWithTail(m.message.Message, uint(statusWidth-errorStyle.GetHorizontalPadding()), common.Ellipsis),
			)
		} else {
			status = messageStyle.Width(statusWidth).Render(
				truncate.StringWithTail(m.message.Message, uint(statusWidth-messageStyle.GetHorizontalPadding()), common.Ellipsis),
			)
		}
	}

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		title,
		status,
		help,
	)

	if m.showHelp {
		return lipgloss.JoinVertical(lipgloss.Top,
			bar,
			helpViewStyle.Width(m.width).Render(m.helpView()),
		)
	}
	return bar
}

func (m Model) helpView() string {
	return m.help.View(m)
}

func (m Model) FullHelp() [][]key.Binding {
	return m.helpKeys
}

func (m Model) ShortHelp() []key.Binding {
	return nil
}
