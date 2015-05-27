package model

import (
	"errors"
	"time"
)

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

func FindLeaveRequest(id, jobID uint64) (l LeaveRequest, err error) {
	res := db.Where("job_id = ?", jobID).First(&l, id)
	return l, res.Error
}

func DeleteLeaveRequest(id, jobID uint64) (err error) {
	err = db.Where("id = ?", id).Where("job_id = ?", jobID).Delete(&LeaveRequest{}).Error
	return err
}
