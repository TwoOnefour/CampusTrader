package response

import "github.com/gin-gonic/gin"

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

func Error(ctx *gin.Context, code int, msg string) {
	ctx.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
	})
}
