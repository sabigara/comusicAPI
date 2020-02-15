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

func (tu *TakeUsecase) Create(trackID, name string, src comusic.FileSrc) (*comusic.Take, *comusic.File, error) {
	take := comusic.NewTake(trackID, name)
	file := comusic.NewFile(take.ID, src)
	take.FileID = file.ID
	file, err := tu.FileRepository.Upload(file)
	if err != nil {
		return nil, nil, fmt.Errorf("interactor.take_usecase.Create: %w", err)
	}
	if err := tu.TakeRepository.Create(take); err != nil {
		return nil, nil, fmt.Errorf("interactor.take_usecase.Create: %w", err)
	}
	err = tu.TrackUsecase.Update(
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

func (tu *TakeUsecase) GetByID(id string) (*comusic.Take, error) {
	tr, err := tu.TakeRepository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("interactor.take_usecase.GetByID: %w", err)
	}
	return tr, nil
}

func (tu *TakeUsecase) Delete(takeID string) error {
	tk, err := tu.GetByID(takeID)
	if err != nil {
		return fmt.Errorf("interactor.take_usecase.Delete: %w", err)
	}
	err = tu.TakeRepository.Delete(takeID)
	if err != nil {
		return fmt.Errorf("interactor.take_usecase.Delete: %w", err)
	}
	return tu.FileRepository.Delete(tk.FileID)
}
