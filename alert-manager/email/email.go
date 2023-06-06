package email

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"text/template"
)

func SendEmail(to []string, subject string, server string, error string, date string, templatePath string) {
	from := "dev.vyctor@gmail.com"
	password := os.Getenv("GMAIL_PASSWORD")

	if password == "" {
		panic("GMAIL_PASSWORD environment variable is not set")
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)
	t, _ := template.ParseFiles(templatePath)

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", subject, mimeHeaders)))
	t.Execute(&body, struct {
		Server string
		Error  string
		Date   string
	}{
		Server: server,
		Error:  error,
		Date:   date,
	})

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())

	if err != nil {
		fmt.Printf("Erro ao enviar email: %s", err)
		return
	}

	fmt.Println("Email enviado com sucesso!")
}
