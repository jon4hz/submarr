package season

import (
	"sort"

	"github.com/charmbracelet/bubbles/list"
	sonarrAPI "github.com/jon4hz/submarr/pkg/sonarr"
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
		// Check if episode is in download queue.
		// If the episode is in the queue multiple times,
		// the one with the most download progress is used.
		var queueItem *sonarrAPI.QueueResource
		for _, q := range queue {
			if episode.ID == q.EpisodeID {
				if queueItem != nil {
					np := (q.Size - q.Sizeleft) / q.Size
					op := (queueItem.Size - queueItem.Sizeleft) / queueItem.Size
					if np > op {
						queueItem = q
					}
				} else {
					queueItem = q
				}
			}
		}
		items[i] = NewEpisodeItem(episode, queueItem)
	}
	return items
}
