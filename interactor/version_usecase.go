package interactor

import (
	"fmt"

	comusic "github.com/sabigara/comusicAPI"
)

type VersionUsecase struct {
	comusic.VersionRepository
	comusic.FileRepository
}

func NewVersionUsecase(pr comusic.VersionRepository, fr comusic.FileRepository) *VersionUsecase {
	return &VersionUsecase{VersionRepository: pr, FileRepository: fr}
}

func (su *VersionUsecase) Create(songID, nickname string) (*comusic.Version, error) {
	version := comusic.NewVersion(songID, nickname)
	err := su.VersionRepository.Create(version)
	if err != nil {
		return nil, fmt.Errorf("interactor.version_usecase.Create: %w", err)
	}
	return version, nil
}

func (tu *VersionUsecase) GetContents(verID string) ([]*comusic.Track, []*comusic.Take, []*comusic.File, error) {
	tracks, takes, err := tu.VersionRepository.GetContents(verID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("interactor.track_usecase.FilterByStudioIDWithVersions: %w", err)
	}
	files := []*comusic.File{}
	for _, tk := range takes {
		files = append(files, &comusic.File{
			Meta: &comusic.Meta{
				ID: tk.FileID,
			},
			URL: tu.FileRepository.URL(tk.FileID),
		})
	}
	return tracks, takes, files, nil
}
