package models

import (
	"gorm.io/gorm"
)

// Lamaran represents a job application submitted by a Pelamar.
type Lamaran struct {
	gorm.Model
	PelamarID      uint   `gorm:"not null" json:"pelamar_id"`
	NamaLengkap    string `gorm:"size:255;not null" json:"nama_lengkap"`
	JenisKelamin   string `gorm:"size:20;not null" json:"jenis_kelamin"`
	NoWA           string `gorm:"size:20;not null" json:"no_wa"`
	AlamatDomisili string `gorm:"type:text;not null" json:"alamat_domisili"`
	KotaID         uint   `gorm:"not null" json:"kota_id"`
	KecamatanID    uint   `gorm:"not null" json:"kecamatan_id"`

	ProgramStudi    string `gorm:"size:255;not null" json:"program_studi"`
	Universitas     string `gorm:"size:255;not null" json:"universitas"`
	JenjangDitempuh string `gorm:"size:100;not null" json:"jenjang_ditempuh"`
	Semester        string `gorm:"size:50;not null" json:"semester"`

	TargetJenjangID  uint   `gorm:"not null" json:"target_jenjang_id"`           // Refers to JenisPendidikan
	TargetMapelID    uint   `gorm:"not null" json:"target_mapel_id"`             // Refers to Mapel
	JangkauanWilayah string `gorm:"type:text;not null" json:"jangkauan_wilayah"` // Comma-separated kecamatan IDs
	Ketersediaan     string `gorm:"size:50;not null" json:"ketersediaan"`        // Online, Offline, Online & Offline
	JadwalFree       string `gorm:"type:text;not null" json:"jadwal_free"`       // JSON payload of schedule
	FeeHarapan       string `gorm:"size:255;not null" json:"fee_harapan"`
	MulaiMengajar    string `gorm:"size:100;not null" json:"mulai_mengajar"`

	Pengalaman string `gorm:"type:text" json:"pengalaman"`
	Kelebihan  string `gorm:"type:text" json:"kelebihan"`
	Kekurangan string `gorm:"type:text" json:"kekurangan"`
	Prioritas  string `gorm:"type:text" json:"prioritas"`

	NamaOrtu     string `gorm:"size:255;not null" json:"nama_ortu"`
	NoHPOrtu     string `gorm:"size:20;not null" json:"no_hp_ortu"`
	InfoLowongan string `gorm:"type:text" json:"info_lowongan"` // Comma-separated

	TranskripURL string `gorm:"type:text;not null" json:"transkrip_url"`
	CVURL        string `gorm:"type:text;not null" json:"cv_url"`

	Status string `gorm:"size:50;default:'Pending'" json:"status"` // Pending, Diterima, Ditolak, dll

	KoreksiNilai string `gorm:"type:text" json:"koreksi_nilai"` // JSON string for manual scoring

	// Relationships
	Pelamar       Pelamar         `gorm:"foreignKey:PelamarID" json:"pelamar"`
	Kota          Kota            `gorm:"foreignKey:KotaID" json:"kota"`
	Kecamatan     Kecamatan       `gorm:"foreignKey:KecamatanID" json:"kecamatan"`
	TargetJenjang JenisPendidikan `gorm:"foreignKey:TargetJenjangID" json:"target_jenjang"`
	TargetMapel   MataPelajaran   `gorm:"foreignKey:TargetMapelID" json:"target_mapel"`
}

// CheckIfPelamarHasApplied checks if a pelamar has already submitted an application.
func CheckIfPelamarHasApplied(db *gorm.DB, pelamarID uint) (bool, *Lamaran, error) {
	var lamaran Lamaran
	err := db.Where("pelamar_id = ?", pelamarID).Preload("TargetJenjang").Preload("TargetMapel").First(&lamaran).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil, nil
		}
		return false, nil, err
	}
	return true, &lamaran, nil
}
