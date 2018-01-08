package utils

import (
	//"fmt"
	"log"
	"os"
	"time"
)

/**
 * WriteLog
 * 日志记录函数
 * @param string logConfig 日志配置
 * @param string logNote 日志内容
 */
func WriteLog(logConfigName string, logNote ...interface{}) {
	config, err := LoadConfig("log", logConfigName)
	if err != nil {
		return
	}
	//组装日志文件名（基础名+日期）
	logName := config["logpath"] + config["filename"] + "-" + time.Now().Format("2006-01-02") + ".log"
	//初始化日志对象
	var Filemange *os.File
	//判断文件是否存在
	if CheckFileIsExist(logName) == true {
		//打开文件
		Filemange, err = os.OpenFile(logName, os.O_RDWR|os.O_APPEND, 0777)
	} else {
		//不存在创建文件
		Filemange, err = os.Create(logName)
	}
	//判断是否获取到文件句柄
	if Filemange == nil {
		return
	}
	//关闭日志
	defer Filemange.Close()

	//初始化日志对象
	logger := log.New(Filemange, "\r\n", log.Ldate|log.Ltime)
	//设置日志级别
	logger.SetPrefix("[loglevel:" + config["loglevel"] + "] ")
	//记录日志
	logger.Println(logNote...)
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
