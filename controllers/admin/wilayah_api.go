package admin

import (
	"glk-web-app/config"
	"glk-web-app/models"

	"github.com/gofiber/fiber/v2"
)

// --- Kota API Handlers ---

func GetKotasList(c *fiber.Ctx) error {
	var kotas []models.Kota
	if err := config.DB.Find(&kotas).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": kotas})
}

func CreateKota(c *fiber.Ctx) error {
	var payload struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	kota := models.Kota{Name: payload.Name}
	if err := config.DB.Create(&kota).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": kota})
}

func DeleteKota(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := config.DB.Delete(&models.Kota{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	// cascade delete? In GORM soft delete usually only deletes the parent. 
	// To be safe we should also delete kecamatans.
	config.DB.Where("kota_id = ?", id).Delete(&models.Kecamatan{})

	return c.JSON(fiber.Map{"status": "success", "message": "Kota deleted"})
}

func UpdateKota(c *fiber.Ctx) error {
	id := c.Params("id")
	var payload struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	var kota models.Kota
	if err := config.DB.First(&kota, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Kota not found"})
	}

	kota.Name = payload.Name
	if err := config.DB.Save(&kota).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": kota})
}

// --- Kecamatan API Handlers ---

func GetKecamatansList(c *fiber.Ctx) error {
	var kecamatans []models.Kecamatan
	if err := config.DB.Preload("Kota").Find(&kecamatans).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": kecamatans})
}

func CreateKecamatan(c *fiber.Ctx) error {
	var payload struct {
		KotaID uint   `json:"kota_id"`
		Name   string `json:"name"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	kec := models.Kecamatan{
		KotaID: payload.KotaID,
		Name:   payload.Name,
	}
	if err := config.DB.Create(&kec).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": kec})
}

func DeleteKecamatan(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := config.DB.Delete(&models.Kecamatan{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Kecamatan deleted"})
}

func UpdateKecamatan(c *fiber.Ctx) error {
	id := c.Params("id")
	var payload struct {
		KotaID uint   `json:"kota_id"`
		Name   string `json:"name"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	var kec models.Kecamatan
	if err := config.DB.First(&kec, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Kecamatan not found"})
	}

	kec.KotaID = payload.KotaID
	kec.Name = payload.Name
	if err := config.DB.Save(&kec).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": kec})
}

// --- Wilayah Page Render ---

func ShowWilayahPage(c *fiber.Ctx) error {
	return c.Render("admin/master/wilayah", contextData(c, fiber.Map{
		"Title": "Master Wilayah (Kota & Kecamatan)",
	}), "layouts/base")
}
