package firebase

import (
	"context"
	"fmt"
	"time"

	fb "firebase.google.com/go"
	comusic "github.com/sabigara/comusicAPI"
)

type UserRepository struct {
	fbapp *fb.App
}

func NewUserRepository(fbapp *fb.App) *UserRepository {
	return &UserRepository{fbapp: fbapp}
}

func (r *UserRepository) GetByEmail(email string) (*comusic.User, error) {
	ctx := context.Background()
	client, err := r.fbapp.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("firebase.user_repository.GetByEmail: %w", err)
	}
	user, err := client.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("firebase.user_repository.GetByEmail: %w", err)
	}
	return &comusic.User{
		Meta: &comusic.Meta{
			ID: user.UID,
			// CreationTimestamp is in milliseconds so divide by 1000.
			CreatedAt: time.Unix(user.UserMetadata.CreationTimestamp/1000, 0),
		},
		Email: user.Email,
	}, nil
}
