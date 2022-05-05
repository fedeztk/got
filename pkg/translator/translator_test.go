package translator

import (
	"encoding/json"
	"testing"
)

func TestTranslation(t *testing.T) {
	testCases := []struct {
		text, source, target string
	}{
		{"coperchio", "it", "en"},
		{"Hello World!", "en", "it"},
		{"corsa", "it", "en"},
		{"exit", "en", "it"},
	}
	testEngines := []string{"google", "deepl", "libre", "iciba", "reverso"}
	for _, tc := range testCases {
		for _, engine := range testEngines {
			res, err := Translate(tc.text, tc.source, tc.target, engine)
			if err != nil {
				t.Error(err)
			}
			resJSON, err := json.MarshalIndent(res, "", "    ")
			if err != nil {
				t.Error(err)
			}
			t.Log(string(resJSON))
		}
	}
}

func TestPrettyPrint(t *testing.T) {
	testCases := []struct {
		text, source, target string
	}{
		{"coperchio", "it", "en"},
		{"Hello World!", "en", "it"},
		{"corsa", "it", "en"},
		{"exit", "en", "it"},
		{"the", "en", "it"},
	}
	testEngines := []string{"google", "deepl", "libre", "iciba", "reverso"}
	for _, tc := range testCases {
		for _, engine := range testEngines {
			res, err := Translate(tc.text, tc.source, tc.target, engine)
			if err != nil {
				t.Error(err)
			}
			t.Log(res.PrettyPrint())
		}
	}
}
