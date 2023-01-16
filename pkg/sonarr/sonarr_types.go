package sonarr

import "time"

// Ping is the response from the ping endpoint
type Ping struct {
	Status string `json:"status"`
}

// Series is the response from the series endpoint
type Serie struct {
	ID                int32                      `json:"id"`
	Title             string                     `json:"title"`
	AlternateTitles   []AlternativeTitleResource `json:"alternateTitles"`
	SortTitle         string                     `json:"sortTitle"`
	Status            SeriesStatusType           `json:"status"`
	Ended             bool                       `json:"ended"`
	ProfileName       string                     `json:"profileName"`
	Overview          string                     `json:"overview"`
	NextAiring        time.Time                  `json:"nextAiring"`
	PreviousAiring    time.Time                  `json:"previousAiring"`
	Network           string                     `json:"network"`
	AirTime           string                     `json:"airTime"`
	Images            []MediaCover               `json:"images"`
	OriginalLanguage  *Language                  `json:"originalLanguage"`
	RemotePoster      string                     `json:"remotePoster"`
	Seasons           []SeasonResource           `json:"seasons"`
	Year              int32                      `json:"year"`
	Path              string                     `json:"path"`
	QualityProfileID  int32                      `json:"qualityProfileId"`
	SeasonFolder      bool                       `json:"seasonFolder"`
	Monitored         bool                       `json:"monitored"`
	UseSceneNumbering bool                       `json:"useSceneNumbering"`
	Runtime           int32                      `json:"runtime"`
	TVDBID            int32                      `json:"tvdbId"`
	TVRageID          int32                      `json:"tvRageId"`
	TVMAZEID          int32                      `json:"tvMazeId"`
	FirstAired        time.Time                  `json:"firstAired"`
	SeriesType        SeriesType                 `json:"seriesType"`
	CleanTitle        string                     `json:"cleanTitle"`
	ImdbID            string                     `json:"imdbId"`
	TitleSlug         string                     `json:"titleSlug"`
	RootFolderPath    string                     `json:"rootFolderPath"`
	Folder            string                     `json:"folder"`
	Certification     string                     `json:"certification"`
	Genres            []string                   `json:"genres"`
	Tags              []int32                    `json:"tags"`
	Added             time.Time                  `json:"added"`
	AddOptions        *AddSeriesOptions          `json:"addOptions"`
	Statistics        *SeriesStatisticsResource  `json:"statistics"`
	EpisodesChanged   bool                       `json:"episodesChanged"`
}

type AlternativeTitleResource struct {
	Title             string `json:"title"`
	SeasonNumber      int32  `json:"seasonNumber"`
	SceneSeasonNumber int32  `json:"sceneSeasonNumber"`
	SceneOrigin       string `json:"sceneOrigin"`
	Comment           string `json:"comment"`
}

type SeriesStatusType string

const (
	Continuing SeriesStatusType = "continuing"
	Ended      SeriesStatusType = "ended"
	Upcoming   SeriesStatusType = "upcoming"
	Deleted    SeriesStatusType = "deleted"
)

type MediaCover struct {
	CoverType MediaCoverType `json:"coverType"`
	URL       string         `json:"url"`
	RemoteURL string         `json:"remoteUrl"`
}

type MediaCoverType string

const (
	UnknownMediaCoverType MediaCoverType = "unknown"
	Poster                MediaCoverType = "poster"
	Banner                MediaCoverType = "banner"
	Fanart                MediaCoverType = "fanart"
	Screenshot            MediaCoverType = "screenshot"
	Headshot              MediaCoverType = "headshot"
)

type Language struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type SeasonResource struct {
	SeasonNumber int32                     `json:"seasonNumber"`
	Monitored    bool                      `json:"monitored"`
	Statistics   *SeasonStatisticsResource `json:"statistics"`
	Images       []MediaCover              `json:"images"`
}

type SeasonStatisticsResource struct {
	NextAiring        time.Time `json:"nextAiring"`
	PreviousAiring    time.Time `json:"previousAiring"`
	EpisodeFileCount  int32     `json:"episodeFileCount"`
	EpisodeCount      int32     `json:"episodeCount"`
	TotalEpisodeCount int32     `json:"totalEpisodeCount"`
	SizeOnDisk        int64     `json:"sizeOnDisk"`
	ReleaseGroups     []string  `json:"releaseGroups"`
	PercentOfEpisodes float64   `json:"percentOfEpisodes"`
}

type SeriesType string

const (
	Standard SeriesType = "standard"
	Daily    SeriesType = "daily"
	Anime    SeriesType = "anime"
)

type AddSeriesOptions struct {
	IgnoreEpisodesWithFiles      bool        `json:"ignoreEpisodesWithFiles"`
	IgnoreEpisodesWithoutFiles   bool        `json:"ignoreEpisodesWithoutFiles"`
	Monitor                      MonitorType `json:"monitor"`
	SearchForMissingEpisodes     bool        `json:"searchForMissingEpisodes"`
	SearchForCutoffUnmetEpisodes bool        `json:"searchForCutoffUnmetEpisodes"`
}

type MonitorType string

const (
	UnknownMonitorType MonitorType = "unknown"
	All                MonitorType = "all"
	Future             MonitorType = "future"
	Missing            MonitorType = "missing"
	Existing           MonitorType = "existing"
	FirstSeason        MonitorType = "firstSeason"
	LastSeason         MonitorType = "lastSeason"
	Pilot              MonitorType = "pilot"
	None               MonitorType = "none"
)

type Ratings struct {
	Votes int32   `json:"votes"`
	Value float64 `json:"value"`
}

type SeriesStatisticsResource struct {
	SeasonCount       int32    `json:"seasonCount"`
	EpisodeFileCount  int32    `json:"episodeFileCount"`
	EpisodeCount      int32    `json:"episodeCount"`
	TotalEpisodeCount int32    `json:"totalEpisodeCount"`
	SizeOnDisk        int64    `json:"sizeOnDisk"`
	ReleaseGroups     []string `json:"releaseGroups"`
	PercentOfEpisodes float64  `json:"percentOfEpisodes"`
}
