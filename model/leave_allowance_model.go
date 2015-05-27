package model

import (
	"errors"
	"math"
	"time"

	"github.com/jinzhu/gorm"
)

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

func FindLeaveAllowance(id, jobID uint64) (l LeaveAllowance, err error) {
	res := db.Where("job_id = ?", jobID).First(&l, id)
	return l, res.Error
}

func DeleteLeaveAllowance(id, jobID uint64) (err error) {
	err = db.Where("id = ?", id).Where("job_id = ?", jobID).Delete(&LeaveAllowance{}).Error
	return err
}
