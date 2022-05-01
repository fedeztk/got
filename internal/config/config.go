package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	REPO       = "https://github.com/fedeztk/got"
	REPOCONFIG = "https://github.com/fedeztk/got/blob/master/config.yml"
)

type Config struct {
	sourceLang, targetLang string
}

func NewConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/got")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("%s\ngot needs a config file to work! Copy the sample in ~/.config/got/config.yml from %s", err, REPOCONFIG))
	}
	return &Config{
		sourceLang: viper.GetString("source"),
		targetLang: viper.GetString("target"),
	}
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
