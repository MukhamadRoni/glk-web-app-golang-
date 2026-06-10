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
	}), "layouts/horizontal")
}

// ShowApply renders the job application form.
func ShowApply(c *fiber.Ctx) error {
	pelamarID, ok := c.Locals("pelamar_id").(uint)
	if !ok {
		return c.Redirect("/login")
	}

	pelamar, err := models.GetPelamarByID(config.DB, pelamarID)
	if err != nil {
		return c.Redirect("/login")
	}

	return c.Render("pelamar/apply", contextData(c, fiber.Map{
		"Title":       "Form Lamaran",
		"Breadcrumb":  "Lamaran",
		"Description": "Isi form lamaran di bawah ini",
		"Name":        pelamar.Name,
		"Email":       pelamar.Email,
	}), "layouts/horizontal")
}

// ProcessApply handles the job application submission.
func ProcessApply(c *fiber.Ctx) error {
	// TODO: Handle file upload (CV) and save to database
	// Also parse Kota, Kecamatan, Jenjang, Mapel, Jangkauan Mengajar
	return c.Redirect("/dashboard?success=1")
}
