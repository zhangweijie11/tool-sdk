package schemas

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	// 响应状态码
	CurdStatusOkCode  int    = 200
	CurdStatusOkMsg   string = "Success"
	CurdCreatFailCode int    = -400200
)

func ReturnJson(Context *gin.Context, httpCode int, dataCode int, msg string, data interface{}) {
	//Context.Header("key2020","value2020")  	//可以根据实际情况在头部添加额外的其他信息
	Context.JSON(httpCode, gin.H{
		"code": dataCode,
		"msg":  msg,
		"data": data,
	})
}

// SuccessGet 获取成功
func SuccessGet(ctx *gin.Context, data interface{}) {
	ReturnJson(ctx, http.StatusOK, CurdStatusOkCode, CurdStatusOkMsg, data)
}

// SuccessDelete 删除成功
func SuccessDelete(ctx *gin.Context, data interface{}) {
	ReturnJson(ctx, http.StatusNoContent, CurdStatusOkCode, CurdStatusOkMsg, data)
}

// SuccessUpdate 更新成功
func SuccessUpdate(ctx *gin.Context, data interface{}) {
	ReturnJson(ctx, http.StatusCreated, CurdStatusOkCode, CurdStatusOkMsg, data)
}

// SuccessCreate 创建成功
func SuccessCreate(ctx *gin.Context, data interface{}) {
	ReturnJson(ctx, http.StatusCreated, CurdStatusOkCode, CurdStatusOkMsg, data)
}

// Fail 失败的请求
func Fail(ctx *gin.Context, message string) {
	ReturnJson(ctx, http.StatusBadRequest, CurdCreatFailCode, message, nil)
}
