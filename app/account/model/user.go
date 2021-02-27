package model


import (
	// "fmt"
	"time"
	"errors"
	"encoding/base64"
	"gin-server/utils"
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


func (u *User) Update(state int8, roleID, groupID int, name, pwd, avatar, email string) (id int, err error) {
	if u.ID == 0 {
		return 0, errors.New("user is gone")
	}
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
	if name != "" {
		d["Name"] = name
	}
	if pwd != "" {
		tmp, err := base64.StdEncoding.DecodeString(pwd)
		if err != nil {
			return 0, errors.New("decode password error")
		}
		spwd, err := utils.GeneratePasswordHash(string(tmp))
		if err != nil {
			return 0, errors.New("encode password error")
		}
		d["Password"] = spwd
	}
	if avatar != "" {
		d["Avatar"] = avatar
	}
	if email != "" {
		d["Email"] = email
	}
	database.MysqlDB.Model(&u).Updates(d)
	return u.ID, nil
}


func (u *User) Add(roleID, groupID int, name, pwd, avatar, email string) (id int, err error) {
	u.State = STATUS_VALID
	u.RoleID = roleID
	u.Password = pwd
	u.Avatar = avatar
	u.Email = email
	u.GroupID = groupID
	u.Name = name
	database.MysqlDB.Create(&u)
	return u.ID, nil
}


func (u *User) Del(IDs []int, group_id int) (err error) {
	sql := database.MysqlDB.Table("user").Where("id in (?)", IDs)
	if group_id != GROUP_SYS_ADMIN_ID {
		sql = sql.Where("group_id=?", group_id)
	}
	sql.Update("state", STATUS_DELETED)
	return nil
}


func (u *User) PasswordUpdate(pwd string) (err error) {
	if _, _err := u.Update(0, 0, 0, "", pwd, "", ""); _err != nil {
		return _err
	}
	return nil
}


func (u *User) ParseState() string {
	if u.State == STATUS_VALID {
		return "已激活"
	} else if u.State == STATUS_DELETED {
		return "已拉黑"
	} else {
		return "未激活"
	}
}

