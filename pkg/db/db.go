package db

import (
	"task/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	DB *gorm.DB
}

func NewDB(config *config.Config) *Db {
	db, err := gorm.Open(postgres.Open(config.Db.Dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return &Db{DB: db}
}
