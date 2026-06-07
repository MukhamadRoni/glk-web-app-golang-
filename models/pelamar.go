package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Pelamar represents a job applicant account.
type Pelamar struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Email     string         `gorm:"size:150;uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	CVURL     string         `gorm:"size:500" json:"cv_url"`
	Status    string         `gorm:"size:20;default:'pending'" json:"status"` // pending | accepted | rejected
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// HashPassword hashes plain-text password and stores it on the struct.
func (p *Pelamar) HashPassword(plain string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.Password = string(hashed)
	return nil
}

// CheckPassword compares a plain-text password against the stored hash.
func (p *Pelamar) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(plain))
	return err == nil
}

// --- CRUD Functions ---

// CreatePelamar inserts a new pelamar record into the database.
func CreatePelamar(db *gorm.DB, p *Pelamar) error {
	return db.Create(p).Error
}

// GetPelamarByEmail retrieves a pelamar by their email address.
func GetPelamarByEmail(db *gorm.DB, email string) (*Pelamar, error) {
	var p Pelamar
	err := db.Where("email = ?", email).First(&p).Error
	return &p, err
}

// GetPelamarByID retrieves a pelamar by their primary key.
func GetPelamarByID(db *gorm.DB, id uint) (*Pelamar, error) {
	var p Pelamar
	err := db.First(&p, id).Error
	return &p, err
}

// GetAllPelamar returns all pelamar records, ordered by newest first.
func GetAllPelamar(db *gorm.DB) ([]Pelamar, error) {
	var list []Pelamar
	err := db.Order("created_at DESC").Find(&list).Error
	return list, err
}

// UpdatePelamarStatus changes the status of a pelamar.
func UpdatePelamarStatus(db *gorm.DB, id uint, status string) error {
	return db.Model(&Pelamar{}).Where("id = ?", id).Update("status", status).Error
}

// DeletePelamar soft-deletes a pelamar record.
func DeletePelamar(db *gorm.DB, id uint) error {
	return db.Delete(&Pelamar{}, id).Error
}
