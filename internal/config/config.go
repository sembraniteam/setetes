package config

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
			SSLMode               string        `mapstructure:"ssl_mode"`
			SSLCert               string        `mapstructure:"ssl_cert"`
			MaxOpenConnections    int           `mapstructure:"max_open_connections"`
			MaxIdleConnections    int           `mapstructure:"max_idle_connections"`
			ConnectionMaxLifetime time.Duration `mapstructure:"connection_max_lifetime"`
			ConnectionMaxIdleTime time.Duration `mapstructure:"connection_max_idle_time"`
		} `mapstructure:"postgres"`

		Redis struct {
			Host            string        `mapstructure:"host"`
			Port            int           `mapstructure:"port"`
			Username        string        `mapstructure:"username"`
			Password        string        `mapstructure:"password"`
			Database        int           `mapstructure:"database"`
			MaxRetries      int           `mapstructure:"max_retries"`
			DialTimeout     time.Duration `mapstructure:"dial_timeout"`
			ReadTimeout     time.Duration `mapstructure:"read_timeout"`
			WriteTimeout    time.Duration `mapstructure:"write_timeout"`
			SSLCert         string        `mapstructure:"ssl_cert"`
			PoolSize        int           `mapstructure:"pool_size"`
			PoolTimeout     time.Duration `mapstructure:"pool_timeout"`
			MinIdleConns    int           `mapstructure:"min_idle_conns"`
			MaxIdleConns    int           `mapstructure:"max_idle_conns"`
			ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
			ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
		} `mapstructure:"redis"`

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

		ED25519 struct {
			PrivateKeyPath string `mapstructure:"private_key_path"`
			PublicKeyPath  string `mapstructure:"public_key_path"`
		} `mapstructure:"ed25519"`
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
