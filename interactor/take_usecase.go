package interactor

import (
	"fmt"

	comusic "github.com/sabigara/comusicAPI"
)

type TakeUsecase struct {
	comusic.TakeRepository
	comusic.FileRepository
}

func NewTakeUsecase(tr comusic.TakeRepository, fr comusic.FileRepository) *TakeUsecase {
	return &TakeUsecase{
		TakeRepository: tr,
		FileRepository: fr,
	}
}

func (tu *TakeUsecase) Create(trackID, name string, src comusic.FileSrc) (*comusic.Take, error) {
	take := comusic.NewTake(trackID, name)
	file := comusic.NewFile(take.ID, src)
	if err := tu.FileRepository.Upload(file); err != nil {
		return nil, fmt.Errorf("interactor.take_usecase.Create: %w", err)
	}
	if err := tu.TakeRepository.Create(take); err != nil {
		return nil, fmt.Errorf("interactor.take_usecase.Create: %w", err)
	}
	return take, nil
}
