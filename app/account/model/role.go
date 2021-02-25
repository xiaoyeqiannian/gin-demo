package model


import (
	"fmt"
	"time"
	"errors"

	"gin-server/database"
)


type Role struct {
	ID        int       `gorm:"size:11;primary_key;AUTO_INCREMENT;not null" json:"id"`
	Name      string    `gorm:"size:32"   json:"name"`
	Menu      string    `gorm:"type:text"   json:"menu"`
	GroupID   int       `gorm:"size:11"   json:"group_id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	State     int8      `gorm:"size:1; DEFAULT:1" json:"state"`
}


// create or modify
func (r *Role) Merge(groupID int, name string, menu string, state int8) (id int, err error) {
	if r.ID > 0 {
		if _err := database.MysqlDB.First(r, "state=?", 1).Error; _err != nil {
			return 0, _err
		}
		d := make(map[string]interface{})
		if state != 0 {
			d["State"] = state
		}
		if groupID != 0 {
			d["GroupID"] = groupID
		}
		if name != "" {
			d["Name"] = name
		}
		if menu != "" {
			d["Menu"] = menu
		}
		database.MysqlDB.Model(&r).Updates(d)
	}else {
		if _err := database.MysqlDB.First(r, "group_id=? and name=? and state=?", groupID, name, 1).Error; _err == nil {
			if r.ID > 0 {
				return r.ID, errors.New(fmt.Sprintf("role %s in %d is created", name, groupID))
			}
		}
		if groupID == 0 {
			return 0, errors.New("The group_id is require")
		}
		r.State = 1
		r.GroupID = groupID
		r.Name = name
		r.Menu = menu
		database.MysqlDB.Create(&r)
	}
	return r.ID, nil
}

