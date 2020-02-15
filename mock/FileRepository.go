package mock

import (
	"io"
	"os"
	"path/filepath"

	comusic "github.com/sabigara/comusicAPI"
)

type FileRepository struct{}

func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

var UPLOADS_DIR string
var UPLOADS_PATH string

func (r *FileRepository) Upload(f *comusic.File) (*comusic.File, error) {
	url := filepath.Join(UPLOADS_DIR, f.ID)
	fpath := filepath.Base(url)
	dpath := filepath.Dir(url)
	os.MkdirAll(dpath, os.ModePerm)
	dst, err := os.Create(filepath.Join(dpath, fpath))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(dst, f.Src)
	// URL to serve
	f.URL = r.URL(f.ID)
	return f, err
}

func (r *FileRepository) Download(url string) (*comusic.File, error) {
	return nil, nil
}

func (r *FileRepository) Delete(id string) error {
	return os.Remove(filepath.Join(UPLOADS_DIR, id))
}

func (r *FileRepository) URL(fileID string) string {
	return UPLOADS_PATH + "/" + fileID
}

func init() {
	uploadsDir, ok := os.LookupEnv("UPLOADS_DIR")
	if !ok {
		panic("UPLOADS_DIR not specified")
	}
	uploadsPath, ok := os.LookupEnv("UPLOADS_PATH")
	if !ok {
		panic("UPLOADS_PATH not specified")
	}
	UPLOADS_DIR = uploadsDir
	UPLOADS_PATH = uploadsPath
}
