package database

import (
	"fmt"

	"github.com/mateusrlopez/go-market/models"
	"github.com/mateusrlopez/go-market/settings"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetConnection() *gorm.DB {
	db, err := gorm.Open(postgres.Open(formatDSN()), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error opening connection with postgres database: %s", err)
		return nil
	}

	db.AutoMigrate(&models.Admin{}, &models.User{})

	return db
}

func formatDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		settings.Settings.Database.Host,
		settings.Settings.Database.UserName,
		settings.Settings.Database.Password,
		settings.Settings.Database.Name,
		settings.Settings.Database.Port,
	)
}
