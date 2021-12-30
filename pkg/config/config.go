package config

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/fedeztk/got/pkg/model"
	"github.com/spf13/viper"
)

type Config struct {
    languages []list.Item
    sourceLang, targetLang string
}

func NewConfig() *Config {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("$HOME/.config/got")
    return &Config{
        languages: readLanguages(),
        sourceLang: viper.GetString("source"),
        targetLang: viper.GetString("target"),
    }
}

func readLanguages() []list.Item {
    err := viper.ReadInConfig()
    if err != nil {
        panic(err)
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
    for _, item := range w {
        viper.Set(item.key, item.value)
    }
    viper.WriteConfig()
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

func (c *Config) RememberLastLangs (source, target string) {
    writeConfig(
        write{"source", source},
        write{"target", target},
    )
}
