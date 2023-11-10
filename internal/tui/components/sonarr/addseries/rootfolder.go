package addseries

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dustin/go-humanize"
	sonarrAPI "github.com/jon4hz/submarr/pkg/sonarr"
)

type rootFolderItem struct {
	rootFolder *sonarrAPI.RootFolderResource
	triggered  bool
}

func newRootFolderItems(rootFolders []*sonarrAPI.RootFolderResource) []list.Item {
	items := make([]list.Item, len(rootFolders))
	for i, rootFolder := range rootFolders {
		var triggered bool
		if i == 0 {
			triggered = true
		}

		items[i] = rootFolderItem{
			rootFolder: rootFolder,
			triggered:  triggered,
		}
	}

	return items
}

func (d rootFolderItem) FilterValue() string { return "" }

type rootFolderDelegate struct{}

func (d rootFolderDelegate) Height() int { return 1 }

func (d rootFolderDelegate) Spacing() int { return 0 }

var itemStyles = list.NewDefaultItemStyles()

func (d rootFolderDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(rootFolderItem)
	if !ok {
		return
	}

	var text string
	if i.triggered {
		text = fmt.Sprintf("✅ %s (%s free)", i.rootFolder.Path, humanize.IBytes(uint64(i.rootFolder.FreeSpace)))
	} else {
		text = fmt.Sprintf("⬜ %s (%s free)", i.rootFolder.Path, humanize.IBytes(uint64(i.rootFolder.FreeSpace)))
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

func (d rootFolderDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
