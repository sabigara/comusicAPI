package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	comusic "github.com/sabigara/comusicAPI"
)

type Hooks struct {
	comusic.ProfileUsecase
	comusic.StudioUsecase
	comusic.PubSub
}

func NewHooks(
	profileUsecase comusic.ProfileUsecase,
	studioUsecase comusic.StudioUsecase,
	pubsub comusic.PubSub,
) *Hooks {
	return &Hooks{
		ProfileUsecase: profileUsecase,
		StudioUsecase:  studioUsecase,
		PubSub:         pubsub,
	}
}

type NewUserInput struct {
	UserID   string
	Nickname string
}

func (h *Hooks) newUserCreated(c echo.Context) error {
	user := c.Get("user").(*comusic.User)
	_, err := h.ProfileUsecase.Create(user.ID, "", "")
	if err != nil {
		return err
	}
	studio, err := h.StudioUsecase.Create(user.ID, "Your Studio")
	if err != nil {
		return err
	}
	err = h.StudioUsecase.AddMembers(studio.ID, user.ID)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
