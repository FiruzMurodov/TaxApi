package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"taxApi/internal/models"
)

func InitConnectionToDb(config *models.Config) *gorm.DB {
	dbrUri := "host=" + config.Db.Host + " port=" + config.Db.Port + " user=" + config.Db.User + " password=" + config.Db.Password + " dbname=" + config.Db.Database
	db, err := gorm.Open(postgres.Open(dbrUri), &gorm.Config{})
	if err != nil {

		return nil
	}

	return db
}
