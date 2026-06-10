package pelamar

import (
	// "context"
	"encoding/json"
	"fmt"
	"glk-web-app/config"
	"glk-web-app/models"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

// ShowTestIntro renders the introduction for the academic test.
func ShowTestIntro(c *fiber.Ctx) error {
	pelamarID, ok := c.Locals("pelamar_id").(uint)
	if !ok {
		return c.Redirect("/login")
	}

	hasApplied, lamaran, err := models.CheckIfPelamarHasApplied(config.DB, pelamarID)
	if err != nil || !hasApplied {
		return c.Redirect("/dashboard")
	}

	// Cari BankSoalA berdasarkan jenjang dan mapel lamaran
	var bankSoal models.BankSoalA
	err = config.DB.Where("jenis_pendidikan_id = ? AND mata_pelajaran_id = ? AND active = 'T'",
		lamaran.TargetJenjangID, lamaran.TargetMapelID).First(&bankSoal).Error

	if err != nil {
		// Jika tidak ada tes untuk posisi ini, anggap selesai
		return c.Redirect("/dashboard")
	}

	// Cek apakah tes sudah selesai
	if lamaran.Status == "Selesai Tes" {
		return c.Redirect("/dashboard")
	}

	return c.Render("pelamar/test_intro", contextData(c, fiber.Map{
		"Title":      "Pengenalan Tes Akademik",
		"Breadcrumb": "Tes Akademik",
		"BankSoal":   bankSoal,
	}), "layouts/horizontal")
}

// StartTest initializes the test session in Redis.
func StartTest(c *fiber.Ctx) error {
	pelamarID, ok := c.Locals("pelamar_id").(uint)
	if !ok {
		return c.Redirect("/login")
	}

	_, lamaran, err := models.CheckIfPelamarHasApplied(config.DB, pelamarID)
	if err != nil || lamaran == nil {
		return c.Redirect("/dashboard")
	}

	var bankSoal models.BankSoalA
	err = config.DB.Where("jenis_pendidikan_id = ? AND mata_pelajaran_id = ? AND active = 'T'",
		lamaran.TargetJenjangID, lamaran.TargetMapelID).First(&bankSoal).Error
	if err != nil {
		return c.Redirect("/dashboard")
	}

	// Check Redis for existing session
	sessionKey := fmt.Sprintf("test:pelamar:%d:start", pelamarID)
	exists, _ := config.RDB.Exists(config.Ctx, sessionKey).Result()
	if exists == 0 {
		// Start session
		config.RDB.Set(config.Ctx, sessionKey, time.Now().Unix(), time.Duration(bankSoal.DurasiPengerjaan+1)*time.Minute)
	}

	return c.Redirect("/test/soal")
}

// ShowTestSoal renders the test UI.
func ShowTestSoal(c *fiber.Ctx) error {
	pelamarID, ok := c.Locals("pelamar_id").(uint)
	if !ok {
		return c.Redirect("/login")
	}

	_, lamaran, err := models.CheckIfPelamarHasApplied(config.DB, pelamarID)
	if err != nil || lamaran == nil {
		return c.Redirect("/dashboard")
	}

	var bankSoal models.BankSoalA
	err = config.DB.Preload("BankSoalBs.BankSoalCs").Where("jenis_pendidikan_id = ? AND mata_pelajaran_id = ? AND active = 'T'",
		lamaran.TargetJenjangID, lamaran.TargetMapelID).First(&bankSoal).Error
	if err != nil {
		return c.Redirect("/dashboard")
	}

	sessionKey := fmt.Sprintf("test:pelamar:%d:start", pelamarID)
	startStr, err := config.RDB.Get(config.Ctx, sessionKey).Result()
	if err == redis.Nil {
		// Session expired or not started
		return c.Redirect("/test/intro")
	}

	// Hitung sisa waktu
	var startUnix int64
	fmt.Sscanf(startStr, "%d", &startUnix)
	elapsed := time.Now().Unix() - startUnix
	remaining := (int64(bankSoal.DurasiPengerjaan) * 60) - elapsed

	if remaining <= 0 {
		// Time is up, auto finish
		return FinishTest(c)
	}

	// Get answers from Redis
	answersKey := fmt.Sprintf("test:pelamar:%d:answers", pelamarID)
	answersMap, _ := config.RDB.HGetAll(config.Ctx, answersKey).Result()

	answersJSON, _ := json.Marshal(answersMap)

	return c.Render("pelamar/test_soal", contextData(c, fiber.Map{
		"Title":      "Tes Akademik",
		"Breadcrumb": "Tes Akademik",
		"BankSoal":   bankSoal,
		"Remaining":  remaining,
		"Answers":    string(answersJSON),
	}), "layouts/horizontal")
}

// SaveAnswer saves an answer to Redis via AJAX.
func SaveAnswer(c *fiber.Ctx) error {
	pelamarID, ok := c.Locals("pelamar_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req struct {
		QuestionID uint   `json:"question_id"`
		Answer     string `json:"answer"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	answersKey := fmt.Sprintf("test:pelamar:%d:answers", pelamarID)
	config.RDB.HSet(config.Ctx, answersKey, fmt.Sprintf("%d", req.QuestionID), req.Answer)

	return c.JSON(fiber.Map{"status": "success"})
}

// FinishTest finishes the test, calculates score, and updates DB.
func FinishTest(c *fiber.Ctx) error {
	pelamarID, ok := c.Locals("pelamar_id").(uint)
	if !ok {
		return c.Redirect("/login")
	}

	hasApplied, lamaran, err := models.CheckIfPelamarHasApplied(config.DB, pelamarID)
	if err != nil || !hasApplied {
		return c.Redirect("/dashboard")
	}

	var bankSoal models.BankSoalA
	err = config.DB.Preload("BankSoalBs.BankSoalCs").Where("jenis_pendidikan_id = ? AND mata_pelajaran_id = ? AND active = 'T'",
		lamaran.TargetJenjangID, lamaran.TargetMapelID).First(&bankSoal).Error
	if err != nil {
		return c.Redirect("/dashboard")
	}

	answersKey := fmt.Sprintf("test:pelamar:%d:answers", pelamarID)
	answersMap, _ := config.RDB.HGetAll(config.Ctx, answersKey).Result()

	// Hitung nilai, simpan log, dll
	correct := 0
	for _, q := range bankSoal.BankSoalBs {
		ans, ok := answersMap[fmt.Sprintf("%d", q.ID)]
		if ok {
			for _, opt := range q.BankSoalCs {
				if opt.OptionText == ans && opt.IsCorrect == "T" {
					correct++
					break
				}
			}
		}
	}

	// Update lamaran status
	config.DB.Model(&lamaran).Update("status", "Selesai Tes")

	// Simpan detail jawaban ke DB (opsional) - disini kita log saja atau simpan ke field baru.
	answersJSON, _ := json.Marshal(answersMap)
	// Kita update Prioritas atau field lain untuk menyimpan JSON jawaban sementara
	config.DB.Model(&lamaran).Update("Prioritas", string(answersJSON))

	// Bersihkan Redis
	config.RDB.Del(config.Ctx, answersKey)
	config.RDB.Del(config.Ctx, fmt.Sprintf("test:pelamar:%d:start", pelamarID))

	return c.Redirect("/dashboard?test_finished=1")
}
