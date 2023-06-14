package utils

type BackendResponse interface {
	PrettyPrint() string
	ShortTranslatedText() string
}
