package bootstrap

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/gzip"
	"github.com/samber/do/v2"
	"github.com/sembraniteam/setetes/internal/config"
	"github.com/sembraniteam/setetes/internal/cryptox"
	"github.com/sembraniteam/setetes/internal/cryptox/pasetox"
	"github.com/sembraniteam/setetes/internal/database/postgresx"
	"github.com/sembraniteam/setetes/internal/ent"
	"github.com/sembraniteam/setetes/internal/httpx"
	"github.com/sembraniteam/setetes/internal/httpx/handler"
	"github.com/sembraniteam/setetes/internal/httpx/middleware"
	"github.com/sembraniteam/setetes/internal/httpx/web"
	"github.com/sembraniteam/setetes/internal/rbac"
	"github.com/sembraniteam/setetes/internal/seed"
	"github.com/sembraniteam/setetes/internal/service"
)

type (
	App struct {
		configPath string
	}

	Bootstrap interface {
		Init() error
		Seeder() error
	}
)

func New(configPath string) Bootstrap {
	return App{configPath: configPath}
}

func (a App) Init() error {
	injector := do.New(service.Packages, handler.Packages)
	_, err := config.LoadConfig(a.configPath)
	if err != nil {
		return err
	}

	pdb := postgresx.New()
	pcl, err := pdb.Connect()
	if err != nil {
		return err
	}

	do.Provide[*ent.Client](
		injector,
		func(_ do.Injector) (*ent.Client, error) {
			return pcl, nil
		})

	rateLimiter := middleware.DefaultTokenBucket()
	defer rateLimiter.Stop()

	rm, err := rbac.New(pcl)
	if err != nil {
		return err
	}

	do.Provide[*rbac.Manager](
		injector,
		func(_ do.Injector) (*rbac.Manager, error) {
			return rm, nil
		},
	)

	keypair, err := cryptox.DefaultKeypair()
	if err != nil {
		return err
	}

	verifier := pasetox.NewVerifier(keypair)

	auth := middleware.NewAuthorizationConfig(
		rm,
		verifier,
		web.PublicRoutes(),
	)

	server := httpx.NewServer(
		httpx.Injector(injector),
		httpx.Middlewares(
			gzip.Gzip(gzip.DefaultCompression),
			middleware.Timeout(),
			middleware.RateLimitByIP(rateLimiter),
			middleware.RequestID(),
			auth.Authorization(),
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

func (a App) Seeder() error {
	injector := do.New(seed.Packages)
	_, err := config.LoadConfig(a.configPath)
	if err != nil {
		return err
	}

	pdb := postgresx.New()
	pcl, err := pdb.Connect()
	if err != nil {
		return err
	}

	do.Provide[*ent.Client](injector, func(_ do.Injector) (*ent.Client, error) {
		return pcl, nil
	})

	rbacMan, err := rbac.New(pcl)
	if err != nil {
		return err
	}

	do.Provide[*rbac.Manager](
		injector,
		func(_ do.Injector) (*rbac.Manager, error) {
			return rbacMan, nil
		},
	)

	seeder := do.MustInvoke[seed.Seeder](injector)
	seeder.RunAll()

	return nil
}
