package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	comusic "github.com/sabigara/comusicAPI"
)

type VersionHandler struct {
	comusic.VersionUsecase
}

func NewVersionHandler(vu comusic.VersionUsecase) *VersionHandler {
	return &VersionHandler{VersionUsecase: vu}
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

func (h *VersionHandler) delete(c echo.Context) error {
	err := h.VersionUsecase.Delete(c.Param("id"))
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

type VersionContents struct {
	Tracks *RespEntity `json:"tracks"`
	Takes  *RespEntity `json:"takes"`
	Files  *RespEntity `json:"files"`
}

func (h *VersionHandler) get(c echo.Context) error {
	tracks, takes, files, err := h.VersionUsecase.GetContents(c.Param("id"))
	if err != nil {
		return err
	}
	ret := &VersionContents{
		Tracks: NewRespEntity(tracks),
		Takes:  NewRespEntity(takes),
		Files:  NewRespEntity(files),
	}
	return c.JSON(http.StatusOK, ret)
}
