package addseries

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	sonarrAPI "github.com/jon4hz/submarr/pkg/sonarr"
)

type qualityProfileItem struct {
	qualityProfile *sonarrAPI.QualityProfileResource
	triggered      bool
}

func newQualityProfileItems(qualityProfiles []*sonarrAPI.QualityProfileResource) []list.Item {
	items := make([]list.Item, len(qualityProfiles))
	for i, qualityProfile := range qualityProfiles {
		var triggered bool
		if i == 0 {
			triggered = true
		}

		items[i] = qualityProfileItem{
			qualityProfile: qualityProfile,
			triggered:      triggered,
		}
	}

	return items
}

func (d qualityProfileItem) FilterValue() string { return "" }

type qualityProfileDelegate struct{}

func (d qualityProfileDelegate) Height() int { return 1 }

func (d qualityProfileDelegate) Spacing() int { return 0 }

func (d qualityProfileDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(qualityProfileItem)
	if !ok {
		return
	}

	var text string
	if i.triggered {
		text = fmt.Sprintf("✅ %s", i.qualityProfile.Name)
	} else {
		text = fmt.Sprintf("⬜ %s", i.qualityProfile.Name)
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

	fmt.Fprint(w, title)
}

func (d qualityProfileDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
