package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	comusic "github.com/sabigara/comusicAPI"
)

type TakeHandler struct {
	comusic.TakeUsecase
	comusic.FileRepository
}

func NewTakeHandler(tu comusic.TakeUsecase, fr comusic.FileRepository) *TakeHandler {
	return &TakeHandler{TakeUsecase: tu, FileRepository: fr}
}

type TakeCreateData struct {
	Name string
}

func (h *TakeHandler) create(c echo.Context) error {
	trackID := c.QueryParam("track_id")
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	file := form.File["file"]
	if len(file) != 1 {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	name := c.FormValue("name")
	src, err := file[0].Open()
	if err != nil {
		return err
	}
	defer src.Close()
	take, uploadedFile, err := h.TakeUsecase.Create(trackID, name, src)
	if err != nil {
		return err
	}
	ret := map[string]interface{}{
		"take": take,
		"file": uploadedFile,
	}
	return c.JSON(http.StatusCreated, ret)
}

func (h *TakeHandler) delete(c echo.Context) error {
	err := h.TakeUsecase.Delete(c.Param("id"))
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
