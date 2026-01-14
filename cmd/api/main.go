package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Any("url", c.Request().URL),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Any("url", c.Request().URL),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "echo: GET /\n")
	})

	e.Any("/redirect", func(c echo.Context) error {
		target := c.QueryParam("target")
		if target == "" {
			target = "/"
		}
		return c.Redirect(http.StatusFound, target)
	})

	addr := ":5001"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	e.Logger.Infof("starting api server on %s", addr)
	e.Logger.Fatal(e.Start(addr))
}
