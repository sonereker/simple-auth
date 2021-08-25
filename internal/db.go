package internal

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
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=false",
		"localhost",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		port,
	)
	return err, dsn
}
