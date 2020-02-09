package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	comusic "github.com/sabigara/comusicAPI"
)

type TakeHandler struct {
	comusic.TakeUsecase
}

func NewTakeHandler(takeUsecase comusic.TakeUsecase) *TakeHandler {
	return &TakeHandler{TakeUsecase: takeUsecase}
}

type TakeCreateData struct {
	Name string
}

func (h *TakeHandler) create(c echo.Context) error {
	trackID := c.QueryParam("track_id")
	req := &TakeCreateData{}
	c.Bind(req)
	ver, err := h.TakeUsecase.Create(trackID, req.Name)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, ver)
}
