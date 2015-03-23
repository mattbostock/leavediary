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

	db.AutoMigrate(&Job{}, &LeaveYear{}, &Org{}, &User{})

	db.Model(&Job{}).AddIndex("id", "id")
	db.Model(&Job{}).AddIndex("deleted_at", "deleted_at")

	db.Model(&LeaveAlloc{}).AddIndex("start_date", "start_date")
	db.Model(&LeaveAlloc{}).AddIndex("deleted_at", "deleted_at")
	db.Model(&LeaveAlloc{}).AddIndex("added_by", "added_by")

	db.Model(&LeaveYear{}).AddIndex("start_date", "start_date")
	db.Model(&LeaveYear{}).AddIndex("deleted_at", "deleted_at")

	db.Model(&Org{}).AddIndex("id", "id")
	db.Model(&Org{}).AddIndex("deleted_at", "deleted_at")

	db.Model(&User{}).AddIndex("id", "id")
	db.Model(&User{}).AddIndex("email", "email")
	db.Model(&User{}).AddIndex("deleted_at", "deleted_at")
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
