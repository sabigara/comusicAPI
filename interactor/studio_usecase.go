package interactor

import (
	"fmt"

	comusic "github.com/sabigara/comusicAPI"
)

type StudioUsecase struct {
	comusic.StudioRepository
}

func NewStudioUsecase(studioRepo comusic.StudioRepository) *StudioUsecase {
	return &StudioUsecase{StudioRepository: studioRepo}
}
func (u *StudioUsecase) GetByID(id string) (*comusic.Studio, error) {
	studio, err := u.StudioRepository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("interactor.studio_usecase.GetByID: %w", err)
	}
	return studio, err
}

func (u *StudioUsecase) Create(ownerID, nickname string) (*comusic.Studio, error) {
	studio := comusic.NewStudio(ownerID, nickname)
	err := u.StudioRepository.Create(studio)
	if err != nil {
		return nil, fmt.Errorf("interactor.studio_usecase.Create: %w", err)
	}
	return studio, nil
}

func (u *StudioUsecase) FilterByOwnerID(ownerID string) (*[]comusic.Studio, error) {
	user, err := u.StudioRepository.FilterByOwnerID(ownerID)
	if err != nil {
		return nil, fmt.Errorf("interactor.studio_usecase.Get: %w", err)
	}
	return user, err
}

func (u *StudioUsecase) GetContents(studioID string) ([]*comusic.Song, []*comusic.Version, error) {
	songs, vers, err := u.StudioRepository.GetContents(studioID)
	if err != nil {
		return nil, nil, fmt.Errorf("interactor.studio_usecase.GetContents: %w", err)
	}
	return songs, vers, nil
}

func (u *StudioUsecase) AddMembers(studioID string, userID ...string) error {
	err := u.StudioRepository.AddMembers(studioID, userID...)
	if err != nil {
		return fmt.Errorf("interactor.studio_usecase.AddMember: %w", err)
	}
	return nil
}
