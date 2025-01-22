package config

import (
	"log/slog"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Dictionaries            []string
	DictionaryDirectory     string
	CacheDirectory          string
	ListenAddr              string
	FallbackToGoogle        bool
	UpdateSchedule          string
	onConfigChangeCallbacks []func()
}

var (
	instance *Config
	once     sync.Once
)

// GetInstance provides access to the singleton instance
func Shared() *Config {
	once.Do(func() {
		viper.SetConfigName("tskks")
		viper.AddConfigPath("$HOME/.config")
		viper.AddConfigPath(".")

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
