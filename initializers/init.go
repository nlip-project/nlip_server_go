package initializers

import (
	"log"
	"nlip/models"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb() {
	var err error
	DB, err = gorm.Open(postgres.Open(os.Getenv("DB_URL")), &gorm.Config{})

	if err != nil {
		log.Fatal("Could not connect to the database!")
	}
}

func SyncDB() {
	DB.AutoMigrate(&models.User{})
}

func InitEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
