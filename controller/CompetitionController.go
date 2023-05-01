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
	"MGA_OJ/util"
	"MGA_OJ/vo"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ICompetitionController			定义了比赛类接口
type ICompetitionController interface {
	Interface.PasswdInterface     // 包含密码方法
	Interface.RestInterface       // 包含了增删查改功能
	Interface.SearchInterface     // 搜索功能
	Interface.LabelInterface      // 标签功能
	Interface.RejudgeInterface    // 包含重判相关功能
	RankList(ctx *gin.Context)    // 获取比赛排名情况
	RankMember(ctx *gin.Context)  // 获取某用户的排名情况
	RollingList(ctx *gin.Context) // 订阅比赛滚榜
	MemberShow(ctx *gin.Context)  // 获取某成员每道题的罚时情况
}

// CompetitionController			定义了个人比赛工具类
type CompetitionController struct {
	DB       *gorm.DB            // 含有一个数据库指针
	Redis    *redis.Client       // 含有一个redis指针
	UpGrader *websocket.Upgrader // 用于持久化连接
	RabbitMq *common.RabbitMQ    // 一个消息队列的指针
}

var InitCompetition = map[string]func(ctx *gin.Context, redis *redis.Client, db *gorm.DB, competition model.Competition){}

var FinishCompetition = map[string]func(ctx *gin.Context, redis *redis.Client, db *gorm.DB, competition model.Competition){}

// @title    CreatePasswd
// @description   为某个比赛创建密码
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionController) CreatePasswd(ctx *gin.Context) {

	// TODO 查找对应比赛
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

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
	if c.DB.Where("id = ?", id).First(&competition) != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(competition)
		c.Redis.HSet(ctx, "Competition", id, v)
	}

leap:

	// TODO 查看比赛是否已经进行过
	if time.Now().After(time.Time(competition.StartTime)) {
		response.Fail(ctx, nil, "比赛已经开始")
		return
	}

	// TODO 查看是否为比赛作者
	if user.ID != competition.UserId {
		response.Fail(ctx, nil, "不是比赛作者，请勿非法操作")
		return
	}

	var passwd model.Passwd
	// TODO 数据验证
	if err := ctx.ShouldBind(&passwd); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	passwd.ID = uuid.NewV4()

	// TODO 创建密码
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(passwd.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, 201, 201, nil, "加密错误")
		return
	}

	passwd.Password = string(hasedPassword)

	if c.DB.Create(&passwd).Error != nil {
		response.Fail(ctx, nil, "密码上传出错，数据库存储错误")
		return
	}

	c.DB.Where("id = ?", competition.PasswdId).Delete(&model.Passwd{})

	// TODO 存储新密码
	competition.PasswdId = passwd.ID
	c.DB.Save(&competition)

	// TODO 返回数据
	response.Success(ctx, nil, "成功")
	return

}

// @title    DeletePasswd
// @description   为某个比赛删除密码
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionController) DeletePasswd(ctx *gin.Context) {

	// TODO 查找对应比赛
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

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
	if c.DB.Where("id = ?", id).First(&competition) != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(competition)
		c.Redis.HSet(ctx, "Competition", id, v)
	}

leap:

	// TODO 查看比赛是否已经进行过
	if time.Now().After(time.Time(competition.StartTime)) {
		response.Fail(ctx, nil, "比赛已经开始")
		return
	}

	// TODO 查看是否为比赛作者
	if user.ID != competition.UserId {
		response.Fail(ctx, nil, "不是比赛作者，请勿非法操作")
		return
	}

	c.DB.Where("id = ?", competition.PasswdId).Delete(&model.Passwd{})

	// TODO 返回数据
	response.Success(ctx, nil, "成功")
	return

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

	// TODO 验证起始时间与终止时间是否合法
	if competitionRequest.StartTime.After(competitionRequest.EndTime) {
		response.Fail(ctx, nil, "起始时间大于了终止时间")
		return
	}
	if time.Now().After(time.Time(competitionRequest.StartTime)) {
		response.Fail(ctx, nil, "起始时间大于了当前时间")
		return
	}
	if time.Time(competitionRequest.EndTime).After(time.Now().Add(30 * 24 * time.Hour)) {
		response.Fail(ctx, nil, "终止时间不可设置为30日后")
		return
	}
	if time.Time(competitionRequest.HackTime).After(time.Now().Add(30 * 24 * time.Hour)) {
		response.Fail(ctx, nil, "hack不可设置为30日后")
		return
	}
	if _, ok := InitCompetition[competitionRequest.Type]; !ok {
		response.Fail(ctx, nil, "比赛类型不合法")
		return
	}
	// TODO 查看是否为小组组长
	if (competitionRequest.GroupId != uuid.UUID{}) {
		var group model.Group

		// TODO 先看redis中是否存在
		if ok, _ := c.Redis.HExists(ctx, "Group", competitionRequest.GroupId.String()).Result(); ok {
			cate, _ := c.Redis.HGet(ctx, "Group", competitionRequest.GroupId.String()).Result()
			if json.Unmarshal([]byte(cate), &group) == nil {
				goto leep
			} else {
				// TODO 移除损坏数据
				c.Redis.HDel(ctx, "Group", competitionRequest.GroupId.String())
			}
		}

		// TODO 查看用户组是否在数据库中存在
		if c.DB.Where("id = ?", competitionRequest.GroupId.String()).First(&group).Error != nil {
			response.Fail(ctx, nil, "用户组不存在")
			return
		}
		{
			// TODO 将用户组存入redis供下次使用
			v, _ := json.Marshal(group)
			c.Redis.HSet(ctx, "Group", competitionRequest.GroupId.String(), v)
		}
	leep:
		if user.ID != group.LeaderId {
			response.Fail(ctx, nil, "不是用户组组长")
			return
		}
	}

	if competitionRequest.LessNum > competitionRequest.UpNum {
		response.Fail(ctx, nil, "人数限制有误")
		return
	}

	// TODO 比赛创建
	competition := model.Competition{
		UserId:    user.ID,
		StartTime: competitionRequest.StartTime,
		EndTime:   competitionRequest.EndTime,
		Title:     competitionRequest.Title,
		Content:   competitionRequest.Content,
		Reslong:   competitionRequest.Reslong,
		Resshort:  competitionRequest.Resshort,
		HackTime:  competitionRequest.HackTime,
		HackScore: competitionRequest.HackScore,
		HackNum:   competitionRequest.HackNum,
		Type:      competitionRequest.Type,
		GroupId:   competitionRequest.GroupId,
		LessNum:   competitionRequest.LessNum,
		UpNum:     competitionRequest.UpNum,
	}

	// TODO 插入数据
	if err := c.DB.Create(&competition).Error; err != nil {
		response.Fail(ctx, nil, "比赛上传出错，数据库存储错误")
		return
	}

	// TODO 成功
	response.Success(ctx, gin.H{"competition": competition}, "创建成功")

	// TODO 等待直至比赛结束
	StartTimer(ctx, c.Redis, c.DB, competition)
}

// @title    Update
// @description   更新一篇比赛的内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionController) Update(ctx *gin.Context) {
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

	// TODO 验证起始时间与终止时间是否合法
	if competitionRequest.StartTime.After(competitionRequest.EndTime) {
		response.Fail(ctx, nil, "起始时间大于了终止时间")
		return
	}
	if time.Now().After(time.Time(competitionRequest.StartTime)) {
		response.Fail(ctx, nil, "起始时间大于了当前时间")
		return
	}
	if time.Time(competitionRequest.EndTime).After(time.Now().Add(30 * 24 * time.Hour)) {
		response.Fail(ctx, nil, "终止时间不可设置为30日后")
		return
	}
	if time.Time(competitionRequest.HackTime).After(time.Now().Add(30 * 24 * time.Hour)) {
		response.Fail(ctx, nil, "hack不可设置为30日后")
		return
	}

	// TODO 查找对应比赛
	id := ctx.Params.ByName("id")

	var competition model.Competition
	if c.DB.Where("id = ?", id).First(&competition) != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}

	// TODO 查看比赛是否已经进行过
	if time.Now().After(time.Time(competition.StartTime)) {
		response.Fail(ctx, nil, "比赛已经进行过")
		return
	}

	if user.ID != competition.UserId {
		response.Fail(ctx, nil, "不是作者，无法修改比赛信息")
		return
	}
	// TODO 查看是否为小组组长
	if (competitionRequest.GroupId != uuid.UUID{}) {
		var group model.Group

		// TODO 先看redis中是否存在
		if ok, _ := c.Redis.HExists(ctx, "Group", competitionRequest.GroupId.String()).Result(); ok {
			cate, _ := c.Redis.HGet(ctx, "Group", competitionRequest.GroupId.String()).Result()
			if json.Unmarshal([]byte(cate), &group) == nil {
				goto leep
			} else {
				// TODO 移除损坏数据
				c.Redis.HDel(ctx, "Group", competitionRequest.GroupId.String())
			}
		}

		// TODO 查看用户组是否在数据库中存在
		if c.DB.Where("id = ?", competitionRequest.GroupId.String()).First(&group).Error != nil {
			response.Fail(ctx, nil, "用户组不存在")
			return
		}
		{
			// TODO 将用户组存入redis供下次使用
			v, _ := json.Marshal(group)
			c.Redis.HSet(ctx, "Group", competitionRequest.GroupId.String(), v)
		}
	leep:
		if user.ID != group.LeaderId {
			response.Fail(ctx, nil, "不是用户组组长")
			return
		}
	}
	if competitionRequest.LessNum > competitionRequest.UpNum {
		response.Fail(ctx, nil, "人数限制有误")
		return
	}

	competitionUpdate := model.Competition{
		StartTime: competitionRequest.StartTime,
		EndTime:   competitionRequest.EndTime,
		Title:     competitionRequest.Title,
		Content:   competitionRequest.Content,
		Reslong:   competitionRequest.Reslong,
		Resshort:  competitionRequest.Resshort,
		HackTime:  competitionRequest.HackTime,
		HackScore: competitionRequest.HackScore,
		HackNum:   competitionRequest.HackNum,
		GroupId:   competitionRequest.GroupId,
		LessNum:   competitionRequest.LessNum,
		UpNum:     competitionRequest.UpNum,
	}

	// TODO 更新比赛内容
	c.DB.Where("id = ?", id).Updates(competitionUpdate)

	// TODO 更新定时器
	util.TimerMap[competition.ID].Reset(time.Until(time.Time(competitionUpdate.StartTime)))

	// TODO 找到那些参赛组并修改endtime
	if competition.Type == "Group" || competition.Type == "Match" {
		var competitionRanks []model.CompetitionRank

		// TODO 查找所有分页中可见的条目
		c.DB.Where("competition_id = ?", competition.ID).Find(&competitionRanks)

		// TODO 查找那些组并更新他们
		for i := range competitionRanks {
			c.DB.Model(model.Group{}).Where("id = ?", competitionRanks[i].MemberId).Update("competition_at", competitionUpdate.EndTime)
			c.Redis.HDel(ctx, "Group", competitionRanks[i].MemberId.String())
		}
	}

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
	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Competition", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Competition", id).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Competition", id)
		}
	}

	// TODO 查看比赛是否存在
	if c.DB.Where("id = ?", id).First(&competition).Error != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
leep:
	// TODO 判断当前用户是否为比赛的作者
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有操作比赛的权力
	if user.ID != competition.ID && user.Level < 4 {
		response.Fail(ctx, nil, "比赛不属于您，请勿非法操作")
		return
	}

	// TODO 删除比赛
	c.DB.Delete(&competition)

	// TODO 移除损坏数据
	c.Redis.HDel(ctx, "Competition", id)

	// TODO 移除计时器
	util.TimerMap[competition.ID].Stop()
	// TODO 此处会使比赛强行退出并出现panic
	delete(util.TimerMap, competition.ID)

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

// @title    LabelCreate
// @description   标签创建
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionController) LabelCreate(ctx *gin.Context) {
	// TODO 获取指定比赛
	id := ctx.Params.ByName("id")

	// TODO 获取标签
	label := ctx.Params.ByName("label")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看比赛是否存在
	var competition model.Competition

	// TODO 先尝试在redis中寻找
	if ok, _ := c.Redis.HExists(ctx, "Competition", id).Result(); ok {
		art, _ := c.Redis.HGet(ctx, "Competition", id).Result()
		if json.Unmarshal([]byte(art), &competition) == nil {
			goto leep
		} else {
			// TODO 解码失败，删除字段
			c.Redis.HDel(ctx, "Competition", id)
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if c.DB.Where("id = ?", id).First(&competition).Error != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
	{
		// TODO 将比赛存入redis供下次使用
		v, _ := json.Marshal(competition)
		c.Redis.HSet(ctx, "Competition", id, v)
	}
leep:

	// TODO 查看是否为比赛作者
	if competition.UserId != user.ID {
		response.Fail(ctx, nil, "不是比赛作者，请勿非法操作")
		return
	}

	// TODO 创建标签
	competitionLabel := model.CompetitionLabel{
		Label:         label,
		CompetitionId: competition.ID,
	}

	// TODO 插入数据
	if err := c.DB.Create(&competitionLabel).Error; err != nil {
		response.Fail(ctx, nil, "竞赛标签上传出错，数据验证有误")
		return
	}

	// TODO 解码失败，删除字段
	c.Redis.HDel(ctx, "CompetitionLabel", id)

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    LabelDelete
// @description   标签删除
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionController) LabelDelete(ctx *gin.Context) {
	// TODO 获取指定比赛
	id := ctx.Params.ByName("id")

	// TODO 获取标签
	label := ctx.Params.ByName("label")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看比赛是否存在
	var competition model.Competition

	// TODO 先尝试在redis中寻找
	if ok, _ := c.Redis.HExists(ctx, "Competition", id).Result(); ok {
		art, _ := c.Redis.HGet(ctx, "Competition", id).Result()
		if json.Unmarshal([]byte(art), &competition) == nil {
			goto leep
		} else {
			// TODO 解码失败，删除字段
			c.Redis.HDel(ctx, "Competition", id)
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if c.DB.Where("id = ?", id).First(&competition).Error != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
	{
		// TODO 将比赛存入redis供下次使用
		v, _ := json.Marshal(competition)
		c.Redis.HSet(ctx, "Competition", id, v)
	}
leep:

	// TODO 查看是否为比赛作者
	if competition.UserId != user.ID {
		response.Fail(ctx, nil, "不是比赛作者，请勿非法操作")
		return
	}

	// TODO 删除比赛标签
	if c.DB.Where("id = ?", label).First(&model.CompetitionLabel{}).Error != nil {
		response.Fail(ctx, nil, "标签不存在")
		return
	}

	c.DB.Where("id = ?", label).Delete(&model.CompetitionLabel{})

	// TODO 解码失败，删除字段
	c.Redis.HDel(ctx, "CompetitionLabel", id)

	// TODO 成功
	response.Success(ctx, nil, "删除成功")
}

// @title    LabelShow
// @description   标签查看
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionController) LabelShow(ctx *gin.Context) {
	// TODO 获取指定表单
	id := ctx.Params.ByName("id")

	// TODO 查找数据
	var competitionLabels []model.CompetitionLabel
	// TODO 先尝试在redis中寻找
	if ok, _ := c.Redis.HExists(ctx, "CompetitionLabel", id).Result(); ok {
		art, _ := c.Redis.HGet(ctx, "CompetitionLabel", id).Result()
		if json.Unmarshal([]byte(art), &competitionLabels) == nil {
			goto leap
		} else {
			// TODO 解码失败，删除字段
			c.Redis.HDel(ctx, "CompetitionLabel", id)
		}
	}

	// TODO 在数据库中查找
	c.DB.Where("competition_id = ?", id).Find(&competitionLabels)
	{
		// TODO 将题目标签存入redis供下次使用
		v, _ := json.Marshal(competitionLabels)
		c.Redis.HSet(ctx, "CompetitionLabel", id, v)
	}

leap:

	// TODO 成功
	response.Success(ctx, gin.H{"competitionLabels": competitionLabels}, "查看成功")
}

// @title    Search
// @description   文本搜索
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionController) Search(ctx *gin.Context) {
	// TODO 获取文本
	text := ctx.Params.ByName("text")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	var competitions []model.Competition

	// TODO 模糊匹配
	c.DB.Where("match(title,content,res_long,res_short) against(? in boolean mode)", text+"*").Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&competitions)

	// TODO 查看查询总数
	var total int64
	c.DB.Where("match(title,content,res_long,res_short) against(? in boolean mode)", text+"*").Model(model.Competition{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"competitions": competitions, "total": total}, "成功")
}

// @title    SearchLabel
// @description   指定标签的搜索
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionController) SearchLabel(ctx *gin.Context) {

	var requestLabels vo.LabelsRequest

	// TODO 获取标签
	if err := ctx.ShouldBind(&requestLabels); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 通过标签寻找
	var competitionIds []struct {
		CompetitionId uuid.UUID `json:"competition_id"` // 竞赛外键
	}

	// TODO 进行标签匹配
	c.DB.Distinct("competition_id").Where("label in (?)", requestLabels.Labels).Model(model.CompetitionLabel{}).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&competitionIds)

	// TODO 查看查询总数
	var total int64
	c.DB.Distinct("competition_id").Where("label in (?)", requestLabels.Labels).Model(model.CompetitionLabel{}).Count(&total)

	// TODO 查找对应表单
	var competitions []model.Competition

	c.DB.Where("id in (?)", competitionIds).Find(&competitions)

	// TODO 返回数据
	response.Success(ctx, gin.H{"competitions": competitions, "total": total}, "成功")
}

// @title    SearchWithLabel
// @description   指定标签与文本的搜索
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionController) SearchWithLabel(ctx *gin.Context) {

	// TODO 获取文本
	text := ctx.Params.ByName("text")

	var requestLabels vo.LabelsRequest

	// TODO 获取标签
	if err := ctx.ShouldBind(&requestLabels); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 通过标签寻找
	var competitionIds []struct {
		CompetitionId uuid.UUID `json:"competition_id"` // 竞赛外键
	}

	// TODO 进行标签匹配
	c.DB.Distinct("competition_id").Where("label in (?)", requestLabels.Labels).Model(model.CompetitionLabel{}).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&competitionIds)

	// TODO 查找对应表单
	var competitions []model.Competition

	// TODO 模糊匹配
	c.DB.Where("id in (?) and match(title,content,res_long,res_short) against(? in boolean mode)", competitionIds, text+"*").Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&competitions)

	// TODO 查看查询总数
	var total int64
	c.DB.Where("id in (?) and match(title,content,res_long,res_short) against(? in boolean mode)", competitionIds, text+"*").Model(model.Competition{}).Count(&total)

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

	// TODO 查找所有分页中可见的条目
	mems, err := c.Redis.ZRevRangeWithScores(ctx, "CompetitionR"+id, int64(pageNum-1)*int64(pageSize), int64(pageNum-1)*int64(pageSize)+int64(pageSize)-1).Result()

	if err != nil {
		// TODO 尝试从数据库中找出相关数据
		var members []model.CompetitionRank
		var total int64
		c.DB.Where("competition_id = ?", id).Order("score desc penalties asc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&members)
		// TODO 返回数据
		response.Success(ctx, gin.H{"members": members, "total": total}, "成功")
		return
	} else {
		// TODO 将redis中的数据取出
		total, _ := c.Redis.ZCard(ctx, "CompetitionR"+id).Result()
		members := make([]model.CompetitionRank, pageSize)

		for i := range mems {
			members[i].CompetitionId = uuid.FromStringOrNil(id)
			members[i].MemberId = mems[i].Member.(uuid.UUID)
			members[i].Score = uint(math.Ceil(mems[i].Score))
			members[i].Penalties = time.Duration((float64(members[i].Score) - mems[i].Score) * 10000000000)
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
	rank, err := c.Redis.ZRevRank(ctx, "CompetitionR"+competition_id, member_id).Result()

	if err != nil {
		// 从数据库中取出
		c.DB.Table("competition_ranks").Select("RANK() OVER(partition by competition_id order by score desc penalties asc)").Where("competition_id = ? and member_id = ?", competition_id, member_id).Scan(&rank)
	}

	// TODO 返回数据
	response.Success(ctx, gin.H{"rank": rank}, "成功")
}

// @title    RollingList
// @description   监听滚榜
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionController) RollingList(ctx *gin.Context) {
	// TODO 获取指定比赛
	id := ctx.Params.ByName("id")

	var competition model.Competition

	// TODO 查看比赛是否存在
	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Competition", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Competition", id).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			goto leep
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
	{
		// TODO 将比赛存入redis供下次使用
		v, _ := json.Marshal(competition)
		c.Redis.HSet(ctx, "Competition", id, v)
	}
leep:

	// TODO 查看比赛是否在进行中
	if !time.Now().After(time.Time(competition.StartTime)) {
		response.Fail(ctx, nil, "比赛不在进行中")
		return
	}
	if competition.HackTime.After(competition.HackTime) {
		if time.Now().After(time.Time(competition.HackTime)) {
			response.Fail(ctx, nil, "比赛不在进行中")
			return
		}
	} else if time.Now().After(time.Time(competition.EndTime)) {
		response.Fail(ctx, nil, "比赛不在进行中")
		return
	}

	// TODO 订阅消息
	pubSub := c.Redis.Subscribe(ctx, "CompetitionChan"+id)
	defer pubSub.Close()
	// TODO 获得消息管道
	ch := pubSub.Channel()

	// TODO 升级get请求为webSocket协议
	ws, err := c.UpGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	// TODO 监听消息
	for msg := range ch {
		// TODO 读取ws中的数据
		_, _, err := ws.ReadMessage()
		// TODO 断开连接
		if err != nil {
			break
		}
		var rk vo.RankList
		json.Unmarshal([]byte(msg.Payload), &rk)
		// TODO 写入ws数据
		ws.WriteJSON(rk)
	}
}

// @title    MemberShow
// @description   获取某成员每道题的罚时情况
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionController) MemberShow(ctx *gin.Context) {
	// TODO 获取比赛id
	competition_id := ctx.Params.ByName("competition")

	// TODO 获取成员id
	member_id := ctx.Params.ByName("member")

	var competitionMembers []model.CompetitionMember

	cM, err := c.Redis.HGet(ctx, "Competition"+competition_id, member_id).Result()

	if err != nil {
		// TODO 去数据库中找
		c.DB.Where("competition_id = ? and member_id = ?", competition_id, member_id).Find(&competitionMembers)
		// TODO 返回数据
		response.Success(ctx, gin.H{"competitionMembers": competitionMembers}, "成功")
	} else {
		json.Unmarshal([]byte(cM), &competitionMembers)
		// TODO 返回数据
		response.Success(ctx, gin.H{"competitionMembers": competitionMembers}, "成功")
	}
}

// @title    Rejudge
// @description   进行重判
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionController) Rejudge(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取比赛id
	id := ctx.Params.ByName("id")

	// TODO 查看比赛是否存在
	var competition model.Competition
	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Competition", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Competition", id).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			goto leep
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
	{
		// TODO 将比赛存入redis供下次使用
		v, _ := json.Marshal(competition)
		c.Redis.HSet(ctx, "Competition", id, v)
	}
leep:

	// TODO 取出用户权限
	if user.Level < 4 || competition.UserId != user.ID {
		response.Fail(ctx, nil, "用户权限不足")
		return
	}
	// TODO 查看比赛是否在进行中
	if competition.HackTime.After(competition.HackTime) {
		if time.Now().Before(time.Time(competition.HackTime)) {
			response.Fail(ctx, nil, "比赛未结束")
			return
		}
	} else if time.Now().Before(time.Time(competition.EndTime)) {
		response.Fail(ctx, nil, "比赛未结束")
		return
	}

	// TODO 获取重判条件
	problem_id := ctx.DefaultQuery("problem_id", "")
	user_id := ctx.DefaultQuery("user_id", "")
	start_time := ctx.DefaultQuery("start_time", "")
	end_time := ctx.DefaultQuery("end_time", "")
	language := ctx.DefaultQuery("language", "")
	condition := ctx.DefaultQuery("condition", "")
	t := ctx.DefaultQuery("time", "0")

	db := common.GetDB()

	var records []model.RecordCompetition

	// TODO 搜索对应问题
	if problem_id != "" {
		db = db.Where("problem_id = ?", problem_id)
	}

	// TODO 搜索对应用户
	if user_id != "" {
		db = db.Where("user_id = ?", user_id)
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

	// TODO 获取rejudge的启动时间
	next, err := time.ParseDuration(t)
	if err != nil {
		response.Fail(ctx, nil, "时间格式有误")
		return
	}
	if next > 48*time.Hour {
		response.Fail(ctx, nil, "定时不可超过48小时")
		return
	}
	timer := time.NewTimer(next)
	response.Success(ctx, nil, "定时成功")
	// TODO 等待
	<-timer.C

	// TODO 加入消息队列
	for _, record := range records {
		{
			// TODO 将提交存入redis供判题机使用
			v, _ := json.Marshal(record)
			c.Redis.HSet(ctx, "RecordCompetition", fmt.Sprint(record.ID), v)
		}

		// TODO 加入消息队列
		recordRabbit := vo.RecordRabbit{
			RecordId: record.ID,
			Type:     competition.Type,
		}
		v, _ := json.Marshal(recordRabbit)
		if err := c.RabbitMq.PublishSimple(string(v)); err != nil {
			log.Println("消息队列出错", err)
		}
	}
}

// @title    CompetitionDataDelete
// @description   对比赛结果清空
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionController) CompetitionDataDelete(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 取出比赛id
	id := ctx.Params.ByName("id")
	// TODO 查看比赛是否存在
	var competition model.Competition
	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Competition", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Competition", id).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			goto leep
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
	{
		// TODO 将比赛存入redis供下次使用
		v, _ := json.Marshal(competition)
		c.Redis.HSet(ctx, "Competition", id, v)
	}
leep:

	// TODO 取出用户权限
	if user.Level < 4 || competition.UserId != user.ID {
		response.Fail(ctx, nil, "用户权限不足")
		return
	}
	// TODO 查看比赛是否在进行中
	if competition.HackTime.After(competition.HackTime) {
		if time.Now().Before(time.Time(competition.HackTime)) {
			response.Fail(ctx, nil, "比赛未结束")
			return
		}
	} else if time.Now().Before(time.Time(competition.EndTime)) {
		response.Fail(ctx, nil, "比赛未结束")
		return
	}

	member_id := ctx.DefaultQuery("member_id", "")

	db := common.GetDB()

	db = db.Where("competition_id = ?", id)

	// TODO 搜索对应用户
	if member_id != "" {
		db = db.Where("member_id = ?", member_id)
	}

	// TODO 清空db中的比赛排行
	db.Model(&model.CompetitionRank{}).Updates(map[string]interface{}{
		"score":     0,
		"penalties": 0,
	})

	// TODO 查看hack数量
	var hackNums []model.HackNum

	db.Find(hackNums)

	for i := range hackNums {
		c.DB.Where("competition_id = ? and member_id = ?", hackNums[i].CompetitionId, hackNums[i].MemberId).Update("score", hackNums[i].Score)
	}

	// TODO 删除比赛中的通过情况
	db.Delete(&model.CompetitionMember{})

	// TODO 查找分数变化情况
	var userScoreChanges []model.UserScoreChange
	db.Find(&userScoreChanges)

	// TODO 删除分数变化
	db.Delete(&model.UserScoreChange{})

	// TODO 回滚分数
	for _, userScoreChange := range userScoreChanges {
		var user model.User
		if c.DB.Where("id = ?", userScoreChange.UserId).First(&user).Error != nil {
			continue
		}
		user.Score -= userScoreChange.ScoreChange
		c.DB.Save(&user)
	}

	// TODO 清空redis中的比赛排行
	c.Redis.Del(ctx, "Competition"+id)
	c.Redis.Del(ctx, "CompetitionR"+id)

	response.Success(ctx, nil, "清除完成")
}

// @title    NewCompetitionController
// @description   新建一个ICompetitionController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   ICompetitionController		返回一个ICompetitionController用于调用各种函数
func NewCompetitionController() ICompetitionController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	rabbitmq := common.GetRabbitMq()
	upGrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	db.AutoMigrate(model.Competition{})
	db.AutoMigrate(model.CompetitionRank{})
	db.AutoMigrate(model.CompetitionMember{})
	db.AutoMigrate(model.CompetitionLabel{})
	db.AutoMigrate(model.RecordCompetition{})
	db.AutoMigrate(model.Passwd{})
	return CompetitionController{DB: db, Redis: redis, UpGrader: upGrader, RabbitMq: rabbitmq}
}

// @title    StartTimer
// @description   建立一个比赛开始定时器
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    competitionId uuid.UUID	比赛id
// @return   void
func StartTimer(ctx *gin.Context, redis *redis.Client, db *gorm.DB, competition model.Competition) {

	util.TimerMap[competition.ID] = time.NewTimer(time.Until(time.Time(competition.StartTime)))
	// TODO 等待比赛开始
	<-util.TimerMap[competition.ID].C
	// TODO 比赛初始事项
	InitCompetition[competition.Type](ctx, redis, db, competition)

	// TODO 创建比赛结束定时器
	util.TimerMap[competition.ID] = time.NewTimer(time.Until(time.Time(competition.EndTime)))

	// TODO 等待比赛结束
	<-util.TimerMap[competition.ID].C

	// TODO 等待hack时间结束
	if competition.HackTime.After(competition.EndTime) {
		util.TimerMap[competition.ID] = time.NewTimer(time.Until(time.Time(competition.HackTime)))
		<-util.TimerMap[competition.ID].C
	}

	FinishCompetition[competition.Type](ctx, redis, db, competition)
}

// @title    CompetitionProblemSubmit
// @description   将自定义的题目，连同提交记录一起提交到题库
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    competition 		对应比赛
// @return   void
func CompetitionProblemSubmit(ctx *gin.Context, redis *redis.Client, db *gorm.DB, competition model.Competition) {
	var problemNews []model.ProblemNew

	db.Where("competition_id = ?", competition.ID).Find(&problemNews)

	for i := range problemNews {
		// TODO 创建题目
		problem := model.Problem{
			Title:        problemNews[i].Title,
			TimeLimit:    problemNews[i].TimeLimit,
			MemoryLimit:  problemNews[i].MemoryLimit,
			Description:  problemNews[i].Description,
			Reslong:      problemNews[i].Reslong,
			Resshort:     problemNews[i].Resshort,
			Input:        problemNews[i].Input,
			Output:       problemNews[i].Output,
			Hint:         problemNews[i].Hint,
			Source:       problemNews[i].Source,
			UserId:       competition.UserId,
			SpecialJudge: problemNews[i].SpecialJudge,
		}

		// TODO 插入数据
		if err := db.Create(&problem).Error; err != nil {
			continue
		}

		var caseSamples []model.CaseSample
		// TODO 先看redis中是否存在
		if ok, _ := redis.HExists(ctx, "CaseSample", problemNews[i].ID.String()).Result(); ok {
			cate, _ := redis.HGet(ctx, "CaseSample", problemNews[i].ID.String()).Result()
			if json.Unmarshal([]byte(cate), &caseSamples) == nil {
				// TODO 跳过数据库搜寻caseSample过程
				goto levp
			} else {
				// TODO 移除损坏数据
				redis.HDel(ctx, "CaseSample", problemNews[i].ID.String())
			}
		}
		db.Where("problem_id = ?", problem.ID).Find(&caseSamples)

		// TODO 将题目存入redis供下次使用
		{
			v, _ := json.Marshal(caseSamples)
			redis.HSet(ctx, "CaseSample", problemNews[i].ID.String(), v)
		}

	levp:

		// TODO 存储测试样例
		for i := range caseSamples {
			// TODO 尝试存入数据库
			cas := model.CaseSample{
				ProblemId: problem.ID,
				Input:     caseSamples[i].Input,
				Output:    caseSamples[i].Output,
				CID:       uint(i + 1),
			}
			// TODO 插入数据
			db.Create(&cas)
		}

		// TODO 查找用例
		var cases []model.Case
		if ok, _ := redis.HExists(ctx, "Case", problemNews[i].ID.String()).Result(); ok {
			cate, _ := redis.HGet(ctx, "Case", problemNews[i].ID.String()).Result()
			if json.Unmarshal([]byte(cate), &cases) == nil {
				// TODO 跳过数据库搜寻testInputs过程
				goto Case
			} else {
				// TODO 移除损坏数据
				redis.HDel(ctx, "Case", problemNews[i].ID.String())
			}
		}

		// TODO 查看题目是否在数据库中存在
		if db.Where("id = ?", problemNews[i].ID.String()).Find(&cases).Error != nil {
			continue
		}
		// TODO 将题目存入redis供下次使用
		{
			v, _ := json.Marshal(cases)
			redis.HSet(ctx, "Case", problemNews[i].ID.String(), v)
		}
	Case:

		// TODO 存储测试用例
		for i := range cases {
			// TODO 尝试存入数据库
			cas := model.Case{
				ProblemId: problem.ID,
				Input:     cases[i].Input,
				Output:    cases[i].Output,
				CID:       uint(i + 1),
			}
			db.Create(&cas)
		}
		var recordSingles []model.RecordCompetition
		db.Where("problem_id = ?", problemNews[i].ID).Find(&recordSingles)
		for i := range recordSingles {
			record := model.Record{
				UserId:    recordSingles[i].UserId,
				ProblemId: problem.ID,
				Code:      recordSingles[i].Code,
				Language:  recordSingles[i].Language,
				Condition: recordSingles[i].Condition,
				Pass:      recordSingles[i].Pass,
				HackId:    recordSingles[i].HackId,
				CreatedAt: recordSingles[i].CreatedAt,
			}
			db.Create(&record)
		}
	}
}

// @title    CompetitionFinish
// @description   整理比赛结果
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    competition 		对应比赛
// @return   void
func CompetitionFinish(ctx context.Context, redis *redis.Client, db *gorm.DB, competition model.Competition) {
	// TODO 整理比赛结果
	competitionMemberMap, _ := redis.HGetAll(ctx, "Competition"+competition.ID.String()).Result()
	competitionRankrs, _ := redis.ZRevRangeWithScores(ctx, "CompetitionR"+competition.ID.String(), 0, -1).Result()

	// TODO 将具体罚时信息全部读出并存入数据库
	for i := range competitionMemberMap {
		var competitionMember []model.CompetitionMember
		json.Unmarshal([]byte(competitionMemberMap[i]), &competitionMember)
		for j := range competitionMember {
			db.Create(&competitionMember[j])
		}
	}
	// TODO 将排名信息读出并存入数据库
	for i := range competitionRankrs {
		competitionRank := model.CompetitionRank{
			Score:     uint(math.Ceil(competitionRankrs[i].Score)),
			Penalties: time.Duration((float64(uint(math.Ceil(competitionRankrs[i].Score))) - competitionRankrs[i].Score) * 10000000000),
		}
		db.Where("member_id = ?", competitionRankrs[i].Member).Updates(competitionRank)
	}

}
