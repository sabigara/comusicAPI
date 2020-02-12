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

func (r *FileRepository) Upload(f *comusic.File) error {
	url := filepath.Join(UPLOADS_DIR, f.URL)
	fpath := filepath.Base(url)
	dpath := filepath.Dir(url)
	os.MkdirAll(dpath, os.ModePerm)
	dst, err := os.Create(filepath.Join(dpath, fpath))
	if err != nil {
		return err
	}
	_, err = io.Copy(dst, f.Src)
	return err
}

func (r *FileRepository) Download(url string) (*comusic.File, error) {
	return nil, nil
}

func init() {
	uploads_dir, ok := os.LookupEnv("UPLOADS_DIR")
	if !ok {
		panic("UPLOADS_DIR not specified")
	}
	UPLOADS_DIR = uploads_dir
}
