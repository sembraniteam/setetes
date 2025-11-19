package bootstrap

import (
	"net/http"

	"github.com/megalodev/setetes/internal"
	"github.com/megalodev/setetes/internal/ent"
)

type App struct {
	Config internal.Config
	DB     *ent.Client
	Server *http.Server
}

func NewApp(config internal.Config) *App {
	return &App{Config: config}
}
