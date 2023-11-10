package addseries

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/submarr/internal/core/sonarr"
	"github.com/jon4hz/submarr/internal/tui/common"
	sonarr_list "github.com/jon4hz/submarr/internal/tui/components/sonarr/list"
	"github.com/jon4hz/submarr/internal/tui/components/statusbar"
	"github.com/jon4hz/submarr/internal/tui/components/toggle"
	sonarrAPI "github.com/jon4hz/submarr/pkg/sonarr"
)

type Model struct {
	common.EmbedableModel

	client *sonarr.Client
	series *sonarrAPI.SeriesResource

	rootFolder                   list.Model
	monitor                      list.Model
	qualityProfile               list.Model
	languageProfile              list.Model
	seriesType                   list.Model
	seasonFolder                 toggle.Model
	searchForMissingEpisodes     toggle.Model
	searchForCutoffUnmetEpisodes toggle.Model

	selectedOption     addOption
	showOptions        bool
	longestOptionWidth int
}

type addOption int

const (
	addOptionRootFolder addOption = iota + 1
	addOptionMonitor
	addOptionQualityProfile
	addOptionLanguageProfile
	addOptionSeriesType
	addOptionSeasonFolder
	addOptionSearchForMissingEpisodes
	addOptionSearchForCutoffUnmetEpisodes
	addOptionAddSeries
)

var addOptions = map[addOption]string{
	addOptionRootFolder:                   "Root Folder",
	addOptionMonitor:                      "Monitor",
	addOptionQualityProfile:               "Quality Profile",
	addOptionLanguageProfile:              "Language Profile",
	addOptionSeriesType:                   "Series Type",
	addOptionSeasonFolder:                 "Season Folder",
	addOptionSearchForMissingEpisodes:     "Search missing",
	addOptionSearchForCutoffUnmetEpisodes: "Search cutoff unmet",
	addOptionAddSeries:                    "",
}

func New(client *sonarr.Client, series *sonarrAPI.SeriesResource, width, height int) common.SubModel {
	setDefaults(client, series)

	m := Model{
		client:             client,
		series:             series,
		selectedOption:     1,
		longestOptionWidth: getLongestOptionWidth(),
		rootFolder: sonarr_list.New(
			"Select Root Folder",
			newRootFolderItems(client.GetRootFolders()),
			rootFolderDelegate{},
			width, height,
		),
		monitor: sonarr_list.New(
			"Select Monitor Type",
			newMonitorItems(),
			monitorDelegate{},
			width, height,
		),
		qualityProfile: sonarr_list.New(
			"Select Quality Profile",
			newQualityProfileItems(client.GetQualityProfiles()),
			qualityProfileDelegate{},
			width, height,
		),
		languageProfile: sonarr_list.New(
			"Select Language Profile",
			newLanguageProfileItems(client.GetLanguageProfiles()),
			languageProfileDelegate{},
			width, height,
		),
		seriesType: sonarr_list.New(
			"Select Series Type",
			newSeriesTypeItems(),
			seriesTypeDelegate{},
			width, height,
		),
		seasonFolder:                 toggle.New(true),
		searchForMissingEpisodes:     toggle.New(true),
		searchForCutoffUnmetEpisodes: toggle.New(true),
	}

	// make sure id is 0
	m.series.ID = 0

	m.Width = width
	m.Height = height

	sublists := []list.Model{
		m.rootFolder,
		m.monitor,
		m.qualityProfile,
		m.languageProfile,
		m.seriesType,
	}
	for _, sublist := range sublists {
		sublist.SetShowFilter(false)
		sublist.SetShowStatusBar(false)
	}

	return &m
}

func getLongestOptionWidth() int {
	var longest int
	for _, option := range addOptions {
		if len(option) > longest {
			longest = len(option)
		}
	}
	return longest
}

func setDefaults(client *sonarr.Client, series *sonarrAPI.SeriesResource) {
	rootFolders := client.GetRootFolders()
	if len(rootFolders) > 0 {
		series.RootFolderPath = rootFolders[0].Path
	}

	qualityProfiles := client.GetQualityProfiles()
	if client.Config.DefaultQualityProfile != "" {
		for _, qualityProfile := range qualityProfiles {
			if qualityProfile.Name == client.Config.DefaultQualityProfile {
				series.QualityProfileID = qualityProfile.ID
				break
			}
		}
	} else {
		if len(qualityProfiles) > 0 {
			series.QualityProfileID = qualityProfiles[0].ID
		}
	}

	languageProfiles := client.GetLanguageProfiles()
	if client.Config.DefaultLanguageProfile != "" {
		for _, languageProfile := range languageProfiles {
			if languageProfile.Name == client.Config.DefaultLanguageProfile {
				series.LanguageProfileID = languageProfile.ID
				break
			}
		}
	} else {
		if len(languageProfiles) > 0 {
			series.LanguageProfileID = languageProfiles[0].ID
		}
	}

	series.AddOptions = &sonarrAPI.AddSeriesOptions{
		Monitor:                      sonarrAPI.All,
		SearchForMissingEpisodes:     true,
		SearchForCutoffUnmetEpisodes: true,
	}
	series.Monitored = true
	series.SeriesType = sonarrAPI.Standard

	series.SeasonFolder = true
}

func (m Model) Init() tea.Cmd {
	return statusbar.NewHelpCmd(DefaultKeyMap.FullHelp())
}

func (m *Model) Update(msg tea.Msg) (common.SubModel, tea.Cmd) {
	if m.showOptions {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.showOptions = false
				return m, nil
			}
		}

		switch m.selectedOption {
		case addOptionRootFolder:
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.String() {
				case " ":
					for i, item := range m.rootFolder.Items() {
						rfItem := item.(rootFolderItem)
						if i == m.rootFolder.Index() {
							rfItem.triggered = true
							m.series.RootFolderPath = rfItem.rootFolder.Path
						} else {
							rfItem.triggered = false
						}
						m.rootFolder.SetItem(i, rfItem)
					}
				}
			}

			var cmd tea.Cmd
			m.rootFolder, cmd = m.rootFolder.Update(msg)
			return m, cmd

		case addOptionMonitor:
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.String() {
				case " ":
					for i, item := range m.monitor.Items() {
						monitorItem := item.(monitorItem)
						if i == m.monitor.Index() {
							monitorItem.triggered = true
							m.series.AddOptions.Monitor = monitorItem.monitor
							if monitorItem.monitor == sonarrAPI.None {
								m.series.Monitored = false
							}
						} else {
							monitorItem.triggered = false
						}
						m.monitor.SetItem(i, monitorItem)
					}
				}
			}

			var cmd tea.Cmd
			m.monitor, cmd = m.monitor.Update(msg)
			return m, cmd

		case addOptionQualityProfile:
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.String() {
				case " ":
					for i, item := range m.qualityProfile.Items() {
						qpItem := item.(qualityProfileItem)
						if i == m.qualityProfile.Index() {
							qpItem.triggered = true
							m.series.QualityProfileID = qpItem.qualityProfile.ID
						} else {
							qpItem.triggered = false
						}
						m.qualityProfile.SetItem(i, qpItem)
					}
				}
			}

			var cmd tea.Cmd
			m.qualityProfile, cmd = m.qualityProfile.Update(msg)
			return m, cmd

		case addOptionLanguageProfile:
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.String() {
				case " ":
					for i, item := range m.languageProfile.Items() {
						lpItem := item.(languageProfileItem)
						if i == m.languageProfile.Index() {
							lpItem.triggered = true
							m.series.LanguageProfileID = lpItem.languageProfile.ID
						} else {
							lpItem.triggered = false
						}
						m.languageProfile.SetItem(i, lpItem)
					}
				}
			}

			var cmd tea.Cmd
			m.languageProfile, cmd = m.languageProfile.Update(msg)
			return m, cmd

		case addOptionSeriesType:
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.String() {
				case " ":
					for i, item := range m.seriesType.Items() {
						stItem := item.(seriesTypeItem)
						if i == m.seriesType.Index() {
							stItem.triggered = true
							m.series.SeriesType = stItem.seriesType
						} else {
							stItem.triggered = false
						}
						m.seriesType.SetItem(i, stItem)
					}
				}
			}

			var cmd tea.Cmd
			m.seriesType, cmd = m.seriesType.Update(msg)
			return m, cmd

		}

		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Down):
			m.nextOption()
		case key.Matches(msg, DefaultKeyMap.Up):
			m.previousOption()
		case key.Matches(msg, DefaultKeyMap.Select):
			switch m.selectedOption {
			case
				addOptionSeasonFolder,
				addOptionSearchForMissingEpisodes,
				addOptionSearchForCutoffUnmetEpisodes:
				// no-op
			case addOptionAddSeries:
				return m, m.addSeries()

			default:
				m.showOptions = true
			}
		case key.Matches(msg, DefaultKeyMap.Add):
			return m, m.addSeries()
		case key.Matches(msg, DefaultKeyMap.Back):
			m.IsBack = true
		case key.Matches(msg, DefaultKeyMap.Quit):
			m.IsQuit = true
		}
	}

	switch m.selectedOption {
	case addOptionSeasonFolder:
		var cmd tea.Cmd
		m.seasonFolder, cmd = m.seasonFolder.Update(msg)
		if m.seasonFolder.Toggled() {
			m.series.SeasonFolder = true
		} else {
			m.series.SeasonFolder = false
		}
		return m, cmd

	case addOptionSearchForMissingEpisodes:
		var cmd tea.Cmd
		m.searchForMissingEpisodes, cmd = m.searchForMissingEpisodes.Update(msg)
		if m.searchForMissingEpisodes.Toggled() {
			m.series.AddOptions.SearchForMissingEpisodes = true
		} else {
			m.series.AddOptions.SearchForMissingEpisodes = false
		}
		return m, cmd

	case addOptionSearchForCutoffUnmetEpisodes:
		var cmd tea.Cmd
		m.searchForCutoffUnmetEpisodes, cmd = m.searchForCutoffUnmetEpisodes.Update(msg)
		if m.searchForCutoffUnmetEpisodes.Toggled() {
			m.series.AddOptions.SearchForCutoffUnmetEpisodes = true
		} else {
			m.series.AddOptions.SearchForCutoffUnmetEpisodes = false
		}
		return m, cmd
	}

	return m, nil
}

func (m *Model) nextOption() {
	m.selectedOption++
	if int(m.selectedOption) > len(addOptions) {
		m.selectedOption = 1
	}
}

func (m *Model) previousOption() {
	m.selectedOption--
	if m.selectedOption < 1 {
		m.selectedOption = addOption(len(addOptions))
	}
}

func (m Model) addSeries() tea.Cmd {
	return m.client.PostSeries(m.series)
}

func (m Model) View() string {
	if !m.showOptions {
		return m.optionsView()
	}

	switch m.selectedOption {
	case addOptionRootFolder:
		return boxStyle.Width(m.Width - 2).Render(m.rootFolder.View())
	case addOptionMonitor:
		return boxStyle.Width(m.Width - 2).Render(m.monitor.View())
	case addOptionQualityProfile:
		return boxStyle.Width(m.Width - 2).Render(m.qualityProfile.View())
	case addOptionLanguageProfile:
		return boxStyle.Width(m.Width - 2).Render(m.languageProfile.View())
	case addOptionSeriesType:
		return boxStyle.Width(m.Width - 2).Render(m.seriesType.View())
	default:
		return m.optionsView()
	}
}

var (
	boxStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).Padding(1, 2, 1, 2)
	titleStyle  = lipgloss.NewStyle().Align(lipgloss.Center).Bold(true).Underline(true)
	keyStyle    = lipgloss.NewStyle().Align(lipgloss.Right).Margin(1, 2, 1, 0)
	valueStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).Padding(0, 1, 0)
	buttonStyle = lipgloss.NewStyle().Align(lipgloss.Center).Border(lipgloss.RoundedBorder(), true).Padding(0, 1, 0)
)

func (m Model) optionsView() string {
	var s strings.Builder

	kvs := [][]string{
		{
			addOptions[addOptionRootFolder],
			m.series.RootFolderPath,
		},
		{
			addOptions[addOptionMonitor],
			string(m.series.AddOptions.Monitor),
		},
		{
			addOptions[addOptionQualityProfile],
			m.client.GetQualityProfileByID(m.series.QualityProfileID).Name,
		},
		{
			addOptions[addOptionLanguageProfile],
			m.client.GetLanguageProfileByID(m.series.LanguageProfileID).Name,
		},
		{
			addOptions[addOptionSeriesType],
			string(m.series.SeriesType),
		},
		{
			addOptions[addOptionSeasonFolder],
			m.seasonFolder.View(),
		},
		{
			addOptions[addOptionSearchForMissingEpisodes],
			m.searchForMissingEpisodes.View(),
		},
		{
			addOptions[addOptionSearchForCutoffUnmetEpisodes],
			m.searchForCutoffUnmetEpisodes.View(),
		},
	}

	lines := make([]string, len(kvs))
	for i, kv := range kvs {
		var color lipgloss.TerminalColor = common.SubtileColor
		if i == int(m.selectedOption)-1 {
			color = lipgloss.Color("#00CCFF")
		}
		lines[i] = lipgloss.JoinHorizontal(lipgloss.Left,
			keyStyle.Width(m.longestOptionWidth).Render(kv[0]),
			valueStyle.Width(m.longestOptionWidth).BorderForeground(color).Render(kv[1]),
		)
	}

	options := lipgloss.JoinVertical(lipgloss.Right,
		lines...,
	)

	width := lipgloss.Width(options)
	s.WriteString(
		titleStyle.Width(width).Render(fmt.Sprintf("%s (%d)", m.series.Title, m.series.Year)),
	)
	s.WriteByte('\n')
	s.WriteByte('\n')

	s.WriteString(options)

	s.WriteByte('\n')
	s.WriteByte('\n')

	var color lipgloss.TerminalColor = common.SubtileColor
	if m.selectedOption == addOptionAddSeries {
		color = lipgloss.Color("#00CCFF")
	}
	s.WriteString(
		lipgloss.Place(width, 1, lipgloss.Center,
			lipgloss.Top, buttonStyle.BorderForeground(color).Render("Add Series")),
	)

	return boxStyle.MaxWidth(m.Width).Render(s.String())
}

func (m *Model) SetSize(width, height int) {
	width -= boxStyle.GetHorizontalFrameSize()
	height -= boxStyle.GetVerticalFrameSize()

	m.Width = width
	m.Height = height

	m.rootFolder.SetSize(width, height)
	m.monitor.SetSize(width, height)
	m.qualityProfile.SetSize(width, height)
	m.languageProfile.SetSize(width, height)
	m.seriesType.SetSize(width, height)
}
