package admin

import (
	"glk-web-app/config"
	"glk-web-app/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetBankSoalList retrieves all BankSoal headers
func GetBankSoalList(c *fiber.Ctx) error {
	var banks []models.BankSoalA
	if err := config.DB.Preload("JenisPendidikan").Preload("MataPelajaran").Order("created_at desc").Find(&banks).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": banks})
}

// GetBankSoalDetail retrieves a specific BankSoal with questions and options
func GetBankSoalDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	var bank models.BankSoalA
	if err := config.DB.Preload("JenisPendidikan").Preload("MataPelajaran").
		Preload("BankSoalBs", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index asc")
		}).
		Preload("BankSoalBs.BankSoalCs").
		First(&bank, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Bank Soal not found"})
	}
	return c.JSON(fiber.Map{"status": "success", "data": bank})
}

// SaveBankSoalPayload represents the nested JSON structure
type SaveBankSoalPayload struct {
	JenisPendidikanID uint                     `json:"jenis_pendidikan_id"`
	MataPelajaranID   uint                     `json:"mata_pelajaran_id"`
	Title             string                   `json:"title"`
	Active            string                   `json:"active"`
	SaveAsNewVersion  bool                     `json:"save_as_new_version"`
	Questions         []SaveBankSoalQuestion   `json:"questions"`
}

type SaveBankSoalQuestion struct {
	ID            uint                  `json:"id"`
	QuestionType  string                `json:"question_type"`
	QuestionText  string                `json:"question_text"`
	QuestionImage string                `json:"question_image"`
	OrderIndex    int                   `json:"order_index"`
	Options       []SaveBankSoalOption  `json:"options"`
}

type SaveBankSoalOption struct {
	ID          uint   `json:"id"`
	OptionText  string `json:"option_text"`
	OptionImage string `json:"option_image"`
	IsCorrect   string `json:"is_correct"`
}

// SaveBankSoal creates or updates a Bank Soal
func SaveBankSoal(c *fiber.Ctx) error {
	id := c.Params("id") // Optional, empty means create
	var payload SaveBankSoalPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON payload"})
	}

	if payload.Title == "" || payload.JenisPendidikanID == 0 || payload.MataPelajaranID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Title, Jenis Pendidikan, and Mata Pelajaran are required"})
	}

	tx := config.DB.Begin()

	var bankA models.BankSoalA
	isNewRecord := true

	if id != "" && !payload.SaveAsNewVersion {
		// Edit existing
		if err := tx.First(&bankA, id).Error; err != nil {
			tx.Rollback()
			return c.Status(404).JSON(fiber.Map{"error": "Bank Soal not found"})
		}
		isNewRecord = false
	}

	if payload.SaveAsNewVersion || id == "" {
		// Find max version for this Title and Mapel if saving as new version
		var maxVersion int64 = 0
		tx.Model(&models.BankSoalA{}).
			Where("title = ? AND mata_pelajaran_id = ?", payload.Title, payload.MataPelajaranID).
			Select("COALESCE(MAX(version), 0)").Scan(&maxVersion)

		bankA = models.BankSoalA{
			JenisPendidikanID: payload.JenisPendidikanID,
			MataPelajaranID:   payload.MataPelajaranID,
			Title:             payload.Title,
			Version:           int(maxVersion) + 1,
			Active:            payload.Active,
		}
		if bankA.Active == "" {
			bankA.Active = "T"
		}
		if err := tx.Create(&bankA).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create header"})
		}
	} else {
		// Update header
		bankA.JenisPendidikanID = payload.JenisPendidikanID
		bankA.MataPelajaranID = payload.MataPelajaranID
		bankA.Title = payload.Title
		if payload.Active != "" {
			bankA.Active = payload.Active
		}
		if err := tx.Save(&bankA).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update header"})
		}

		// Delete old B and C to rebuild cleanly
		// First get existing Bs to delete Cs
		var oldBs []models.BankSoalB
		tx.Where("bank_soal_a_id = ?", bankA.ID).Find(&oldBs)
		for _, b := range oldBs {
			tx.Where("bank_soal_b_id = ?", b.ID).Delete(&models.BankSoalC{})
		}
		tx.Where("bank_soal_a_id = ?", bankA.ID).Delete(&models.BankSoalB{})
	}

	// Insert new questions and options
	for i, q := range payload.Questions {
		newB := models.BankSoalB{
			BankSoalAID:   bankA.ID,
			QuestionType:  q.QuestionType,
			QuestionText:  q.QuestionText,
			QuestionImage: q.QuestionImage,
			OrderIndex:    i, // Force order based on array
		}
		if err := tx.Create(&newB).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save question"})
		}

		for _, opt := range q.Options {
			newC := models.BankSoalC{
				BankSoalBID: newB.ID,
				OptionText:  opt.OptionText,
				OptionImage: opt.OptionImage,
				IsCorrect:   opt.IsCorrect,
			}
			if newC.IsCorrect == "" {
				newC.IsCorrect = "F"
			}
			if err := tx.Create(&newC).Error; err != nil {
				tx.Rollback()
				return c.Status(500).JSON(fiber.Map{"error": "Failed to save option"})
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Transaction failed"})
	}

	msg := "Bank Soal created successfully"
	if !isNewRecord {
		msg = "Bank Soal updated successfully"
	}
	if payload.SaveAsNewVersion {
		msg = "New version created successfully"
	}

	return c.JSON(fiber.Map{"status": "success", "message": msg, "data": bankA})
}

// DeleteBankSoal deletes a BankSoal header and its contents
func DeleteBankSoal(c *fiber.Ctx) error {
	id := c.Params("id")
	
	tx := config.DB.Begin()

	var oldBs []models.BankSoalB
	tx.Where("bank_soal_a_id = ?", id).Find(&oldBs)
	for _, b := range oldBs {
		tx.Where("bank_soal_b_id = ?", b.ID).Delete(&models.BankSoalC{})
	}
	tx.Where("bank_soal_a_id = ?", id).Delete(&models.BankSoalB{})
	
	if err := tx.Delete(&models.BankSoalA{}, id).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	tx.Commit()
	return c.JSON(fiber.Map{"status": "success", "message": "Bank Soal deleted"})
}

// UpdateBankSoalActive toggles the active status
func UpdateBankSoalActive(c *fiber.Ctx) error {
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

	if err := config.DB.Model(&models.BankSoalA{}).Where("id = ?", id).Update("active", payload.Active).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Status updated successfully"})
}

// ShowBankSoalPage renders the index page
func ShowBankSoalPage(c *fiber.Ctx) error {
	return c.Render("admin/master/bank_soal", contextData(c, fiber.Map{
		"Title": "Master Bank Soal",
	}), "layouts/base")
}

// ShowBankSoalFormPage renders the create/edit page
func ShowBankSoalFormPage(c *fiber.Ctx) error {
	return c.Render("admin/master/bank_soal_form", contextData(c, fiber.Map{
		"Title": "Form Bank Soal",
	}), "layouts/base")
}
