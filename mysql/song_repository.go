package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	comusic "github.com/sabigara/comusicAPI"
)

type SongRepository struct {
	*sqlx.DB
}

func NewSongRepository(db *sqlx.DB) *SongRepository {
	return &SongRepository{DB: db}
}

func (r *SongRepository) Create(song *comusic.Song) error {
	_, err := r.Exec(
		`INSERT INTO songs (id, studio_id, created_at, updated_at, name)
		 VALUES (?, ?, ?, ?, ?)`,
		song.ID, song.StudioID, song.CreatedAt, song.UpdatedAt, song.Name,
	)
	if err != nil {
		return fmt.Errorf("mysql.song_repository.Create: %w", err)
	}
	return nil
}

type nullableVersion struct {
	ID        sql.NullString
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	SongID    sql.NullString
	Name      sql.NullString
}

func (nv *nullableVersion) valid() bool {
	return nv.ID.Valid && nv.CreatedAt.Valid &&
		nv.UpdatedAt.Valid && nv.SongID.Valid && nv.Name.Valid
}

func (nv *nullableVersion) convert() *comusic.Version {
	return &comusic.Version{
		Meta: &comusic.Meta{
			ID:        nv.ID.String,
			CreatedAt: nv.CreatedAt.Time,
			UpdatedAt: nv.CreatedAt.Time,
		},
		SongID: nv.SongID.String,
		Name:   nv.Name.String,
	}
}

func (r *SongRepository) FilterByStudioIDWithVersions(studioID string) (comusic.SongVerMap, error) {
	dict := make(comusic.SongVerMap)
	rows, err := r.Query(
		`SELECT
		 songs.id, songs.studio_id, songs.created_at, songs.updated_at, songs.name,
		 vers.id, vers.song_id, vers.created_at, vers.updated_at, vers.name
		 FROM songs
		 LEFT OUTER JOIN versions AS vers ON songs.id = vers.song_id
		 WHERE songs.studio_id = ?`,
		studioID,
	)
	if err != nil {
		return nil, fmt.Errorf("mysql.song_repository.FilterByStudioIDWithVersions: %w", err)
	}
	for rows.Next() {
		s := &comusic.Song{Meta: &comusic.Meta{}}
		v := &nullableVersion{}
		err := rows.Scan(
			&s.ID, &s.StudioID, &s.CreatedAt, &s.UpdatedAt, &s.Name,
			&v.ID, &v.SongID, &v.CreatedAt, &v.UpdatedAt, &v.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("mysql.song_repository.FilterByStudioIDWithVersions: %w", err)
		}
		if val, ok := dict[s.ID]; !ok {
			dict[s.ID] = &comusic.SongVer{
				Data: s,
			}
			if v.valid() {
				dict[s.ID].Versions = []*comusic.Version{v.convert()}
			}
		} else {
			if v.valid() {
				val.Versions = append(val.Versions, v.convert())
			}
		}
	}
	return dict, nil
}
