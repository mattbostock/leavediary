package model

import (
	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db  gorm.DB
	log *logrus.Logger
)

func InitDB(driver, database string) {
	var err error

	db, err = gorm.Open(driver, database)

	if err != nil {
		panic(err)
	}

	// enable SQL logging
	db.LogMode(true)
	db.SetLogger(&debugLogger{})
}

func SetLogger(l *logrus.Logger) {
	log = l
}

// debugLogger satisfies Gorm's logger interface
// so that we can log SQL queries at Logrus' debug level
type debugLogger struct{}

func (*debugLogger) Print(msg ...interface{}) {
	log.Debug(msg)
}
