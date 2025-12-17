package httpx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/sembraniteam/setetes/internal"
)

const Timeout = time.Second * 30

type (
	Router func(e *gin.Engine, i do.Injector)

	Option func(c *Config)

	Config struct {
		*http.Server
		engine     *gin.Engine
		mode       string
		middleware []gin.HandlerFunc
		routerFunc Router
		injector   do.Injector
	}

	Server interface {
		Run() error
		Stop() error
	}
)

func NewServer(c *internal.Config, opts ...Option) Server {
	config := &Config{
		mode: c.App.Mode,
		Server: &http.Server{
			Addr:              fmt.Sprintf("%s:%d", c.App.Host, c.App.Port),
			ReadTimeout:       Timeout,
			WriteTimeout:      Timeout,
			ReadHeaderTimeout: Timeout,
			IdleTimeout:       Timeout,
		},
	}

	for _, opt := range opts {
		opt(config)
	}

	config.engine = config.buildEngine()
	config.Handler = config.engine

	return config
}

func (c *Config) Run() error {
	var err error

	if c.TLSConfig != nil {
		err = c.ListenAndServeTLS("", "")
	} else {
		err = c.ListenAndServe()
	}

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (c *Config) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	err := c.Shutdown(ctx)
	if err != nil {
		_ = c.Close()
	}

	return nil
}

func Middlewares(middleware ...gin.HandlerFunc) Option {
	return func(c *Config) {
		c.middleware = append(c.middleware, middleware...)
	}
}

func UseRouter(r Router) Option {
	return func(c *Config) {
		c.routerFunc = r
	}
}

func Injector(i do.Injector) Option {
	return func(c *Config) {
		c.injector = i
	}
}

func (c *Config) buildEngine() *gin.Engine {
	switch c.mode {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "development":
		gin.SetMode(gin.DebugMode)
	default:
		panic("unsupported mode: " + c.mode)
	}

	g := gin.Default()
	g.Use(c.middleware...)
	c.routerFunc(g, c.injector)

	return g
}
