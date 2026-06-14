package utils

import (
	"crypto/tls"
	"fmt"
	"glk-web-app/config"
	"log"
	"net/smtp"
)

// SendMagicLinkEmail sends an email to the candidate with the magic login link.
func SendMagicLinkEmail(toEmail, magicLink string) error {
	host := config.GetEnv("SMTP_HOST", "smtp.gmail.com")
	port := config.GetEnv("SMTP_PORT", "587")
	user := config.GetEnv("SMTP_USER", "")
	pass := config.GetEnv("SMTP_PASS", "")

	if user == "" || pass == "" {
		log.Println("[EMAIL] SMTP credentials are not set in .env")
		return fmt.Errorf("SMTP credentials are not set in .env")
	}

	auth := smtp.PlainAuth("", user, pass, host)

	subject := "Masuk ke Portal Pelamar - Gurulesku"

	// Menggunakan desain modern, minimalis, dan clean
	body := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<style>
			/* Efek interaktif untuk email client yang mendukung */
			.btn-login:hover {
				background-color: #E05D2D !important;
				box-shadow: 0 4px 12px rgba(248, 114, 66, 0.3) !important;
			}
			@media only screen and (max-width: 600px) {
				.main-container {
					width: 100%% !important;
					border-radius: 0 !important;
				}
				.content-padding {
					padding: 30px 20px !important;
				}
			}
		</style>
	</head>
	<body style="margin: 0; padding: 0; font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; background-color: #F8FAFC; -webkit-font-smoothing: antialiased;">
		<table border="0" cellpadding="0" cellspacing="0" width="100%%" style="background-color: #F8FAFC; padding: 40px 15px;">
			<tr>
				<td align="center">
					<table class="main-container" border="0" cellpadding="0" cellspacing="0" width="600" style="background-color: #ffffff; border-radius: 12px; overflow: hidden; box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.05), 0 8px 10px -6px rgba(0, 0, 0, 0.01);">

						<tr>
							<td align="center" style="padding: 40px 0 30px 0; background-color: #ffffff;">
								<div style="width: 60px; height: 6px; background-color: #F87242; border-radius: 3px; margin-bottom: 20px;"></div>
								<h2 style="margin: 0; color: #1E293B; font-weight: 700; font-size: 24px; letter-spacing: -0.5px;">Portal Pelamar</h2>
								<p style="margin: 5px 0 0 0; color: #64748B; font-size: 14px; font-weight: 500;">PT Gurulesku Nusantara</p>
							</td>
						</tr>

						<tr>
							<td class="content-padding" style="padding: 0 40px 30px 40px;">
								<table border="0" cellpadding="0" cellspacing="0" width="100%%">
									<tr>
										<td style="color: #334155; font-size: 16px; line-height: 1.6;">
											<p style="margin-top: 0; font-weight: 600; color: #0F172A; font-size: 18px;">Halo,</p>
											<p>Kami menerima permintaan untuk masuk ke akun Anda di portal rekrutmen kami. Silakan klik tombol di bawah ini untuk melanjutkan.</p>
										</td>
									</tr>

									<tr>
										<td align="center" style="padding: 35px 0;">
											<a href="%s" class="btn-login" style="background-color: #F87242; color: #ffffff; padding: 16px 40px; text-decoration: none; border-radius: 8px; font-weight: 600; font-size: 16px; display: inline-block; transition: all 0.2s ease;">Masuk ke Akun Saya</a>
										</td>
									</tr>

									<tr>
										<td style="background-color: #F1F5F9; border-radius: 8px; padding: 20px;">
											<p style="margin: 0 0 10px 0; color: #64748B; font-size: 13px; font-weight: 600;">Jika tombol di atas tidak berfungsi, salin link berikut:</p>
											<p style="margin: 0; word-break: break-all;"><a href="%s" style="color: #3B82F6; text-decoration: underline; font-size: 13px;">%s</a></p>
										</td>
									</tr>

									<tr>
										<td style="padding-top: 30px;">
											<p style="margin: 0; color: #94A3B8; font-size: 13px; line-height: 1.5;">*Link ajaib ini hanya berlaku selama 15 menit dan hanya dapat digunakan satu kali demi keamanan data Anda.</p>
										</td>
									</tr>
								</table>
							</td>
						</tr>

						<tr>
							<td align="center" style="padding: 25px 40px; background-color: #F8FAFC; border-top: 1px solid #E2E8F0;">
								<p style="margin: 0; color: #94A3B8; font-size: 12px; font-weight: 500;">&copy; 2026 PT Gurulesku Nusantara. Hak cipta dilindungi.</p>
							</td>
						</tr>
					</table>
				</td>
			</tr>
		</table>
	</body>
	</html>
	`, magicLink, magicLink, magicLink)

	msg := []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		body)

	addr := host + ":" + port

	// Logic pengiriman email tetap sama
	err := smtp.SendMail(addr, auth, user, []string{toEmail}, msg)
	if err != nil {
		log.Printf("[EMAIL ERROR] Standard SendMail failed: %v", err)
		if port == "465" {
			tlsconfig := &tls.Config{
				InsecureSkipVerify: false,
				ServerName:         host,
			}
			conn, err := tls.Dial("tcp", addr, tlsconfig)
			if err != nil {
				log.Printf("[EMAIL ERROR] TLS Dial failed: %v", err)
				return err
			}
			c, err := smtp.NewClient(conn, host)
			if err != nil {
				log.Printf("[EMAIL ERROR] NewClient failed: %v", err)
				return err
			}
			if err = c.Auth(auth); err != nil {
				log.Printf("[EMAIL ERROR] SMTP Auth failed: %v", err)
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
