package postgresx

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sembraniteam/setetes/internal"
	"github.com/sembraniteam/setetes/internal/ent"
)

type (
	client struct {
		config internal.Config
	}

	Postgres interface {
		Connect() (*ent.Client, error)
	}
)

func New() Postgres {
	return client{config: *internal.Get()}
}

func (c client) Connect() (*ent.Client, error) {
	postgres := c.config.Postgres

	db, err := sql.Open("pgx", c.dsn())
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(postgres.MaxIdleConnections)
	db.SetMaxOpenConns(postgres.MaxOpenConnections)
	db.SetConnMaxLifetime(postgres.ConnectionMaxLifetime * time.Minute)
	db.SetConnMaxIdleTime(postgres.ConnectionMaxIdleTime * time.Minute)

	drv := entsql.OpenDB(dialect.Postgres, db)

	return ent.NewClient(ent.Driver(drv)), nil
}

func (c client) dsn() string {
	q := url.Values{}
	postgres := c.config.Postgres
	q.Add("sslmode", postgres.SSLMode)

	if postgres.SSLMode == "verify-ca" || postgres.SSLMode == "verify-full" {
		q.Add("sslrootcert", postgres.SSLCert)
	}

	u := &url.URL{
		Scheme:   dialect.Postgres,
		User:     url.UserPassword(postgres.Username, postgres.Password),
		Host:     fmt.Sprintf("%s:%d", postgres.Host, postgres.Port),
		Path:     postgres.Database,
		RawQuery: q.Encode(),
	}

	return u.String()
}
