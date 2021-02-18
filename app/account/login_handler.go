package account

import (
	"encoding/base64"
	"net/http"

	// "log"
	"fmt"
	"time"
	"go-server/app/account/middleware"
	. "go-server/app/account/model"
	. "go-server/database"
	"go-server/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)


type LoginParam struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}


func Login(c *gin.Context) {
	if c.Request.Method == "POST" {
		var params LoginParam
		if err := c.Bind(&params); err != nil {
			fmt.Println("Login参数错误")
			c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, "参数错误", nil))
			return
		}
		var u User
		MysqlDB.First(&u, "name=?", params.Name)
		if u.ID <= 0 {
			c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, "未找到此用户，请先注册", nil))
			return
		}
		pwd, _ := base64.StdEncoding.DecodeString(params.Password)
		fmt.Println("pwd", pwd)
		if isOK := utils.CheckPasswordHash(u.Password, string(pwd)); isOK != true {
			c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, "账号或密码错误", nil))
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, middleware.CustomClaims{
			middleware.PayLoad{UserID:u.ID, UserName:u.Name, Avatar:u.Avatar, RoleID:u.RoleID, GroupID:u.GroupID},
			jwt.StandardClaims{ExpiresAt: int64(time.Now().Unix() + middleware.JWT_EXPIRE_TIME), Issuer:"go-server"},
		})
		tokenString, err := token.SignedString([]byte(middleware.JWT_KEY))
		if err != nil {
			c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, "登录失败，请重试!", nil))
			return
		}
		fmt.Println("login ok ", tokenString)
		c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", gin.H{"token": tokenString}))
		return
	} else {
		c.HTML(http.StatusOK, "login.html", nil)
		return
	}
}


func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", nil))
	return
}


func Regist(c *gin.Context) {
	var params LoginParam
	if err := c.Bind(&params); err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, "参数错误", nil))
		return
	}
	var u User
	MysqlDB.First(&u, "name=?", params.Name)
	if u.ID > 0 {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, "此账号已注册，请直接登录", nil))
		return
	}
	pwd, err := utils.GeneratePasswordHash(params.Password)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, "注册失败", nil))
		return
	}
	var g Group
	if err := g.Merge(0, params.Name, GROUP_KIND_PERSONAL, 1); err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, err.Error(), nil))
		return
	}
	manager := User{Name: params.Name,
		            Password: pwd,
					RoleID: ROLE_ADMIN_ID,
					GroupID: g.ID}
	MysqlDB.Create(&manager)
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", nil))
	return
}
