package config

import (
	"log/slog"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Config struct {
	Dictionaries            []string
	DictionaryDirectory     string
	ListenAddr              string
	UpdateSchedule          string
	FallbackToGoogle        bool
	UseDiskCache            bool
	onConfigChangeCallbacks []func()
}

var (
	instance *Config
	once     sync.Once
)

func Shared() *Config {
	once.Do(func() {
		viper.SetConfigName("toyskkserv")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.config")

		err := viper.ReadInConfig()
		if err != nil {
			slog.Error("Failed to read config file", "err", err)
			panic(err)
		}

		instance = &Config{}
		err = viper.Unmarshal(instance)
		if err != nil {
			slog.Error("Failed to unmarshal config", "err", err)
			panic(err)
		}
		dir, err := homedir.Expand(instance.DictionaryDirectory)
		if err == nil {
			instance.DictionaryDirectory = dir
		} else {
			slog.Error("Failed to expand home dir", "err", err)
		}

		instance.onConfigChangeCallbacks = []func(){}
		viper.OnConfigChange(func(e fsnotify.Event) {
			slog.Info("Configuration file changed", "file", viper.ConfigFileUsed())

			// only support dictionary change
			instance.Dictionaries = viper.GetStringSlice("Dictionaries")

			for _, f := range instance.onConfigChangeCallbacks {
				f()
			}
		})
		viper.WatchConfig()
	})
	return instance
}

func (c *Config) OnConfigChange(f func()) {
	c.onConfigChangeCallbacks = append(c.onConfigChangeCallbacks, f)
}
