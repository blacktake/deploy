package controllers

import (
	"bytes"
	"fmt"
	//"io"
	//"strings"
	"encoding/json"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	//"io/ioutil"
	"mint-platform/platform/forms"
	"mint-platform/platform/models"
	"mint-platform/platform/utils"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"time"
	//"bufio"
)

type ResultInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ProjectController struct{}

var projectModel = new(models.ProjectModel)

//var deployScriptPath = "/data/go/workspace/src/mint-deploy/deploy"
var deployScriptPath = "/root/tool/mia_release"

//var deployScriptPath = "/data/python/deploy_api"
//var deployScriptPath = "/root/tool/mia_deploy/api_repo"

func (ctrl ProjectController) ProjectHandler(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")
	pageparam := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageparam)
	if err != nil {
		utils.WriteLog("log_project", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	pageSize := 20
	projectList, err := projectModel.GetProjectList(page, pageSize)
	if err != nil {
		utils.WriteLog("log_project", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	total, err := projectModel.GetProjectTotal()
	if err != nil {
		utils.WriteLog("log_project", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	paginator := ctrl.setPaginator(c, pageSize, total)
	menu, err := ctrl.getMenu(roleId.(int))
	if err != nil {
		utils.WriteLog("log_project", "get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}

	promissionTab := 0
	roleList := []int{1, 2}
	for _, val := range roleList {
		if val == roleId {
			promissionTab = 1
		}
	}
	c.HTML(200, "project/projectlist.html", gin.H{
		"username":      userName,
		"moduleName":    "deploy",
		"ctrName":       "project",
		"ctrNameZn":     "发布列表",
		"promissionTab": promissionTab,
		"projectList":   projectList,
		"paginator":     paginator,
		"menu":          menu,
	})
}

func (ctrl ProjectController) ShowAddWorker(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")

	chs := make(chan int)
	defer close(chs)
	var branchstr string
	go func() {
		idrsafilelist, err := utils.LoadConfig("fileconfig", "IDRSAFILE")
		if err != nil {
			chs <- 0
			return
		}
		//测试环境
		//idraspath := idrsafilelist["path_test"]
		//正式环境
		idraspath := idrsafilelist["path"]
		//fmt.Println(idraspath)
		apilistpath, err := utils.LoadConfig("fileconfig", "APIFILEPATH")
		if err != nil {
			chs <- 0
			return
		}
		//测试环境
		//apifilePath := apilistpath["path_test"]
		//正式环境
		apifilePath := apilistpath["path"]
		//fmt.Println(apifilePath)
		gitbranchCmd := fmt.Sprintf("cd %s;ssh-agent bash -c 'ssh-add %s;git fetch origin -p';git branch -r | cut -d'/' -f2 | sed '1d'", apifilePath, idraspath)
		//fmt.Println(gitbranchCmd)
		res, err := utils.ExcuteCmd(gitbranchCmd)
		if err != nil {
			chs <- 0
		} else {
			chs <- 1
			branchstr = res
		}
	}()
	result := <-chs
	if result == 0 {
		utils.WriteLog("log_project", "get git branch err")
		//		c.HTML(200, "empty_page.html", gin.H{})
		//		return
	}
	branchstr = "aaa"
	branchList := strings.Split(branchstr, "\n")
	menu, err := ctrl.getMenu(roleId.(int))
	if err != nil {
		utils.WriteLog("log_project", "get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	c.HTML(200, "project/addproject.html", gin.H{
		"username":   userName,
		"moduleName": "deploy",
		"branchList": branchList[0 : len(branchList)-1],
		"ctrName":    "project",
		"ctrNameZn":  "添加工单",
		"menu":       menu,
	})
}

func (ctrl ProjectController) AddWorker(c *gin.Context) {
	var projectForm forms.ProjectForm
	if c.BindJSON(&projectForm) != nil {
		c.JSON(406, gin.H{"code": "0", "form": projectForm})
		c.Abort()
		return
	}
	_, err := projectModel.AddWorker(projectForm)
	if err != nil {
		utils.WriteLog("log_project", err)
		c.JSON(200, gin.H{"code": "0", "desc": "fail"})
	} else {
		c.JSON(200, gin.H{"code": "1", "desc": "success"})
	}
	return
}

func (ctrl ProjectController) CancelWorker(c *gin.Context) {
	orderidParam := c.PostForm("orderid")
	orderId, _ := strconv.Atoi(orderidParam)
	_, err := projectModel.CancelWorker(orderId)
	if err != nil {
		utils.WriteLog("log_project", err)
		c.JSON(200, gin.H{"code": "0", "desc": "fail"})
	} else {
		c.JSON(200, gin.H{"code": "1", "desc": "success"})
	}
	return
}

func (ctrl ProjectController) AccessWorker(c *gin.Context) {
	orderidParam := c.PostForm("orderid")
	orderId, _ := strconv.Atoi(orderidParam)
	session := sessions.Default(c)
	userName := session.Get("username").(string)
	_, err := projectModel.AccessWorker(orderId, userName)
	if err != nil {
		utils.WriteLog("log_project", err)
		c.JSON(200, gin.H{"code": "0", "desc": "fail"})
	} else {
		c.JSON(200, gin.H{"code": "1", "desc": "success"})
	}
	return
}
func (ctrl ProjectController) ShowPublish(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")
	orderidParam := c.Query("id")
	orderId, _ := strconv.Atoi(orderidParam)
	_, err := projectModel.GetAvailableDeployOrderInfo(orderId)
	if err != nil {
		utils.WriteLog("log_project", "get available deploy order err， err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	menu, err := ctrl.getMenu(roleId.(int))
	if err != nil {
		utils.WriteLog("log_project", "get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	c.HTML(200, "project/publish.html", gin.H{
		"username":   userName,
		"orderId":    orderId,
		"moduleName": "deploy",
		"ctrName":    "project",
		"ctrNameZn":  "代码发布",
		"menu":       menu,
	})
}

func (ctrl ProjectController) Publish(c *gin.Context) {
	orderidParam := c.PostForm("orderid")
	orderId, _ := strconv.Atoi(orderidParam)

	redisHander := utils.CreateRedis("lockRedis")
	if redisHander.Pool == nil {
		utils.WriteLog("log_project", "redis连接为空")
		c.JSON(200, gin.H{"code": "0", "desc": "redis conn is empty"})
		return
	}
	lock, err := redisHander.SetNx("lock", 1)
	if err != nil {
		utils.WriteLog("log_project", "redis setnx err, err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "redis setnx err"})
		return
	}
	if lock == false {
		c.JSON(200, gin.H{"code": "4", "desc": "some on online deploying"})
		return
	}
	project, err := projectModel.GetDeployOrderInfo(orderId)

	if err != nil {
		utils.WriteLog("log_project", "get deploy order err， err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "no project"})
		return
	}
	createStyleTime := project.DateCreate.Format("20060102_150405")
	t := time.Now()
	deployStyleTime := t.Format("20060102_150405")
	deployTime := t.Format("2006-01-02 15:04:05")
	deployLogName := fmt.Sprintf("%s_deploy_%s_%s.log", project.UserName, project.GroupName, createStyleTime)
	logFilePath := fmt.Sprintf("%s/temlogs/%s", deployScriptPath, deployLogName)
	//cmdStr := fmt.Sprintf("`python %s/admin_sync_deploy_api.py -t %s -p api -D %s -v %s > %s &`", deployScriptPath, project.GroupName, deployStyleTime, project.Version, logFilePath)
	cmdStr := fmt.Sprintf("`%s/deploy -t %s -p api -D %s -v %s > %s &`", deployScriptPath, project.GroupName, deployStyleTime, project.Version, logFilePath)
	//fmt.Println(cmdStr)
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	var out bytes.Buffer
	cmd.Stdout = &out //标准输出
	err = cmd.Run()   //运行指令 ，做判断
	if err != nil {
		utils.WriteLog("log_project", "cmd Run err,err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "cmd run err"})
		return
	}
	projectInfo, err := projectModel.UpdateDeployOrderInfo(orderId, t)
	if err != nil {
		utils.WriteLog("log_project", "project update deploy order info err ,err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "project update deploy order info err"})
		return
	}
	_, err = redisHander.Delete("lock")
	if err != nil {
		utils.WriteLog("log_project", "redis delere key err, err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "redis delere key err"})
		return
	}
	deployStr := fmt.Sprintf("[%s]于%s在%s分组发布了分支%s, 理由:[%s]", projectInfo.UserName, deployTime, projectInfo.GroupName, projectInfo.Version, projectInfo.TaskName)
	deployStr = strings.Replace(deployStr, "&", " and ", -1)
	requestUrl := "http://wxpush.miyabaobei.com"
	parseRequesUrl, err := url.Parse(requestUrl)
	if err != nil {
		utils.WriteLog("log_project", "url parse wx url err,err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "url parse wx url err"})
		return
	}
	// 需要添加的 GET 参数
	extraParams := url.Values{
		"q": {deployStr},
	}
	// 更改 URL Struct 中的 RawQuery 为 Encode 后的 Query string
	parseRequesUrl.RawQuery = extraParams.Encode()
	requestUrl = parseRequesUrl.String()

	result, err := utils.GetUrl(requestUrl)
	if err != nil {
		utils.WriteLog("log_project", "get wxpush err,err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "get wxpush err"})
		return
	}
	resultStr := &ResultInfo{}
	err = json.Unmarshal(result, resultStr)
	if err != nil {
		utils.WriteLog("log_project", "json unmarshal err, err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "json unmarshal err"})
		return
	}
	if resultStr.Code == 1 {
		utils.MakeTemplateToMail(project.UserName, project.FunctionalIntroduction, project.EmailList, project.Auditor, project.GroupName)
		c.JSON(200, gin.H{"code": "1", "desc": "success"})
	} else {
		c.JSON(200, gin.H{"code": "0", "desc": "fail"})
	}
	return
}

func (ctrl ProjectController) ShowInfo(c *gin.Context) {
	orderidParam := c.PostForm("orderid")
	orderId, _ := strconv.Atoi(orderidParam)
	project, err := projectModel.GetDeployOrderInfo(orderId)
	if err != nil {
		utils.WriteLog("log_project", "get deploy order err， err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "no project"})
		return
	}
	createStyleTime := project.DateCreate.Format("20060102_150405")
	deployLogName := fmt.Sprintf("%s_deploy_%s_%s.log", project.UserName, project.GroupName, createStyleTime)
	logFilePath := fmt.Sprintf("%s/temlogs/%s", deployScriptPath, deployLogName)
	fileInfo, err := utils.ReadLine(logFilePath)
	if err != nil {
		utils.WriteLog("log_project", "read file err， err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "read file err"})
		return
	}
	c.JSON(200, gin.H{"code": "1", "desc": string(fileInfo)})
}

func (ctrl ProjectController) setPaginator(c *gin.Context, per int, nums int64) *utils.Paginator {
	p := utils.NewPaginator(c, per, nums)
	return p
}

func (ctrl ProjectController) getMenu(roleId int) (menus []forms.ParentMenu, err error) {
	var menuModel = new(models.MenuModel)
	return menuModel.GetMenu(roleId)
}
