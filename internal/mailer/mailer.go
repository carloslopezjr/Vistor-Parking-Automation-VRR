package mailer

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
)

// Config holds SMTP configuration parameters.
type Config struct {
	Host string
	Port int
	User string
	Pass string
	From string
}

// Service sends emails using SMTP.
type Service interface {
	// Send sends a plain-text email to the given recipients.
	Send(ctx context.Context, to []string, subject, body string) error
}

// New returns a new SMTP-backed mailer service.
func New(cfg Config) Service {
	return &smtpService{cfg: cfg}
}

type smtpService struct {
	cfg Config
}

func (s *smtpService) Send(ctx context.Context, to []string, subject, body string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if len(to) == 0 {
		return fmt.Errorf("no recipients provided")
	}
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

	auth := smtp.PlainAuth("", s.cfg.User, s.cfg.Pass, s.cfg.Host)

	headers := map[string]string{
		"From":         s.cfg.From,
		"To":           strings.Join(to, ", "),
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/plain; charset=utf-8",
	}

	var sb strings.Builder
	for k, v := range headers {
		sb.WriteString(k)
		sb.WriteString(": ")
		sb.WriteString(v)
		sb.WriteString("\r\n")
	}
	sb.WriteString("\r\n")
	sb.WriteString(body)

	msg := []byte(sb.String())

	// net/smtp does not support context cancellation; we honor ctx upfront
	// and rely on the underlying network stack for timeouts.
	return smtp.SendMail(addr, auth, s.cfg.From, to, msg)
}
