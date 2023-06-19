// @Title  HackController
// @Description  该文件提供关于操作hack的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"

	"encoding/json"

	"github.com/go-redis/redis/v9"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IHackController			定义了hack类接口
type IHackController interface {
	Show(ctx *gin.Context)    // 查询指定的hack
	ShowNum(ctx *gin.Context) // 查看某个比赛中某用户的hack数量
}

// HackController			定义了文章工具类
type HackController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    ShowNum
// @description  查看某个比赛中某用户的hack数量
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (h HackController) ShowNum(ctx *gin.Context) {
	// TODO 获取path中的id
	member_id := ctx.Params.ByName("member_id")
	competition_id := ctx.Params.ByName("competition_id")

	var hack model.HackNum
	if ok, _ := h.Redis.HExists(ctx, "HackNum", competition_id+member_id).Result(); ok {
		cate, _ := h.Redis.HGet(ctx, "HackNum", competition_id+member_id).Result()
		if json.Unmarshal([]byte(cate), &hack) == nil {
			// TODO 跳过数据库搜寻Hack过程
			goto leap
		} else {
			// TODO 移除损坏数据
			h.Redis.HDel(ctx, "HackNum", competition_id+member_id)
		}
	}

	// TODO 查看hack是否在数据库中存在
	if h.DB.Where("competition_id = (?) and member_id = (?)", competition_id, member_id).First(&hack).Error != nil {
		response.Fail(ctx, nil, "hackNum不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(hack)
		h.Redis.HSet(ctx, "HackNum", competition_id+member_id, v)
	}

leap:
	response.Success(ctx, gin.H{"hackNum": hack}, "查看成功")
}

// @title    Show
// @description   查看Hack功能
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (h HackController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	var hack model.Hack
	if ok, _ := h.Redis.HExists(ctx, "Hack", id).Result(); ok {
		cate, _ := h.Redis.HGet(ctx, "Hack", id).Result()
		if json.Unmarshal([]byte(cate), &hack) == nil {
			// TODO 跳过数据库搜寻Hack过程
			goto leap
		} else {
			// TODO 移除损坏数据
			h.Redis.HDel(ctx, "Hack", id)
		}
	}

	// TODO 查看hack是否在数据库中存在
	if h.DB.Where("id = (?)", id).First(&hack).Error != nil {
		response.Fail(ctx, nil, "hack不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(hack)
		h.Redis.HSet(ctx, "Hack", id, v)
	}

leap:
	response.Success(ctx, gin.H{"hack": hack}, "查看成功")
}

// @title    NewHackController
// @description   新建一个IHackController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   IHackController		返回一个IHackController用于调用各种函数
func NewHackController() IHackController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.Hack{})
	db.AutoMigrate(model.HackNum{})
	return HackController{DB: db, Redis: redis}
}
