package interactor

import (
	"fmt"

	comusic "github.com/sabigara/comusicAPI"
)

type SongUsecase struct {
	comusic.SongRepository
}

func NewSongUsecase(pr comusic.SongRepository) *SongUsecase {
	return &SongUsecase{SongRepository: pr}
}

func (u *SongUsecase) Create(studioID, nickname string) (*comusic.Song, error) {
	song := comusic.NewSong(studioID, nickname)
	err := u.SongRepository.Create(song)
	if err != nil {
		return nil, fmt.Errorf("interactor.song_usecase.Create: %w", err)
	}
	return song, nil
}
