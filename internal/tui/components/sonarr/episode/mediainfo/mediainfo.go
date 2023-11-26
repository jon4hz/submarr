package mediainfo

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/submarr/internal/tui/common"
	"github.com/jon4hz/submarr/internal/tui/components/statusbar"
	"github.com/jon4hz/submarr/internal/tui/styles"
	sonarrAPI "github.com/jon4hz/submarr/pkg/sonarr"
)

type kv struct {
	key   string
	value string
}

type Model struct {
	common.EmbedableModel

	mediaInfo  *sonarrAPI.MediaInfoResource
	kvs        [13]kv
	longestKey int
}

func New(episode *sonarrAPI.EpisodeResource, width, height int) common.SubModel {
	m := Model{
		mediaInfo: episode.EpisodeFile.MediaInfo,
	}

	m.setKVs()
	m.SetSize(width, height)

	return &m
}

func (m *Model) setKVs() {
	m.kvs = [...]kv{
		{"Audio Bitrate", fmt.Sprintf("%d", m.mediaInfo.AudioBitrate)},
		{"Audio Channels", strconv.FormatFloat(m.mediaInfo.AudioChannels, 'f', -1, 64)},
		{"Audio Codec", m.mediaInfo.AudioCodec},
		{"Audio Languages", m.mediaInfo.AudioLanguages},
		{"Audio Stream Count", fmt.Sprintf("%d", m.mediaInfo.AudioStreamCount)},
		{"Video Bit Depth", fmt.Sprintf("%d", m.mediaInfo.VideoBitDepth)},
		{"Video Bitrate", fmt.Sprintf("%d", m.mediaInfo.VideoBitrate)},
		{"Video Codec", m.mediaInfo.VideoCodec},
		{"Video Fps", strconv.FormatFloat(m.mediaInfo.VideoFps, 'f', -1, 64)},
		{"Resolution", m.mediaInfo.Resolution},
		{"Run Time", m.mediaInfo.RunTime},
		{"Scan Type", m.mediaInfo.ScanType},
		{"Subtitles", wrap(m.mediaInfo.Subtitles, 27)},
	}

	for _, kv := range &m.kvs {
		if len(kv.key) > m.longestKey {
			m.longestKey = len(kv.key)
		}
	}
}

func wrap(s string, width int) string {
	if lipgloss.Width(s) <= width {
		return s
	}
	return lipgloss.NewStyle().Width(width).Render(s)
}

func (m Model) Init() tea.Cmd {
	return statusbar.NewHelpCmd(DefaultKeyMap.FullHelp())
}

func (m *Model) Update(_ tea.Msg) (common.SubModel, tea.Cmd) {
	return m, nil
}

var (
	titleStyle = lipgloss.NewStyle().
			Padding(0, 1, 0, 1).
			Margin(0, 0, 1, 0).
			Bold(true).
			Background(styles.PurpleColor)
	kvKeyStyle = lipgloss.NewStyle().
			Align(lipgloss.Right).
			Padding(0, 2, 0, 0)
	kvValStyle = lipgloss.NewStyle()
)

func (m Model) View() string {
	var s strings.Builder
	keyStyle := kvKeyStyle.Copy().Width(m.longestKey + 2)

	keys := make([]string, 0, len(m.kvs))
	values := make([]string, 0, len(m.kvs))

	for _, kv := range &m.kvs {
		keys = append(keys, keyStyle.Render(kv.key))
		values = append(values, kvValStyle.Render(kv.value))
	}

	kvs := lipgloss.JoinHorizontal(lipgloss.Left,
		lipgloss.JoinVertical(lipgloss.Top, keys...),
		lipgloss.JoinVertical(lipgloss.Top, values...),
	)

	s.WriteString(
		lipgloss.Place(
			lipgloss.Width(kvs), 1,
			lipgloss.Center, lipgloss.Center,
			titleStyle.Render("Media Info"),
		),
	)
	s.WriteByte('\n')
	s.WriteString(kvs)

	return s.String()
}

func (m *Model) SetSize(width, height int) {
	m.Width = width
	m.Height = height
}
