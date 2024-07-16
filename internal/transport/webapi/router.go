package webapi

import (
	"context"
	"log/slog"
	"net"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
	"github.com/vitaliy-ukiru/test-bank/internal/config"
	"github.com/vitaliy-ukiru/test-bank/internal/transport/webapi/controllers"
	"github.com/vitaliy-ukiru/test-bank/internal/transport/webapi/middlewares"
	"github.com/vitaliy-ukiru/test-bank/pkg/logging"
)

type ApiRouter struct {
	e                 *echo.Echo
	cfg               config.Config
	accountController *controllers.AccountController
}

func configureEcho(e *echo.Echo, cfg config.Config, logger logging.Logger) {
	e.HideBanner = true
	stdLog := logging.ConfigureLogLogger(logger, slog.LevelInfo)
	e.Logger.SetLevel(99)
	e.Logger.SetOutput(stdLog.Writer())
	e.Server.ErrorLog = stdLog

	level := slog.LevelError
	if cfg.Env == config.EnvDev {
		level = slog.LevelDebug
	}

	e.Use(slogecho.NewWithConfig(logger.ToStd(), slogecho.Config{
		DefaultLevel:    level,
		WithRequestBody: true,
	}))

	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.RequestID())
	e.Use(middlewares.WrapRequestContextWithLogger(logger))

}

func New(
	cfg config.Config,
	controller *controllers.AccountController,
	logger logging.Logger,
) *ApiRouter {
	e := echo.New()
	configureEcho(e, cfg, logger)
	controller.Bind(e)

	return &ApiRouter{
		e:                 e,
		cfg:               cfg,
		accountController: controller,
	}
}

func (a *ApiRouter) Start() error {
	return a.e.Start(net.JoinHostPort(a.cfg.Server.Host, strconv.Itoa(a.cfg.Server.Port)))
}

func (a *ApiRouter) Shutdown(ctx context.Context) error {
	return a.e.Shutdown(ctx)
}
