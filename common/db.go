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
	err, dsn := prepareConnectionParams()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func prepareConnectionParams() (error, string) {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s port=%d",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
		port,
	)
	return err, dsn
}
