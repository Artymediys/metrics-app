package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

type Server struct {
	host string
	port string
}

func NewServer(host, port string) *Server {
	return &Server{
		host: host,
		port: port,
	}
}

func (eServer *Server) Send(from string, to []string, subject string, body string) error {
	message := []byte("From: " + from + "\r\n" +
		"To: " + strings.Join(to, ", ") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		body + "\r\n")

	client, err := smtp.Dial(eServer.host + ":" + eServer.port)
	if err != nil {
		return fmt.Errorf("services.email.Send -> smtp.Dial: %w", err)
	}
	defer client.Close()
	defer client.Quit()

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("services.email.Send -> client.Mail: %w", err)
	}
	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("services.email.Send -> client.Rcpt: %w", err)
		}
	}

	dataWriter, err := client.Data()
	if err != nil {
		return fmt.Errorf("services.email.Send -> client.Data: %w", err)
	}

	if _, err = dataWriter.Write(message); err != nil {
		return fmt.Errorf("services.email.Send -> dataWriter.Write: %w", err)
	}

	if err = dataWriter.Close(); err != nil {
		return fmt.Errorf("services.email.Send -> dataWriter.Close: %w", err)
	}
	
	return nil
}
