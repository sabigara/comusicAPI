package interactor

import (
	"fmt"

	comusic "github.com/sabigara/comusicAPI"
)

type VersionUsecase struct {
	comusic.VersionRepository
}

func NewVersionUsecase(pr comusic.VersionRepository) *VersionUsecase {
	return &VersionUsecase{VersionRepository: pr}
}

func (su *VersionUsecase) Create(songID, nickname string) (*comusic.Version, error) {
	version := comusic.NewVersion(songID, nickname)
	err := su.VersionRepository.Create(version)
	if err != nil {
		return nil, fmt.Errorf("interactor.version_usecase.Create: %w", err)
	}
	return version, nil
}

func (tu *VersionUsecase) GetContents(verID string) ([]*comusic.Track, []*comusic.Take, error) {
	tracks, takes, err := tu.VersionRepository.GetContents(verID)
	if err != nil {
		return nil, nil, fmt.Errorf("interactor.track_usecase.FilterByStudioIDWithVersions: %w", err)
	}
	return tracks, takes, nil
}
