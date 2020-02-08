package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	comusic "github.com/sabigara/comusicAPI"
)

type ProfileHandler struct {
	comusic.ProfileUsecase
}

func NewProfileHandler(userUsecase comusic.ProfileUsecase) *ProfileHandler {
	return &ProfileHandler{ProfileUsecase: userUsecase}
}

type ProfileCreateData struct {
	Nickname string
	Bio      string
}

func (h *ProfileHandler) create(c echo.Context) error {
	user := c.Get("user").(*comusic.User)
	req := &ProfileCreateData{}
	c.Bind(req)
	profile, err := h.ProfileUsecase.Create(user.ID, req.Nickname, req.Bio)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, profile)
}

type ProfileUpdateData struct {
	Nickname *string
	Bio      *string
}

func (h *ProfileHandler) update(c echo.Context) error {
	user := c.Get("user").(*comusic.User)
	req := &ProfileUpdateData{}
	if err := c.Bind(req); err != nil {
		return err
	}
	err := h.ProfileUsecase.Update(user.ID, req.Nickname, req.Bio)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (h *ProfileHandler) get(c echo.Context) error {
	user := c.Get("user").(*comusic.User)
	profile, err := h.GetByUserID(user.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, profile)
}
