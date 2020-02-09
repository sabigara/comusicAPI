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
			id, track_id, created_at, updated_at, name,
			file_name, file_url
		)
		VALUES (
			?, ?, ?, ?, ?,
			?, ?
		)`,
		take.ID, take.TrackID, take.CreatedAt, take.UpdatedAt, take.Name,
		take.FileName, take.FileURL,
	)
	if err != nil {
		return fmt.Errorf("mysql.take_repository.Create: %w", err)
	}
	return nil
}
