package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Admin represents an administrator account.
type Admin struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"size:100;uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	RoleID    uint           `gorm:"index" json:"role_id"`
	Role      Role           `gorm:"foreignKey:RoleID" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// HashPassword hashes plain-text password and stores it on the struct.
func (a *Admin) HashPassword(plain string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.Password = string(hashed)
	return nil
}

// CheckPassword compares a plain-text password against the stored hash.
func (a *Admin) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(plain))
	return err == nil
}

// --- CRUD Functions ---

// CreateAdmin inserts a new admin record into the database.
func CreateAdmin(db *gorm.DB, a *Admin) error {
	return db.Create(a).Error
}

// GetAdminByUsername retrieves an admin by their username.
func GetAdminByUsername(db *gorm.DB, username string) (*Admin, error) {
	var a Admin
	err := db.Where("username = ?", username).First(&a).Error
	return &a, err
}

// GetAdminByID retrieves an admin by their primary key.
func GetAdminByID(db *gorm.DB, id uint) (*Admin, error) {
	var a Admin
	err := db.First(&a, id).Error
	return &a, err
}

// GetAllAdmins returns all admin records.
func GetAllAdmins(db *gorm.DB) ([]Admin, error) {
	var list []Admin
	err := db.Order("created_at DESC").Find(&list).Error
	return list, err
}
