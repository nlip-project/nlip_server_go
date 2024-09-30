package initializers

import "nlip/models"

func SyncDB() {
	DB.AutoMigrate(&models.User{})
}
