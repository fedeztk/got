package config

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/fedeztk/got/internal/model"
	"github.com/spf13/viper"
)

const (
	REPO       = "https://github.com/fedeztk/got"
	REPOCONFIG = "https://github.com/fedeztk/got/blob/master/config.yml"
)

type Config struct {
	languages              []list.Item
	sourceLang, targetLang string
}

func NewConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/got")
	return &Config{
		languages:  readLanguages(),
		sourceLang: viper.GetString("source"),
		targetLang: viper.GetString("target"),
	}
}

func readLanguages() []list.Item {
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("%s\ngot needs a config file to work! Copy the sample in ~/.config/got/config.yml from %s", err, REPOCONFIG))
	}

	items := make([]list.Item, 0)

	for _, line := range strings.Split(viper.GetString("languages"), "\n") {
		if line[0] == '#' || len(line) == 0 {
			continue
		}
		fields := strings.Split(line, "-")
		items = append(items, model.NewListItem(strings.TrimSpace(fields[0]), strings.TrimSpace(fields[1])))
	}
	return items
}

type write struct {
	key, value string
}

func writeConfig(w ...write) {
	modified := false
	for _, item := range w {
		if viper.GetString(item.key) != item.value {
			viper.Set(item.key, item.value)
			modified = true
		}
	}
	if modified {
		viper.WriteConfig()
	}
}

func (c *Config) Langs() []list.Item {
	return c.languages
}

func (c *Config) Source() string {
	return c.sourceLang
}

func (c *Config) Target() string {
	return c.targetLang
}

func (c *Config) RememberLastLangs(source, target string) {
	writeConfig(
		write{"source", source},
		write{"target", target},
	)
}
