package model

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyBindingMgr struct {
	Bindings map[int][]key.Binding
	state    int
}

func (kbm keyBindingMgr) ShortHelp() []key.Binding {
	if kbm.state != CHOOSING {
		return append(globalKeyMap, kbm.Bindings[kbm.state]...)
	}
	return append(globalKeyMap, key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle full help")))
}

// only used when in CHOOSING state
func (kbm keyBindingMgr) FullHelp() [][]key.Binding {
	keys := append(globalKeyMap, kbm.Bindings[kbm.state]...)
	// group them 2 per line
	groups := [][]key.Binding{}
	for i := 0; i < len(keys); i += 2 {
		groups = append(groups, keys[i:i+2])
	}
	return groups
}

func newKeyBindingMgr(listKeyMaps [][]key.Binding) keyBindingMgr {
	gbm := keyBindingMgr{
		Bindings: make(map[int][]key.Binding, 4),
	}
	gbm.Bindings[TYPING] = typingKeyMap
	gbm.Bindings[LOADING] = gbm.Bindings[TYPING]  // no particular keys for loading
	gbm.Bindings[TRANSLATING] = translatingKeyMap // provided by the list component

	// get keys from bubbles.list component
	listMapping := []key.Binding{}
	for _, list := range listKeyMaps {
		for _, k := range list {
			if k.Enabled() {
				listMapping = append(listMapping, k)
			}
		}
	}
	gbm.Bindings[CHOOSING] = listMapping

	return gbm
}

var (
	globalKeyMap = []key.Binding{
		key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next tab"),
		),
		key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift-tab", "previous tab"),
		),
		key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc/ctrl+c", "exit"),
		),
	}

	typingKeyMap = []key.Binding{
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "submit"),
		),
	}

	translatingKeyMap = []key.Binding{
		key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "copy to clipboard"),
		),
		key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "play translation"),
		),
	}

	getListAdditionalKeyMap = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "choose source language"),
			),
			key.NewBinding(
				key.WithKeys("t"),
				key.WithHelp("t", "choose target language"),
			),
			key.NewBinding(
				key.WithKeys("i"),
				key.WithHelp("i", "invert languages"),
			),
		}
	}
)
