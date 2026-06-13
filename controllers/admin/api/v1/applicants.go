package v1

import (
	"glk-web-app/config"
	"glk-web-app/models"
	"time"

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

// -----------------------------------------------------------------------
// GET /api/v1/recruitment/pelamar
// Returns a list of lamaran with date range filtering, search, and pagination.
// Query Params:
//   - start_date: YYYY-MM-DD (default: today)
//   - end_date: YYYY-MM-DD (default: today)
//   - search: string
//   - limit: int
//   - offset: int
// -----------------------------------------------------------------------
func ListRecruitmentPelamar(c *fiber.Ctx) error {
	startDate := c.Query("start_date", "")
	endDate := c.Query("end_date", "")
	search := c.Query("search", "")
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	// Default to today if dates are not provided
	today := time.Now().Format("2006-01-02")
	if startDate == "" {
		startDate = today
	}
	if endDate == "" {
		endDate = today
	}

	query := config.DB.Model(&models.Lamaran{}).
		Preload("Pelamar").
		Preload("Kota").
		Preload("Kecamatan").
		Preload("TargetJenjang").
		Preload("TargetMapel")

	// Filter by date range
	query = query.Where("DATE(created_at) BETWEEN ? AND ?", startDate, endDate)

	// Search by name or email
	if search != "" {
		query = query.Joins("JOIN pelamars ON pelamars.id = lamarans.pelamar_id").
			Where("pelamars.name LIKE ? OR pelamars.email LIKE ? OR lamarans.nama_lengkap LIKE ?",
				"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	var items []models.Lamaran
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&items).Error
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil data recruitment pelamar")
	}

	return successResponse(c, fiber.Map{
		"total":      total,
		"items":      items,
		"start_date": startDate,
		"end_date":   endDate,
		"limit":      limit,
		"offset":     offset,
	})
}

// -----------------------------------------------------------------------
// GET /api/v1/recruitment/pelamar/:id
// Returns full details of a lamaran including test results.
// -----------------------------------------------------------------------
func GetRecruitmentPelamarDetail(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "ID tidak valid")
	}

	var lamaran models.Lamaran
	err = config.DB.Preload("Pelamar").
		Preload("Kota").
		Preload("Kecamatan").
		Preload("TargetJenjang").
		Preload("TargetMapel").
		First(&lamaran, id).Error
	if err != nil {
		return errorResponse(c, fiber.StatusNotFound, "Lamaran tidak ditemukan")
	}

	// Fetch BankSoal for test results
	var bankSoal models.BankSoalA
	err = config.DB.Preload("BankSoalBs.BankSoalCs").
		Where("jenis_pendidikan_id = ? AND mata_pelajaran_id = ?",
			lamaran.TargetJenjangID, lamaran.TargetMapelID).
		First(&bankSoal).Error

	testResults := fiber.Map{
		"finished": lamaran.Status == "Selesai Tes",
		"answers":  lamaran.Prioritas, // JSON string
		"bankSoal": bankSoal,
	}

	return successResponse(c, fiber.Map{
		"lamaran":     lamaran,
		"testResults": testResults,
	})
}

// -----------------------------------------------------------------------
// PATCH /api/v1/recruitment/pelamar/:id/correction
// Body: { "corrections": { "qID": "T" | "F" } }
// Updates manual corrections for test results.
// -----------------------------------------------------------------------
func UpdateRecruitmentCorrection(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "ID tidak valid")
	}

	type Body struct {
		Corrections string `json:"corrections"` // JSON string
	}
	var body Body
	if err := c.BodyParser(&body); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Body tidak valid")
	}

	if err := config.DB.Model(&models.Lamaran{}).Where("id = ?", id).Update("koreksi_nilai", body.Corrections).Error; err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Gagal menyimpan koreksi")
	}

	return successResponse(c, fiber.Map{"id": id, "updated": true})
}
