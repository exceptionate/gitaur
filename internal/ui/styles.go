package ui

import "github.com/charmbracelet/lipgloss"

var (
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	Label = lipgloss.NewStyle().
		Bold(true)

	Success = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))

	Warning = lipgloss.NewStyle().
		Foreground(lipgloss.Color("214"))

	Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196"))

	Info = lipgloss.NewStyle().
		Foreground(lipgloss.Color("39"))
)
