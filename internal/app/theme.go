package app

import "github.com/charmbracelet/lipgloss"

// Theme holds all styles used by the application.
type Theme struct {
	StatusBar    lipgloss.Style
	StatusText   lipgloss.Style
	SearchBar    lipgloss.Style
	SearchPrompt lipgloss.Style
	LineNumber   lipgloss.Style
	MatchHL      lipgloss.Style
	Normal       lipgloss.Style
}

// DefaultTheme returns the default color theme.
func DefaultTheme() Theme {
	return Theme{
		StatusBar: lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("252")).
			Padding(0, 1),
		StatusText: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")),
		SearchBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")),
		SearchPrompt: lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true),
		LineNumber: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Width(6).
			Align(lipgloss.Right).
			MarginRight(1),
		MatchHL: lipgloss.NewStyle().
			Background(lipgloss.Color("214")).
			Foreground(lipgloss.Color("0")).
			Bold(true),
		Normal: lipgloss.NewStyle(),
	}
}
