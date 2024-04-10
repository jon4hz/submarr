package sonarr

import (
	"context"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jon4hz/submarr/internal/httpclient"
	"github.com/jon4hz/submarr/internal/logging"
	sonarrAPI "github.com/jon4hz/submarr/pkg/sonarr"
)

type EpisodeHistoryResult struct {
	Error   error
	Episode *sonarrAPI.EpisodeResource
	History []*sonarrAPI.HistoryResource
}

func (c *Client) GetEpisodeHistory(episode *sonarrAPI.EpisodeResource) tea.Cmd {
	return func() tea.Msg {
		res, err := c.sonarr.GetHistory(context.Background(),
			httpclient.WithPage(1),
			httpclient.WithPageSize(1000),
			httpclient.WithSortKey("date"),
			httpclient.WithSortDirection(httpclient.Descending),
			httpclient.WithParams(map[string]string{"episodeId": strconv.Itoa(int(episode.ID))}),
		)
		if err != nil {
			logging.Log.Error("Failed to fetch episode history", "id", strconv.Itoa(int(episode.ID)), "series", episode.SeriesTitle, "err", err)
			return EpisodeHistoryResult{Error: err}
		}
		return EpisodeHistoryResult{
			Episode: episode,
			History: res.Records,
		}
	}
}

type EpisodeDeleteResult struct {
	Error error
}

func (c *Client) DeleteEpisodeFile(episode *sonarrAPI.EpisodeResource) tea.Cmd {
	return func() tea.Msg {
		err := c.sonarr.DeleteEpisodeFile(context.Background(), episode.EpisodeFileID)
		if err != nil {
			logging.Log.Error("Failed to delete episode file", "id", strconv.Itoa(int(episode.ID)), "series", episode.SeriesTitle, "err", err)
			return EpisodeDeleteResult{Error: err}
		}
		return EpisodeDeleteResult{}
	}
}
