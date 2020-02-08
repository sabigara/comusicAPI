package interactor

import (
	"fmt"

	comusic "github.com/sabigara/comusicAPI"
)

type ProfileUsecase struct {
	comusic.ProfileRepository
}

func NewProfileUsecase(pr comusic.ProfileRepository) *ProfileUsecase {
	return &ProfileUsecase{ProfileRepository: pr}
}

func (pu *ProfileUsecase) Create(userID, nickname, bio string) (*comusic.Profile, error) {
	profile := comusic.NewProfile(userID, nickname, bio)
	err := pu.ProfileRepository.Create(profile)
	if err != nil {
		return nil, fmt.Errorf("interactor.profile_usecase.Create: %w", err)
	}
	return profile, nil
}

func (pu *ProfileUsecase) Update(userID string, nickname, bio *string) error {
	err := pu.ProfileRepository.Update(userID, nickname, bio)
	if err != nil {
		return fmt.Errorf("interactor.profile_usecase.Create: %w", err)
	}
	return nil
}

func (pu *ProfileUsecase) GetByUserID(userID string) (*comusic.Profile, error) {
	user, err := pu.ProfileRepository.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("interactor.profile_usecase.Get: %w", err)
	}
	return user, err
}
