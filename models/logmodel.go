package models

import (
	"mint-platform/platform/db"
	"strconv"
	"time"
	//"fmt"
)

type Api_log struct {
	Log_id          int     `xorm:"pk;auto"`
	Dt              string  `xorm:"dt varchar(30)"`
	Api             string  `xorm:"api varchar(100)"`
	Ios_request     int64   `xorm:"ios_request bigint(20)"`
	Android_request int64   `xorm:"android_request bigint(20)"`
	Other_request   int64   `xorm:"other_request bigint(20)"`
	Total_request   int64   `xorm:"total_request bigint(20)"`
	Min_time        float64 `xorm:"min_time float"`
	Max_time        float64 `xorm:"max_time float"`
	Avg_time        float64 `xorm:"avg_time float"`
}

type App_login_register_log struct {
	Id           int       `xorm:"id pk;auto"`
	User_id      int       `xorm:"user_id int(10)"`
	Cell_phone   string    `xorm:"cell_phone varchar(11)"`
	Device       string    `xorm:"device varchar(20)"`
	Session      string    `xorm:"session varchar(50)"`
	Auth_session string    `xorm:"auth_session varchar(300)"`
	Version      string    `xorm:"version varchar(40)"`
	Channel_code string    `xorm:"channel_code varchar(300)"`
	Idfa         string    `xorm:"idfa varchar(300)"`
	Dvc_id       string    `xorm:"dvc_id varchar(300)"`
	Ip           string    `xorm:"ip varchar(50)"`
	Ua           string    `xorm:"ua varchar(500)"`
	Verify_code  string    `xorm:"verify_code varchar(20)"`
	Log_type     int       `xorm:"log_type tinyint(4)"`
	Create_time  time.Time `xorm:"create_time"`
}

type LogModel struct{}

func (l *LogModel) TableName() string {
	return "api_log"
}

func (l *LogModel) ReadList(dt string, api string, page int, pageSize int) (apilogs []Api_log, err error) {
	//db.Init("deploy_apibi")
	engine := db.GetDB("deploy_apibi")
	sql := "select dt,api,ios_request,android_request,other_request,total_request,min_time,max_time,avg_time from api_log where 1=1 "
	where := ""
	if len(dt) > 0 {
		where += " and dt='" + dt + "'"
	}
	if len(api) > 0 {
		where += " and api like '%" + api + "%'"
	}
	where += " order by dt desc,total_request desc"
	offset := strconv.Itoa((page - 1) * pageSize)
	limit := " limit " + offset + ", " + strconv.Itoa(pageSize)

	sql += where + limit
	err = engine.Sql(sql).Find(&apilogs)
	if err != nil {
		return apilogs, err
	}
	return apilogs, nil
}

/**
获取所有的log数量
*/
func (l *LogModel) GetLogTotal(dt string, api string) (num int64, err error) {
	engine := db.GetDB("deploy_apibi")
	sql := "select count(*) as total from api_log where 1=1 "
	where := " "
	if len(dt) > 0 {
		where += " and dt='" + dt + "'"
	}
	if len(api) > 0 {
		where += " and api like '%" + api + "%'"
	}
	sql += where
	var apilog Api_log
	num, err = engine.Sql(sql).Count(&apilog)
	if err != nil {
		return 0, err
	}

	return num, nil
}

/*
获取所有的登录日志信息
*/
func (l *LogModel) GetLoginList(userId string, cell_phone string, ip string, page int, pageSize int) (apiloginlogs []App_login_register_log, err error) {
	engine := db.GetDB("deploy_log")
	sql := "select id,user_id,cell_phone,device,session,version,channel_code,idfa,dvc_id,ip,ua, log_type, create_time  from app_login_register_log where 1=1 "
	where := ""
	if len(userId) > 0 {
		where += " and user_id='" + userId + "'"
	}
	if len(cell_phone) > 0 {
		where += " and cell_phone = '" + cell_phone + "'"
	}
	if len(ip) > 0 {
		where += " and ip = '" + ip + "'"
	}
	where += " order by id desc"
	offset := strconv.Itoa((page - 1) * pageSize)
	limit := " limit " + offset + ", " + strconv.Itoa(pageSize)

	sql += where + limit
	err = engine.Sql(sql).Find(&apiloginlogs)
	if err != nil {
		return apiloginlogs, err
	}
	return apiloginlogs, nil
}

/**
获取所有的登录log数量
*/
func (l *LogModel) GetLoginLogTotal(userId string, cell_phone string, ip string) (num int64, err error) {
	engine := db.GetDB("deploy_log")
	sql := "select count(*) as total from app_login_register_log where 1=1 "
	where := " "
	if len(userId) > 0 {
		where += " and user_id='" + userId + "'"
	}
	if len(cell_phone) > 0 {
		where += " and cell_phone = '" + cell_phone + "'"
	}
	if len(ip) > 0 {
		where += " and ip = '" + ip + "'"
	}
	sql += where
	var apiloginlog App_login_register_log
	num, err = engine.Sql(sql).Count(&apiloginlog)
	if err != nil {
		return 0, err
	}

	return num, nil
}
