package interactor

import (
	"fmt"
	"time"

	comusic "github.com/sabigara/comusicAPI"
)

type TakeUsecase struct {
	comusic.TrackUsecase
	comusic.TakeRepository
	comusic.FileRepository
}

func NewTakeUsecase(
	tru comusic.TrackUsecase,
	tr comusic.TakeRepository,
	fr comusic.FileRepository,
) *TakeUsecase {
	return &TakeUsecase{
		TrackUsecase:   tru,
		TakeRepository: tr,
		FileRepository: fr,
	}
}

func (u *TakeUsecase) Create(trackID, name string, src comusic.FileSrc) (*comusic.Take, *comusic.File, error) {
	take := comusic.NewTake(trackID, name)
	file := comusic.NewFile(take.ID, src)
	take.FileID = file.ID
	file, err := u.FileRepository.Upload(file)
	if err != nil {
		return nil, nil, fmt.Errorf("interactor.take_usecase.Create: %w", err)
	}
	if err := u.TakeRepository.Create(take); err != nil {
		return nil, nil, fmt.Errorf("interactor.take_usecase.Create: %w", err)
	}
	err = u.TrackUsecase.Update(
		&comusic.TrackUpdateInput{
			ID:         trackID,
			UpdatedAt:  time.Now().UTC(),
			ActiveTake: &take.ID,
		},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("interactor.take_usecase.Create: %w", err)
	}
	return take, file, nil
}

func (u *TakeUsecase) GetByID(id string) (*comusic.Take, error) {
	tr, err := u.TakeRepository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("interactor.take_usecase.GetByID: %w", err)
	}
	return tr, nil
}

func (u *TakeUsecase) FilterByTrackID(trackID string) ([]*comusic.Take, error) {
	takes, err := u.TakeRepository.FilterByTrackID(trackID)
	if err != nil {
		return nil, fmt.Errorf("interactor.take_usecase.FilterByTrackID: %w", err)
	}
	return takes, nil
}

func (u *TakeUsecase) Delete(takeID string) error {
	tk, err := u.GetByID(takeID)
	if err != nil {
		return fmt.Errorf("interactor.take_usecase.Delete: %w", err)
	}
	err = u.TakeRepository.Delete(takeID)
	if err != nil {
		return fmt.Errorf("interactor.take_usecase.Delete: %w", err)
	}
	return u.FileRepository.Delete(tk.FileID)
}
