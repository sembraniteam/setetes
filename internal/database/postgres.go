package database

import (
	"net/url"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/megalodev/setetes/internal"
	"github.com/megalodev/setetes/internal/ent"
)

type (
	connector struct {
		config internal.Config
	}

	Postgres interface {
		Connect() (*ent.Client, error)
	}
)

func NewPostgres(config internal.Config) Postgres {
	return connector{config: config}
}

func (c connector) Connect() (*ent.Client, error) {
	postgres := c.config.Postgres

	drv, err := entsql.Open(dialect.Postgres, c.dsn())
	if err != nil {
		return nil, err
	}

	db := drv.DB()
	db.SetMaxIdleConns(postgres.MaxIdleConnections)
	db.SetMaxOpenConns(postgres.MaxOpenConnections)
	db.SetConnMaxLifetime(postgres.ConnectionMaxLifetime * time.Hour)
	db.SetConnMaxIdleTime(postgres.ConnectionMaxIdleTime * time.Hour)

	return ent.NewClient(ent.Driver(drv)), nil
}

func (c connector) dsn() string {
	q := url.Values{}
	postgres := c.config.Postgres
	q.Add("sslmode", postgres.SSLMode)

	if postgres.SSLMode == "verify-ca" || postgres.SSLMode == "verify-full" {
		q.Add("sslrootcert", postgres.SSLCert)
	}

	u := &url.URL{
		Scheme:   dialect.Postgres,
		User:     url.UserPassword(postgres.Username, postgres.Password),
		Host:     postgres.Host,
		Path:     postgres.Database,
		RawQuery: q.Encode(),
	}

	return u.String()
}
