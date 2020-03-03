package http

import (
	"net/http"
	"os"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	comusic "github.com/sabigara/comusicAPI"
	"github.com/sabigara/comusicAPI/utils"
)

type strKeyMap map[string]interface{}

// RespEntity represents json response for a query for
// an object type, keeping order in AllIds while
// having capability to select one ByID.
type RespEntity struct {
	ByID   strKeyMap `json:"byId"`
	AllIDs []string  `json:"allIds"`
}

func NewRespEntity(data interface{}) *RespEntity {
	ret := &RespEntity{}
	ret.ByID = strKeyMap{}
	ret.AllIDs = []string{}

	s := reflect.ValueOf(data)
	for i := 0; i < s.Len(); i++ {
		d := s.Index(i).Elem()
		// Should do type guard and nil check to ensure having "ID" field.
		id := d.FieldByName("ID").String()
		if !utils.Contains(ret.AllIDs, id) {
			ret.AllIDs = append(ret.AllIDs, id)
		}
		ret.ByID[id] = d.Interface()
	}
	return ret
}

var profHandler *ProfileHandler
var studioHandler *StudioHandler
var songHandler *SongHandler
var verHandler *VersionHandler
var trackHandler *TrackHandler
var takeHandler *TakeHandler
var hooks *Hooks

var authenticate func(...interface{}) (*comusic.User, error)

// SetHandlers sets all handlers with their all dependencies injected.
func SetHandlers(
	prof *ProfileHandler,
	studio *StudioHandler,
	song *SongHandler,
	ver *VersionHandler,
	track *TrackHandler,
	take *TakeHandler,
	h *Hooks,
) {
	profHandler = prof
	studioHandler = studio
	songHandler = song
	verHandler = ver
	trackHandler = track
	takeHandler = take
	hooks = h
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
		uploadsDir, ok := os.LookupEnv("UPLOADS_DIR")
		if !ok {
			panic("UPLOADS_DIR not specified")
		}
		e.Static("uploads", uploadsDir)
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
	e.GET("studios/:id/contents", studioHandler.getContents)

	e.POST("songs", songHandler.create)
	e.DELETE("songs/:id", songHandler.delete)

	e.POST("versions", verHandler.create)
	e.GET("versions/:id/contents", verHandler.get)
	e.DELETE("versions/:id", verHandler.delete)

	e.POST("tracks", trackHandler.create)
	e.DELETE("tracks/:id", trackHandler.delete)

	e.POST("takes", takeHandler.create)
	e.DELETE("takes/:id", takeHandler.delete)

	// Hooks for handling events.
	// TODO: Disable authMiddleware
	e.POST("hooks/new-user", hooks.newUserCreated)

	e.Logger.Fatal(e.Start(addr))
}
