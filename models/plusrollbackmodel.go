package models

import (
	"errors"
	//"fmt"
	"mint-platform/platform/db"
	"mint-platform/platform/forms"
	"strconv"
	"time"
)

type PlusapiRollback_workorder struct {
	Id              int       `xorm:"id int(11) pk not null"`
	UserName        string    `xorm:"user_name varchar(50) not null"`
	TaskName        string    `xorm:"task_name varchar(100) not null"`
	EnvironmentType int       `xorm:"environment_type smallint(6) not null"`
    GroupName       string    `xorm:"group_name varchar(20) not null"`
	Version         string    `xorm:"rollback_version varchar(100) not null"`
	Status          int       `xorm:"rollback_status smallint(6) not null"`
	DateCreate      time.Time `xorm:"created"`
	DateCancel      time.Time `xorm:"date_cancel"`
	DateRollback    time.Time `xorm:"date_rollback"`
}

type PlusRollbackModel struct{}

func (r *PlusRollbackModel) TableName() string {
	return "plusapi_rollback_workorder"
}

func (r *PlusRollbackModel) GetRollbackList(page int, pageSize int) (rollbacks []PlusapiRollback_workorder, err error) {
	engine := db.GetDB("deploy_online")
	//engine.Sync2(new(Plusapi_Rollback_workorder))
	startNum := (page - 1) * pageSize
	err = engine.In("rollback_status", 1, 2).Desc("id").Limit(pageSize, startNum).Find(&rollbacks)
	if err != nil {
		return rollbacks, err
	}
	return rollbacks, nil
}

func (r *PlusRollbackModel) GetRollbackTotal() (num int64, err error) {
	engine := db.GetDB("deploy_online")
	rollback := new(PlusapiRollback_workorder)
	num, err = engine.In("rollback_status", 1, 2).Count(rollback)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (r *PlusRollbackModel) AddWorker(form forms.PlusProjectForm) (rollback PlusapiRollback_workorder, err error) {
	engine := db.GetDB("deploy_online")
	rollback.UserName = form.UserName
	rollback.TaskName = form.TaskName
	rollback.EnvironmentType, _ = strconv.Atoi(form.EnvironmentType)
	rollback.GroupName = form.Group
	rollback.Version = form.Version
	rollback.Status = 1

	affected, err := engine.Insert(&rollback)
	if err != nil {
		return rollback, err
	}
	if affected > 0 {
		return rollback, nil
	}
	return rollback, errors.New("添加回滚工单失败！")
}

func (r *PlusRollbackModel) CancelWorker(orderId int) (rollback PlusapiRollback_workorder, err error) {
	engine := db.GetDB("deploy_online")

	has, err := engine.Where("id = ?", orderId).Get(&rollback)
	if err != nil {
		return rollback, err
	}
	if has {
		rollback.Status = 0
		rollback.DateCancel = time.Now()

		affected, err := engine.Id(orderId).Cols("rollback_status", "date_cancel").Update(&rollback)
		if err != nil {
			return rollback, err
		}
		if affected > 0 {
			return rollback, nil
		}
		return rollback, errors.New("更新失败")
	}
	return rollback, errors.New("没有这条工单信息")
}

func (r *PlusRollbackModel) GetAvailableRollbackOrderInfo(orderId int) (rollback PlusapiRollback_workorder, err error) {
	engine := db.GetDB("deploy_online")
	has, err := engine.Where("id = ? and rollback_status = ?", orderId, 1).Get(&rollback)
	if err != nil {
		return rollback, err
	}
	if has {
		return rollback, nil
	}
	return rollback, errors.New("没有这条工单信息")
}

func (r *PlusRollbackModel) GetRollbackOrderInfo(orderId int) (rollback PlusapiRollback_workorder, err error) {
	engine := db.GetDB("deploy_online")

	has, err := engine.Where("id = ?", orderId).Get(&rollback)
	if err != nil {
		return rollback, err
	}
	if has {
		return rollback, nil
	}
	return rollback, errors.New("没有这条工单信息")
}

func (r *PlusRollbackModel) UpdateRollbackOrderInfo(orderId int, rollbackTime time.Time) (rollback PlusapiRollback_workorder, err error) {
	engine := db.GetDB("deploy_online")
	has, err := engine.Where("id = ?", orderId).Get(&rollback)
	if err != nil {
		return rollback, err
	}
	if has {
		rollback.DateRollback = rollbackTime
		rollback.Status = 2
		affected, err := engine.Id(orderId).Cols("rollback_status", "date_rollback").Update(&rollback)
		if err != nil {
			return rollback, err
		}
		if affected > 0 {
			return rollback, nil
		}
	}
	return rollback, errors.New("没有这条工单信息")
}
