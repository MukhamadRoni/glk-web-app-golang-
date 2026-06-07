package pelamar

import (
	"glk-web-app/config"
	"glk-web-app/models"

	"github.com/gofiber/fiber/v2"
)

// ShowLogin renders the pelamar login page.
func ShowLogin(c *fiber.Ctx) error {
	return c.Render("pelamar/login", fiber.Map{
		"Title":       "Masuk",
		"Description": "Login portal pelamar GLK",
	}, "layouts/auth")
}

// ProcessLogin handles the pelamar login form submission.
func ProcessLogin(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	pelamar, err := models.GetPelamarByEmail(config.DB, email)
	if err != nil || !pelamar.CheckPassword(password) {
		return c.Render("pelamar/login", fiber.Map{
			"Title": "Masuk",
			"Error": "Email atau password salah.",
		}, "layouts/auth")
	}

	// Store session cookie (simple approach — replace with JWT in Phase 3)
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

// ShowRegister renders the pelamar registration page.
func ShowRegister(c *fiber.Ctx) error {
	return c.Render("pelamar/register", fiber.Map{
		"Title":       "Daftar",
		"Description": "Buat akun pelamar GLK",
	}, "layouts/auth")
}

// ProcessRegister handles the pelamar registration form submission.
func ProcessRegister(c *fiber.Ctx) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	// Check if email already exists
	if _, err := models.GetPelamarByEmail(config.DB, email); err == nil {
		return c.Render("pelamar/register", fiber.Map{
			"Title": "Daftar",
			"Error": "Email sudah terdaftar, gunakan email lain.",
		}, "layouts/auth")
	}

	pelamar := &models.Pelamar{
		Name:  name,
		Email: email,
	}
	if err := pelamar.HashPassword(password); err != nil {
		return c.Render("pelamar/register", fiber.Map{
			"Title": "Daftar",
			"Error": "Terjadi kesalahan, coba lagi.",
		}, "layouts/auth")
	}

	if err := models.CreatePelamar(config.DB, pelamar); err != nil {
		return c.Render("pelamar/register", fiber.Map{
			"Title": "Daftar",
			"Error": "Gagal membuat akun, coba lagi.",
		}, "layouts/auth")
	}

	return c.Redirect("/login?registered=1")
}

// Logout clears the pelamar session and redirects to login.
func Logout(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err == nil {
		_ = sess.Destroy()
	}
	return c.Redirect("/login")
}
