package toggle

import tea "github.com/charmbracelet/bubbletea"

type Model struct {
	toggled bool
}

func New(isActive ...bool) Model {
	if len(isActive) > 0 {
		return Model{toggled: isActive[0]}
	}
	return Model{}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeySpace {
			m.toggled = !m.toggled
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.toggled {
		return "✅"
	}

	return "⬜"
}

func (m Model) Toggled() bool {
	return m.toggled
}
