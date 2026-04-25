package mr

import "github.com/charmbracelet/lipgloss"

// Styles for the list UI
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("212")).
			Bold(true)

	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	DimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	BorderColor = lipgloss.Color("62")

	SectionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true)

	CodeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	AuthTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true)
)
