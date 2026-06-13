package main
import (
	"log"
	"glk-web-app/config"
	"github.com/joho/godotenv"
)
func main() {
	godotenv.Load()
	config.ConnectDB()
	sql := `
CREATE OR REPLACE VIEW view_jadwal_pelamar AS
SELECT 
    l.id AS lamaran_id,
    l.pelamar_id,
    l.nama_lengkap,
    day_data.key AS hari,
    jam_data.value AS jam_prakiraan
FROM 
    lamarans l,
    LATERAL jsonb_each(l.jadwal_free::jsonb) AS day_data(key, value),
    LATERAL jsonb_array_elements_text(day_data.value) AS jam_data(value)
WHERE 
    l.deleted_at IS NULL;`
	err := config.DB.Exec(sql).Error
	if err != nil {
		log.Fatal(err)
	}
	log.Println("View view_jadwal_pelamar created successfully")
}

