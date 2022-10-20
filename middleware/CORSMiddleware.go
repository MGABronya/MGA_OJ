// @Title  CORSMiddleware
// @Description  该中间件用于处理跨域问题
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:41
package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @title    CORSMiddleware
// @description   该中间件用于处理跨域问题
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     void        void    		  无入参
// @return    HandlerFunc        gin.HandlerFunc            返回一个响应函数
func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method

		// TODO 请求头部
		origin := ctx.Request.Header.Get("Origin")

		// TODO 处理头部
		if origin != "" {
			// TODO 接收客户端发送的origin （重要！）
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			// TODO 服务器支持的所有跨域请求的方法
			ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			// TODO 允许跨域设置可以返回其他子段，可以自定义字段
			ctx.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			// TODO 允许浏览器（客户端）可以解析的头部 （重要）
			ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			// TODO 设置缓存时间
			ctx.Header("Access-Control-Max-Age", "172800")
			// TODO 允许客户端传递校验信息比如 cookie (重要)
			ctx.Header("Access-Control-Allow-Credentials", "true")
		}

		// TODO 允许类型校验
		if method == "OPTIONS" {
			ctx.JSON(http.StatusOK, "ok!")
		}

		// TODO 捕获panic
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()
	}
}
