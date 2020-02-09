package interactor

import (
	"fmt"

	comusic "github.com/sabigara/comusicAPI"
)

type TakeUsecase struct {
	comusic.TakeRepository
}

func NewTakeUsecase(tr comusic.TakeRepository) *TakeUsecase {
	return &TakeUsecase{TakeRepository: tr}
}

func (su *TakeUsecase) Create(trackID, name string) (*comusic.Take, error) {
	take := comusic.NewTake(trackID, name)
	err := su.TakeRepository.Create(take)
	if err != nil {
		return nil, fmt.Errorf("interactor.take_usecase.Create: %w", err)
	}
	return take, nil
}
