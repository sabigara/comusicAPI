package interactor

import (
	"fmt"

	comusic "github.com/sabigara/comusicAPI"
)

type TrackUsecase struct {
	comusic.TrackRepository
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

func (tu *TrackUsecase) FilterByVersionIDWithTakes(verID string) (comusic.TrackTakeMap, error) {
	trackTakeMap, err := tu.TrackRepository.FilterByVersionIDWithTakes(verID)
	if err != nil {
		return nil, fmt.Errorf("interactor.track_usecase.FilterByStudioIDWithVersions: %w", err)
	}
	return trackTakeMap, nil
}
