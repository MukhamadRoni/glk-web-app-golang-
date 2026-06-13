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
	Requirements      string          `gorm:"type:text" json:"requirements"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	DeletedAt         gorm.DeletedAt  `gorm:"index" json:"-"`
	Active            string          `gorm:"size:1;default:'T'" json:"active"`
}

// BankSoalA represents the header/version of the question bank
type BankSoalA struct {
	ID                uint            `gorm:"primarykey" json:"id"`
	JenisPendidikanID uint            `gorm:"index" json:"jenis_pendidikan_id"`
	JenisPendidikan   JenisPendidikan `gorm:"foreignKey:JenisPendidikanID" json:"jenis_pendidikan"`
	MataPelajaranID   uint            `gorm:"index" json:"mata_pelajaran_id"`
	MataPelajaran     MataPelajaran   `gorm:"foreignKey:MataPelajaranID" json:"mata_pelajaran"`
	Title             string          `gorm:"size:255;not null" json:"title"`
	DurasiPengerjaan  int             `gorm:"default:20" json:"durasi_pengerjaan"` // Durasi dalam menit
	Version           int             `gorm:"not null;default:1" json:"version"`
	Active            string          `gorm:"size:1;default:'T'" json:"active"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	DeletedAt         gorm.DeletedAt  `gorm:"index" json:"-"`
	BankSoalBs        []BankSoalB     `gorm:"foreignKey:BankSoalAID" json:"questions"`
}

// BankSoalB represents the questions
type BankSoalB struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	BankSoalAID   uint           `gorm:"index;not null" json:"bank_soal_a_id"`
	QuestionType  string         `gorm:"size:50;not null" json:"question_type"` // e.g., 'MULTIPLE_CHOICE', 'ESSAY'
	QuestionText  string         `gorm:"type:text;not null" json:"question_text"`
	QuestionImage string         `gorm:"size:255" json:"question_image"`
	OrderIndex    int            `gorm:"default:0" json:"order_index"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	BankSoalCs    []BankSoalC    `gorm:"foreignKey:BankSoalBID" json:"options"`
}

// BankSoalC represents the answers/options
type BankSoalC struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	BankSoalBID   uint           `gorm:"index;not null" json:"bank_soal_b_id"`
	OptionText    string         `gorm:"type:text" json:"option_text"`
	OptionImage   string         `gorm:"size:255" json:"option_image"`
	IsCorrect     string         `gorm:"size:1;default:'F'" json:"is_correct"` // 'T' or 'F'
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// ConfidenceScore represents scoring categories for test results
type ConfidenceScore struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Color     string         `gorm:"size:20;not null" json:"color"` // Hex format e.g. #FF0000
	MinScore  int            `gorm:"not null" json:"min_score"`
	MaxScore  int            `gorm:"not null" json:"max_score"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
