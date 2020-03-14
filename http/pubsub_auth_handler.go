package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	comusic "github.com/sabigara/comusicAPI"
)

type PubSubAuthHandler struct {
	comusic.PubSub
}

func NewPubSubAuthHandler(pubsub comusic.PubSub) *PubSubAuthHandler {
	return &PubSubAuthHandler{PubSub: pubsub}
}

type PubSubTokenResp struct {
	PubSubToken string `json:"pubsubToken"`
}

func (h *PubSubAuthHandler) get(c echo.Context) error {
	user := c.Get("user").(*comusic.User)
	token, err := h.PubSub.GenerateToken(user)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &PubSubTokenResp{PubSubToken: token})
}
