package admin

import (
	"bytes"
	"encoding/json"
	"glk-web-app/config"
	"glk-web-app/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// ShowProfilingPage renders the AI Profiling management page
func ShowProfilingPage(c *fiber.Ctx) error {
	return c.Render("admin/ai/profiling", contextData(c, fiber.Map{
		"Title":      "AI Profiling Skills",
		"Breadcrumb": "AI Mode / Profiling",
	}), "layouts/base")
}

// GetProfilingList retrieves all AIProfilingSkill records
func GetProfilingList(c *fiber.Ctx) error {
	var skills []models.AIProfilingSkill
	if err := config.DB.Order("id desc").Find(&skills).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": skills})
}

// CreateProfiling creates a new AIProfilingSkill record
func CreateProfiling(c *fiber.Ctx) error {
	var skill models.AIProfilingSkill
	if err := c.BodyParser(&skill); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := config.DB.Create(&skill).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": skill})
}

// UpdateProfiling updates an existing AIProfilingSkill record
func UpdateProfiling(c *fiber.Ctx) error {
	id := c.Params("id")
	var skill models.AIProfilingSkill
	if err := config.DB.First(&skill, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Skill record not found"})
	}

	if err := c.BodyParser(&skill); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := config.DB.Save(&skill).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": skill})
}

// DeleteProfiling soft-deletes an AIProfilingSkill record
func DeleteProfiling(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := config.DB.Delete(&models.AIProfilingSkill{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Skill record deleted"})
}

// UploadProfilingProxy handles file upload to GDrive via Apps Script specifically for Profiling
func UploadProfilingProxy(c *fiber.Ctx) error {
	var payload struct {
		Filename string `json:"filename"`
		FileData string `json:"fileData"`
		MimeType string `json:"mimeType"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid payload"})
	}

	gasURL := config.GetEnv("GAS_PROFILING_URL", config.GetEnv("GAS_UPLOAD_URL", ""))
	if gasURL == "" {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "GAS_PROFILING_URL or GAS_UPLOAD_URL not set"})
	}

	payloadBytes, _ := json.Marshal(payload)
	resp, err := http.Post(gasURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to decode GAS response"})
	}

	if result["status"] == "success" {
		return c.JSON(fiber.Map{"success": true, "url": result["fileUrl"]})
	}

	return c.Status(500).JSON(fiber.Map{"success": false, "message": result["message"]})
}
