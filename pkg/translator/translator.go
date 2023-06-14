// Package translator provides a simple api for simplytranslate and lingvatranslate
package translator

import (
	"errors"

	"github.com/fedeztk/got/pkg/translator/lingvatranslate"
	"github.com/fedeztk/got/pkg/translator/simplytranslate"
	"github.com/fedeztk/got/pkg/translator/utils"
)

type Backend interface {
	Translate(text, source, target, engine string) (utils.BackendResponse, error)
	TextToSpeech(text, language string) ([]byte, error)
}

func NewBackend(backend string) (Backend, error) {
	switch backend {
	case "lingvatranslate":
		return lingvatranslate.New(), nil
	case "simplytranslate":
		return simplytranslate.New(), nil
	default:
		return nil, errors.New("backend not supported, please use one of the following: lingvatranslate, simplytranslate")
	}
}
