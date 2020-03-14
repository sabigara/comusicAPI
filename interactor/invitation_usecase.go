package interactor

import (
	"fmt"

	comusic "github.com/sabigara/comusicAPI"
)

type InvitationUsecase struct {
	comusic.InvitationRepository
	comusic.UserRepository
	comusic.StudioUsecase
	comusic.SongUsecase
	comusic.MailUsecase
	comusic.PubSub
}

func NewInvitationUsecase(
	inviteRepo comusic.InvitationRepository,
	userRepo comusic.UserRepository,
	studioUsecase comusic.StudioUsecase,
	songUsecase comusic.SongUsecase,
	mailUsecase comusic.MailUsecase,
	pubsub comusic.PubSub,
) *InvitationUsecase {
	return &InvitationUsecase{
		InvitationRepository: inviteRepo,
		UserRepository:       userRepo,
		StudioUsecase:        studioUsecase,
		SongUsecase:          songUsecase,
		MailUsecase:          mailUsecase,
		PubSub:               pubsub,
	}
}

func (u *InvitationUsecase) Filter(email, groupID string) ([]*comusic.Invitation, error) {
	invites, err := u.InvitationRepository.Filter(email, groupID)
	if err != nil {
		return nil, err
	}
	return invites, err
}

func (u *InvitationUsecase) Create(email, groupID string, groupType comusic.GroupType) (err error) {
	// Check if provided groupID exists as studioID or songID.
	// If not, return error.
	switch groupType {
	case comusic.StudioGroupType:
		_, err = u.StudioUsecase.GetByID(groupID)
	case comusic.SongGroupType:
		_, err = u.SongUsecase.GetByID(groupID)
	default:
		return fmt.Errorf("interactor.invitation_usecase.Create: Invalid GroupType")
	}
	if err != nil {
		return err
	}

	invitation := comusic.NewInvitation(email, groupID, groupType)

	err = u.InvitationRepository.Create(invitation)
	if err != nil {
		return err
	}
	// Check if the user already exists.
	user, err := u.UserRepository.GetByEmail(email)
	if err != nil {
		err = u.MailUsecase.InviteToStudioNew(email, "studio_name", "http://localhost:3000/login")
	} else {
		err = u.MailUsecase.InviteToStudio(user, "studio_name", "http://localhost:3000/login")
		if err != nil {
			return err
		}
		err = u.PubSub.Publish("invitation#"+user.ID, comusic.NewPublication("invitation", invitation))
		if err != nil {
			return err
		}
	}
	return err
}

func (u *InvitationUsecase) Accept(email, groupID string) error {
	invite, err := u.InvitationRepository.GetByIDs(email, groupID)
	if err != nil {
		return err
	}
	err = u.InvitationRepository.Accept(email, groupID)
	if err != nil {
		return err
	}
	user, err := u.UserRepository.GetByEmail(email)
	if err != nil {
		return err
	}
	switch invite.GroupType {
	case comusic.StudioGroupType:
		return u.StudioUsecase.AddMembers(groupID, user.ID)
	case comusic.SongGroupType:
		return u.SongUsecase.AddGuests(groupID, user.ID)
	default:
		return nil
	}
}
