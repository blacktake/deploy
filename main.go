package main

import (
	"encoding/json"
	"fmt"
	"mint-platform/platform/controllers"
	"mint-platform/platform/db"
	"mint-platform/platform/forms"
	"mint-platform/platform/utils"
	"net/http"
	"runtime"
	"strings"

	"github.com/braintree/manners"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	ConfigRuntime()
	StartGin()
}

func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func CheckLoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userName := session.Get("username")
		userid := session.Get("userid")
		roleId := session.Get("roleid")
		if userName == nil {
			if strings.Index(c.Request.URL.Path, "ajax") != -1 {
				c.JSON(200, gin.H{"code": "3", "desc": "login"})
			} else {
				c.Redirect(http.StatusTemporaryRedirect, "/login")
			}

			c.Abort()
		}
		res := CheckUrlValidity(c.Request.URL.Path, roleId.(int))
		if !res {
			c.HTML(200, "no_permission.html", gin.H{"content": "你没有权限"})
			c.Abort()
		}
		// 请求前
		c.Next() //处理请求
		// 请求后
		var datajson []byte
		if strings.Index(c.Request.URL.Path, "ajax") != -1 {
			if len(c.Request.Form) > 0 {
				datajson, _ = json.Marshal(c.Request.Form)
			}
			SaveUserBehavior(c.Request.URL.Path, userid.(int), userName.(string), string(datajson))
		}
	}

}

//检查url有效性
func CheckUrlValidity(url string, roleId int) bool {
	engine := db.GetDB("deploy_online")
	var menu forms.Menu
	has, err := engine.Where("url = ?", url).Get(&menu)
	if err != nil {
		utils.WriteLog("log_preloading", "get menu url err，err:", err)
		return false
	}
	if has {
		var rolemenu forms.Admin_role_menu_relation
		has, err := engine.Where("roleid=? and menuid=?", roleId, menu.Id).Get(&rolemenu)
		if err != nil {
			utils.WriteLog("log_preloading", "get rolemenu relation err，err:", err)
			return false
		}
		if has {
			return true
		}
		return false
	}
	return true
}

//记录用户对数据操作
func SaveUserBehavior(url string, uid int, username string, datajson string) {
	engine := db.GetDB("deploy_online")
	userbehavior := new(forms.UserBehavior)
	userbehavior.Uid = uid
	userbehavior.Username = username
	userbehavior.Operate = "用户操作了:" + url
	userbehavior.Datajson = datajson
	_, err := engine.Insert(userbehavior)
	if err != nil {
		utils.WriteLog("log_user", "记录用户行为err，err:", err)
	}
}

func StartGin() {
	gin.SetMode(gin.ReleaseMode)
	//gin.SetMode(gin.DebugMode)
	//router := gin.New()
	router := gin.Default()

	store := sessions.NewCookieStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	//store, _ := sessions.NewRedisStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	//router.Use(sessions.Sessions("mysession", store))
	//db.Init("deploy_online")
	db.NewEngineMap()
	db.Init("deploy_online")
	db.Init("deploy_apibi")
	db.Init("deploy_log")
	//router.Use(gin.Logger())
	//router.Use(gin.Recovery())
	router.Use(CORSMiddleware())
	//router.Use(RateLimit, gin.Recovery())

	router.LoadHTMLGlob("templates/**/*")
	router.Static("/static", "static")

	router.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{})
	})

	index := new(controllers.IndexController)
	router.GET("/", index.IndexHandler)

	user := new(controllers.UserController)
	router.POST("/siginajax", user.Signin)
	router.GET("/login", user.Login)
	router.GET("/logout", user.Logout)

	dashboardGroup := router.Group("/dashboard")
	dashboardGroup.Use(CheckLoginMiddleware())
	{
		statistics := new(controllers.StatisticsController)
		dashboardGroup.GET("/showdashboard", statistics.ShowDashboard)
	}

	deployGroup := router.Group("/deploy")
	deployGroup.Use(CheckLoginMiddleware())
	{
		project := new(controllers.ProjectController)
		deployGroup.GET("/project", project.ProjectHandler)
		deployGroup.GET("/addworker", project.ShowAddWorker)
		deployGroup.POST("/addworkerajax", project.AddWorker)
		deployGroup.POST("/cancelworkerajax", project.CancelWorker)
		deployGroup.POST("/accessworkerajax", project.AccessWorker)
		deployGroup.GET("/publish", project.ShowPublish)
		deployGroup.POST("/excuteajax", project.Publish)
		deployGroup.POST("/showinfoajax", project.ShowInfo)

		rollback := new(controllers.RollbackController)
		deployGroup.GET("/rollbacklist", rollback.RollbackList)
		deployGroup.GET("/addrollbackorder", rollback.ShowAddWorker)
		deployGroup.POST("/showversion", rollback.ShowRollbackVersion)
		deployGroup.POST("/addrollbackorderajax", rollback.AddWorker)
		deployGroup.POST("/cancelrollbackajax", rollback.CancelWorker)
		deployGroup.GET("/rollback", rollback.ShowRollback)
		deployGroup.POST("/excuterollbackajax", rollback.Rollback)
		deployGroup.POST("/showrollbackajax", rollback.ShowInfo)
	}

	plusdeployGroup := router.Group("/plusdeploy")
	plusdeployGroup.Use(CheckLoginMiddleware())
	{
		plusproject := new(controllers.PlusProjectController)
		deployGroup.GET("/plusproject", plusproject.ProjectHandler)
		deployGroup.GET("/plusaddworker", plusproject.ShowAddWorker)
		deployGroup.POST("/plusaddworkerajax", plusproject.AddWorker)
		deployGroup.POST("/pluscancelworkerajax", plusproject.CancelWorker)
		deployGroup.POST("/plusaccessworkerajax", plusproject.AccessWorker)
		deployGroup.GET("/pluspublish", plusproject.ShowPublish)
		deployGroup.POST("/plusexcuteajax", plusproject.Publish)
		deployGroup.POST("/plusshowinfoajax", plusproject.ShowInfo)

		plusrollback := new(controllers.PlusRollbackController)
		deployGroup.GET("/plusrollbacklist", plusrollback.RollbackList)
		deployGroup.GET("/plusaddrollbackorder", plusrollback.ShowAddWorker)
		deployGroup.POST("/plusshowversion", plusrollback.ShowRollbackVersion)
		deployGroup.POST("/plusaddrollbackorderajax", plusrollback.AddWorker)
		deployGroup.POST("/pluscancelrollbackajax", plusrollback.CancelWorker)
		deployGroup.GET("/plusrollback", plusrollback.ShowRollback)
		deployGroup.POST("/plusexcuterollbackajax", plusrollback.Rollback)
		deployGroup.POST("/plusshowrollbackajax", plusrollback.ShowInfo)
	}

	logGroup := router.Group("/log")
	logGroup.Use(CheckLoginMiddleware())
	{
		logctr := new(controllers.LogController)
		logGroup.GET("/loglist", logctr.LogList)
		logGroup.GET("/loginlist", logctr.LoginList)
	}

	toolGroup := router.Group("/tool")
	toolGroup.Use(CheckLoginMiddleware())
	{
		toolctr := new(controllers.ToolController)
		toolGroup.GET("/navigation", toolctr.Navigation)
		toolGroup.GET("/delvarnish", toolctr.Delvarnish)
		toolGroup.POST("/delvarnishajax", toolctr.ExcuteDelvarnish)
		toolGroup.GET("/showrestartphp", toolctr.ShowRestartphp)
		toolGroup.POST("/restartphpajax", toolctr.Restartphp)
		toolGroup.POST("/restartsinglephpajax", toolctr.RestartSinglephp)
	}

	monitorGroup := router.Group("/monitor")
	monitorGroup.Use(CheckLoginMiddleware())
	{
		monitorctr := new(controllers.MonitorController)
		monitorGroup.GET("/mq", monitorctr.Mqlist)
		monitorGroup.GET("/php", monitorctr.Showphp)
	}

	menuGroup := router.Group("/menu")
	menuGroup.Use(CheckLoginMiddleware())
	{
		menu := new(controllers.MenuController)
		menuGroup.GET("/list", menu.MenuList)
		menuGroup.GET("/edit", menu.ShowEditMenu)
		menuGroup.POST("/addmenuajax", menu.AddMenu)
		menuGroup.POST("/changemenuajax", menu.ChangeMenuOrder)
		menuGroup.POST("/deletemenuajax", menu.DeleteMenu)
		menuGroup.POST("/editmenuajax", menu.EditMenu)
	}

	permissionGroup := router.Group("/permission")
	permissionGroup.Use(CheckLoginMiddleware())
	{
		permission := new(controllers.PermissionController)
		permissionGroup.GET("/rolelist", permission.RoleList)
		permissionGroup.GET("/permissionlist", permission.PermissionList)
		permissionGroup.GET("/userlist", permission.Userlist)
		permissionGroup.POST("/changeuserajax", permission.ChangeUser)
		permissionGroup.POST("/deleteuserajax", permission.DeleteUser)
		permissionGroup.GET("/adduser", permission.ShowAddUser)
		permissionGroup.GET("/edituser", permission.ShowEditUser)
		permissionGroup.POST("/adduserajax", permission.AddUser)
		permissionGroup.POST("/edituserajax", permission.EditUser)
		permissionGroup.GET("/addrole", permission.ShowAddRole)
		permissionGroup.POST("/addroleajax", permission.AddRole)
		permissionGroup.GET("/editpermission", permission.ShowEditPermission)
		permissionGroup.POST("/updatepermissionajax", permission.UpdatePermission)
	}

	manners.ListenAndServe(":8077", router)
	//router.Run(":8002")
}
