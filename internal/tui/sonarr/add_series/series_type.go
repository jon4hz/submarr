package addseries

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	sonarrAPI "github.com/jon4hz/subrr/pkg/sonarr"
)

type seriesTypeItem struct {
	seriesType sonarrAPI.SeriesType
	triggered  bool
}

func newSeriesTypeItems() []list.Item {
	items := [...]sonarrAPI.SeriesType{
		sonarrAPI.Standard,
		sonarrAPI.Daily,
		sonarrAPI.Anime,
	}

	listItems := make([]list.Item, len(items))
	for i, item := range items {
		var triggered bool
		if i == 0 {
			triggered = true
		}
		listItems[i] = seriesTypeItem{
			seriesType: item,
			triggered:  triggered,
		}
	}
	return listItems
}

func (d seriesTypeItem) FilterValue() string { return "" }

type seriesTypeDelegate struct{}

func (d seriesTypeDelegate) Height() int { return 1 }

func (d seriesTypeDelegate) Spacing() int { return 0 }

func (d seriesTypeDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(seriesTypeItem)
	if !ok {
		return
	}

	var text string
	if i.triggered {
		text = fmt.Sprintf("✅ %s", i.seriesType)
	} else {
		text = fmt.Sprintf("⬜ %s", i.seriesType)
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

func (d seriesTypeDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
