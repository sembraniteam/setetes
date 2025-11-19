package internal

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Postgres struct {
		Host                  string        `mapstructure:"host"`
		Port                  string        `mapstructure:"port"`
		Username              string        `mapstructure:"username"`
		Password              string        `mapstructure:"password"`
		Database              string        `mapstructure:"database"`
		SSLMode               string        `mapstructure:"sslmode"`
		SSLCert               string        `mapstructure:"sslcert"`
		MaxOpenConnections    int           `mapstructure:"max_open_connections"`
		MaxIdleConnections    int           `mapstructure:"max_idle_connections"`
		ConnectionMaxLifetime time.Duration `mapstructure:"connection_max_lifetime"`
		ConnectionMaxIdleTime time.Duration `mapstructure:"connection_max_idle_time"`
	} `mapstructure:"postgres"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigType("yml")
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		var newConfig Config
		if err := viper.Unmarshal(&newConfig); err != nil {
			panic(err)
		}

		config = newConfig
	})

	viper.WatchConfig()

	return &config, nil
}
