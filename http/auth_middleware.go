package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	comusic "github.com/sabigara/comusicAPI"
)

type (
	// authMiddlewareConfig defines the config for Session middleware.
	authMiddlewareConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper     middleware.Skipper
		AuthUsecase comusic.AuthUsecase
	}
)

var (
	// DefaultConfig is the default Session middleware config.
	authMiddlewareDefaultConfig = authMiddlewareConfig{
		Skipper: middleware.DefaultSkipper,
	}
)

func authMiddlewareWithConfig(config authMiddlewareConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = authMiddlewareDefaultConfig.Skipper
	}
	if config.AuthUsecase == nil {
		panic("Authenticate function is required for authMiddleware")
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			idToken := c.Request().Header.Get("Authorization")
			user, err := config.AuthUsecase.Authenticate(idToken)
			if err != nil {
				if errors.Is(err, comusic.ErrUnauthenticated) {
					return echo.NewHTTPError(http.StatusUnauthorized).SetInternal(err)
				}
				return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
			}
			c.Set("user", user)
			return next(c)
		}
	}
}
