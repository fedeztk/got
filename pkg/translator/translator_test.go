package translator

import (
	"encoding/json"
	"testing"
)

func TestTranslation(t *testing.T) {
    var resJSON []byte

	res, err := Translate("Hello World!", "en", "it")
	if err != nil {
		t.Error(err)
	}
    resJSON, err = json.MarshalIndent(res, "", "    ")
	if err != nil {
        panic(err)
	}
	t.Log(string(resJSON))

	res, err = Translate("coperchio", "it", "en")
	if err != nil {
		t.Error(err)
	}
	resJSON, err = json.MarshalIndent(res, "", "    ")
	if err != nil {
        panic(err)
	}
	t.Log(string(resJSON))

	res, err = Translate("corsa", "it", "en")
	if err != nil {
		t.Error(err)
	}
	resJSON, err = json.MarshalIndent(res, "", "    ")
	if err != nil {
        panic(err)
	}
	t.Log(string(resJSON))
}
