package serie

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/stickers/flexbox"
	"github.com/jon4hz/subrr/internal/core/sonarr"
	"github.com/jon4hz/subrr/internal/tui/common"
	"github.com/jon4hz/subrr/internal/tui/sonarr/seasons"
	"github.com/jon4hz/subrr/internal/tui/statusbar"
	sonarrAPI "github.com/jon4hz/subrr/pkg/sonarr"
	zone "github.com/lrstanley/bubblezone"
)

var (
	subtleColor   = lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}
	selectedColor = lipgloss.Color("#00CCFF")
)

type cell int

const (
	infoCell cell = iota
	//statsCell
	seasonsCell
	historyCell
)

var cellMap = map[cell]*flexbox.Cell{}

type Model struct {
	common.EmbedableModel

	flexBox      *flexbox.HorizontalFlexBox
	cellmap      map[cell]*flexbox.Cell
	selectedCell cell

	client *sonarr.Client
	serie  *sonarrAPI.SeriesResource

	infoViewport viewport.Model
	seasonsList  list.Model
}

var (
	cellStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true).
			BorderForeground(subtleColor)

	infoCellStyle = cellStyle.Copy().
			MarginLeft(2).
			Padding(0, 2, 1, 2)

	historyCellStyle = cellStyle.Copy().
				MarginLeft(2)

	statsCellStyle = cellStyle.Copy()

	seasonsCellStyle = cellStyle.Copy()
)

func New(sonarr *sonarr.Client, serie *sonarrAPI.SeriesResource, width, heigth int) *Model {
	m := Model{
		flexBox:      flexbox.NewHorizontal(width, heigth),
		client:       sonarr,
		serie:        serie,
		infoViewport: viewport.New(0, 0),
		selectedCell: infoCell,
		seasonsList:  list.New(newSeasonsItems(serie), seasons.Delegate{}, 0, 0),
	}

	m.Width = width
	m.Height = heigth

	m.cellmap = map[cell]*flexbox.Cell{
		infoCell: flexbox.NewCell(1, 1).SetStyle(infoCellStyle).SetID("info"),
		//statsCell:   flexbox.NewCell(1, 1).SetStyle(statsCellStyle).SetID("stats"),
		historyCell: flexbox.NewCell(1, 1).SetStyle(historyCellStyle).SetID("history").SetContent("History"),
		seasonsCell: flexbox.NewCell(1, 1).SetStyle(seasonsCellStyle).SetID("Seasons"),
	}

	columns := []*flexbox.Column{
		m.flexBox.NewColumn().AddCells(
			m.cellmap[infoCell],
			m.cellmap[historyCell],
		),
		m.flexBox.NewColumn().AddCells(
			//m.cellmap[statsCell],
			m.cellmap[seasonsCell],
		),
	}
	m.flexBox.AddColumns(columns)

	m.seasonsList.SetShowHelp(false)

	// initial render
	m.SetSize(width, heigth)
	m.updateFocus()
	m.redraw()

	return &m
}

func newSeasonsItems(serie *sonarrAPI.SeriesResource) []list.Item {
	sort.Slice(serie.Seasons, func(i, j int) bool {
		return serie.Seasons[i].SeasonNumber > serie.Seasons[j].SeasonNumber
	})
	items := make([]list.Item, len(serie.Seasons))
	for i, season := range serie.Seasons {
		items[i] = seasons.NewItem(season)
	}
	return items
}

func (m Model) Init() tea.Cmd {
	return statusbar.NewHelpCmd(DefaultKeyMap.FullHelp())
}

func (m *Model) Update(msg tea.Msg) (common.SubModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Back):
			m.IsBack = true
			return m, nil

		case key.Matches(msg, DefaultKeyMap.Quit):
			m.IsQuit = true
			return m, nil

		case key.Matches(msg, DefaultKeyMap.Tab):
			m.focusNext()
			return m, nil

		case key.Matches(msg, DefaultKeyMap.ShiftTab):
			m.focusPrev()
			return m, nil
		}

	case tea.MouseMsg:
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
							/* cmd := m.selectSeries(&item.Series)
							return m, cmd */
						}
						// else select the item
						m.seasonsList.Select(i)
						break
					}
				}
			}
		}
	}

	switch m.selectedCell {
	case infoCell:
		var cmd tea.Cmd
		m.infoViewport, cmd = m.infoViewport.Update(msg)
		return m, cmd
	/* case statsCell:
	return m, nil */
	case historyCell:
		return m, nil
	case seasonsCell:
		var cmd tea.Cmd
		m.seasonsList, cmd = m.seasonsList.Update(msg)
		return m, cmd
	}
	return m, nil
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

	seasonListHeight := max(m.cellmap[seasonsCell].GetContentHeight()-seasonsCellStyle.GetVerticalPadding(), 0)
	seasonListWidth := max(m.cellmap[seasonsCell].GetContentWidth()-seasonsCellStyle.GetHorizontalPadding(), 0)
	m.seasonsList.SetSize(seasonListWidth, seasonListHeight)
}

// nolint:unparam
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
		infoS = infoCellStyle.Copy()
		//statsS   = statsCellStyle.Copy()
		historyS = historyCellStyle.Copy()
		seasonS  = seasonsCellStyle.Copy()
	)

	switch m.selectedCell {
	case infoCell:
		infoS = infoS.Copy().BorderForeground(selectedColor)
	/* case statsCell:
	statsS = statsS.Copy().BorderForeground(selectedColor) */
	case historyCell:
		historyS = historyS.Copy().BorderForeground(selectedColor)
	case seasonsCell:
		seasonS = seasonS.Copy().BorderForeground(selectedColor)
	}

	m.cellmap[infoCell].SetStyle(infoS)
	//m.cellmap[statsCell].SetStyle(statsS)
	m.cellmap[historyCell].SetStyle(historyS)
	m.cellmap[seasonsCell].SetStyle(seasonS)
}

func (m *Model) redraw() {
	m.cellmap[infoCell].SetContent(m.renderInfoCell())
	//m.cellmap[statsCell].     // .SetContent(m.renderStatsCell())
	//m.cellmap[historyCell].// .SetContent(m.renderHistoryCell())
	m.cellmap[seasonsCell].SetContent(m.renderSeasonsCell())
}

var (
	titleStyle = lipgloss.NewStyle().
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

	s.WriteString(titleStyle.Width(contentWidth).Render(m.serie.Title))
	s.WriteByte('\n')
	s.WriteString(descStyle.Width(contentWidth).Render(m.serie.Overview))

	m.infoViewport.SetContent(s.String())

	return m.infoViewport.View()
}

func (m *Model) renderSeasonsCell() string {
	return m.seasonsList.View()
}

func (m Model) View() string {
	m.redraw()
	return m.flexBox.Render()
}
