package router

import (
	"backend/control"
	"backend/logic"

	"github.com/gin-gonic/gin"
)

func RouterInit(r *gin.RouterGroup) {

	r.Static("/logo", "./image")
	api := r.Group("api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", control.UserLogin)
			auth.POST("/signup", control.SignUp)
		}

		api.Use(logic.AuthMiddleware()) //应该只影响后面的，如果前面的也受影响，可能是gin版本不同
		user := api.Group("/user")
		{
			user.PUT("/password", control.UpdatePassword)
		}
		message := api.Group("/message")
		{
			message.DELETE("", control.DeleteMessage)
			message.GET("", control.GetMessage)
			message.GET("/pages", control.GetMessagePageNum)
			message.POST("/search", control.GetMessageSearch)
			message.POST("/search/pages", control.GetMessagePageNumSearch)
		}
		api.POST("/applyauth", control.ApplyForTableAuth)
		project := api.Group("/project")
		{
			project.GET("", control.GetProjectBrief)
			project.GET("/pages", control.GetProjectPageNum)
			project.POST("/search", control.GetProjectBriefSearch)
			project.POST("/search/pages", control.GetProjectPageNumSearch)
			project.POST("/newprojecttable", control.NewProjectTable)
			project.GET("/getprojecttable", control.GetProjectTable)
			project.GET("/getallprojecttable", control.GetAllProjectTable)
			project.DELETE("/deleteprojecttable", control.DeleteProjectTable)
		}
		api.POST("/projectdetail", control.GetProjectDetail)
		apiinfo := api.Group("/apiinfo")
		{
			apiinfo.GET("", control.GetApiBrief)
			apiinfo.GET("/pages", control.GetApiPageNum)
			apiinfo.POST("/search/pages", control.GetApiPageNumSearch)
			apiinfo.POST("/search", control.GetApiSearch)
			apiinfo.GET("/details", control.GetApiDetail)
			apiinfo.POST("", control.SaveApi)
		}
		developer := api.Group("/developer/project")
		{
			developer.POST("/img", control.UploadLogo)
			developer.POST("/intro", control.UpdateProjectNameDesc)
			developer.POST("/member", control.UpdateProjectMember)
			developer.DELETE("/member", control.DeleteMember)
		}
		admin := api.Group("/admin")
		{
			admin.GET("/users", control.GetAllUser)
			admin.POST("/publish", control.SaveNotification)
			admin.DELETE("/apiinfo", control.DeleteApi)
			admin.GET("/users/details", control.GetAllUserDetail)
			admin.POST("/invitecode", control.GenInviteCode)
			admin.GET("/tables", control.GetTable)
			admin.POST("/auth", control.UpdatePermission)
			admin.POST("/auth/search", control.GetPermission)
			admin.GET("/requests/pages", control.GetPermissionRequestNum)
			admin.GET("/requests", control.GetPermissionRequest)
			admin.POST("/requests/search/pages", control.GetPermissionRequestNumSearch)
			admin.POST("/requests/search", control.GetPermissionRequestSearch)
			admin.POST("/requests/reject", control.RejectPermissionRequest)
			admin.POST("/requests/approve", control.ApprovePermissionRequest)
		}
	}
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	// 添加CORS中间件
	// config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://localhost:2020", "http://localhost:8080", "http://localhost:8081", "http://localhost:8082", "http://localhost:8083",
	// 	"http://localhost:8084", "http://localhost:8085"} // 允许访问的域名
	// config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"} // 允许的HTTP方法
	// router.Use(cors.New(config))
	api := router.Group("")
	RouterInit(api)
	// UserRouterInit(api)
	//NewsRouterInit(api)
	//CommentRouterInit(api)
	return router
}
