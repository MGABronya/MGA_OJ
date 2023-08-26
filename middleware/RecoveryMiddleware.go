// @Title  RecoveryMiddleware
// @Description  该中间件用于“拦截”运行时恐慌的内建函数,防止程序崩溃
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:46
package middleware

import (
	"MGA_OJ/response"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// @title    RecoveryMiddleware
// @description   该中间件用于“拦截”运行时恐慌的内建函数,防止程序崩溃
// @auth      MGAronya             2022-9-16 10:19
// @param     void        void    		  无入参
// @return    HandlerFunc        gin.HandlerFunc            返回一个响应函数
func RecoveryMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			// TODO recover用于“拦截”运行时恐慌的内建函数,防止程序崩溃
			if err := recover(); err != nil {
				response.Fail(ctx, nil, fmt.Sprint(err))
				log.Println(err)
			}
		}()
		ctx.Next()
	}
}
