package mysql

import (
	"database/sql"
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
	// Set null if zero value
	var activeTake *string = nil
	if t.ActiveTake != "" {
		activeTake = &t.ActiveTake
	}
	_, err := r.Exec(
		`INSERT INTO tracks (
			id, version_id, created_at, updated_at, name,
			volume, pan, is_muted, is_soloed, icon, active_take
		)
		VALUES (
			?, ?, ?, ?, ?, 
			?, ?, ?, ?, ?, ?
		)`,
		t.ID, t.VersionID, t.CreatedAt, t.UpdatedAt, t.Name,
		t.Volume, t.Pan, t.IsMuted, t.IsSoloed, t.Icon, activeTake,
	)
	if err != nil {
		return fmt.Errorf("mysql.track_repository.Create: %w", err)
	}
	return nil
}

func (r *TrackRepository) GetByID(id string) (*comusic.Track, error) {
	tr := &comusic.Track{Meta: &comusic.Meta{}}
	active_take := &sql.NullString{}
	row := r.DB.QueryRow(
		`SELECT id, version_id, active_take, created_at, updated_at, name,
		volume, pan, is_muted, is_soloed, icon
		FROM tracks WHERE id = ?`,
		id,
	)
	err := row.Scan(
		&tr.Meta.ID, &tr.VersionID, active_take, &tr.CreatedAt, &tr.UpdatedAt, &tr.Name,
		&tr.Volume, &tr.Pan, &tr.IsMuted, &tr.IsSoloed, &tr.Icon,
	)
	if err != nil {
		return nil, fmt.Errorf("mysql.track_repository.GetByID: %w", err)
	}
	if active_take.Valid {
		tr.ActiveTake = active_take.String
	} else {
		tr.ActiveTake = ""
	}
	return tr, nil
}

func (r *TrackRepository) Update(in *comusic.TrackUpdateInput) error {
	_, err := r.Exec(`
		UPDATE tracks SET
		updated_at = ?,
		version_id = COALESCE(?, version_id),
		active_take = COALESCE(?, active_take),
		name = COALESCE(?, name),
		volume = COALESCE(?, volume),
		pan = COALESCE(?, pan),
		is_muted = COALESCE(?, is_muted),
		is_soloed = COALESCE(?, is_soloed),
		icon = COALESCE(?, icon)
		WHERE id = ?`,
		in.UpdatedAt, in.VerID, in.ActiveTake, in.Name,
		in.Vol, in.Pan, in.IsMuted, in.IsSoloed, in.Icon,
		in.ID,
	)
	if err != nil {
		return fmt.Errorf("mysql.track_repository.Update: %w", err)
	}
	return nil
}

func (r *TrackRepository) Delete(id string) error {
	_, err := r.Exec(`DELETE FROM tracks WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("mysql.track_repository.Delete: %w", err)
	}
	return nil
}
