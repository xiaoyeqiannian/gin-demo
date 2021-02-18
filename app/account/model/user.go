package model


import (
	// "fmt"
	"time"

	"go-server/database"
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


// func init() {
// 	m := database.MysqlDB.AutoMigrate(&User{})
// 	fmt.Println("create User table", &m)
// }


func (u *User) Update(state int8, roleID int, pwd, avatar, email string) (code int, err error) {
	if err = database.MysqlDB.First(u, u.ID).Error; err != nil {
		return
	}
	d := make(map[string]interface{})
	if state != 0 {
		d["State"] = state
	}
	if roleID != 0 {
		d["RoleID"] = roleID
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
	return
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

