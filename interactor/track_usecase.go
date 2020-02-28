package interactor

import (
	"fmt"

	comusic "github.com/sabigara/comusicAPI"
)

type TrackUsecase struct {
	comusic.TrackRepository
	comusic.TakeUsecase
}

func NewTrackUsecase(pr comusic.TrackRepository) *TrackUsecase {
	return &TrackUsecase{TrackRepository: pr}
}

func (u *TrackUsecase) Create(verID, name string) (*comusic.Track, error) {
	track := comusic.NewTrack(verID, name)
	err := u.TrackRepository.Create(track)
	if err != nil {
		return nil, fmt.Errorf("interactor.track_usecase.Create: %w", err)
	}
	return track, nil
}

func (u *TrackUsecase) GetByID(id string) (*comusic.Track, error) {
	tr, err := u.TrackRepository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("interactor.track_usecase.GetByID: %w", err)
	}
	return tr, nil
}

func (u *TrackUsecase) Update(in *comusic.TrackUpdateInput) error {
	err := u.TrackRepository.Update(in)
	if err != nil {
		return fmt.Errorf("interactor.track_usecase.Update: %w", err)
	}
	return nil
}

func (u *TrackUsecase) Delete(id string) error {
	errf := "interactor.track_usecase.Delete: %w"
	track, err := u.TrackRepository.GetByID(id)
	if err != nil {
		return fmt.Errorf(errf, err)
	}
	takes, err := u.TakeUsecase.FilterByTrackID(track.ID)
	if err != nil {
		return fmt.Errorf(errf, err)
	}
	for _, tk := range takes {
		if err := u.TakeUsecase.Delete(tk.ID); err != nil {
			return fmt.Errorf(errf, err)
		}
	}
	if err := u.TrackRepository.Delete(id); err != nil {
		return fmt.Errorf(errf, err)
	}
	return nil
}
