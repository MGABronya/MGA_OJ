// @Title  CompetitionController
// @Description  该文件提供关于操作比赛的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"MGA_OJ/Interface"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"
	"MGA_OJ/vo"
	"encoding/json"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

// ICompetitionController			定义了比赛类接口
type ICompetitionController interface {
	Interface.RestInterface      // 包含增删查改功能
	RankList(ctx *gin.Context)   // 获取比赛排名情况
	RankMember(ctx *gin.Context) // 获取某用户的排名情况
	MemberShow(ctx *gin.Context) // 获取某成员每道题的罚时情况
}

// CompetitionController			定义了比赛工具类
type CompetitionController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Create
// @description   创建一篇比赛
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionController) Create(ctx *gin.Context) {
	var competitionRequest vo.CompetitionRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&competitionRequest); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 取出用户权限
	if user.Level < 2 {
		response.Fail(ctx, nil, "用户权限不足")
		return
	}

	// TODO 查看表单是否存在
	var set model.Set
	if c.DB.Where("id = ?", competitionRequest.SetId).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}

	competition := model.Competition{
		UserId:    user.ID,
		SetId:     competitionRequest.SetId,
		StartTime: competitionRequest.StartTime,
		EndTime:   competitionRequest.EndTime,
		Type:      competitionRequest.Type,
	}

	// TODO 插入数据
	if err := c.DB.Create(&competition).Error; err != nil {
		response.Fail(ctx, nil, "比赛上传出错，数据库存储错误")
		return
	}

	// TODO 将所有竞赛内题目做上标记
	var topicLists []model.TopicList
	c.DB.Where("set_id = ?", set.ID).Find(&topicLists)
	for _, topicList := range topicLists {
		var problemLists []model.ProblemList
		c.DB.Where("topic_id = ?", topicList.TopicId).Find(&problemLists)
		for _, problemList := range problemLists {
			var problem model.Problem
			c.DB.Where("id = ?", problemList.ProblemId).First(&problem)
			// TODO 如果竞赛标记已经被标记
			if problem.CompetitionId != 0 {
				response.Fail(ctx, nil, "题目不存在")
				return
			}
			problem.CompetitionId = competition.ID
			// TODO 做上标记
			c.DB.Save(&problem)
		}
	}

	// TODO 成功
	response.Success(ctx, gin.H{"competition": competition}, "创建成功")
}

// @title    Update
// @description   更新一篇比赛的内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionController) Update(ctx *gin.Context) {
	var competitionUpdate vo.CompetitionUpdate
	// TODO 数据验证
	if err := ctx.ShouldBind(&competitionUpdate); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 取出用户权限
	if user.Level < 2 {
		response.Fail(ctx, nil, "用户权限不足")
		return
	}

	// TODO 查找对应比赛
	id := ctx.Params.ByName("id")

	var competition model.Competition
	if c.DB.Where("id = ?", id).First(&competition) != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}

	// TODO 更新比赛内容
	c.DB.Where("id = ?", id).Updates(competitionUpdate)

	// TODO 移除损坏数据
	c.Redis.HDel(ctx, "Competition", id)

	// TODO 成功
	response.Success(ctx, nil, "更新成功")
}

// @title    Show
// @description   查看一篇比赛的内容
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var competition model.Competition

	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Competition", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Competition", id).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			response.Success(ctx, gin.H{"competition": competition}, "成功")
			return
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Competition", id)
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if c.DB.Where("id = ?", id).First(&competition).Error != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}

	response.Success(ctx, gin.H{"competition": competition}, "成功")

	// TODO 将竞赛存入redis供下次使用
	v, _ := json.Marshal(competition)
	c.Redis.HSet(ctx, "Competition", id, v)
}

// @title    Delete
// @description   删除一篇比赛
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionController) Delete(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	var competition model.Competition

	// TODO 查看比赛是否存在
	if c.DB.Where("id = ?", id).First(&competition).Error != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}

	// TODO 判断当前用户是否为比赛的作者
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有操作比赛的权力
	if user.Level < 2 {
		response.Fail(ctx, nil, "用户权限不足")
		return
	}

	// TODO 删除比赛
	c.DB.Delete(&competition)

	// TODO 移除损坏数据
	c.Redis.HDel(ctx, "Competition", id)

	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description   获取多篇比赛
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionController) PageList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var competitions []model.Competition

	// TODO 查找所有分页中可见的条目
	c.DB.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&competitions)

	var total int64
	c.DB.Model(model.Competition{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"competitions": competitions, "total": total}, "成功")
}

// @title    RankList
// @description   获取当前比赛排名，包含ac题目数量和罚时
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionController) RankList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取比赛id
	id := ctx.Params.ByName("id")

	var err error

	// TODO 查找所有分页中可见的条目
	mems, err := c.Redis.ZRevRangeWithScores(ctx, "Competition"+id, int64(pageNum)*int64(pageSize), int64(pageNum)*int64(pageSize)+int64(pageSize)-1).Result()

	if err != nil {
		// TODO 尝试从数据库中找出相关数据
		var members []model.CompetitionRank
		var total int64
		c.DB.Where("competition_id = ?", id).Order("accept_num desc penalties asc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&members)
		// TODO 返回数据
		response.Success(ctx, gin.H{"members": members, "total": total}, "成功")
		return
	} else {
		// TODO 将redis中的数据取出
		total, _ := c.Redis.ZCard(ctx, "Competition"+id).Result()
		members := make([]model.CompetitionRank, pageSize)

		for i := range mems {
			id, _ := strconv.Atoi(id)
			members[i].CompetitionId = uint(id)
			members[i].MemberId = mems[i].Member.(uint)
			members[i].AcceptNum = uint(math.Ceil(mems[i].Score))
			members[i].Penalties = time.Unix(int64((float64(members[i].AcceptNum)-mems[i].Score))*10000000000, 0)
		}
		// TODO 返回数据
		response.Success(ctx, gin.H{"members": members, "total": total}, "成功")
		return
	}
}

// @title    RankMember
// @description   获取当前某成员的比赛排名信息
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionController) RankMember(ctx *gin.Context) {
	// TODO 获取比赛id
	competition_id := ctx.Params.ByName("competition")

	// TODO 获取成员id
	member_id := ctx.Params.ByName("member")

	var err error

	// TODO 获得当前排名
	rank, err := c.Redis.ZRevRank(ctx, "Competition"+competition_id, member_id).Result()

	if err != nil {
		// 从数据库中取出
		c.DB.Table("competition_ranks").Select("RANK() OVER(partition by competition_id order by accept_num desc penalties asc)").Where("competition_id = ? and member_id = ?", competition_id, member_id).Scan(&rank)
	}

	// TODO 返回数据
	response.Success(ctx, gin.H{"rank": rank}, "成功")
}

// @title    MemberShow
// @description   获取当前某成员的比赛排名信息
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionController) MemberShow(ctx *gin.Context) {
	// TODO 获取比赛id
	competition_id := ctx.Params.ByName("competition")

	// TODO 获取成员id
	member_id := ctx.Params.ByName("member")

	var competitionMember []model.CompetitionMember

	var err error

	cM, err := c.Redis.HGet(ctx, "competition"+competition_id, member_id).Result()

	if err != nil {
		// TODO 去数据库中找
		c.DB.Where("competition_id = ? and member_id = ?", competition_id, member_id).Find(&competitionMember)
		// TODO 返回数据
		response.Success(ctx, gin.H{"competitionMember": competitionMember}, "成功")
	} else {
		json.Unmarshal([]byte(cM), &competitionMember)
		// TODO 返回数据
		response.Success(ctx, gin.H{"competitionMember": competitionMember}, "成功")
	}
}

// @title    NewCompetitionController
// @description   新建一个ICompetitionController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   ICompetitionController		返回一个ICompetitionController用于调用各种函数
func NewCompetitionController() ICompetitionController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.Competition{})
	db.AutoMigrate(model.CompetitionRank{})
	db.AutoMigrate(model.CompetitionMember{})
	return CompetitionController{DB: db, Redis: redis}
}
