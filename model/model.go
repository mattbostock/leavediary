package model

import (
	"errors"
	"time"
)

type Job struct {
	ID         uint64    `gorm:"column:id; primary_key:yes"`
	JoinDate   time.Time `sql:"DEFAULT:null"`
	EndDate    time.Time `sql:"DEFAULT:null"`
	LeaveYears []LeaveYear
	LeaveAlloc []LeaveAlloc
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
}

type LeaveAlloc struct {
	ID          uint64 `gorm:"column:id; primary_key:yes"`
	AddedBy     User
	Days        float64
	Description string
	StartDate   time.Time `sql:"DEFAULT:null"`
}

type LeaveYear struct {
	ID        uint64    `gorm:"column:id; primary_key:yes"`
	StartDate time.Time `sql:"DEFAULT:null"`
	Days      float64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type User struct {
	ID        uint64 `gorm:"column:id; primary_key:yes"`
	Name      string `sql:"type:text;"`
	GitHubID  uint64
	Email     string `sql:"type:text;"`
	TimeZone  int
	Jobs      []Job
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
