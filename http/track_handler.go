package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	comusic "github.com/sabigara/comusicAPI"
)

type TrackHandler struct {
	comusic.TrackUsecase
	comusic.FileRepository
}

func NewTrackHandler(tu comusic.TrackUsecase, fr comusic.FileRepository) *TrackHandler {
	return &TrackHandler{TrackUsecase: tu, FileRepository: fr}
}

type TrackCreateData struct {
	Name string
}

func (h *TrackHandler) create(c echo.Context) error {
	studioID := c.QueryParam("version_id")
	req := &TrackCreateData{}
	c.Bind(req)
	profile, err := h.TrackUsecase.Create(studioID, req.Name)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, profile)
}

type getReturn struct {
	Tracks *Entity `json:"tracks"`
	Takes  *Entity `json:"takes"`
	Files  *Entity `json:"files"`
}

func (h *TrackHandler) get(c echo.Context) error {
	trackVersMap, err := h.FilterByVersionIDWithTakes(c.QueryParam("version_id"))
	if err != nil {
		return err
	}
	ret := &getReturn{}
	ret.Tracks = NewEntity()
	ret.Takes = NewEntity()
	ret.Files = NewEntity()

	for k, v := range trackVersMap {
		ret.Tracks.ByID[k] = v.Data
		ret.Tracks.AllIDs = append(ret.Tracks.AllIDs, k)
		for _, tk := range v.Takes {
			ret.Takes.ByID[tk.ID] = tk
			ret.Takes.AllIDs = append(ret.Takes.AllIDs, tk.ID)
			ret.Files.ByID[tk.FileID] = &comusic.File{
				URL: h.FileRepository.URL(tk.FileID),
			}
			ret.Files.AllIDs = append(ret.Files.AllIDs, tk.FileID)
		}
	}

	return c.JSON(http.StatusOK, ret)
}
