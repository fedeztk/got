package utils

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	IndentOne    = lipgloss.NewStyle().Margin(0, 0, 0, 2)
	IndentTwo    = IndentOne.Copy().Margin(0, 0, 0, 4)
	IndentThree  = IndentTwo.Copy().Margin(0, 0, 0, 6).MaxWidth(80)
	Title        = IndentOne.Copy().Bold(true).Background(lipgloss.Color("12")).Padding(0, 1).Foreground(lipgloss.Color("0"))
	TitleSec     = IndentTwo.Copy().Bold(true).Background(lipgloss.Color("13")).Padding(0, 1).Foreground(lipgloss.Color("0")).MarginBottom(1).MarginTop(1)
	ListItem     = IndentTwo.Copy().Bold(true)
	TitleSecAlt  = IndentTwo.Copy().Bold(true).Background(lipgloss.Color("14")).Padding(0, 1).Foreground(lipgloss.Color("0")).MarginBottom(1).MarginTop(1)
	TitleSecAlt2 = IndentTwo.Copy().Bold(true).Background(lipgloss.Color("11")).Padding(0, 1).Foreground(lipgloss.Color("0")).MarginBottom(1).MarginTop(1)
)

func PrintList(list []string) string {
	builder := strings.Builder{}
	for _, item := range list[:len(list)-1] {
		builder.WriteString(item + ", ")
	}
	builder.WriteString(list[len(list)-1])
	builder.WriteString("\n")
	return builder.String()
}
