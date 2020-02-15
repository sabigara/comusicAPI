package comusic

import (
	"time"

	"github.com/google/uuid"
)

type Meta struct {
	ID        string    `json:"id" db:"id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
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
	UserID   string `json:"user_id" db:"userId"`
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
	OwnerID string `json:"owner_id" db:"ownerId"`
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
	StudioID string `json:"studioId"`
	Name     string `json:"name"`
}

func NewSong(studioID, name string) *Song {
	return &Song{
		Meta:     NewMeta(),
		StudioID: studioID,
		Name:     name,
	}
}

type SongVer struct {
	Data     *Song
	Versions []*Version
}
type SongVerMap map[string]*SongVer

type SongUsecase interface {
	Create(studioID, name string) (*Song, error)
	FilterByStudioIDWithVersions(studioID string) (SongVerMap, error)
}

type SongRepository interface {
	Create(*Song) error
	// Prefetch versions
	FilterByStudioIDWithVersions(studioID string) (SongVerMap, error)
}

type Version struct {
	*Meta
	SongID string `json:"songId"`
	Name   string `json:"name"`
}

func NewVersion(songID, name string) *Version {
	return &Version{
		Meta:   NewMeta(),
		SongID: songID,
		Name:   name,
	}
}

type VersionUsecase interface {
	Create(songID, name string) (*Version, error)
	GetContents(verID string) ([]*Track, []*Take, error)
}

type VersionRepository interface {
	Create(*Version) error
	GetContents(verID string) ([]*Track, []*Take, error)
}

type Track struct {
	*Meta
	VersionID  string  `json:"versionId" db:"version_id"`
	Name       string  `json:"name" db:"name"`
	Volume     float32 `json:"volume" db:"volume"`
	Pan        float32 `json:"pan" db:"pan"`
	IsMuted    bool    `json:"isMuted" db:"is_muted"`
	IsSoloed   bool    `json:"isSoloed" db:"is_soloed"`
	Icon       int     `json:"icon" db:"icon"`
	ActiveTake string  `json:"activeTake" db:"active_take"`
}

func NewTrack(verID, name string) *Track {
	return &Track{
		Meta:      NewMeta(),
		VersionID: verID,
		Name:      name,
		Volume:    0.7,
	}
}

const (
	DrumsIcon = iota
	BassIcon
	GuitarIcon
	VocalIcon
	KeyboardIcon
)

type TrackTake struct {
	Data  *Track
	Takes []*Take
}
type TrackTakeMap map[string]*TrackTake
type TrackUpdateInput struct {
	ID         string
	UpdatedAt  time.Time
	VerID      *string
	ActiveTake *string
	Name       *string
	Vol        *float32
	Pan        *float32
	IsMuted    *bool
	IsSoloed   *bool
	Icon       *int
}

type TrackUsecase interface {
	Create(verID, name string) (*Track, error)
	GetByID(id string) (*Track, error)
	Update(*TrackUpdateInput) error
	Delete(id string) error
}

type TrackRepository interface {
	Create(*Track) error
	GetByID(id string) (*Track, error)
	Update(*TrackUpdateInput) error
	Delete(id string) error
}

type Take struct {
	*Meta
	TrackID string `json:"trackId" db:"track_id"`
	Name    string `json:"name" db:"name"`
	FileID  string `json:"fileId" db:"file_id"`
}

func NewTake(trackID, name string) *Take {
	return &Take{
		Meta:    NewMeta(),
		TrackID: trackID,
		Name:    name,
	}
}

type TakeUsecase interface {
	Create(trackID, name string, src FileSrc) (*Take, *File, error)
	GetByID(id string) (*Take, error)
	FilterByTrackID(trackID string) ([]*Take, error)
	Delete(takeId string) error
}

type TakeRepository interface {
	Create(*Take) error
	GetByID(id string) (*Take, error)
	FilterByTrackID(trackID string) ([]*Take, error)
	Delete(takeId string) error
}

type File struct {
	*Meta
	URL string  `json:"url"`
	Src FileSrc `json:"-"`
}

type FileSrc interface {
	Read(p []byte) (n int, err error)
}

func NewFile(url string, src FileSrc) *File {
	return &File{
		Meta: NewMeta(),
		Src:  src,
	}
}

type FileRepository interface {
	Upload(file *File) (*File, error)
	Download(url string) (*File, error)
	Delete(id string) error
	URL(fileID string) string
}
