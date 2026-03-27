package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/mail"
	"net/smtp"

	"github.com/yourusername/navhub/internal/config"
)

// EmailService handles email operations
type EmailService struct {
	from     string
	smtpHost string
	smtpPort int
	username string
	password string
}

// NewEmailService creates a new email service
func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		from:     cfg.SMTPFrom,
		smtpHost: cfg.SMTPHost,
		smtpPort: cfg.SMTPPort,
		username: cfg.SMTPUser,
		password: cfg.SMTPPassword,
	}
}

// SendVerificationEmail sends an email verification link
func (s *EmailService) SendVerificationEmail(email, token string) error {
	// TODO: Implement actual email sending
	fmt.Printf("📧 [EMAIL] Verification email sent to %s with token: %s\n", email, token)
	return nil
}

// SendPasswordResetEmail sends a password reset link
func (s *EmailService) SendPasswordResetEmail(email, token string) error {
	// Create the reset URL
	resetURL := fmt.Sprintf("http://localhost:5173/reset-password?token=%s", token)

	// Compose email
	subject := "重置您的 NavHub 密码"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>重置密码</title>
</head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="background-color: #f4f4f4; padding: 30px; border-radius: 5px;">
        <h2 style="color: #333;">重置您的密码</h2>
        <p>您好，</p>
        <p>我们收到了您的密码重置请求。请点击下面的按钮重置您的密码：</p>

        <div style="text-align: center; margin: 30px 0;">
            <a href="%s"
               style="background-color: #4CAF50; color: white; padding: 15px 30px;
                      text-decoration: none; border-radius: 5px; display: inline-block;
                      font-weight: bold;">
                重置密码
            </a>
        </div>

        <p>或者复制以下链接到浏览器：</p>
        <p style="background-color: #fff; padding: 10px; border-radius: 3px; word-break: break-all;">
            %s
        </p>

        <p style="color: #666; font-size: 12px;">
            此链接将在1小时后过期。如果您没有请求重置密码，请忽略此邮件。
        </p>

        <hr style="border: none; border-top: 1px solid #ddd; margin: 20px 0;">

        <p style="color: #999; font-size: 12px;">
            © 2024 NavHub. All rights reserved.
        </p>
    </div>
</body>
</html>
`, resetURL, resetURL)

	// Send email
	return s.send(email, subject, body)
}

// send sends an email using SMTP
func (s *EmailService) send(to, subject, body string) error {
	// Validate email address
	if _, err := mail.ParseAddress(to); err != nil {
		return fmt.Errorf("invalid email address: %w", err)
	}

	// Format SMTP address
	smtpAddr := fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort)

	// Prepare email content
	var auth smtp.Auth
	if s.username != "" && s.password != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.smtpHost)
	}

	// Build email headers and body
	fromAddr := s.from
	if _, err := mail.ParseAddress(fromAddr); err != nil {
		return fmt.Errorf("invalid from address: %w", err)
	}

	headers := make(map[string]string)
	headers["From"] = fromAddr
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Build message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Send email
	var err error
	if auth != nil {
		err = smtp.SendMail(smtpAddr, auth, fromAddr, []string{to}, []byte(message))
	} else {
		err = smtp.SendMail(smtpAddr, nil, fromAddr, []string{to}, []byte(message))
	}

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Printf("📧 [EMAIL] Email sent successfully to %s\n", to)
	return nil
}

// GenerateToken generates a secure random token
func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
