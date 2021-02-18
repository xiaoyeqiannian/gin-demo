package proc

import (
	"strings"
	. "go-server/database"

	"github.com/casbin/casbin"
	"github.com/casbin/casbin/util"
	gormadapter "github.com/casbin/gorm-adapter"
	"github.com/spf13/viper"
)

type Permission struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

func UpdateCasbin(RoleID string, permissions []Permission) error {
	ClearCasbin(0, RoleID)
	for _, v := range permissions {
		e := Casbin()
		addflag := e.AddPolicy(RoleID, v.Path, v.Method)
		if addflag == false {
			continue
		}
	}
	return nil
}

func GetPolicyPathByRoleID(RoleID string) (permissions []Permission) {
	e := Casbin()
	list := e.GetFilteredPolicy(0, RoleID)
	for _, v := range list {
		permissions = append(permissions, Permission{
			Path:   v[1],
			Method: v[2],
		})
	}
	return permissions
}

func ClearCasbin(v int, p ...string) bool {
	e := Casbin()
	return e.RemoveFilteredPolicy(v, p...)

}

func Casbin() *casbin.Enforcer {
	a := gormadapter.NewAdapterByDB(MysqlDB)
	e := casbin.NewEnforcer(viper.Get("casbin.model-path"), a)
	e.AddFunction("ParamsMatch", ParamsMatchFunc)
	_ = e.LoadPolicy()
	return e
}

func ParamsMatch(fullNameKey1 string, key2 string) bool {
	key1 := strings.Split(fullNameKey1, "?")[0]
	// 剥离路径后再使用casbin的keyMatch2
	return util.KeyMatch2(key1, key2)
}

func ParamsMatchFunc(args ...interface{}) (interface{}, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)

	return ParamsMatch(name1, name2), nil
}
