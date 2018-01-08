package controllers

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"mint-platform/platform/forms"
	"mint-platform/platform/models"
	"mint-platform/platform/utils"
	"strconv"
	"strings"
	//"fmt"
)

type LogController struct{}

var logModel = new(models.LogModel)

func (ctrl LogController) LogList(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")

	dt := c.DefaultQuery("dt", "")
	url := c.DefaultQuery("api", "")
	pageparam := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageparam)
	if err != nil {
		utils.WriteLog("log_log", "logList strconv atoi err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	pageSize := 25
	logList, err := logModel.ReadList(dt, url, page, pageSize)
	if err != nil {
		utils.WriteLog("log_log", "LogList get ReadList err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	total, err := logModel.GetLogTotal(dt, url)
	if err != nil {
		utils.WriteLog("log_log", "LogList GetLogTotal err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	menu, err := ctrl.getMenu(roleId.(int))
	if err != nil {
		utils.WriteLog("log_menu", "LogController LogList get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	paginator := ctrl.setPaginator(c, pageSize, total)
	c.HTML(200, "log/loglist.html", gin.H{
		"username":   userName,
		"moduleName": "log",
		"ctrName":    "loglist",
		"ctrNameZn":  "日志列表",
		"loglist":    logList,
		"paginator":  paginator,
		"menu":       menu,
	})
}

func (ctrl LogController) LoginList(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")

	userId := c.DefaultQuery("user_id", "")
	cellPhone := c.DefaultQuery("cell_phone", "")
	ip := c.DefaultQuery("ip", "")
	userId = strings.Replace(userId, " ", "", -1)
	cellPhone = strings.Replace(cellPhone, " ", "", -1)
	ip = strings.Replace(ip, " ", "", -1)
	pageparam := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageparam)
	if err != nil {
		utils.WriteLog("log_log", "LoginList strconv atoi err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	pageSize := 25
	logList, err := logModel.GetLoginList(userId, cellPhone, ip, page, pageSize)
	if err != nil {
		utils.WriteLog("log_log", "LoginList get ReadList err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	total, err := logModel.GetLoginLogTotal(userId, cellPhone, ip)
	if err != nil {
		utils.WriteLog("log_log", "LoginList GetLogTotal err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	menu, err := ctrl.getMenu(roleId.(int))
	if err != nil {
		utils.WriteLog("log_menu", "LogController LoginList get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	paginator := ctrl.setPaginator(c, pageSize, total)
	c.HTML(200, "log/loginlist.html", gin.H{
		"username":   userName,
		"moduleName": "log",
		"ctrName":    "loginlist",
		"ctrNameZn":  "登录列表",
		"loglist":    logList,
		"paginator":  paginator,
		"menu":       menu,
	})
}

func (ctrl LogController) setPaginator(c *gin.Context, per int, nums int64) *utils.Paginator {
	p := utils.NewPaginator(c, per, nums)
	return p
}

func (ctrl LogController) getMenu(roleId int) (menus []forms.ParentMenu, err error) {
	var menuModel = new(models.MenuModel)
	return menuModel.GetMenu(roleId)
}
