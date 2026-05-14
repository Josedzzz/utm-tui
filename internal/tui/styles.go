// Package tui handles the screens
package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	primaryColor   = lipgloss.Color("#778DA9")
	secondaryColor = lipgloss.Color("#415A77")
	whiteColor     = lipgloss.Color("#FFFFFF")
	grayColor      = lipgloss.Color("#A0A0A0")
	greenColor     = lipgloss.Color("#00FF00")
	redColor       = lipgloss.Color("#FF0000")
	yellowColor    = lipgloss.Color("#FFFF00")

	// General styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(whiteColor).
			Background(primaryColor).
			Padding(0, 2).
			MarginBottom(1)

	miniLogoStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(secondaryColor).
			PaddingLeft(1).
			MarginBottom(1)

	asciiStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			MarginBottom(1)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				PaddingRight(1).
				Foreground(whiteColor).
				Background(secondaryColor).
				Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(grayColor).
			MarginTop(1).
			Italic(true)

	processingStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Italic(true)

	footerStyle = lipgloss.NewStyle().
			Foreground(grayColor).
			Italic(true).
			MarginTop(1)

	containerStyle = lipgloss.NewStyle().
			Padding(1, 2)

	successStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(redColor).
			Bold(true)

	inputStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

	// VM list specific
	vmStatusStyle = map[string]lipgloss.Style{
		"running":   lipgloss.NewStyle().Foreground(greenColor),
		"stopped":   lipgloss.NewStyle().Foreground(primaryColor),
		"suspended": lipgloss.NewStyle().Foreground(yellowColor),
	}
)
