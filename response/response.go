// @Title  response
// @Description  用于统一封装各种返回格式，方便调用
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @title    Response
// @description   对返回值的简单封装
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     ctx *gin.Context, httpStatus int, code int, data gin.H, msg string  入参为响应需写进上下文的基本信息
// @return    void void    没有返回值
func Response(ctx *gin.Context, httpStatus int, code int, data gin.H, msg string) {
	ctx.JSON(httpStatus, gin.H{"code": code, "data": data, "msg": msg})
}

// @title    Response
// @description   对成功状态返回值的简单封装
// @auth      MGAronya（张健）             2022-9-16 10:22
// @param     ctx *gin.Context, data gin.H, msg string	入参为响应需写进上下文的基本信息
// @return    void void    没有返回值
func Success(ctx *gin.Context, data gin.H, msg string) {
	Response(ctx, http.StatusOK, 200, data, msg)
}

// @title    Response
// @description   对失败状态返回值的简单封装
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     ctx *gin.Context, data gin.H, msg string  入参为响应需写进上下文的基本信息
// @return    void void    没有返回值
func Fail(ctx *gin.Context, data gin.H, msg string) {
	Response(ctx, http.StatusOK, 400, data, msg)
}
