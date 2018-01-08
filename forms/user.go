package forms

import "time"

//SigninForm ...
type SigninForm struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

//SignupForm ...
type SignupForm struct {
	Username string `form:"username" json:"username" binding:"required,max=100"`
	Password string `form:"password" json:"password" binding:"required"`
}

type AdduserForm struct {
	Username string `form:"username" json:"username" binding:"required,max=100"`
	Password string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email"    json:"email" binding:"required"`
	RoleId   string `form:"role_id"  json:"role_id" binding:"required"`
}

type EditUserForm struct {
	Username string `form:"username" json:"username" binding:"required,max=100"`
	Password string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email"    json:"email" binding:"required"`
	RoleId   string `form:"role_id"  json:"role_id" binding:"required"`
	UserId   string `form:"uid"      json:"uid" binding:"required"`
}

type UserBehavior struct {
	Id         int       `xorm:"id int(11) pk not null"`
	Uid        int       `xorm:"uid int(11) not null"`
	Username   string    `xorm:"username varchar(50)"`
	Operate    string    `xorm:"operate  varchar(100)"`
	Datajson   string    `xorm:"datajson  varchar(500)"`
	DateCreate time.Time `xorm:"created"`
}
