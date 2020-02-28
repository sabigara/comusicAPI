package interactor

import (
	"fmt"

	comusic "github.com/sabigara/comusicAPI"
)

type StudioUsecase struct {
	comusic.StudioRepository
}

func NewStudioUsecase(pr comusic.StudioRepository) *StudioUsecase {
	return &StudioUsecase{StudioRepository: pr}
}

func (su *StudioUsecase) Create(ownerID, nickname string) (*comusic.Studio, error) {
	studio := comusic.NewStudio(ownerID, nickname)
	err := su.StudioRepository.Create(studio)
	if err != nil {
		return nil, fmt.Errorf("interactor.studio_usecase.Create: %w", err)
	}
	return studio, nil
}

func (su *StudioUsecase) FilterByOwnerID(ownerID string) (*[]comusic.Studio, error) {
	user, err := su.StudioRepository.FilterByOwnerID(ownerID)
	if err != nil {
		return nil, fmt.Errorf("interactor.studio_usecase.Get: %w", err)
	}
	return user, err
}

func (su *StudioUsecase) GetContents(studioID string) ([]*comusic.Song, []*comusic.Version, error) {
	songs, vers, err := su.StudioRepository.GetContents(studioID)
	if err != nil {
		return nil, nil, fmt.Errorf("interactor.studio_usecase.GetContents: %w", err)
	}
	return songs, vers, nil
}
