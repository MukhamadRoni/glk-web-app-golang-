package admin

import (
	"glk-web-app/config"
	"glk-web-app/models"

	"github.com/gofiber/fiber/v2"
)

// ShowConfidenceScorePage renders the master page for confidence scores
func ShowConfidenceScorePage(c *fiber.Ctx) error {
	return c.Render("admin/master/confidence_score", contextData(c, fiber.Map{
		"Title":      "Master Confidence Score",
		"Breadcrumb": "Master / Confidence Score",
	}), "layouts/base")
}

// GetConfidenceScoreList retrieves all confidence scores
func GetConfidenceScoreList(c *fiber.Ctx) error {
	var scores []models.ConfidenceScore
	if err := config.DB.Order("min_score asc").Find(&scores).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": scores})
}

// CreateConfidenceScore creates a new confidence score record
func CreateConfidenceScore(c *fiber.Ctx) error {
	var cs models.ConfidenceScore
	if err := c.BodyParser(&cs); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := config.DB.Create(&cs).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": cs})
}

// UpdateConfidenceScore updates an existing confidence score record
func UpdateConfidenceScore(c *fiber.Ctx) error {
	id := c.Params("id")
	var cs models.ConfidenceScore
	if err := config.DB.First(&cs, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Confidence score not found"})
	}

	if err := c.BodyParser(&cs); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := config.DB.Save(&cs).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": cs})
}

// DeleteConfidenceScore deletes a confidence score record
func DeleteConfidenceScore(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := config.DB.Delete(&models.ConfidenceScore{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Confidence score deleted"})
}
