package controllers

import (
	//"html"
	//"io"
	//"strings"
	"encoding/json"
	"mint-platform/platform/forms"
	"mint-platform/platform/models"
	"mint-platform/platform/utils"
	"strconv"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	//"fmt"
	//"reflect"
)

type MenuController struct{}

var menuModel = new(models.MenuModel)

func (ctrl MenuController) MenuList(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")
	parentMenuList, err := menuModel.GetParentMenuList()
	if err != nil {
		utils.WriteLog("log_menu", "get parent menu list err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	menuList, err := menuModel.GetMenuList()
	if err != nil {
		utils.WriteLog("log_menu", "get menu list err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	newMenuList := menuModel.GetAllMenuList(menuList)
	menu, err := ctrl.getMenu(roleId.(int))
	if err != nil {
		utils.WriteLog("log_menu", "get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	c.HTML(200, "menu/list.html", gin.H{
		"username":       userName,
		"moduleName":     "menu",
		"ctrName":        "menu",
		"ctrNameZn":      "菜单列表",
		"parentMenuList": parentMenuList,
		"menuList":       newMenuList,
		"menu":           menu,
	})
}

func (ctrl MenuController) AddMenu(c *gin.Context) {
	var menuForm forms.MenuForm
	if c.BindJSON(&menuForm) != nil {
		c.JSON(406, gin.H{"code": "0", "form": menuForm})
		c.Abort()
		return
	}
	_, err := menuModel.AddMenu(menuForm)
	if err != nil {
		utils.WriteLog("log_menu", err)
		c.JSON(200, gin.H{"code": "0", "desc": "fail"})
	} else {
		c.JSON(200, gin.H{"code": "1", "desc": "success"})
	}
	return
}

func (ctrl MenuController) ChangeMenuOrder(c *gin.Context) {
	params := c.PostForm("data")
	var menuidlist []forms.MenuId
	err := json.Unmarshal([]byte(params), &menuidlist)
	if err != nil {
		utils.WriteLog("log_menu", "json unmarshal err, err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "json unmarshal err"})
		return
	}
	err = menuModel.ChangeMenuOrder(menuidlist)
	if err != nil {
		utils.WriteLog("log_menu", err)
		c.JSON(200, gin.H{"code": "0", "desc": err.Error()})
	} else {
		c.JSON(200, gin.H{"code": "1", "desc": "success"})
	}
	return
}

func (ctrl MenuController) DeleteMenu(c *gin.Context) {
	menuIdparams := c.PostForm("menuid")
	menuId, _ := strconv.Atoi(menuIdparams)
	err := menuModel.DeleteMenu(menuId)
	if err != nil {
		utils.WriteLog("log_menu", "delete menu err,err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": err.Error()})
	} else {
		c.JSON(200, gin.H{"code": "1", "desc": "success"})
	}
	return
}

func (ctrl MenuController) ShowEditMenu(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")
	idparam := c.Query("id")
	id, err := strconv.Atoi(idparam)
	if err != nil {
		utils.WriteLog("log_menu", "strconv atoi err,err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	parentMenuList, err := menuModel.GetParentMenuList()
	if err != nil {
		utils.WriteLog("log_menu", "get parent menu list err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	menuInfo, err := menuModel.GetMenuInfo(id)
	if err != nil {
		utils.WriteLog("log_menu", "get menu info err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	menu, err := ctrl.getMenu(roleId.(int))
	if err != nil {
		utils.WriteLog("log_menu", "get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	c.HTML(200, "menu/edit.html", gin.H{
		"username":       userName,
		"moduleName":     "menu",
		"ctrName":        "menu",
		"ctrNameZn":      "修改菜单",
		"parentMenuList": parentMenuList,
		"menuInfo":       menuInfo,
		"menu":           menu,
	})
}

func (ctrl MenuController) EditMenu(c *gin.Context) {
	var menuForm forms.AllMenuForm
	if c.BindJSON(&menuForm) != nil {
		c.JSON(406, gin.H{"code": "0", "form": menuForm})
		c.Abort()
		return
	}
	_, err := menuModel.EditMenu(menuForm)
	if err != nil {
		utils.WriteLog("log_menu", err)
		c.JSON(200, gin.H{"code": "0", "desc": err.Error()})
	} else {
		c.JSON(200, gin.H{"code": "1", "desc": "success"})
	}
	return
}

func (ctrl MenuController) getMenu(roleId int) (menus []forms.ParentMenu, err error) {
	return menuModel.GetMenu(roleId)
}
