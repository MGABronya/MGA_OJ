// @Title  RejudgeController
// @Description  该文件提供关于操作重判的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"MGA_OJ/common"
	"MGA_OJ/model"
	rabbitMq "MGA_OJ/rabbitMq"
	"MGA_OJ/response"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v9"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IRejudgeController			定义了重判类接口
type IRejudgeController interface {
	Do(ctx *gin.Context)                // 进行重判
	CompetitionDelete(ctx *gin.Context) // 对比赛结果进行清空
	CompetitionScore(ctx *gin.Context)  // 对比赛结果重新进行分数统计
}

// RejudgeController			定义了重判工具类
type RejudgeController struct {
	DB       *gorm.DB           // 含有一个数据库指针
	Redis    *redis.Client      // 含有一个redis指针
	Rabbitmq *rabbitMq.RabbitMQ // 含有一个消息中间件
}

// @title    Do
// @description   进行重判
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RejudgeController) Do(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 取出用户权限
	if user.Level < 2 {
		response.Fail(ctx, nil, "用户权限不足")
		return
	}

	// TODO 获取查询条件
	problem_id := ctx.DefaultQuery("problem_id", "")
	user_id := ctx.DefaultQuery("user_id", "")
	competition_id := ctx.DefaultQuery("competition_id", "")
	start_time := ctx.DefaultQuery("start_time", "")
	end_time := ctx.DefaultQuery("end_time", "")
	language := ctx.DefaultQuery("language", "")
	condition := ctx.DefaultQuery("condition", "")

	db := common.GetDB()

	var records []model.Record

	// TODO 搜索对应问题
	if problem_id != "" {
		db = db.Where("problem_id = ?", problem_id)
	}

	// TODO 搜索对应用户
	if user_id != "" {
		db = db.Where("user_id = ?", user_id)
	}

	// TODO 搜索对应比赛
	if competition_id != "" {
		db = db.Where("competition_id = ?", competition_id)
	}

	// TODO 搜索对应起始时间
	if start_time != "" {
		db = db.Where("created_at >= ?", start_time)
	}

	// TODO 搜索对应截至时间
	if end_time != "" {
		db = db.Where("created_at <= ?", end_time)
	}

	// TODO 搜索对应语言
	if language != "" {
		db = db.Where("language = ?", language)
	}

	// TODO 搜索对应状态
	if condition != "" {
		db = db.Where("condition = ?", condition)
	}

	// TODO 查找记录组
	db.Find(&records)

	// TODO 加入消息队列
	for _, record := range records {
		// TODO 加入消息队列
		if err := r.Rabbitmq.PublishSimple(fmt.Sprint(record.ID)); err != nil {
			response.Fail(ctx, nil, "消息队列出错")
			return
		}

		{
			// TODO 将提交存入redis供判题机使用
			v, _ := json.Marshal(record)
			r.Redis.HSet(ctx, "Record", fmt.Sprint(record.ID), v)
		}
	}
	response.Success(ctx, nil, "已全部重新提交")
}

// @title    CompetitionDelete
// @description   对比赛结果清空
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RejudgeController) CompetitionDelete(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 取出用户权限
	if user.Level < 2 {
		response.Fail(ctx, nil, "用户权限不足")
		return
	}

	// TODO 取出比赛id
	id := ctx.Params.ByName("id")

	// TODO 清空db中的比赛排行
	r.DB.Where("competition_id = ?", id).Delete(&model.CompetitionRank{})
	r.DB.Where("competition_id = ?", id).Delete(&model.CompetitionMember{})

	// TODO 清空redis中的比赛排行
	r.Redis.Del(ctx, "competition"+id)
	r.Redis.Del(ctx, "Competition"+id)

	response.Success(ctx, nil, "清除完成")
}

// @title    CompetitionScore
// @description   对比赛分数重新计算
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RejudgeController) CompetitionScore(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 取出用户权限
	if user.Level < 2 {
		response.Fail(ctx, nil, "用户权限不足")
		return
	}

	// TODO 取出比赛id
	id := ctx.Params.ByName("id")

	var competition model.Competition
	// TODO 先看redis中是否存在
	if ok, _ := r.Redis.HExists(ctx, "Competition", id).Result(); ok {
		cate, _ := r.Redis.HGet(ctx, "Competition", id).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			goto leap
		} else {
			// TODO 移除损坏数据
			r.Redis.HDel(ctx, "Competition", id)
		}
	}
	// TODO 在数据库中查找
	if r.DB.Where("id = ?", id).First(&competition).Error != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
leap:
	// TODO 查找分数变化情况
	var userScoreChanges []model.UserScoreChange
	r.DB.Where("competition_id = ?", id).Find(&userScoreChanges)

	// TODO 删除分数变化
	r.DB.Where("competition_id = ?", id).Delete(&model.UserScoreChange{})

	// TODO 回滚分数
	for _, userScoreChange := range userScoreChanges {
		var user model.User
		if r.DB.Where("id = ?", userScoreChange.UserId).First(&user).Error != nil {
			continue
		}
		user.Score -= userScoreChange.ScoreChange
		r.DB.Save(&user)
	}

	// TODO 整理比赛结果
	CompetitionFinish(ctx, r.Redis, r.DB, competition)

	response.Success(ctx, nil, "重新计算完成")
}

// @title    NewRejudgeController
// @description   新建一个IRejudgeController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   IRejudgeController		返回一个IRejudgeController用于调用各种函数
func NewRejudgeController() IRejudgeController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	rabbitmq := rabbitMq.NewRabbitMQSimple("MGAronya")
	return RejudgeController{DB: db, Redis: redis, Rabbitmq: rabbitmq}
}
