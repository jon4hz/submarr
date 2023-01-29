package sonarr

import (
	"strings"
	"time"

	"github.com/jon4hz/subrr/internal/httpclient"
)

// Ping is the response from the ping endpoint
type Ping struct {
	Status string `json:"status"`
}

// SeriesResource is the response from the series endpoint
type SeriesResource struct {
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

type QueueResourcePagingResource struct {
	Page          int32                    `json:"page"`
	PageSize      int32                    `json:"pageSize"`
	SortKey       string                   `json:"sortKey"`
	SortDirection httpclient.SortDirection `json:"sortDirection"`
	Filters       []PagingResourceFilter   `json:"filters"`
	TotalRecords  int32                    `json:"totalRecords"`
	Records       []QueueResource          `json:"records"`
}

type PagingResourceFilter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type QueueResource struct {
	ID                      int32                          `json:"id"`
	SeriesID                int32                          `json:"seriesId"`
	EpisodeID               int32                          `json:"episodeId"`
	Series                  *SeriesResource                `json:"series"`
	Episode                 *EpisodeResource               `json:"episode"`
	Languages               []Language                     `json:"languages"`
	Quality                 *QualityModel                  `json:"quality"`
	CustomFormats           []CustomFormatResource         `json:"customFormats"`
	Size                    float64                        `json:"size"`
	Title                   string                         `json:"title"`
	Sizeleft                float64                        `json:"sizeleft"`
	Timeleft                TimeLeft                       `json:"timeleft"`
	EstimatedCompletionTime time.Time                      `json:"estimatedCompletionTime"`
	Status                  string                         `json:"status"`
	TrackedDownloadStatus   TrackedDownloadStatus          `json:"trackedDownloadStatus"`
	TrackedDownloadState    TrackedDownloadState           `json:"trackedDownloadState"`
	StatusMessages          []TrackedDownloadStatusMessage `json:"statusMessages"`
	ErrorMessage            string                         `json:"errorMessage"`
	DownloadID              string                         `json:"downloadId"`
	Protocol                DownloadProtocol               `json:"protocol"`
	DownloadClient          string                         `json:"downloadClient"`
	Indexer                 string                         `json:"indexer"`
	OutputPath              string                         `json:"outputPath"`
}

// TimeLeft is a custom type to handle the timeleft field
type TimeLeft time.Time

func (tl *TimeLeft) UnmarshalJSON(b []byte) (err error) {
	value := strings.Trim(string(b), `"`) // get rid of "
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse("15:04:05", value) // parse time
	if err != nil {
		return err
	}
	*tl = TimeLeft(t) // set result using the pointer
	return nil
}

func (tl TimeLeft) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(tl).Format("15:04:05") + `"`), nil
}

type EpisodeResource struct {
	ID                         int32                `json:"id"`
	SeriesID                   int32                `json:"seriesId"`
	TVDBID                     int32                `json:"tvdbId"`
	EpisodeFileID              int32                `json:"episodeFileId"`
	SeasonNumber               int32                `json:"seasonNumber"`
	EpisodeNumber              int32                `json:"episodeNumber"`
	Title                      string               `json:"title"`
	AirDate                    CivilTime            `json:"airDate"`
	AirDateUTC                 time.Time            `json:"airDateUtc"`
	Overview                   string               `json:"overview"`
	EpisodeFile                *EpisodeFileResource `json:"episodeFile"`
	HasFile                    bool                 `json:"hasFile"`
	Monitored                  bool                 `json:"monitored"`
	AbsoluteEpisodeNumber      int32                `json:"absoluteEpisodeNumber"`
	SceneAbsoluteEpisodeNumber int32                `json:"sceneAbsoluteEpisodeNumber"`
	SceneEpisodeNumber         int32                `json:"sceneEpisodeNumber"`
	SceneSeasonNumber          int32                `json:"sceneSeasonNumber"`
	UnverifiedSceneNumbering   bool                 `json:"unverifiedSceneNumbering"`
	EndTime                    time.Time            `json:"endTime"`
	GrabDate                   time.Time            `json:"grabDate"`
	SeriesTitle                string               `json:"seriesTitle"`
	Series                     *SeriesResource      `json:"series"`
	Images                     []MediaCover         `json:"images"`
	Grabbed                    bool                 `json:"grabbed"`
}

// CivilTime implements a custom time format for JSON marshalling/unmarshalling
type CivilTime time.Time

func (c *CivilTime) UnmarshalJSON(b []byte) (err error) {
	value := strings.Trim(string(b), `"`) // get rid of "
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse("2006-01-02", value) // parse time
	if err != nil {
		return err
	}
	*c = CivilTime(t) // set result using the pointer
	return nil
}

func (c CivilTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(c).Format("2006-01-02") + `"`), nil
}

type EpisodeFileResource struct {
	ID                  int32                  `json:"id"`
	SeriesID            int32                  `json:"seriesId"`
	SeasonNumber        int32                  `json:"seasonNumber"`
	RelativePath        string                 `json:"relativePath"`
	Path                string                 `json:"path"`
	Size                int64                  `json:"size"`
	DateAdded           time.Time              `json:"dateAdded"`
	SceneName           string                 `json:"sceneName"`
	ReleaseGroup        string                 `json:"releaseGroup"`
	Languages           []Language             `json:"languages"`
	Quality             *QualityModel          `json:"quality"`
	CustomFormats       []CustomFormatResource `json:"customFormats"`
	MediaInfo           *MediaInfoResource     `json:"mediaInfo"`
	QualityCutoffNotMet bool                   `json:"qualityCutoffNotMet"`
}

type QualityModel struct {
	Quality  *Quality  `json:"quality"`
	Revision *Revision `json:"revision"`
}

type Quality struct {
	ID         int32         `json:"id"`
	Name       string        `json:"name"`
	Source     QualitySource `json:"source"`
	Resolution int32         `json:"resolution"`
}

type QualitySource string

const (
	UnknownQualitySource QualitySource = "unknown"
	Television           QualitySource = "television"
	TelevisionRaw        QualitySource = "televisionRaw"
	Web                  QualitySource = "web"
	WebRip               QualitySource = "webRip"
	Dvd                  QualitySource = "dvd"
	Bluray               QualitySource = "bluray"
	BlurayRaw            QualitySource = "blurayRaw"
)

type Revision struct {
	Version  int32 `json:"version"`
	Real     int32 `json:"real"`
	IsRepack bool  `json:"isRepack"`
}

type CustomFormatResource struct {
	ID                              int32                             `json:"id"`
	Name                            string                            `json:"name"`
	IncludeCustomFormatWhenRenaming bool                              `json:"includeCustomFormatWhenRenaming"`
	Specifications                  []CustomFormatSpecificationSchema `json:"specifications"`
}

type CustomFormatSpecificationSchema struct {
	ID                 int32   `json:"id"`
	Name               string  `json:"name"`
	Implementation     string  `json:"implementation"`
	ImplementationName string  `json:"implementationName"`
	InfoLink           string  `json:"infoLink"`
	Negate             bool    `json:"negate"`
	Required           bool    `json:"required"`
	Fields             []Field `json:"fields"`
	Presets            []any   `json:"presets"`
}

type Field struct {
	Order                       int32          `json:"order"`
	Name                        string         `json:"name"`
	Label                       string         `json:"label"`
	Unit                        string         `json:"unit"`
	HelpText                    string         `json:"helpText"`
	HelpLink                    string         `json:"helpLink"`
	Value                       any            `json:"value"`
	Type                        string         `json:"type"`
	Advanced                    bool           `json:"advanced"`
	SelectOptions               []SelectOption `json:"selectOptions"`
	SelectOptionsProviderAction string         `json:"selectOptionsProviderAction"`
	Section                     string         `json:"section"`
	Hidden                      string         `json:"hidden"`
	Privacy                     PrivacyLevel   `json:"privacy"`
}

type SelectOption struct {
	Value int32  `json:"value"`
	Name  string `json:"name"`
	Order int32  `json:"order"`
	Hint  string `json:"hint"`
}

type PrivacyLevel string

const (
	Normal   PrivacyLevel = "normal"
	Password PrivacyLevel = "password"
	APIKey   PrivacyLevel = "apiKey"
	UserName PrivacyLevel = "userName"
)

type MediaInfoResource struct {
	ID                    int32   `json:"id"`
	AudioBitrate          int64   `json:"audioBitrate"`
	AudioChannels         float64 `json:"audioChannels"`
	AudioCodec            string  `json:"audioCodec"`
	AudioLanguages        string  `json:"audioLanguages"`
	AudioStreamCount      int32   `json:"audioStreamCount"`
	VideoBitDepth         int32   `json:"videoBitDepth"`
	VideoBitrate          int64   `json:"videoBitrate"`
	VideoCodec            string  `json:"videoCodec"`
	VideoFps              float64 `json:"videoFps"`
	VideoDynamicRange     string  `json:"videoDynamicRange"`
	VideoDynamicRangeType string  `json:"videoDynamicRangeType"`
	Resolution            string  `json:"resolution"`
	RunTime               string  `json:"runTime"`
	ScanType              string  `json:"scanType"`
	Subtitles             string  `json:"subtitles"`
}

type TimeSpan struct {
	Ticks             int64   `json:"ticks"`
	Days              int32   `json:"days"`
	Hours             int32   `json:"hours"`
	Milliseconds      int32   `json:"milliseconds"`
	Minutes           int32   `json:"minutes"`
	Seconds           int32   `json:"seconds"`
	TotalDays         float64 `json:"totalDays"`
	TotalHours        float64 `json:"totalHours"`
	TotalMilliseconds float64 `json:"totalMilliseconds"`
	TotalMinutes      float64 `json:"totalMinutes"`
	TotalSeconds      float64 `json:"totalSeconds"`
}

type TrackedDownloadStatus string

const (
	OK      TrackedDownloadStatus = "ok"
	Warning TrackedDownloadStatus = "warning"
	Error   TrackedDownloadStatus = "error"
)

type TrackedDownloadState string

const (
	Downloading   TrackedDownloadState = "downloading"
	ImportPending TrackedDownloadState = "importPending"
	Importing     TrackedDownloadState = "importing"
	Imported      TrackedDownloadState = "imported"
	FailedPending TrackedDownloadState = "failedPending"
	Failed        TrackedDownloadState = "failed"
	Ignored       TrackedDownloadState = "ignored"
)

type TrackedDownloadStatusMessage struct {
	Title    string   `json:"title"`
	Messages []string `json:"messages"`
}

type DownloadProtocol string

const (
	UnknownDownloadProtocol DownloadProtocol = "unknown"
	Usenet                  DownloadProtocol = "usenet"
	Torrent                 DownloadProtocol = "torrent"
)

type QualityProfileResource struct {
	ID                int32  `json:"id"`
	Name              string `json:"name"`
	UpgradeAllowed    bool   `json:"upgradeAllowed"`
	Cutoff            int32
	Items             []QualityProfileQualityItemResource `json:"items"`
	MinFormatScore    int32                               `json:"minFormatScore"`
	CutoffFormatScore int32                               `json:"cutoffFormatScore"`
	FormatItems       []ProfileFormatItemResource         `json:"formatItems"`
}

type QualityProfileQualityItemResource struct {
	ID      int32   `json:"id"`
	Name    string  `json:"name"`
	Quality Quality `json:"quality"`
	Items   []any   `json:"items"`
	Allowed bool    `json:"allowed"`
}

type ProfileFormatItemResource struct {
	ID     int32  `json:"id"`
	Format int32  `json:"format"`
	Name   string `json:"name"`
	Score  int32  `json:"score"`
}
