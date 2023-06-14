package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	sourceLang, targetLang, engine, backend string
}

func NewConfig() *Config {
	home, _ := os.UserHomeDir()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(home + "/.config/got")

	err := viper.ReadInConfig()
	if err != nil {
		writeDefaultConfig()
	}
	return &Config{
		sourceLang: viper.GetString("source"),
		targetLang: viper.GetString("target"),
		engine:     viper.GetString("engine"),
		backend:    viper.GetString("backend"),
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

func (c *Config) Engine() string {
	return c.engine
}

func (c *Config) Backend() string {
	return c.backend
}

func (c *Config) SetEngine(engine string) {
	c.engine = engine
}

func (c *Config) SetBackend(backend string) {
	c.backend = backend
}

func (c *Config) RememberLastSettings(source, target string) {
	writeConfig(
		write{"source", source},
		write{"target", target},
		write{"engine", c.engine},
		write{"backend", c.backend},
	)
}

func writeDefaultConfig() {
	home, _ := os.UserHomeDir()
	os.MkdirAll(home+"/.config/got", os.ModePerm)
	viper.SafeWriteConfig()
	writeConfig(
		write{"source", "en"},
		write{"target", "it"},
		write{"engine", "google"},
		write{"backend", "lingvatranslate"},
	)
}
