package comusic

import (
	"time"

	"github.com/google/uuid"
)

type Meta struct {
	ID        string    `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewMeta() *Meta {
	return &Meta{
		ID:        uuid.New().String(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

type User struct {
	*Meta
	Email string `json:"email"`
}

func NewUser(id, displayName, email string) *User {
	return &User{
		Meta:  NewMeta(),
		Email: email,
	}
}

type Profile struct {
	*Meta
	UserID   string `json:"user_id" db:"user_id"`
	Nickname string `json:"nickname" db:"nickname"`
	Bio      string `json:"bio" db:"bio"`
}

func NewProfile(userID, nickname, bio string) *Profile {
	return &Profile{
		Meta:     NewMeta(),
		UserID:   userID,
		Nickname: nickname,
		Bio:      bio,
	}
}

type ProfileUsecase interface {
	Create(userID, nickname, bio string) (*Profile, error)
	Update(userID string, nickname, bio *string) error
	GetByUserID(userID string) (*Profile, error)
}

type ProfileRepository interface {
	Create(*Profile) error
	Update(userID string, nickname, bio *string) error
	GetByUserID(userID string) (*Profile, error)
}

type Studio struct {
	*Meta
	OwnerID string `json:"owner_id" db:"owner_id"`
	Name    string `json:"name" db:"name"`
}

func NewStudio(owenerID, name string) *Studio {
	return &Studio{
		Meta:    NewMeta(),
		OwnerID: owenerID,
		Name:    name,
	}
}

type StudioUsecase interface {
	Create(ownerID, name string) (*Studio, error)
	FilterByOwnerID(ownerID string) (*[]Studio, error)
}

type StudioRepository interface {
	Create(*Studio) error
	FilterByOwnerID(id string) (*[]Studio, error)
}

type Song struct {
	*Meta
	StudioID string
	Name     string
}

type SongUsecase struct {
	FilterByStudioIDWithVersions func(studioID string) ([]*Song, []*Version, error)
}

type SongRepository struct {
	// Prefetch versions
	FilterByStudioIDWithVersions func(studioID string) ([]*Song, []*Version, error)
}

type Version struct {
	*Meta
	SongID string
	Name   string
}

type Track struct {
	*Meta
	Name       string
	Pan        int
	IsMuted    bool
	IsSoloed   bool
	Icon       int
	ActiveTake string
	Takes      []string
}

type TrackUsecase struct {
	FilterByVersionIDWithTakes func(verID string) ([]*Track, []*Take, error)
}

type Take struct {
	*Meta
	Name     string
	FileName string
	FileURL  string
}
