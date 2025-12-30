package models

import (
	"sync"

	"github.com/jinzhu/gorm"
)

var (
	dbInstance *gorm.DB
	dbOnce     sync.Once
	dbError    error
)

func GetDB() (*gorm.DB, error) {
	dbOnce.Do(func() {
		dbInstance, dbError = InitDB()
	})
	return dbInstance, dbError
}

func GetDBConfig() (*DBConfig, error) {
	return LoadConfig()
}