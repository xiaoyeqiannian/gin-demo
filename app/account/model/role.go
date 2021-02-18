package model


import (
	// "fmt"
	"time"

	"go-server/database"
)


type Role struct {
	ID        int       `gorm:"size:11;primary_key;AUTO_INCREMENT;not null" json:"id"`
	Name      string    `gorm:"size:32"   json:"name"`
	GroupID   int       `gorm:"size:11"   json:"group_id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	State     int8      `gorm:"size:1; DEFAULT:1" json:"state"`
}


// create or modify
func (r *Role) Merge(id, group_id int, name string, state int8) error {
	if id > 0{
		if err := database.MysqlDB.First(r, id).Error; err != nil {
			return err
		}
	}
	d := make(map[string]interface{})
	if state != 0 {
		d["State"] = state
	}
	if group_id != 0 {
		d["GroupID"] = group_id
	}
	if name != "" {
		d["Name"] = name
	}
	database.MysqlDB.Model(&r).Updates(d)
	return nil
}

