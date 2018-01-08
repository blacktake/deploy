package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mint-platform/platform/forms"
	"mint-platform/platform/models"
	"mint-platform/platform/utils"
	"net/url"
	"os/exec"
	"regexp"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type ToolController struct{}

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var restartPhpPath = "/root/tool/script"

//var restartPhpPath = "/data/bin/script"

func (ctrl ToolController) Delvarnish(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")
	menu, err := ctrl.getMenu(roleId.(int))
	if err != nil {
		utils.WriteLog("log_tool", "ToolController Delvarnish get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	c.HTML(200, "tool/delvarnish.html", gin.H{
		"username":   userName,
		"moduleName": "tool",
		"ctrName":    "delvarnish",
		"ctrNameZn":  "清除varnish",
		"menu":       menu,
	})
}

func (ctrl ToolController) ExcuteDelvarnish(c *gin.Context) {
	urlstr := c.PostForm("url")
	var digitsRegexp = regexp.MustCompile(`/.*?/.*?/`)
	regFlag := digitsRegexp.MatchString(urlstr)
	if !regFlag {
		c.JSON(200, gin.H{"code": "0", "desc": "url is wrong"})
		return
	}
	urlstr = strings.Replace(urlstr, " ", "", -1)
	urlstr = strings.Replace(urlstr, "\n", "", -1)
	delvarnishUrl := "http://10.1.51.244:8076"
	//delvarnishUrl := "http://127.0.0.1:8070"
	apiUrl := fmt.Sprintf("http://api.miyabaobei.com%s", urlstr)
	params := url.Values{}
	params.Add("url", apiUrl)
	params.Add("secure_key", "mia.com")
	result, err := utils.PostForm(delvarnishUrl, params)
	if err != nil {
		utils.WriteLog("log_tool", "post delvarnish err,err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "post delvarnish err"})
		return
	}
	resultStr := &Result{}
	err = json.Unmarshal(result, resultStr)
	if err != nil {
		utils.WriteLog("log_tool", "tool ExcuteDelvarnish json unmarshal err, err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "json unmarshal err"})
		return
	}
	if resultStr.Code == 1 {
		c.JSON(200, gin.H{"code": "1", "desc": "success"})
	} else {
		c.JSON(200, gin.H{"code": "0", "desc": resultStr.Message})
	}
	return
}

func (ctrl ToolController) ShowRestartphp(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")
	menu, err := ctrl.getMenu(roleId.(int))
	if err != nil {
		utils.WriteLog("log_tool", "ToolController ShowRestartphp get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	c.HTML(200, "tool/showrestartphp.html", gin.H{
		"username":   userName,
		"moduleName": "tool",
		"ctrName":    "showrestartphp",
		"ctrNameZn":  "重启php",
		"menu":       menu,
	})
}

func (ctrl ToolController) Restartphp(c *gin.Context) {
	typestr := c.PostForm("type")
	var fileName = ""
	if typestr == "api" {
		fileName = "restart_api_php.sh"
	} else if typestr == "cart" {
		fileName = "restart_api_cart_php.sh"
	} else if typestr == "order" {
		fileName = "restart_api_order_php.sh"
	} else if typestr == "all" {
		fileName = "restart_all_php.sh"
	} else if typestr == "recommend" {
		fileName = "restart_api_recommend_php.sh"
	} else if typestr == "thirdservice" {
		fileName = "restart_api_third_php.sh"
	}

	cmdStr := fmt.Sprintf("sh %s/%s", restartPhpPath, fileName)
	chs := make(chan int)
	var res string
	go func() {
		cmd := exec.Command("/bin/sh", "-c", cmdStr)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			utils.WriteLog("log_tool", "cmd stdoutPipe err,err:", err)
			chs <- 0
			return
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			utils.WriteLog("log_tool", "cmd StderrPipe err,err:", err)
			chs <- 0
			return
		}
		//执行命令
		if err := cmd.Start(); err != nil {
			utils.WriteLog("log_tool", "cmd Start err,err:", err)
			chs <- 0
			return
		}
		bytesErr, err := ioutil.ReadAll(stderr)
		if err != nil {
			utils.WriteLog("log_tool", "ioutil ReadAll stderr err,err:", err)
			chs <- 0
			return
		}
		if len(bytesErr) != 0 {
			utils.WriteLog("log_tool", "bytesErr len != 0")
			chs <- 0
			return
		}
		bytestr, err := ioutil.ReadAll(stdout)
		if err != nil {
			utils.WriteLog("log_tool", "ioutil readall err，err:", err)
			chs <- 0
			return
		}
		if err := cmd.Wait(); err != nil {
			utils.WriteLog("log_tool", "cmd wait err, err:", err)
			chs <- 0
			return
		}
		res = string(bytestr)
		chs <- 1
	}()
	result := <-chs
	if result == 0 {
		c.JSON(200, gin.H{"code": "0", "desc": "fail"})
		return
	}
	c.JSON(200, gin.H{"code": "1", "desc": "success", "content": res})
	return
}

func (ctrl ToolController) RestartSinglephp(c *gin.Context) {
	ipstr := c.PostForm("ip")
	cmdStr := fmt.Sprintf("sh %s/restart_single_php.sh %s", restartPhpPath, ipstr)
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		utils.WriteLog("log_tool", "RestartSinglephp cmd stdoutPipe err,err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "cmd stdoutPipe err"})
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		utils.WriteLog("log_tool", "RestartSinglephp cmd StderrPipe err,err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "cmd StderrPipe err"})
		return
	}
	//执行命令
	if err := cmd.Start(); err != nil {
		utils.WriteLog("log_tool", "RestartSinglephp cmd Start err,err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "cmd start err"})
		return
	}
	bytesErr, err := ioutil.ReadAll(stderr)
	if err != nil {
		utils.WriteLog("log_tool", "RestartSinglephp ioutil ReadAll stderr err,err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "ioutil ReadAll stderr err"})
		return
	}
	if len(bytesErr) != 0 {
		utils.WriteLog("log_tool", "RestartSinglephp bytesErr len != 0")
		c.JSON(200, gin.H{"code": "0", "desc": "bytesErr len not 0"})
		return
	}
	bytestr, err := ioutil.ReadAll(stdout)
	if err != nil {
		utils.WriteLog("log_tool", "RestartSinglephp ioutil readall stdout err，err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "ioutil ReadAll err"})
		return
	}
	if err := cmd.Wait(); err != nil {
		utils.WriteLog("log_tool", "RestartSinglephp cmd wait err, err:", err)
		c.JSON(200, gin.H{"code": "0", "desc": "cmd Wait err"})
		return
	}
	c.JSON(200, gin.H{"code": "1", "desc": "success", "content": string(bytestr)})
	return
}

func (ctrl ToolController) Navigation(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")
	menu, err := ctrl.getMenu(roleId.(int))
	if err != nil {
		utils.WriteLog("log_tool", "ToolController Navigation get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}

	cardList := `[
	{"name":"Git","url":"http://dev45.gitlab.miyabaobei.com"},
	{"name":"ApiDocs","url":"http://apidocs.miyabaobei.com"},
	{"name":"Splunk","url":"http://splunk111.miyabaobei.com"},
	{"name":"Ums","url":"http://ums.intra.miyabaobei.com"},
	{"name":"Trac","url":"http://trac.intra.miyabaobei.com"},
	{"name":"Worktile","url":"https://my.worktile.com"},
	{"name":"DeployTest","url":"http://deploytest.miyabaobei.com"},
	{"name":"Confluence","url":"http://wiki.mia.com"},
	{"name":"禅道","url":"http://pm.mia.com"}
	]`
	json, _ := simplejson.NewJson([]byte(cardList))
	arr := json.MustArray()
	c.HTML(200, "tool/navigation.html", gin.H{
		"username":   userName,
		"moduleName": "tool",
		"cardList":   arr,
		"ctrName":    "navigation",
		"menu":       menu,
	})
}

func (ctrl ToolController) getMenu(roleId int) (menus []forms.ParentMenu, err error) {
	var menuModel = new(models.MenuModel)
	return menuModel.GetMenu(roleId)
}
