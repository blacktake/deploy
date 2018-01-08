package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mint-platform/platform/forms"
	"mint-platform/platform/models"
	"mint-platform/platform/utils"
	"os/exec"
	"strconv"
	"strings"
)

type MonitorController struct{}

type RabbitmqInfo struct {
	Messages  float64
	Consumers float64
	MqIp      string
}

type MqInfo struct {
	Mqinfos []RabbitmqInfo
	MqName  string
}

type PhpProcess struct {
	Ip     string
	Count  int
	Danger bool
}

var checkPhpPath = "/root/tool/script"

//var checkPhpPath = "/data/bin//script"

func (ctrl MonitorController) Mqlist(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")

	var mqinfo map[string]string
	var queueslist map[string]string

	mqinfo = make(map[string]string)
	queueslist = make(map[string]string)
	queueslist["coupon_create_bind_coupon_que"] = "10.1.50.215,10.1.50.216"
	queueslist["redbag_bind_redbag_que"] = "10.1.50.215,10.1.50.216"
	queueslist["index_activation_que"] = "10.1.52.76,10.1.52.77,10.1.52.78"
	queueslist["amqp_write_user_id_que"] = "10.1.52.76,10.1.52.77,10.1.52.78"
	queueslist["save_device_info_que"] = "10.1.52.76,10.1.52.77,10.1.52.78"

	mqinfo["10.1.50.215"] = "miya_amqp_admin:miya_admin_pwd"
	mqinfo["10.1.50.216"] = "miya_amqp_admin:miya_admin_pwd"
	mqinfo["10.1.52.76"] = "miya_amqp_admin:miya_admin_pwd"
	mqinfo["10.1.52.77"] = "miya_amqp_admin:miya_admin_pwd"
	mqinfo["10.1.52.78"] = "miya_amqp_admin:miya_admin_pwd"

	couponMqName := "coupon_create_bind_coupon_que"
	couponMq, err := ctrl.getMqInfos(queueslist[couponMqName], couponMqName, mqinfo, queueslist)
	if err != nil {
		utils.WriteLog("log_monitor", "MonitorController MqList coupon mq getMqInfos err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	couponMq.MqName = couponMqName

	redbagBindMqName := "redbag_bind_redbag_que"
	redbagMq, err := ctrl.getMqInfos(queueslist[redbagBindMqName], redbagBindMqName, mqinfo, queueslist)
	if err != nil {
		utils.WriteLog("log_monitor", "MonitorController MqList redbag mq getMqInfos err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	redbagMq.MqName = redbagBindMqName

	indexActivationMqName := "index_activation_que"
	indexActivationMq, err := ctrl.getMqInfos(queueslist[indexActivationMqName], indexActivationMqName, mqinfo, queueslist)
	if err != nil {
		utils.WriteLog("log_monitor", "MonitorController MqList index activation mq getMqInfos err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	indexActivationMq.MqName = indexActivationMqName

	amqpWriteUserIdMqName := "amqp_write_user_id_que"
	amqpWriteUserIdMq, err := ctrl.getMqInfos(queueslist[amqpWriteUserIdMqName], amqpWriteUserIdMqName, mqinfo, queueslist)
	if err != nil {
		utils.WriteLog("log_monitor", "MonitorController MqList amqp_write_user_id mq getMqInfos err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	amqpWriteUserIdMq.MqName = amqpWriteUserIdMqName

	saveDeviceInfoMqName := "save_device_info_que"
	saveDeviceInfoMq, err := ctrl.getMqInfos(queueslist[saveDeviceInfoMqName], saveDeviceInfoMqName, mqinfo, queueslist)
	if err != nil {
		utils.WriteLog("log_monitor", "MonitorController MqList save_device_info_que getMqInfos err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	saveDeviceInfoMq.MqName = saveDeviceInfoMqName

	menu, err := ctrl.getMenu(roleId.(int))
	if err != nil {
		utils.WriteLog("log_monitor", "MonitorController MqList get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	c.HTML(200, "monitor/mqlist.html", gin.H{
		"username":          userName,
		"moduleName":        "monitor",
		"ctrName":           "mq",
		"ctrNameZn":         "mq列表",
		"couponMq":          couponMq,
		"redbagMq":          redbagMq,
		"indexActivationMq": indexActivationMq,
		"amqpWriteUserIdMq": amqpWriteUserIdMq,
		"saveDeviceInfoMq":  saveDeviceInfoMq,
		"menu":              menu,
	})
}

func (ctrl MonitorController) getMqInfos(ips string, mqname string, mqinfo map[string]string, queueslist map[string]string) (mq MqInfo, err error) {
	if len(queueslist[mqname]) == 0 {
		return mq, errors.New("queueslist is empty_pagepty！")
	}
	iplists := strings.Split(ips, ",")
	for _, ip := range iplists {
		ip = strings.Replace(ip, " ", "", -1)
		urlstr := fmt.Sprintf("http://%s@%s:15672/api/queues/%s/%s", mqinfo[ip], ip, "%2f", mqname)
		result, err := utils.GetUrl(urlstr)
		if err != nil {
			utils.WriteLog("log_monitor", "MonitorController MqList utils GetUrl err, err:", err)
			return mq, err
		}
		resinfo, err := ctrl.getRabbitmqInfo(result, ip)
		if err != nil {
			utils.WriteLog("log_monitor", "MonitorController MqList ctrl getRabbitmqInfo err, err:", err)
			return mq, err
		}
		mq.Mqinfos = append(mq.Mqinfos, resinfo)
	}
	return mq, nil
}

func (ctrl MonitorController) getRabbitmqInfo(result []byte, ip string) (mqinfo RabbitmqInfo, err error) {
	var f interface{}
	err = json.Unmarshal(result, &f)
	if err != nil {
		return mqinfo, err
	}
	m := f.(map[string]interface{})
	for k, v := range m {
		if k == "messages" {
			mqinfo.Messages = v.(float64)
		}
		if k == "consumers" {
			mqinfo.Consumers = v.(float64)
		}
	}
	mqinfo.MqIp = ip
	return mqinfo, nil
}

func (ctrl MonitorController) Showphp(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")
	menu, err := ctrl.getMenu(roleId.(int))
	if err != nil {
		utils.WriteLog("log_monitor", "MonitorController Showphp get menu err, err:", err)
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	typestr := c.DefaultQuery("type", "api")
	var fileName = ""
	if typestr == "api" {
		fileName = "check_api_php.sh"
	} else if typestr == "cart" {
		fileName = "check_cart_php.sh"
	} else if typestr == "order" {
		fileName = "check_order_php.sh"
	} else if typestr == "recommend" {
		fileName = "check_recommend_php.sh"
	} else if typestr == "thirdservice" {
		fileName = "check_third_php.sh"
	}
	cmdStr := fmt.Sprintf("sh %s/%s", checkPhpPath, fileName)
	chs := make(chan int)
	var phpprolist []PhpProcess
	go func() {
		cmd := exec.Command("/bin/sh", "-c", cmdStr)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			utils.WriteLog("log_monitor", "cmd stdoutPipe err,err:", err, " type:", typestr)
			chs <- 0
			return
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			utils.WriteLog("log_monitor", "cmd StderrPipe err,err:", err, " type:", typestr)
			chs <- 0
			return
		}
		//执行命令
		if err := cmd.Start(); err != nil {
			utils.WriteLog("log_monitor", "cmd Start err,err:", err, " type:", typestr)
			chs <- 0
			return
		}
		bytesErr, err := ioutil.ReadAll(stderr)
		if err != nil {
			utils.WriteLog("log_monitor", "ioutil ReadAll stderr err,err:", err, " type:", typestr)
			chs <- 0
			return
		}
		if len(bytesErr) != 0 {
			utils.WriteLog("log_monitor", "bytesErr len != 0", " type", typestr)
			chs <- 0
			return
		}
		bytestr, err := ioutil.ReadAll(stdout)
		if err != nil {
			utils.WriteLog("log_monitor", "ioutil readall err，err:", err, " type:", typestr)
			chs <- 0
			return
		}
		if err := cmd.Wait(); err != nil {
			utils.WriteLog("log_monitor", "cmd wait err, err:", err, " type:", typestr)
			chs <- 0
			return
		}
		phpprocesslist := ctrl.getphpprocesslist(string(bytestr))
		phpprolist = phpprocesslist
		chs <- 1
	}()
	result := <-chs
	if result == 0 {
		c.HTML(200, "empty_page.html", gin.H{})
		return
	}
	c.HTML(200, "monitor/showphp.html", gin.H{
		"username":       userName,
		"moduleName":     "monitor",
		"ctrName":        "showphp",
		"ctrNameZn":      "展示php进程",
		"typestr":        typestr,
		"phpprocesslist": phpprolist,
		"menu":           menu,
	})
}

func (ctrl MonitorController) getphpprocesslist(phpstr string) (phplist []PhpProcess) {
	phpres := strings.Split(phpstr, "#")
	for _, v := range phpres {
		if len(v) > 1 {
			ipcount := strings.Split(v, ":")
			phpprocess := new(PhpProcess)
			phpprocess.Ip = strings.Replace(ipcount[0], "\n", "", -1)
			countnum := strings.Replace(ipcount[1], "\n", "", -1)
			phpcount, _ := strconv.Atoi(countnum)
			phpprocess.Count = phpcount
			if phpcount > 200 {
				phpprocess.Danger = true
			} else {
				phpprocess.Danger = false
			}
			phplist = append(phplist, *phpprocess)
		}
	}
	return phplist
}

func (ctrl MonitorController) getMenu(roleId int) (menus []forms.ParentMenu, err error) {
	var menuModel = new(models.MenuModel)
	return menuModel.GetMenu(roleId)
}
