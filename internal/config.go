package internal

import "time"

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
