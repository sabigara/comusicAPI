package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	comusic "github.com/sabigara/comusicAPI"
)

type StudioRepository struct {
	*sqlx.DB
}

func NewStudioRepository(db *sqlx.DB) *StudioRepository {
	return &StudioRepository{DB: db}
}

func (r *StudioRepository) Create(p *comusic.Studio) error {
	_, err := r.Exec(`
		INSERT INTO studios (id, owner_id, created_at, updated_at, name)
		VALUES (?, ?, ?, ?, ?)`,
		p.ID, p.OwnerID, p.CreatedAt, p.UpdatedAt, p.Name,
	)
	if err != nil {
		return fmt.Errorf("mysql.studio_repository.Save: %w", err)
	}
	return nil
}

func (r *StudioRepository) Update(id string, nickname, bio *string) error {
	_, err := r.Exec(`
		UPDATE profiles SET
		nickname = COALESCE(?, nickname),
		bio = COALESCE(?, bio)
		WHERE user_id = ?`,
		nickname, bio, id,
	)
	if err != nil {
		return fmt.Errorf("mysql.profile_repository.Update: %w", err)
	}
	return nil
}

func (r *StudioRepository) FilterByOwnerID(ownerID string) (*[]comusic.Studio, error) {
	p := &[]comusic.Studio{}
	err := r.Select(
		p,
		`SELECT * FROM studios
		 WHERE owner_id = ?`,
		ownerID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("mysql.profile_repository.Get: %w", comusic.ErrResourceNotFound)
		}
		return nil, fmt.Errorf("mysql.profile_repository.Get: %w", err)
	}
	return p, nil
}
