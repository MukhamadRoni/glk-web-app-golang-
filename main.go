package main

import (
	"log"
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
	// ── 1. Load environment variables ────────────────────────────────────
	if err := godotenv.Load(); err != nil {
		log.Println("[WARN] .env not found, falling back to system environment")
	}

	// ── 2. Connect to database & auto-migrate models ──────────────────────
	config.ConnectDB(
		&models.Admin{},
		&models.Pelamar{},
	)

	// ── 3. Shared session store (in-memory; swap for Redis/DB in production) ─
	sessionStore := session.New(session.Config{
		Expiration:     24 * time.Hour,
		CookieHTTPOnly: true,
	})

	// Inject session store into both controller packages
	adminCtrl.InitStore(sessionStore)
	pelamarCtrl.InitStore(sessionStore)

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
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		AppName: "GLK Web Pelamar",
		Views:   engine,
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
	app.Get("/register", pelamarCtrl.ShowRegister)
	app.Post("/register", pelamarCtrl.ProcessRegister)
	app.Get("/logout", pelamarCtrl.Logout)

	// ── Protected Routes (require pelamar session) ───────────────────────
	protected := app.Group("/", pelamarCtrl.AuthRequired)
	protected.Get("/dashboard", pelamarCtrl.ShowDashboard)

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
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		AppName: "GLK Web Admin",
		Views:   engine,
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

	// ── Protected Web Routes ─────────────────────────────────────────────
	web := app.Group("/admin", adminCtrl.AuthRequired)
	web.Get("/dashboard", adminCtrl.ShowDashboard)
	web.Get("/applicants", adminCtrl.ShowApplicants)
	web.Get("/applicants/:id", adminCtrl.ShowApplicant)

	// ── API v1 Routes (JSON, protected via APIAuthRequired) ──────────────
	api := app.Group("/api/v1", adminCtrl.APIAuthRequired)

	// Applicants API
	api.Get("/applicants", adminAPIv1.ListApplicants)
	api.Get("/applicants/:id", adminAPIv1.GetApplicant)
	api.Patch("/applicants/:id/status", adminAPIv1.UpdateApplicantStatus)
	api.Delete("/applicants/:id", adminAPIv1.DeleteApplicant)

	// ── 404 Handler ──────────────────────────────────────────────────────
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("Halaman tidak ditemukan")
	})

	return app
}
