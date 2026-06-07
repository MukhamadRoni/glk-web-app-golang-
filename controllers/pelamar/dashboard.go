package pelamar

import (
	"glk-web-app/config"
	"glk-web-app/models"

	"github.com/gofiber/fiber/v2"
)

// ShowDashboard renders the pelamar dashboard page.
func ShowDashboard(c *fiber.Ctx) error {
	pelamarID, ok := c.Locals("pelamar_id").(uint)
	if !ok {
		return c.Redirect("/login")
	}

	pelamar, err := models.GetPelamarByID(config.DB, pelamarID)
	if err != nil {
		return c.Redirect("/login")
	}

	return c.Render("pelamar/dashboard", contextData(c, fiber.Map{
		"Title":       "Dashboard Saya",
		"Breadcrumb":  "Dashboard",
		"Description": "Pantau status lamaran Anda",
		"Name":        pelamar.Name,
		// Applications will be added once the Application model is created.
		"Applications": []interface{}{},
	}), "layouts/base")
}
