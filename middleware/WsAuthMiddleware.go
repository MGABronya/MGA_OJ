// @Title  WsAuthMiddleware
// @Description  中间件，用于解析ws的token
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package middleware

import (
	"MGA_OJ/common"
	"MGA_OJ/model"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

// @title    AuthMiddleware
// @description   中间件，用于ws解析token
// @auth      MGAronya       2022-9-16 12:15
// @param    void
// @return   gin.HandlerFunc	将token解析完毕后传回上下文
func WsAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO 获取 token
		tokenString := ctx.Query("token")

		fmt.Print("请求token", tokenString)

		// TODO validate token formate
		if tokenString == "" {
			ctx.JSON(201, gin.H{
				"code": 201,
				"msg":  "格式错误，权限不足",
			})
			ctx.Abort()
			return
		}

		token, claims, err := common.ParseToken(tokenString)

		if err != nil || !token.Valid {
			ctx.JSON(201, gin.H{
				"code": 201,
				"msg":  "解析错误，权限不足",
			})
			fmt.Print("解析错误，权限不足")
			ctx.Abort()
			return
		}

		// TODO token通过验证, 获取claims中的UserID
		id := fmt.Sprint(claims.UserId)
		DB := common.GetDB()
		Redis := common.GetRedisClient(0)
		var user model.User

		// TODO 先看redis中是否存在
		if ok, _ := Redis.HExists(ctx, "User", id).Result(); ok {
			cate, _ := Redis.HGet(ctx, "User", id).Result()
			if json.Unmarshal([]byte(cate), &user) == nil {
				goto leep
			} else {
				// TODO 移除损坏数据
				Redis.HDel(ctx, "User", id)
			}
		}

		// TODO 查看用户是否在数据库中存在
		if DB.Where("id = (?)", id).First(&user).Error != nil {
			ctx.JSON(201, gin.H{
				"code": 201,
				"msg":  "用户不存在，权限不足",
			})
			fmt.Print("用户不存在，权限不足")
			ctx.Abort()
			return
		}

	leep:
		// TODO 用户存在 将user信息写入上下文
		ctx.Set("user", user)

		ctx.Next()
	}
}
