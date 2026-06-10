package utils

import (
	"crypto/tls"
	"fmt"
	"glk-web-app/config"
	"net/smtp"
)

// SendMagicLinkEmail sends an email to the candidate with the magic login link.
func SendMagicLinkEmail(toEmail, magicLink string) error {
	host := config.GetEnv("SMTP_HOST", "smtp.gmail.com")
	port := config.GetEnv("SMTP_PORT", "587")
	user := config.GetEnv("SMTP_USER", "")
	pass := config.GetEnv("SMTP_PASS", "")

	if user == "" || pass == "" {
		return fmt.Errorf("SMTP credentials are not set in .env")
	}

	auth := smtp.PlainAuth("", user, pass, host)

	subject := "Login Portal Pelamar - GLK"
	body := fmt.Sprintf(`
	<html>
	<body>
		<h2>Halo Pelamar!</h2>
		<p>Anda telah meminta link untuk masuk ke Portal Pelamar GLK.</p>
		<p>Silakan klik tombol di bawah ini untuk masuk secara otomatis:</p>
		<a href="%s" style="display:inline-block;padding:10px 20px;background-color:#F87242;color:#ffffff;text-decoration:none;border-radius:5px;font-weight:bold;">Masuk ke Portal</a>
		<br><br>
		<p>Atau copy paste link berikut di browser Anda:</p>
		<p><a href="%s">%s</a></p>
		<p>Link ini hanya berlaku untuk sementara waktu.</p>
		<p>Terima kasih,</p>
		<p>Tim Rekrutmen GLK</p>
	</body>
	</html>
	`, magicLink, magicLink, magicLink)

	msg := []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		body)

	addr := host + ":" + port

	// We use tls for port 465, but for 587 we use starttls. 
	// smtp.SendMail handles STARTTLS automatically if the server supports it.
	err := smtp.SendMail(addr, auth, user, []string{toEmail}, msg)
	if err != nil {
		// fallback to TLS if port is 465
		if port == "465" {
			tlsconfig := &tls.Config{
				InsecureSkipVerify: false,
				ServerName:         host,
			}
			conn, err := tls.Dial("tcp", addr, tlsconfig)
			if err != nil {
				return err
			}
			c, err := smtp.NewClient(conn, host)
			if err != nil {
				return err
			}
			if err = c.Auth(auth); err != nil {
				return err
			}
			if err = c.Mail(user); err != nil {
				return err
			}
			if err = c.Rcpt(toEmail); err != nil {
				return err
			}
			w, err := c.Data()
			if err != nil {
				return err
			}
			_, err = w.Write(msg)
			if err != nil {
				return err
			}
			err = w.Close()
			if err != nil {
				return err
			}
			return c.Quit()
		}
		return err
	}
	return nil
}
