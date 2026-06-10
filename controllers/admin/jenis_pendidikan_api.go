package admin

import (
	"glk-web-app/config"
	"glk-web-app/models"

	"github.com/gofiber/fiber/v2"
)

// GetJenisPendidikanList retrieves all JenisPendidikan
func GetJenisPendidikanList(c *fiber.Ctx) error {
	var jps []models.JenisPendidikan
	if err := config.DB.Find(&jps).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": jps})
}

// CreateJenisPendidikan creates a new JenisPendidikan
func CreateJenisPendidikan(c *fiber.Ctx) error {
	var payload struct {
		JenisPendidikan string `json:"jenis_pendidikan"`
		Name            string `json:"name"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	jp := models.JenisPendidikan{
		JenisPendidikan: payload.JenisPendidikan,
		Name:            payload.Name,
	}
	if err := config.DB.Create(&jp).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": jp})
}

// UpdateJenisPendidikan updates a JenisPendidikan
func UpdateJenisPendidikan(c *fiber.Ctx) error {
	id := c.Params("id")
	var payload struct {
		JenisPendidikan string `json:"jenis_pendidikan"`
		Name            string `json:"name"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	var jp models.JenisPendidikan
	if err := config.DB.First(&jp, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Jenis Pendidikan not found"})
	}

	jp.JenisPendidikan = payload.JenisPendidikan
	jp.Name = payload.Name
	if err := config.DB.Save(&jp).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": jp})
}

// UpdateJenisPendidikanActive toggles the active status
func UpdateJenisPendidikanActive(c *fiber.Ctx) error {
	id := c.Params("id")
	var payload struct {
		Active string `json:"active"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	if payload.Active != "T" && payload.Active != "F" {
		return c.Status(400).JSON(fiber.Map{"error": "Active must be 'T' or 'F'"})
	}

	var jp models.JenisPendidikan
	if err := config.DB.First(&jp, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Jenis Pendidikan not found"})
	}

	jp.Active = payload.Active
	if err := config.DB.Save(&jp).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Status updated successfully"})
}

// DeleteJenisPendidikan deletes a JenisPendidikan
func DeleteJenisPendidikan(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := config.DB.Delete(&models.JenisPendidikan{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	// Cascade? Let's leave it simple for now
	return c.JSON(fiber.Map{"status": "success", "message": "Jenis Pendidikan deleted"})
}

// ShowJenisPendidikanPage renders the page
func ShowJenisPendidikanPage(c *fiber.Ctx) error {
	return c.Render("admin/master/jenis_pendidikan", contextData(c, fiber.Map{
		"Title": "Master Jenis Pendidikan",
	}), "layouts/base")
}
