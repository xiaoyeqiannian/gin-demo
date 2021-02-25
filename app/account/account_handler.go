package account

import (
	"net/http"
	"strconv"
	"fmt"
	"encoding/base64"
	"strings"
	. "gin-server/app/account/model"
	"gin-server/app/account/proc"
	. "gin-server/database"
	"gin-server/utils"
	"github.com/gin-gonic/gin"
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
	var currentUser gin.H
	if tmp, ok := c.Get("claims"); ok {
		if claims, ok := tmp.(gin.H); ok {
			currentUser = claims
		}
	}
	var p AccountParam
	if _err := c.Bind(&p); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, _err.Error(), nil))
		return
	}
	var u User
	if p.ID > 0 {
		u.ID = p.ID
	} else {
		u.ID = currentUser["user_id"].(int)
	}
	if _, _err := u.Merge(p.State, p.RoleID, 0, p.Name, p.Password, p.Avatar, p.Email); _err != nil {
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
	var currentUser gin.H
	if tmp, ok := c.Get("claims"); ok {
		if claims, ok := tmp.(gin.H); ok {
			currentUser = claims
		}
	}
	var u User
	if _, _err := u.Merge(p.State, p.RoleID, currentUser["group_id"].(int), p.Name, p.Password, p.Avatar, p.Email); _err != nil {
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
	var currentUser gin.H
	if tmp, ok := c.Get("claims"); ok {
		if claims, ok := tmp.(gin.H); ok {
			currentUser = claims
		}
	}
	var u User
	if MysqlDB.First(&u, currentUser["user_id"].(int)); u.ID == 0 {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, "Can not find this user, Please regist", nil))
		return
	}
	pwd, _ := base64.StdEncoding.DecodeString(p.Password)
	if isOK := utils.CheckPasswordHash(u.Password, string(pwd)); isOK != true {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, "Name or password error", nil))
		return
	}
	tmp, err := base64.StdEncoding.DecodeString(p.NewPassword)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, "Parse password error!", nil))
		return
	}
	newPwd, err := utils.GeneratePasswordHash(string(tmp))
	if err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, "Generate password error!", nil))
		return
	}
	if currentUser["group_id"]==GROUP_SYS_ADMIN_ID || currentUser["role_id"].(int)==ROLE_ADMIN_ID{
		if p.ID > 0{
			u.ID = p.ID
		}
	}
	if _, _err := u.Merge(0, 0, 0, "", newPwd, "", ""); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.REQERR, _err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", nil))
}


func AccountDel(c *gin.Context) {
	type Param struct {
		ID       int    `json:"id" binding:"required"`
	}
	var p Param
	if _err := c.Bind(&p); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, _err.Error(), nil))
		return
	}
	u := User{ID: p.ID}
	if _, _err := u.Merge(2, 0, 0, "","", "", ""); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.REQERR, _err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", nil))
}


func AccountMenu(c *gin.Context){
	var currentUser gin.H
	if tmp, ok := c.Get("claims"); ok {
		if claims, ok := tmp.(gin.H); ok {
			currentUser = claims
		}
	}
	var u User
	if MysqlDB.First(&u, currentUser["user_id"].(int)); u.ID == 0 {
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
	var currentUser gin.H
	if tmp, ok := c.Get("claims"); ok {
		if claims, ok := tmp.(gin.H); ok {
			currentUser = claims
		}
	}
	var p Param
	if err := c.Bind(&p); err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, err.Error(), nil))
		return
	}
	var g Group
	if p.ID > 0 {
		g.ID = p.ID
	} else {
		g.ID = currentUser["group_id"].(int)
	}
	if _, _err := g.Merge(p.Name, p.Kind, p.State); _err != nil {
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
		fmt.Println("searchKey:", searchKey)
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
		ID       int    `json:"id" binding:"required"`
	}
	var p Param
	if _err := c.Bind(&p); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, _err.Error(), nil))
		return
	}
	g := Group{ID: p.ID}
	if _, _err := g.Merge("", 0, 2); _err != nil {
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
	var currentUser gin.H
	if tmp, ok := c.Get("claims"); ok {
		if claims, ok := tmp.(gin.H); ok {
			currentUser = claims
		}
	}
	var r Role
	if p.ID > 0{
		r.ID = p.ID
	}
	var groupID int
	if p.GroupID > 0 {
		groupID = p.GroupID
	} else {
		groupID = currentUser["group_id"].(int)
	}
	if _, _err := r.Merge(groupID, p.Name, p.Menu, p.State); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.REQERR, _err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", gin.H{"id": r.ID}))
}


func RoleList(c *gin.Context) {
	var roles []Role
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	var total int
	var currentUser gin.H
	if tmp, ok := c.Get("claims"); ok {
		if claims, ok := tmp.(gin.H); ok {
			currentUser = claims
		}
	}
	groupID := currentUser["group_id"].(int)
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
		ID       int    `json:"id"`
	}
	var p Param
	if _err := c.Bind(&p); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, _err.Error(), nil))
		return
	}
	r := Role{ID: p.ID}
	if _, _err := r.Merge(p.ID, "", "", 2); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.REQERR, _err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", nil))
}

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
