package model


import (
	// "fmt"
	"time"
	// "errors"
	"gin-server/database"
)


type User struct {
	ID        int       `gorm:"size:11;primary_key;AUTO_INCREMENT;not null" json:"id"`
	Name      string    `gorm:"size:32"  json:"name"`
	Avatar    string    `gorm:"size:128"  json:"avatar"`
	Password  string    `gorm:"size:128" json:"-"`
	Email     string    `gorm:"size:64" json:"email"`
	RoleID    int       `gorm:"size:11"  json:"role_id"`
	GroupID   int       `gorm:"size:11"  json:"group_id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	State     int8      `gorm:"size:1; DEFAULT:1; COMMENT:'0:未激活，1:有效，2:已拉黑'" json:"state"`
}

// TableName overrides the table name used by User to `user`
func (User) TableName() string {
	return "user"
}

// func init() {
// 	m := database.MysqlDB.AutoMigrate(&User{})
// 	fmt.Println("create User table", &m)
// }


func (u *User) Merge(state int8, roleID, groupID int, name, pwd, avatar, email string) (id int, err error) {
	if u.ID > 0 {
		if _err := database.MysqlDB.First(u, u.ID).Error; _err != nil {
			return 0, _err
		}
		d := make(map[string]interface{})
		if state != 0 {
			d["State"] = state
		}
		if roleID != 0 {
			d["RoleID"] = roleID
		}
		if groupID != 0 {
			d["GroupID"] = groupID
		}
		if pwd != "" {
			d["Name"] = name
		}
		if pwd != "" {
			d["Password"] = pwd
		}
		if avatar != "" {
			d["Avatar"] = avatar
		}
		if email != "" {
			d["Email"] = email
		}
		database.MysqlDB.Model(&u).Updates(d)
	} else {
		u.State = state
		u.RoleID = roleID
		u.Password = pwd
		u.Avatar = avatar
		u.Email = email
		u.GroupID = groupID
		u.Name = name
		database.MysqlDB.Create(&u)
	}
	return u.ID, nil
}


func (u *User) ParseState() string {
	if u.State == 1 {
		return "已激活"
	} else if u.State == 2 {
		return "已拉黑"
	} else {
		return "未激活"
	}
}

