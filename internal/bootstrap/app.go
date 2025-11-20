package bootstrap

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/gzip"
	"github.com/megalodev/setetes/internal"
	"github.com/megalodev/setetes/internal/database/postgresx"
	"github.com/megalodev/setetes/internal/httpx"
	"github.com/megalodev/setetes/internal/httpx/middleware"
	"github.com/megalodev/setetes/internal/httpx/web"
	"github.com/samber/do/v2"
)

type (
	App struct {
		configPath string
	}

	Bootstrap interface {
		Init() error
	}
)

func New(configPath string) Bootstrap {
	return App{configPath: configPath}
}

func (a App) Init() error {
	config, err := internal.LoadConfig(a.configPath)
	if err != nil {
		return err
	}

	pdb := postgresx.New(*config)
	_, err = pdb.Connect()
	if err != nil {
		return err
	}

	injector := do.New()
	rateLimiter := middleware.Default()
	defer rateLimiter.Stop()

	server := httpx.NewServer(
		config,
		httpx.Injector(injector),
		httpx.Middlewares(
			gzip.Gzip(gzip.DefaultCompression),
			middleware.Timeout(),
			middleware.RateLimitByIP(rateLimiter),
			middleware.RequestID(),
		),
		httpx.UseRouter(web.Routes),
	)
	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.Run()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case _ = <-quit:
		fmt.Println("shutting down...")

	case err = <-serverErrors:
		if err != nil {
			fmt.Println("server exited with error:", err)
		}
	}

	if err = server.Stop(); err != nil {
		return fmt.Errorf("force to shutdown server because error: %w", err)
	}

	return nil
}
