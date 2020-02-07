package mock

import (
	comusic "github.com/sabigara/comusicAPI"
)

type UserRepository struct {
	SaveRet     func() error
	SaveInvoked bool
	GetRet      func() (*comusic.User, error)
	GetInvoked  bool
}

func (ur *UserRepository) Save(user *comusic.User) error {
	ur.SaveInvoked = true
	return ur.SaveRet()
}

func (ur *UserRepository) GetById(id string) (*comusic.User, error) {
	ur.GetInvoked = true
	return ur.GetRet()
}

func (ur *UserRepository) GetByEmail(email string) (*comusic.User, error) {
	return ur.GetRet()
}
