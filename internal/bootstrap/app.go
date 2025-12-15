package bootstrap

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/gzip"
	"github.com/megalodev/setetes/internal"
	"github.com/megalodev/setetes/internal/database/postgresx"
	"github.com/megalodev/setetes/internal/ent"
	"github.com/megalodev/setetes/internal/httpx"
	"github.com/megalodev/setetes/internal/httpx/handler"
	"github.com/megalodev/setetes/internal/httpx/middleware"
	"github.com/megalodev/setetes/internal/httpx/web"
	"github.com/megalodev/setetes/internal/service"
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
	injector := do.New(service.Packages, handler.Packages)
	config, err := internal.LoadConfig(a.configPath)
	if err != nil {
		return err
	}

	pdb := postgresx.New(*config)
	pcl, err := pdb.Connect()
	if err != nil {
		return err
	}

	do.Provide[*ent.Client](injector, func(_ do.Injector) (*ent.Client, error) {
		return pcl, nil
	})

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
	case <-quit:
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
