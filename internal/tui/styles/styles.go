package styles

import "github.com/charmbracelet/lipgloss"

var (
	// universal
	ErrorColor  = lipgloss.AdaptiveColor{Light: "#F08F89", Dark: "#F08F89"}
	OkColor     = lipgloss.AdaptiveColor{Light: "#89F0CB", Dark: "#89F0CB"}
	SubtleColor = lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}
	BlueColor   = lipgloss.AdaptiveColor{Light: "#11488f", Dark: "#11488f"}
	PurpleColor = lipgloss.Color("#7B61FF")

	// sonarr
	SonarrBlue = lipgloss.Color("#00CCFF")
)
