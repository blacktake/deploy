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

type Admin_role struct {
	RoleId      int       `xorm:"roleid tinyint(3) pk not null"`
	RoleName    string    `xorm:"rolename varchar(50) not null"`
	Description string    `xorm:"description text"`
	RoleStatus  int       `xorm:"role_status tinyint(1)"`
	CreateTime  time.Time `xorm:"created"`
	UpdateTime  time.Time `xorm:"update_time"`
}

type Admin_role_menu_relation struct {
	RoleId int `xorm:"roleid tinyint(3) pk not null"`
	MenuId int `xorm:"menuid smallint(6) pk not null"`
}

type PermissionModel struct{}

/**
获取有效的角色列表
*/
func (m *PermissionModel) GetRoleList() (roles []Admin_role, err error) {
	engine := db.GetDB("deploy_online")
	//engine.Sync2(new(Menu))url

	err = engine.Where("role_status = ?", 1).Find(&roles)
	if err != nil {
		return roles, err
	}
	return roles, nil
}

/**
根据角色id获取角色信息
*/
func (m *PermissionModel) GetRoleInfo(roleId int) (role Admin_role, err error) {
	engine := db.GetDB("deploy_online")
	//engine.Sync2(new(Menu))url
	_, err = engine.Where("roleid = ? and role_status = ?", roleId, 1).Get(&role)
	if err != nil {
		return role, err
	}
	return role, nil
}

/**
添加角色
*/
func (m *PermissionModel) AddRole(form forms.AddRoleForm) (role Admin_role, err error) {
	engine := db.GetDB("deploy_online")
	has, err := engine.Where("rolename = ? ", form.Rolename).Get(&role)
	if err != nil {
		return role, err
	}
	if has {
		return role, errors.New("角色已存在！")
	}
	role.RoleName = form.Rolename
	role.Description = form.Description
	role.RoleStatus = 1
	affected, err := engine.Insert(&role)
	if err != nil {
		return role, err
	}
	if affected > 0 {
		return role, nil
	}
	return role, errors.New("添加角色失败！")
}

/**
根据角色id获取menuIds
*/
func (m *PermissionModel) GetRoleIdAndMenuIdReleation(roleId int) (rolemenus []Admin_role_menu_relation, err error) {
	engine := db.GetDB("deploy_online")
	//engine.Sync2(new(Menu))url
	err = engine.Where("roleid = ?", roleId).Find(&rolemenus)
	if err != nil {
		return rolemenus, err
	}
	return rolemenus, nil
}

/**
更新浏览菜单权限
*/
func (m *PermissionModel) UpdatePermission(roleId int, ids string) (result int, err error) {
	engine := db.GetDB("deploy_online")
	roleMenuIds, err := m.GetRoleIdAndMenuIdReleation(roleId)
	if err != nil {
		return 0, err
	}
	if len(roleMenuIds) > 0 {
		sql := "DELETE FROM `admin_role_menu_relation` WHERE `roleid` = ? "
		_, err := engine.Exec(sql, roleId)
		if err != nil {
			return 0, err
		}
	}
	idsArr := strings.Split(ids, ",")
	rolemenus := make([]Admin_role_menu_relation, 0)
	for _, id := range idsArr {
		menuId, _ := strconv.Atoi(id)
		var rolemenu Admin_role_menu_relation
		rolemenu.RoleId = roleId
		rolemenu.MenuId = menuId
		rolemenus = append(rolemenus, rolemenu)
	}
	if len(rolemenus) == 0 {
		return 2, nil
	}
	_, err = engine.Insert(&rolemenus)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

/**
根据角色id和menuid获取角色是否拥有该菜单权限
*/
func (m *PermissionModel) GetRoleMenuRelation(roleId int, menuId int) (rolemenu Admin_role_menu_relation, err error) {
	engine := db.GetDB("deploy_online")
	_, err = engine.Where("roleid = ? and menuid = ?", roleId, menuId).Get(&rolemenu)
	if err != nil {
		return rolemenu, err
	}
	return rolemenu, nil
}

/**
整合权限menu和整个menu，在menu中标记哪些是此角色有权限的
*/
func (m *PermissionModel) FormatMenuList(roleId int) (menus []forms.PermissionParentMenu, err error) {
	roleMenuIds, err := m.GetRoleIdAndMenuIdReleation(roleId)
	if err != nil {
		return menus, err
	}
	menuModel := new(MenuModel)
	menulist, err := menuModel.GetMenuList()
	if err != nil {
		return menus, err
	}
	var tmpMenuList []forms.PermissionParentMenu
	var newMenuList []forms.PermissionParentMenu
	var tmpChildrenMenuList []forms.PermissionChildrenMenu
	var newChildrenMenuList []forms.PermissionChildrenMenu
	if len(menulist) > 0 {
		for _, value := range menulist {
			if value.Parentid == 0 {
				var newParentMenu forms.PermissionParentMenu
				newParentMenu.Id = value.Id
				newParentMenu.Name = value.Name
				newParentMenu.ModuleName = value.ModuleName
				newParentMenu.Listorder = value.Listorder
				newParentMenu.Icon = value.Icon
				newParentMenu.ParentId = value.Parentid
				tmpMenuList = append(tmpMenuList, newParentMenu)
			} else {
				var tmpChildrenMenu forms.PermissionChildrenMenu
				tmpChildrenMenu.Id = value.Id
				tmpChildrenMenu.Name = value.Name
				tmpChildrenMenu.ModuleName = value.ModuleName
				tmpChildrenMenu.Url = value.Url
				tmpChildrenMenu.CName = value.CName
				tmpChildrenMenu.Listorder = value.Listorder
				tmpChildrenMenu.ParentId = value.Parentid
				tmpChildrenMenuList = append(tmpChildrenMenuList, tmpChildrenMenu)
			}
		}
		if len(roleMenuIds) > 0 {
			for _, value := range tmpChildrenMenuList {
				var childrenMenuStruct forms.PermissionChildrenMenu
				childrenMenuStruct.Id = value.Id
				childrenMenuStruct.Name = value.Name
				childrenMenuStruct.Url = value.Url
				childrenMenuStruct.ModuleName = value.ModuleName
				childrenMenuStruct.CName = value.CName
				childrenMenuStruct.Listorder = value.Listorder
				childrenMenuStruct.ParentId = value.ParentId
				for _, childrenMenu := range roleMenuIds {
					if childrenMenu.MenuId == value.Id {
						childrenMenuStruct.RoleId = childrenMenu.RoleId
					}
				}
				newChildrenMenuList = append(newChildrenMenuList, childrenMenuStruct)
			}
		} else {
			newChildrenMenuList = tmpChildrenMenuList
		}

		//fmt.Println(newChildrenMenuList, "**********", tmpChildrenMenuList, "***************", roleMenuIds)
		for _, parentMenu := range tmpMenuList {
			for _, value := range newChildrenMenuList {
				if parentMenu.Id == value.ParentId {
					newChildrenMenu := new(forms.PermissionChildrenMenu)
					newChildrenMenu.Id = value.Id
					newChildrenMenu.Name = value.Name
					newChildrenMenu.Url = value.Url
					newChildrenMenu.ModuleName = value.ModuleName
					newChildrenMenu.CName = value.CName
					newChildrenMenu.Listorder = value.Listorder
					newChildrenMenu.ParentId = value.ParentId
					if value.RoleId > 0 {
						newChildrenMenu.RoleId = value.RoleId
						parentMenu.RoleId = value.RoleId
					}
					parentMenu.ChildrenList = append(parentMenu.ChildrenList, newChildrenMenu)
				}
			}
			newMenuList = append(newMenuList, parentMenu)
		}

	}
	return newMenuList, nil
}
