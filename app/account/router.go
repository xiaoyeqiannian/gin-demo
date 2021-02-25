package account

import (
	_ "gin-server/app/account/docs"

	"gin-server/app/account/middleware"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRouter(router *gin.Engine) *gin.Engine {

	router.POST("/account/login", Login)
	router.POST("/account/logout", middleware.JwtAuth(), Logout)
	router.POST("/account/regist", Regist)
	router.POST("/account/password/modify", middleware.JwtAuth(), AccountPasswordModify)

	apiAccount := router.Group("/account").Use(middleware.JwtAuth()).Use(middleware.PermissionHandler())
	{
		apiAccount.POST("/modify", AccountModify)
		apiAccount.GET("/list", AccountList)
		apiAccount.POST("/del", AccountDel)
		apiAccount.POST("/sub/add", AccountSubAdd)

		apiAccount.POST("/menu", AccountMenu)

		apiAccount.POST("/role/modify", RoleModify)
		apiAccount.GET("/role/list", RoleList)
		apiAccount.POST("/role/del", RoleDel)
		
		apiAccount.POST("/group/modify", GroupModify)
		apiAccount.GET("/group/list", GroupList)
		apiAccount.POST("/group/del", GroupDel)
	}

	apiFile := router.Group("/file").Use(middleware.JwtAuth()).Use(middleware.PermissionHandler())
	{
		apiFile.POST("/upload", AppUploadFile)
	}

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
