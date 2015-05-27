package model

import (
	"errors"
	"time"
)

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

func (u *User) CurrentJob() Job {
	// FIXME: Support multiple jobs per user in future; for now, just support one per user
	if len(u.Jobs) > 0 {
		return u.Jobs[0]
	}

	return Job{}
}

func (u *User) Save() error {
	res := db.Save(u)
	return res.Error
}

func (u *User) TZLocation() *time.Location {
	return time.FixedZone("User-defined", int(u.TZOffset))
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
	res := db.Preload("Jobs").First(&user, id)
	return user, res.Error
}
