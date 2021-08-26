package common

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strconv"
)

//NewDBConnection establishes a new DB connection
func NewDBConnection() (*gorm.DB, error) {
	dsn, err := prepareConnectionParams()
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func prepareConnectionParams() (string, error) {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return "", err
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s port=%d",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
		port,
	)
	return dsn, nil
}
