package models

import (
	//"fmt"
	"encoding/base64"
	"errors"
	"mint-platform/platform/db"
	"mint-platform/platform/forms"
	"strconv"
	"time"
)

type Plusapi_workorder struct {
	Id                     int       `xorm:"id int(11) pk not null"`
	UserName               string    `xorm:"user_name varchar(50) not null"`
	TaskName               string    `xorm:"task_name varchar(100) not null"`
	EnvironmentType        int       `xorm:"environment_type smallint(6) not null"`
	GroupName              string    `xorm:"group_name varchar(20) not null"`
	Version                string    `xorm:"git_version varchar(100) not null"`
	Status                 int       `xorm:"deploy_status smallint(6) not null"`
	DateCreate             time.Time `xorm:"created"`
	DateDeployed           time.Time `xorm:"date_deployed not null"`
	DateCancel             time.Time `xorm:"date_cancel not null"`
	EmailList              string    `xorm:"email_list varchar(255) not null"`
	FunctionalIntroduction string    `xorm:"functional_introduction text not null"`
	IsCheck                int       `xorm:"is_check tinyint(2) not null"`
	Auditor                string    `xorm:"auditor varchar(50) not null"`
}

type PlusProjectModel struct{}

func (p *PlusProjectModel) TableName() string {
	return "plusapi_workorder"
}

func (p *PlusProjectModel) GetProjectList(page int, pageSize int) (projects []Plusapi_workorder, err error) {
	engine := db.GetDB("deploy_online")
	projects = make([]Plusapi_workorder, 0)
	startNum := (page - 1) * pageSize
	err = engine.In("deploy_status", 1, 2).Desc("id").Limit(pageSize, startNum).Find(&projects)
	if err != nil {
		return projects, err
	}
	return projects, nil
}

func (p *PlusProjectModel) GetProjectTotal() (num int64, err error) {
	engine := db.GetDB("deploy_online")
	project := new(Plusapi_workorder)
	num, err = engine.In("deploy_status", 1, 2).Count(project)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (p *PlusProjectModel) AddWorker(form forms.PlusProjectForm) (project Plusapi_workorder, err error) {
	engine := db.GetDB("deploy_online")
	/*
		    has, err := engine.Where("deploy_status = ? ", 1).Limit(1).Get(&project)
			if err != nil {
				return user, err
			}
			if has {
				return user, errors.New("还有未发布的工单！")
			}
	*/
	project.UserName = form.UserName
	project.TaskName = form.TaskName
	project.EnvironmentType, _ = strconv.Atoi(form.EnvironmentType)
	project.GroupName = form.Group
	project.Version = form.Version
	project.Status = 1

	//正式环境才发邮件
	if project.EnvironmentType == 1 {
		project.EmailList = form.EmailList
		FunctionalIntroductionDecode, _ := base64.StdEncoding.DecodeString(form.FunctionalIntroduction)
		project.FunctionalIntroduction = string(FunctionalIntroductionDecode)
	}
	affected, err := engine.Insert(&project)
	if err != nil {
		return project, err
	}
	if affected > 0 {
		return project, nil
	}
	return project, errors.New("添加发布工单失败！")
}

func (p *PlusProjectModel) CancelWorker(orderId int) (project Plusapi_workorder, err error) {
	engine := db.GetDB("deploy_online")

	has, err := engine.Where("id = ?", orderId).Get(&project)
	if err != nil {
		return project, err
	}
	if has {
		project.Status = 0
		project.DateCancel = time.Now()

		affected, err := engine.Id(orderId).Cols("deploy_status", "date_cancel").Update(&project)
		if err != nil {
			return project, err
		}
		if affected > 0 {
			return project, nil
		}
		return project, errors.New("更新失败")
	}
	return project, errors.New("没有这条工单信息")
}

func (p *PlusProjectModel) AccessWorker(orderId int, userName string) (project Plusapi_workorder, err error) {
	engine := db.GetDB("deploy_online")

	has, err := engine.Where("id = ?", orderId).Get(&project)
	if err != nil {
		return project, err
	}
	if has {
		project.IsCheck = 1
		project.Auditor = userName

		affected, err := engine.Id(orderId).Cols("is_check", "auditor").Update(&project)
		if err != nil {
			return project, err
		}
		if affected > 0 {
			return project, nil
		}
		return project, errors.New("更新失败")
	}
	return project, errors.New("没有这条工单信息")
}

func (p *PlusProjectModel) GetDeployOrderInfo(orderId int) (project Plusapi_workorder, err error) {
	engine := db.GetDB("deploy_online")

	has, err := engine.Where("id = ?", orderId).Get(&project)
	if err != nil {
		return project, err
	}
	if has {
		return project, nil
	}
	return project, errors.New("没有这条工单信息")
}

func (p *PlusProjectModel) UpdateDeployOrderInfo(orderId int, deployTime time.Time) (project Plusapi_workorder, err error) {
	engine := db.GetDB("deploy_online")
	has, err := engine.Where("id = ?", orderId).Get(&project)
	if err != nil {
		return project, err
	}
	if has {
		project.DateDeployed = deployTime
		project.Status = 2
		affected, err := engine.Id(orderId).Cols("deploy_status", "date_deployed").Update(&project)
		if err != nil {
			return project, err
		}
		if affected > 0 {
			return project, nil
		}
	}
	return project, errors.New("没有这条工单信息")
}

func (p *PlusProjectModel) GetAvailableDeployOrderInfo(orderId int) (plusproject Plusapi_workorder, err error) {
	engine := db.GetDB("deploy_online")

	has, err := engine.Where("id = ? and deploy_status=?", orderId, 1).Get(&plusproject)
	if err != nil {
		return plusproject, err
	}
	if has {
		return plusproject, nil
	}
	return plusproject, errors.New("没有这条工单信息")
}

func (p *PlusProjectModel) GetHistoryDeployVersion(groupName string) (projects []Plusapi_workorder, err error) {
	engine := db.GetDB("deploy_online")

	err = engine.Where("deploy_status =? and group_name=?", 2, groupName).Desc("id").Limit(5, 1).Find(&projects)
	if err != nil {
		return projects, err
	}
	return projects, nil
}
