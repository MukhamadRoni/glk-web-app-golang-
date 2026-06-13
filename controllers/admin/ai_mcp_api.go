package admin

import (
	"bytes"
	"encoding/json"
	"glk-web-app/config"
	"glk-web-app/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// ShowMCPPage renders the AI MCP management page
func ShowMCPPage(c *fiber.Ctx) error {
	return c.Render("admin/ai/mcp", contextData(c, fiber.Map{
		"Title":      "AI Master Context Provider (MCP)",
		"Breadcrumb": "AI Mode / MCP",
	}), "layouts/base")
}

// GetMCPList retrieves all CompanyMCP records
func GetMCPList(c *fiber.Ctx) error {
	var mcps []models.CompanyMCP
	if err := config.DB.Order("id desc").Find(&mcps).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": mcps})
}

// CreateMCP creates a new CompanyMCP record
func CreateMCP(c *fiber.Ctx) error {
	var mcp models.CompanyMCP
	if err := c.BodyParser(&mcp); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := config.DB.Create(&mcp).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": mcp})
}

// UpdateMCP updates an existing CompanyMCP record
func UpdateMCP(c *fiber.Ctx) error {
	id := c.Params("id")
	var mcp models.CompanyMCP
	if err := config.DB.First(&mcp, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "MCP record not found"})
	}

	if err := c.BodyParser(&mcp); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := config.DB.Save(&mcp).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": mcp})
}

// DeleteMCP soft-deletes a CompanyMCP record
func DeleteMCP(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := config.DB.Delete(&models.CompanyMCP{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "MCP record deleted"})
}

// UploadMCPProxy handles file upload to GDrive via Apps Script specifically for MCP
func UploadMCPProxy(c *fiber.Ctx) error {
	var payload struct {
		Filename string `json:"filename"`
		FileData string `json:"fileData"`
		MimeType string `json:"mimeType"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid payload"})
	}

	// We'll use a specific GAS URL for MCP if defined, otherwise fallback
	gasURL := config.GetEnv("GAS_MCP_URL", config.GetEnv("GAS_UPLOAD_URL", ""))
	if gasURL == "" {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "GAS_MCP_URL or GAS_UPLOAD_URL not set"})
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
