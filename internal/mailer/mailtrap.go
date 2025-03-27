package mailer

import (
	"crypto/tls"
	"errors"

	gomail "gopkg.in/mail.v2"
)

type MailTrap struct {
	fromEmail, apiKey string
}

func NewMailTrapMailer(fromEmail, apiKey string) (MailTrap, error) {
	if apiKey == "" {
		return MailTrap{}, errors.New("api key is required")
	}

	return MailTrap{
		fromEmail: fromEmail,
		apiKey:    apiKey,
	}, nil
}

func (m MailTrap) SendMailTrap(templateFile, username, email string, data any, isSandBox bool) (int, error) {

	// parse template

	subject, body, err := parseSendEmail(templateFile, data)
	if err != nil {
		return -1, err
	}
	// send welcome email to user

	message := gomail.NewMessage()
	message.SetHeader("From", m.fromEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject)

	message.AddAlternative("text/html", body)

	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "smtp@mailtrap.io", m.apiKey)

	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true} // Debug TLS issues

	// Enable Debug Logging
	dialer.SSL = false
	dialer.LocalName = "localhost"

	if err := dialer.DialAndSend(message); err != nil {
		return -1, err
	}

	// initiate max retries

	// for i := range MaxRetries {
	// 	err := dialer.DialAndSend(message)
	// 	if err != nil {
	// 		log.Printf("Failed to send email, attempt %d of %d", i+1, MaxRetries)
	// 		log.Printf("Error: %v\n", err)

	// 		// exponential backoff
	// 		time.Sleep(time.Second * time.Duration(i+1))
	// 		continue
	// 	}
	// }

	return 200, nil

}
