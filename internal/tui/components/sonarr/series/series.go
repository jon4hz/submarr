package series

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/stickers/flexbox"
	"github.com/jon4hz/submarr/internal/core/sonarr"
	"github.com/jon4hz/submarr/internal/tui/common"
	sonarr_list "github.com/jon4hz/submarr/internal/tui/components/sonarr/list"
	"github.com/jon4hz/submarr/internal/tui/components/sonarr/removeseries"
	"github.com/jon4hz/submarr/internal/tui/components/sonarr/seasons"
	"github.com/jon4hz/submarr/internal/tui/components/statusbar"
	"github.com/jon4hz/submarr/internal/tui/overlay"
	sonarrAPI "github.com/jon4hz/submarr/pkg/sonarr"
	zone "github.com/lrstanley/bubblezone"
)

type SelectSeasonMsg int

var (
	subtleColor   = lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}
	selectedColor = lipgloss.Color("#00CCFF")
)

type cell int

const (
	infoCell cell = iota
	seasonsCell
	statsCell
)

type state int

const (
	stateSeries state = iota + 1
	stateDelete
)

type Model struct {
	common.EmbedableModel

	flexBox      *flexbox.HorizontalFlexBox
	cellmap      map[cell]*flexbox.Cell
	selectedCell cell

	client *sonarr.Client

	infoViewport  viewport.Model
	statsViewport viewport.Model
	seasonsList   list.Model
	state         state
	delete        common.SubModel
}

var (
	cellStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true).
			BorderForeground(subtleColor)

	infoCellStyle = cellStyle.Copy().
			MarginLeft(2).
			Padding(0, 2, 1, 2)

	statsCellStyle = cellStyle.Copy().
			MarginLeft(2)

	seasonsCellStyle = cellStyle.Copy()
)

func New(sonarr *sonarr.Client, width, height int) *Model {
	m := Model{
		flexBox:       flexbox.NewHorizontal(width, height),
		client:        sonarr,
		infoViewport:  viewport.New(0, 0),
		statsViewport: viewport.New(0, 0),
		selectedCell:  seasonsCell,
		state:         stateSeries,
		seasonsList:   sonarr_list.New("Seasons", newSeasonsItems(sonarr.GetSerie()), seasons.Delegate{}, 0, 0),
	}

	m.Width = width
	m.Height = height

	m.cellmap = map[cell]*flexbox.Cell{
		infoCell:    flexbox.NewCell(1, 2).SetStyle(infoCellStyle).SetID("info"),
		statsCell:   flexbox.NewCell(1, 5).SetStyle(statsCellStyle).SetID("stats"),
		seasonsCell: flexbox.NewCell(1, 1).SetStyle(seasonsCellStyle).SetID("Seasons"),
	}

	columns := []*flexbox.Column{
		m.flexBox.NewColumn().AddCells(
			m.cellmap[infoCell],
			m.cellmap[statsCell],
		),
		m.flexBox.NewColumn().AddCells(
			m.cellmap[seasonsCell],
		),
	}
	m.flexBox.AddColumns(columns)

	// initial render
	m.SetSize(width, height)
	m.updateFocus()
	m.redraw()
	m.updateStatsViewport()

	return &m
}

func newSeasonsItems(serie *sonarrAPI.SeriesResource) []list.Item {
	sort.Slice(serie.Seasons, func(i, j int) bool {
		return serie.Seasons[i].SeasonNumber > serie.Seasons[j].SeasonNumber
	})
	items := make([]list.Item, len(serie.Seasons))
	for i, season := range serie.Seasons {
		items[i] = seasons.NewItem(i, season)
	}
	return items
}

func (m Model) Init() tea.Cmd {
	return statusbar.NewHelpCmd(DefaultKeyMap.FullHelp())
}

func (m *Model) Update(msg tea.Msg) (common.SubModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case stateSeries:
			switch {
			case key.Matches(msg, DefaultKeyMap.Back):
				if !m.seasonsList.SettingFilter() {
					m.IsBack = true
					return m, nil
				}

			case key.Matches(msg, DefaultKeyMap.Quit):
				if !m.seasonsList.SettingFilter() {
					m.IsQuit = true
					return m, nil
				}

			case key.Matches(msg, DefaultKeyMap.Tab):
				if !m.seasonsList.SettingFilter() {
					m.focusNext()
					return m, nil
				}

			case key.Matches(msg, DefaultKeyMap.ShiftTab):
				if !m.seasonsList.SettingFilter() {
					m.focusPrev()
					return m, nil
				}

			case key.Matches(msg, DefaultKeyMap.Reload):
				if !m.seasonsList.SettingFilter() {
					return m, tea.Batch(
						m.client.ReloadSerie(),
						m.seasonsList.StartSpinner(),
						statusbar.NewMessageCmd("Reloading series...", statusbar.WithMessageTimeout(2)),
					)
				}

			case key.Matches(msg, DefaultKeyMap.ToggleMonitor):
				if !m.seasonsList.SettingFilter() {
					return m, tea.Batch(
						m.client.ToggleMonitorSeason(m.seasonsList.Index()),
						m.seasonsList.StartSpinner(),
						statusbar.NewMessageCmd("Toggling season monitor...", statusbar.WithMessageTimeout(2)),
					)
				}

			case key.Matches(msg, DefaultKeyMap.ToggleMonitorSeries):
				if !m.seasonsList.SettingFilter() {
					return m, tea.Batch(
						m.client.ToggleMonitorSeries(),
						m.seasonsList.StartSpinner(),
						statusbar.NewMessageCmd("Toggling series monitor...", statusbar.WithMessageTimeout(2)),
					)
				}

			case key.Matches(msg, DefaultKeyMap.Select):
				if !m.seasonsList.SettingFilter() {
					item := m.seasonsList.SelectedItem().(seasons.SeasonItem)
					return m, func() tea.Msg { return SelectSeasonMsg(item.Index) }
				}

			case key.Matches(msg, DefaultKeyMap.Refresh):
				return m, tea.Batch(
					m.client.RefreshSeries(),
					statusbar.NewMessageCmd("Refreshing series...", statusbar.WithMessageTimeout(2)),
				)

			case key.Matches(msg, DefaultKeyMap.AutomaticSearch):
				if !m.seasonsList.SettingFilter() {
					season := m.seasonsList.SelectedItem().(seasons.SeasonItem)
					return m, tea.Batch(
						m.client.AutomaticSearchSeason(season.Season.SeasonNumber),
						statusbar.NewMessageCmd(fmt.Sprintf("Searching for season %d...", season.Season.SeasonNumber), statusbar.WithMessageTimeout(2)),
					)
				}

			case key.Matches(msg, DefaultKeyMap.AutomaticSearchAll):
				return m, tea.Batch(
					m.client.AutomaticSearchSeries(),
					statusbar.NewMessageCmd("Searching for all seasons...", statusbar.WithMessageTimeout(2)),
				)

			case key.Matches(msg, DefaultKeyMap.Delete):
				return m, m.deleteSeries()
			}
		}

	case sonarr.FetchSerieResult:
		m.seasonsList.StopSpinner()
		if msg.Error != nil {
			return m, statusbar.NewMessageCmd(msg.Error.Error(), statusbar.WithMessageTimeout(3))
		}
		m.updateStatsViewport()
		return m, m.seasonsList.SetItems(newSeasonsItems(msg.Serie))

	case tea.MouseMsg:
		switch m.state {
		case stateSeries:
			switch m.selectedCell {
			case seasonsCell:
				switch msg.Type {
				case tea.MouseWheelUp:
					m.seasonsList.CursorUp()
					return m, nil

				case tea.MouseWheelDown:
					m.seasonsList.CursorDown()
					return m, nil

				case tea.MouseLeft:
					for i, listItem := range m.seasonsList.VisibleItems() {
						item, _ := listItem.(seasons.SeasonItem)
						if zone.Get(fmt.Sprintf("Season %d", item.Season.SeasonNumber)).InBounds(msg) {
							// if we click on an already selected item, open the details
							if i == m.seasonsList.Index() {
								return m, func() tea.Msg { return SelectSeasonMsg(item.Index) }
							}
							// else select the item
							m.seasonsList.Select(i)
							break
						}
					}
				}
			}
		}
	}

	switch m.state {
	case stateSeries:

		switch m.selectedCell {
		case infoCell:
			var cmd tea.Cmd
			m.infoViewport, cmd = m.infoViewport.Update(msg)
			return m, cmd

		case statsCell:
			var cmd tea.Cmd
			m.statsViewport, cmd = m.statsViewport.Update(msg)
			return m, cmd

		case seasonsCell:
			var cmd tea.Cmd
			m.seasonsList, cmd = m.seasonsList.Update(msg)
			return m, cmd
		}
	case stateDelete:
		var cmd tea.Cmd
		m.delete, cmd = m.delete.Update(msg)

		if m.delete.Quit() {
			m.IsQuit = true
			return m, nil
		}

		if m.delete.Back() {
			m.state = stateSeries
			return m, statusbar.NewHelpCmd(DefaultKeyMap.FullHelp())
		}

		return m, cmd
	}
	return m, nil
}

func (m *Model) deleteSeries() tea.Cmd {
	m.state = stateDelete
	m.delete = removeseries.New(m.client, m.client.GetSerie(), m.Width, m.Height)
	return m.delete.Init()
}

func (m *Model) SetSize(width, height int) {
	m.Width = width
	m.Height = height
	m.flexBox.SetWidth(width)
	m.flexBox.SetHeight(height)

	// if the layout changes, we need to force a recalculation.
	// otherwise we cant use the size of the cells when rendering the content.
	m.flexBox.ForceRecalculate()

	m.infoViewport.Height = max(m.cellmap[infoCell].GetContentHeight()-infoCellStyle.GetVerticalPadding(), 0) // make sure we dont get a negative value
	m.infoViewport.Width = max(m.cellmap[infoCell].GetContentWidth()-infoCellStyle.GetHorizontalPadding(), 0) // make sure we dont get a negative value

	m.statsViewport.Height = max(m.cellmap[statsCell].GetContentHeight()-statsCellStyle.GetVerticalPadding(), 0)
	m.statsViewport.Width = max(m.cellmap[statsCell].GetContentWidth()-statsCellStyle.GetHorizontalPadding(), 0)

	seasonListHeight := max(m.cellmap[seasonsCell].GetContentHeight()-seasonsCellStyle.GetVerticalPadding(), 0)
	seasonListWidth := max(m.cellmap[seasonsCell].GetContentWidth()-seasonsCellStyle.GetHorizontalPadding(), 0)
	m.seasonsList.SetSize(seasonListWidth, seasonListHeight)

	if m.state == stateDelete {
		m.delete.SetSize(width, height)
	}
}

func (m *Model) focusNext() {
	m.selectedCell++
	if m.selectedCell > cell(len(m.cellmap)-1) {
		m.selectedCell = 0
	}
	m.updateFocus()
}

func (m *Model) focusPrev() {
	m.selectedCell--
	if m.selectedCell < 0 {
		m.selectedCell = cell(len(m.cellmap) - 1)
	}
	m.updateFocus()
}

func (m *Model) updateFocus() {
	var (
		infoS   = infoCellStyle.Copy()
		statsS  = statsCellStyle.Copy()
		seasonS = seasonsCellStyle.Copy()
	)

	switch m.selectedCell {
	case infoCell:
		infoS = infoS.Copy().BorderForeground(selectedColor)
	case statsCell:
		statsS = statsS.Copy().BorderForeground(selectedColor)
	case seasonsCell:
		seasonS = seasonS.Copy().BorderForeground(selectedColor)
	}

	m.cellmap[infoCell].SetStyle(infoS)
	m.cellmap[statsCell].SetStyle(statsS)
	m.cellmap[seasonsCell].SetStyle(seasonS)
}

func (m *Model) updateStatsViewport() {
	serie := m.client.GetSerie()
	if serie == nil {
		return
	}
	var s strings.Builder

	s.WriteString("Monitoring: ")
	if serie.Monitored {
		s.WriteString("Yes")
	} else {
		s.WriteString("No")
	}
	s.WriteByte('\n')

	s.WriteString("Type: ")
	s.WriteString(common.Title(string(serie.SeriesType)))
	s.WriteByte('\n')

	s.WriteString("Path: ")
	s.WriteString(fmt.Sprintf("%q", serie.Path))
	s.WriteByte('\n')

	s.WriteString("Quality: ")

	qp := m.client.GetSerieQualityProfile()
	if qp != nil {
		s.WriteString(qp.Name)
	} else {
		s.WriteString("Unknown")
	}
	s.WriteByte('\n')

	s.WriteString("Language: ")
	if serie.OriginalLanguage != nil {
		s.WriteString(serie.OriginalLanguage.Name)
	} else {
		s.WriteString("Unknown")
	}
	s.WriteByte('\n')

	s.WriteByte('\n')

	s.WriteString("Status: ")
	s.WriteString(common.Title(string(serie.Status)))
	s.WriteByte('\n')

	s.WriteString("Next Airing: ")
	if serie.Status == "continuing" {
		if serie.NextAiring.IsZero() {
			s.WriteString("Unknown")
		} else {
			s.WriteString(serie.NextAiring.Local().Format("January 2, 2006 - 15:04"))
		}
	} else {
		s.WriteString("Series Ended")
	}
	s.WriteByte('\n')

	s.WriteString("Added on: ")
	s.WriteString(serie.Added.Local().Format("January 2, 2006"))
	s.WriteByte('\n')

	s.WriteByte('\n')

	s.WriteString("Year: ")
	s.WriteString(strconv.FormatInt(int64(serie.Year), 10))
	s.WriteByte('\n')

	s.WriteString("Network: ")
	s.WriteString(serie.Network)
	s.WriteByte('\n')

	s.WriteString("Runtime: ")
	s.WriteString(strconv.FormatInt(int64(serie.Runtime), 10) + "m")
	s.WriteByte('\n')

	s.WriteString("Rating: ")
	s.WriteString(serie.Certification)
	s.WriteByte('\n')

	s.WriteString("Genres: ")
	for i, genre := range serie.Genres {
		s.WriteString(genre)
		if i < len(serie.Genres)-1 {
			s.WriteString(", ")
		}
	}
	s.WriteByte('\n')

	s.WriteString("Alternate Titles: ")
	for i, title := range serie.AlternateTitles {
		s.WriteString(title.Title)
		if i < len(serie.AlternateTitles)-1 {
			s.WriteString(", ")
		}
	}

	m.statsViewport.SetContent(s.String())
}

func (m *Model) redraw() {
	m.cellmap[infoCell].SetContent(m.renderInfoCell())
	m.cellmap[statsCell].SetContent(m.renderStatsCell())
	m.cellmap[seasonsCell].SetContent(m.renderSeasonsCell())
}

var (
	seriesTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Border(lipgloss.NormalBorder(), false, false, true, false). // top, right, bottom, left
				BorderForeground(subtleColor).
		//Margin(0, 2).
		Align(lipgloss.Center)

	descStyle = lipgloss.NewStyle()
)

func (m *Model) renderInfoCell() string {
	var (
		s            strings.Builder
		contentWidth = m.infoViewport.Width - 1
	)

	s.WriteString(seriesTitleStyle.Width(contentWidth).Render(m.client.GetSerie().Title))
	s.WriteByte('\n')
	s.WriteString(descStyle.Width(contentWidth).Render(m.client.GetSerie().Overview))

	m.infoViewport.SetContent(s.String())

	return m.infoViewport.View()
}

func (m Model) renderSeasonsCell() string {
	return m.seasonsList.View()
}

func (m Model) renderStatsCell() string {
	return m.statsViewport.View()
}

func (m Model) View() string {
	switch m.state {
	case stateSeries:
		m.redraw()
		return m.flexBox.Render()

	case stateDelete:
		m.redraw()
		fg := m.delete.View()
		x := ((m.Width - lipgloss.Width(fg)) / 2)
		y := ((m.Height - lipgloss.Height(fg)) / 2)
		// make sure background fills the whole screen
		bg := lipgloss.NewStyle().Width(m.Width).Height(m.Height).Render(m.flexBox.Render())
		return overlay.PlaceOverlay(x, y, fg, bg)
	}

	return ":("
}
