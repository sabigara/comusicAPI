package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	comusic "github.com/sabigara/comusicAPI"
)

type TrackRepository struct {
	*sqlx.DB
}

func NewTrackRepository(db *sqlx.DB) *TrackRepository {
	return &TrackRepository{DB: db}
}

func (r *TrackRepository) Create(t *comusic.Track) error {
	_, err := r.Exec(
		`INSERT INTO tracks (
			id, version_id, created_at, updated_at, name,
			pan, is_muted, is_soloed, icon, active_take
		)
		VALUES (
			?, ?, ?, ?, ?, 
			?, ?, ?, ?, ?
		)`,
		t.ID, t.VersionID, t.CreatedAt, t.UpdatedAt, t.Name,
		t.Pan, t.IsMuted, t.IsSoloed, t.Icon, t.ActiveTake,
	)
	if err != nil {
		return fmt.Errorf("mysql.track_repository.Create: %w", err)
	}
	return nil
}

func (r *TrackRepository) FilterByVersionIDWithTakes(versionID string) (comusic.TrackTakeMap, error) {
	dict := make(comusic.TrackTakeMap)
	rows, err := r.Query(`
		SELECT
		tr.id, tr.version_id, tr.created_at, tr.updated_at, tr.name,
		tr.pan, tr.is_muted, tr.is_soloed, tr.icon, tr.active_take,
			COALESCE(takes.id, ''),
			COALESCE(takes.track_id, ''),
			COALESCE(takes.created_at, NOW()),
			COALESCE(takes.updated_at, NOW()),
			COALESCE(takes.name, '')
		FROM tracks as tr
		LEFT OUTER JOIN takes ON tr.id = takes.track_id
		WHERE tr.version_id = ?`,
		versionID,
	)
	if err != nil {
		return nil, fmt.Errorf("mysql.track_repository.FilterByStudioIDWithVersions: %w", err)
	}
	for rows.Next() {
		tr := &comusic.Track{Meta: &comusic.Meta{}}
		tk := &comusic.Take{Meta: &comusic.Meta{}}
		err := rows.Scan(
			&tr.ID, &tr.VersionID, &tr.CreatedAt, &tr.UpdatedAt, &tr.Name,
			&tr.Pan, &tr.IsMuted, &tr.IsSoloed, &tr.Icon, &tr.ActiveTake,
			&tk.ID, &tk.TrackID, &tk.CreatedAt, &tk.UpdatedAt, &tk.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("mysql.track_repository.FilterByStudioIDWithVersions: %w", err)
		}
		if val, ok := dict[tr.ID]; !ok {
			dict[tr.ID] = &comusic.TrackTake{
				Data: tr,
			}
			if tk.ID != "" {
				dict[tr.ID].Takes = []*comusic.Take{tk}
			}
		} else {
			if tk.ID != "" {
				val.Takes = append(val.Takes, tk)
			}
		}
	}
	return dict, nil
}
