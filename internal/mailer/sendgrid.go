package mailer

import (
	"fmt"
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
	subject, body, err := parseSendEmail(templateFile, data)
	if err != nil {
		return err
	}

	// send welcome email to user

	message := mail.NewSingleEmail(from, subject, to, "", body)

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
			log.Printf("Failed to send email, attempt %d of %d", i+1, MaxRetries)
			log.Printf("Error: %v", err.Error())

			// exponential backoff
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		log.Printf("Email sent with statusCode %v", response.StatusCode)
		return nil
	}

	return fmt.Errorf("error sending an email to: %v", to)
}
