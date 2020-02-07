package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	comusic "github.com/sabigara/comusicAPI"
)

// UserHandler implements HandlerFunc methods for User domain.
type UserHandler struct {
	comusic.UserUsecase
}

// NewUserHandler returns new UserHandler
func NewUserHandler(userUsecase comusic.UserUsecase) *UserHandler {
	return &UserHandler{UserUsecase: userUsecase}
}

func (h *UserHandler) get(c echo.Context) error {
	user := c.Get("user")
	c.JSON(http.StatusOK, user)
	return nil
}
