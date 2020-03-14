package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	comusic "github.com/sabigara/comusicAPI"
)

type StudioHandler struct {
	comusic.StudioUsecase
	comusic.ProfileUsecase
}

func NewStudioHandler(
	studioUsecase comusic.StudioUsecase,
	profileUsecase comusic.ProfileUsecase,
) *StudioHandler {
	return &StudioHandler{
		StudioUsecase:  studioUsecase,
		ProfileUsecase: profileUsecase,
	}
}

type StudioCreateData struct {
	Name string
}

func (h *StudioHandler) create(c echo.Context) error {
	user := c.Get("user").(*comusic.User)
	req := &StudioCreateData{}
	c.Bind(req)
	profile, err := h.StudioUsecase.Create(user.ID, req.Name)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, profile)
}

type StudioFilterResp struct {
	Studios *RespEntity `json:"studios"`
}

func (h *StudioHandler) filter(c echo.Context) error {
	ownerID := c.QueryParam("owner_id")
	memberID := c.QueryParam("member_id")
	studios, err := h.Filter(ownerID, memberID)
	if err != nil {
		return err
	}
	resp := &StudioFilterResp{Studios: NewRespEntity(studios)}
	return c.JSON(http.StatusOK, resp)
}

type StudioContents struct {
	Songs   *RespEntity `json:"songs"`
	Version *RespEntity `json:"versions"`
}

// Plus studio members and song guests?
func (h *StudioHandler) getContents(c echo.Context) error {
	songs, vers, err := h.StudioUsecase.GetContents(c.Param("id"))
	if err != nil {
		return err
	}
	ret := &StudioContents{
		Songs:   NewRespEntity(songs),
		Version: NewRespEntity(vers),
	}
	return c.JSON(http.StatusOK, ret)
}

type GetStudioMembersResp struct {
	Members *RespEntity `json:"members"`
}

func (h *StudioHandler) getMembers(c echo.Context) error {
	id := c.Param("id")
	members, err := h.ProfileUsecase.GetStudioMembers(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &GetStudioMembersResp{Members: NewRespEntity(members)})
}
