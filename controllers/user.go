package controllers

import (
	//"fmt"
	"mint-platform/platform/forms"
	"mint-platform/platform/models"
	"mint-platform/platform/utils"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	//"reflect"
)

type UserController struct{}

var UserModel = new(models.UserModel)

func (ctrl UserController) Signin(c *gin.Context) {
	var signinForm forms.SigninForm

	if c.BindJSON(&signinForm) != nil {
		c.JSON(406, gin.H{"code": "0", "form": signinForm})
		c.Abort()
		return
	}
	user, err := UserModel.Signin(signinForm)
	if (err == nil) && (user.Id > 0) {
		session := sessions.Default(c)
		session.Set("username", user.UserName)
		session.Set("userid", user.Id)
		session.Set("roleid", user.RoleId)
		session.Save()
		c.JSON(200, gin.H{"code": "1", "desc": "success"})
	} else {
		utils.WriteLog("log_user", err)
		c.JSON(200, gin.H{"code": "0", "desc": "fail"})
	}
	return
}

func (ctrl UserController) Login(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	if userName != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
	c.HTML(200, "login.html", gin.H{})
}

func (ctrl UserController) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusTemporaryRedirect, "/login")
	return
}
