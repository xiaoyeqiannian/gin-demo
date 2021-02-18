package proc

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CurrentUser(c *gin.Context) (gin.H, error) {
	if tmp, ok := c.Get("claims"); ok {
		if claims, ok := tmp.(gin.H); ok {
			return claims, nil
		}
	}
	return nil, errors.New("用户信息获取失败!")
}
