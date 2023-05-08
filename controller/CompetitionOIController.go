// @Title  CompetitionOIController
// @Description  该文件提供关于操作个人比赛的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"MGA_OJ/Interface"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"
	"MGA_OJ/vo"
	"context"
	"encoding/json"
	"fmt"
	"log"
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

// ICompetitionOIController			定义了个人比赛类接口
type ICompetitionOIController interface {
	Interface.RecordInterface // 包含代码提交记录相关功能
	Interface.EnterInterface  // 包含报名方法
}

// CompetitionOIController			定义了个人比赛工具类
type CompetitionOIController struct {
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
func (c CompetitionOIController) Enter(ctx *gin.Context) {
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
	if time.Now().After(time.Time(competition.EndTime)) {
		response.Fail(ctx, nil, "比赛已经结束")
		return
	}

	// TODO 查看比赛是否为OI比赛
	if competition.Type != "OI" {
		response.Fail(ctx, nil, "并非OI比赛，无法报名")
		return
	}

	var competitionRank model.CompetitionRank

	// TODO 查看是否已经报名
	if _, err := c.Redis.ZScore(ctx, "CompetitionR"+id, user.ID.String()).Result(); err == nil {
		response.Fail(ctx, nil, "已报名")
		return
	}
	if c.DB.Where("member_id = ? and competition_id = ?", user.ID, competition.ID).First(&competitionRank).Error == nil {
		response.Fail(ctx, nil, "已报名")
		// TODO 加入redis
		c.Redis.ZAdd(ctx, "CompetitionR"+id, redis.Z{Member: user.ID.String(), Score: 0})
		return
	}

	// TODO 查看比赛是否需要密码
	var passwd model.Passwd
	if c.DB.Where("id = ?", competition.PasswdId).First(&passwd).Error == nil {
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
func (c CompetitionOIController) EnterCondition(ctx *gin.Context) {
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
	// TODO 检查比赛类型
	if competition.Type != "OI" {
		response.Fail(ctx, nil, "不是OI比赛，无法报名")
		return
	}

	var competitionRank model.CompetitionRank

	// TODO 查看是否已经报名
	if _, err := c.Redis.ZScore(ctx, "CompetitionR"+id, user.ID.String()).Result(); err == nil {
		response.Success(ctx, gin.H{"enter": true}, "已报名")
		return
	}
	if c.DB.Where("member_id = ? and competition_id = ?", user.ID, competition.ID).First(&competitionRank).Error == nil {
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
func (c CompetitionOIController) CancelEnter(ctx *gin.Context) {
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
	if competition.Type != "OI" {
		response.Fail(ctx, nil, "不是OI比赛，无法报名")
		return
	}

	// TODO 查看比赛是否已经进行过
	if time.Now().After(time.Time(competition.StartTime)) {
		response.Fail(ctx, nil, "比赛已经开始")
		return
	}

	var competitionRank model.CompetitionRank

	// TODO 查看是否已经报名
	if c.DB.Where("member_id = ? and competition_id = ?", user.ID, competition.ID).First(&competitionRank).Error != nil {
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
func (c CompetitionOIController) EnterPage(ctx *gin.Context) {

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

	var competitionRanks []model.CompetitionRank

	// TODO 查找所有分页中可见的条目
	c.DB.Where("competition_id = ?", competition.ID).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&competitionRanks)

	var total int64
	c.DB.Where("competition_id = ?", competition.ID).Model(model.CompetitionRank{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"competitionRanks": competitionRanks, "total": total}, "成功")
}

// @title    Submit
// @description   用户进行提交操作
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionOIController) Submit(ctx *gin.Context) {

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
	if !time.Now().After(time.Time(competition.StartTime)) || time.Now().After(time.Time(competition.EndTime)) {
		response.Fail(ctx, nil, "比赛不在进行中")
		return
	}

	var competitionRank model.CompetitionRank
	// TODO 查看是否已经报名
	// TODO 先看redis中是否存在
	if _, err := c.Redis.ZScore(ctx, "CompetitionR"+id, user.ID.String()).Result(); err != nil {
		if c.DB.Where("member_id = ? and competition_id = ?", user.ID, competition.ID).First(&competitionRank).Error != nil {
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
	if c.DB.Where("id = ?", requestRecord.ProblemId).First(&problem).Error != nil {
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

	// TODO 插入数据，但不提交到判题机
	if err := c.DB.Create(&record).Error; err != nil {
		response.Fail(ctx, nil, "提交上传出错，数据验证有误")
		return
	}

	{
		// TODO 将提交存入redis供判题机使用
		v, _ := json.Marshal(record)
		c.Redis.HSet(ctx, "RecordCompetition", fmt.Sprint(record.ID), v)
	}

	// TODO 成功
	response.Success(ctx, nil, "提交成功")
}

// @title    ShowRecord
// @description   查看一篇提交的内容
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionOIController) ShowRecord(ctx *gin.Context) {
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
	if c.DB.Where("id = ?", id).First(&record).Error != nil {
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
	if c.DB.Where("id = ?", record.CompetitionId.String()).First(&competition).Error != nil {
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
func (c CompetitionOIController) SearchList(ctx *gin.Context) {

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

	db = db.Where("competition_id = ?", id)

	// TODO 根据参数设置where条件
	if Language != "" {
		db = db.Where("language = ?", Language)
	}
	if UserId != "" {
		db = db.Where("user_id = ?", UserId)
	}
	if ProblemId != "" {
		db = db.Where("problem_id = ?", ProblemId)
	}
	if StartTime != "" {
		db = db.Where("created_at >= ?", StartTime)
	}
	if EndTime != "" {
		db = db.Where("created_at <= ?", EndTime)
	}
	if Condition != "" {
		db = db.Where("condition = ?", Condition)
	}
	if PassLow != "" {
		db = db.Where("pass >= ?", PassLow)
	}
	if PassTop != "" {
		db = db.Where("pass <= ?", PassTop)
	}
	if Hack != "" {
		db = db.Where("hack_id != ?", uuid.UUID{})
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
func (c CompetitionOIController) PublishPageList(ctx *gin.Context) {

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
		// TODO 读取ws中的数据
		_, _, err := ws.ReadMessage()
		// TODO 断开连接
		if err != nil {
			break
		}
		var recordList vo.RecordList
		json.Unmarshal([]byte(msg.Payload), &recordList)
		// TODO 写入ws数据
		ws.WriteJSON(recordList)
	}
}

// @title    CaseList
// @description   查看一篇提交的测试通过情况
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionOIController) CaseList(ctx *gin.Context) {
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
	if c.DB.Where("id = ?", id).First(&record).Error != nil {
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
	c.DB.Where("record_id = ?", record.ID).Order("id asc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&cases)

	c.DB.Where("record_id = ?", record.ID).Model(model.CaseCondition{}).Count(&total)

	response.Success(ctx, gin.H{"cases": cases}, "成功")
}

// @title    Case
// @description   查看一篇测试通过情况
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionOIController) Case(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var cas model.CaseCondition

	// TODO 查找所有分页中可见的条目
	if c.DB.Where("id = ?", id).First(&cas).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}

	response.Success(ctx, gin.H{"case": cas}, "成功")
}

// @title    NewCompetitionController
// @description   新建一个ICompetitionOIController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   ICompetitionOIController		返回一个ICompetitionOIController用于调用各种函数
func NewCompetitionOIController() ICompetitionOIController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	rabbitmq := common.GetRabbitMq()
	upGrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	InitCompetition["OI"] = initOICompetition
	FinishCompetition["OI"] = finishOICompetition

	return CompetitionOIController{DB: db, Redis: redis, UpGrader: upGrader, RabbitMq: rabbitmq}
}

func initOICompetition(ctx context.Context, redis *redis.Client, db *gorm.DB, competition model.Competition) {
	if competition.Type != "Single" {
		log.Println("single competition's type is error!")
	} else {
		log.Println("single competition start!", competition)
	}
}

func finishOICompetition(ctx context.Context, redis *redis.Client, db *gorm.DB, competition model.Competition) {
	if competition.Type != "Single" {
		log.Println("single competition's type is error!")
	} else {
		log.Println("single competition finish!", competition)
	}
	RabbitMq := common.GetRabbitMq()
	var records []model.RecordCompetition
	db.Where("competition_id = ?", competition.ID).Find(&records)
	for i := range records {
		// TODO 加入消息队列
		recordRabbit := vo.RecordRabbit{
			RecordId: records[i].ID,
			Type:     competition.Type,
		}
		v, _ := json.Marshal(recordRabbit)
		if err := RabbitMq.PublishSimple(string(v)); err != nil {
			log.Println(err)
			return
		}
		// TODO 发布订阅用于提交列表
		recordList := vo.RecordList{
			RecordId: records[i].ID,
		}
		// TODO 将recordlist打包
		v, _ = json.Marshal(recordList)
		redis.Publish(ctx, "RecordCompetitionChan", v)
	}
	CompetitionFinish(ctx, redis, db, competition)
}
