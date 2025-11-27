package email

import (
	"fmt"
	"net/smtp"
	"os"
	"strconv"
)

type SMTPSender struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

func NewSMTPSender() *SMTPSender {
	host := os.Getenv("SMTP_HOST")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	portStr := os.Getenv("SMTP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil || port == 0 {
		port = 587
	}

	return &SMTPSender{
		Host:     host,
		Port:     port,
		User:     user,
		Password: pass,
		From:     from,
	}
}

func (s *SMTPSender) SendEmail(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)

	msg := "MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		fmt.Sprintf("From: %s\r\n", s.From) +
		fmt.Sprintf("To: %s\r\n", to) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"\r\n" +
		body

	auth := smtp.PlainAuth("", s.User, s.Password, s.Host)

	return smtp.SendMail(
		addr,
		auth,
		s.User,
		[]string{to},
		[]byte(msg),
	)
}
