package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGridMailer(fromEmail, apiKey string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (s *SendGridMailer) Send(templateFile, username, email string, data any, isSandBox bool) error {
	from := mail.NewEmail(FromName, s.fromEmail)
	to := mail.NewEmail(username, email)

	// parse template
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}

	// send welcome email to user

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	// set mail setting to allow email to be sent/not sent in development or production environment
	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandBox,
		},
	})

	// initiate max retries

	for i := range MaxRetries {
		response, err := s.client.Send(message)
		if err != nil {
			log.Printf("Failed to send email to %v, attempt %d of %d", email, i+1, MaxRetries)
			log.Printf("Error: %v", err.Error())
		}
		log.Printf("Email sent with statusCode %v", response.StatusCode)
	}

	return fmt.Errorf("error sending an email to: %v", to)
}
