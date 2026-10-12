package config

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConnectDB membuka koneksi MySQL lewat GORM dan mengatur connection pool.
func ConnectDB() (*gorm.DB, error) {

	// parseTime=true  -> kolom DATETIME dibaca sebagai time.Time
	// charset=utf8mb4 -> mendukung Unicode penuh
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		Env("DB_USERNAME", "root"),
		Env("DB_PASSWORD", ""),
		Env("DB_HOST", "mysql"),
		Env("DB_PORT", "3306"),
		Env("DB_DATABASE", ""),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxLifetime(3 * time.Minute)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(10)

	return db, nil
}
