package infra

import (
	"voltaserve/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDb() *gorm.DB {
	if db != nil {
		return db
	}
	var err error
	db, err = gorm.Open(postgres.Open(config.GetConfig().DatabaseUrl), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
