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

type HTTP struct {
	comusic.AuthUsecase
	*ProfileHandler
	*StudioHandler
	*SongHandler
	*VersionHandler
	*TrackHandler
	*TakeHandler
	*InvitationHandler
	*PubSubAuthHandler
	*Hooks
}

func New(
	auth comusic.AuthUsecase,
	prof *ProfileHandler,
	studio *StudioHandler,
	song *SongHandler,
	ver *VersionHandler,
	track *TrackHandler,
	take *TakeHandler,
	invite *InvitationHandler,
	pubsubAuth *PubSubAuthHandler,
	h *Hooks,
) *HTTP {
	return &HTTP{
		AuthUsecase:       auth,
		ProfileHandler:    prof,
		StudioHandler:     studio,
		SongHandler:       song,
		VersionHandler:    ver,
		TrackHandler:      track,
		TakeHandler:       take,
		InvitationHandler: invite,
		PubSubAuthHandler: pubsubAuth,
		Hooks:             h,
	}
}

// Start starts server after settings routes.
func (httpIns *HTTP) Start(addr string, debug bool) {
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
	e.Use(authMiddlewareWithConfig(
		authMiddlewareConfig{AuthUsecase: httpIns.AuthUsecase},
	))

	e.GET("profile", httpIns.ProfileHandler.get)
	e.POST("profile", httpIns.ProfileHandler.create)
	e.PATCH("profile", httpIns.ProfileHandler.update)

	e.GET("studios", httpIns.StudioHandler.filter)
	e.POST("studios", httpIns.StudioHandler.create)
	e.GET("studios/:id/contents", httpIns.StudioHandler.getContents)
	e.GET("studios/:id/members", httpIns.StudioHandler.getMembers)

	e.GET("invitations", httpIns.InvitationHandler.filter)
	e.PUT("invitations", httpIns.InvitationHandler.create)
	e.PATCH("invitations", httpIns.InvitationHandler.accept)

	e.GET("songs", httpIns.SongHandler.filter)
	e.POST("songs", httpIns.SongHandler.create)
	e.DELETE("songs/:id", httpIns.SongHandler.delete)
	e.GET("songs/:id/guests", nil)

	e.POST("versions", httpIns.VersionHandler.create)
	e.GET("versions/:id/contents", httpIns.VersionHandler.get)
	e.DELETE("versions/:id", httpIns.VersionHandler.delete)

	e.POST("tracks", httpIns.TrackHandler.create)
	e.DELETE("tracks/:id", httpIns.TrackHandler.delete)

	e.POST("takes", httpIns.TakeHandler.create)
	e.DELETE("takes/:id", httpIns.TakeHandler.delete)

	e.GET("pubsub/token", httpIns.PubSubAuthHandler.get)

	// Hooks for handling events.
	// TODO: Disable authMiddleware
	e.POST("hooks/new-user", httpIns.Hooks.newUserCreated)

	e.Logger.Fatal(e.Start(addr))
}
