// @Title  CompetitionStandardSingleController
// @Description  该文件提供关于操作标准个人比赛的各种方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package controller

import (
	"MGA_OJ/Interface"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"
	"MGA_OJ/util"
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ICompetitionStandardSingleController			定义了标准个人比赛类接口
type ICompetitionStandardSingleController interface {
	Interface.EnterInterface // 包含报名方法
}

// CompetitionStandardSingleController			定义了个人比赛工具类
type CompetitionStandardSingleController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Enter
// @description   报名一篇比赛
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionStandardSingleController) Enter(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应比赛
	id := ctx.Params.ByName("id")

	// TODO 获取参赛人数
	userNum, _ := strconv.Atoi(ctx.DefaultQuery("userNum", "50"))

	var competition model.Competition
	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Competition", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Competition", id).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			// TODO 跳过数据库搜寻competition过程
			goto leap
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Competition", id)
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if c.DB.Where("id = (?)", id).First(&competition).Error != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(competition)
		c.Redis.HSet(ctx, "Competition", id, v)
	}

leap:
	// TODO 查看是否为比赛的创建者
	if user.ID != competition.UserId || user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	// TODO 查看比赛是否已经进行过
	if time.Now().After(time.Time(competition.StartTime)) {
		response.Fail(ctx, nil, "比赛已经超过开始时间")
		return
	}

	// TODO 查看比赛是否为单人比赛
	if competition.Type != "Single" {
		response.Fail(ctx, nil, "并非单人比赛，无法报名")
		return
	}

	// TODO 随机生成指定数量的user
	for i := 0; i < userNum; i++ {
		// TODO 创建用户
		newName := "参赛账号" + util.RandomString(8)
		password := util.RandomString(6)
		hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			response.Response(ctx, 201, 201, nil, "加密错误")
			return
		}
		newUser := model.User{
			Name:       newName,
			Email:      newName,
			Password:   string(hasedPassword),
			Icon:       "MGA" + strconv.Itoa(rand.Intn(9)+1) + ".jpg",
			Score:      1500,
			LikeNum:    0,
			UnLikeNum:  0,
			CollectNum: 0,
			VisitNum:   0,
			Level:      0,
		}
		newUserStandard := model.UserStandard{
			Email:    newName,
			Password: password,
			CID:      competition.ID,
		}
		c.DB.Create(&newUser)
		c.DB.Create(&newUserStandard)
		// TODO 为指定竞赛报名
		var competitionRank model.CompetitionRank
		competitionRank.CompetitionId = competition.ID
		competitionRank.MemberId = user.ID
		competitionRank.Score = 0
		competitionRank.Penalties = 0
		// TODO 插入数据
		if err := c.DB.Create(&competitionRank).Error; err != nil {
			response.Fail(ctx, nil, "数据库存储错误")
			return
		}
		// TODO 加入redis
		c.Redis.ZAdd(ctx, "CompetitionR"+id, redis.Z{Member: user.ID.String(), Score: 0})
	}
	response.Success(ctx, nil, "用户生成成功成功")
}

// @title    EnterCondition
// @description   查看报名状态
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionStandardSingleController) EnterCondition(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应比赛
	id := ctx.Params.ByName("id")

	var competition model.Competition
	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Competition", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Competition", id).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			// TODO 跳过数据库搜寻competition过程
			goto leap
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Competition", id)
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if c.DB.Where("id = (?)", id).First(&competition).Error != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(competition)
		c.Redis.HSet(ctx, "Competition", id, v)
	}

leap:
	// TODO 检查比赛类型
	if competition.Type != "Single" {
		response.Fail(ctx, nil, "不是个人比赛，无法报名")
		return
	}

	var competitionRank model.CompetitionRank

	// TODO 查看是否已经报名
	if _, err := c.Redis.ZScore(ctx, "CompetitionR"+id, user.ID.String()).Result(); err == nil {
		response.Success(ctx, gin.H{"enter": true}, "已报名")
		return
	}
	if c.DB.Where("member_id = (?) and competition_id = (?)", user.ID, competition.ID).First(&competitionRank).Error == nil {
		response.Success(ctx, gin.H{"enter": true}, "已报名")
		{
			// TODO 加入redis
			c.Redis.ZAdd(ctx, "CompetitionR"+id, redis.Z{Member: user.ID.String(), Score: 0})
		}
		return
	}

	response.Success(ctx, gin.H{"enter": false}, "未报名")
	return

}

// @title    CancelEnter
// @description   取消报名一篇比赛
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionStandardSingleController) CancelEnter(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应比赛
	id := ctx.Params.ByName("id")

	var competition model.Competition
	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Competition", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Competition", id).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			// TODO 跳过数据库搜寻competition过程
			goto leap
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Competition", id)
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if c.DB.Where("id = (?)", id).First(&competition).Error != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(competition)
		c.Redis.HSet(ctx, "Competition", id, v)
	}

leap:
	if competition.Type != "Single" {
		response.Fail(ctx, nil, "不是个人比赛，无法报名")
		return
	}

	// TODO 查看比赛是否已经进行过
	if time.Now().After(time.Time(competition.StartTime)) {
		response.Fail(ctx, nil, "比赛已经开始")
		return
	}

	var competitionRank model.CompetitionRank

	// TODO 查看是否已经报名
	if c.DB.Where("member_id = (?) and competition_id = (?)", user.ID, competition.ID).First(&competitionRank).Error != nil {
		response.Fail(ctx, nil, "未报名")
		return
	}

	c.DB.Delete(&competitionRank)
	c.Redis.ZRem(ctx, "CompetitionR"+id, user.ID.String())
	response.Success(ctx, nil, "取消报名成功")
	return
}

// @title    EnterPage
// @description   查看报名列表
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionStandardSingleController) EnterPage(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应比赛
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	var competition model.Competition
	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Competition", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Competition", id).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			// TODO 跳过数据库搜寻competition过程
			goto leap
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Competition", id)
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if c.DB.Where("id = (?)", id).First(&competition).Error != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(competition)
		c.Redis.HSet(ctx, "Competition", id, v)
	}

leap:
	// TODO 查看是否为比赛的创建者
	if user.ID != competition.UserId {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	var userStandards []model.UserStandard

	// TODO 查找所有分页中可见的条目
	c.DB.Where("cid = (?)", competition.ID).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&userStandards)

	var total int64
	c.DB.Where("cid = (?)", competition.ID).Model(model.UserStandard{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"userStandards": userStandards, "total": total}, "成功")
}

// @title    NewCompetitionController
// @description   新建一个ICompetitionStandardSingleController
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   ICompetitionStandardSingleController		返回一个ICompetitionStandardSingleController用于调用各种函数
func NewCompetitionStandardSingleController() ICompetitionStandardSingleController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.UserStandard{})

	return CompetitionStandardSingleController{DB: db, Redis: redis}
}
