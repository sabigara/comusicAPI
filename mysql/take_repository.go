package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	comusic "github.com/sabigara/comusicAPI"
)

type TakeRepository struct {
	*sqlx.DB
}

func NewTakeRepository(db *sqlx.DB) *TakeRepository {
	return &TakeRepository{DB: db}
}

func (r *TakeRepository) Create(take *comusic.Take) error {
	_, err := r.Exec(`
		INSERT INTO takes (
			id, track_id, file_id, created_at, updated_at, name
		)
		VALUES (
			?, ?, ?, ?, ?, ?
		)`,
		take.ID, take.TrackID, take.FileID, take.CreatedAt, take.UpdatedAt, take.Name,
	)
	if err != nil {
		return fmt.Errorf("mysql.take_repository.Create: %w", err)
	}
	return nil
}

func (r *TakeRepository) GetByID(id string) (*comusic.Take, error) {
	tk := &comusic.Take{}
	err := r.DB.Get(
		tk,
		`SELECT * FROM takes WHERE id = ?`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("mysql.take_repository.GetByID: %w", err)
	}
	return tk, nil
}

func (r *TakeRepository) FilterByTrackID(trackID string) ([]*comusic.Take, error) {
	takes := &[]*comusic.Take{}
	err := r.DB.Select(
		takes,
		`SELECT * FROM takes WHERE track_id = ?`,
		trackID,
	)
	if err != nil {
		return nil, fmt.Errorf("mysql.take_repository.FilterByTrackID: %w", err)
	}
	return *takes, nil
}

func (r *TakeRepository) Delete(takeID string) error {
	_, err := r.Exec(`DELETE FROM takes WHERE id = ?`, takeID)
	if err != nil {
		return fmt.Errorf("mysql.take_repository.Delete: %w", err)
	}
	return nil
}
