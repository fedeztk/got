package simplytranslate

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/fedeztk/got/pkg/translator/utils"
)

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

type SimplyTranslate struct {
	languages map[string]string
	engines   []string
	client    http.Client
	baseURL   string
}

func New() SimplyTranslate {
	return SimplyTranslate{
		languages: utils.GetAllLanguages(),
		engines:   []string{"google", "deepl", "libre", "iciba", "reverso"},
		client:    http.Client{},
		baseURL:   "https://simplytranslate.org/api",
	}
}

func (b SimplyTranslate) Translate(text, source, target, engine string) (utils.BackendResponse, error) {
	r := Response{}

	if text == "" {
		return r, nil
	}
	if source == "" {
		source = "auto"
	}
	if target == "" {
		target = "en"
	}

	if _, ok := b.languages[source]; !ok {
		return r, errors.New("source language not supported")
	}
	if _, ok := b.languages[target]; !ok {
		return r, errors.New("target language not supported")
	}

	checkEngine := func(engine string) bool {
		for _, e := range b.engines {
			if e == engine {
				return true
			}
		}
		return false
	}
	if !checkEngine(engine) {
		return r, errors.New("engine not supported, please use one of the following: " + strings.Join(b.engines, ", "))
	}

	var translateURL = b.baseURL + "/translate/?engine="
	req, err := http.NewRequest("GET", translateURL+engine, nil)
	if err != nil {
		return r, err
	}

	query := req.URL.Query()
	query.Add("from", source)
	query.Add("to", target)
	query.Add("text", text)
	req.URL.RawQuery = query.Encode()

	res, err := b.client.Do(req)
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

func (b SimplyTranslate) TextToSpeech(text, lang string) ([]byte, error) {
	var r []byte
	// only google tts is supported
	engine := b.engines[0]

	if text == "" {
		return r, nil
	}
	if lang == "" {
		lang = "en"
	}

	if _, ok := b.languages[lang]; !ok {
		return r, errors.New("language not supported")
	}

	var ttsURL = b.baseURL + "/api/tts/?engine="
	req, err := http.NewRequest("GET", ttsURL+engine, nil)
	if err != nil {
		return r, err
	}

	query := req.URL.Query()
	query.Add("lang", lang)
	query.Add("text", text)
	req.URL.RawQuery = query.Encode()

	res, err := b.client.Do(req)
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
