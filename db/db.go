package db

import (
	"chatapp/user"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

func NewDatabase() (*Database, error) {
	db, err := gorm.Open(sqlite.Open("go_chat.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// map User struct to create table in SQLite
	db.AutoMigrate(&user.User{})

	return &Database{db: db}, nil
}

func (d *Database) GetDB() *gorm.DB {
	return d.db
}
