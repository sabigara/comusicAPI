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

func (pu *StudioUsecase) Create(ownerID, nickname string) (*comusic.Studio, error) {
	studio := comusic.NewStudio(ownerID, nickname)
	err := pu.StudioRepository.Create(studio)
	if err != nil {
		return nil, fmt.Errorf("interactor.profile_usecase.Create: %w", err)
	}
	return studio, nil
}

func (pu *StudioUsecase) FilterByOwnerID(ownerID string) (*[]comusic.Studio, error) {
	user, err := pu.StudioRepository.FilterByOwnerID(ownerID)
	if err != nil {
		return nil, fmt.Errorf("interactor.profile_usecase.Get: %w", err)
	}
	return user, err
}
