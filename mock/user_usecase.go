package mock

import comusic "github.com/sabigara/comusicAPI"

type UserInteractor struct {
	CreateRet     func() (*comusic.User, error)
	CreateInvoked bool

	GetRet     func() (*comusic.User, error)
	GetInvoked bool
}

func (ui *UserInteractor) Create(name, email, password string) (*comusic.User, error) {
	ui.CreateInvoked = true
	return ui.CreateRet()
}

func (ui *UserInteractor) GetById(id string) (*comusic.User, error) {
	ui.GetInvoked = true
	return ui.GetRet()
}
