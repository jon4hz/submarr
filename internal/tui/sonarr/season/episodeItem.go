package season

import (
	"sort"

	"github.com/charmbracelet/bubbles/list"
	sonarrAPI "github.com/jon4hz/subrr/pkg/sonarr"
)

type EpisodeItem struct {
	episode *sonarrAPI.EpisodeResource
	queue   *sonarrAPI.QueueResource
}

func (e EpisodeItem) FilterValue() string {
	return e.episode.Title
}

func NewEpisodeItem(episode *sonarrAPI.EpisodeResource, queue *sonarrAPI.QueueResource) EpisodeItem {
	return EpisodeItem{episode, queue}
}

func episodeToItems(episodes []*sonarrAPI.EpisodeResource, queue []*sonarrAPI.QueueResource) []list.Item {
	sort.Slice(episodes, func(i, j int) bool {
		return episodes[i].EpisodeNumber > episodes[j].EpisodeNumber
	})

	items := make([]list.Item, len(episodes))
	for i, episode := range episodes {
		// check if episode is in download queue
		var queueItem *sonarrAPI.QueueResource
		for _, q := range queue {
			if episode.ID == q.EpisodeID {
				queueItem = q
			}
		}
		items[i] = NewEpisodeItem(episode, queueItem)
	}
	return items
}
