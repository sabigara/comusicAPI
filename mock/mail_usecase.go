package mock

import (
	"fmt"

	comusic "github.com/sabigara/comusicAPI"
)

type MailUsecase struct {
}

func NewMailUsecase() *MailUsecase {
	return &MailUsecase{}
}

func (u *MailUsecase) Print(content string) {
	fmt.Printf("\n---------- mail start ----------\n%s\n----------  mail end  ----------\n\n", content)
}

func (u *MailUsecase) InviteToStudioNew(to, studio_name, signupURL string) error {
	u.Print(fmt.Sprintf("<InviteToStudioNew>\nto: %s\nstudio_name: %s\nsignup_url: %s\n", to, studio_name, signupURL))
	return nil
}

func (u *MailUsecase) InviteToStudio(user *comusic.User, studio_name, signInURL string) error {
	u.Print(fmt.Sprintf("<InviteToStudio>\nuser: %s\nstudio_name: %s\nsignin_url: %s", user.Email, studio_name, signInURL))
	return nil
}
