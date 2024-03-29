package common

import tea "github.com/charmbracelet/bubbletea"

type TabModel interface {
	SubModelBase

	Title() string
	Update(msg tea.Msg) (TabModel, tea.Cmd)
}

type SubModel interface {
	SubModelBase

	Update(msg tea.Msg) (SubModel, tea.Cmd)
}

type SubModelBase interface {
	// some bubbleteaish Model methods
	Init() tea.Cmd
	View() string

	// some custom methods
	// SetSize sets the width and height of the model
	SetSize(width, height int)

	// Quit returns true if the model should quit
	Quit() bool

	// Back returns true if the model should go back
	Back() bool
}

// EmbedableModel is a model that can be embedded in other models.
// It provides some basic functionality that is common to all models.
type EmbedableModel struct {
	Width  int
	Height int

	IsQuit bool
	IsBack bool
}

func (m EmbedableModel) Quit() bool { return m.IsQuit }

func (m EmbedableModel) Back() bool { return m.IsBack }
