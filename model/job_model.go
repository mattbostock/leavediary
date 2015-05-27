package model

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Job struct {
	ID              uint64 `gorm:"column:id; primary_key:yes"`
	ExportSecret    string
	User            User
	UserID          uint64    `gorm:"column:user_id"`
	StartTime       time.Time `sql:"DEFAULT:null"`
	EndTime         time.Time `sql:"DEFAULT:null"`
	EmployerName    string
	LeaveAllowances []LeaveAllowance
	LeaveRequests   []LeaveRequest

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func (j *Job) Save() error {
	res := db.Save(j)
	return res.Error
}

func (j *Job) CurrentLeaveAllowance() (l LeaveAllowance, _ error) {
	// FIXME See if there's a neater way to eagerly load user data (aside from user ID)
	if err := db.Model(&j).Related(&j.User).Error; err != nil {
		return l, err
	}

	if j.User.ID == 0 {
		return l, errors.New("No user data loaded")
	}

	now := time.Now().In(j.User.TZLocation())
	res := db.First(&l, "is_adjustment = ? AND start_time <= ? AND end_time >= ? AND job_id = ?",
		false, now, now, j.ID)

	if res.Error == gorm.RecordNotFound {
		return l, nil
	}

	return l, res.Error
}

func (j *Job) LeavePeriods() (l []LeaveAllowance, _ error) {
	res := db.Model(j).Related(&l)
	return l, res.Error
}

func (j *Job) RequestsLastYearAndFuture() (requests []LeaveRequest, err error) {
	// FIXME See if there's a neater way to eagerly load user data (aside from user ID)
	if err := db.Model(&j).Related(&j.User).Error; err != nil {
		return requests, err
	}

	if j.User.ID == 0 {
		return requests, errors.New("No user data loaded")
	}

	yearAgo := time.Now().AddDate(-1, 0, 0).In(j.User.TZLocation())
	err = db.Order("start_time DESC").Order("end_time DESC").
		Where("end_time >= ? AND job_id = ?", yearAgo, j.ID).
		Find(&requests).Error

	return requests, err
}

func FindJobFromExportSecret(secret string) (j Job, err error) {
	res := db.Where("export_secret = ?", secret).First(&j)
	return j, res.Error
}
