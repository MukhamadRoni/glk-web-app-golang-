package pelamar

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// store is the shared session store for the pelamar app.
// It is initialised once by InitStore() and used by all handlers.
var store *session.Store

// InitStore sets the session store that auth and middleware handlers will use.
// Call this once during app bootstrap in main.go.
func InitStore(s *session.Store) {
	store = s
}

// AuthRequired is a middleware that checks for an active pelamar session.
// Unauthenticated requests are redirected to the login page.
func AuthRequired(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil || sess.Get("pelamar_id") == nil {
		return c.Redirect("/login")
	}
	// Inject user info into locals so views can access it
	c.Locals("pelamar_id", sess.Get("pelamar_id"))
	c.Locals("pelamar_name", sess.Get("pelamar_name"))
	return c.Next()
}

// contextData returns common template data merged with the user's session.
func contextData(c *fiber.Ctx, extra fiber.Map) fiber.Map {
	data := fiber.Map{
		"Username":  c.Locals("pelamar_name"),
		"Timestamp": time.Now(),
	}
	for k, v := range extra {
		data[k] = v
	}
	return data
}
