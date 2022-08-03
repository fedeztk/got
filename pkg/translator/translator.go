// Package translator provides a simple api for simplytranslate
package translator

import (
	_ "embed"
	"encoding/json"
	"errors"
	"io/ioutil"
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

func GetAllLanguages() map[string]string {
	return languageMap
}

func Translate(text, source, target, engine string) (Response, error) {
	r := Response{}
	t := struct {
		client    http.Client
		languages map[string]string
		engines   []string
	}{
		client:    http.Client{},
		languages: languageMap,
		engines:   []string{"google", "deepl", "libre", "iciba", "reverso"},
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

	checkEngine := func(engine string) bool {
		for _, e := range t.engines {
			if e == engine {
				return true
			}
		}
		return false
	}
	if !checkEngine(engine) {
		return r, errors.New("engine not supported, please use one of the following: " + strings.Join(t.engines, ", "))
	}

	const translateURL = "https://simplytranslate.org/api/translate/?engine="
	req, err := http.NewRequest("GET", translateURL+engine, nil)
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

	if res.StatusCode != 200 {
		return r, errors.New("Unable to translate! Status code received from server: " + res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return r, err
	}

	return r, nil
}

type Response struct {
	DefinitionsByCategory map[string][]struct {
		Definition    string              `json:"definition,omitempty"`
		Dictionary    string              `json:"dictionary,omitempty"`
		UseInSentence string              `json:"use-in-sentence,omitempty"`
		Synonyms      map[string][]string `json:"synonyms,omitempty"`
		Informal      string              `json:"informal,omitempty"`
	} `json:"definitions"`
	TranslatedText    string `json:"translated-text"`
	SingleTranslation map[string]map[string]struct {
		Words     []string `json:"words,omitempty"`
		Frequency string   `json:"frequency,omitempty"`
	} `json:"translations,omitempty"`
}

func TextToSpeech(text, lang string) ([]byte, error) {
	var r []byte
	t := struct {
		client    http.Client
		languages map[string]string
		engine    string
	}{
		client:    http.Client{},
		languages: languageMap,
		engine:    "google",
	}

	// only google tts is supported
	engine := t.engine

	if text == "" {
		return r, nil
	}
	if lang == "" {
		lang = "en"
	}

	if _, ok := t.languages[lang]; !ok {
		return r, errors.New("language not supported")
	}

	const ttsURL = "https://simplytranslate.org/api/tts/?engine="
	req, err := http.NewRequest("GET", ttsURL+engine, nil)
	if err != nil {
		return r, err
	}

	query := req.URL.Query()
	query.Add("lang", lang)
	query.Add("text", text)
	req.URL.RawQuery = query.Encode()

	res, err := t.client.Do(req)
	if err != nil {
		return r, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return r, errors.New("Unable to get TextToSpeech! Status code received from server: " + res.Status)
	}

	r, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return r, err
	}

	return r, nil
}
