package lingvatranslate

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/fedeztk/got/pkg/translator/utils"
)

type Response struct {
	Translation string `json:"translation,omitempty"`
	Info        struct {
		Pronunciation struct {
			Query string `json:"query,omitempty"`
		} `json:"pronunciation,omitempty"`
		Definitions []struct {
			Type string `json:"type,omitempty"`
			List []struct {
				Definition string `json:"definition,omitempty"`
				Example    string `json:"example,omitempty"`
				Synonyms   []any  `json:"synonyms,omitempty"`
				Field      string `json:"field,omitempty"`
			} `json:"list,omitempty"`
		} `json:"definitions,omitempty"`
		Examples          []string `json:"examples,omitempty"`
		Similar           []any    `json:"similar,omitempty"`
		ExtraTranslations []struct {
			Type string `json:"type,omitempty"`
			List []struct {
				Word      string   `json:"word,omitempty"`
				Meanings  []string `json:"meanings,omitempty"`
				Frequency int      `json:"frequency,omitempty"`
			} `json:"list,omitempty"`
		} `json:"extraTranslations,omitempty"`
	} `json:"info,omitempty"`
}
type LingvaTranslate struct {
	languages map[string]string
	client    http.Client
	baseURL   string
}

func New() LingvaTranslate {
	return LingvaTranslate{
		languages: utils.GetAllLanguages(),
		client:    http.Client{},
		baseURL:   "https://lingva.ml/api/v1/",
	}
}

func (b LingvaTranslate) Translate(text, source, target, engine string) (utils.BackendResponse, error) {
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

	var translateURL = b.baseURL + "/" + source + "/" + target + "/" + url.QueryEscape(text)
	req, err := http.NewRequest("GET", translateURL, nil)
	if err != nil {
		return r, err
	}

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

func (b LingvaTranslate) TextToSpeech(text, lang string) ([]byte, error) {
	var r []byte

	if text == "" {
		return r, nil
	}
	if lang == "" {
		lang = "en"
	}

	if _, ok := b.languages[lang]; !ok {
		return r, errors.New("language not supported")
	}

	var ttsURL = b.baseURL + "/audio/" + lang + "/" + url.QueryEscape(text)
	req, err := http.NewRequest("GET", ttsURL, nil)
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
