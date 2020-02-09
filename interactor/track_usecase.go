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

func (tu *TrackUsecase) FilterByVersionIDWithTakes(verID string) (comusic.TrackTakeMap, error) {
	trackTakeMap, err := tu.TrackRepository.FilterByVersionIDWithTakes(verID)
	if err != nil {
		return nil, fmt.Errorf("interactor.track_usecase.FilterByStudioIDWithVersions: %w", err)
	}
	return trackTakeMap, nil
}
