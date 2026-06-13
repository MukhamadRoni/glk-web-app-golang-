package pelamar

import (
	"encoding/json"
	"glk-web-app/config"
	"glk-web-app/models"
	"glk-web-app/utils"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ShowDashboard renders the pelamar dashboard page.
func ShowDashboard(c *fiber.Ctx) error {
	pelamarID, ok := c.Locals("pelamar_id").(uint)
	if !ok {
		return c.Redirect("/login")
	}

	pelamar, err := models.GetPelamarByID(config.DB, pelamarID)
	if err != nil {
		return c.Redirect("/login")
	}

	hasApplied, lamaran, err := models.CheckIfPelamarHasApplied(config.DB, pelamarID)

	var applications []models.Lamaran
	hasPendingTest := false

	if err == nil && hasApplied && lamaran != nil {
		applications = append(applications, *lamaran)

		// Check if test exists for this lamaran
		if lamaran.Status != "Selesai Tes" {
			var count int64
			config.DB.Model(&models.BankSoalA{}).Where("jenis_pendidikan_id = ? AND mata_pelajaran_id = ? AND active = 'T'",
				lamaran.TargetJenjangID, lamaran.TargetMapelID).Count(&count)
			if count > 0 {
				hasPendingTest = true
			}
		}
	}

	return c.Render("pelamar/dashboard", contextData(c, fiber.Map{
		"Title":          "Dashboard Saya",
		"Breadcrumb":     "Dashboard",
		"Description":    "Pantau status lamaran Anda",
		"Name":           pelamar.Name,
		"Applications":   applications,
		"HasPendingTest": hasPendingTest,
	}), "layouts/horizontal")
}

// ShowApply renders the job application form.
func ShowApply(c *fiber.Ctx) error {
	pelamarID, ok := c.Locals("pelamar_id").(uint)
	if !ok {
		return c.Redirect("/login")
	}

	pelamar, err := models.GetPelamarByID(config.DB, pelamarID)
	if err != nil {
		return c.Redirect("/login")
	}

	hasApplied, lamaran, err := models.CheckIfPelamarHasApplied(config.DB, pelamarID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Render("pelamar/apply", contextData(c, fiber.Map{
		"Title":       "Form Lamaran",
		"Breadcrumb":  "Lamaran",
		"Description": "Isi form lamaran di bawah ini",
		"Name":        pelamar.Name,
		"Email":       pelamar.Email,
		"HasApplied":  hasApplied,
		"Lamaran":     lamaran,
	}), "layouts/horizontal")
}

// ProcessApply handles the job application submission.
func ProcessApply(c *fiber.Ctx) error {
	pelamarID, ok := c.Locals("pelamar_id").(uint)
	if !ok {
		return c.Redirect("/login")
	}

	pelamar, _ := models.GetPelamarByID(config.DB, pelamarID)

	// Prevent double submission
	hasApplied, _, err := models.CheckIfPelamarHasApplied(config.DB, pelamarID)
	if err != nil || hasApplied {
		return c.Redirect("/dashboard")
	}

	// Parse the multipart form
	if err := c.BodyParser(&fiber.Map{}); err != nil { // Ensure body is parsed
		log.Println("Error parsing body:", err)
	}

	// 1. Ambil Nama untuk penamaan file
	namaLengkap := c.FormValue("namaLengkap")
	if namaLengkap == "" && pelamar != nil {
		namaLengkap = pelamar.Name
	}
	safeNama := strings.ReplaceAll(namaLengkap, " ", "_")
	timestamp := time.Now().Format("20060102_150405")

	// 1. Ambil File Transkrip
	transkripFile, err := c.FormFile("transkrip")
	if err != nil {
		log.Println("Transkrip file error:", err)
		return c.Redirect("/apply?error=Transkrip file is required")
	}

	// Upload Transkrip ke Google Drive
	transkripURL, err := utils.UploadToGDrive(transkripFile, "TRANSKRIP_"+safeNama+"_"+timestamp)
	if err != nil {
		log.Println("Failed to upload transkrip:", err)
		return c.Redirect("/apply?error=Failed to upload transkrip")
	}

	// 2. Ambil File CV
	cvFile, err := c.FormFile("cv")
	if err != nil {
		log.Println("CV file error:", err)
		return c.Redirect("/apply?error=CV file is required")
	}

	// Upload CV ke Google Drive
	cvURL, err := utils.UploadToGDrive(cvFile, "CV_"+safeNama+"_"+timestamp)
	if err != nil {
		log.Println("Failed to upload CV:", err)
		return c.Redirect("/apply?error=Failed to upload CV")
	}

	// 3. Ambil data form
	form, err := c.MultipartForm()
	if err != nil || form == nil {
		log.Println("MultipartForm error:", err)
		return c.Redirect("/apply?error=Invalid form data")
	}

	kotaID, _ := strconv.Atoi(c.FormValue("kotaDomisili"))
	kecamatanID, _ := strconv.Atoi(c.FormValue("kecamatanDomisili"))
	targetJenjangID, _ := strconv.Atoi(c.FormValue("jenjang"))
	targetMapelID, _ := strconv.Atoi(c.FormValue("mapel"))

	jangkauanArray := form.Value["jangkauanMengajar[]"]
	jangkauan := strings.Join(jangkauanArray, ",")

	infoArray := form.Value["infoLowongan[]"]
	infoLowongan := strings.Join(infoArray, ",")

	// Proses Jadwal Free
	jadwal := map[string][]string{}
	days := []string{"Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu", "Minggu"}
	for _, day := range days {
		jadwal[day] = form.Value["jadwal_"+day+"[]"]
	}
	jadwalJSON, _ := json.Marshal(jadwal)

	// Proses Ketersediaan
	ketersediaanMap := map[string]string{
		"Online":  c.FormValue("ketersediaanOnline"),
		"Offline": c.FormValue("ketersediaanOffline"),
	}
	ketersediaanJSON, _ := json.Marshal(ketersediaanMap)

	lamaran := models.Lamaran{
		PelamarID:        pelamarID,
		NamaLengkap:      c.FormValue("namaLengkap"),
		JenisKelamin:     c.FormValue("jenisKelamin"),
		NoWA:             c.FormValue("noHp"),
		AlamatDomisili:   c.FormValue("alamat"),
		KotaID:           uint(kotaID),
		KecamatanID:      uint(kecamatanID),
		ProgramStudi:     c.FormValue("programStudi"),
		Universitas:      c.FormValue("programStudi"), // Default to same field as ProgramStudi if not split
		JenjangDitempuh:  c.FormValue("jenjangDitempuh"),
		Semester:         c.FormValue("semester"),
		TargetJenjangID:  uint(targetJenjangID),
		TargetMapelID:    uint(targetMapelID),
		JangkauanWilayah: jangkauan,
		Ketersediaan:     string(ketersediaanJSON),
		JadwalFree:       string(jadwalJSON),
		FeeHarapan:       c.FormValue("feeHarapan"),
		MulaiMengajar:    c.FormValue("mulaiMengajar"),
		Pengalaman:       c.FormValue("pengalamanMengajar"),
		Kelebihan:        c.FormValue("kelebihanDiri"),
		Kekurangan:       c.FormValue("kekuranganDiri"),
		Prioritas:        c.FormValue("aktivitasPrioritas"),
		NamaOrtu:         c.FormValue("namaOrangTua"),
		NoHPOrtu:         c.FormValue("hpOrangTua"),
		InfoLowongan:     infoLowongan,
		TranskripURL:     transkripURL,
		CVURL:            cvURL,
		Status:           "Pending",
	}

	if err := config.DB.Create(&lamaran).Error; err != nil {
		log.Println("Database error:", err)
		return c.Redirect("/apply?error=Failed to save data")
	}

	// Check if BankSoal exists
	var count int64
	config.DB.Model(&models.BankSoalA{}).Where("jenis_pendidikan_id = ? AND mata_pelajaran_id = ? AND active = 'T'",
		lamaran.TargetJenjangID, lamaran.TargetMapelID).Count(&count)

	if count > 0 {
		return c.Redirect("/test/intro")
	}

	return c.Redirect("/dashboard?success=1")
}
