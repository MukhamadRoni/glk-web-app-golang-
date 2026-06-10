package main

import (
    "glk-web-app/config"
    "glk-web-app/models"
    "log"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Error loading .env file")
	}
    config.ConnectDB(
        &models.Admin{},
        &models.Pelamar{},
        &models.Lamaran{},
    )
    if err := config.DB.Exec("TRUNCATE TABLE lamarans CASCADE").Error; err != nil { log.Println(err) }
    if err := config.DB.Exec("TRUNCATE TABLE pelamars CASCADE").Error; err != nil { log.Println(err) }
    log.Println("Truncated")
}
