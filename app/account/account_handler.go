package account

import (
	"net/http"
	"strconv"
	"fmt"
	. "go-server/app/account/model"
	"go-server/app/account/proc"
	. "go-server/database"
	"go-server/utils"
	"github.com/gin-gonic/gin"
)


func AccountModify(c *gin.Context) {
	type Param struct {
        Name     string `json:"name"`
		RoleID   int    `json:"role_id"`
		State    int8   `json:"state"`
		Avatar   string `json:"avatar"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	var currentUser gin.H
	if tmp, ok := c.Get("claims"); ok {
		if claims, ok := tmp.(gin.H); ok {
			currentUser = claims
		}
	}
	var param Param
	if err := c.Bind(&param); err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, err.Error(), nil))
		return
	}
	u := User{ID: currentUser["user_id"].(int)}
	if _, _err := u.Update(param.State, param.RoleID, param.Password, param.Avatar, param.Email); _err != nil {
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
		fmt.Println("searchKey:", searchKey)
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
			roleName = role.Name
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


func AccountDel(c *gin.Context) {

}


func AccountMenu(c *gin.Context){

}


func GroupModify(c *gin.Context){
	type Param struct {
		ID       int    `json:"id"`
        Name     string `json:"name"`
		Kind     int8    `json:"kind"`
		State    int8   `json:"state"`
	}
	var currentUser gin.H
	if tmp, ok := c.Get("claims"); ok {
		if claims, ok := tmp.(gin.H); ok {
			currentUser = claims
		}
	}
	var param Param
	if err := c.Bind(&param); err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, err.Error(), nil))
		return
	}
	var g Group
	fmt.Println("new group", g, g.ID)
	if param.ID <= 0 {
		param.ID = currentUser["group_id"].(int)
	}
	if _err := g.Merge(param.ID, param.Name, param.Kind, param.State); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.REQERR, _err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", nil))
	return
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
	return
}


func GroupDel(c *gin.Context){

}


func RoleModify(c *gin.Context) {
	type Param struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		GroupID  int    `json:"group_id"`
		State    int8   `json:"state"`
		Permissions []proc.Permission `json:"permissions"`
	}
	var param Param
	if err := c.Bind(&param); err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, err.Error(), nil))
		return
	}
	var role Role
	if _err := role.Merge(param.ID, param.GroupID, param.Name, param.State); _err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.REQERR, _err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", nil))
	return
}


func RoleList(c *gin.Context) {
	var roles []Role
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	var total int
	sql := MysqlDB.Table("roles")
	sql = sql.Count(&total)
	sql.Limit(pageSize).Offset((page - 1) * pageSize).Find(&roles)
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", gin.H{"items": roles, "total": total, "page": page, "per_page": pageSize}))
}


func RoleDel(c *gin.Context) {
	currentUser, err := proc.CurrentUser(c)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, err.Error(), nil))
		return
	}
	fmt.Println(currentUser["user_id"])
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", nil))
}

func ManagerRolePermission(c *gin.Context) {
	if c.Request.Method == "POST" {
		type Param struct {
			RoleID      string            `json:"role_id"`
			Permissions []proc.Permission `json:"permissions"`
		}
		var param Param
		if err := c.Bind(&param); err != nil {
			c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, err.Error(), nil))
			return
		}
		err := proc.UpdateCasbin(param.RoleID, param.Permissions)
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
