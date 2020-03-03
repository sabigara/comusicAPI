package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	comusic "github.com/sabigara/comusicAPI"
)

type Hooks struct {
	comusic.ProfileUsecase
}

func NewHooks(profileUsecase comusic.ProfileUsecase) *Hooks {
	return &Hooks{ProfileUsecase: profileUsecase}
}

type NewUserInput struct {
	UserID   string
	Nickname string
}

func (h *Hooks) newUserCreated(c echo.Context) error {
	req := &NewUserInput{}
	c.Bind(req)
	profile, err := h.ProfileUsecase.Create(req.UserID, req.Nickname, "")
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, profile)
}
