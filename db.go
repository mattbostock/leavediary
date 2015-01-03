package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

func initDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "sqlite.db")

	if err != nil {
		panic(err)
	}

	// enable SQL logging
	db.LogMode(true)
	db.SetLogger(&debugLogger{})

	return &db
}

// debugLogger satisfies Gorm's logger interface
// so that we can log SQL queries at Logrus' debug level
type debugLogger struct{}

func (*debugLogger) Print(msg ...interface{}) {
	log.Debug(msg)
}
