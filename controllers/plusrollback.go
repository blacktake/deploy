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

type PlusRollbackController struct{}

type DeployPlusHistoryVersion struct {
	Id      int    `json:"orderid"`
	Version string `json:"version"`
}

var PlusRollbackModel = new(models.PlusRollbackModel)

//var deployScriptPath = "/data/python/deploy_api"
//var deployScriptPath = "/root/tool/mia_deploy/api_repo"

func (ctrl PlusRollbackController) RollbackList(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")
	pageparam := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageparam)
	if err != nil {
		utils.WriteLog("log_plusrollback", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	pageSize := 20
	rollbackList, err := PlusRollbackModel.GetRollbackList(page, pageSize)
	if err != nil {
		utils.WriteLog("log_plusrollback", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	total, err := PlusRollbackModel.GetRollbackTotal()
	if err != nil {
		utils.WriteLog("log_plusrollback", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	paginator := ctrl.setPaginator(c, pageSize, total)
	menu, err := ctrl.getMenu(roleId.(int))
	if err != nil {
		utils.WriteLog("log_plusrollback", "get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	c.HTML(200, "plusrollback/rollbacklist.html", gin.H{
		"username":     userName,
		"moduleName":   "deploy",
		"ctrName":      "plusrollback",
		"ctrNameZn":    "plus回滚列表",
		"rollbackList": rollbackList,
		"paginator":    paginator,
		"menu":         menu,
	})
}

func (ctrl PlusRollbackController) ShowAddWorker(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")
	menu, err := ctrl.getMenu(roleId.(int))
	if err != nil {
		utils.WriteLog("log_plusrollback", "get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	c.HTML(200, "plusrollback/addrollback.html", gin.H{
		"username":   userName,
		"moduleName": "deploy",
		"ctrName":    "plusrollback",
		"ctrNameZn":  "添加回滚工单",
		"menu":       menu,
	})
}

func (ctrl PlusRollbackController) AddWorker(c *gin.Context) {
	var projectForm forms.PlusProjectForm
	if c.BindJSON(&projectForm) != nil {
		c.JSON(406, gin.H{"code": "0", "form": projectForm})
		c.Abort()
		return
	}
	_, err := PlusRollbackModel.AddWorker(projectForm)
	if err != nil {
		utils.WriteLog("log_plusrollback", err)
		c.JSON(200, gin.H{"code": "0", "desc": "fail"})
	} else {
		c.JSON(200, gin.H{"code": "1", "desc": "success"})
	}
	return
}

func (ctrl PlusRollbackController) ShowRollbackVersion(c *gin.Context) {
	groupName := c.PostForm("groupname")
	var plusprojectModel = new(models.PlusProjectModel)
	projectList, err := plusprojectModel.GetHistoryDeployVersion(groupName)
	if err != nil {
		utils.WriteLog("log_plusrollback", "GetHistoryDeployVersion  err , err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "get history deploy version err"})
		return
	}
	deployVersionList := make([]DeployPlusHistoryVersion, len(projectList))
	var deployVersion DeployPlusHistoryVersion
	for k, project := range projectList {
		deployVersion.Id = project.Id
		deployVersion.Version = project.DateDeployed.Format("20060102_150405")
		deployVersionList[k] = deployVersion
	}
	c.JSON(200, gin.H{"code": "1", "desc": "success", "projectList": deployVersionList})
	return
}

func (ctrl PlusRollbackController) CancelWorker(c *gin.Context) {
	orderidParam := c.PostForm("orderid")
	orderId, _ := strconv.Atoi(orderidParam)
	_, err := PlusRollbackModel.CancelWorker(orderId)
	if err != nil {
		utils.WriteLog("log_plusrollback", "cancel rollback order err, err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "fail"})
	} else {
		c.JSON(200, gin.H{"code": "1", "desc": "success"})
	}
	return
}

func (ctrl PlusRollbackController) ShowRollback(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")
	orderidParam := c.Query("id")
	orderId, _ := strconv.Atoi(orderidParam)
	_, err := PlusRollbackModel.GetAvailableRollbackOrderInfo(orderId)
	if err != nil {
		utils.WriteLog("log_plusrollback", "get available rollback order err， err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	menu, err := ctrl.getMenu(roleId.(int))
	if err != nil {
		utils.WriteLog("log_plusrollback", "get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	c.HTML(200, "plusrollback/rollback.html", gin.H{
		"username":   userName,
		"orderId":    orderId,
		"moduleName": "deploy",
		"ctrName":    "plusrollback",
		"ctrNameZn":  "代码回滚",
		"menu":       menu,
	})
}

func (ctrl PlusRollbackController) Rollback(c *gin.Context) {
	orderidParam := c.PostForm("orderid")
	orderId, _ := strconv.Atoi(orderidParam)

	redisHander := utils.CreateRedis("lockRedis")
	if redisHander.Pool == nil {
		utils.WriteLog("log_plusrollback", "redis连接为空")
		c.JSON(200, gin.H{"code": "0", "desc": "redis conn is empty"})
		return
	}
	lock, err := redisHander.SetNx("pluslock", 1)
	if err != nil {
		utils.WriteLog("log_plusrollback", "redis setnx err, err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "redis setnx err"})
		return
	}
	if lock == false {
		c.JSON(200, gin.H{"code": "4", "desc": "some on online deploying"})
		return
	}
	rollbackOrder, err := PlusRollbackModel.GetRollbackOrderInfo(orderId)
	if err != nil {
		utils.WriteLog("log_plusrollback", "get rollback order err， err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "no project"})
		return
	}
	createStyleTime := rollbackOrder.DateCreate.Format("20060102_150405")
	t := time.Now()
	rollbackTime := t.Format("2006-01-02 15:04:05")
	deployLogName := fmt.Sprintf("%s_rollback_%s_%s.log", rollbackOrder.UserName, rollbackOrder.GroupName, createStyleTime)
	logFilePath := fmt.Sprintf("%s/temlogs/%s", plusdeployScriptPath, deployLogName)
	//cmdStr := fmt.Sprintf("`python %s/admin_sync_deploy_api.py -t %s -p api -R %s > %s &`", deployScriptPath, rollbackOrder.GroupName, rollbackOrder.Version, logFilePath)
	cmdStr := fmt.Sprintf("`%s/plusdeploy  -t %s -p mia_plus -R %s > %s &`", plusdeployScriptPath, rollbackOrder.GroupName, rollbackOrder.Version, logFilePath)
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		utils.WriteLog("log_plusrollback", "cmd stdoutPipe err,err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "cmd stdoutPipe err"})
		return
	}
	//执行命令
	if err := cmd.Start(); err != nil {
		utils.WriteLog("log_plusrollback", "cmd Start err,err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "cmd start err"})
		return
	}
	if _, err := ioutil.ReadAll(stdout); err != nil {
		utils.WriteLog("log_plusrollback", "ioutil readall err，err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "ioutil ReadAll err"})
		return
	}
	if err := cmd.Wait(); err != nil {
		utils.WriteLog("log_plusrollback", "cmd wait err, err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "cmd Wait err"})
		return
	}
	rollbackInfo, err := PlusRollbackModel.UpdateRollbackOrderInfo(orderId, t)
	if err != nil {
		utils.WriteLog("log_plusrollback", "rollback update rollback order info err ,err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "rollback update rollback order info err"})
		return
	}
	_, err = redisHander.Delete("pluslock")
	if err != nil {
		utils.WriteLog("log_plusrollback", "redis delere key err, err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "redis delere key err"})
		return
	}
	rollbackStr := fmt.Sprintf("[plus回滚]-[%s]于%s在%s分组回滚了分支%s, 理由:[%s]", rollbackInfo.UserName, rollbackTime, rollbackInfo.GroupName, rollbackInfo.Version, rollbackInfo.TaskName)
	result, err := utils.Get("http://wxpush.miyabaobei.com", rollbackStr)
	if err != nil {
		utils.WriteLog("log_plusrollback", "get wxpush err,err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "get wxpush err"})
		return
	}
	resultStr := &ResultInfo{}
	err = json.Unmarshal(result, resultStr)
	if err != nil {
		utils.WriteLog("log_plusrollback", "json unmarshal err, err:", err)
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

func (ctrl PlusRollbackController) ShowInfo(c *gin.Context) {
	orderidParam := c.PostForm("orderid")
	orderId, _ := strconv.Atoi(orderidParam)
	rollbackInfo, err := PlusRollbackModel.GetRollbackOrderInfo(orderId)
	if err != nil {
		utils.WriteLog("log_plusrollback", "get rollback order err， err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "no rollback"})
		return
	}
	createStyleTime := rollbackInfo.DateCreate.Format("20060102_150405")
	deployLogName := fmt.Sprintf("%s_rollback_%s_%s.log", rollbackInfo.UserName, rollbackInfo.GroupName, createStyleTime)
	logFilePath := fmt.Sprintf("%s/temlogs/%s", plusdeployScriptPath, deployLogName)
	fileInfo, err := utils.ReadLine(logFilePath)
	if err != nil {
		utils.WriteLog("log_plusrollback", "read file err， err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "read file err"})
		return
	}
	c.JSON(200, gin.H{"code": "1", "desc": string(fileInfo)})
}

func (ctrl PlusRollbackController) setPaginator(c *gin.Context, per int, nums int64) *utils.Paginator {
	p := utils.NewPaginator(c, per, nums)
	return p
}

func (ctrl PlusRollbackController) getMenu(roleId int) (menus []forms.ParentMenu, err error) {
	var menuModel = new(models.MenuModel)
	return menuModel.GetMenu(roleId)
}
