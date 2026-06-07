package admin

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"glk-web-app/models"
)

// store is the shared session store for the admin app.
var store *session.Store

// InitStore sets the session store for admin handlers and middleware.
// Call once during app bootstrap in main.go.
func InitStore(s *session.Store) {
	store = s
}

// AuthRequired is a middleware that protects admin web routes.
// Unauthenticated requests are redirected to /admin/login.
func AuthRequired(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil || sess.Get("admin_id") == nil {
		return c.Redirect("/admin/login")
	}
	c.Locals("admin_id", sess.Get("admin_id"))
	c.Locals("admin_username", sess.Get("admin_username"))
	return c.Next()
}

// APIAuthRequired is a middleware that protects admin API v1 routes.
// Returns JSON 401 instead of a redirect for unauthenticated API calls.
func APIAuthRequired(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil || sess.Get("admin_id") == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}
	c.Locals("admin_id", sess.Get("admin_id"))
	c.Locals("admin_username", sess.Get("admin_username"))
	return c.Next()
}

// contextData merges common template variables (from session) with page-specific data.
func contextData(c *fiber.Ctx, extra fiber.Map) fiber.Map {
	data := fiber.Map{
		"Username":  c.Locals("admin_username"),
		"Timestamp": time.Now(),
        "Menus":     models.AdminMenus,
	}
	for k, v := range extra {
		data[k] = v
	}
	return data
}
