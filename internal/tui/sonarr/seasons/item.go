package seasons

import (
	"fmt"

	sonarrAPI "github.com/jon4hz/subrr/pkg/sonarr"
)

type SeasonItem struct {
	Season sonarrAPI.SeasonResource
}

func (s SeasonItem) FilterValue() string {
	return fmt.Sprintf("Season %d", s.Season.SeasonNumber)
}

func NewItem(season sonarrAPI.SeasonResource) SeasonItem {
	return SeasonItem{season}
}
