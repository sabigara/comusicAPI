package mock

import (
	"time"

	comusic "github.com/sabigara/comusicAPI"
)

func Aunthenticate(credentials ...interface{}) (*comusic.User, error) {
	return &comusic.User{
		Meta: &comusic.Meta{
			ID:        "4148e7cc-a5f0-4fb4-9392-ee82f0e324d1",
			CreatedAt: time.Date(2019, time.May, 5, 13, 15, 10, 0, time.UTC),
			UpdatedAt: time.Date(2020, time.February, 5, 11, 9, 13, 0, time.UTC),
		},
		Email: "email@example.com",
	}, nil
}
