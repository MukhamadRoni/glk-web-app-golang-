package pelamar

import (
	// "fmt"
	"glk-web-app/config"
	"glk-web-app/models"
	"glk-web-app/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ShowLogin renders the pelamar login page.
func ShowLogin(c *fiber.Ctx) error {
	return c.Render("pelamar/login", fiber.Map{
		"Title":       "Masuk",
		"Description": "Login portal pelamar GLK",
	}, "layouts/auth")
}

// ProcessLogin handles the pelamar login form submission by sending a magic link.
func ProcessLogin(c *fiber.Ctx) error {
	email := c.FormValue("email")
	if email == "" {
		return c.Render("pelamar/login", fiber.Map{
			"Title": "Masuk",
			"Error": "Email wajib diisi.",
		}, "layouts/auth")
	}

	token, err := utils.GenerateMagicLinkToken(email)
	if err != nil {
		return c.Render("pelamar/login", fiber.Map{
			"Title": "Masuk",
			"Error": "Gagal membuat link login.",
		}, "layouts/auth")
	}

	// FOR DEVELOPMENT: Bypass email sending and redirect directly to magic link
	// baseURL := config.GetEnv("APP_URL", fmt.Sprintf("http://localhost:%s", config.GetEnv("PELAMAR_PORT", "8081")))
	// magicLink := fmt.Sprintf("%s/magic-link?token=%s", baseURL, token)

	// err = utils.SendMagicLinkEmail(email, magicLink)
	// if err != nil {
	// 	fmt.Println("Error sending email:", err)
	// 	return c.Render("pelamar/login", fiber.Map{
	// 		"Title": "Masuk",
	// 		"Error": "Gagal mengirim email, periksa konfigurasi SMTP Anda.",
	// 	}, "layouts/auth")
	// }

	// return c.Render("pelamar/login_sent", fiber.Map{
	// 	"Title": "Cek Email Anda",
	// 	"Email": email,
	// }, "layouts/auth")

	return c.Redirect("/magic-link?token=" + token)
}

// VerifyMagicLink handles the callback from the email link.
func VerifyMagicLink(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Redirect("/login?error=Token tidak valid")
	}

	email, err := utils.VerifyMagicLinkToken(token)
	if err != nil {
		return c.Redirect("/login?error=Link telah kedaluwarsa atau tidak valid")
	}

	pelamar, err := models.GetPelamarByEmail(config.DB, email)
	if err != nil {
		// Auto-register pelamar if not exists
		namePart := strings.Split(email, "@")[0]
		pelamar = &models.Pelamar{
			Name:  namePart,
			Email: email,
		}
		if err := models.CreatePelamar(config.DB, pelamar); err != nil {
			return c.Redirect("/login?error=Gagal membuat akun otomatis")
		}
	}

	// Store session cookie
	sess, err := store.Get(c)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	sess.Set("pelamar_id", pelamar.ID)
	sess.Set("pelamar_name", pelamar.Name)
	if err := sess.Save(); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Redirect("/dashboard")
}

// Logout clears the pelamar session and redirects to login.
func Logout(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err == nil {
		_ = sess.Destroy()
	}
	return c.Redirect("/login")
}
