package sendgrid

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailUsecase struct {
	apiKey string
}

func NewMailUsecase() *MailUsecase {
	key, ok := os.LookupEnv("SENDGRID_API_KEY")
	if !ok {
		panic("SENDGRID_API_KEY is not set")
	}
	return &MailUsecase{
		apiKey: key,
	}
}

type TemplateData = map[string]string

func (u *MailUsecase) sendTemplate(to, templateID string, templateData TemplateData) error {
	m := mail.NewV3Mail()
	address := "noreply@comusic.com"
	name := "comusic"
	e := mail.NewEmail(name, address)
	m.SetFrom(e)
	m.SetTemplateID(templateID)

	p := mail.NewPersonalization()
	for k, v := range templateData {
		p.SetDynamicTemplateData(k, v)
	}
	p.AddTos(mail.NewEmail("To User", to))
	m.AddPersonalizations(p)

	req := sendgrid.GetRequest(u.apiKey, "/v3/mail/send", "https://api.sendgrid.com")
	req.Method = "POST"
	req.Body = mail.GetRequestBody(m)
	_, err := sendgrid.API(req)
	if err != nil {
		return fmt.Errorf("sendgrid.SendTemplate: %w", err)
	}

	return nil
}

func (u *MailUsecase) InviteToStudioNew(to, studio_name, signupURL string) error {
	return u.sendTemplate(to, "d-c0d11e8aede54caf9728b8fd11b3b09e", TemplateData{
		"studio_name": studio_name,
		"signup_url":  signupURL,
	})
}
