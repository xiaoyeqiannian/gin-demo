package middleware

import (
	"fmt"
	"strings"
	"net/http"
	"gin-server/app/account/proc"
	"gin-server/app/account/model"
	"gin-server/utils"

	"github.com/gin-gonic/gin"
)

func PermissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var currentUser gin.H
		if tmp, ok := c.Get("claims"); ok {
			if claims, ok := tmp.(gin.H); ok {
				currentUser = claims
			}
		}
		// 获取请求的URI
		obj := c.Request.URL.RequestURI()
		// 获取请求方法
		act := strings.ToLower(c.Request.Method)
		// 获取用户的角色
		roleID := currentUser["role_id"].(int)
		var sub = ""
		if roleID == model.ROLE_ROOT_ID{
			sub = "root"
		}else{
			sub = fmt.Sprint(roleID)
		}
		e := proc.Casbin()
		fmt.Println(sub, obj, act)
		// 判断策略中是否存在
		if e.Enforce(fmt.Sprint(sub), obj, act) {
			c.Next()
		} else {
			c.JSON(http.StatusOK, utils.RespJson(utils.ROLEERR, "权限不足", nil))
			c.Abort()
			return
		}
	}
}
