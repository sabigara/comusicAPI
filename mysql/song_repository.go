package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	comusic "github.com/sabigara/comusicAPI"
)

type SongRepository struct {
	*sqlx.DB
}

func NewSongRepository(db *sqlx.DB) *SongRepository {
	return &SongRepository{DB: db}
}

func (r *SongRepository) Create(song *comusic.Song) error {
	_, err := r.Exec(`
		INSERT INTO songs (id, studio_id, created_at, updated_at, name)
		VALUES (?, ?, ?, ?, ?)`,
		song.ID, song.StudioID, song.CreatedAt, song.UpdatedAt, song.Name,
	)
	if err != nil {
		return fmt.Errorf("mysql.song_repository.Create: %w", err)
	}
	return nil
}

func (r *SongRepository) FilterByStudioIDWithVersions(studioID string) (comusic.SongVerMap, error) {
	dict := make(comusic.SongVerMap)
	rows, err := r.Query(`
		SELECT songs.id, songs.studio_id, songs.created_at, songs.updated_at, songs.name,
		COALESCE(vers.id, ''),
		COALESCE(vers.song_id, ''),
		COALESCE(vers.created_at, NOW()),
		COALESCE(vers.updated_at, NOW()),
		COALESCE(vers.name, '')
		FROM songs
		LEFT OUTER JOIN versions AS vers ON songs.id = vers.song_id
		WHERE songs.studio_id = ?`,
		studioID,
	)
	if err != nil {
		return nil, fmt.Errorf("mysql.song_repository.FilterByStudioIDWithVersions: %w", err)
	}
	for rows.Next() {
		s := &comusic.Song{Meta: &comusic.Meta{}}
		v := &comusic.Version{Meta: &comusic.Meta{}}
		err := rows.Scan(
			&s.ID, &s.StudioID, &s.CreatedAt, &s.UpdatedAt, &s.Name,
			&v.ID, &v.SongID, &v.CreatedAt, &v.UpdatedAt, &v.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("mysql.song_repository.FilterByStudioIDWithVersions: %w", err)
		}
		if val, ok := dict[s.ID]; !ok {
			dict[s.ID] = &comusic.SongVer{
				Data: s,
			}
			if v.ID != "" {
				dict[s.ID].Versions = []*comusic.Version{v}
			}
		} else {
			if v.ID != "" {
				val.Versions = append(val.Versions, v)
			}
		}
	}
	return dict, nil
}
