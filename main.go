package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"

	"glk-web-app/config"
	adminCtrl "glk-web-app/controllers/admin"
	adminAPIv1 "glk-web-app/controllers/admin/api/v1"
	pelamarCtrl "glk-web-app/controllers/pelamar"
	"glk-web-app/models"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("[WARN] .env not found, falling back to system environment")
	}
	config.ConnectRedis()

	// Ensure the log directory exists
	if err := os.MkdirAll("logs", 0755); err != nil {
	}

	// ── 2. Connect to database & auto-migrate models ──────────────────────
	config.ConnectDB(
		&models.Role{},
		&models.RoleMenu{},
		&models.Menu{},
		&models.Kota{},
		&models.Kecamatan{},
		&models.MataPelajaran{},
		&models.BankSoalA{},
		&models.BankSoalB{},
		&models.BankSoalC{},
		&models.Admin{},
		&models.Pelamar{},
		&models.Lamaran{},
		&models.ConfidenceScore{},
		&models.CompanyMCP{},
		&models.AIProfilingSkill{},
	)
	// Seed Mega Project Menus
	var menuCount int64
	config.DB.Model(&models.Menu{}).Count(&menuCount)
	if menuCount == 0 {
		// Root menus (Modules)
		recModule := models.Menu{Module: "Recruitment", Code: "MOD_RECRUITMENT", Name: "Recruitment", Icon: "bx bx-briefcase"}
		admModule := models.Menu{Module: "Admkar", Code: "MOD_ADMKAR", Name: "Admkar", Icon: "bx bx-id-card"}
		payModule := models.Menu{Module: "Payroll", Code: "MOD_PAYROLL", Name: "Payroll", Icon: "bx bx-money"}
		config.DB.Create(&recModule)
		config.DB.Create(&admModule)
		config.DB.Create(&payModule)

		menus := []models.Menu{
			// Recruitment Submenus
			{Module: "Recruitment", ParentID: &recModule.ID, Code: "MNU_REC_DASHBOARD", Name: "Dashboard", URL: "/admin/recruitment/dashboard", Icon: "bx bx-home-alt"},
			{Module: "Recruitment", ParentID: &recModule.ID, Code: "MNU_REC_MASTER", Name: "Master", Icon: "bx bx-data"},
			{Module: "Recruitment", ParentID: &recModule.ID, Code: "MNU_REC_TRANSAKSI", Name: "Transaksi", URL: "/admin/recruitment/transaksi", Icon: "bx bx-cart"},
			{Module: "Recruitment", ParentID: &recModule.ID, Code: "MNU_REC_LAPORAN", Name: "Laporan", URL: "/admin/recruitment/laporan", Icon: "bx bx-file"},
			{Module: "Recruitment", ParentID: &recModule.ID, Code: "MNU_REC_PERM", Name: "Permission", Icon: "bx bx-shield-quarter"},
			// Admkar Submenus
			{Module: "Admkar", ParentID: &admModule.ID, Code: "MNU_ADM_DASHBOARD", Name: "Dashboard", URL: "/admin/admkar/dashboard", Icon: "bx bx-home-alt"},
			{Module: "Admkar", ParentID: &admModule.ID, Code: "MNU_ADM_MASTER", Name: "Master", Icon: "bx bx-data"},
			{Module: "Admkar", ParentID: &admModule.ID, Code: "MNU_ADM_TRANSAKSI", Name: "Transaksi", URL: "/admin/admkar/transaksi", Icon: "bx bx-cart"},
			{Module: "Admkar", ParentID: &admModule.ID, Code: "MNU_ADM_LAPORAN", Name: "Laporan", URL: "/admin/admkar/laporan", Icon: "bx bx-file"},
			{Module: "Admkar", ParentID: &admModule.ID, Code: "MNU_ADM_PERM", Name: "Permission", Icon: "bx bx-shield-quarter"},
			// Payroll Submenus
			{Module: "Payroll", ParentID: &payModule.ID, Code: "MNU_PAY_DASHBOARD", Name: "Dashboard", URL: "/admin/payroll/dashboard", Icon: "bx bx-home-alt"},
			{Module: "Payroll", ParentID: &payModule.ID, Code: "MNU_PAY_MASTER", Name: "Master", Icon: "bx bx-data"},
			{Module: "Payroll", ParentID: &payModule.ID, Code: "MNU_PAY_TRANSAKSI", Name: "Transaksi", URL: "/admin/payroll/transaksi", Icon: "bx bx-cart"},
			{Module: "Payroll", ParentID: &payModule.ID, Code: "MNU_PAY_LAPORAN", Name: "Laporan", URL: "/admin/payroll/laporan", Icon: "bx bx-file"},
			{Module: "Payroll", ParentID: &payModule.ID, Code: "MNU_PAY_PERM", Name: "Permission", Icon: "bx bx-shield-quarter"},
		}

		for i := range menus {
			config.DB.Create(&menus[i])
		}

		// Master submenus for Recruitment
		var recMaster models.Menu
		config.DB.Where("code = ?", "MNU_REC_MASTER").First(&recMaster)
		config.DB.Create(&models.Menu{Module: "Recruitment", ParentID: &recMaster.ID, Code: "MNU_REC_WILAYAH", Name: "Wilayah", URL: "/admin/recruitment/master/wilayah"})
		config.DB.Create(&models.Menu{Module: "Recruitment", ParentID: &recMaster.ID, Code: "MNU_REC_BANK_SOAL", Name: "Bank Soal", URL: "/admin/recruitment/master/bank-soal"})
		config.DB.Create(&models.Menu{Module: "Recruitment", ParentID: &recMaster.ID, Code: "MNU_REC_MAPEL", Name: "Mata Pelajaran", URL: "/admin/recruitment/master/mapel"})

		// Permission submenus for Recruitment
		var recPerm models.Menu
		config.DB.Where("code = ?", "MNU_REC_PERM").First(&recPerm)
		config.DB.Create(&models.Menu{Module: "Recruitment", ParentID: &recPerm.ID, Code: "MNU_REC_ROLE", Name: "Role", URL: "/admin/permission/role"})
		config.DB.Create(&models.Menu{Module: "Recruitment", ParentID: &recPerm.ID, Code: "MNU_REC_MENU_ROLE", Name: "Menu Role", URL: "/admin/permission/menu-role"})
		config.DB.Create(&models.Menu{Module: "Recruitment", ParentID: &recPerm.ID, Code: "MNU_REC_USERS", Name: "Users", URL: "/admin/permission/users"})

		// Permission submenus for Admkar
		var admPerm models.Menu
		config.DB.Where("code = ?", "MNU_ADM_PERM").First(&admPerm)
		config.DB.Create(&models.Menu{Module: "Admkar", ParentID: &admPerm.ID, Code: "MNU_ADM_ROLE", Name: "Role", URL: "/admin/permission/role"})
		config.DB.Create(&models.Menu{Module: "Admkar", ParentID: &admPerm.ID, Code: "MNU_ADM_MENU_ROLE", Name: "Menu Role", URL: "/admin/permission/menu-role"})
		config.DB.Create(&models.Menu{Module: "Admkar", ParentID: &admPerm.ID, Code: "MNU_ADM_USERS", Name: "Users", URL: "/admin/permission/users"})
	}

	// Seed default role and admin user
	var roleCount int64
	config.DB.Model(&models.Role{}).Count(&roleCount)
	if roleCount == 0 {
		role := models.Role{Name: "Super Admin"}
		config.DB.Create(&role)

		// Grant all menus to Super Admin
		var allMenus []models.Menu
		config.DB.Find(&allMenus)
		for _, m := range allMenus {
			config.DB.Create(&models.RoleMenu{RoleID: role.ID, MenuCode: m.Code})
		}

		var adminCount int64
		config.DB.Model(&models.Admin{}).Count(&adminCount)
		if adminCount == 0 {
			admin := &models.Admin{
				Username: "admin",
				RoleID:   role.ID,
			}
			admin.HashPassword("admin")
			config.DB.Create(admin)
			log.Println("[INFO] Seeded default admin user (admin/admin) with Super Admin role")
		} else {
			// If admin exists but role is empty, update it
			config.DB.Model(&models.Admin{}).Where("username = ?", "admin").Update("role_id", role.ID)
		}
	}

	// ── 3. Session stores ────────────────────────────────────────────────
	adminSessionStore := session.New(session.Config{
		CookieName:     "admin_session",
		Expiration:     24 * time.Hour,
		CookieHTTPOnly: true,
	})

	pelamarSessionStore := session.New(session.Config{
		CookieName:     "pelamar_session",
		Expiration:     24 * time.Hour,
		CookieHTTPOnly: true,
	})

	// Inject session stores
	adminCtrl.InitStore(adminSessionStore)
	pelamarCtrl.InitStore(pelamarSessionStore)

	// ── 4. Start both servers concurrently ───────────────────────────────
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := newPelamarApp().Listen(":" + config.GetEnv("PELAMAR_PORT", "8081")); err != nil {
			log.Fatalf("[Pelamar] Server error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := newAdminApp().Listen(":" + config.GetEnv("ADMIN_PORT", "8082")); err != nil {
			log.Fatalf("[Admin] Server error: %v", err)
		}
	}()

	log.Println("┌─────────────────────────────────────────")
	log.Println("│  GLK Web App")
	log.Printf("│  Web Pelamar → http://localhost:%s", config.GetEnv("PELAMAR_PORT", "8081"))
	log.Printf("│  Web Admin   → http://localhost:%s", config.GetEnv("ADMIN_PORT", "8082"))
	log.Printf("│  Admin API   → http://localhost:%s/api/v1", config.GetEnv("ADMIN_PORT", "8082"))
	log.Println("└─────────────────────────────────────────")

	wg.Wait()
}

// ─────────────────────────────────────────────────────────────────────────────
// newPelamarApp builds the Fiber application for Web Pelamar (port 8081).
// ─────────────────────────────────────────────────────────────────────────────
func newPelamarApp() *fiber.App {
	engine := html.NewFileSystem(http.Dir("./views"), ".html")
	engine.Reload(true) // Enable hot-reload for templates
	engine.AddFunc("add", func(x, y int) int {
		return x + y
	})

	app := fiber.New(fiber.Config{
		AppName:   "GLK Web Pelamar",
		Views:     engine,
		BodyLimit: 50 * 1024 * 1024, // 50 MB
	})

	// ── Global Middleware ────────────────────────────────────────────────
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[Pelamar] ${time} | ${status} | ${latency} | ${method} ${path}\n",
	}))

	// ── Static Files (shared with admin app, same physical folder) ───────
	app.Static("/static", "./static", fiber.Static{
		Compress: true,
		MaxAge:   86400, // 1 day browser cache
	})

	// ── Public Routes ────────────────────────────────────────────────────
	app.Get("/", func(c *fiber.Ctx) error { return c.Redirect("/login") })
	app.Get("/login", pelamarCtrl.ShowLogin)
	app.Post("/login", pelamarCtrl.ProcessLogin)
	app.Get("/magic-link", pelamarCtrl.VerifyMagicLink)
	app.Get("/logout", pelamarCtrl.Logout)

	// ── Protected Routes (require pelamar session) ───────────────────────
	protected := app.Group("/", pelamarCtrl.AuthRequired)
	protected.Get("/dashboard", pelamarCtrl.ShowDashboard)
	protected.Get("/apply", pelamarCtrl.ShowApplyPage)
	protected.Post("/apply", pelamarCtrl.ProcessApply)

	// Test Akademik
	protected.Get("/test/intro", pelamarCtrl.ShowTestIntro)
	protected.Post("/test/start", pelamarCtrl.StartTest)
	protected.Get("/test/soal", pelamarCtrl.ShowTestSoal)
	protected.Post("/test/answer", pelamarCtrl.SaveAnswer)
	protected.Post("/test/finish", pelamarCtrl.FinishTest)

	// ── API Routes (require pelamar session) ─────────────────────────────
	pelamarAPI := app.Group("/api/v1", pelamarCtrl.AuthRequired)
	pelamarAPI.Get("/kota", adminCtrl.GetKotasList)
	pelamarAPI.Get("/kecamatan", adminCtrl.GetKecamatansList)
	pelamarAPI.Get("/jenis-pendidikan", adminCtrl.GetJenisPendidikanList)
	pelamarAPI.Get("/mapel", adminCtrl.GetMapelsList)

	// ── 404 Handler ──────────────────────────────────────────────────────
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("Halaman tidak ditemukan")
	})

	return app
}

// ─────────────────────────────────────────────────────────────────────────────
// newAdminApp builds the Fiber application for Web Admin (port 8082).
// Includes both web routes and API v1 routes.
// ─────────────────────────────────────────────────────────────────────────────
func newAdminApp() *fiber.App {
	engine := html.NewFileSystem(http.Dir("./views"), ".html")
	engine.Reload(true) // Enable hot-reload for templates
	engine.AddFunc("add", func(x, y int) int {
		return x + y
	})

	app := fiber.New(fiber.Config{
		AppName:   "GLK Web Admin",
		Views:     engine,
		BodyLimit: 50 * 1024 * 1024, // 50 MB
	})

	// ── Global Middleware ────────────────────────────────────────────────
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[Admin] ${time} | ${status} | ${latency} | ${method} ${path}\n",
	}))

	// ── Static Files (shared, same physical folder as pelamar app) ───────
	app.Static("/static", "./static", fiber.Static{
		Compress: true,
		MaxAge:   86400,
	})

	// ── Public Web Routes ────────────────────────────────────────────────
	app.Get("/", func(c *fiber.Ctx) error { return c.Redirect("/admin/login") })
	app.Get("/admin/login", adminCtrl.ShowLogin)
	app.Post("/admin/login", adminCtrl.ProcessLogin)
	app.Get("/admin/logout", adminCtrl.Logout)

	api := app.Group("/api/v1", adminCtrl.APIAuthRequired)

	// API Wilayah (Kota & Kecamatan)
	api.Get("/kota", adminCtrl.GetKotasList)
	api.Post("/kota", adminCtrl.CreateKota)
	api.Put("/kota/:id", adminCtrl.UpdateKota)
	api.Delete("/kota/:id", adminCtrl.DeleteKota)

	api.Get("/kecamatan", adminCtrl.GetKecamatansList)
	api.Post("/kecamatan", adminCtrl.CreateKecamatan)
	api.Put("/kecamatan/:id", adminCtrl.UpdateKecamatan)
	api.Delete("/kecamatan/:id", adminCtrl.DeleteKecamatan)

	// API Mata Pelajaran (Raw SQL)
	api.Get("/mapel", adminCtrl.GetMapelsList)
	api.Post("/mapel", adminCtrl.CreateMapel)
	api.Put("/mapel/:id", adminCtrl.UpdateMapel)
	api.Patch("/mapel/:id/active", adminCtrl.UpdateMapelActive)
	api.Delete("/mapel/:id", adminCtrl.DeleteMapel)

	// API Jenis Pendidikan
	api.Get("/jenis-pendidikan", adminCtrl.GetJenisPendidikanList)
	api.Post("/jenis-pendidikan", adminCtrl.CreateJenisPendidikan)
	api.Put("/jenis-pendidikan/:id", adminCtrl.UpdateJenisPendidikan)
	api.Patch("/jenis-pendidikan/:id/active", adminCtrl.UpdateJenisPendidikanActive)
	api.Delete("/jenis-pendidikan/:id", adminCtrl.DeleteJenisPendidikan)

	// API Bank Soal
	api.Get("/bank-soal", adminCtrl.GetBankSoalList)
	api.Get("/bank-soal/:id", adminCtrl.GetBankSoalDetail)
	api.Post("/bank-soal", adminCtrl.SaveBankSoal)
	api.Put("/bank-soal/:id", adminCtrl.SaveBankSoal)
	api.Patch("/bank-soal/:id/active", adminCtrl.UpdateBankSoalActive)
	api.Delete("/bank-soal/:id", adminCtrl.DeleteBankSoal)

	// ── Protected Web Routes ─────────────────────────────────────────────
	web := app.Group("/admin", adminCtrl.AuthRequired)
	web.Get("/dashboard", adminCtrl.ShowDashboard)
	web.Get("/recruitment/dashboard", adminCtrl.ShowRecruitmentDashboard)
	web.Get("/applicants", adminCtrl.ShowApplicants)
	web.Get("/applicants/:id", adminCtrl.ShowApplicant)

	web.Get("/recruitment/master/wilayah", adminCtrl.ShowWilayahPage)
	web.Get("/recruitment/master/mapel", adminCtrl.ShowMapelPage)
	web.Get("/recruitment/master/jenis-pendidikan", adminCtrl.ShowJenisPendidikanPage)
	web.Get("/recruitment/master/bank-soal", adminCtrl.ShowBankSoalPage)
	web.Get("/recruitment/master/bank-soal/form", adminCtrl.ShowBankSoalFormPage)
	web.Get("/recruitment/master/bank-soal/form/:id", adminCtrl.ShowBankSoalFormPage)
	web.Get("/recruitment/master/confidence-score", adminCtrl.ShowConfidenceScorePage)
	web.Get("/ai/mcp", adminCtrl.ShowMCPPage)
	web.Get("/ai/profiling", adminCtrl.ShowProfilingPage)
	web.Get("/ai/chat", adminCtrl.ShowChatPage)

	// Recruitment Transactions
	web.Get("/recruitment/transaksi/pelamar", adminCtrl.ShowRecruitmentPelamar)

	// Permission Routes
	permission := web.Group("/permission")
	permission.Get("/role", adminCtrl.ShowRoles)
	permission.Post("/role", adminCtrl.ProcessRole)
	permission.Get("/menu-role", adminCtrl.ShowMenuRoles)
	permission.Post("/menu-role", adminCtrl.ProcessMenuRole)
	permission.Get("/users", adminCtrl.ShowUsers)
	permission.Post("/users", adminCtrl.ProcessUser)

	// ── API v1 Routes (JSON, protected via APIAuthRequired) ──────────────
	// ── API v1 Routes (JSON, protected via APIAuthRequired) ──────────────

	// Applicants API
	api.Get("/applicants", adminAPIv1.ListApplicants)
	api.Get("/applicants/:id", adminAPIv1.GetApplicant)
	api.Patch("/applicants/:id/status", adminAPIv1.UpdateApplicantStatus)
	api.Delete("/applicants/:id", adminAPIv1.DeleteApplicant)

	// Recruitment Transactions API
	api.Get("/recruitment/pelamar", adminAPIv1.ListRecruitmentPelamar)
	api.Get("/recruitment/pelamar/:id", adminAPIv1.GetRecruitmentPelamarDetail)
	api.Patch("/recruitment/pelamar/:id/correction", adminAPIv1.UpdateRecruitmentCorrection)
	api.Patch("/recruitment/pelamar/:id/status", adminAPIv1.UpdateRecruitmentStatus)

	// Master Confidence Score API
	api.Get("/confidence-score", adminCtrl.GetConfidenceScoreList)
	api.Post("/confidence-score", adminCtrl.CreateConfidenceScore)
	api.Put("/confidence-score/:id", adminCtrl.UpdateConfidenceScore)
	api.Delete("/confidence-score/:id", adminCtrl.DeleteConfidenceScore)

	// AI MCP API
	api.Get("/ai/mcp", adminCtrl.GetMCPList)
	api.Post("/ai/mcp", adminCtrl.CreateMCP)
	api.Put("/ai/mcp/:id", adminCtrl.UpdateMCP)
	api.Delete("/ai/mcp/:id", adminCtrl.DeleteMCP)
	api.Post("/ai/mcp/upload-proxy", adminCtrl.UploadMCPProxy)

	// AI Profiling API
	api.Get("/ai/profiling", adminCtrl.GetProfilingList)
	api.Post("/ai/profiling", adminCtrl.CreateProfiling)
	api.Put("/ai/profiling/:id", adminCtrl.UpdateProfiling)
	api.Delete("/ai/profiling/:id", adminCtrl.DeleteProfiling)
	api.Post("/ai/profiling/upload-proxy", adminCtrl.UploadProfilingProxy)

	// AI Chat API
	api.Post("/ai/chat", adminCtrl.ProcessChat)

	// ── 404 Handler ──────────────────────────────────────────────────────
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("Halaman tidak ditemukan")
	})

	return app
}
