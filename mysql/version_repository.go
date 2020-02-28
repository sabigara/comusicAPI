package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	comusic "github.com/sabigara/comusicAPI"
)

type VersionRepository struct {
	*sqlx.DB
}

func NewVersionRepository(db *sqlx.DB) *VersionRepository {
	return &VersionRepository{DB: db}
}

func (r *VersionRepository) Create(version *comusic.Version) error {
	_, err := r.Exec(`
		INSERT INTO versions (id, song_id, created_at, updated_at, name)
		VALUES (?, ?, ?, ?, ?)`,
		version.ID, version.SongID, version.CreatedAt, version.UpdatedAt, version.Name,
	)
	if err != nil {
		return fmt.Errorf("mysql.version_repository.Create: %w", err)
	}
	return nil
}

func (r *VersionRepository) Delete(verID string) error {
	_, err := r.Exec(`DELETE FROM versions WHERE id = ?`, verID)
	if err != nil {
		return fmt.Errorf("mysql.version_repository.Delete: %w", err)
	}
	return nil
}

func (r *VersionRepository) GetContents(versionID string) ([]*comusic.Track, []*comusic.Take, error) {
	// tracks.active_take is nullable
	var activeTake sql.NullString
	tracks := []*comusic.Track{}
	takes := []*comusic.Take{}
	rows, err := r.Query(`
		SELECT
		tr.id, tr.version_id, tr.created_at, tr.updated_at, tr.name,
		tr.volume, tr.pan, tr.is_muted, tr.is_soloed, tr.icon, tr.active_take,
			COALESCE(takes.id, ''),
			COALESCE(takes.track_id, ''),
			COALESCE(takes.file_id, ''),
			COALESCE(takes.created_at, NOW()),
			COALESCE(takes.updated_at, NOW()),
			COALESCE(takes.name, '')
		FROM tracks as tr
		LEFT OUTER JOIN takes ON tr.id = takes.track_id
		WHERE tr.version_id = ?
		ORDER BY tr.id`,
		versionID,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("mysql.version_repository.GetContents: %w", err)
	}
	for rows.Next() {
		tr := &comusic.Track{Meta: &comusic.Meta{}}
		tk := &comusic.Take{Meta: &comusic.Meta{}}
		err := rows.Scan(
			&tr.ID, &tr.VersionID, &tr.CreatedAt, &tr.UpdatedAt, &tr.Name,
			&tr.Volume, &tr.Pan, &tr.IsMuted, &tr.IsSoloed, &tr.Icon, &activeTake,
			&tk.ID, &tk.TrackID, &tk.FileID, &tk.CreatedAt, &tk.UpdatedAt, &tk.Name,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("mysql.version_repository.GetContents: %w", err)
		}
		// Replace null for DB with zero value for string in Go
		if activeTake.Valid {
			tr.ActiveTake = activeTake.String
		} else {
			tr.ActiveTake = ""
		}
		tracks = append(tracks, tr)
		if tk.ID != "" {
			takes = append(takes, tk)
		}
	}
	return tracks, takes, nil
}
