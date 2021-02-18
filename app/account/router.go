package account

import (
	_ "go-server/app/account/docs"

	"go-server/app/account/middleware"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRouter(router *gin.Engine) *gin.Engine {

	router.POST("/login", Login)
	router.POST("/logout", middleware.JwtAuth(), Logout)
	router.POST("/regist", Regist)

	apiAccount := router.Group("/account").Use(middleware.JwtAuth()).Use(middleware.PermissionHandler())
	{
		apiAccount.POST("/modify", AccountModify)
		apiAccount.GET("/list", AccountList)
		apiAccount.POST("/del", AccountDel)

		apiAccount.POST("/menu", AccountMenu)

		apiAccount.POST("/role/modify", RoleModify)
		apiAccount.GET("/role/list", RoleList)
		apiAccount.GET("/role/del", RoleDel)
		
		apiAccount.POST("/group/modify", GroupModify)
		apiAccount.GET("/group/list", GroupList)
		apiAccount.GET("/group/del", GroupDel)
	}

	apiFile := router.Group("/file").Use(middleware.JwtAuth()).Use(middleware.PermissionHandler())
	{
		apiFile.POST("/upload", AppUploadFile)
	}

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
