package models

import (
	"errors"
	"mint-platform/platform/db"
	"mint-platform/platform/forms"
	"strconv"
	"time"
	//"fmt"
)

type Deployapi_user struct {
	Id         int       `xorm:"id int(11) pk not null"`
	UserName   string    `xorm:"username varchar(50) not null"`
	PassWord   string    `xorm:"password varchar(50) not null"`
	Email      string    `xorm:"email varchar(32) not null"`
	UserType   int       `xorm:"user_type int(11) not null"`
	UserStatus int       `xorm:"user_status smallint(6) not null"`
	RoleId     int       `xorm:"role_id smallint(6)"`
	RoleName   string    `xorm:"role_name varchar(40)"`
	LastLogin  time.Time `xorm:"last_login"`
	DateJoined time.Time `xorm:"created"`
}

type UserModel struct{}

func (u *UserModel) TableName() string {
	return "deployapi_user"
}

func (u *UserModel) Signin(form forms.SigninForm) (user Deployapi_user, err error) {
	engine := db.GetDB("deploy_online")
	engine.Sync2(new(Deployapi_user))

	has, err := engine.Where("username = ? and password = ? and user_status = ? ", form.Username, form.Password, 1).Get(&user)
	if err != nil {
		return user, err
	}
	if has {
		user.LastLogin = time.Now()
		affected, err := engine.Id(user.Id).Cols("last_login").Update(&user)
		if affected > 0 {
			return user, nil
		}
		return user, err
	}
	return user, errors.New("user信息为空")
}

func (u *UserModel) GetUserList(page int, pageSize int) (user []Deployapi_user, err error) {
	engine := db.GetDB("deploy_online")
	users := make([]Deployapi_user, 0)
	err = engine.Where("user_status != ?", 2).Desc("id").Limit(pageSize, page).Find(&users)
	if err != nil {
		return users, err
	}
	return users, nil
}

func (u *UserModel) ChangeUserStatus(uid int) (user Deployapi_user, err error) {
	engine := db.GetDB("deploy_online")
	has, err := engine.Where("id = ?", uid).Get(&user)
	if err != nil {
		return user, err
	}
	if has {
		if user.UserStatus == 1 {
			user.UserStatus = 0
		} else {
			user.UserStatus = 1
		}
		affected, err := engine.Id(uid).Cols("user_status").Update(&user)
		if err != nil {
			return user, err
		}
		if affected > 0 {
			return user, nil
		}
		return user, errors.New("更新失败")
	}
	return user, errors.New("user信息为空")
}

func (u *UserModel) DeleteUser(uid int) (user Deployapi_user, err error) {
	engine := db.GetDB("deploy_online")
	has, err := engine.Where("id = ?", uid).Get(&user)
	if err != nil {
		return user, err
	}
	if has {
		user.UserStatus = 2
		affected, err := engine.Id(uid).Cols("user_status").Update(&user)
		if err != nil {
			return user, err
		}
		if affected > 0 {
			return user, nil
		}
		return user, errors.New("更新失败")
	}
	return user, errors.New("user信息为空")
}

func (u *UserModel) AddUser(form forms.AdduserForm) (user Deployapi_user, err error) {
	engine := db.GetDB("deploy_online")
	has, err := engine.Where("username = ? or email = ? ", form.Username, form.Email).Get(&user)
	if err != nil {
		return user, err
	}
	if has {
		return user, errors.New("用户已存在！")
	}
	permissionModel := new(PermissionModel)
	user.UserName = form.Username
	user.PassWord = form.Password
	user.Email = form.Email
	user.RoleId, _ = strconv.Atoi(form.RoleId)
	roleInfo, err := permissionModel.GetRoleInfo(user.RoleId)
	if err != nil {
		return user, nil
	}
	user.RoleName = roleInfo.RoleName
	user.UserType = 1
	affected, err := engine.Insert(&user)
	if err != nil {
		return user, err
	}
	if affected > 0 {
		return user, nil
	}
	return user, errors.New("添加用户失败！")
}

/**
修改用户信息
*/
func (u *UserModel) EditUser(form forms.EditUserForm) (user Deployapi_user, err error) {
	engine := db.GetDB("deploy_online")
	uid, _ := strconv.Atoi(form.UserId)
	has, err := engine.Where("id = ?", uid).Get(&user)
	if err != nil {
		return user, err
	}
	if has {
		permissionModel := new(PermissionModel)
		user.UserName = form.Username
		user.PassWord = form.Password
		user.Email = form.Email
		user.RoleId, _ = strconv.Atoi(form.RoleId)
		roleInfo, err := permissionModel.GetRoleInfo(user.RoleId)
		if err != nil {
			return user, nil
		}
		user.RoleName = roleInfo.RoleName
		affected, err := engine.Id(uid).Update(&user)
		if err != nil {
			return user, err
		}
		if affected > 0 {
			return user, nil
		}
		return user, errors.New("修改用户失败！")
	}
	return user, errors.New("没有这条用户信息")
}

func (u *UserModel) GetUsersTotal() (num int64, err error) {
	engine := db.GetDB("deploy_online")
	user := new(Deployapi_user)
	num, err = engine.Count(user)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (u *UserModel) GetUserInfo(uid int) (user Deployapi_user, err error) {
	engine := db.GetDB("deploy_online")
	has, err := engine.Where("id = ?", uid).Get(&user)
	if err != nil {
		return user, err
	}
	if has {
		return user, nil
	}
	return user, errors.New("user信息为空")
}
