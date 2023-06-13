package simplytranslate

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	indentOne   = lipgloss.NewStyle().Margin(0, 0, 0, 2)
	indentTwo   = indentOne.Copy().Margin(0, 0, 0, 4)
	indentThree = indentTwo.Copy().Margin(0, 0, 0, 6).MaxWidth(80)
	title       = indentOne.Copy().Bold(true).Background(lipgloss.Color("12")).Padding(0, 1).Foreground(lipgloss.Color("0"))
	titleSec    = indentTwo.Copy().Bold(true).Background(lipgloss.Color("13")).Padding(0, 1).Foreground(lipgloss.Color("0")).MarginBottom(1).MarginTop(1)
	listItem    = indentTwo.Copy().Bold(true)
)

func (r Response) ShortTranslatedText() string {
	return r.TranslatedText
}

func (r Response) PrettyPrint() string {
	builder := strings.Builder{}
	builder.WriteString(title.Render("Translated text: "+r.TranslatedText) + "\n")
	for category, defByCategory := range r.DefinitionsByCategory {
		builder.WriteString(titleSec.Render("Part of speech: "+
			getPartOfSpeechOrUndefined(category)) + "\n")
		for _, def := range defByCategory {
			if definition := def.Definition; definition != "" {
				builder.WriteString(indentTwo.Render("Definition:"))
				builder.WriteString("\n" + indentThree.Render("- "+definition) + "\n")
			}
			if dictionary := def.Dictionary; dictionary != "" {
				builder.WriteString(indentTwo.Render("Dictionary:"))
				builder.WriteString("\n" + indentThree.Render("- "+dictionary) + "\n")
			}
			if useInSentence := def.UseInSentence; useInSentence != "" {
				builder.WriteString(indentTwo.Render("Use in sentence:"))
				builder.WriteString("\n" + indentThree.Render("- "+useInSentence) + "\n")
			}
			for key, synonymsList := range def.Synonyms {
				builder.WriteString(indentTwo.Render("Synonyms:"))
				builder.WriteString("\n" + indentThree.Render("- "))
				for _, synonym := range synonymsList[:len(synonymsList)-1] {
					builder.WriteString(synonym + ", ")
				}
				builder.WriteString(synonymsList[len(synonymsList)-1])
				if key != "" {
					builder.WriteString(" (" + key + ")")
				}
				builder.WriteString("\n")
			}
			if informal := def.Informal; informal != "" {
				builder.WriteString(indentTwo.Render("Informal: "+informal) + "\n")
			}
			builder.WriteString("\n")
		}
	}
	for category, translationsByCategory := range r.SingleTranslation {
		builder.WriteString(titleSec.Render("Part of speech: "+
			getPartOfSpeechOrUndefined(category)) + "\n")
		for singleWord, translations := range translationsByCategory {
			builder.WriteString(listItem.Render("- "+singleWord) + ": ")
			for _, word := range translations.Words[:len(translations.Words)-1] {
				builder.WriteString(word + ", ")
			}
			builder.WriteString(translations.Words[len(translations.Words)-1] + "\n")
		}
	}
	return builder.String()
}

func getPartOfSpeechOrUndefined(s string) string {
	if s == "" || s == "null" {
		return "undefined"
	}
	return s
}
