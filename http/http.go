package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	comusic "github.com/sabigara/comusicAPI"
)

var profileHandler *ProfileHandler
var studioHandler *StudioHandler
var authenticate func(...interface{}) (*comusic.User, error)

// SetHandlers sets all handlers with their all dependencies injected.
func SetHandlers(
	profile *ProfileHandler,
	studio *StudioHandler,
) {
	profileHandler = profile
	studioHandler = studio
}

func SetAuthenticate(f func(...interface{}) (*comusic.User, error)) {
	authenticate = f
}

func errorHandler(e *echo.Echo) func(error, echo.Context) {
	return func(err error, c echo.Context) {
		he, ok := err.(*echo.HTTPError)
		if !ok {
			he = &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
			}
			e.Logger.Error(err.Error())
		}
		code := he.Code
		message := he.Message

		if e.Debug {
			message = err.Error()
		}

		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead {
				err = c.NoContent(he.Code)
			} else {
				err = c.JSON(code, message)
			}
			if err != nil {
				e.Logger.Error(err)
			}
		}
	}
}

// Start starts server after settings routes.
func Start(addr string, debug bool) {
	e := echo.New()
	e.HideBanner = true
	if debug {
		e.Debug = true
	}
	e.HTTPErrorHandler = errorHandler(e)

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(authMiddlewareWithConfig(authMiddlewareConfig{Authenticate: authenticate}))

	e.GET("profile", profileHandler.get)
	e.POST("profile", profileHandler.create)
	e.PATCH("profile", profileHandler.update)

	e.GET("studios", studioHandler.get)
	e.POST("studios", studioHandler.create)

	e.Logger.Fatal(e.Start(addr))
}
