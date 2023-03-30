// @Title  LevelController
// @Description  该文件提供关于操作权限的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"
	"strconv"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"

	"gorm.io/gorm"
)

// ILevelController			定义了权限类接口
type ILevelController interface {
	Update(ctx *gin.Context) // 更新某用户的权限
}

// LevelController			定义了权限工具类
type LevelController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Update
// @description   更新某用户的权限
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (l LevelController) Update(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取指定等级
	level, err := strconv.Atoi(ctx.Params.ByName("level"))

	if err != nil {
		response.Fail(ctx, nil, "权限等级有误")
		return
	}

	if level >= user.Level {
		response.Fail(ctx, nil, "权限等级大于等于了你的权限等级")
		return
	}

	// TODO 获取指定用户
	id := ctx.Params.ByName("id")

	var userb model.User

	// TODO 先看redis中是否存在
	if ok, _ := l.Redis.HExists(ctx, "User", id).Result(); ok {
		cate, _ := l.Redis.HGet(ctx, "User", id).Result()
		if json.Unmarshal([]byte(cate), &userb) == nil {
			goto leap
		} else {
			// TODO 移除损坏数据
			l.Redis.HDel(ctx, "User", id)
		}
	}

	// TODO 查看用户是否在数据库中存在
	if l.DB.Where("id = ?", id).First(&userb).Error != nil {
		response.Fail(ctx, nil, "用户不存在")
		return
	}
	{
		// TODO 将用户存入redis供下次使用
		v, _ := json.Marshal(userb)
		l.Redis.HSet(ctx, "User", id, v)
	}
leap:

	// TODO 查看指定用户的等级
	if userb.Level >= user.Level {
		response.Fail(ctx, nil, "无法修改该用户的权限等级")
		return
	}

	// TODO 更新
	userb.Level = level
	l.DB.Save(&userb)

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    NewLevelController
// @description   新建一个ILevelController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   ILetterController		返回一个ILetterController用于调用各种函数
func NewLevelController() ILevelController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	return LevelController{DB: db, Redis: redis}
}
