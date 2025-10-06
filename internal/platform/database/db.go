package database

import (
	"fmt"

	"github.com/rezbow/ecommerce/internal/platform/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	return db, err
}
