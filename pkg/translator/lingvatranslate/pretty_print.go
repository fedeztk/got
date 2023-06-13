package lingvatranslate

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
	return r.Translation
}

func (r Response) PrettyPrint() string {
	builder := strings.Builder{}
	builder.WriteString(title.Render("Translated text: "+r.Translation) + "\n")
	builder.WriteString(titleSec.Render("Pronunciation: "+r.Info.Pronunciation.Query) + "\n")
	builder.WriteString(titleSec.Render("Definitions: ") + "\n")
	for _, def := range r.Info.Definitions {
		builder.WriteString(indentTwo.Render("Type: "+def.Type) + "\n")
		for _, list := range def.List {
			builder.WriteString(indentThree.Render("Definition: "+list.Definition) + "\n")
			builder.WriteString(indentThree.Render("Example: "+list.Example) + "\n")
			// builder.WriteString(indentThree.Render("Synonyms: ") + "\n")
			// for _, synonym := range list.Synonyms {
			//   builder.WriteString(indentThree.Render("- "+synonym) + "\n")
			// }
			builder.WriteString(indentThree.Render("Field: "+list.Field) + "\n")
		}
	}
	builder.WriteString(titleSec.Render("Examples: ") + "\n")
	for _, example := range r.Info.Examples {
		sanitizedExample := strings.ReplaceAll(example, "<b>", "")
		sanitizedExample = strings.ReplaceAll(sanitizedExample, "</b>", "")
		builder.WriteString(indentTwo.Render("- "+sanitizedExample) + "\n")
	}
	builder.WriteString(titleSec.Render("Similar: ") + "\n")
	// for _, similar := range r.Info.Similar {
	//   builder.WriteString(indentTwo.Render("- "+similar) + "\n")
	// }
	builder.WriteString(titleSec.Render("Extra translations: ") + "\n")
	for _, extraTranslation := range r.Info.ExtraTranslations {
		builder.WriteString(indentTwo.Render("Type: "+extraTranslation.Type) + "\n")
		for _, list := range extraTranslation.List {
			builder.WriteString(indentThree.Render("Word: "+list.Word) + "\n")
			builder.WriteString(indentThree.Render("Meanings: ") + "\n")
			for _, meaning := range list.Meanings {
				builder.WriteString(indentThree.Render("- "+meaning) + "\n")
			}
		}
	}
	builder.WriteString(titleSec.Render("Synonyms: ") + "\n")
	// for _, synonym := range r.Info.Synonyms {
	//   builder.WriteString(indentTwo.Render("- "+synonym) + "\n")
	// }
	// builder.WriteString(titleSec.Render("Antonyms: ") + "\n")
	// for _, antonym := range r.Info.Antonyms {
	//   builder.WriteString(indentTwo.Render("- "+antonym) + "\n")
	// }
	// builder.WriteString(titleSec.Render("Definitions: ") + "\n")
	// for _, definition := range r.Info.Definitions {
	//   builder.WriteString(indentTwo.Render("- "+definition) + "\n")
	// }
	builder.WriteString(titleSec.Render("Translations: ") + "\n")

	return builder.String()
	// builder := strings.Builder{}
	// builder.WriteString(title.Render("Translated text: "+r.TranslatedText) + "\n")
	// for category, defByCategory := range r.DefinitionsByCategory {
	// 	builder.WriteString(titleSec.Render("Part of speech: "+
	// 		getPartOfSpeechOrUndefined(category)) + "\n")
	// 	for _, def := range defByCategory {
	// 		if definition := def.Definition; definition != "" {
	// 			builder.WriteString(indentTwo.Render("Definition:"))
	// 			builder.WriteString("\n" + indentThree.Render("- "+definition) + "\n")
	// 		}
	// 		if dictionary := def.Dictionary; dictionary != "" {
	// 			builder.WriteString(indentTwo.Render("Dictionary:"))
	// 			builder.WriteString("\n" + indentThree.Render("- "+dictionary) + "\n")
	// 		}
	// 		if useInSentence := def.UseInSentence; useInSentence != "" {
	// 			builder.WriteString(indentTwo.Render("Use in sentence:"))
	// 			builder.WriteString("\n" + indentThree.Render("- "+useInSentence) + "\n")
	// 		}
	// 		for key, synonymsList := range def.Synonyms {
	// 			builder.WriteString(indentTwo.Render("Synonyms:"))
	// 			builder.WriteString("\n" + indentThree.Render("- "))
	// 			for _, synonym := range synonymsList[:len(synonymsList)-1] {
	// 				builder.WriteString(synonym + ", ")
	// 			}
	// 			builder.WriteString(synonymsList[len(synonymsList)-1])
	// 			if key != "" {
	// 				builder.WriteString(" (" + key + ")")
	// 			}
	// 			builder.WriteString("\n")
	// 		}
	// 		if informal := def.Informal; informal != "" {
	// 			builder.WriteString(indentTwo.Render("Informal: "+informal) + "\n")
	// 		}
	// 		builder.WriteString("\n")
	// 	}
	// }
	// for category, translationsByCategory := range r.SingleTranslation {
	// 	builder.WriteString(titleSec.Render("Part of speech: "+
	// 		getPartOfSpeechOrUndefined(category)) + "\n")
	// 	for singleWord, translations := range translationsByCategory {
	// 		builder.WriteString(listItem.Render("- "+singleWord) + ": ")
	// 		for _, word := range translations.Words[:len(translations.Words)-1] {
	// 			builder.WriteString(word + ", ")
	// 		}
	// 		builder.WriteString(translations.Words[len(translations.Words)-1] + "\n")
	// 	}
	// }
	// return builder.String()
}

func getPartOfSpeechOrUndefined(s string) string {
	if s == "" || s == "null" {
		return "undefined"
	}
	return s
}
