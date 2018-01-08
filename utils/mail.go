package utils

/*
	mail类
*/
import (
	//"fmt"
	"net/smtp"
	"strings"
	"time"
)

const (
	user     = "notice@mia.com"
	password = "123qweASD"
	host     = "smtp.exmail.qq.com:25"
	//	to       = "yulei1@mia.com"
	subject = "上线通知"
)

/**
 * SendToMail
 * @param string to 收件人
 * @param string body 邮件内容
 * @param string mailtype 邮件内容类型
 * @return error 操作结果
 */
func SendToMail(to, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

func MakeTemplateToMail(user_name, functional_introduction, email_list, auditor, group_name string) error {
	str_functional_introduction := strings.Replace(functional_introduction, ";", "<br />", -1)
	html := "<html><meta charset='utf-8'>    <meta name='viewport' content='width=device-width, initial-scale=1.0'>    <meta http-equiv='X-UA-Compatible' content='IE=edge'><table width='93%' border='1' cellspacing='1' cellpadding='1'><thead><tr ><th colspan='4'>上线通知</th></tr>"
	html += "<tr><th>上线人</th><th>" + user_name + "</th> <th>上线时间</th><th>" + time.Now().Format("2006-01-02 15:04:05") + "</th></tr>"
	html += "<tr><th>上线审核人</th><th>" + auditor + "</th> <th>分组</th><th>" + group_name + "</th></tr>"
	html += "<tr><th>上线内容</th> <th colspan='3'><p>" + str_functional_introduction + "</p></th></tr>"
	html += "</thead><html>"
	err := SendToMail(email_list, html, "html")
	return err
}
