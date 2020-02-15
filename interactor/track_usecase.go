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

func (tu *TrackUsecase) Create(verID, name string) (*comusic.Track, error) {
	track := comusic.NewTrack(verID, name)
	err := tu.TrackRepository.Create(track)
	if err != nil {
		return nil, fmt.Errorf("interactor.track_usecase.Create: %w", err)
	}
	return track, nil
}

func (tu *TrackUsecase) GetByID(id string) (*comusic.Track, error) {
	tr, err := tu.TrackRepository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("interactor.track_usecase.GetByID: %w", err)
	}
	return tr, nil
}

func (tu *TrackUsecase) Update(in *comusic.TrackUpdateInput) error {
	err := tu.TrackRepository.Update(in)
	if err != nil {
		return fmt.Errorf("interactor.track_usecase.Update: %w", err)
	}
	return nil
}

func (tu *TrackUsecase) Delete(id string) error {
	errf := "interactor.track_usecase.Delete: %w"
	track, err := tu.TrackRepository.GetByID(id)
	if err != nil {
		return fmt.Errorf(errf, err)
	}
	takes, err := tu.TakeUsecase.FilterByTrackID(track.ID)
	if err != nil {
		return fmt.Errorf(errf, err)
	}
	for _, tk := range takes {
		if err := tu.TakeUsecase.Delete(tk.ID); err != nil {
			return fmt.Errorf(errf, err)
		}
	}
	if err := tu.TrackRepository.Delete(id); err != nil {
		return fmt.Errorf(errf, err)
	}
	return nil
}
