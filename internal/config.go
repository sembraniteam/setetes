package internal

import (
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	config = new(Config)
	mu     = new(sync.RWMutex)
)

type (
	Config struct {
		App struct {
			Mode string `mapstructure:"mode"`
			Host string `mapstructure:"host"`
			Port int    `mapstructure:"port"`
		} `mapstructure:"app"`

		Postgres struct {
			Host                  string        `mapstructure:"host"`
			Port                  int           `mapstructure:"port"`
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

		Password struct {
			Pepper string `mapstructure:"pepper"`
			Argon2 struct {
				Memory      uint32 `mapstructure:"memory"`
				Iterations  uint32 `mapstructure:"iterations"`
				Parallelism uint8  `mapstructure:"parallelism"`
				SaltLength  uint32 `mapstructure:"salt_length"`
				KeyLength   uint32 `mapstructure:"key_length"`
			} `mapstructure:"argon2"`
		} `mapstructure:"password"`
	}
)

func Get() *Config {
	mu.RLock()
	defer mu.RUnlock()
	return config
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigType("yml")
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	mu.Lock()
	if err := viper.Unmarshal(&config); err != nil {
		mu.Unlock()
		return nil, err
	}
	mu.Unlock()

	viper.OnConfigChange(func(e fsnotify.Event) {
		var newConfig Config
		if err := viper.Unmarshal(&newConfig); err != nil {
			panic(err)
		}
		mu.Lock()
		*config = newConfig
		mu.Unlock()
	})

	viper.WatchConfig()

	return config, nil
}
