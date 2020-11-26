package config

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	db *gorm.DB
)

func Connect() {
	// Please define your user name and password for my sql.
	// d, err := gorm.Open("mysql", "root:root@/myrestcrud?charset=utf8&parseTime=True&loc=Local")
	d, err := gorm.Open("sqlite3", "./mvc/sqlite-data/db.sqlite")
	if err != nil {
		panic(err)
	}
	db = d
}

func GetDB() *gorm.DB {
	return db
}
