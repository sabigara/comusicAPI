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

type UserRepository interface {
	GetByEmail(email string) (*User, error)
}

type AuthUsecase interface {
	Authenticate(credentials ...interface{}) (*User, error)
}

type Profile struct {
	*Meta
	UserID   string `json:"userId" db:"user_id"`
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
	OwnerID string `json:"ownerId" db:"owner_id"`
	Name    string `json:"name" db:"name"`
}

func NewStudio(ownerID, name string) *Studio {
	return &Studio{
		Meta:    NewMeta(),
		OwnerID: ownerID,
		Name:    name,
	}
}

type StudioUsecase interface {
	GetByID(id string) (*Studio, error)
	Create(ownerID, name string) (*Studio, error)
	Filter(ownerID, memberID string) ([]*Studio, error)
	GetContents(studioID string) ([]*Song, []*Version, error)
	AddMembers(studioID string, userID ...string) error
}

type StudioRepository interface {
	GetByID(id string) (*Studio, error)
	Create(*Studio) error
	FilterByOwnerID(id string) ([]*Studio, error)
	FilterByMemberID(id string) ([]*Studio, error)
	GetContents(studioID string) ([]*Song, []*Version, error)
	AddMembers(studioID string, userID ...string) error
}

type Song struct {
	*Meta
	StudioID string `json:"studioId" db:"studio_id"`
	Name     string `json:"name" db:"name"`
}

func NewSong(studioID, name string) *Song {
	return &Song{
		Meta:     NewMeta(),
		StudioID: studioID,
		Name:     name,
	}
}

type SongUsecase interface {
	Filter(guestID string) ([]*Song, error)
	GetByID(songID string) (*Song, error)
	Create(studioID, name string) (*Song, error)
	Delete(songID string) error
	AddGuests(songID string, userID ...string) error
}

type SongRepository interface {
	FilterByGuestID(guestID string) ([]*Song, error)
	GetByID(songID string) (*Song, error)
	Create(*Song) error
	Delete(songID string) error
	AddGuests(songID string, userID ...string) error
}

type Version struct {
	*Meta
	SongID string `json:"songId" db:"song_id"`
	Name   string `json:"name" db:"name"`
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
	Delete(verId string) error
	GetContents(verID string) ([]*Track, []*Take, []*File, error)
}

type VersionRepository interface {
	Create(*Version) error
	Delete(verId string) error
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

type MailUsecase interface {
	InviteToStudioNew(to, studio_name, signupURL string) error
	InviteToStudio(user *User, studio_name, signInURL string) error
}

type GroupType = int

const (
	ErrGroupType GroupType = iota
	StudioGroupType
	SongGroupType
)

func NewGroupType(str string) GroupType {
	switch str {
	case "studio":
		return StudioGroupType
	case "song":
		return SongGroupType
	default:
		return ErrGroupType
	}
}

type Invitation struct {
	Email      string `json:"email" db:"email"`
	GroupID    string `json:"groupId" db:"group_id"`
	GroupType  `json:"groupType" db:"group_type"`
	IsAccepted bool `json:"isAccepted" db:"is_accepted"`
}

type InvitationUsecase interface {
	Filter(email, groupID string) ([]*Invitation, error)
	Create(email, groupID string, groupType GroupType) error
	Accept(email, groupID string) error
}

type InvitationRepository interface {
	GetByIDs(email, groupID string) (*Invitation, error)
	Filter(email, groupID string) ([]*Invitation, error)
	Create(email, groupID string, groupType GroupType) error
	Accept(email, groupID string) error
}
