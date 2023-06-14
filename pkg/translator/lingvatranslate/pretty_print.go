package lingvatranslate

import (
	"fmt"
	"strings"

	"github.com/fedeztk/got/pkg/translator/utils"
)

func (r Response) ShortTranslatedText() string {
	return r.Translation
}

func (r Response) PrettyPrint() string {
	builder := strings.Builder{}
	builder.WriteString(utils.Title.Render("Translated text: "+r.Translation) + "\n")
	if pronunciation := r.Info.Pronunciation; pronunciation.Query != "" {
		builder.WriteString(utils.TitleSecAlt.Render("Pronunciation: "+pronunciation.Query) + "\n")
	}
	for _, def := range r.Info.Definitions {
		if t := def.Type; t != "" {
			builder.WriteString(utils.TitleSec.Render("Part of speech: "+t) + "\n")
		}
		for _, list := range def.List {
			if definition := list.Definition; definition != "" {
				builder.WriteString(utils.IndentTwo.Render("Definition:"))
				builder.WriteString("\n" + utils.IndentThree.Render("- "+definition) + "\n")
			}
			if example := list.Example; example != "" {
				builder.WriteString(utils.IndentTwo.Render("Example:"))
				builder.WriteString("\n" + utils.IndentThree.Render("- "+example) + "\n")
			}
			if synonyms := list.Synonyms; len(synonyms) > 0 {
				builder.WriteString(utils.IndentTwo.Render("Synonyms:"))
				builder.WriteString("\n" + utils.IndentThree.Render("- "))
				builder.WriteString(utils.PrintList(synonyms))
			}
			if field := list.Field; field != "" {
				builder.WriteString(utils.IndentTwo.Render("Field:"))
				builder.WriteString("\n" + utils.IndentThree.Render("- "+field) + "\n")
			}
			builder.WriteString("\n")
		}
	}
	if examples := r.Info.Examples; len(examples) > 0 {
		builder.WriteString(utils.TitleSecAlt2.Render("Examples:") + "\n")
		for _, example := range examples {
			sanitizedExample := strings.ReplaceAll(example, "<b>", "")
			sanitizedExample = strings.ReplaceAll(sanitizedExample, "</b>", "")
			builder.WriteString(utils.IndentTwo.Render("- "+sanitizedExample) + "\n")
		}
	}
	for _, similar := range r.Info.Similar {
		if similar == nil {
			continue
		}
		builder.WriteString(utils.TitleSec.Render(fmt.Sprintf("%v", similar)) + "\n")
	}
	if extraTranslations := r.Info.ExtraTranslations; len(extraTranslations) > 0 {
		builder.WriteString("\n" + utils.Title.Render("Extra translations:") + "\n")
		for _, extraTranslation := range extraTranslations {
			if t := extraTranslation.Type; t != "" {
				builder.WriteString(utils.TitleSec.Render("Part of speech: "+t) + "\n")
			}
			for _, list := range extraTranslation.List {
				if word := list.Word; word != "" {
					builder.WriteString(utils.IndentTwo.Render("Word:"))
					builder.WriteString("\n" + utils.IndentThree.Render("- "+word) + "\n")
				}
				if meanings := list.Meanings; len(meanings) > 0 {
					builder.WriteString(utils.IndentTwo.Render("Meaning:"))
					builder.WriteString("\n" + utils.IndentThree.Render("- "))
					builder.WriteString(utils.PrintList(meanings) + "\n")
				}
			}
		}
	}

	return builder.String()
}
