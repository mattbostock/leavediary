package model

import (
	"errors"
	"math"
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

func FindLeaveAllowance(id uint64) (l LeaveAllowance, err error) {
	res := db.First(&l, id)
	return l, res.Error
}

func FindLeaveRequest(id uint64) (l LeaveRequest, err error) {
	res := db.First(&l, id)
	return l, res.Error
}

func DeleteLeaveAllowance(id uint64) (err error) {
	err = db.Where("id = ?", id).Delete(&LeaveAllowance{}).Error
	return err
}

func DeleteLeaveRequest(id uint64) (err error) {
	err = db.Where("id = ?", id).Delete(&LeaveRequest{}).Error
	return err
}

func FindUser(id uint64) (user User, err error) {
	res := db.Preload("Jobs").First(&user, id)
	return user, res.Error
}

func FindJobFromExportSecret(secret string) (j Job, err error) {
	res := db.Where("export_secret = ?", secret).First(&j)
	return j, res.Error
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
		Find(&requests).Where("end_time >= ? AND job_id = ?", yearAgo, j.ID).Error

	return requests, err
}

func (l *LeaveAllowance) RemainingTime() (minutes int32, err error) {
	allotted, err := l.allottedTime()
	if err != nil {
		return 0, err
	}

	used, err := l.usedTime()
	if err != nil {
		return 0, err
	}

	return allotted - used, nil
}

func (l *LeaveAllowance) allottedTime() (minutes int32, err error) {
	// FIXME See if there's a neater way to eagerly load job data (aside from job ID)
	if err := db.Model(&l).Related(&l.Job).Error; err != nil {
		return minutes, err
	}

	if l.Job.ID == 0 {
		return minutes, errors.New("Job ID is zero")
	}

	err = db.Table("leave_allowances").
		Select("COALESCE(SUM(minutes),0)").
		Where("is_adjustment = ? AND start_time = ? AND end_time = ? AND job_id = ?", true, l.StartTime, l.EndTime, l.JobID).
		Where("deleted_at IS NULL OR deleted_at <= '0001-01-02'").
		Row().Scan(&minutes)
	minutes += l.Minutes

	jobStartAdjust := math.Min(1.0, l.EndTime.Sub(l.Job.StartTime).Minutes()/l.EndTime.Sub(l.StartTime).Minutes())
	minutes = int32(float64(minutes) * jobStartAdjust)

	// FIXME handle job leave date

	return minutes, err
}

// FIXME: take into account the job start date
func (l *LeaveAllowance) usedTime() (minutes int32, err error) {
	if l.JobID == 0 {
		return minutes, errors.New("Job ID is zero")
	}

	// Assumes leave doesn't span allocated periods
	err = db.Table("leave_requests").
		Select("COALESCE(SUM(minutes),0)").
		Where("job_id = ? AND start_time >= ? AND end_time <= ?", l.JobID, l.StartTime, l.EndTime).
		Where("deleted_at IS NULL OR deleted_at <= '0001-01-02'").
		Row().Scan(&minutes)

	return minutes, err
}

// IntersectsLeaveRequests checks if a leave allowance starts or ends part-way through
// any existing leave requests for the same job.
//
// This query can be used to ensure that changes to a leave allowance will not
// cause leave requests to span leave allowance periods, which is currently not supported.
func (a *LeaveAllowance) IntersectsLeaveRequest() (intersects LeaveRequest, _ error) {
	if a.JobID == 0 {
		return intersects, errors.New("Job ID is zero")
	}

	res := db.First(&intersects,
		`
		id <> ? AND
		(start_time < ? AND end_time > ?)
		OR
		(start_time < ? AND end_time > ?)
		AND job_id = ?`,

		a.ID,
		a.StartTime,
		a.StartTime,
		a.EndTime,
		a.EndTime,
		a.JobID)

	if res.Error == gorm.RecordNotFound {
		return intersects, nil
	}

	return intersects, res.Error
}

func (a *LeaveAllowance) OverlapsAnother() (overlaps LeaveAllowance, _ error) {
	if a.JobID == 0 {
		return overlaps, errors.New("Job ID is zero")
	}

	res := db.First(&overlaps,
		`id <> ? AND is_adjustment = ?
		AND start_time >= ? AND start_time < ?
		AND end_time > ? AND end_time <= ?
		AND job_id = ?`,

		a.ID,
		false,
		a.StartTime,
		a.EndTime,
		a.StartTime,
		a.EndTime,
		a.JobID)

	if res.Error == gorm.RecordNotFound {
		return overlaps, nil
	}

	return overlaps, res.Error
}

func (a *LeaveAllowance) Save() error {
	res := db.Save(a)
	return res.Error
}

func (r *LeaveRequest) After(t time.Time) bool {
	return r.StartTime.After(t)
}

func (r *LeaveRequest) Before(t time.Time) bool {
	return r.StartTime.Before(t)
}

func (r *LeaveRequest) Save() error {
	res := db.Save(r)
	return res.Error
}

func (r *LeaveRequest) FitsExistingAllowancePeriod() (result int, err error) {
	if r.JobID == 0 {
		return result, errors.New("Job ID is zero")
	}

	err = db.Table("leave_allowances").
		Select("COUNT(id)").
		Where("is_adjustment = ? AND start_time <= ? AND end_time >= ? AND job_id = ?", false, r.StartTime, r.EndTime, r.JobID).
		Where("deleted_at IS NULL OR deleted_at <= '0001-01-02'").
		Row().Scan(&result)

	return result, err
}
