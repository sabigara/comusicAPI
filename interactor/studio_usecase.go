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
		return nil, err
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

func (u *StudioUsecase) Filter(ownerID, memberID string) (studios []*comusic.Studio, err error) {
	switch {
	case ownerID != "":
		studios, err = u.StudioRepository.FilterByOwnerID(ownerID)
	case memberID != "":
		studios, err = u.StudioRepository.FilterByMemberID(memberID)
	default:
		return nil, fmt.Errorf("interactor.studio_usecase.Filter: no query provided.")
	}
	if err != nil {
		return nil, err
	}
	return
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
