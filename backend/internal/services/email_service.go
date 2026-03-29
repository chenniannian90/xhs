package services

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net/mail"
	"net/smtp"

	"github.com/yourusername/navhub/internal/config"
)

// EmailService handles email operations
type EmailService struct {
	from         string
	smtpHost     string
	smtpPort     int
	username     string
	password     string
	frontendURL  string
}

// NewEmailService creates a new email service
func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		from:        cfg.SMTPFrom,
		smtpHost:    cfg.SMTPHost,
		smtpPort:    cfg.SMTPPort,
		username:    cfg.SMTPUser,
		password:    cfg.SMTPPassword,
		frontendURL: cfg.FrontendURL,
	}
}

// SendVerificationEmail sends an email verification link
func (s *EmailService) SendVerificationEmail(email, token string) error {
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", s.frontendURL, token)

	subject := "验证您的 NavHub 邮箱"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>邮箱验证</title>
</head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="background-color: #f4f4f4; padding: 30px; border-radius: 5px;">
        <h2 style="color: #333;">验证您的邮箱</h2>
        <p>您好，</p>
        <p>感谢您注册 NavHub！请点击下面的按钮验证您的邮箱地址：</p>

        <div style="text-align: center; margin: 30px 0;">
            <a href="%s"
               style="background-color: #4CAF50; color: white; padding: 15px 30px;
                      text-decoration: none; border-radius: 5px; display: inline-block;
                      font-weight: bold;">
                验证邮箱
            </a>
        </div>

        <p>或者复制以下链接到浏览器：</p>
        <p style="background-color: #fff; padding: 10px; border-radius: 3px; word-break: break-all;">
            %s
        </p>

        <p style="color: #666; font-size: 12px;">
            此链接将在24小时后过期。如果您没有注册 NavHub，请忽略此邮件。
        </p>

        <hr style="border: none; border-top: 1px solid #ddd; margin: 20px 0;">
        <p style="color: #999; font-size: 12px;">© 2025 NavHub. All rights reserved.</p>
    </div>
</body>
</html>
`, verifyURL, verifyURL)

	return s.send(email, subject, body)
}

// SendPasswordResetEmail sends a password reset link
func (s *EmailService) SendPasswordResetEmail(email, token string) error {
	// Create the reset URL
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.frontendURL, token)

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

	// Format SMTP address
	smtpAddr := fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort)

	// Use SSL/TLS for port 465, STARTTLS for port 587/25
	if s.smtpPort == 465 {
		return s.sendWithTLS(smtpAddr, fromAddr, to, []byte(message))
	}

	// Default: STARTTLS
	var auth smtp.Auth
	if s.username != "" && s.password != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.smtpHost)
	}

	if err := smtp.SendMail(smtpAddr, auth, fromAddr, []string{to}, []byte(message)); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Printf("📧 [EMAIL] Email sent successfully to %s\n", to)
	return nil
}

// sendWithTLS sends email using direct TLS connection (port 465)
func (s *EmailService) sendWithTLS(smtpAddr, from, to string, message []byte) error {
	tlsConfig := &tls.Config{
		ServerName: s.smtpHost,
	}

	conn, err := tls.Dial("tcp", smtpAddr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.smtpHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	// Authenticate
	auth := smtp.PlainAuth("", s.username, s.password, s.smtpHost)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}

	// Set sender
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipient
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send message body
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to start data: %w", err)
	}

	if _, err := w.Write(message); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close data: %w", err)
	}

	client.Quit()

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
