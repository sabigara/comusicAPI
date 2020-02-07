package interactor

import (
	"fmt"

	comusic "github.com/sabigara/comusicAPI"
)

type UserInteractor struct {
	comusic.UserRepository
}

func NewUserInteractor(userRepository comusic.UserRepository) *UserInteractor {
	return &UserInteractor{UserRepository: userRepository}
}

func (ui *UserInteractor) Create(name, email, password string) (*comusic.User, error) {
	user := comusic.NewUser("", name, email)
	user.Password = password
	err := ui.UserRepository.Save(user)
	if err != nil {
		return nil, fmt.Errorf("interactor.user_usecase.Create: %w", err)
	}
	return user, nil
}

func (ui *UserInteractor) Update(id, name string) (*comusic.User, error) {
	user, err := ui.UserRepository.GetById(id)
	if err != nil {
		return nil, fmt.Errorf("interactor.user_usecase.Update: %w", err)
	}
	user.Name = name
	if err := ui.Save(user); err != nil {
		return nil, fmt.Errorf("interactor.user_usecase.Update: %w", err)
	}
	return user, nil
}

func (ui *UserInteractor) GetById(id string) (*comusic.User, error) {
	user, err := ui.UserRepository.GetById(id)
	if err != nil {
		return nil, fmt.Errorf("interactor.user_usecase.Get: %w", err)
	}
	return user, err
}

func (ui *UserInteractor) GetByEmail(email string) (*comusic.User, error) {
	user, err := ui.UserRepository.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("interactor.user_usecase.Get: %w", err)
	}
	return user, err
}
