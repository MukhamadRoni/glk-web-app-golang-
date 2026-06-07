package models

import (
	"time"

	"gorm.io/gorm"
)

// Wilayah represents region/area master data
type Wilayah struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Code      string         `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// MataPelajaran represents subjects master data
type MataPelajaran struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Code      string         `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BankSoal represents a question bank item
type BankSoal struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	MataPelajaranID uint           `gorm:"index;not null" json:"mata_pelajaran_id"`
	MataPelajaran   MataPelajaran  `gorm:"foreignKey:MataPelajaranID" json:"mata_pelajaran"`
	Question        string         `gorm:"type:text;not null" json:"question"`
	OptionsJSON     string         `gorm:"type:text;not null" json:"options_json"` // JSON array of options
	CorrectAnswer   string         `gorm:"size:255;not null" json:"correct_answer"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
