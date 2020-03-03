package mock

import (
	"fmt"
)

type MailUsecase struct {
}

func NewMailUsecase() *MailUsecase {
	return &MailUsecase{}
}

func (u *MailUsecase) InviteToStudioNew(to, studio_name, signupURL string) error {
	fmt.Printf("\n---------- mail start ----------\n<InviteToStudioNew>\nto: %s\nstudio_name: %s\nsignup_url: %s\n----------  mail end  ----------\n\n", to, studio_name, signupURL)
	return nil
}
