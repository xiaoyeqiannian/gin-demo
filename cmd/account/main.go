package main

import (
	"fmt"
	"go-server/app/account"
	"go-server/database"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("account")
	viper.AddConfigPath("configs")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("获取配置文件失败!!", err.Error())
	}

	database.InitMysql()

	r := gin.Default()
	defer database.MysqlDB.Close()
	// Logging to a file.
	// f, _ := os.Create("gin.log")
	// gin.DefaultWriter = io.MultiWriter(f)
	// gin.DefaultWriter = io.MultiWriter(f, os.Stdout)// 如果需要同时将日志写入文件和控制台，请使用以下代码。
	r = account.InitRouter(r)
	r.Run(fmt.Sprintf(":%s", viper.Get("server.port")))
}
