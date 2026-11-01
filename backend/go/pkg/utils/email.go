package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
)

// SendPasswordResetEmail envia um código de 6 dígitos por SMTP para redefinição de senha.
// Caso SMTP_HOST não esteja configurado (ex.: ambiente local de desenvolvimento),
// registra o código no log do servidor para facilitar testes.
func SendPasswordResetEmail(toEmail, userName, code string) error {
	host := strings.TrimSpace(os.Getenv("SMTP_HOST"))
	port := strings.TrimSpace(os.Getenv("SMTP_PORT"))
	if port == "" {
		port = "587"
	}
	username := strings.TrimSpace(os.Getenv("SMTP_USERNAME"))
	password := strings.TrimSpace(os.Getenv("SMTP_PASSWORD"))
	from := strings.TrimSpace(os.Getenv("SMTP_FROM"))
	if from == "" {
		from = username
	}
	if from == "" {
		from = "no-reply@speakup.local"
	}

	if host == "" {
		log.Printf("[SMTP não configurado] Código de recuperação para %s (%s): %s", toEmail, userName, code)
		return nil
	}

	displayName := strings.TrimSpace(userName)
	if displayName == "" {
		displayName = "Estudante"
	}

	subject := "SpeakUp - Código de Recuperação de Senha"
	body := fmt.Sprintf(
		"Olá, %s!\r\n\r\n"+
			"Recebemos uma solicitação para redefinir a senha da sua conta no SpeakUp.\r\n\r\n"+
			"Seu código de verificação é: %s\r\n\r\n"+
			"Este código expira em 15 minutos.\r\n"+
			"Se você não solicitou essa alteração, ignore este e-mail.\r\n\r\n"+
			"Equipe SpeakUp",
		displayName,
		code,
	)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n"+
			"%s",
		from,
		toEmail,
		subject,
		body,
	))

	addr := fmt.Sprintf("%s:%s", host, port)
	var auth smtp.Auth
	if username != "" && password != "" {
		auth = smtp.PlainAuth("", username, password, host)
	}

	if err := smtp.SendMail(addr, auth, from, []string{toEmail}, msg); err != nil {
		return fmt.Errorf("falha ao enviar email SMTP: %w", err)
	}

	return nil
}

