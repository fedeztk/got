package utils

import (
	_ "embed"
	"strings"
)

//go:embed languages.txt
var languages string
var languageMap map[string]string = loadLanguageMap()

func loadLanguageMap() map[string]string {
	m := make(map[string]string)
	for _, line := range strings.Split(languages, "\n") {
		parts := strings.Split(strings.ReplaceAll(line, " ", ""), ",")
		if line == "" {
			break
		}
		m[parts[0]] = parts[1]
	}
	return m
}

func GetAllLanguages() map[string]string {
	return languageMap
}
