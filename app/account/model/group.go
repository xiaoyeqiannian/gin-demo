package model


import (
	"fmt"
	"time"
	"errors"
	"go-server/database"
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
func (g *Group) Merge(id int, name string, kind, state int8) error {
	if id > 0 {
		if err := database.MysqlDB.First(g, id).Error; err != nil {
			return errors.New(fmt.Sprintf("can not find group %d", id))
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
	}else {
		g.Name = name
		g.Kind = kind
		database.MysqlDB.Create(&g)
	}
	return nil
}


func (g *Group) ParseState() string {
	if g.State == 1 {
		return "有效"
	} else if g.State == 2 {
		return "无效"
	} else {
		return "未知"
	}
}
