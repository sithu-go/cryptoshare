package ds

import (
	"cryptoshare/model"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func LoadDB() (*gorm.DB, error) {
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASS")
	name := os.Getenv("MYSQL_NAME")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(err.Error(), "DDDS", dsn)
		return nil, err
	}

	log.Println("Successfully connected to MySQL")

	// migrate DB
	err = db.AutoMigrate(
		&model.Bank{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
