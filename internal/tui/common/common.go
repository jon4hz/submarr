package common

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const Ellipsis = "…"

func Title(s string) string {
	return cases.Title(language.English).String(s)
}

type Spinner struct {
	spinner.Model
	Message string
}

func NewSpinner() Spinner {
	s := Spinner{
		Message: GetRandomLoadingMessage(),
	}
	s.Model = spinner.New(spinner.WithSpinner(spinner.Points))
	return s
}

func (s Spinner) View() string {
	return s.Model.View() + "  " + s.Message
}

var SubtileColor = lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}
