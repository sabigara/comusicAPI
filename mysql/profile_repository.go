package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	comusic "github.com/sabigara/comusicAPI"
)

type ProfileRepository struct {
	*sqlx.DB
}

func NewProfileRepository(db *sqlx.DB) *ProfileRepository {
	return &ProfileRepository{DB: db}
}

func (r *ProfileRepository) Create(p *comusic.Profile) error {
	_, err := r.Exec(
		`INSERT INTO profiles (id, user_id, created_at, updated_at, nickname, bio)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		p.ID, p.UserID, p.CreatedAt, p.UpdatedAt, p.Nickname, p.Bio,
	)
	if err != nil {
		return fmt.Errorf("mysql.profile_repository.Save: %w", err)
	}
	return nil
}

func (r *ProfileRepository) Update(id string, nickname, bio *string) error {
	_, err := r.Exec(
		`UPDATE profiles SET
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

func (r *ProfileRepository) GetByUserID(userID string) (*comusic.Profile, error) {
	p := &comusic.Profile{}
	err := r.Get(
		p,
		`SELECT * FROM profiles
		 WHERE user_id = ?`,
		userID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("mysql.profile_repository.Get: %w", comusic.ErrResourceNotFound)
		}
		return nil, fmt.Errorf("mysql.profile_repository.Get: %w", err)
	}
	return p, nil
}
