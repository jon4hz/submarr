package addseries

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	sonarrAPI "github.com/jon4hz/submarr/pkg/sonarr"
)

type monitorItem struct {
	monitor   sonarrAPI.MonitorType
	triggered bool
}

func newMonitorItems() []list.Item {
	items := [...]sonarrAPI.MonitorType{
		sonarrAPI.All,
		sonarrAPI.Future,
		sonarrAPI.Missing,
		sonarrAPI.Existing,
		sonarrAPI.FirstSeason,
		sonarrAPI.LastSeason,
		sonarrAPI.Pilot,
		sonarrAPI.None,
	}

	listItems := make([]list.Item, len(items))
	for i, item := range items {
		var triggered bool
		if i == 0 {
			triggered = true
		}
		listItems[i] = monitorItem{
			monitor:   item,
			triggered: triggered,
		}
	}
	return listItems
}

func (d monitorItem) FilterValue() string { return "" }

type monitorDelegate struct{}

func (d monitorDelegate) Height() int { return 1 }

func (d monitorDelegate) Spacing() int { return 0 }

func (d monitorDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(monitorItem)
	if !ok {
		return
	}

	var text string
	if i.triggered {
		text = fmt.Sprintf("✅ %s", i.monitor)
	} else {
		text = fmt.Sprintf("⬜ %s", i.monitor)
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

func (d monitorDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
