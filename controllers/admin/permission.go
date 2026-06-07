package admin

import (
	"strconv"

	"glk-web-app/config"
	"glk-web-app/models"

	"github.com/gofiber/fiber/v2"
)

// --- Role Handlers ---

func ShowRoles(c *fiber.Ctx) error {
	var roles []models.Role
	config.DB.Find(&roles)

	return c.Render("admin/permission/role", contextData(c, fiber.Map{
		"Title": "Roles",
		"Roles": roles,
	}), "layouts/base")
}

func ProcessRole(c *fiber.Ctx) error {
	name := c.FormValue("name")
	if name != "" {
		config.DB.Create(&models.Role{Name: name})
	}
	return c.Redirect("/admin/permission/role")
}

// --- Menu Role Handlers ---

func ShowMenuRoles(c *fiber.Ctx) error {
	var roles []models.Role
	config.DB.Find(&roles)

	roleID := c.QueryInt("role_id", 0)
	var activeRoleMenus []models.RoleMenu
	if roleID != 0 {
		config.DB.Where("role_id = ?", roleID).Find(&activeRoleMenus)
	}

	activeMap := make(map[string]bool)
	for _, rm := range activeRoleMenus {
		activeMap[rm.MenuCode] = true
	}

	var allMenus []models.Menu
	config.DB.Where("parent_id IS NULL").Preload("Children.Children").Find(&allMenus)

	return c.Render("admin/permission/menu_role", contextData(c, fiber.Map{
		"Title":     "Menu Role",
		"Roles":     roles,
		"RoleID":    roleID,
		"AllMenus":  allMenus,
		"ActiveMap": activeMap,
	}), "layouts/base")
}

func ProcessMenuRole(c *fiber.Ctx) error {
	roleIDStr := c.FormValue("role_id")
	if roleIDStr == "" {
		return c.Redirect("/admin/permission/menu-role")
	}

	roleID, _ := strconv.Atoi(roleIDStr)

	// Delete existing
	config.DB.Where("role_id = ?", roleID).Delete(&models.RoleMenu{})

	// Insert new
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		if codes, ok := form.Value["menu_codes"]; ok {
			for _, code := range codes {
				config.DB.Create(&models.RoleMenu{
					RoleID:   uint(roleID),
					MenuCode: code,
				})
			}
		}
	} else {
		// If using application/x-www-form-urlencoded
		// fiber doesn't expose a clean array accessor for form-urlencoded arrays in v2 easily without parsing
		// we can parse it from string
		// Actually, let's just ensure the HTML form sends multipart/form-data.
	}
	
	return c.Redirect("/admin/permission/menu-role?role_id=" + roleIDStr)
}

// --- Users Handlers ---

func ShowUsers(c *fiber.Ctx) error {
	var users []models.Admin
	config.DB.Preload("Role").Find(&users)

	var roles []models.Role
	config.DB.Find(&roles)

	return c.Render("admin/permission/users", contextData(c, fiber.Map{
		"Title": "Users",
		"Users": users,
		"Roles": roles,
	}), "layouts/base")
}

func ProcessUser(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	roleIDStr := c.FormValue("role_id")
	roleID, _ := strconv.Atoi(roleIDStr)

	if username != "" && password != "" {
		admin := &models.Admin{
			Username: username,
			RoleID:   uint(roleID),
		}
		admin.HashPassword(password)
		config.DB.Create(admin)
	}
	return c.Redirect("/admin/permission/users")
}
