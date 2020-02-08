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

// type StudioUpdateData struct {
// 	Nickname *string
// 	Bio      *string
// }

// func (h *StudioHandler) update(c echo.Context) error {
// 	user := c.Get("user").(*comusic.User)
// 	req := &ProfileUpdateData{}
// 	if err := c.Bind(req); err != nil {
// 		return err
// 	}
// 	err := h.StudioUsecase.Update(user.ID, req.Nickname, req.Bio)
// 	if err != nil {
// 		return err
// 	}
// 	return c.NoContent(http.StatusOK)
// }

func (h *StudioHandler) get(c echo.Context) error {
	user := c.Get("user").(*comusic.User)
	studios, err := h.FilterByOwnerID(user.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, studios)
}
