package account

import (
	// "encoding/json"

	// "fmt"
	"net/http"
	// "strconv"
	// "time"
	"go-server/utils"

	// "go-server/app/account/middleware"
	// . "go-server/app/account/model"
	// . "go-server/database"

	// "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)


// @Tags 文件
// @Summary 上传文件
// @version 1.0
// @Description 上传文件
// @Accept  multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Success 0 {string} json "{"code": 0,"msg": "ok","data":{"src": "http://xx.xxx.com/xxx/df13e16abdd0c8317966dbe06cb20778"}}"
// @Failure -100 {string} json "{"code": -100,"msg": "参数错误","data": null}"
// @Router /api/v1/file/upload [post]
func AppUploadFile(c *gin.Context) {
	file, _ := c.FormFile("file")
	if file == nil {
		c.JSON(http.StatusOK, utils.RespJson(utils.PARAMERR, "参数错误", nil))
		return
	}

	fileID := "xxx/xxxx"
	if b, err := file.Open(); err == nil {
		out := make([]byte, file.Size)
		_, err := b.Read(out)
		if err != nil {
			c.JSON(http.StatusOK, utils.RespJson(utils.UNKOWNERR, err.Error(), nil))
			return
		}
		fileID, err = utils.StorageSave(out)
		if err != nil {
			c.JSON(http.StatusOK, utils.RespJson(utils.UNKOWNERR, err.Error(), nil))
			return
		}
	}
	c.JSON(http.StatusOK, utils.RespJson(utils.REQOK, "ok", gin.H{"src": "http://xxx.xxx.com/" + fileID}))
}
