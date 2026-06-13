package admin

import (
	"glk-web-app/config"
	"glk-web-app/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

// ShowDashboard renders the global HRIS dashboard with a welcome message.
func ShowDashboard(c *fiber.Ctx) error {
	return c.Render("admin/dashboard", contextData(c, fiber.Map{
		"Title":      "Dashboard Global",
		"Breadcrumb": "Dashboard",
		"Welcome":    "Selamat datang di HRIS Gurulesku",
	}), "layouts/base")
}

// ShowRecruitmentDashboard renders the recruitment specific dashboard with summary statistics.
func ShowRecruitmentDashboard(c *fiber.Ctx) error {
	var lamarans []models.Lamaran
	// Gunakan Order ID DESC sebagai cadangan jika created_at bermasalah,
	// dan pastikan Preload Pelamar sukses
	err := config.DB.Preload("Pelamar").Order("id DESC").Find(&lamarans).Error
	if err != nil {
		log.Println("[Dashboard] Error fetching lamarans:", err)
		lamarans = []models.Lamaran{}
	}

	log.Printf("[Dashboard] Total lamarans found: %d", len(lamarans))

	// Calculate stats
	var accepted, pending, rejected int
	for _, l := range lamarans {
		if l.Status == "Diterima" {
			accepted++
		} else if l.Status == "Ditolak" {
			rejected++
		} else {
			pending++
		}
	}

	// Show only the 10 most recent lamaran in the table
	recent := lamarans
	if len(recent) > 10 {
		recent = recent[:10]
	}

	return c.Render("admin/recruitment/dashboard", contextData(c, fiber.Map{
		"Title":      "Recruitment Dashboard",
		"Breadcrumb": "Recruitment / Dashboard",
		"Stats": fiber.Map{
			"TotalApplicants": len(lamarans),
			"Accepted":        accepted,
			"Pending":         pending,
			"Rejected":        rejected,
		},
		"Lamarans": recent,
	}), "layouts/base")
}
