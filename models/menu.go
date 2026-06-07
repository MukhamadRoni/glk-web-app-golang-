package models

import (
	"gorm.io/gorm"
)

// Menu represents a dynamic database-driven menu for the mega project.
type Menu struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Module    string         `gorm:"size:100;not null;index" json:"module"` // Grouping e.g. Recruitment, Admkar
	ParentID  *uint          `gorm:"index" json:"parent_id"`
	Code      string         `gorm:"size:100;uniqueIndex;not null" json:"code"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	URL       string         `gorm:"size:255;default:'#'" json:"url"`
	Icon      string         `gorm:"size:100" json:"icon"`
	OrderNum  int            `gorm:"default:0" json:"order_num"`
	Children  []Menu         `gorm:"foreignKey:ParentID" json:"children"`
	CreatedAt int64          `gorm:"autoCreateTime" json:"-"`
	UpdatedAt int64          `gorm:"autoUpdateTime" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BuildMenuTree builds a nested menu tree from a flat list of menus.
func BuildMenuTree(menus []Menu, parentID *uint) []Menu {
	var tree []Menu
	for _, m := range menus {
		if (m.ParentID == nil && parentID == nil) || (m.ParentID != nil && parentID != nil && *m.ParentID == *parentID) {
			m.Children = BuildMenuTree(menus, &m.ID)
			tree = append(tree, m)
		}
	}
	return tree
}

// FilterMenus recursively filters the given menu tree based on allowed menu codes.
func FilterMenus(menus []Menu, allowedCodes map[string]bool) []Menu {
	var filtered []Menu
	for _, m := range menus {
		// Process children first
		var children []Menu
		if len(m.Children) > 0 {
			children = FilterMenus(m.Children, allowedCodes)
		}
		
		// Include this menu if it has children that are allowed OR if it's explicitly allowed
		if len(children) > 0 || allowedCodes[m.Code] {
			copied := m
			copied.Children = children
			filtered = append(filtered, copied)
		}
	}
	return filtered
}
