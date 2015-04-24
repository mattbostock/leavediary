package model

import (
	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
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

	db.AutoMigrate(&Job{}, &LeaveAllowance{}, &LeaveRequest{}, &User{})

	db.Model(&Job{}).AddIndex("start_time", "start_time")
	db.Model(&Job{}).AddIndex("end_time", "end_time")
	db.Model(&Job{}).AddIndex("deleted_at", "deleted_at")

	db.Model(&LeaveAllowance{}).AddIndex("start_time", "start_time")
	db.Model(&LeaveAllowance{}).AddIndex("end_time", "end_time")
	db.Model(&LeaveAllowance{}).AddIndex("deleted_at", "deleted_at")

	db.Model(&LeaveRequest{}).AddIndex("start_time", "start_time")
	db.Model(&LeaveRequest{}).AddIndex("end_time", "end_time")
	db.Model(&LeaveRequest{}).AddIndex("deleted_at", "deleted_at")

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
