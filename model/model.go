package model

import (
	"errors"
	"time"
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

type LeaveAllowance struct {
	ID           uint64 `gorm:"column:id; primary_key:yes"`
	Job          Job
	JobID        uint64    `gorm:"column:job_id"`
	StartTime    time.Time `sql:"DEFAULT:null"`
	EndTime      time.Time `sql:"DEFAULT:null"`
	Minutes      int32
	Description  string
	IsAdjustment bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type LeaveRequest struct {
	ID          uint64 `gorm:"column:id; primary_key:yes"`
	Job         Job
	JobID       uint64 `gorm:"column:job_id"`
	Minutes     uint32
	Description string
	StartTime   time.Time `sql:"DEFAULT:null"`
	EndTime     time.Time `sql:"DEFAULT:null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type User struct {
	ID       uint64 `gorm:"column:id; primary_key:yes"`
	Name     string `sql:"type:text;"`
	GitHubID uint64 `gorm:"column:github_id"`
	Email    string `sql:"type:text;"`
	TZOffset int16  `gorm:"column:tz_offset"` // time zone as seconds east of UTC
	Jobs     []Job

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func (u *User) UpdateOrCreate() error {
	if u.GitHubID == 0 {
		return errors.New("GitHub user ID was set to zero; cannot match")
	}

	res := db.Where(User{GitHubID: u.GitHubID}).FirstOrInit(u)

	if res.Error != nil {
		return res.Error
	}

	res = db.Save(u)
	return res.Error
}

func FindUser(id uint64) (user User, err error) {
	res := db.First(&user, id)
	return user, res.Error
}
