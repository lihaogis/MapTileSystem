package database

import (
	"fmt"
	"log"
	"map-tile-system/internal/model"
	"map-tile-system/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(
		&model.DataSource{},
		&model.ApiKey{},
		&model.CallLog{},
		&model.User{},
		&model.StatisticsSummary{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database connected and migrated successfully")
	return db, nil
}
