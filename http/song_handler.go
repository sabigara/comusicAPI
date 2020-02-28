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

func (h *SongHandler) delete(c echo.Context) error {
	err := h.SongUsecase.Delete(c.Param("id"))
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
