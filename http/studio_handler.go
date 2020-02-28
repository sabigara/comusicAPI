package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	comusic "github.com/sabigara/comusicAPI"
)

type StudioHandler struct {
	comusic.StudioUsecase
}

func NewStudioHandler(studioUsecase comusic.StudioUsecase) *StudioHandler {
	return &StudioHandler{StudioUsecase: studioUsecase}
}

type StudioCreateData struct {
	Name string
}

func (h *StudioHandler) create(c echo.Context) error {
	user := c.Get("user").(*comusic.User)
	req := &StudioCreateData{}
	c.Bind(req)
	profile, err := h.StudioUsecase.Create(user.ID, req.Name)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, profile)
}

func (h *StudioHandler) get(c echo.Context) error {
	user := c.Get("user").(*comusic.User)
	studios, err := h.FilterByOwnerID(user.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, studios)
}

type StudioContents struct {
	Songs   *RespEntity `json:"songs"`
	Version *RespEntity `json:"versions"`
}

func (h *StudioHandler) getContents(c echo.Context) error {
	songs, vers, err := h.StudioUsecase.GetContents(c.Param("id"))
	if err != nil {
		return err
	}
	ret := &StudioContents{
		Songs:   NewRespEntity(songs),
		Version: NewRespEntity(vers),
	}
	return c.JSON(http.StatusOK, ret)
}
