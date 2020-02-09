package mysql

import (
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
	_, err := r.Exec(
		`INSERT INTO versions (id, song_id, created_at, updated_at, name)
		 VALUES (?, ?, ?, ?, ?)`,
		version.ID, version.SongID, version.CreatedAt, version.UpdatedAt, version.Name,
	)
	if err != nil {
		return fmt.Errorf("mysql.version_repository.Create: %w", err)
	}
	return nil
}
