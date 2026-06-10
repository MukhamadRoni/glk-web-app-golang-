package models

import (
	"time"

	"gorm.io/gorm"
)

// Kota represents City master data
type Kota struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Kecamatan represents District master data
type Kecamatan struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	KotaID    uint           `gorm:"index;not null" json:"kota_id"`
	Kota      Kota           `gorm:"foreignKey:KotaID" json:"kota"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// JenisPendidikan represents education types
type JenisPendidikan struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	JenisPendidikan string         `gorm:"size:50;not null" json:"jenis_pendidikan"`
	Name            string         `gorm:"size:100;not null" json:"name"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Active          string         `gorm:"size:1;default:'T'" json:"active"`
}

// MataPelajaran represents subjects master data
type MataPelajaran struct {
	ID                uint            `gorm:"primarykey" json:"id"`
	JenisPendidikanID uint            `gorm:"index" json:"jenis_pendidikan_id"`
	JenisPendidikan   JenisPendidikan `gorm:"foreignKey:JenisPendidikanID" json:"jenis_pendidikan"`
	Code              string          `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Name              string          `gorm:"size:100;not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Active    string         `gorm:"size:1;default:'T'" json:"active"`
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
