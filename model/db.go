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

	db.Model(&Job{}).AddIndex("job_start_time", "start_time")
	db.Model(&Job{}).AddIndex("job_end_time", "end_time")
	db.Model(&Job{}).AddIndex("job_deleted_at", "deleted_at")

	db.Model(&LeaveAllowance{}).AddIndex("leave_allowance_start_time", "start_time")
	db.Model(&LeaveAllowance{}).AddIndex("leave_allowance_end_time", "end_time")
	db.Model(&LeaveAllowance{}).AddIndex("leave_allowance_deleted_at", "deleted_at")

	db.Model(&LeaveRequest{}).AddIndex("leave_request_start_time", "start_time")
	db.Model(&LeaveRequest{}).AddIndex("leave_request_end_time", "end_time")
	db.Model(&LeaveRequest{}).AddIndex("leave_request_deleted_at", "deleted_at")

	db.Model(&User{}).AddIndex("user_email", "email")
	db.Model(&User{}).AddIndex("user_deleted_at", "deleted_at")
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
