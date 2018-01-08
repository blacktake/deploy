package models

import (
	"errors"
	//"fmt"
	"mint-platform/platform/db"
	"mint-platform/platform/forms"
	"strconv"
	"strings"
	"time"
)

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
	UpdateDate time.Time `xorm:"update_date not null"`
}

type MenuModel struct{}

func (m *MenuModel) TableName() string {
	return "menu"
}

/**
获取有效的父菜单
*/
func (m *MenuModel) GetParentMenuList() (menus []Menu, err error) {
	engine := db.GetDB("deploy_online")
	//engine.Sync2(new(Menu))url

	err = engine.Where("parentid = ? and menu_status = ?", 0, 1).Asc("listorder").Find(&menus)
	if err != nil {
		return menus, err
	}
	return menus, nil
}

/**
增加菜单
*/
func (m *MenuModel) AddMenu(form forms.MenuForm) (menu Menu, err error) {
	engine := db.GetDB("deploy_online")
	menu.Name = form.Name
	menu.Parentid, _ = strconv.Atoi(form.Parentid)
	if len(form.Url) > 0 {
		menu.Islink = 1
		menu.Url = form.Url
	} else {
		menu.Islink = 0
		menu.Url = ""
	}
	menu.ModuleName = form.ModuleName
	menu.CName = form.CName
	menu.MenuStatus = 1
	menu.Icon = form.Icon
	affected, err := engine.Insert(&menu)
	if err != nil {
		return menu, err
	}
	if affected > 0 {
		return menu, nil
	}
	return menu, errors.New("添加菜单失败！")
}

/**
获取有效的菜单列表
*/
func (m *MenuModel) GetMenuList() (menus []Menu, err error) {
	engine := db.GetDB("deploy_online")
	//engine.Sync2(new(Menu))

	err = engine.Where("menu_status = ?", 1).Asc("listorder").Find(&menus)
	if err != nil {
		return menus, err
	}
	return menus, nil
}

/**
更新菜单信息
*/
func (m *MenuModel) UpdateMenuOrder(id int, order int) (menu Menu, err error) {
	engine := db.GetDB("deploy_online")

	has, err := engine.Where("id = ?", id).Get(&menu)
	if err != nil {
		return menu, err
	}
	if has {
		menu.Listorder = order
		menu.UpdateDate = time.Now()

		affected, err := engine.Id(id).Cols("listorder", "update_date").Update(&menu)
		if err != nil {
			return menu, err
		}
		if affected > 0 {
			return menu, nil
		}
		return menu, errors.New("更新失败")
	}
	return menu, errors.New("没有这条菜单信息")
}

/**
改变菜单顺寻
*/
func (m *MenuModel) ChangeMenuOrder(menuidlist []forms.MenuId) (err error) {
	listsort := 0
	for _, value := range menuidlist {
		if strings.Contains(value.Id, "children") {
			return errors.New("子菜单不能作为父菜单！")
		}
		id := strings.Replace(value.Id, "parent_", "", -1)
		idint, _ := strconv.Atoi(id)
		listsort += 1
		_, err := m.UpdateMenuOrder(idint, listsort)
		if err != nil {
			return err
		}

		for _, cv := range value.Childrens {
			idstr := strings.Replace(cv.Id, "children_", "", -1)
			childrenid, _ := strconv.Atoi(idstr)
			listsort += 1
			_, err := m.UpdateMenuOrder(childrenid, listsort)
			if err != nil {
				return err
			}
		}

	}
	return nil

}

/**
获取所有的menu列表
*/
func (m *MenuModel) GetAllMenuList(menuList []Menu) []forms.ParentMenu {
	var tmpMenuList []forms.ParentMenu
	var newMenuList []forms.ParentMenu
	var newParentMenu forms.ParentMenu
	if len(menuList) > 0 {
		for _, value := range menuList {
			if value.Parentid == 0 {
				newParentMenu.Id = value.Id
				newParentMenu.Name = value.Name
				newParentMenu.ModuleName = value.ModuleName
				newParentMenu.Listorder = value.Listorder
				newParentMenu.Icon = value.Icon
				tmpMenuList = append(tmpMenuList, newParentMenu)
			}
		}
		for _, parentMenu := range tmpMenuList {
			for _, value := range menuList {
				if parentMenu.Id == value.Parentid {
					newChildrenMenu := new(forms.ChildrenMenu)
					newChildrenMenu.Id = value.Id
					newChildrenMenu.Name = value.Name
					newChildrenMenu.Url = value.Url
					newChildrenMenu.ModuleName = value.ModuleName
					newChildrenMenu.CName = value.CName
					newChildrenMenu.Listorder = value.Listorder
					newChildrenMenu.ParentId = value.Parentid
					parentMenu.ChildrenList = append(parentMenu.ChildrenList, newChildrenMenu)
				}
			}
			newMenuList = append(newMenuList, parentMenu)
		}
	}
	return newMenuList
}

/**
删除子菜单
*/
func (m *MenuModel) DeleteMenu(menuId int) (err error) {
	engine := db.GetDB("deploy_online")

	menu := new(Menu)
	has, err := engine.Where("id = ?", menuId).Get(menu)
	if err != nil {
		return err
	}
	if has {
		menu.MenuStatus = 0
		menu.UpdateDate = time.Now()

		affected, err := engine.Id(menuId).Cols("menu_status", "update_date").Update(menu)
		if err != nil {
			return err
		}
		if affected > 0 {
			return nil
		}
		return errors.New("更新失败")
	}
	return errors.New("没有这条菜单信息")
}

/**
根据菜单id获取有效的菜单信息
*/
func (m *MenuModel) GetMenuInfo(id int) (menu Menu, err error) {
	engine := db.GetDB("deploy_online")
	//engine.Sync2(new(Menu))url

	has, err := engine.Where("id = ? and menu_status=?", id, 1).Get(&menu)
	if err != nil {
		return menu, err
	}
	if has {
		return menu, nil
	}
	return menu, errors.New("没有这条菜单信息")
}

/**
更新单个菜单信息
*/
func (m *MenuModel) EditMenu(form forms.AllMenuForm) (menu Menu, err error) {
	engine := db.GetDB("deploy_online")

	has, err := engine.Where("id = ?", form.Id).Get(&menu)
	if err != nil {
		return menu, err
	}
	if has {
		menu.Name = form.Name
		menu.Url = form.Url
		menu.ModuleName = form.ModuleName
		menu.Parentid, _ = strconv.Atoi(form.Parentid)
		menu.CName = form.CName
		menu.Icon = form.Icon
		menu.UpdateDate = time.Now()
		affected, err := engine.Id(form.Id).Update(&menu)
		if err != nil {
			return menu, err
		}
		if affected > 0 {
			return menu, nil
		}
		return menu, errors.New("更新失败")
	}
	return menu, errors.New("没有这条菜单信息")
}

/**
根据子菜单ids获取所有菜单列表
*/
func (m *MenuModel) getMenuList(menuIds []int) (menus []Menu, err error) {
	engine := db.GetDB("deploy_online")
	err = engine.In("id", menuIds).Asc("listorder").Find(&menus)
	if err != nil {
		return menus, err
	}
	return menus, nil
}

/**
根据子菜单ids获取所有子菜单列表
*/
func (m *MenuModel) regroup(parentMenuList []Menu, childrenMenuList []Menu) []forms.ParentMenu {
	var tmpMenuList []forms.ParentMenu
	var newMenuList []forms.ParentMenu
	var newParentMenu forms.ParentMenu
	if len(parentMenuList) > 0 && len(childrenMenuList) > 0 {
		for _, value := range parentMenuList {
			newParentMenu.Id = value.Id
			newParentMenu.Name = value.Name
			newParentMenu.ModuleName = value.ModuleName
			newParentMenu.Listorder = value.Listorder
			newParentMenu.Icon = value.Icon
			newParentMenu.ParentId = value.Parentid
			tmpMenuList = append(tmpMenuList, newParentMenu)
		}
		for _, parentMenu := range tmpMenuList {
			for _, value := range childrenMenuList {
				if parentMenu.Id == value.Parentid {
					newChildrenMenu := new(forms.ChildrenMenu)
					newChildrenMenu.Id = value.Id
					newChildrenMenu.Name = value.Name
					newChildrenMenu.Url = value.Url
					newChildrenMenu.ModuleName = value.ModuleName
					newChildrenMenu.CName = value.CName
					newChildrenMenu.Listorder = value.Listorder
					newChildrenMenu.ParentId = value.Parentid
					parentMenu.ChildrenList = append(parentMenu.ChildrenList, newChildrenMenu)
				}
			}
			newMenuList = append(newMenuList, parentMenu)
		}
	}
	return newMenuList
}

/*
func (m *MenuModel) GetMenu() (menu []forms.ParentMenu, err error) {
	menuList, err := m.GetMenuList()
	if err != nil {
		return menu, err
	}
	menu = m.GetAllMenuList(menuList)
	return menu, nil
}
*/

/**
封装统一获取菜单的方法
*/
func (m *MenuModel) GetMenu(roleId int) (menu []forms.ParentMenu, err error) {
	permissionModel := new(PermissionModel)
	rolemenulist, err := permissionModel.GetRoleIdAndMenuIdReleation(roleId)
	if err != nil {
		return menu, err
	}
	var childrenMenuIds []int
	for _, v := range rolemenulist {
		childrenMenuIds = append(childrenMenuIds, v.MenuId)
	}
	childrenMenuList, err := m.getMenuList(childrenMenuIds)
	if err != nil {
		return menu, err
	}
	var parentMenuIds []int
	for _, menuInfo := range childrenMenuList {
		parentMenuIds = append(parentMenuIds, menuInfo.Parentid)
	}
	parentMenuList, err := m.getMenuList(parentMenuIds)
	if err != nil {
		return menu, err
	}
	menu = m.regroup(parentMenuList, childrenMenuList)
	return menu, nil
}
