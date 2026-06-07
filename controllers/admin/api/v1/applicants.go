package v1

import (
	"glk-web-app/config"
	"glk-web-app/models"

	"github.com/gofiber/fiber/v2"
)

// successResponse is a helper to return a consistent JSON success envelope.
func successResponse(c *fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// errorResponse is a helper to return a consistent JSON error envelope.
func errorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}

// -----------------------------------------------------------------------
// GET /api/v1/applicants
// Returns a paginated/full list of all pelamar.
// -----------------------------------------------------------------------
func ListApplicants(c *fiber.Ctx) error {
	applicants, err := models.GetAllPelamar(config.DB)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil data pelamar")
	}
	return successResponse(c, fiber.Map{
		"total": len(applicants),
		"items": applicants,
	})
}

// -----------------------------------------------------------------------
// GET /api/v1/applicants/:id
// Returns a single pelamar by ID.
// -----------------------------------------------------------------------
func GetApplicant(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "ID tidak valid")
	}

	pelamar, err := models.GetPelamarByID(config.DB, uint(id))
	if err != nil {
		return errorResponse(c, fiber.StatusNotFound, "Pelamar tidak ditemukan")
	}

	return successResponse(c, pelamar)
}

// -----------------------------------------------------------------------
// PATCH /api/v1/applicants/:id/status
// Body: { "status": "accepted" | "rejected" | "pending" }
// Updates the status of a pelamar.
// -----------------------------------------------------------------------
func UpdateApplicantStatus(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "ID tidak valid")
	}

	type Body struct {
		Status string `json:"status"`
	}
	var body Body
	if err := c.BodyParser(&body); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Body tidak valid")
	}

	allowed := map[string]bool{"accepted": true, "rejected": true, "pending": true}
	if !allowed[body.Status] {
		return errorResponse(c, fiber.StatusBadRequest, "Status tidak valid. Gunakan: accepted, rejected, pending")
	}

	if err := models.UpdatePelamarStatus(config.DB, uint(id), body.Status); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Gagal memperbarui status")
	}

	return successResponse(c, fiber.Map{
		"id":     id,
		"status": body.Status,
	})
}

// -----------------------------------------------------------------------
// DELETE /api/v1/applicants/:id
// Soft-deletes a pelamar record.
// -----------------------------------------------------------------------
func DeleteApplicant(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "ID tidak valid")
	}

	if err := models.DeletePelamar(config.DB, uint(id)); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Gagal menghapus pelamar")
	}

	return successResponse(c, fiber.Map{"id": id, "deleted": true})
}
