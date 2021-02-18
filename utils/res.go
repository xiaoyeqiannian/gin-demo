package utils

const (
	REQOK                  = "0000" // 正常
    DBERR               = "2000"
    THIRDERR            = "2001"
    DATAERR             = "2002"
    IOERR               = "2003"

    LOGINERR            = "2100" // 登陆错误
    PARAMERR            = "2101" // 参数错误
    USERERR             = "2102" // 用户异常
    ROLEERR             = "2103" // 权限错误
    PWDERR              = "2104" // 密码错误
    VERIFYERR           = "2105" // 验证错误

    REQERR              = "2200"

    NODATA              = "2300" // 无数据
    UNDERDEBUG          = "2301" // debug模式下无法使用

    UNKOWNERR           = "2400"
)

type Res struct {
	Code string         `json:"code"`
	Message  string  `json:"message"`
	Data interface{} `json:"data"`
}

func RespJson(code string, desc string, data interface{}) Res {
	var res Res
	res.Code = code
	res.Message = desc
	res.Data = data
	return res
}
