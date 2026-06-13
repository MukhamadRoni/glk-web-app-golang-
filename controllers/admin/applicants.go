package admin

import (
	"glk-web-app/config"
	"glk-web-app/models"

	"github.com/gofiber/fiber/v2"
)

// ShowApplicants renders the full applicants list page.
func ShowApplicants(c *fiber.Ctx) error {
	applicants, err := models.GetAllPelamar(config.DB)
	if err != nil {
		applicants = []models.Pelamar{}
	}

	return c.Render("admin/applicants", contextData(c, fiber.Map{
		"Title":      "Daftar Pelamar",
		"Breadcrumb": "Pelamar",
		"Applicants": applicants,
	}), "layouts/base")
}

// ShowApplicant renders the detail page for a single applicant.
func ShowApplicant(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("ID tidak valid")
	}

	pelamar, err := models.GetPelamarByID(config.DB, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Pelamar tidak ditemukan")
	}

	return c.Render("admin/applicant_detail", contextData(c, fiber.Map{
		"Title":      pelamar.Name,
		"Breadcrumb": "Detail Pelamar",
		"Pelamar":    pelamar,
	}), "layouts/base")
}

// ShowRecruitmentPelamar renders the recruitment transaction page for applicants.
func ShowRecruitmentPelamar(c *fiber.Ctx) error {
	return c.Render("admin/recruitment/transaksi/pelamar", contextData(c, fiber.Map{
		"Title":      "Transaksi Pelamar",
		"Breadcrumb": "Recruitment / Transaksi / Pelamar",
	}), "layouts/base")
}
