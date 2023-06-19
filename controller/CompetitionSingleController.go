// @Title  CompetitionSingleController
// @Description  该文件提供关于操作个人比赛的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"MGA_OJ/Interface"
	TQ "MGA_OJ/Test-request"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"
	"MGA_OJ/util"
	"MGA_OJ/vo"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ICompetitionSingleController			定义了个人比赛类接口
type ICompetitionSingleController interface {
	Interface.RecordInterface          // 包含代码提交记录相关功能
	Interface.EnterInterface           // 包含报名方法
	Interface.HackInterface            // 包含hack方法
	CompetitionScore(ctx *gin.Context) // 对比赛结果进行分数统计
}

// CompetitionSingleController			定义了个人比赛工具类
type CompetitionSingleController struct {
	DB       *gorm.DB            // 含有一个数据库指针
	Redis    *redis.Client       // 含有一个redis指针
	UpGrader *websocket.Upgrader // 用于持久化连接
	RabbitMq *common.RabbitMQ    // 一个消息队列的指针
}

// @title    Enter
// @description   报名一篇比赛
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionSingleController) Enter(ctx *gin.Context) {
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

	// TODO 查看比赛是否已经进行过
	if time.Now().After(time.Time(competition.EndTime)) {
		response.Fail(ctx, nil, "比赛已经结束")
		return
	}

	// TODO 查看比赛是否为单人比赛
	if competition.Type != "Single" {
		response.Fail(ctx, nil, "并非单人比赛，无法报名")
		return
	}

	var competitionRank model.CompetitionRank

	// TODO 查看是否已经报名
	if _, err := c.Redis.ZScore(ctx, "CompetitionR"+id, user.ID.String()).Result(); err == nil {
		response.Fail(ctx, nil, "已报名")
		return
	}
	if c.DB.Where("member_id = (?) and competition_id = (?)", user.ID, competition.ID).First(&competitionRank).Error == nil {
		response.Fail(ctx, nil, "已报名")
		// TODO 加入redis
		c.Redis.ZAdd(ctx, "CompetitionR"+id, redis.Z{Member: user.ID.String(), Score: 0})
		return
	}

	// TODO 查看比赛是否需要密码
	var passwd model.Passwd
	if c.DB.Where("id = (?)", competition.PasswdId).First(&passwd).Error == nil {
		var pass model.Passwd
		// TODO 数据验证
		if err := ctx.ShouldBind(&pass); err != nil {
			log.Print(err.Error())
			response.Fail(ctx, nil, "数据验证错误")
			return
		}
		// TODO 判断密码是否正确
		if err := bcrypt.CompareHashAndPassword([]byte(passwd.Password), []byte(pass.Password)); err != nil {
			response.Fail(ctx, nil, "密码错误")
			return
		}
	}

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
	response.Success(ctx, gin.H{"competitionRank": competitionRank}, "报名成功")
	return
}

// @title    EnterCondition
// @description   查看报名状态
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionSingleController) EnterCondition(ctx *gin.Context) {
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
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionSingleController) CancelEnter(ctx *gin.Context) {
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
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionSingleController) EnterPage(ctx *gin.Context) {

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

	var competitionRanks []model.CompetitionRank

	// TODO 查找所有分页中可见的条目
	c.DB.Where("competition_id = (?)", competition.ID).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&competitionRanks)

	var total int64
	c.DB.Where("competition_id = (?)", competition.ID).Model(model.CompetitionRank{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"competitionRanks": competitionRanks, "total": total}, "成功")
}

// @title    Submit
// @description   用户进行提交操作
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionSingleController) Submit(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取比赛id
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
	if c.DB.Where("id = (?)", id).First(&competition).Error != nil {
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
	if !time.Now().After(time.Time(competition.StartTime)) || time.Now().After(time.Time(competition.EndTime)) {
		response.Fail(ctx, nil, "比赛不在进行中")
		return
	}

	var competitionRank model.CompetitionRank
	// TODO 查看是否已经报名
	// TODO 先看redis中是否存在
	if _, err := c.Redis.ZScore(ctx, "CompetitionR"+id, user.ID.String()).Result(); err != nil {
		if c.DB.Where("member_id = (?) and competition_id = (?)", user.ID, competition.ID).First(&competitionRank).Error != nil {
			response.Success(ctx, nil, "未报名")
			return
		}
		// TODO 加入redis
		c.Redis.ZAdd(ctx, "CompetitionR"+id, redis.Z{Member: user.ID.String(), Score: 0})
	}

	var requestRecord vo.RecordRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestRecord); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 查看当前problem状态
	var problem model.ProblemNew

	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "ProblemNew", fmt.Sprint(requestRecord.ProblemId)).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "ProblemNew", fmt.Sprint(requestRecord.ProblemId)).Result()
		if json.Unmarshal([]byte(cate), &problem) == nil {
			// TODO 跳过数据库搜寻problem过程
			goto leap
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "ProblemNew", fmt.Sprint(requestRecord.ProblemId))
		}
	}

	// TODO 查看题目是否在数据库中存在
	if c.DB.Where("id = (?)", requestRecord.ProblemId).First(&problem).Error != nil {
		response.Fail(ctx, nil, "题目不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(problem)
		c.Redis.HSet(ctx, "ProblemNew", fmt.Sprint(requestRecord.ProblemId), v)
	}

leap:

	// TODO 查看改题目是否为该比赛题目
	if problem.CompetitionId != competition.ID {
		response.Fail(ctx, nil, "题目不存在")
		return
	}

	// TODO 创建提交
	record := model.RecordCompetition{
		UserId:        user.ID,
		ProblemId:     requestRecord.ProblemId,
		Language:      requestRecord.Language,
		Code:          requestRecord.Code,
		Condition:     "Waiting",
		Pass:          0,
		CompetitionId: competition.ID,
	}

	// TODO 插入数据
	if err := c.DB.Create(&record).Error; err != nil {
		response.Fail(ctx, nil, "提交上传出错，数据验证有误")
		return
	}

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
		response.Fail(ctx, nil, "消息队列出错")
		return
	}
	// TODO 发布订阅用于提交列表
	recordList := vo.RecordList{
		RecordId: record.ID,
	}
	// TODO 将recordlist打包
	v, _ = json.Marshal(recordList)
	c.Redis.Publish(ctx, "RecordCompetitionChan", v)

	// TODO 成功
	response.Success(ctx, nil, "提交成功")
}

// @title    ShowRecord
// @description   查看一篇提交的内容
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionSingleController) ShowRecord(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var record model.RecordCompetition

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "RecordCompetition", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "RecordCompetition", id).Result()
		if json.Unmarshal([]byte(cate), &record) == nil {
			// TODO 跳过数据库搜寻过程
			goto leep
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "RecordCompetition", id)
		}
	}

	// TODO 查看提交是否在数据库中存在
	if c.DB.Where("id = (?)", id).First(&record).Error != nil {
		response.Fail(ctx, nil, "提交不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(record)
		c.Redis.HSet(ctx, "RecordCompetition", id, v)
	}
leep:
	// TODO 查找比赛
	var competition model.Competition

	// TODO 查看比赛是否存在
	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Competition", record.CompetitionId.String()).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Competition", record.CompetitionId.String()).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			goto leap
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Competition", record.CompetitionId.String())
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if c.DB.Where("id = (?)", record.CompetitionId.String()).First(&competition).Error != nil {
		response.Fail(ctx, nil, "提交对应比赛不存在")
		return
	}
	{
		// TODO 将比赛存入redis供下次使用
		v, _ := json.Marshal(competition)
		c.Redis.HSet(ctx, "Competition", record.CompetitionId.String(), v)
	}
leap:
	// TODO 如果不是自己的提交且比赛正在进行
	if record.UserId != user.ID && time.Now().After(time.Time(competition.StartTime)) && time.Now().Before(time.Time(competition.EndTime)) {
		record.Code = ""
	}

	response.Success(ctx, gin.H{"record": record}, "成功")
}

// @title    SearchList
// @description   获取多篇提交
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionSingleController) SearchList(ctx *gin.Context) {

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取查询条件
	Language := ctx.DefaultQuery("language", "")
	UserId := ctx.DefaultQuery("user_id", "")
	ProblemId := ctx.DefaultQuery("problem_id", "")
	StartTime := ctx.DefaultQuery("start_time", "")
	EndTime := ctx.DefaultQuery("end_time", "")
	Condition := ctx.DefaultQuery("condition", "")
	PassLow := ctx.DefaultQuery("pass_low", "")
	PassTop := ctx.DefaultQuery("pass_top", "")
	Hack := ctx.DefaultQuery("hack", "")

	db := common.GetDB()

	db = db.Where("`competition_id` = (?)", id)

	// TODO 根据参数设置where条件
	if Language != "" {
		db = db.Where("`language` = (?)", Language)
	}
	if UserId != "" {
		db = db.Where("`user_id` = (?)", UserId)
	}
	if ProblemId != "" {
		db = db.Where("`problem_id` = (?)", ProblemId)
	}
	if StartTime != "" {
		db = db.Where("`created_at` >= (?)", StartTime)
	}
	if EndTime != "" {
		db = db.Where("`created_at` <= (?)", EndTime)
	}
	if Condition != "" {
		db = db.Where("`condition` = (?)", Condition)
	}
	if PassLow != "" {
		db = db.Where("`pass` >= (?)", PassLow)
	}
	if PassTop != "" {
		db = db.Where("`pass` <= (?)", PassTop)
	}
	if Hack != "" {
		db = db.Where("`hack_id` != (?)", uuid.UUID{})
	}

	// TODO 分页
	var records []model.RecordCompetition

	var total int64

	// TODO 查找所有分页中可见的条目
	db.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&records)

	db.Model(model.Record{}).Count(&total)

	// TODO 查找比赛
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
	if c.DB.Where("id = (?)", id).First(&competition).Error != nil {
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
	if time.Now().After(time.Time(competition.StartTime)) && time.Now().Before(time.Time(competition.EndTime)) {
		for i := range records {
			if records[i].UserId != user.ID {
				records[i].Code = ""
			}
		}
	}

	// TODO 返回数据
	response.Success(ctx, gin.H{"records": records, "total": total}, "成功")
}

// @title    PublishPageList
// @description  订阅提交列表
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionSingleController) PublishPageList(ctx *gin.Context) {

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 查找比赛
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
	if c.DB.Where("id = (?)", id).First(&competition).Error != nil {
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
	if time.Now().Before(time.Time(competition.StartTime)) {
		if competition.HackTime.After(competition.EndTime) {
			if time.Now().After(time.Time(competition.HackTime)) {
				response.Fail(ctx, nil, "比赛未开始")
				return
			}
		} else if time.Now().After(time.Time(competition.EndTime)) {
			response.Fail(ctx, nil, "比赛未开始")
			return
		}

	}
	// TODO 订阅消息
	pubSub := c.Redis.Subscribe(ctx, "RecordCompetitionChan"+id)
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
		var recordList vo.RecordList
		json.Unmarshal([]byte(msg.Payload), &recordList)
		// TODO 写入ws数据
		if err := ws.WriteJSON(recordList); err != nil {
			break
		}
	}
}

// @title    CaseList
// @description   查看一篇提交的测试通过情况
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionSingleController) CaseList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var record model.RecordCompetition

	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "RecordCompetition", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "RecordCompetition", id).Result()
		if json.Unmarshal([]byte(cate), &record) == nil {
			// TODO 跳过数据库搜寻过程
			goto leep
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "RecordCompetition", id)
		}
	}

	// TODO 查看提交是否在数据库中存在
	if c.DB.Where("id = (?)", id).First(&record).Error != nil {
		response.Fail(ctx, nil, "提交不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(record)
		c.Redis.HSet(ctx, "RecordCompetition", id, v)
	}
leep:
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var cases []model.CaseCondition

	var total int64

	// TODO 查找所有分页中可见的条目
	c.DB.Where("record_id = (?)", record.ID).Order("id asc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&cases)

	c.DB.Where("record_id = (?)", record.ID).Model(model.CaseCondition{}).Count(&total)

	response.Success(ctx, gin.H{"cases": cases}, "成功")
}

// @title    Case
// @description   查看一篇测试通过情况
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionSingleController) Case(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var cas model.CaseCondition

	// TODO 查找所有分页中可见的条目
	if c.DB.Where("id = (?)", id).First(&cas).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}

	response.Success(ctx, gin.H{"case": cas}, "成功")
}

// @title    Hack
// @description   Hack比赛功能
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionSingleController) Hack(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var hackRequest vo.HackRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&hackRequest); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	var record model.Record

	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "RecordCompetition", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "RecordCompetition", id).Result()
		if json.Unmarshal([]byte(cate), &record) == nil {
			// TODO 跳过数据库搜寻过程
			goto leep
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "RecordCompetition", id)
		}
	}

	// TODO 查看提交是否在数据库中存在
	if c.DB.Where("id = (?)", id).First(&record).Error != nil {
		response.Fail(ctx, nil, "提交不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(record)
		c.Redis.HSet(ctx, "RecordCompetition", id, v)
	}
leep:

	if record.Condition != "Accepted" {
		response.Fail(ctx, nil, "提交未通过")
		return
	}

	if (record.HackId != uuid.UUID{}) {
		response.Fail(ctx, nil, "已经被hack了")
		return
	}

	// TODO 查看题目
	// TODO 先看redis中是否存在
	var problem model.ProblemNew
	if ok, _ := c.Redis.HExists(ctx, "ProblemNew", record.ProblemId.String()).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "ProblemNew", record.ProblemId.String()).Result()
		if json.Unmarshal([]byte(cate), &problem) == nil {
			// TODO 跳过数据库搜寻problem过程
			goto leap
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "ProblemNew", record.ProblemId.String())
		}
	}

	// TODO 查看题目是否在数据库中存在
	if c.DB.Where("id = (?)", record.ProblemId.String()).First(&problem).Error != nil {
		response.Fail(ctx, nil, "题目不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(problem)
		c.Redis.HSet(ctx, "ProblemNew", record.ProblemId.String(), v)
	}

leap:
	// TODO 查看题目对应比赛是否存在
	var competition model.Competition

	// TODO 查看比赛是否存在
	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Competition", problem.CompetitionId.String()).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Competition", problem.CompetitionId.String()).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			goto comp
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Competition", problem.CompetitionId.String())
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if c.DB.Where("id = (?)", problem.CompetitionId.String()).First(&competition).Error != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
	{
		// TODO 将比赛存入redis供下次使用
		v, _ := json.Marshal(competition)
		c.Redis.HSet(ctx, "Competition", problem.CompetitionId.String(), v)
	}
comp:
	if competition.Type != "Single" {
		response.Fail(ctx, nil, "非个人比赛")
		return
	}
	// TODO 查看比赛是否在hack时间
	if !time.Now().After(time.Time(competition.EndTime)) || !time.Now().Before(time.Time(competition.HackTime)) {
		response.Fail(ctx, nil, "不在hack时间")
		return
	}

	var competitionRank model.CompetitionRank
	// TODO 查看是否已经报名
	// TODO 先看redis中是否存在
	if _, err := c.Redis.ZScore(ctx, "CompetitionR"+id, user.ID.String()).Result(); err != nil {
		if c.DB.Where("member_id = (?) and competition_id = (?)", user.ID, competition.ID).First(&competitionRank).Error != nil {
			response.Success(ctx, nil, "未报名")
			return
		}
		// TODO 加入redis
		c.Redis.ZAdd(ctx, "CompetitionR"+id, redis.Z{Member: user.ID.String(), Score: 0})
	}

	// TODO 查看题目是否有输入测试程序
	var inputCheckProgram model.Program

	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Program", problem.InputCheck.String()).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Program", problem.InputCheck.String()).Result()
		if json.Unmarshal([]byte(cate), &inputCheckProgram) == nil {
			// TODO 跳过数据库搜寻program过程
			goto inputCheck
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Program", problem.InputCheck.String())
		}
	}

	// TODO 查看程序是否在数据库中存在
	if c.DB.Where("id = (?)", problem.InputCheck.String()).First(&inputCheckProgram).Error != nil {
		response.Fail(ctx, nil, "输入检查程序不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(inputCheckProgram)
		c.Redis.HSet(ctx, "Program", problem.InputCheck.String(), v)
	}
	// TODO 查看是否通过输入检查程序
	if condition, output := TQ.JudgeRun(inputCheckProgram.Language, inputCheckProgram.Code, hackRequest.Input, problem.MemoryLimit*2, problem.TimeLimit*2); condition != "ok" || output != "ok" {
		response.Fail(ctx, nil, "输入检查程序未通过")
		return
	}
inputCheck:
	// TODO 查看题目是否有标准程序
	var standardProgram model.Program

	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Program", problem.Standard.String()).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Program", problem.Standard.String()).Result()
		if json.Unmarshal([]byte(cate), &standardProgram) == nil {
			// TODO 跳过数据库搜寻program过程
			goto special
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Program", problem.Standard.String())
		}
	}

	// TODO 查看程序是否在数据库中存在
	if c.DB.Where("id = (?)", problem.Standard.String()).First(&standardProgram).Error != nil {
		response.Fail(ctx, nil, "标准程序不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(standardProgram)
		c.Redis.HSet(ctx, "Program", problem.Standard.String(), v)
	}
special:
	// TODO 查看是否通过标准程序
	var hackoutput, condition, recordoutput string
	if condition, hackoutput = TQ.JudgeRun(standardProgram.Language, standardProgram.Code, hackRequest.Input, problem.MemoryLimit, problem.TimeLimit); condition != "ok" {
		response.Fail(ctx, nil, "未通过标准程序")
		return
	}
	if condition, recordoutput = TQ.JudgeRun(record.Language, record.Code, hackRequest.Input, problem.MemoryLimit, problem.TimeLimit); condition != "ok" {
		goto success
	}

	{
		// TODO 查看题目是否有特判程序
		var specialJudgeProgram model.Program

		// TODO 先看redis中是否存在
		if ok, _ := c.Redis.HExists(ctx, "Program", problem.SpecialJudge.String()).Result(); ok {
			cate, _ := c.Redis.HGet(ctx, "Program", problem.SpecialJudge.String()).Result()
			if json.Unmarshal([]byte(cate), &specialJudgeProgram) == nil {
				// TODO 跳过数据库搜寻program过程
				goto standard
			} else {
				// TODO 移除损坏数据
				c.Redis.HDel(ctx, "Program", problem.SpecialJudge.String())
			}
		}

		// TODO 查看程序是否在数据库中存在
		if c.DB.Where("id = (?)", problem.SpecialJudge.String()).First(&specialJudgeProgram).Error != nil {
			if recordoutput != hackoutput {
				goto success
			}
			response.Fail(ctx, gin.H{"output": hackoutput}, "与标准程序输出一致")
			return
		}
		// TODO 将题目存入redis供下次使用
		{
			v, _ := json.Marshal(specialJudgeProgram)
			c.Redis.HSet(ctx, "Program", problem.SpecialJudge.String(), v)
		}
	standard:
		if condition, hackoutput = TQ.JudgeRun(specialJudgeProgram.Language, specialJudgeProgram.Code, hackRequest.Input+"\n"+hackoutput, problem.MemoryLimit*3, problem.TimeLimit*3); condition != "ok" || hackoutput != "ok" {
			response.Fail(ctx, nil, "输入未通过特殊裁判")
			return
		}
		if condition, recordoutput = TQ.JudgeRun(specialJudgeProgram.Language, specialJudgeProgram.Code, hackRequest.Input+"\n"+recordoutput, problem.MemoryLimit*3, problem.TimeLimit*3); condition == "ok" && recordoutput == "ok" {
			response.Fail(ctx, nil, "通过了特殊裁判")
			return
		}
	}

success:
	// TODO 成功hack
	hack := model.Hack{
		UserId:   user.ID,
		Type:     competition.Type,
		Input:    hackRequest.Input,
		RecordId: record.ID,
	}

	// TODO 插入数据
	if err := c.DB.Create(&hack).Error; err != nil {
		response.Fail(ctx, nil, "记录上传出错，数据验证有误")
		return
	}

	record.HackId = hack.ID
	c.DB.Save(&record)

	// TODO 分数提升
	// TODO 查看用户hack数量
	var hackNum model.HackNum

	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "HackNum", competition.ID.String()+user.ID.String()).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "HackNum", competition.ID.String()+user.ID.String()).Result()
		if json.Unmarshal([]byte(cate), &hackNum) == nil {
			// TODO 跳过数据库搜寻hackNum过程
			goto hacknum
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "HackNum", competition.ID.String()+user.ID.String())
		}
	}

	// TODO 查看hackNum是否在数据库中存在
	if c.DB.Where("member_id = (?) and competition_id = (?)", user.ID, competition.ID).First(&hackNum).Error != nil {
		hackNum = model.HackNum{
			MemberId:      user.ID,
			CompetitionId: competition.ID,
			Num:           0,
		}
		c.DB.Create(&hackNum)
	}
hacknum:
	// TODO 分数变化
	hackNum.Num++
	if hackNum.Num <= competition.HackNum {
		hackNum.Score += competition.HackScore
		c.Redis.ZIncrBy(ctx, "CompetitionR"+competition.ID.String(), float64(competition.HackScore), user.ID.String())
		// TODO 发布订阅用于滚榜
		rankList := vo.RankList{
			MemberId: user.ID,
		}
		// TODO 将ranklist打包
		v, _ := json.Marshal(rankList)
		c.Redis.Publish(ctx, "CompetitionChan", v)
	}
	c.Redis.ZIncrBy(ctx, "CompetitionR"+competition.ID.String(), -float64(competition.HackScore), record.UserId.String())
	// TODO 发布订阅用于滚榜
	rankList := vo.RankList{
		MemberId: record.UserId,
	}
	// TODO 将ranklist打包
	v, _ := json.Marshal(rankList)
	c.Redis.Publish(ctx, "CompetitionChan", v)
	// TODO 将hackNum存入数据库
	c.DB.Save(&hackNum)
	// TODO 移除损坏数据
	c.Redis.HDel(ctx, "HackNum", competition.ID.String()+user.ID.String())

	response.Success(ctx, gin.H{"hack": hack}, "成功")
}

// @title    CompetitionScore
// @description   对比赛分数计算
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionSingleController) CompetitionScore(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 取出比赛id
	id := ctx.Params.ByName("id")

	var competition model.Competition
	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Competition", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Competition", id).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			goto leap
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Competition", id)
		}
	}
	// TODO 在数据库中查找
	if c.DB.Where("id = (?)", id).First(&competition).Error != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
leap:
	// TODO 取出用户权限
	if user.Level < 4 && competition.UserId != user.ID {
		response.Fail(ctx, nil, "用户权限不足")
		return
	}

	// TODO 整理比赛结果
	var competitionRankrs []model.CompetitionRank

	c.DB.Where("competition_id = (?)", id).Order("score desc penalties asc").Find(&competitionRankrs)

	// TODO 用户分数总和
	var sum float64

	// TODO 记录分数
	var scores []float64

	// TODO 处理用户分数
	// TODO 用户字典
	userMap := make(map[uuid.UUID]model.User)

	// TODO 按序取出所有用户
	for i := range competitionRankrs {
		id := competitionRankrs[i].MemberId
		var user model.User

		// TODO 先看redis中是否存在
		if ok, _ := c.Redis.HExists(ctx, "User", id.String()).Result(); ok {
			cate, _ := c.Redis.HGet(ctx, "User", id.String()).Result()
			if json.Unmarshal([]byte(cate), &user) == nil {
				goto leep
			} else {
				// TODO 移除损坏数据
				c.Redis.HDel(ctx, "User", id.String())
			}
		}

		// TODO 查看用户是否在数据库中存在
		if c.DB.Where("id = (?)", id).First(&user).Error != nil {
			continue
		}
	leep:
		// TODO 存入字典
		userMap[user.ID] = user
		// TODO 统计分数
		scores = append(scores, user.Score)
		sum += user.Score
	}

	// TODO 将用户按原预期排名排序
	sort.Sort(sort.Float64Slice(scores))

	// TODO 遍历比赛结果，计算每个用户的预期排名差
	for i := range competitionRankrs {
		id := competitionRankrs[i].MemberId
		// TODO 二分查找实际排名
		j := sort.Search(len(scores), func(i int) bool {
			return scores[i] <= userMap[id].Score
		})
		// TODO 计算该用户的期望排名差
		del := j - i
		// TODO 查看该用户的参赛次数
		var fre int64
		c.DB.Where("user_id = (?)", id).Model(model.UserScoreChange{}).Count(&fre)
		// TODO 查看本次比赛人数
		total := len(scores)
		// TODO 带入公式计算分数变化
		scoreChange := util.ScoreChange(float64(fre), sum, float64(del), float64(total))

		// TODO 将分数变化存入数据库
		userScoreChange := model.UserScoreChange{
			ScoreChange:   scoreChange,
			CompetitionId: competition.ID,
			UserId:        id,
			Type:          competition.Type,
		}
		c.DB.Create(&userScoreChange)

		// TODO 将用户信息更新存入数据库
		var user model.User
		user = userMap[id]
		user.Score += scoreChange
		c.DB.Save(&user)
	}

	response.Success(ctx, nil, "计算完成")
}

// @title    NewCompetitionController
// @description   新建一个ICompetitionSingleController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   ICompetitionSingleController		返回一个ICompetitionSingleController用于调用各种函数
func NewCompetitionSingleController() ICompetitionSingleController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	rabbitmq := common.GetRabbitMq()
	upGrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	InitCompetition["Single"] = initSingleCompetition
	FinishCompetition["Single"] = finishSingleCompetition

	return CompetitionSingleController{DB: db, Redis: redis, UpGrader: upGrader, RabbitMq: rabbitmq}
}

func initSingleCompetition(ctx context.Context, redis *redis.Client, db *gorm.DB, competition model.Competition) {
	if competition.Type != "Single" {
		log.Println("single competition's type is error!")
	} else {
		log.Println("single competition start!", competition)
	}
}

func finishSingleCompetition(ctx context.Context, redis *redis.Client, db *gorm.DB, competition model.Competition) {
	if competition.Type != "Single" {
		log.Println("single competition's type is error!")
	} else {
		log.Println("single competition finish!", competition)
	}
	CompetitionFinish(ctx, redis, db, competition)
}
