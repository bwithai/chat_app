package db

import (
	"chatapp/user"
	"chatapp/ws"
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
	//db.AutoMigrate(&ws.DbClient{})
	//db.AutoMigrate(&ws.DbRoom{})
	db.AutoMigrate(&ws.DbMessage{})

	err = db.Exec("DELETE FROM users").Error
	if err != nil {
		return nil, err
	}
	err = db.Exec("DELETE FROM db_messages").Error
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) GetDB() *gorm.DB {
	return d.db
}
