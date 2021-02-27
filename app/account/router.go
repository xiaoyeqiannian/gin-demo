package account

import (
	_ "gin-server/app/account/docs"

	"gin-server/app/account/middleware"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRouter(router *gin.Engine) *gin.Engine {

	apiA := router.Group("/account")
	{
		apiA.POST("/login", Login)
		apiA.POST("/logout", middleware.JwtAuth(), Logout)
		apiA.POST("/regist", Regist)
		apiA.POST("/code/verify", CodeVerify)
		apiA.POST("/password/forget", PasswordForget)
		apiA.POST("/password/modify", middleware.JwtAuth(), AccountPasswordModify)
	}

	apiAPro := router.Group("/account").Use(middleware.JwtAuth()).Use(middleware.PermissionHandler())
	{
		apiAPro.POST("/modify", AccountModify)
		apiAPro.GET("/list", AccountList)
		apiAPro.POST("/del", AccountDel)
		apiAPro.POST("/sub/add", AccountSubAdd)

		apiAPro.POST("/menu", AccountMenu)

		apiAPro.POST("/role/modify", RoleModify)
		apiAPro.GET("/role/list", RoleList)
		apiAPro.POST("/role/del", RoleDel)
		
		apiAPro.POST("/group/modify", GroupModify)
		apiAPro.GET("/group/list", GroupList)
		apiAPro.POST("/group/del", GroupDel)
	}

	apiFile := router.Group("/file").Use(middleware.JwtAuth()).Use(middleware.PermissionHandler())
	{
		apiFile.POST("/upload", AppUploadFile)
	}

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
