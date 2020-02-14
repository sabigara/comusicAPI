package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	comusic "github.com/sabigara/comusicAPI"
	"github.com/sabigara/comusicAPI/utils"
)

type VersionHandler struct {
	comusic.VersionUsecase
	comusic.FileRepository
}

func NewVersionHandler(vu comusic.VersionUsecase, fr comusic.FileRepository) *VersionHandler {
	return &VersionHandler{VersionUsecase: vu, FileRepository: fr}
}

type VersionCreateData struct {
	Name string
}

func (h *VersionHandler) create(c echo.Context) error {
	songID := c.QueryParam("song_id")
	req := &VersionCreateData{}
	c.Bind(req)
	ver, err := h.VersionUsecase.Create(songID, req.Name)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, ver)
}

type VersionContents struct {
	Tracks *Entity `json:"tracks"`
	Takes  *Entity `json:"takes"`
	Files  *Entity `json:"files"`
}

func (h *VersionHandler) get(c echo.Context) error {
	tracks, takes, err := h.VersionUsecase.GetContents(c.Param("id"))
	if err != nil {
		return err
	}
	ret := &VersionContents{}
	ret.Tracks = NewEntity()
	ret.Takes = NewEntity()
	ret.Files = NewEntity()

	for _, tr := range tracks {
		if !utils.Contains(ret.Tracks.AllIDs, tr.ID) {
			ret.Tracks.AllIDs = append(ret.Tracks.AllIDs, tr.ID)
		}
		ret.Tracks.ByID[tr.ID] = tr
	}
	for _, tk := range takes {
		if !utils.Contains(ret.Takes.AllIDs, tk.ID) {
			ret.Takes.AllIDs = append(ret.Takes.AllIDs, tk.ID)
		}
		if !utils.Contains(ret.Files.AllIDs, tk.FileID) {
			ret.Files.AllIDs = append(ret.Files.AllIDs, tk.FileID)
		}
		ret.Takes.ByID[tk.ID] = tk
		ret.Files.ByID[tk.FileID] = &comusic.File{URL: h.FileRepository.URL(tk.FileID)}
	}

	return c.JSON(http.StatusOK, ret)
}
