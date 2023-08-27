package seasons

import (
	"fmt"

	sonarrAPI "github.com/jon4hz/submarr/pkg/sonarr"
)

type SeasonItem struct {
	Index  int
	Season *sonarrAPI.SeasonResource
}

func (s SeasonItem) FilterValue() string {
	return fmt.Sprintf("Season %d", s.Season.SeasonNumber)
}

func NewItem(index int, season *sonarrAPI.SeasonResource) SeasonItem {
	return SeasonItem{index, season}
}
