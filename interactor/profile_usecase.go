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

func (u *ProfileUsecase) Create(userID, nickname, bio string) (*comusic.Profile, error) {
	profile := comusic.NewProfile(userID, nickname, bio)
	err := u.ProfileRepository.Create(profile)
	if err != nil {
		return nil, fmt.Errorf("interactor.profile_usecase.Create: %w", err)
	}
	return profile, nil
}

func (u *ProfileUsecase) Update(userID string, nickname, bio *string) error {
	err := u.ProfileRepository.Update(userID, nickname, bio)
	if err != nil {
		return fmt.Errorf("interactor.profile_usecase.Create: %w", err)
	}
	return nil
}

func (u *ProfileUsecase) GetByUserID(userID string) (*comusic.Profile, error) {
	user, err := u.ProfileRepository.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("interactor.profile_usecase.Get: %w", err)
	}
	return user, err
}
