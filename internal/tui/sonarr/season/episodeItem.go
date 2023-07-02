package season

import (
	"sort"

	"github.com/charmbracelet/bubbles/list"
	sonarrAPI "github.com/jon4hz/subrr/pkg/sonarr"
)

type EpisodeItem struct {
	episode *sonarrAPI.EpisodeResource
}

func (e EpisodeItem) FilterValue() string {
	return e.episode.Title
}

func NewEpisodeItem(episode *sonarrAPI.EpisodeResource) EpisodeItem {
	return EpisodeItem{episode}
}

func episodeToItems(episodes []*sonarrAPI.EpisodeResource) []list.Item {
	sort.Slice(episodes, func(i, j int) bool {
		return episodes[i].EpisodeNumber > episodes[j].EpisodeNumber
	})
	items := make([]list.Item, len(episodes))
	for i, episode := range episodes {
		items[i] = NewEpisodeItem(episode)
	}
	return items
}
