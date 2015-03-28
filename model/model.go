package model

import (
	"errors"
	"time"
)

type Job struct {
	ID              uint64 `gorm:"column:id; primary_key:yes"`
	ExportSecret    string
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
	ID        uint64    `gorm:"column:id; primary_key:yes"`
	StartTime time.Time `sql:"DEFAULT:null"`
	EndTime   time.Time `sql:"DEFAULT:null"`
	Minutes   int32
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type LeaveRequest struct {
	ID          uint64 `gorm:"column:id; primary_key:yes"`
	Minutes     uint32
	Description string
	StartTime   time.Time `sql:"DEFAULT:null"`
	EndTime     time.Time `sql:"DEFAULT:null"`
}

type User struct {
	ID       uint64 `gorm:"column:id; primary_key:yes"`
	Name     string `sql:"type:text;"`
	GitHubID uint64
	Email    string `sql:"type:text;"`
	TZOffset int16  // time zone as seconds east of UTC
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
