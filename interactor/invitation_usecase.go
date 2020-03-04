package interactor

import (
	"fmt"

	comusic "github.com/sabigara/comusicAPI"
)

type InvitationUsecase struct {
	comusic.InvitationRepository
	comusic.MailUsecase
}

func NewInvitationUsecase(inviteRepo comusic.InvitationRepository, mailUsecase comusic.MailUsecase) *InvitationUsecase {
	return &InvitationUsecase{InvitationRepository: inviteRepo, MailUsecase: mailUsecase}
}

func (u *InvitationUsecase) Filter(email, groupID string) ([]*comusic.Invitation, error) {
	invites, err := u.InvitationRepository.Filter(email, groupID)
	if err != nil {
		return nil, fmt.Errorf("interactor.invitation_usecase.Filter: %w", err)
	}
	return invites, err
}

func (u *InvitationUsecase) Create(email, groupID string, groupType comusic.GroupType) error {
	err := u.InvitationRepository.Create(email, groupID, groupType)
	if err != nil {
		return fmt.Errorf("interactor.invitation_usecase.Create: %w", err)
	}
	u.MailUsecase.InviteToStudioNew(email, "studio_name", "http://localhost:3000/login")
	return nil
}

func (u *InvitationUsecase) Accept(email, groupID string) error {
	err := u.InvitationRepository.Accept(email, groupID)
	if err != nil {
		return fmt.Errorf("interactor.invitation_usecase.Update: %w", err)
	}
	return nil
}
