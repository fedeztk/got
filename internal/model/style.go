package model

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// prompt
	promptStyleIndicator = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	placeholderStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	promptStyleUpperText = lipgloss.NewStyle().Background(lipgloss.Color("6")).Bold(true).MarginLeft(2).Padding(0, 1).Foreground(lipgloss.Color("0"))
	promptStyleSelLang   = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).MarginLeft(2)
	// spinner
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	// list
	titleStyle         = lipgloss.NewStyle().Background(lipgloss.Color("11")).Bold(true).Padding(0, 1).Foreground(lipgloss.Color("0"))
	statusMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).MarginLeft(2).Bold(true)
	// errors
	ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
	// tabs
	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}
	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}
	highlight = lipgloss.Color("4")
	tab       = lipgloss.NewStyle().
			Border(tabBorder, true).
			BorderForeground(highlight).
			Padding(0, 2)
	activeTab = tab.Copy().Border(activeTabBorder, true).Bold(true)
	tabGap    = tab.Copy().
			BorderTop(false).
			BorderLeft(false).
			BorderRight(false)
	// footer
	footerStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	footerTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
)
