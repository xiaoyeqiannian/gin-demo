package account

import (
	"encoding/base64"
	"net/http"

	// "fmt"
	"time"
	"gin-server/app/account/middleware"
	. "gin-server/app/account/model"
	. "gin-server/database"
	"gin-server/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)


type LoginParam struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}


// @Tags login
// @Summary login
// @version 1.0
// @Description login
// @Accept  application/json
// @Produce json
// @Param username body string true "登陆账号"
// @Param password body string true "密码"
// @Success 0000 {string} json "{"code":"0000","message":"ok","data":{"token":"xxx.xx.xx"}}"
// @Failure 2101 {string} json "{"code":"2101","message":"name or password error","data":null}"
// @Router /account/login [post]
func Login(c *gin.Context) {
	var p LoginParam
	if err := c.Bind(&p); err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, "Params error!", nil))
		return
	}
	var u User
	if MysqlDB.First(&u, "name=?", p.Name); u.ID<= 0 {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, "Can not find this user, Please regist", nil))
		return
	}
	pwd, _ := base64.StdEncoding.DecodeString(p.Password)
	if isOK := utils.CheckPasswordHash(u.Password, string(pwd)); isOK != true {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, "Name or password error", nil))
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, middleware.CustomClaims{
		middleware.PayLoad{ UserID:u.ID, UserName:u.Name, Avatar:u.Avatar, RoleID:u.RoleID, GroupID:u.GroupID },
		jwt.StandardClaims{ ExpiresAt: int64(time.Now().Unix() + middleware.JWT_EXPIRE_TIME), Issuer:"gin-server" },
	})
	if tokenString, err := token.SignedString([]byte(middleware.JWT_KEY)); err == nil {
		c.JSON( http.StatusOK, utils.RespJson(utils.REQOK, "ok", gin.H{ "token": tokenString }) )
		return
	}
	c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, "login error, try again!", nil))
	return
}


func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", nil))
	return
}


// @Tags regist
// @Summary regist
// @version 1.0
// @Description regist
// @Accept  application/json
// @Produce json
// @Param username body string true "登陆账号"
// @Param password body string true "密码"
// @Success 0000 {string} json "{"code":"0000","message":"ok","data":{"token":"xxx.xx.xx"}}"
// @Failure 2101 {string} json "{"code":"2101","message":"name or password error","data":null}"
// @Router /account/regist [post]
func Regist(c *gin.Context) {
	var p LoginParam
	if err := c.Bind(&p); err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, "Params error!", nil))
		return
	}
	var u User
	MysqlDB.First(&u, "name=?", p.Name)
	if u.ID > 0 {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, "Registed, Please login", nil))
		return
	}
	var g Group
	if _, err := g.Add(p.Name, GROUP_KIND_PERSONAL); err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, err.Error(), nil))
		return
	}
	manager := User{Name: p.Name, Password: p.Password, RoleID: ROLE_ADMIN_ID, GroupID: g.ID}
	if MysqlDB.Create(&manager); manager.ID == 0{
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, "Insert data error!", nil))
		return
	}
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", nil))
	return
}
