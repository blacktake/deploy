package controllers

import (
	//"fmt"
	//"io"
	//"strings"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"mint-platform/platform/forms"
	"mint-platform/platform/models"
	"mint-platform/platform/utils"
	"strconv"
	//"net/http"
)

type PermissionController struct{}

var permissionModel = new(models.PermissionModel)

func (ctrl PermissionController) RoleList(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")
	rolelist, err := permissionModel.GetRoleList()
	if err != nil {
		utils.WriteLog("log_permission", "get role list err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	menu, _ := ctrl.getMenu(roleId.(int))
	c.HTML(200, "permission/rolelist.html", gin.H{
		"username":   userName,
		"moduleName": "permission",
		"ctrName":    "rolelist",
		"rolelist":   rolelist,
		"menu":       menu,
		"ctrNameZn":  "角色列表",
	})
}

func (ctrl PermissionController) PermissionList(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")
	menu, _ := ctrl.getMenu(roleId.(int))
	c.HTML(200, "index/index.html", gin.H{
		"username":   userName,
		"moduleName": "permission",
		"ctrName":    "permissionlist",
		"menu":       menu,
		"ctrNameZn":  "权限列表",
	})
}

func (ctrl PermissionController) Userlist(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")
	pageparam := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageparam)
	if err != nil {
		utils.WriteLog("log_user", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	//fmt.Println("type:", reflect.TypeOf(page))
	pageSize := 20
	startPage := (page - 1) * 20
	userlist, err := UserModel.GetUserList(startPage, pageSize)
	if err != nil {
		utils.WriteLog("log_user", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	total, err := UserModel.GetUsersTotal()
	if err != nil {
		utils.WriteLog("log_user", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	paginator := ctrl.setPaginator(c, pageSize, total)
	menu, err := ctrl.getMenu(roleId.(int))
	if err != nil {
		utils.WriteLog("log_user", "get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	//fmt.Println("type:", reflect.TypeOf(userlist))
	c.HTML(200, "permission/userlist.html", gin.H{
		"username":   userName,
		"roleId":     roleId,
		"moduleName": "permission",
		"ctrName":    "user",
		"ctrNameZn":  "用户列表",
		"userList":   userlist,
		"paginator":  paginator,
		"menu":       menu,
	})
}

func (ctrl PermissionController) ChangeUser(c *gin.Context) {
	uidParam := c.PostForm("userid")
	userId, _ := strconv.Atoi(uidParam)
	user, err := UserModel.ChangeUserStatus(userId)
	if (err == nil) && (user.Id > 0) {
		c.JSON(200, gin.H{"code": "1", "desc": "success"})
	} else {
		utils.WriteLog("log_user", err)
		c.JSON(200, gin.H{"code": "0", "desc": "fail"})
	}
	return
}

func (ctrl PermissionController) DeleteUser(c *gin.Context) {
	uidParam := c.PostForm("userid")
	userId, _ := strconv.Atoi(uidParam)
	user, err := UserModel.DeleteUser(userId)
	if (err == nil) && (user.Id > 0) {
		c.JSON(200, gin.H{"code": "1", "desc": "success"})
	} else {
		utils.WriteLog("log_user", err)
		c.JSON(200, gin.H{"code": "0", "desc": "fail"})
	}
	return
}

func (ctrl PermissionController) ShowAddUser(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")
	rolelist, err := permissionModel.GetRoleList()
	if err != nil {
		utils.WriteLog("log_user", "permission get role list err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	menu, err := ctrl.getMenu(roleId.(int))
	if err != nil {
		utils.WriteLog("log_user", "get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	c.HTML(200, "permission/adduser.html", gin.H{
		"username":   userName,
		"moduleName": "permission",
		"ctrName":    "user",
		"ctrNameZn":  "添加用户",
		"rolelist":   rolelist,
		"menu":       menu,
	})
}

func (ctrl PermissionController) ShowEditUser(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")
	userIdParam := c.Query("userid")
	if len(userIdParam) == 0 {
		utils.WriteLog("log_user", "ShowEditUser.userIdParam is empty")
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	uid, err := strconv.Atoi(userIdParam)
	if err != nil {
		utils.WriteLog("log_user", "strconv.Atoi err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	userInfo, err := UserModel.GetUserInfo(uid)
	if err != nil {
		utils.WriteLog("log_user", "get user info err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	rolelist, err := permissionModel.GetRoleList()
	if err != nil {
		utils.WriteLog("log_user", "permission get role list err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	roleInfo, err := permissionModel.GetRoleInfo(userInfo.RoleId)
	if err != nil {
		utils.WriteLog("log_user", "permission get role info err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	menu, err := ctrl.getMenu(roleId.(int))
	if err != nil {
		utils.WriteLog("log_user", "get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	c.HTML(200, "permission/edituser.html", gin.H{
		"username":   userName,
		"moduleName": "permission",
		"ctrName":    "user",
		"ctrNameZn":  "编辑用户",
		"rolelist":   rolelist,
		"roleinfo":   roleInfo,
		"userinfo":   userInfo,
		"menu":       menu,
	})
}

func (ctrl PermissionController) AddUser(c *gin.Context) {
	var adduserForm forms.AdduserForm
	if c.BindJSON(&adduserForm) != nil {
		c.JSON(406, gin.H{"code": "0", "form": adduserForm})
		c.Abort()
		return
	}
	_, err := UserModel.AddUser(adduserForm)
	if err != nil {
		utils.WriteLog("log_user", err)
		c.JSON(200, gin.H{"code": "0", "desc": "fail"})
	} else {
		c.JSON(200, gin.H{"code": "1", "desc": "success"})
	}
	return
}

func (ctrl PermissionController) EditUser(c *gin.Context) {
	var edituserForm forms.EditUserForm
	if c.BindJSON(&edituserForm) != nil {
		c.JSON(406, gin.H{"code": "0", "form": edituserForm})
		c.Abort()
		return
	}
	_, err := UserModel.EditUser(edituserForm)
	if err != nil {
		utils.WriteLog("log_user", err)
		c.JSON(200, gin.H{"code": "0", "desc": "fail"})
	} else {
		c.JSON(200, gin.H{"code": "1", "desc": "success"})
	}
	return
}

func (ctrl PermissionController) ShowAddRole(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")
	menu, err := ctrl.getMenu(roleId.(int))
	if err != nil {
		utils.WriteLog("log_permission", "get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	c.HTML(200, "permission/addrole.html", gin.H{
		"username":   userName,
		"moduleName": "permission",
		"ctrName":    "user",
		"ctrNameZn":  "添加角色",
		"menu":       menu,
	})
}

func (ctrl PermissionController) AddRole(c *gin.Context) {
	var addroleForm forms.AddRoleForm
	if c.BindJSON(&addroleForm) != nil {
		c.JSON(406, gin.H{"code": "0", "form": addroleForm})
		c.Abort()
		return
	}
	_, err := permissionModel.AddRole(addroleForm)
	if err != nil {
		utils.WriteLog("log_permission", "permission add role err, err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "fail"})
	} else {
		c.JSON(200, gin.H{"code": "1", "desc": "success"})
	}
	return
}

func (ctrl PermissionController) ShowEditPermission(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	sessionRoleId := session.Get("roleid")
	roleIdParam := c.Query("roleid")
	if len(roleIdParam) == 0 {
		utils.WriteLog("log_user", "ShowEditPermission.roleIdParam is empty")
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	roleId, err := strconv.Atoi(roleIdParam)
	if err != nil {
		utils.WriteLog("log_permission", "strconv.Atoi err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	menu, err := ctrl.getMenu(sessionRoleId.(int))
	if err != nil {
		utils.WriteLog("log_permission", "get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	newMenuList, err := permissionModel.FormatMenuList(roleId)
	c.HTML(200, "permission/permissionlist.html", gin.H{
		"username":   userName,
		"moduleName": "permission",
		"ctrName":    "rolelist",
		"ctrNameZn":  "编辑权限",
		"roleId":     roleId,
		"menulist":   newMenuList,
		"menu":       menu,
	})
}

func (ctrl PermissionController) UpdatePermission(c *gin.Context) {
	ids := c.PostForm("ids")
	roleIdparam := c.PostForm("roleid")
	roleId, _ := strconv.Atoi(roleIdparam)
	_, err := permissionModel.UpdatePermission(roleId, ids)
	if err != nil {
		utils.WriteLog("log_permission", "permission update role menu relation err, err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "fail"})
	} else {
		c.JSON(200, gin.H{"code": "1", "desc": "success"})
	}
	return
}

func (ctrl PermissionController) setPaginator(c *gin.Context, per int, nums int64) *utils.Paginator {
	p := utils.NewPaginator(c, per, nums)
	return p
}

func (ctrl PermissionController) getMenu(roleId int) (menus []forms.ParentMenu, err error) {
	var menuModel = new(models.MenuModel)
	return menuModel.GetMenu(roleId)
}
