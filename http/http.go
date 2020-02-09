package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	comusic "github.com/sabigara/comusicAPI"
)

var profHandler *ProfileHandler
var studioHandler *StudioHandler
var songHandler *SongHandler
var verHandler *VersionHandler
var trackHandler *TrackHandler
var takeHandler *TakeHandler

var authenticate func(...interface{}) (*comusic.User, error)

// SetHandlers sets all handlers with their all dependencies injected.
func SetHandlers(
	prof *ProfileHandler,
	studio *StudioHandler,
	song *SongHandler,
	ver *VersionHandler,
	track *TrackHandler,
	take *TakeHandler,
) {
	profHandler = prof
	studioHandler = studio
	songHandler = song
	verHandler = ver
	trackHandler = track
	takeHandler = take
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

	e.GET("profile", profHandler.get)
	e.POST("profile", profHandler.create)
	e.PATCH("profile", profHandler.update)

	e.GET("studios", studioHandler.get)
	e.POST("studios", studioHandler.create)

	e.POST("songs", songHandler.create)
	e.GET("songs", songHandler.get)

	e.POST("versions", verHandler.create)

	e.POST("tracks", trackHandler.create)
	e.GET("tracks", trackHandler.get)

	e.POST("takes", takeHandler.create)

	e.Logger.Fatal(e.Start(addr))
}
