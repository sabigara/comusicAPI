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
