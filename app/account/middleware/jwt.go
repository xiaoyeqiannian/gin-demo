package middleware

import (
	"fmt"
	"strings"
	"net/http"
	"go-server/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	JWT_KEY               = []byte("demodemodeomdemo")
	JWT_EXPIRE_TIME int64 = 3600 * 24
)

type PayLoad struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
	Avatar   string `json:"avatar"`
	RoleID   int    `json:"role_id"`
	GroupID  int    `json:"group_id"`
}

type CustomClaims struct {
	Identity PayLoad `json:"identity"`
	jwt.StandardClaims
}

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.Abort()
			c.JSON(http.StatusOK, utils.RespJson(utils.LOGINERR, "login error", nil))
			return
		}
		t := strings.Split(tokenString, " ")
		if len(t) < 2{
			c.Abort()
			c.JSON(http.StatusOK, utils.RespJson(utils.LOGINERR, "token error", nil))
			return
		}
		token, err := jwt.ParseWithClaims(t[1], &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return JWT_KEY, nil
		})
		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					fmt.Println("That's not even a token")
				} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
					fmt.Println("Token is expired")
				} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
					fmt.Println("Token not active yet")
				} else {
					fmt.Println("Couldn't handle this token")
				}
				c.Abort()
				c.JSON(http.StatusOK, utils.RespJson(utils.LOGINERR, "login error", nil))
				return
			}
		}
		claims, ok := token.Claims.(*CustomClaims)
		if !ok || !token.Valid {
			c.Abort()
			fmt.Println("token is invalid")
			c.JSON(http.StatusOK, utils.RespJson(utils.LOGINERR, "login error", nil))
			return
		}
		
		c.Set("claims", gin.H{
						"user_name":   claims.Identity.UserName,
						"user_id":     claims.Identity.UserID,
						"avatar":      claims.Identity.Avatar,
						"role_id":     claims.Identity.RoleID,
						"group_id":    claims.Identity.GroupID,
						"current_url": c.Request.URL.Path})
		c.Next()
	}
}
