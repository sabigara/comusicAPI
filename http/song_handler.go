package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	comusic "github.com/sabigara/comusicAPI"
)

type SongHandler struct {
	comusic.SongUsecase
}

func NewSongHandler(songUsecase comusic.SongUsecase) *SongHandler {
	return &SongHandler{SongUsecase: songUsecase}
}

type SongCreateData struct {
	Name string
}

func (h *SongHandler) create(c echo.Context) error {
	studioID := c.QueryParam("studio_id")
	req := &SongCreateData{}
	c.Bind(req)
	profile, err := h.SongUsecase.Create(studioID, req.Name)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, profile)
}

func (h *SongHandler) get(c echo.Context) error {
	songVersMap, err := h.FilterByStudioIDWithVersions(c.QueryParam("studio_id"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, songVersMap)
}
