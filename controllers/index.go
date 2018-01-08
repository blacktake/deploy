package controllers

import (
	//"html"
	//"io"
	//"strings"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"mint-platform/platform/forms"
	"mint-platform/platform/models"
	"net/http"
)

type IndexController struct{}

func (ctrl IndexController) IndexHandler(c *gin.Context) {
	session := sessions.Default(c)
	userName := session.Get("username")
	roleId := session.Get("roleid")
	if userName == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
	menu, _ := ctrl.getMenu(roleId.(int))
	c.HTML(200, "index/index.html", gin.H{
		"username":   userName,
		"moduleName": "index",
		"ctrName":    "index",
		"menu":       menu,
		"ctrNameZn":  "发布列表",
	})
}

func (ctrl IndexController) getMenu(roleId int) (menus []forms.ParentMenu, err error) {
	var menuModel = new(models.MenuModel)
	return menuModel.GetMenu(roleId)
}
