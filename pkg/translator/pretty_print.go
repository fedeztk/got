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
		markdown.WriteString("## " + "Part of speech: *" + category + "*\n")
		for _, def := range defByCategory {
			if definition := def.Definition; definition != "" {
				markdown.WriteString("> Definition:\n\t\t " + definition + "\n")
			}
			if dictionary := def.Dictionary; dictionary != "" {
				markdown.WriteString("> Dictionary: \n\t\t" + dictionary + "\n")
			}
			if useInSentence := def.UseInSentence; useInSentence != "" {
				markdown.WriteString("> Use in sentence: \n\t\t " + useInSentence + "\n")
			}
			for key, synonymsList := range def.Synonyms {
				markdown.WriteString("> Synonyms: \n\t\t - ")
				for _, synonym := range synonymsList[:len(synonymsList)-1] {
					markdown.WriteString(synonym + ", ")
				}
				markdown.WriteString(synonymsList[len(synonymsList)-1])
				if key != "" {
					markdown.WriteString(" (" + key + ")")
				}
				markdown.WriteString("\n")
			}
			if informal := def.Informal; informal != "" {
				markdown.WriteString("> Informal: " + informal + "\n")
			}
			markdown.WriteString("\n\n")
		}
		markdown.WriteString("\n" + "---" + "\n")
	}
	for category, translationsByCategory := range r.SingleTranslation {
		markdown.WriteString("## " + "Part of speech: *" + category + "*\n")
		for singleWord, translations := range translationsByCategory {
			markdown.WriteString("- **" + singleWord + "**: ")
			for _, word := range translations.Words[:len(translations.Words)-1] {
				markdown.WriteString(word + ", ")
			}
			markdown.WriteString(translations.Words[len(translations.Words)-1] + "\n")
		}
	}
	render, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(120),
	)
	pretty, err := render.Render(markdown.String())
	if err != nil {
		panic(err)
	}
	return pretty
}
