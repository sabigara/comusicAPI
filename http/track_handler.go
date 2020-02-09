package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	comusic "github.com/sabigara/comusicAPI"
)

type TrackHandler struct {
	comusic.TrackUsecase
}

func NewTrackHandler(trackUsecase comusic.TrackUsecase) *TrackHandler {
	return &TrackHandler{TrackUsecase: trackUsecase}
}

type TrackCreateData struct {
	Name string
}

func (h *TrackHandler) create(c echo.Context) error {
	studioID := c.QueryParam("version_id")
	req := &TrackCreateData{}
	c.Bind(req)
	profile, err := h.TrackUsecase.Create(studioID, req.Name)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, profile)
}

func (h *TrackHandler) get(c echo.Context) error {
	trackVersMap, err := h.FilterByVersionIDWithTakes(c.QueryParam("version_id"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, trackVersMap)
}
