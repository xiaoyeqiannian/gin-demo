package account

import (
	"net/http"
	"strconv"
	// "fmt"
	"encoding/base64"
	"strings"
	"github.com/gin-gonic/gin"

	. "gin-server/app/account/model"
	"gin-server/app/account/proc"
	. "gin-server/database"
	"gin-server/utils"
	. "gin-server/app/account/middleware"
	
)

type AccountParam struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	RoleID   int    `json:"role_id"`
	State    int8   `json:"state"`
	Avatar   string `json:"avatar"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
func AccountModify(c *gin.Context) {
	var currentUser PayLoad
	GetCurrentUser(c, &currentUser)
	var p AccountParam
	if _err := c.Bind(&p); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, _err.Error(), nil))
		return
	}
	var u User
	if p.ID > 0 {
		u.ID = p.ID
	} else {
		u.ID = currentUser.UserID
	}
	if _, _err := u.Update(p.State, p.RoleID, 0, p.Name, p.Password, p.Avatar, p.Email); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.REQERR, _err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", nil))
}
func AccountSubAdd(c *gin.Context) {
	var p AccountParam
	if _err := c.Bind(&p); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, _err.Error(), nil))
		return
	}
	var currentUser PayLoad
	GetCurrentUser(c, &currentUser)
	var u User
	if _, _err := u.Add(p.RoleID, currentUser.GroupID, p.Name, p.Password, p.Avatar, p.Email); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.REQERR, _err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", nil))
}


func AccountList(c *gin.Context) {
	var users []User
	searchKey := c.Query("searchKey")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	var total int
	sql := MysqlDB.Table("user")
	if searchKey != "" {
		sql = sql.Where("name LIKE ?", "%"+searchKey+"%").Count(&total)
	} else {
		sql = sql.Count(&total)
	}
	sql.Limit(pageSize).Offset((page - 1) * pageSize).Find(&users)
	type Resp struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Avatar    string `json:"avatar"`
		RoleID    int    `json:"role_id"`
		RoleName  string `json:"role_name"`
		UpdatedAt string `json:"updated_at"`
		CreatedAt string `json:"created_at"`
		State     int8   `json:"state"`
		StateCN   string `json:"state_cn"`
	}
	items := make([]Resp, len(users))
	for i, v := range users {
		var roleName = "-"
		if v.RoleID > 0 {
			var role = Role{ID: v.RoleID}
			MysqlDB.First(&role)
			if role.ID > 0{
				roleName = role.Name
			}
		}
		r := Resp{
			ID:        v.ID,
			Name:      v.Name,
			Avatar:    v.Avatar,
			RoleName:  roleName,
			RoleID:    v.RoleID,
			UpdatedAt: v.UpdatedAt.Format("2006-01-02 15:04:05"),
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
			State:     v.State,
			StateCN:   v.ParseState(),
		}
		items[i] = r
	}
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", gin.H{"items": items, "total": total, "page": page, "per_page": pageSize}))
}


func AccountPasswordModify(c *gin.Context) {
	type Param struct {
		ID       int    `json:"id"`
		Password string `json:"password"`
		NewPassword string `json:"new_password"`
	}
	var p Param
	if _err := c.Bind(&p); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, _err.Error(), nil))
		return
	}
	var currentUser PayLoad
	GetCurrentUser(c, &currentUser)
	var u User
	if MysqlDB.First(&u, currentUser.UserID); u.ID == 0 {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, "Can not find this user, Please regist", nil))
		return
	}
	pwd, _ := base64.StdEncoding.DecodeString(p.Password)
	if isOK := utils.CheckPasswordHash(u.Password, string(pwd)); isOK != true {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, "Name or password error", nil))
		return
	}
	if currentUser.GroupID==GROUP_SYS_ADMIN_ID || currentUser.RoleID==ROLE_ADMIN_ID{
		if p.ID > 0{
			u.ID = p.ID
		}
	}
	if _err := u.PasswordUpdate(p.NewPassword); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.REQERR, _err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", nil))
}


func AccountDel(c *gin.Context) {
	type Param struct {
		IDs       []int    `json:"ids" binding:"required"`
	}
	var p Param
	if _err := c.Bind(&p); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, _err.Error(), nil))
		return
	}
	var currentUser PayLoad
	GetCurrentUser(c, &currentUser)
	var u User
	if _err := u.Del(p.IDs, currentUser.GroupID); _err!=nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, _err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", nil))
}


func AccountMenu(c *gin.Context){
	var currentUser PayLoad
	GetCurrentUser(c, &currentUser)
	var u User
	if MysqlDB.First(&u, currentUser.UserID); u.ID == 0 {
		c.JSON(http.StatusOK, utils.RespJson(utils.REQERR, "Can not find user!", nil))
		return
	}
	var r Role
	if MysqlDB.First(&r, u.RoleID); r.ID == 0 {
		c.JSON(http.StatusOK, utils.RespJson(utils.REQERR, "Role error!", nil))
		return
	}
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", gin.H{"menu": strings.Split(r.Menu, ",")}))
}


func GroupModify(c *gin.Context){
	type Param struct {
		ID       int    `json:"id"`
        Name     string `json:"name"`
		Kind     int8   `json:"kind"`
		State    int8   `json:"state"`
	}
	var currentUser PayLoad
	GetCurrentUser(c, &currentUser)
	var p Param
	if err := c.Bind(&p); err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, err.Error(), nil))
		return
	}
	var g Group
	if p.ID > 0 {
		g.ID = p.ID
	} else {
		g.ID = currentUser.GroupID
	}
	if _, _err := g.Update(p.Name, p.Kind, p.State); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.REQERR, _err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", nil))
}


func GroupList(c *gin.Context){
	var groups []Group
	searchKey := c.Query("searchKey")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	var total int
	sql := MysqlDB.Table("group")
	if searchKey != "" {
		sql = sql.Where("name LIKE ?", "%"+searchKey+"%").Count(&total)
	} else {
		sql = sql.Count(&total)
	}
	sql.Limit(pageSize).Offset((page - 1) * pageSize).Find(&groups)
	type Resp struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Kind      int8    `json:"kind"`
		UpdatedAt string `json:"updated_at"`
		CreatedAt string `json:"created_at"`
		State     int8   `json:"state"`
		StateCN   string `json:"state_cn"`
	}
	items := make([]Resp, len(groups))
	for i, v := range groups {
		r := Resp{
			ID:        v.ID,
			Name:      v.Name,
			Kind:      v.Kind,
			UpdatedAt: v.UpdatedAt.Format("2006-01-02 15:04:05"),
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
			State:     v.State,
			StateCN:   v.ParseState(),
		}
		items[i] = r
	}
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", gin.H{"items": items, "total": total, "page": page, "per_page": pageSize}))
}


func GroupDel(c *gin.Context){
	type Param struct {
		IDs       []int    `json:"ids" binding:"required"`
	}
	var p Param
	if _err := c.Bind(&p); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, _err.Error(), nil))
		return
	}
	var currentUser PayLoad
	GetCurrentUser(c, &currentUser)
	var g Group
	if _err := g.Del(p.IDs, currentUser.GroupID); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.REQERR, _err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", nil))
}


func RoleModify(c *gin.Context) {
	type Param struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Menu     string `json:"menu"`
		GroupID  int    `json:"group_id"`
		State    int8   `json:"state"`
		Permissions []proc.Permission `json:"permissions"`
	}
	var p Param
	if err := c.Bind(&p); err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, err.Error(), nil))
		return
	}
	var currentUser PayLoad
	GetCurrentUser(c, &currentUser)
	var r = Role{ID: p.ID}
	if r.ID > 0{
		MysqlDB.First(&r)
	}
	var groupID int
	if p.GroupID > 0 {
		groupID = p.GroupID
	} else {
		groupID = currentUser.GroupID
	}
	if r.ID > 0{
		if _, _err := r.Update(groupID, p.Name, p.Menu, p.State); _err != nil {
			c.JSON(http.StatusOK, utils.RespJson(utils.REQERR, _err.Error(), nil))
			return
		}
	} else {
		if _, _err := r.Add(groupID, p.Name, p.Menu); _err != nil {
			c.JSON(http.StatusOK, utils.RespJson(utils.REQERR, _err.Error(), nil))
			return
		}
	}
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", gin.H{"id": r.ID}))
}


func RoleList(c *gin.Context) {
	var roles []Role
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	var total int
	var currentUser PayLoad
	GetCurrentUser(c, &currentUser)
	groupID := currentUser.GroupID
	if groupID == GROUP_SYS_ADMIN_ID {
		groupID, _ = strconv.Atoi(c.DefaultQuery("group_id", "0"))
	}
	sql := MysqlDB.Table("role")
	if groupID > 0 {
		sql = sql.Where("group_id=? and state=1", groupID)
	}
	sql = sql.Count(&total)
	sql.Limit(pageSize).Offset((page - 1) * pageSize).Find(&roles)
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", gin.H{"items": roles, "total": total, "page": page, "per_page": pageSize}))
}


func RoleDel(c *gin.Context) {
	type Param struct {
		IDs       []int    `json:"ids"`
	}
	var p Param
	if _err := c.Bind(&p); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, _err.Error(), nil))
		return
	}
	var currentUser PayLoad
	GetCurrentUser(c, &currentUser)
	var r Role
	if _err := r.Del(p.IDs, currentUser.GroupID); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.REQERR, _err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", nil))
}


func PasswordForget(c *gin.Context) {
	type Param struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
	}
	var p Param
	if _err := c.Bind(&p); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, _err.Error(), nil))
		return
	}
	// TODO send email and create code
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", nil))
}


func CodeVerify(c *gin.Context) {
	type Param struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Code     string `json:"code"`
		Password string `json:"password"`
	}
	var p Param
	if _err := c.Bind(&p); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, _err.Error(), nil))
		return
	}
	var u User
	if MysqlDB.First(&u, "name=?", p.Name); u.ID==0 {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, "Registed, Please login", nil))
		return
	}
	if _err := u.PasswordUpdate(p.Password); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, _err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", nil))
}


// TODO
func ManagerRolePermission(c *gin.Context) {
	if c.Request.Method == "POST" {
		type Param struct {
			RoleID      string            `json:"role_id"`
			Permissions []proc.Permission `json:"permissions"`
		}
		var p Param
		if err := c.Bind(&p); err != nil {
			c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, err.Error(), nil))
			return
		}
		err := proc.UpdateCasbin(p.RoleID, p.Permissions)
		if err != nil {
			c.JSON(http.StatusOK, utils.RespJson(utils.REQERR, "操作权限失败", nil))
		} else {
			c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", gin.H{}))
		}
	} else {
		roleID := c.Query("role_id")
		permissions := proc.GetPolicyPathByRoleID(roleID)
		c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", gin.H{"permission": permissions}))
	}
}
