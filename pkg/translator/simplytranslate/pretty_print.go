package simplytranslate

import (
	"strings"

	"github.com/fedeztk/got/pkg/translator/utils"
)

func (r Response) ShortTranslatedText() string {
	return r.TranslatedText
}

func (r Response) PrettyPrint() string {
	builder := strings.Builder{}
	builder.WriteString(utils.Title.Render("Translated text: "+r.TranslatedText) + "\n")
	for category, defByCategory := range r.DefinitionsByCategory {
		builder.WriteString(utils.TitleSec.Render("Part of speech: "+
			getPartOfSpeechOrUndefined(category)) + "\n")
		for _, def := range defByCategory {
			if definition := def.Definition; definition != "" {
				builder.WriteString(utils.IndentTwo.Render("Definition:"))
				builder.WriteString("\n" + utils.IndentThree.Render("- "+definition) + "\n")
			}
			if dictionary := def.Dictionary; dictionary != "" {
				builder.WriteString(utils.IndentTwo.Render("Dictionary:"))
				builder.WriteString("\n" + utils.IndentThree.Render("- "+dictionary) + "\n")
			}
			if useInSentence := def.UseInSentence; useInSentence != "" {
				builder.WriteString(utils.IndentTwo.Render("Use in sentence:"))
				builder.WriteString("\n" + utils.IndentThree.Render("- "+useInSentence) + "\n")
			}
			for key, synonymsList := range def.Synonyms {
				builder.WriteString(utils.IndentTwo.Render("Synonyms:"))
				builder.WriteString("\n" + utils.IndentThree.Render("- "))
				if key != "" {
					builder.WriteString(" (" + key + ") ")
				}
				builder.WriteString(utils.PrintList(synonymsList))
				builder.WriteString("\n")
			}
			if informal := def.Informal; informal != "" {
				builder.WriteString(utils.IndentTwo.Render("Informal: "+informal) + "\n")
			}
			builder.WriteString("\n")
		}
	}
	for category, translationsByCategory := range r.SingleTranslation {
		builder.WriteString(utils.TitleSec.Render("Part of speech: "+
			getPartOfSpeechOrUndefined(category)) + "\n")
		for singleWord, translations := range translationsByCategory {
			builder.WriteString(utils.ListItem.Render("- "+singleWord) + ": ")
			builder.WriteString(utils.PrintList(translations.Words))
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
