// @Title  CompetitionStandardGroupController
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

// ICompetitionStandardGroupController			定义了标准个人比赛类接口
type ICompetitionStandardGroupController interface {
	Interface.EnterInterface // 包含报名方法
}

// CompetitionStandardGroupController			定义了个人比赛工具类
type CompetitionStandardGroupController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Enter
// @description   报名一篇比赛
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionStandardGroupController) Enter(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应比赛
	id := ctx.Params.ByName("id")

	// TODO 获取小组数量
	groupNum, _ := strconv.Atoi(ctx.DefaultQuery("groupNum", "20"))
	// TODO 获取小组人数
	userNum, _ := strconv.Atoi(ctx.DefaultQuery("userNum", "3"))

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

	// TODO 查看比赛是否为小组比赛
	if competition.Type != "Group" {
		response.Fail(ctx, nil, "并非小组比赛，无法报名")
		return
	}

	// TODO 生成指定数量的小组
	for i := 0; i < groupNum; i++ {
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
		group := model.Group{
			Title:         competition.Title + "参赛组",
			Content:       competition.Title + "参赛组",
			Auto:          false,
			LeaderId:      user.ID,
			CompetitionAt: competition.EndTime,
		}
		if competition.HackTime.After(competition.EndTime) {
			group.CompetitionAt = competition.HackTime
		}
		c.DB.Create(&group)
		groupList := model.UserList{
			GroupId: group.ID,
			UserId:  user.ID,
		}
		c.DB.Create(&groupList)
		// TODO 随机生成指定数量的user
		for i := 1; i < userNum; i++ {
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
			groupList := model.UserList{
				GroupId: group.ID,
				UserId:  user.ID,
			}
			c.DB.Create(&groupList)
		}
		// TODO 为指定竞赛报名
		var competitionRank model.CompetitionRank
		competitionRank.CompetitionId = competition.ID
		competitionRank.MemberId = group.ID
		competitionRank.Score = 0
		competitionRank.Penalties = 0
		// TODO 插入数据
		if err := c.DB.Create(&competitionRank).Error; err != nil {
			response.Fail(ctx, nil, "数据库存储错误")
			return
		}
		// TODO 加入redis
		c.Redis.ZAdd(ctx, "CompetitionR"+id, redis.Z{Member: group.ID.String(), Score: 0})
	}
	response.Success(ctx, nil, "用户小组生成成功成功")
}

// @title    EnterCondition
// @description   查看报名状态
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionStandardGroupController) EnterCondition(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应比赛
	competition_id := ctx.Params.ByName("id")

	var competition model.Competition
	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Competition", competition_id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Competition", competition_id).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			// TODO 跳过数据库搜寻competition过程
			goto leap
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Competition", competition_id)
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if c.DB.Where("id = (?)", competition_id).First(&competition).Error != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(competition)
		c.Redis.HSet(ctx, "Competition", competition_id, v)
	}

leap:

	groups, _ := c.Redis.ZRange(ctx, "CompetitionR"+competition.ID.String(), 0, -1).Result()

	for i := range groups {
		// TODO 查找组
		var group model.Group

		// TODO 先看redis中是否存在
		if ok, _ := c.Redis.HExists(ctx, "Group", groups[i]).Result(); ok {
			cate, _ := c.Redis.HGet(ctx, "Group", groups[i]).Result()
			if json.Unmarshal([]byte(cate), &group) == nil {
				goto leep
			} else {
				// TODO 移除损坏数据
				c.Redis.HDel(ctx, "Group", groups[i])
			}
		}

		// TODO 查看用户组是否在数据库中存在
		c.DB.Where("id = (?)", groups[i]).First(&group)
		{
			// TODO 将用户组存入redis供下次使用
			v, _ := json.Marshal(group)
			c.Redis.HSet(ctx, "Group", groups[i], v)
		}
	leep:
		if c.DB.Where("group_id = (?) and user_id = (?)", group.ID, user.ID).First(&model.UserList{}).Error == nil {
			response.Success(ctx, gin.H{"group": group}, "已报名")
			return
		}
	}

	response.Success(ctx, nil, "未报名")
	return

}

// @title    CancelEnter
// @description   取消报名一篇比赛
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionStandardGroupController) CancelEnter(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取小组id
	group_id := ctx.Params.ByName("group_id")

	// TODO 查找对应比赛
	competition_id := ctx.Params.ByName("competition_id")

	var competition model.Competition
	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Competition", competition_id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Competition", competition_id).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			// TODO 跳过数据库搜寻competition过程
			goto leap
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Competition", competition_id)
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if c.DB.Where("id = (?)", competition_id).First(&competition).Error != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(competition)
		c.Redis.HSet(ctx, "Competition", competition_id, v)
	}

leap:

	// TODO 查看比赛是否已经进行过
	if time.Now().After(time.Time(competition.StartTime)) {
		response.Fail(ctx, nil, "比赛已经开始")
		return
	}

	// TODO 检查比赛类型
	if competition.Type != "Group" {
		response.Fail(ctx, nil, "不是组队比赛，无法报名")
		return
	}

	// TODO 查找组
	var group model.Group

	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Group", group_id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Group", group_id).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Group", group_id)
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if c.DB.Where("id = (?)", group_id).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		c.Redis.HSet(ctx, "Group", group_id, v)
	}
leep:
	if group.LeaderId != user.ID {
		response.Fail(ctx, nil, "不是用户组组长")
		return
	}

	var competitionRank model.CompetitionRank

	// TODO 查看是否已经报名
	if c.DB.Where("member_id = (?) and competition_id = (?)", group.ID, competition.ID).First(&competitionRank).Error != nil {
		response.Fail(ctx, nil, "未报名")
		return
	}

	c.DB.Delete(&competitionRank)
	c.Redis.ZRem(ctx, "CompetitionR"+competition.ID.String(), group.ID.String())
	// TODO 修改小组的比赛时间并保存
	group.CreatedAt = model.Time(time.Now())
	c.DB.Save(&group)
	c.Redis.HDel(ctx, "Group", group.ID.String())
	response.Success(ctx, nil, "取消报名成功")
	return
}

// @title    EnterPage
// @description   查看报名列表
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionStandardGroupController) EnterPage(ctx *gin.Context) {
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
// @description   新建一个ICompetitionStandardGroupController
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   ICompetitionStandardGroupController		返回一个ICompetitionStandardGroupController用于调用各种函数
func NewCompetitionStandardGroupController() ICompetitionStandardGroupController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.UserStandard{})

	return CompetitionStandardGroupController{DB: db, Redis: redis}
}
