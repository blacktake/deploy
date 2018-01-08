package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mint-platform/platform/forms"
	"mint-platform/platform/models"
	"mint-platform/platform/utils"
	"os/exec"
	"strconv"
	"time"
)

type RollbackController struct{}

type DeployHistoryVersion struct {
	Id      int    `json:"orderid"`
	Version string `json:"version"`
}

var RollbackModel = new(models.RollbackModel)

//var deployScriptPath = "/data/python/deploy_api"
//var deployScriptPath = "/root/tool/mia_deploy/api_repo"

func (ctrl RollbackController) RollbackList(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")
	pageparam := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageparam)
	if err != nil {
		utils.WriteLog("log_rollback", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	pageSize := 20
	rollbackList, err := RollbackModel.GetRollbackList(page, pageSize)
	if err != nil {
		utils.WriteLog("log_rollback", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	total, err := RollbackModel.GetRollbackTotal()
	if err != nil {
		utils.WriteLog("log_rollback", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	paginator := ctrl.setPaginator(c, pageSize, total)
	menu, err := ctrl.getMenu(roleId.(int))
	if err != nil {
		utils.WriteLog("log_rollback", "get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	c.HTML(200, "rollback/rollbacklist.html", gin.H{
		"username":     userName,
		"moduleName":   "deploy",
		"ctrName":      "rollback",
		"ctrNameZn":    "回滚列表",
		"rollbackList": rollbackList,
		"paginator":    paginator,
		"menu":         menu,
	})
}

func (ctrl RollbackController) ShowAddWorker(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")
	menu, err := ctrl.getMenu(roleId.(int))
	if err != nil {
		utils.WriteLog("log_rollback", "get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	c.HTML(200, "rollback/addrollback.html", gin.H{
		"username":   userName,
		"moduleName": "deploy",
		"ctrName":    "rollback",
		"ctrNameZn":  "添加回滚工单",
		"menu":       menu,
	})
}

func (ctrl RollbackController) AddWorker(c *gin.Context) {
	var projectForm forms.ProjectForm
	if c.BindJSON(&projectForm) != nil {
		c.JSON(406, gin.H{"code": "0", "form": projectForm})
		c.Abort()
		return
	}
	_, err := RollbackModel.AddWorker(projectForm)
	if err != nil {
		utils.WriteLog("log_rollback", err)
		c.JSON(200, gin.H{"code": "0", "desc": "fail"})
	} else {
		c.JSON(200, gin.H{"code": "1", "desc": "success"})
	}
	return
}

func (ctrl RollbackController) ShowRollbackVersion(c *gin.Context) {
	groupName := c.PostForm("groupname")
	var projectModel = new(models.ProjectModel)
	projectList, err := projectModel.GetHistoryDeployVersion(groupName)
	if err != nil {
		utils.WriteLog("log_rollback", "GetHistoryDeployVersion  err , err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "get history deploy version err"})
		return
	}
	deployVersionList := make([]DeployHistoryVersion, len(projectList))
	var deployVersion DeployHistoryVersion
	for k, project := range projectList {
		deployVersion.Id = project.Id
		deployVersion.Version = project.DateDeployed.Format("20060102_150405")
		deployVersionList[k] = deployVersion
	}
	c.JSON(200, gin.H{"code": "1", "desc": "success", "projectList": deployVersionList})
	return
}

func (ctrl RollbackController) CancelWorker(c *gin.Context) {
	orderidParam := c.PostForm("orderid")
	orderId, _ := strconv.Atoi(orderidParam)
	_, err := RollbackModel.CancelWorker(orderId)
	if err != nil {
		utils.WriteLog("log_rollback", "cancel rollback order err, err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "fail"})
	} else {
		c.JSON(200, gin.H{"code": "1", "desc": "success"})
	}
	return
}

func (ctrl RollbackController) ShowRollback(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")
	orderidParam := c.Query("id")
	orderId, _ := strconv.Atoi(orderidParam)
	_, err := RollbackModel.GetAvailableRollbackOrderInfo(orderId)
	if err != nil {
		utils.WriteLog("log_rollback", "get available rollback order err， err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	menu, err := ctrl.getMenu(roleId.(int))
	if err != nil {
		utils.WriteLog("log_rollback", "get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	c.HTML(200, "rollback/rollback.html", gin.H{
		"username":   userName,
		"orderId":    orderId,
		"moduleName": "deploy",
		"ctrName":    "rollback",
		"ctrNameZn":  "代码回滚",
		"menu":       menu,
	})
}

func (ctrl RollbackController) Rollback(c *gin.Context) {
	orderidParam := c.PostForm("orderid")
	orderId, _ := strconv.Atoi(orderidParam)

	redisHander := utils.CreateRedis("lockRedis")
	if redisHander.Pool == nil {
		utils.WriteLog("log_rollback", "redis连接为空")
		c.JSON(200, gin.H{"code": "0", "desc": "redis conn is empty"})
		return
	}
	lock, err := redisHander.SetNx("lock", 1)
	if err != nil {
		utils.WriteLog("log_rollback", "redis setnx err, err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "redis setnx err"})
		return
	}
	if lock == false {
		c.JSON(200, gin.H{"code": "4", "desc": "some on online deploying"})
		return
	}
	rollbackOrder, err := RollbackModel.GetRollbackOrderInfo(orderId)
	if err != nil {
		utils.WriteLog("log_rollback", "get rollback order err， err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "no project"})
		return
	}
	createStyleTime := rollbackOrder.DateCreate.Format("20060102_150405")
	t := time.Now()
	rollbackTime := t.Format("2006-01-02 15:04:05")
	deployLogName := fmt.Sprintf("%s_rollback_%s_%s.log", rollbackOrder.UserName, rollbackOrder.GroupName, createStyleTime)
	logFilePath := fmt.Sprintf("%s/temlogs/%s", deployScriptPath, deployLogName)
	//cmdStr := fmt.Sprintf("`python %s/admin_sync_deploy_api.py -t %s -p api -R %s > %s &`", deployScriptPath, rollbackOrder.GroupName, rollbackOrder.Version, logFilePath)
	cmdStr := fmt.Sprintf("`%s/deploy  -t %s -p api -R %s > %s &`", deployScriptPath, rollbackOrder.GroupName, rollbackOrder.Version, logFilePath)
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		utils.WriteLog("log_rollback", "cmd stdoutPipe err,err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "cmd stdoutPipe err"})
		return
	}
	//执行命令
	if err := cmd.Start(); err != nil {
		utils.WriteLog("log_rollback", "cmd Start err,err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "cmd start err"})
		return
	}
	if _, err := ioutil.ReadAll(stdout); err != nil {
		utils.WriteLog("log_rollback", "ioutil readall err，err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "ioutil ReadAll err"})
		return
	}
	if err := cmd.Wait(); err != nil {
		utils.WriteLog("log_rollback", "cmd wait err, err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "cmd Wait err"})
		return
	}
	rollbackInfo, err := RollbackModel.UpdateRollbackOrderInfo(orderId, t)
	if err != nil {
		utils.WriteLog("log_rollback", "rollback update rollback order info err ,err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "rollback update rollback order info err"})
		return
	}
	_, err = redisHander.Delete("lock")
	if err != nil {
		utils.WriteLog("log_rollback", "redis delere key err, err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "redis delere key err"})
		return
	}
	rollbackStr := fmt.Sprintf("[%s]于%s在%s分组回滚了分支%s, 理由:[%s]", rollbackInfo.UserName, rollbackTime, rollbackInfo.GroupName, rollbackInfo.Version, rollbackInfo.TaskName)
	result, err := utils.Get("http://wxpush.miyabaobei.com", rollbackStr)
	if err != nil {
		utils.WriteLog("log_rollback", "get wxpush err,err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "get wxpush err"})
		return
	}
	resultStr := &ResultInfo{}
	err = json.Unmarshal(result, resultStr)
	if err != nil {
		utils.WriteLog("log_rollback", "json unmarshal err, err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "json unmarshal err"})
		return
	}
	if resultStr.Code == 1 {
		c.JSON(200, gin.H{"code": "1", "desc": "success"})
	} else {
		c.JSON(200, gin.H{"code": "0", "desc": "fail"})
	}
	return
}

func (ctrl RollbackController) ShowInfo(c *gin.Context) {
	orderidParam := c.PostForm("orderid")
	orderId, _ := strconv.Atoi(orderidParam)
	rollbackInfo, err := RollbackModel.GetRollbackOrderInfo(orderId)
	if err != nil {
		utils.WriteLog("log_rollback", "get rollback order err， err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "no rollback"})
		return
	}
	createStyleTime := rollbackInfo.DateCreate.Format("20060102_150405")
	deployLogName := fmt.Sprintf("%s_rollback_%s_%s.log", rollbackInfo.UserName, rollbackInfo.GroupName, createStyleTime)
	logFilePath := fmt.Sprintf("%s/temlogs/%s", deployScriptPath, deployLogName)
	fileInfo, err := utils.ReadLine(logFilePath)
	if err != nil {
		utils.WriteLog("log_rollback", "read file err， err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "read file err"})
		return
	}
	c.JSON(200, gin.H{"code": "1", "desc": string(fileInfo)})
}

func (ctrl RollbackController) setPaginator(c *gin.Context, per int, nums int64) *utils.Paginator {
	p := utils.NewPaginator(c, per, nums)
	return p
}

func (ctrl RollbackController) getMenu(roleId int) (menus []forms.ParentMenu, err error) {
	var menuModel = new(models.MenuModel)
	return menuModel.GetMenu(roleId)
}
