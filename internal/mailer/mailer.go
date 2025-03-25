package mailer

import "embed"

const (
	FromName            = "SocialGo"
	MaxRetries          = 3
	WelcomeUserTemplate = "user_invitation.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(string, string, string, any, bool) error
}
