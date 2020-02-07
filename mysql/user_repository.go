package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	comusic "github.com/sabigara/comusicAPI"
)

type UserRepository struct {
	*sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (ur *UserRepository) Save(user *comusic.User) error {
	_, err := ur.Exec(
		`INSERT INTO user (id, name, email, password)
		 VALUES (?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE name = ?
		 `,
		user.ID, user.Name, user.Email, user.Password,
		user.Name,
	)
	if err != nil {
		return fmt.Errorf("mysql.user_repository.Create: %w", err)
	}
	return nil
}

func (ur *UserRepository) GetById(id string) (u *comusic.User, err error) {
	row := ur.QueryRow(
		`SELECT id, name, email, password
		 FROM user
		 WHERE id = ?`, id,
	)
	u = &comusic.User{}
	err = row.Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("mysql.user_repository.Get: %w", comusic.ErrResourceNotFound)
		}
		return nil, fmt.Errorf("mysql.user_repository.Get: %w", err)
	}
	return
}

func (ur *UserRepository) GetByEmail(email string) (u *comusic.User, err error) {
	row := ur.QueryRow(
		`SELECT id, name, email, password
		 FROM user
		 WHERE email = ?`, email,
	)
	u = &comusic.User{}
	err = row.Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("mysql.user_repository.Get: %w", comusic.ErrResourceNotFound)
		}
		return nil, fmt.Errorf("mysql.user_repository.Get: %w", err)
	}
	return
}
