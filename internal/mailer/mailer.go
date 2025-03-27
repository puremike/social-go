package mailer

import (
	"bytes"
	"embed"
	"html/template"
)

const (
	FromName            = "SocialGo"
	MaxRetries          = 3
	WelcomeUserTemplate = "user_invitation.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	// Send(string, string, string, any, bool) error
	SendMailTrap(string, string, string, any, bool) (int, error)
}

func parseSendEmail(templateFile string, data any) (string, string, error) {
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return "", "", err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return "", "", err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return "", "", err
	}
	return subject.String(), body.String(), nil
}
