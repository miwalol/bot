package utils

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	url2 "net/url"
	"os"
)

func DbConnect() *gorm.DB {
	url, err := url2.Parse(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to parse \"DATABASE_URL\"!")
	}
	host := url.Host
	user := url.User.Username()
	password, passwordSet := url.User.Password()
	database := url.Path[1:]

	if host == "" || user == "" || password == "" || !passwordSet || database == "" {
		log.Fatal("Cannot connect to database, some environment variables are not set")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s", host, user, password, database,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			NoLowerCase:   true,
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to database: %s", err))
	}

	return db
}
