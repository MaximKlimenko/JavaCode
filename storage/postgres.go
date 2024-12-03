package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/MaximKlimenko/JavaCode/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

func NewConnection(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	err = db.AutoMigrate(&models.Wallet{})
	if err != nil {
		log.Fatal("Failed to run migrations: ", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to initialize DB connection: %v", err)
	}
	sqlDB.SetMaxOpenConns(100)                // Максимальное количество соединений
	sqlDB.SetMaxIdleConns(50)                 // Максимальное количество простаивающих соединений
	sqlDB.SetConnMaxLifetime(time.Minute * 5) // Максимальное время жизни соединения
	return db, err
}
