package admin

import (
	"database/sql"
	"glk-web-app/config"

	// "glk-web-app/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

type MapelResponse struct {
	ID                  uint      `json:"id"`
	Code                string    `json:"code"`
	Name                string    `json:"name"`
	JenisPendidikanID   *uint     `json:"jenis_pendidikan_id"`
	JenisPendidikanName *string   `json:"jenis_pendidikan_name"`
	Active              string    `json:"active"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// GetMapelsList retrieves all MataPelajaran using Raw SQL
func GetMapelsList(c *fiber.Ctx) error {
	var mapels []MapelResponse
	// Fetching active and inactive records, avoiding soft-deleted ones
	err := config.DB.Raw(`
		SELECT m.id, m.code, m.name, m.jenis_pendidikan_id, j.name AS jenis_pendidikan_name, 
		       m.created_at, m.updated_at, m.active 
		FROM mata_pelajarans m 
		LEFT JOIN jenis_pendidikans j ON m.jenis_pendidikan_id = j.id 
		WHERE m.deleted_at IS NULL
	`).Scan(&mapels).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": mapels})
}

// CreateMapel creates a new MataPelajaran using Raw SQL
func CreateMapel(c *fiber.Ctx) error {
	var payload struct {
		JenisPendidikanID uint   `json:"jenis_pendidikan_id"`
		Code              string `json:"code"`
		Name              string `json:"name"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	if payload.Code == "" || payload.Name == "" || payload.JenisPendidikanID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Jenis Pendidikan, Code and Name are required"})
	}

	now := time.Now()
	err := config.DB.Exec("INSERT INTO mata_pelajarans (jenis_pendidikan_id, code, name, created_at, updated_at, active) VALUES (?, ?, ?, ?, ?, 'T')",
		payload.JenisPendidikanID, payload.Code, payload.Name, now, now).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Mata Pelajaran created successfully"})
}

// UpdateMapel updates MataPelajaran using Raw SQL
func UpdateMapel(c *fiber.Ctx) error {
	id := c.Params("id")
	var payload struct {
		JenisPendidikanID uint   `json:"jenis_pendidikan_id"`
		Code              string `json:"code"`
		Name              string `json:"name"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	if payload.Code == "" || payload.Name == "" || payload.JenisPendidikanID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Jenis Pendidikan, Code and Name are required"})
	}

	now := time.Now()
	err := config.DB.Exec("UPDATE mata_pelajarans SET jenis_pendidikan_id = ?, code = ?, name = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL",
		payload.JenisPendidikanID, payload.Code, payload.Name, now, id).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Mata Pelajaran updated successfully"})
}

// UpdateMapelActive toggles the active status using Raw SQL
func UpdateMapelActive(c *fiber.Ctx) error {
	id := c.Params("id")
	var payload struct {
		Active string `json:"active"` // Expected 'T' or 'F'
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	if payload.Active != "T" && payload.Active != "F" {
		return c.Status(400).JSON(fiber.Map{"error": "Active must be 'T' or 'F'"})
	}

	now := time.Now()
	err := config.DB.Exec("UPDATE mata_pelajarans SET active = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL",
		payload.Active, now, id).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Status updated successfully"})
}

// DeleteMapel soft deletes MataPelajaran using Raw SQL
func DeleteMapel(c *fiber.Ctx) error {
	id := c.Params("id")
	now := time.Now()
	// Soft delete by setting deleted_at
	err := config.DB.Exec("UPDATE mata_pelajarans SET deleted_at = ? WHERE id = ?", sql.NullTime{Time: now, Valid: true}, id).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Mata Pelajaran deleted"})
}

// ShowMapelPage renders the page
func ShowMapelPage(c *fiber.Ctx) error {
	return c.Render("admin/master/mapel", contextData(c, fiber.Map{
		"Title": "Master Mata Pelajaran",
	}), "layouts/base")
}
