package admin

import (
	"glk-web-app/config"
	"glk-web-app/models"

	"github.com/gofiber/fiber/v2"
)

// ShowLogin renders the admin login page.
func ShowLogin(c *fiber.Ctx) error {
	return c.Render("admin/login", fiber.Map{
		"Title":       "Admin Login",
		"Description": "Masuk ke panel admin GLK",
	}, "layouts/auth")
}

// ProcessLogin handles the admin login form POST.
func ProcessLogin(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	admin, err := models.GetAdminByUsername(config.DB, username)
	if err != nil || !admin.CheckPassword(password) {
		return c.Render("admin/login", fiber.Map{
			"Title": "Admin Login",
			"Error": "Username atau password salah.",
		}, "layouts/auth")
	}

	sess, err := store.Get(c)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	sess.Set("admin_id", admin.ID)
	sess.Set("admin_username", admin.Username)
	sess.Set("admin_role_id", admin.RoleID)
	if err := sess.Save(); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Redirect("/admin/dashboard")
}

// Logout destroys the admin session and redirects to login.
func Logout(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err == nil {
		_ = sess.Destroy()
	}
	return c.Redirect("/admin/login")
}
