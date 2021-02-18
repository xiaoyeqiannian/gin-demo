package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var MysqlDB *gorm.DB // MySQL

func InitMysql() {
	var err error

	MysqlDB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		viper.Get("mysql.user"),
		viper.Get("mysql.password"),
		viper.Get("mysql.host"),
		viper.Get("mysql.port"),
		viper.Get("mysql.db"),
		viper.Get("mysql.charset"),
		viper.Get("mysql.parseTime"),
		viper.Get("mysql.loc")))
	if err != nil {
		fmt.Println("failed to connect database:", err)
		return
	}
	MysqlDB.LogMode(true)
	fmt.Println("connect database success", &MysqlDB)
	// 如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响
	MysqlDB.SingularTable(true)
	// defer MysqlDB.Close()
}
