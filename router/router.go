package router

import (
	"backend/control"
	"backend/logic"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 用户路由组
func UserRouterInit(r *gin.RouterGroup) {
	userManager := r.Group("/user")
	{
		// 其它的 api 删完了，只放一个示例
		userManager.POST("/register", control.RegisterHandler)
		userManager.Use(logic.AuthMiddleware())
	}
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	// 添加CORS中间件
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8080", "http://localhost:8081", "http://localhost:8082", "http://localhost:8083",
		"http://localhost:8084", "http://localhost:8085"} // 允许访问的域名
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"} // 允许的HTTP方法
	router.Use(cors.New(config))
	api := router.Group("")
	UserRouterInit(api)
	//NewsRouterInit(api)
	//CommentRouterInit(api)
	return router
}
