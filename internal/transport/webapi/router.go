package webapi

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
	"github.com/swaggest/swgui"
	"github.com/swaggest/swgui/v5emb"
	"github.com/vitaliy-ukiru/bank-service/api"
	"github.com/vitaliy-ukiru/bank-service/internal/config"
	"github.com/vitaliy-ukiru/bank-service/internal/transport/webapi/controllers"
	"github.com/vitaliy-ukiru/bank-service/internal/transport/webapi/middlewares"
	"github.com/vitaliy-ukiru/bank-service/internal/transport/webapi/response"
	"github.com/vitaliy-ukiru/bank-service/pkg/logging"
)

type ApiRouter struct {
	e                 *echo.Echo
	cfg               config.Config
	accountController *controllers.AccountController
}

func configureEcho(e *echo.Echo, logger logging.Logger) {
	e.HideBanner = true
	stdLog := logging.ConfigureLogLogger(logger, slog.LevelInfo)
	e.Logger.SetLevel(99)
	e.Logger.SetOutput(stdLog.Writer())
	e.Server.ErrorLog = stdLog

	e.Use(slogecho.NewWithConfig(logger.ToStd(), slogecho.Config{
		DefaultLevel:     slog.LevelDebug,
		ClientErrorLevel: slog.LevelDebug,
		ServerErrorLevel: slog.LevelWarn,
		WithRequestBody:  true,
	}))

	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.RequestID())
	e.Use(middlewares.WrapRequestContextWithLogger(logger))
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		var httpErr *echo.HTTPError
		if errors.As(err, &httpErr) {
			_ = c.JSON(httpErr.Code, response.Fail(fmt.Sprint(httpErr.Message)))
			return
		}
		_ = c.JSON(http.StatusInternalServerError, response.Error(err))
	}

	echo.MethodNotAllowedHandler = plainErrorHandler(http.StatusMethodNotAllowed)
	echo.NotFoundHandler = plainErrorHandler(http.StatusNotFound)
}

func plainErrorHandler(code int) func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.JSON(code, response.Fail(http.StatusText(code)))
	}
}

func New(
	cfg config.Config,
	controller *controllers.AccountController,
	logger logging.Logger,
) *ApiRouter {
	e := echo.New()
	configureEcho(e, logger)

	controller.Bind(e)

	e.Any("/docs*", echo.WrapHandler(
		v5emb.NewHandlerWithConfig(swgui.Config{
			Title:       "Bank Service",
			SwaggerJSON: "/openapi.json",
			BasePath:    "/docs",
			JsonEditor:  true,
		}),
	))

	e.GET("/openapi.json", func(c echo.Context) error {
		return c.JSONBlob(200, api.OpenAPISpec)
	})

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
