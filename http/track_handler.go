package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	comusic "github.com/sabigara/comusicAPI"
)

type TrackHandler struct {
	comusic.TrackUsecase
}

func NewTrackHandler(tu comusic.TrackUsecase) *TrackHandler {
	return &TrackHandler{TrackUsecase: tu}
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
