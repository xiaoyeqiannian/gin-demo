package model


import (
	// "fmt"
	"time"
	"errors"
	"gin-server/database"
)


type Group struct {
	ID        int       `gorm:"size:11;primary_key;AUTO_INCREMENT;not null" json:"id"`
	Name      string    `gorm:"size:32"   json:"name"`
	Kind      int8      `gorm:"size:1"   json:"kind"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	State     int8      `gorm:"size:1; DEFAULT:1" json:"state"`
}

// create or modify
func (g *Group) Update(name string, kind, state int8) (id int, err error) {
	if g.ID == 0 {
		return 0, errors.New("Group is gone")
	}
	if _err := database.MysqlDB.First(g, g.ID).Error; _err != nil {
		return 0, _err
	}
	d := make(map[string]interface{})
	if state != 0 {
		d["State"] = state
	}
	if kind != 0 {
		d["Kind"] = kind
	}
	if name != "" {
		d["Name"] = name
	}
	database.MysqlDB.Model(&g).Updates(d)
	return g.ID, nil
}


func (g *Group) Add(name string, kind int8) (id int, err error) {
	g.Name = name
	g.Kind = kind
	database.MysqlDB.Create(&g)
	return g.ID, nil
}


func (g *Group) Del(IDs []int, group_id int) (err error) {
	if group_id != GROUP_SYS_ADMIN_ID {
		return errors.New("Fail")
	}
	sql := database.MysqlDB.Table("group").Where("id in (?)", IDs)
	sql.Update("state", STATUS_DELETED)
	return nil
}


func (g *Group) ParseState() string {
	if g.State == STATUS_VALID {
		return "有效"
	} else if g.State == STATUS_DELETED {
		return "无效"
	} else {
		return "未知"
	}
}
