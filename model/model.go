package model

import (
	"time"
)

type Org struct {
	ID   uint64 `gorm:"column:id; primary_key:yes"`
	Name string `sql:"type:text;"`
	// default number of days leave for organisation
	DefaultDays float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

type Employment struct {
	ID         uint64    `gorm:"column:id; primary_key:yes"`
	JoinDate   time.Time `sql:"DEFAULT:null"`
	EndDate    time.Time `sql:"DEFAULT:null"`
	LeaveYears []LeaveYear
	LeaveAlloc []LeaveAlloc
	Org        Org
	Managees   []User
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
	ID          uint64 `gorm:"column:id; primary_key:yes"`
	Name        string `sql:"type:text;"`
	Email       string `sql:"type:text;"`
	JobTitle    string `sql:"type:text;"`
	TimeZone    int
	Employments []Employment `gorm:"many2many:users_employments;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}
