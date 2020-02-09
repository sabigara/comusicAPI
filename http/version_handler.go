package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	comusic "github.com/sabigara/comusicAPI"
)

type VersionHandler struct {
	comusic.VersionUsecase
}

func NewVersionHandler(versionUsecase comusic.VersionUsecase) *VersionHandler {
	return &VersionHandler{VersionUsecase: versionUsecase}
}

type VersionCreateData struct {
	Name string
}

func (h *VersionHandler) create(c echo.Context) error {
	songID := c.QueryParam("song_id")
	req := &VersionCreateData{}
	c.Bind(req)
	ver, err := h.VersionUsecase.Create(songID, req.Name)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, ver)
}
