package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendgrid(apiKey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return nil
	}

	subject := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return nil
	}

	body := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return nil
	}
	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())
	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	var retryErr error
	for i := 0; i < maxRetries; i++ {
		response, retryErr := m.client.Send(message)
		if err != nil {

			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		if response.StatusCode >= 200 && response.StatusCode < 300 {
			log.Printf("Email sent to %v with status code %d", email, response.StatusCode)
			return nil
		}

		log.Printf("Failed to send email to %v, attempt %d of %d", email, i+1, maxRetries)

		time.Sleep(time.Second * time.Duration(i+1))
	}

	return fmt.Errorf("failed to send email after %d attempts, error: %v", maxRetries, retryErr)
}
