package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	comusic "github.com/sabigara/comusicAPI"
)

type InvitationHandler struct {
	comusic.InvitationUsecase
}

func NewInvitationHandler(inviteUsecase comusic.InvitationUsecase) *InvitationHandler {
	return &InvitationHandler{InvitationUsecase: inviteUsecase}
}

type InvitationsResp struct {
	Invitations *RespEntity `json:"invitations"`
}

func (h *InvitationHandler) filter(c echo.Context) error {
	email := c.QueryParam("email")
	groupID := c.QueryParam("group_id")
	invites, err := h.InvitationUsecase.Filter(email, groupID)
	if err != nil {
		return err
	}
	ret := &InvitationsResp{
		Invitations: NewRespEntity(invites),
	}
	return c.JSON(http.StatusOK, ret)
}

func (h *InvitationHandler) create(c echo.Context) error {
	// user := c.Get("user").(*comusic.User)
	groupID := c.QueryParam("group_id")
	groupTypeStr := c.QueryParam("group_type")
	email := c.QueryParam("email")
	groupType := comusic.NewGroupType(groupTypeStr)
	if groupType == comusic.ErrGroupType {
		return echo.NewHTTPError(http.StatusBadRequest, "invitation_handler.create: invalid group_type")
	}
	err := h.InvitationUsecase.Create(email, groupID, groupType)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

type InvitationUpdateInput struct {
	IsAccepted bool `json:"isAccepted"`
}

func (h *InvitationHandler) accept(c echo.Context) error {
	user := c.Get("user").(*comusic.User)
	groupID := c.QueryParam("group_id")
	input := &InvitationUpdateInput{}
	err := c.Bind(input)
	if err != nil {
		return err
	}
	err = h.InvitationUsecase.Accept(user.Email, groupID)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
