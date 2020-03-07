package mysql

import (
	"database/sql"
	"errors"
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

func (r *SongRepository) GetByID(id string) (*comusic.Song, error) {
	s := &comusic.Song{}
	err := r.Get(
		s,
		`SELECT * FROM songs 
		 WHERE id = ?`,
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("mysql.song_repository.GetByID: %w", comusic.ErrResourceNotFound)
		}
		return nil, fmt.Errorf("mysql.song_repository.Get: %w", err)
	}
	return s, nil
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

func (r *SongRepository) Delete(songID string) error {
	_, err := r.Exec(`DELETE FROM songs WHERE id = ?`, songID)
	if err != nil {
		return fmt.Errorf("mysql.song_repository.Delete: %w", err)
	}
	return nil
}
