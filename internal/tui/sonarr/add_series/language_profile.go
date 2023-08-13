package addseries

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	sonarrAPI "github.com/jon4hz/subrr/pkg/sonarr"
)

type languageProfileItem struct {
	languageProfile *sonarrAPI.LanguageProfileResource
	triggered       bool
}

func newLanguageProfileItems(languageProfiles []*sonarrAPI.LanguageProfileResource) []list.Item {
	items := make([]list.Item, len(languageProfiles))
	for i, languageProfile := range languageProfiles {
		var triggered bool
		if i == 0 {
			triggered = true
		}

		items[i] = languageProfileItem{
			languageProfile: languageProfile,
			triggered:       triggered,
		}
	}

	return items
}

func (d languageProfileItem) FilterValue() string { return "" }

type languageProfileDelegate struct{}

func (d languageProfileDelegate) Height() int { return 1 }

func (d languageProfileDelegate) Spacing() int { return 0 }

func (d languageProfileDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(languageProfileItem)
	if !ok {
		return
	}

	var text string
	if i.triggered {
		text = fmt.Sprintf("✅ %s", i.languageProfile.Name)
	} else {
		text = fmt.Sprintf("⬜ %s", i.languageProfile.Name)
	}

	var (
		isSelected = index == m.Index()
		title      string
	)

	if isSelected {
		title = itemStyles.SelectedTitle.Render(text)
	} else {
		title = itemStyles.NormalTitle.Render(text)
	}

	fmt.Fprintf(w, "%s", title)
}

func (d languageProfileDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
