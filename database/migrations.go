package database

import (
	"github.com/Rizal-Nurochman/database/entities"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&entities.User{},
		&entities.RefreshToken{},
	)
	
	if err != nil {
		return err
	}

	return nil
}