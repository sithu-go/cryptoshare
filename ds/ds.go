package ds

import (
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

var DB *gorm.DB

type DataSource struct {
	DB  *gorm.DB
	RDB *redis.Client
}

func NewDataSource() (*DataSource, error) {
	db, err := LoadDB()
	if err != nil {
		return nil, err
	}

	DB = db

	go AddDefaultAdmin()

	rdb, err := LoadRDB()
	if err != nil {
		return nil, err
	}

	return &DataSource{
		DB:  db,
		RDB: rdb,
	}, nil
}
