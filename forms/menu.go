package forms

import "time"

type MenuForm struct {
	Name       string `form:"name" json:"name" binding:"required"`
	Url        string `form:"url" json:"url"`
	Parentid   string `form:"parentid" json:"parentid" binding:"required"`
	ModuleName string `form:"module_name" json:"module_name" binding:"required"`
	CName      string `form:"c_name" json:"c_name"`
	Icon       string `form:"icon" json:"icon"`
}

type MenuId struct {
	Id        string       `json:"id"`
	Childrens []ChildrenId `json:"children"`
}

type ChildrenId struct {
	Id string `json:"id"`
}

type ParentMenu struct {
	Id           int             `json:"id"`
	Name         string          `json:"name"`
	ChildrenList []*ChildrenMenu `json:"childrenlist"`
	ModuleName   string          `json:"module_name"`
	CName        string          `json:"c_name"`
	Icon         string          `json:"icon"`
	ParentId     int             `json:"parentid"`
	Listorder    int             `json:"listorder"`
}

type ChildrenMenu struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Url        string `json:"url"`
	ModuleName string `json:"module_name"`
	CName      string `json:"c_name"`
	ParentId   int    `json:"parentid"`
	Listorder  int    `json:"listorder"`
}

type AllMenuForm struct {
	Id         string `form:"id" json:"id" binding:"required"`
	Name       string `form:"name" json:"name" binding:"required"`
	Url        string `form:"url" json:"url"`
	Parentid   string `form:"parentid" json:"parentid"`
	ModuleName string `form:"module_name" json:"module_name" binding:"required"`
	CName      string `form:"c_name" json:"c_name"`
	Icon       string `form:"icon" json:"icon"`
}

type PermissionParentMenu struct {
	Id           int                       `json:"id"`
	Name         string                    `json:"name"`
	ChildrenList []*PermissionChildrenMenu `json:"childrenlist"`
	ModuleName   string                    `json:"module_name"`
	CName        string                    `json:"c_name"`
	Icon         string                    `json:"icon"`
	ParentId     int                       `json:"parentid"`
	Listorder    int                       `json:"listorder"`
	RoleId       int                       `json:"roleid"`
}

type PermissionChildrenMenu struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Url        string `json:"url"`
	ModuleName string `json:"module_name"`
	CName      string `json:"c_name"`
	Listorder  int    `json:"listorder"`
	ParentId   int    `json:"parentid"`
	RoleId     int    `json:"roleid"`
}

type Menu struct {
	Id         int       `xorm:"id int(11) pk not null"`
	Name       string    `xorm:"Name varchar(40) not null"`
	Parentid   int       `xorm:"parentid smallint(6)"`
	Islink     int       `xorm:"islink tinyint(1)"`
	Url        string    `xorm:"url varchar(40)"`
	ModuleName string    `xorm:"module_name varchar(40)" not null`
	CName      string    `xorm:"c_name varchar(40)"`
	Icon       string    `xorm:"icon varchar(40)"`
	MenuStatus int       `xorm:"menu_status tinyint(3) not null"`
	Listorder  int       `xorm:"listorder tinyint(4)"`
	CreateDate time.Time `xorm:"created"`
	UpdateDate time.Time `xorm:"update_date"`
}

type Admin_role_menu_relation struct {
	RoleId int `xorm:"roleid tinyint(3) pk not null"`
	MenuId int `xorm:"menuid smallint(6) pk not null"`
}
