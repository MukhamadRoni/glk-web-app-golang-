package admin

import (
	"glk-web-app/config"
	"glk-web-app/models"

	"github.com/gofiber/fiber/v2"
)

// ShowDashboard renders the admin dashboard with summary statistics.
func ShowDashboard(c *fiber.Ctx) error {
	allPelamar, err := models.GetAllPelamar(config.DB)
	if err != nil {
		allPelamar = []models.Pelamar{}
	}

	// Calculate stats
	var accepted, pending, rejected int
	for _, p := range allPelamar {
		switch p.Status {
		case "accepted":
			accepted++
		case "rejected":
			rejected++
		default:
			pending++
		}
	}

	// Show only the 10 most recent applicants in the table
	recent := allPelamar
	if len(recent) > 10 {
		recent = recent[:10]
	}

	return c.Render("admin/dashboard", contextData(c, fiber.Map{
		"Title":      "Dashboard",
		"Breadcrumb": "Dashboard",
		"Stats": fiber.Map{
			"TotalApplicants": len(allPelamar),
			"Accepted":        accepted,
			"Pending":         pending,
			"Rejected":        rejected,
		},
		"Applicants": recent,
	}), "layouts/base")
}
