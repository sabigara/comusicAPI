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

func (r SongRepository) FilterByGuestID(guestID string) ([]*comusic.Song, error) {
	songs := []*comusic.Song{}
	rows, err := r.Query(`
		SELECT songs.id, songs.studio_id, songs.created_at, songs.updated_at, songs.name
		FROM songs 
		INNER JOIN guest_song ON songs.id = guest_song.song_id
		WHERE guest_song.user_id = ?`,
		guestID,
	)
	if err != nil {
		return nil, fmt.Errorf("mysql.song_repository.FilterByGuestID: %w", err)
	}
	for rows.Next() {
		s := &comusic.Song{Meta: &comusic.Meta{}}
		err := rows.Scan(
			&s.ID, &s.StudioID, &s.CreatedAt, &s.UpdatedAt, &s.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("mysql.studio_repository.FilterByGuestID: %w", err)
		}
		songs = append(songs, s)
	}
	return songs, nil
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

func (r *SongRepository) AddGuests(songID string, userIDs ...string) error {
	for _, userID := range userIDs {
		_, err := r.Exec(`
			INSERT INTO guest_song (user_id, song_id)
			VALUES (?, ?)`,
			userID, songID,
		)
		if err != nil {
			return fmt.Errorf("mysql.song_repository.AddMembers: %w", err)
		}
	}
	return nil
}
