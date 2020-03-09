package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	comusic "github.com/sabigara/comusicAPI"
)

type Hooks struct {
	comusic.ProfileUsecase
	comusic.StudioUsecase
}

func NewHooks(profileUsecase comusic.ProfileUsecase, studioUsecase comusic.StudioUsecase) *Hooks {
	return &Hooks{ProfileUsecase: profileUsecase, StudioUsecase: studioUsecase}
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
	studio, err := h.StudioUsecase.Create(req.UserID, "Your Studio")
	if err != nil {
		return err
	}
	err = h.StudioUsecase.AddMembers(studio.ID, req.UserID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, profile)
}
