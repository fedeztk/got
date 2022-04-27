package translator

import (
	"strings"

	"github.com/charmbracelet/glamour"
)

func (r Response) PrettyPrint() string {
	markdown := strings.Builder{}
    markdown.WriteString("# Translated text: **" + r.TranslatedText + "**\n\n")
	markdown.WriteString("\n" + "---" + "\n")
	for category, defByCategory := range r.DefinitionsByCategory {
		markdown.WriteString("## " + "Definition: *" + category + "*\n")
		for _, def := range defByCategory {
			markdown.WriteString("> " + def.Definition + "\n")
			markdown.WriteString("> " + def.Dictionary + "\n")
			markdown.WriteString("> " + def.UseInSentence + "\n")
		}
	}
	markdown.WriteString("\n" + "---" + "\n")
	for category, translationsByCategory := range r.SingleTranslation {
		markdown.WriteString("## " + "Translation: *" + category + "*\n")
		for singleWord, translations := range translationsByCategory {
			markdown.WriteString("- **" + singleWord + "**: ")
			for _, word := range translations.Words[:len(translations.Words)-1] {
				markdown.WriteString(word + ", ")
			}
			markdown.WriteString(translations.Words[len(translations.Words)-1] + "\n")
		}
	}
	pretty, err := glamour.Render(markdown.String(), "auto")
	if err != nil {
		panic(err)
	}
	return pretty
}
