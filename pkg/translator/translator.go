// Package translator provides a simple api for simplytranslate
package translator

import (
	_ "embed"
	"encoding/json"
	"errors"
	"net/http"
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

func Translate(text, source, target string) (Response, error) {
	r := Response{}
	t := struct {
		client    http.Client
		languages map[string]string
	}{
		client:    http.Client{},
		languages: languageMap,
	}

	if text == "" {
		return r, nil
	}
	if source == "" {
		source = "auto"
	}
	if target == "" {
		target = "en"
	}

	if _, ok := t.languages[source]; !ok {
		return r, errors.New("source language not supported")
	}
	if _, ok := t.languages[target]; !ok {
		return r, errors.New("target language not supported")
	}

	const translateURL = "https://simplytranslate.org/api/translate/?engine=google"
	req, err := http.NewRequest("GET", translateURL, nil)
	if err != nil {
		return r, err
	}

	query := req.URL.Query()
	query.Add("from", source)
	query.Add("to", target)
	query.Add("text", text)
	req.URL.RawQuery = query.Encode()

	res, err := t.client.Do(req)
	if err != nil {
		return r, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return r, err
	}

	return r, nil
}

type Response struct {
	DefinitionsByCategory map[string][]struct {
		Definition    string `json:"definition,omitempty"`
		Dictionary    string `json:"dictionary,omitempty"`
		UseInSentence string `json:"use-in-sentence,omitempty"`
	} `json:"definitions"`
	TranslatedText    string `json:"translated-text"`
	SingleTranslation map[string]map[string]struct {
		Words     []string `json:"words,omitempty"`
		Frequency string   `json:"frequency,omitempty"`
	} `json:"translations,omitempty"`
}
