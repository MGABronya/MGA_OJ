// @Title  RecordController
// @Description  该文件提供关于操作提交的各种方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package controller

import (
	Handle "MGA_OJ/Behavior"
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
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// IRecordController			定义了提交类接口
type IRecordController interface {
	Interface.RecordInterface // 代码提交记录相关功能
	Interface.HackInterface   // 包含hack相关功能
}

// RecordController			定义了提交工具类
type RecordController struct {
	DB       *gorm.DB            // 含有一个数据库指针
	Redis    *redis.Client       // 含有一个redis指针
	Rabbitmq *common.RabbitMQ    // 含有一个消息中间件
	UpGrader *websocket.Upgrader // 用于持久化连接
}

// @title    Submit
// @description   用户进行提交操作
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RecordController) Submit(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)
	var requestRecord vo.RecordRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestRecord); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 查看当前problem状态
	var problem model.Problem

	// TODO 先看redis中是否存在
	if ok, _ := r.Redis.HExists(ctx, "Problem", fmt.Sprint(requestRecord.ProblemId)).Result(); ok {
		cate, _ := r.Redis.HGet(ctx, "Problem", fmt.Sprint(requestRecord.ProblemId)).Result()
		if json.Unmarshal([]byte(cate), &problem) == nil {
			// TODO 跳过数据库搜寻problem过程
			goto leep
		} else {
			// TODO 移除损坏数据
			r.Redis.HDel(ctx, "Problem", fmt.Sprint(requestRecord.ProblemId))
		}
	}

	// TODO 查看题目是否在数据库中存在
	if r.DB.Where("id = (?)", requestRecord.ProblemId).First(&problem).Error != nil {
		response.Fail(ctx, nil, "题目不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(problem)
		r.Redis.HSet(ctx, "Problem", fmt.Sprint(requestRecord.ProblemId), v)
	}

leep:

	// TODO 创建提交
	record := model.Record{
		UserId:    user.ID,
		ProblemId: requestRecord.ProblemId,
		Language:  requestRecord.Language,
		Code:      requestRecord.Code,
		Condition: "Waiting",
		Pass:      0,
	}

	// TODO 插入数据
	if err := r.DB.Create(&record).Error; err != nil {
		response.Fail(ctx, nil, "提交上传出错，数据验证有误")
		return
	}

	recordRaabbit := vo.RecordRabbit{
		RecordId: record.ID,
		Type:     "Normal",
	}
	{
		// TODO 将提交存入redis供判题机使用
		v, _ := json.Marshal(record)
		r.Redis.HSet(ctx, "Record", fmt.Sprint(record.ID), v)
	}

	v, _ := json.Marshal(recordRaabbit)
	// TODO 查看提交题目能否被解析为外站题目
	if _, _, err := util.DeCodeUUID(record.ProblemId); err != nil {
		// TODO 加入消息队列用于本地消费
		if SubmitCheck() {
			// TODO 如果需要转发到云端，交由转发机处理
			r.Redis.Publish(ctx, "Cloud", v)
		} else if err := r.Rabbitmq.PublishSimple(string(v)); err != nil {
			response.Fail(ctx, nil, "消息队列出错")
			return
		}
	} else {
		// TODO 如果存在指定平台，交由转发机处理
		r.Redis.Publish(ctx, "Vjudge", v)
	}

	// TODO 发布订阅用于提交列表
	recordList := vo.RecordList{
		RecordId: record.ID,
	}
	// TODO 将recordlist打包
	v, _ = json.Marshal(recordList)
	r.Redis.Publish(ctx, "RecordChan", v)

	// TODO 成功
	response.Success(ctx, gin.H{"record": record}, "提交成功")
}

// @title    ShowRecord
// @description   查看一篇提交的内容
// @auth      MGAronya       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RecordController) ShowRecord(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var record model.Record

	// TODO 先看redis中是否存在
	if ok, _ := r.Redis.HExists(ctx, "Record", id).Result(); ok {
		cate, _ := r.Redis.HGet(ctx, "Record", id).Result()
		if json.Unmarshal([]byte(cate), &record) == nil {
			// TODO 跳过数据库搜寻过程
			goto leep
		} else {
			// TODO 移除损坏数据
			r.Redis.HDel(ctx, "Record", id)
		}
	}

	// TODO 查看提交是否在数据库中存在
	if r.DB.Where("id = (?)", id).First(&record).Error != nil {
		response.Fail(ctx, nil, "提交不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(record)
		r.Redis.HSet(ctx, "Record", id, v)
	}
leep:

	response.Success(ctx, gin.H{"record": record}, "成功")
}

// @title    SearchList
// @description   获取多篇提交
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RecordController) SearchList(ctx *gin.Context) {

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

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
	var records []model.Record

	var total int64

	// TODO 查找所有分页中可见的条目
	db.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&records)

	db.Model(model.Record{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"records": records, "total": total}, "成功")
}

// @title    PublishPageList
// @description  订阅提交列表
// @auth      MGAronya       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RecordController) PublishPageList(ctx *gin.Context) {

	// TODO 订阅消息
	pubSub := r.Redis.Subscribe(ctx, "RecordChan")
	defer pubSub.Close()
	// TODO 获得消息管道
	ch := pubSub.Channel()

	// TODO 升级get请求为webSocket协议
	ws, err := r.UpGrader.Upgrade(ctx.Writer, ctx.Request, nil)
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

// @title    Publish
// @description  订阅提交
// @auth      MGAronya       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RecordController) Publish(ctx *gin.Context) {

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 订阅消息
	pubSub := r.Redis.Subscribe(ctx, "RecordChan"+id)
	defer pubSub.Close()
	// TODO 获得消息管道
	ch := pubSub.Channel()

	// TODO 升级get请求为webSocket协议
	ws, err := r.UpGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	// TODO 监听消息
	for msg := range ch {
		var recordCase vo.RecordCase
		json.Unmarshal([]byte(msg.Payload), &recordCase)
		// TODO 写入ws数据
		if err := ws.WriteJSON(recordCase); err != nil {
			break
		}
	}
}

// @title    CaseList
// @description   查看一篇提交的测试通过情况
// @auth      MGAronya       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RecordController) CaseList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var record model.Record

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 先看redis中是否存在
	if ok, _ := r.Redis.HExists(ctx, "Record", id).Result(); ok {
		cate, _ := r.Redis.HGet(ctx, "Record", id).Result()
		if json.Unmarshal([]byte(cate), &record) == nil {
			// TODO 跳过数据库搜寻过程
			goto leep
		} else {
			// TODO 移除损坏数据
			r.Redis.HDel(ctx, "Record", id)
		}
	}

	// TODO 查看提交是否在数据库中存在
	if r.DB.Where("id = (?)", id).First(&record).Error != nil {
		response.Fail(ctx, nil, "提交不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(record)
		r.Redis.HSet(ctx, "Record", id, v)
	}
leep:
	if user.ID != record.UserId {
		response.Fail(ctx, nil, "非提交者无法查看")
		return
	}
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var cases []model.CaseCondition

	var total int64

	// TODO 查找所有分页中可见的条目
	r.DB.Where("record_id = (?)", record.ID).Order("c_id desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&cases)

	r.DB.Where("record_id = (?)", record.ID).Model(model.CaseCondition{}).Count(&total)

	response.Success(ctx, gin.H{"cases": cases}, "成功")
}

// @title    Case
// @description   查看一篇测试通过情况
// @auth      MGAronya       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RecordController) Case(ctx *gin.Context) {
	// TODO 获取path中的id
	cid := ctx.Params.ByName("cid")
	id := ctx.Params.ByName("id")
	var cas model.CaseCondition

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var record model.Record

	// TODO 先看redis中是否存在
	if ok, _ := r.Redis.HExists(ctx, "Record", id).Result(); ok {
		cate, _ := r.Redis.HGet(ctx, "Record", id).Result()
		if json.Unmarshal([]byte(cate), &record) == nil {
			// TODO 跳过数据库搜寻过程
			goto leep
		} else {
			// TODO 移除损坏数据
			r.Redis.HDel(ctx, "Record", id)
		}
	}

	// TODO 查看提交是否在数据库中存在
	if r.DB.Where("id = (?)", id).First(&record).Error != nil {
		response.Fail(ctx, nil, "提交不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(record)
		r.Redis.HSet(ctx, "Record", id, v)
	}
leep:
	if user.ID != record.UserId {
		response.Fail(ctx, nil, "非提交者无法查看")
		return
	}

	// TODO 查找所有分页中可见的条目
	if r.DB.Where("record_id = (?) and c_id = (?)", id, cid).First(&cas).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}

	response.Success(ctx, gin.H{"case": cas}, "成功")
}

// @title    Hack
// @description   Hack功能
// @auth      MGAronya       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RecordController) Hack(ctx *gin.Context) {
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
	if ok, _ := r.Redis.HExists(ctx, "Record", id).Result(); ok {
		cate, _ := r.Redis.HGet(ctx, "Record", id).Result()
		if json.Unmarshal([]byte(cate), &record) == nil {
			// TODO 跳过数据库搜寻过程
			goto leep
		} else {
			// TODO 移除损坏数据
			r.Redis.HDel(ctx, "Record", id)
		}
	}

	// TODO 查看提交是否在数据库中存在
	if r.DB.Where("id = (?)", id).First(&record).Error != nil {
		response.Fail(ctx, nil, "提交不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(record)
		r.Redis.HSet(ctx, "Record", id, v)
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
	var problem model.Problem
	if ok, _ := r.Redis.HExists(ctx, "Problem", record.ProblemId.String()).Result(); ok {
		cate, _ := r.Redis.HGet(ctx, "Problem", record.ProblemId.String()).Result()
		if json.Unmarshal([]byte(cate), &problem) == nil {
			// TODO 跳过数据库搜寻problem过程
			goto leap
		} else {
			// TODO 移除损坏数据
			r.Redis.HDel(ctx, "Problem", record.ProblemId.String())
		}
	}

	// TODO 查看题目是否在数据库中存在
	if r.DB.Where("id = (?)", record.ProblemId.String()).First(&problem).Error != nil {
		response.Fail(ctx, nil, "题目不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(problem)
		r.Redis.HSet(ctx, "Problem", record.ProblemId.String(), v)
	}

leap:
	// TODO 查看题目是否有输入测试程序
	var inputCheckProgram model.Program

	// TODO 先看redis中是否存在
	if ok, _ := r.Redis.HExists(ctx, "Program", problem.InputCheck.String()).Result(); ok {
		cate, _ := r.Redis.HGet(ctx, "Program", problem.InputCheck.String()).Result()
		if json.Unmarshal([]byte(cate), &inputCheckProgram) == nil {
			// TODO 跳过数据库搜寻program过程
			goto inputCheck
		} else {
			// TODO 移除损坏数据
			r.Redis.HDel(ctx, "Program", problem.InputCheck.String())
		}
	}

	// TODO 查看程序是否在数据库中存在
	if r.DB.Where("id = (?)", problem.InputCheck.String()).First(&inputCheckProgram).Error != nil {
		response.Fail(ctx, nil, "输入检查程序不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(inputCheckProgram)
		r.Redis.HSet(ctx, "Program", problem.InputCheck.String(), v)
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
	if ok, _ := r.Redis.HExists(ctx, "Program", problem.Standard.String()).Result(); ok {
		cate, _ := r.Redis.HGet(ctx, "Program", problem.Standard.String()).Result()
		if json.Unmarshal([]byte(cate), &standardProgram) == nil {
			// TODO 跳过数据库搜寻program过程
			goto special
		} else {
			// TODO 移除损坏数据
			r.Redis.HDel(ctx, "Program", problem.Standard.String())
		}
	}

	// TODO 查看程序是否在数据库中存在
	if r.DB.Where("id = (?)", problem.Standard.String()).First(&standardProgram).Error != nil {
		response.Fail(ctx, nil, "标准程序不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(standardProgram)
		r.Redis.HSet(ctx, "Program", problem.Standard.String(), v)
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
		if ok, _ := r.Redis.HExists(ctx, "Program", problem.SpecialJudge.String()).Result(); ok {
			cate, _ := r.Redis.HGet(ctx, "Program", problem.SpecialJudge.String()).Result()
			if json.Unmarshal([]byte(cate), &specialJudgeProgram) == nil {
				// TODO 跳过数据库搜寻program过程
				goto standard
			} else {
				// TODO 移除损坏数据
				r.Redis.HDel(ctx, "Program", problem.SpecialJudge.String())
			}
		}

		// TODO 查看程序是否在数据库中存在
		if r.DB.Where("id = (?)", problem.SpecialJudge.String()).First(&specialJudgeProgram).Error != nil {
			if recordoutput != hackoutput {
				goto success
			}
			response.Fail(ctx, gin.H{"output": hackoutput}, "与标准程序输出一致")
			return
		}
		// TODO 将题目存入redis供下次使用
		{
			v, _ := json.Marshal(specialJudgeProgram)
			r.Redis.HSet(ctx, "Program", problem.SpecialJudge.String(), v)
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
		Type:     "Normal",
		Input:    hackRequest.Input,
		RecordId: record.ID,
	}

	// 用户hack行为统计
	Handle.Behaviors["Hack"].PublishBehavior(1, user.ID)

	// TODO 插入数据
	if err := r.DB.Create(&hack).Error; err != nil {
		response.Fail(ctx, nil, "记录上传出错，数据验证有误")
		return
	}

	record.HackId = hack.ID
	r.DB.Save(&record)

	response.Success(ctx, gin.H{"hack": hack}, "成功")
}

// @title    NewRecordController
// @description   新建一个IRecordController
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   IRecordController		返回一个IRecordController用于调用各种函数
func NewRecordController() IRecordController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	rabbitmq := common.GetRabbitMq()
	upGrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	db.AutoMigrate(model.Record{})
	db.AutoMigrate(model.CaseCondition{})
	return RecordController{DB: db, Redis: redis, Rabbitmq: rabbitmq, UpGrader: upGrader}
}

// @title    Transponder
// @description   转发器，用于转发提交
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   void
func Transponder() {
	redis := common.GetRedisClient(0)
	db := common.GetDB()
	ctx := context.Background()
	// TODO 订阅消息
	pubSub := redis.Subscribe(ctx, "Vjudge")
	defer pubSub.Close()
	// TODO 获得消息管道
	ch := pubSub.Channel()
	// TODO 监听消息
	for msg := range ch {
		var recordRaabbit vo.RecordRabbit
		json.Unmarshal([]byte(msg.Payload), &recordRaabbit)
		// TODO 尝试在redis中取出该提交
		var record model.Record
		// TODO 先看redis中是否存在
		if ok, _ := redis.HExists(ctx, "Record", recordRaabbit.RecordId.String()).Result(); ok {
			cate, _ := redis.HGet(ctx, "Record", recordRaabbit.RecordId.String()).Result()
			// TODO 移除损坏数据
			redis.HDel(ctx, "Record", recordRaabbit.RecordId.String())
			if json.Unmarshal([]byte(cate), &record) == nil {
				// TODO 跳过数据库搜寻过程
				goto feep
			}
		}

		// TODO 未能找到提交记录
		if db.Where("id = (?)", recordRaabbit.RecordId).First(&record).Error != nil {
			log.Printf("%s Record Disappear!!\n", recordRaabbit.RecordId.String())
			return
		}
	feep:
		// TODO 尝试将该提交转发至指定平台
		if proid, source, err := util.DeCodeUUID(record.ProblemId); err == nil {
			// TODO 提交至指定平台
			runid, err := common.VjudgeMap[source].Submit(record.Code, proid, record.Language)
			if err != nil {
				log.Println("转发失败\n", err, record)
				continue
			}
			// TODO 跟踪指定提交
			go TrackerMap[recordRaabbit.Type](source, runid, proid, record)
		} else {
			log.Println("错误的目标平台\n", record)
			continue
		}
	}
}

// @title    Clouder
// @description   云转发器，用于转发提交
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   void
func Clouder() {
	redis := common.GetRedisClient(0)
	db := common.GetDB()
	ctx := context.Background()
	// TODO 订阅消息
	pubSub := redis.Subscribe(ctx, "Cloud")
	defer pubSub.Close()
	// TODO 获得消息管道
	ch := pubSub.Channel()
	// TODO 监听消息
	for msg := range ch {
		var recordRaabbit vo.RecordRabbit
		json.Unmarshal([]byte(msg.Payload), &recordRaabbit)
		// TODO 尝试在redis中取出该提交
		var record model.Record
		// TODO 先看redis中是否存在
		if ok, _ := redis.HExists(ctx, "Record", recordRaabbit.RecordId.String()).Result(); ok {
			cate, _ := redis.HGet(ctx, "Record", recordRaabbit.RecordId.String()).Result()
			// TODO 移除损坏数据
			redis.HDel(ctx, "Record", recordRaabbit.RecordId.String())
			if json.Unmarshal([]byte(cate), &record) == nil {
				// TODO 跳过数据库搜寻过程
				goto feep
			}
		}

		// TODO 未能找到提交记录
		if db.Where("id = (?)", recordRaabbit.RecordId).First(&record).Error != nil {
			log.Printf("%s Record Disappear!!\n", recordRaabbit.RecordId.String())
			return
		}
	feep:
		// TODO 尝试将该提交转发至云端
		go ClouderMap[recordRaabbit.Type](record)
	}
}

// TrackerMap		跟踪器映射表
var TrackerMap map[string]func(source string, runid string, proid string, record model.Record) = map[string]func(source string, runid string, proid string, record model.Record){
	"Group":  GroupTracker,
	"Match":  GroupTracker,
	"Normal": NormalTracker,
	"Single": SingleTracker,
	"OI":     SingleTracker,
}

// ClouderMap		云端映射表
var ClouderMap map[string]func(record model.Record) = map[string]func(record model.Record){
	"Group":  GroupClouder,
	"Match":  GroupClouder,
	"Normal": NormalClouder,
	"Single": SingleClouder,
	"OI":     SingleClouder,
}

// @title    NormalTracker
// @description   跟踪指定提交
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   void
func NormalTracker(source string, runid string, proid string, record model.Record) {
	channel := make(chan map[string]string)
	go common.VjudgeMap[source].GetStatus(runid, proid, channel)
	redis := common.GetRedisClient(0)
	db := common.GetDB()
	ctx := context.Background()
	// TODO 发布订阅用于提交列表
	recordList := vo.RecordList{
		RecordId: record.ID,
	}
	// TODO 提交订阅
	var recordCase vo.RecordCase
	recordCase.CaseId = 0
	var result map[string]string
	for result = range channel {
		// TODO 先矫正一波
		result["Result"] = util.StateCorrection(result["Result"])
		if result["Result"] != record.Condition {
			// TODO 更新状态
			record.Condition = result["Result"]
			record.Html = result["Html"]
			// TODO 将recordlist打包
			v, _ := json.Marshal(recordList)
			redis.Publish(ctx, "RecordChan", v)
			// TODO 将record存入redis
			v, _ = json.Marshal(record)
			redis.HSet(ctx, "Record", fmt.Sprint(record.ID), v)
			// TODO 提交长连接
			{
				recordCase.Condition = result["Result"]
				v, _ := json.Marshal(recordCase)
				redis.Publish(ctx, "RecordChan"+record.ID.String(), v)
			}
		}
	}
	// TODO 将record存入数据库
	db.Save(&record)
	if record.Condition == "Accepted" {
		// TODO 检查是否是今日首次通过
		if db.Where("condition = Accepted and to_days(created_at) = to_days(now())").First(&model.Record{}).Error != nil {
			Handle.Behaviors["Accepts"].PublishBehavior(1, record.UserId)
		}
		// TODO 检查该题目是否是首次通过
		if db.Where("condition = Accepted and problem_id = ?", record.ProblemId).First(&model.Record{}).Error != nil {
			Handle.Behaviors["Days"].PublishBehavior(1, record.UserId)
		}
	}
	// TODO 将其花费的时间与空间存入数据库
	t, err := strconv.Atoi(result["Time"])
	if err != nil {
		t = 0
	}
	m, err := strconv.Atoi(result["Memory"])
	if err != nil {
		m = 0
	}
	// TODO 将其作为第一个用例的结果存入数据库
	cas := model.CaseCondition{
		RecordId: record.ID,
		CID:      1,
		Time:     uint(t),
		Memory:   uint(m),
	}
	db.Create(&cas)
}

// @title    NormalClouder
// @description   普通提交到云端
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   void
func NormalClouder(record model.Record) {
	redis := common.GetRedisClient(0)
	db := common.GetDB()
	ctx := context.Background()
	// TODO 提交订阅
	var recordCase vo.RecordCase
	recordCase.CaseId = 0
	// TODO 发布订阅用于提交列表
	recordList := vo.RecordList{
		RecordId: record.ID,
	}
	// TODO 确保信息进入频道
	defer func() {
		// TODO 将recordlist打包
		v, _ := json.Marshal(recordList)
		redis.Publish(ctx, "RecordChan", v)
		// TODO 将record存入redis
		v, _ = json.Marshal(record)
		redis.HSet(ctx, "Record", fmt.Sprint(record.ID), v)
		// TODO 提交长连接
		{
			recordCase.Condition = record.Condition
			v, _ := json.Marshal(recordCase)
			redis.Publish(ctx, "RecordChan"+record.ID.String(), v)
		}
		// TODO 将record存入mysql
		db.Save(&record)
	}()

	// TODO 一些准备工作
	{
		record.Condition = "Preparing"
		// TODO 将recordlist打包
		v, _ := json.Marshal(recordList)
		redis.Publish(ctx, "RecordChan", v)
		// TODO 将record存入redis
		v, _ = json.Marshal(record)
		redis.HSet(ctx, "Record", fmt.Sprint(record.ID), v)
		// TODO 提交长连接
		{
			recordCase.Condition = "Preparing"
			v, _ := json.Marshal(recordCase)
			redis.Publish(ctx, "RecordChan"+record.ID.String(), v)
		}
	}
	// TODO 查看代码是否为空
	if record.Code == "" {
		record.Condition = "Code is empty"
		return
	}
	// TODO 找到提交记录后，开始判题逻辑
	if cmdI, ok := util.LanguageMap[record.Language]; ok {
		// TODO 从数据库中读出输入输出
		var cases []model.Case
		var problem model.Problem
		// TODO 先看redis中是否存在
		id := fmt.Sprint(record.ProblemId)
		if ok, _ := redis.HExists(ctx, "Problem", id).Result(); ok {
			cate, _ := redis.HGet(ctx, "Problem", id).Result()
			if json.Unmarshal([]byte(cate), &problem) == nil {
				// TODO 跳过数据库搜寻problem过程
				goto leep
			} else {
				// TODO 移除损坏数据
				redis.HDel(ctx, "Problem", id)
			}
		}

		// TODO 查看题目是否在数据库中存在
		if db.Where("id = (?)", id).First(&problem).Error != nil {
			record.Condition = "Problem Doesn't Exist"
			return
		}
		// TODO 将题目存入redis供下次使用
		{
			v, _ := json.Marshal(problem)
			redis.HSet(ctx, "Problem", id, v)
		}

	leep:

		// TODO 查找用例
		if ok, _ := redis.HExists(ctx, "Case", id).Result(); ok {
			cate, _ := redis.HGet(ctx, "Case", id).Result()
			if json.Unmarshal([]byte(cate), &cases) == nil {
				// TODO 跳过数据库搜寻testInputs过程
				goto Case
			} else {
				// TODO 移除损坏数据
				redis.HDel(ctx, "Case", id)
			}
		}

		// TODO 查看题目是否在数据库中存在
		if db.Where("problem_id = (?)", id).Find(&cases).Error != nil {
			record.Condition = "Input Doesn't Exist"
			return
		}
		// TODO 将题目存入redis供下次使用
		{
			v, _ := json.Marshal(cases)
			redis.HSet(ctx, "Case", id, v)
		}
	Case:

		// TODO 记录是否通过
		flag := true

		// TODO 开始运行工作
		{
			record.Condition = "Running"
			// TODO 将recordlist打包
			v, _ := json.Marshal(recordList)
			redis.Publish(ctx, "RecordChan", v)
			// TODO 将record存入redis
			v, _ = json.Marshal(record)
			redis.HSet(ctx, "Record", fmt.Sprint(record.ID), v)
			// TODO 提交长连接
			{
				recordCase.Condition = "Running"
				v, _ := json.Marshal(recordCase)
				redis.Publish(ctx, "RecordChan"+record.ID.String(), v)
			}
		}

		// TODO 最终将所有用例填入cases
		caseConditions := make([]model.CaseCondition, 0)

		defer func() {
			for i := range caseConditions {
				db.Create(&caseConditions[i])
			}
		}()

		for i := 0; i < len(cases); i++ {
			// TODO 将用例添加至最终数组
			cas := model.CaseCondition{
				RecordId: record.ID,
				Input:    cases[i].Input,
				Output:   cases[i].Output,
				CID:      uint(i + 1),
			}
			caseConditions = append(caseConditions, cas)

			// TODO 提交到云端
			condition, output, time, memory := TQ.CloudRun(record.Language, record.Code, cases[i].Input, problem.MemoryLimit*3, problem.TimeLimit*3)
			if condition != "ok" {
				record.Condition = condition
				flag = false
				goto final
			}

			// TODO 更新用例通过情况
			caseConditions[i].Time = uint(time)
			caseConditions[i].Memory = uint(memory)
			// TODO 超时
			if caseConditions[i].Time > problem.TimeLimit*cmdI.TimeMultiplier() {
				record.Condition = "Time Limit Exceeded"
				flag = false
				goto final
			}
			// TODO 超出内存限制
			if caseConditions[i].Memory > problem.MemoryLimit*cmdI.MemoryMultiplier() {
				record.Condition = "Memory Limit Exceeded"
				flag = false
				goto final
			}
			// TODO 答案错误
			var specalJudge model.Program
			// TODO 查看题目是否有标准程序

			// TODO 先看redis中是否存在
			if ok, _ := redis.HExists(ctx, "Program", problem.SpecialJudge.String()).Result(); ok {
				cate, _ := redis.HGet(ctx, "Program", problem.SpecialJudge.String()).Result()
				if json.Unmarshal([]byte(cate), &specalJudge) == nil {
					// TODO 跳过数据库搜寻program过程
					goto special
				} else {
					// TODO 移除损坏数据
					redis.HDel(ctx, "Program", problem.SpecialJudge.String())
				}
			}

			// TODO 查看程序是否在数据库中存在
			if db.Where("id = (?)", problem.SpecialJudge.String()).First(&specalJudge).Error != nil {
				goto outPut
			}
			// TODO 将题目存入redis供下次使用
			{
				v, _ := json.Marshal(specalJudge)
				redis.HSet(ctx, "Program", problem.SpecialJudge.String(), v)
			}
		special:
			// TODO 进行特判
			{
				if condition, output := TQ.JudgeRun(specalJudge.Language, specalJudge.Code, cases[i].Input+"\n"+output, problem.MemoryLimit*3, problem.TimeLimit*3); condition != "ok" || output != "ok" {
					record.Condition = condition
					flag = false
					goto final
				}
				goto pass
			}
		outPut:
			// TODO 正常判断
			if output != cases[i].Output {
				// TODO 去除格式后查看是否正确
				if util.RemoveWhiteSpace(output) == util.RemoveWhiteSpace(cases[i].Output) {
					record.Condition = "Presentation Error"
					flag = false
					goto final
				}
				record.Condition = "Wrong Answer"
				flag = false
				goto final
			}
		pass:
			// TODO 通过数量+1
			record.Pass++
			// TODO 长连接返回实时通过用例情况
			recordCase.CaseId++
			// TODO 将recordlist打包
			v, _ := json.Marshal(recordCase)
			redis.Publish(ctx, "RecordChan"+id, v)
		}
	final:
		// TODO 如果提交通过
		if flag {
			record.Condition = "Accepted"
			// TODO 检查是否是今日首次通过
			if db.Where("user_id = ? and condition = Accepted and to_days(created_at) = to_days(now())", record.UserId).First(&model.Record{}).Error != nil {
				Handle.Behaviors["Days"].PublishBehavior(1, record.UserId)
			}
			// TODO 检查该题目是否是首次通过
			if db.Where("user_id = ? and condition = Accepted and problem_id = ?", record.UserId, record.ProblemId).First(&model.Record{}).Error != nil {
				Handle.Behaviors["Accepts"].PublishBehavior(1, record.UserId)
				categoryMap := util.ProblemCategory(problem.ID)
				for category := range categoryMap {
					Handle.Behaviors[category].PublishBehavior(1, record.UserId)
				}
			}

		}
	} else {
		record.Condition = "Language Error"
	}
}

// @title    SingleTracker
// @description   跟踪指定提交
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   void
func SingleTracker(source string, runid string, proid string, record model.Record) {
	redix := common.GetRedisClient(0)
	db := common.GetDB()
	ctx := context.Background()
	var problem model.ProblemNew
	// TODO 先看redis中是否存在
	id := fmt.Sprint(record.ProblemId)
	if ok, _ := redix.HExists(ctx, "ProblemNew", id).Result(); ok {
		cate, _ := redix.HGet(ctx, "ProblemNew", id).Result()
		if json.Unmarshal([]byte(cate), &problem) == nil {
			// TODO 跳过数据库搜寻problem过程
			goto leep
		} else {
			// TODO 移除损坏数据
			redix.HDel(ctx, "ProblemNew", id)
		}
	}
	// TODO 查看题目是否在数据库中存在
	if db.Where("id = (?)", id).First(&problem).Error != nil {
		record.Condition = "Problem Doesn't Exist"
		// TODO 将record存入redis
		v, _ := json.Marshal(record)
		redix.HSet(ctx, "RecordCompetition", fmt.Sprint(record.ID), v)
		// TODO 将record存入mysql
		db.Save(&record)
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(problem)
		redix.HSet(ctx, "Problem", id, v)
	}
leep:
	var competition model.Competition
	// TODO 先看redis中是否存在
	if ok, _ := redix.HExists(ctx, "Competition", problem.CompetitionId.String()).Result(); ok {
		cate, _ := redix.HGet(ctx, "Competition", problem.CompetitionId.String()).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			// TODO 跳过数据库搜寻competition过程
			goto leap
		} else {
			// TODO 移除损坏数据
			redix.HDel(ctx, "Competition", problem.CompetitionId.String())
		}
	}
	// TODO 查看比赛是否在数据库中存在
	if db.Where("id = (?)", problem.CompetitionId.String()).First(&competition).Error != nil {
		record.Condition = "Competition Doesn't Exist"
		// TODO 将record存入redis
		v, _ := json.Marshal(record)
		redix.HSet(ctx, "RecordCompetition", fmt.Sprint(record.ID), v)
		// TODO 将record存入mysql
		db.Save(&record)
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(competition)
		redix.HSet(ctx, "Competition", problem.CompetitionId.String(), v)
	}
leap:
	channel := make(chan map[string]string)
	go common.VjudgeMap[source].GetStatus(runid, proid, channel)
	// TODO 发布订阅用于提交列表
	recordList := vo.RecordList{
		RecordId: record.ID,
	}
	// TODO 提交订阅
	var recordCase vo.RecordCase
	recordCase.CaseId = 0
	var result map[string]string
	for result = range channel {
		// TODO 先矫正一波
		result["Result"] = util.StateCorrection(result["Result"])
		if result["Result"] != record.Condition {
			// TODO 更新状态
			record.Condition = result["Result"]
			// TODO 将recordlist打包
			v, _ := json.Marshal(recordList)
			redix.Publish(ctx, "RecordChan", v)
			// TODO 将record存入redis
			v, _ = json.Marshal(record)
			redix.HSet(ctx, "Record", fmt.Sprint(record.ID), v)
			// TODO 提交长连接
			{
				recordCase.Condition = result["Result"]
				v, _ := json.Marshal(recordCase)
				redix.Publish(ctx, "RecordChan"+record.ID.String(), v)
			}
		}
	}
	// TODO 将record存入数据库
	db.Save(&record)
	if record.Condition == "Accepted" {
		// TODO 检查是否是今日首次通过
		if db.Where("user_id = ? and condition = Accepted and to_days(created_at) = to_days(now())", record.UserId).First(&model.Record{}).Error != nil {
			Handle.Behaviors["Accepts"].PublishBehavior(1, record.UserId)
		}
		// TODO 检查该题目是否是首次通过
		if db.Where("user_id = ? and condition = Accepted and problem_id = ?", record.UserId, record.ProblemId).First(&model.Record{}).Error != nil {
			Handle.Behaviors["Days"].PublishBehavior(1, record.UserId)
		}
	}
	// TODO 将其花费的时间与空间存入数据库
	t, err := strconv.Atoi(result["Time"])
	if err != nil {
		t = 0
	}
	m, err := strconv.Atoi(result["Memory"])
	if err != nil {
		m = 0
	}
	// TODO 将其作为第一个用例的结果存入数据库
	cas := model.CaseCondition{
		RecordId: record.ID,
		CID:      1,
		Time:     uint(t),
		Memory:   uint(m),
	}
	db.Create(&cas)

	if record.CreatedAt.After(competition.StartTime) && record.CreatedAt.Before(competition.EndTime) {
		var competitionMembers []model.CompetitionMember
		// TODO 在redis中取出成员罚时具体数据
		cM, err := redix.HGet(ctx, "Competition"+competition.ID.String(), record.UserId.String()).Result()
		if err == nil {
			json.Unmarshal([]byte(cM), &competitionMembers)
		}
		// TODO 找出数组中对应的题目
		k := -1
		for i := range competitionMembers {
			if competitionMembers[i].ProblemId == record.ProblemId {
				k = i
				break
			}
		}
		if k == -1 {
			k = len(competitionMembers)
			competitionMembers = append(competitionMembers, model.CompetitionMember{
				ID:            uuid.NewV4(),
				MemberId:      record.UserId,
				CompetitionId: competition.ID,
				ProblemId:     record.ProblemId,
				Pass:          0,
				Penalties:     0,
			})
		}
		// TODO 在redis中取出通过、罚时情况
		cR, err := redix.ZScore(ctx, "CompetitionR"+competition.ID.String(), record.UserId.String()).Result()
		if err != nil {
			cR = 0
		}
		// TODO 先前没有通过
		if competitionMembers[k].Condition != "Accepted" {
			// TODO 记录罚时
			competitionMembers[k].Penalties += time.Time(record.CreatedAt).Sub(time.Time(competition.StartTime))
			competitionMembers[k].Condition = record.Condition
			// TODO 通过样例数量增加
			if record.Condition == "Accepted" {
				// TODO 获取用例通过分数
				score := int(problem.Score)
				// TODO 如果分数上升
				if score > 0 {
					competitionMembers[k].RecordId = record.ID
					// TODO 记入罚时
					cR -= float64(competitionMembers[k].Penalties) / 10000000000
					cR += float64(score)
					// TODO 存入redis供下次使用
					v, _ := json.Marshal(competitionMembers)
					redix.HSet(ctx, "Competition"+competition.ID.String(), record.UserId.String(), v)
					redix.ZAdd(ctx, "CompetitionR"+competition.ID.String(), redis.Z{Score: cR, Member: record.UserId.String()})
					// TODO 发布订阅用于滚榜
					rankList := vo.RankList{
						MemberId: record.UserId,
					}
					// TODO 将ranklist打包
					v, _ = json.Marshal(rankList)
					redix.Publish(ctx, "CompetitionChan"+competition.ID.String(), v)
				}
			}
		}
	}
}

// @title    SingleClouder
// @description   单人赛提交到云端
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   void
func SingleClouder(record model.Record) {
	redix := common.GetRedisClient(0)
	db := common.GetDB()
	ctx := context.Background()
	log.Println("Single work for record")

	var problem model.ProblemNew
	// TODO 先看redis中是否存在
	id := fmt.Sprint(record.ProblemId)
	if ok, _ := redix.HExists(ctx, "ProblemNew", id).Result(); ok {
		cate, _ := redix.HGet(ctx, "ProblemNew", id).Result()
		if json.Unmarshal([]byte(cate), &problem) == nil {
			// TODO 跳过数据库搜寻problem过程
			goto leep
		} else {
			// TODO 移除损坏数据
			redix.HDel(ctx, "ProblemNew", id)
		}
	}
	// TODO 查看题目是否在数据库中存在
	if db.Where("id = (?)", id).First(&problem).Error != nil {
		record.Condition = "Problem Doesn't Exist"
		// TODO 将record存入redis
		v, _ := json.Marshal(record)
		redix.HSet(ctx, "RecordCompetition", fmt.Sprint(record.ID), v)
		// TODO 将record存入mysql
		db.Save(&record)
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(problem)
		redix.HSet(ctx, "Problem", id, v)
	}
leep:
	var competition model.Competition
	// TODO 先看redis中是否存在
	if ok, _ := redix.HExists(ctx, "Competition", problem.CompetitionId.String()).Result(); ok {
		cate, _ := redix.HGet(ctx, "Competition", problem.CompetitionId.String()).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			// TODO 跳过数据库搜寻competition过程
			goto leap
		} else {
			// TODO 移除损坏数据
			redix.HDel(ctx, "Competition", problem.CompetitionId.String())
		}
	}
	// TODO 查看比赛是否在数据库中存在
	if db.Where("id = (?)", problem.CompetitionId.String()).First(&competition).Error != nil {
		record.Condition = "Competition Doesn't Exist"
		// TODO 将record存入redis
		v, _ := json.Marshal(record)
		redix.HSet(ctx, "RecordCompetition", fmt.Sprint(record.ID), v)
		// TODO 将record存入mysql
		db.Save(&record)
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(competition)
		redix.HSet(ctx, "Competition", problem.CompetitionId.String(), v)
	}
leap:
	// TODO 发布订阅用于提交列表
	recordList := vo.RecordList{
		RecordId: record.ID,
	}
	// TODO 确保信息进入频道
	defer func() {
		// TODO 将recordlist打包
		v, _ := json.Marshal(recordList)
		redix.Publish(ctx, "RecordCompetitionChan"+problem.CompetitionId.String(), v)
		// TODO 将record存入redis
		v, _ = json.Marshal(record)
		redix.HSet(ctx, "RecordCompetition", fmt.Sprint(record.ID), v)
		// TODO 将record存入mysql
		db.Save(&record)
		log.Println("Record Save", record)
	}()
	// TODO 一些准备工作
	{
		record.Condition = "Preparing"
		// TODO 将recordlist打包
		v, _ := json.Marshal(recordList)
		redix.Publish(ctx, "RecordCompetitionChan"+problem.CompetitionId.String(), v)
		// TODO 将record存入redis
		v, _ = json.Marshal(record)
		redix.HSet(ctx, "RecordCompetition", fmt.Sprint(record.ID), v)
	}
	// TODO 查看代码是否为空
	if record.Code == "" {
		record.Condition = "Code is empty"
		return
	}
	// TODO 找到提交记录后，开始判题逻辑
	if cmdI, ok := util.LanguageMap[record.Language]; ok {
		// TODO 从数据库中读出输入输出
		var cases []model.Case

		// TODO 查找用例
		if ok, _ := redix.HExists(ctx, "Case", id).Result(); ok {
			cate, _ := redix.HGet(ctx, "Case", id).Result()
			if json.Unmarshal([]byte(cate), &cases) == nil {
				// TODO 跳过数据库搜寻testInputs过程
				goto Case
			} else {
				// TODO 移除损坏数据
				redix.HDel(ctx, "Case", id)
			}
		}

		// TODO 查看题目是否在数据库中存在
		if db.Where("problem_id = (?)", id).Find(&cases).Error != nil {
			record.Condition = "Input Doesn't Exist"
			return
		}
		// TODO 将题目存入redis供下次使用
		{
			v, _ := json.Marshal(cases)
			redix.HSet(ctx, "Case", id, v)
		}
	Case:

		// TODO 记录是否通过
		flag := true

		// TODO 开始运行工作
		{
			record.Condition = "Runing"
			// TODO 将recordlist打包
			v, _ := json.Marshal(recordList)
			redix.Publish(ctx, "RecordCompetitionChan"+problem.CompetitionId.String(), v)
			// TODO 将record存入redis
			v, _ = json.Marshal(record)
			redix.HSet(ctx, "RecordCompetition", fmt.Sprint(record.ID), v)
		}

		// TODO 最终将所有用例填入cases
		caseConditions := make([]model.CaseCondition, 0)

		defer func() {
			for i := range caseConditions {
				db.Create(caseConditions[i])
			}
		}()

		for i := 0; i < len(cases); i++ {
			// TODO 将用例添加至最终数组
			cas := model.CaseCondition{
				RecordId: record.ID,
				Input:    cases[i].Input,
				CID:      uint(i + 1),
			}
			caseConditions = append(caseConditions, cas)

			// TODO 提交到云端
			condition, output, time, memory := TQ.CloudRun(record.Language, record.Code, cases[i].Input, problem.MemoryLimit*3, problem.TimeLimit*3)
			if condition != "ok" {
				record.Condition = condition
				flag = false
				goto final
			}
			// TODO 更新用例通过情况
			caseConditions[i].Time = uint(time)
			caseConditions[i].Memory = uint(memory)
			// TODO 超时
			if caseConditions[i].Time > problem.TimeLimit*cmdI.TimeMultiplier() {
				record.Condition = "Time Limit Exceeded"
				flag = false
				goto final
			}
			// TODO 超出内存限制
			if caseConditions[i].Memory > problem.MemoryLimit*cmdI.MemoryMultiplier() {
				record.Condition = "Memory Limit Exceeded"
				flag = false
				goto final
			}
			// TODO 答案错误
			var specalJudge model.Program
			// TODO 查看题目是否有标准程序

			// TODO 先看redis中是否存在
			if ok, _ := redix.HExists(ctx, "Program", problem.SpecialJudge.String()).Result(); ok {
				cate, _ := redix.HGet(ctx, "Program", problem.SpecialJudge.String()).Result()
				if json.Unmarshal([]byte(cate), &specalJudge) == nil {
					// TODO 跳过数据库搜寻program过程
					goto special
				} else {
					// TODO 移除损坏数据
					redix.HDel(ctx, "Program", problem.SpecialJudge.String())
				}
			}

			// TODO 查看程序是否在数据库中存在
			if db.Where("id = (?)", problem.SpecialJudge.String()).First(&specalJudge).Error != nil {
				goto outPut
			}
			// TODO 将题目存入redis供下次使用
			{
				v, _ := json.Marshal(specalJudge)
				redix.HSet(ctx, "Program", problem.SpecialJudge.String(), v)
			}
		special:
			// TODO 进行特判
			{
				if condition, output := TQ.JudgeRun(specalJudge.Language, specalJudge.Code, cases[i].Input+"\n"+output, problem.MemoryLimit*3, problem.TimeLimit*3); condition != "ok" || output != "ok" {
					record.Condition = condition
					flag = false
					goto final
				}
				goto pass
			}
		outPut:
			// TODO 正常判断
			if output != cases[i].Output {
				// TODO 去除格式后查看是否正确
				if util.RemoveWhiteSpace(output) == util.RemoveWhiteSpace(cases[i].Output) {
					record.Condition = " Presentation Error"
					flag = false
					goto final
				}
				record.Condition = "Wrong Answer"
				flag = false
				goto final
			}
		pass:
			// TODO 通过数量+1
			record.Pass++

			// TODO 数据库插入数据错误
			if db.Create(&cas).Error != nil {
				record.Condition = "System error 5"
				return
			}
		}
	final:
		// TODO 如果提交通过
		if flag {
			record.Condition = "Accepted"
			// TODO 检查是否是今日首次通过
			if db.Where("user_id = ? and condition = Accepted and to_days(created_at) = to_days(now())", record.UserId).First(&model.Record{}).Error != nil {
				Handle.Behaviors["Days"].PublishBehavior(1, record.UserId)
			}
			// TODO 检查该题目是否是首次通过
			if db.Where("user_id = ? and condition = Accepted and problem_id = ?", record.UserId, record.ProblemId).First(&model.Record{}).Error != nil {
				Handle.Behaviors["Accepts"].PublishBehavior(1, record.UserId)
				categoryMap := util.ProblemCategory(problem.ID)
				for category := range categoryMap {
					Handle.Behaviors[category].PublishBehavior(1, record.UserId)
				}
			}
		}
		// TODO 查看是否为比赛提交,且比赛已经开始
		if record.CreatedAt.After(competition.StartTime) && record.CreatedAt.Before(competition.EndTime) {
			var competitionMembers []model.CompetitionMember
			// TODO 在redis中取出成员罚时具体数据
			cM, err := redix.HGet(ctx, "Competition"+competition.ID.String(), record.UserId.String()).Result()
			if err == nil {
				json.Unmarshal([]byte(cM), &competitionMembers)
			}
			// TODO 找出数组中对应的题目
			k := -1
			for i := range competitionMembers {
				if competitionMembers[i].ProblemId == record.ProblemId {
					k = i
					break
				}
			}
			if k == -1 {
				k = len(competitionMembers)
				competitionMembers = append(competitionMembers, model.CompetitionMember{
					ID:            uuid.NewV4(),
					MemberId:      record.UserId,
					CompetitionId: competition.ID,
					ProblemId:     record.ProblemId,
					Pass:          0,
					Penalties:     0,
				})
			}
			// TODO 在redis中取出通过、罚时情况
			cR, err := redix.ZScore(ctx, "CompetitionR"+competition.ID.String(), record.UserId.String()).Result()
			if err != nil {
				cR = 0
			}
			// TODO 先前没有通过
			if competitionMembers[k].Condition != "Accepted" {
				// TODO 记录罚时
				competitionMembers[k].Penalties += time.Time(record.CreatedAt).Sub(time.Time(competition.StartTime))
				// TODO 通过样例数量增加
				if competitionMembers[k].Pass < record.Pass {
					competitionMembers[k].Condition = record.Condition
					// TODO 获取用例通过分数
					score := 0
					for i := competitionMembers[k].Pass + 1; i <= record.Pass; i++ {
						score += int(cases[i].Score)
					}
					// TODO 如果分数上升
					if score > 0 {
						competitionMembers[k].RecordId = record.ID
						// TODO 记入罚时
						cR -= float64(competitionMembers[k].Penalties) / 10000000000
						cR += float64(score)
						// TODO 存入redis供下次使用
						v, _ := json.Marshal(competitionMembers)
						redix.HSet(ctx, "Competition"+competition.ID.String(), record.UserId.String(), v)
						redix.ZAdd(ctx, "CompetitionR"+competition.ID.String(), redis.Z{Score: cR, Member: record.UserId.String()})
						// TODO 发布订阅用于滚榜
						rankList := vo.RankList{
							MemberId: record.UserId,
						}
						// TODO 将ranklist打包
						v, _ = json.Marshal(rankList)
						redix.Publish(ctx, "CompetitionChan"+competition.ID.String(), v)
					}
				}
			}
		}
	} else {
		record.Condition = "Language Error"
	}
}

// @title    GroupTracker
// @description   跟踪指定提交
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   void
func GroupTracker(source string, runid string, proid string, record model.Record) {
	redix := common.GetRedisClient(0)
	db := common.GetDB()
	ctx := context.Background()
	var problem model.ProblemNew
	// TODO 先看redis中是否存在
	id := fmt.Sprint(record.ProblemId)
	if ok, _ := redix.HExists(ctx, "ProblemNew", id).Result(); ok {
		cate, _ := redix.HGet(ctx, "ProblemNew", id).Result()
		if json.Unmarshal([]byte(cate), &problem) == nil {
			// TODO 跳过数据库搜寻problem过程
			goto leep
		} else {
			// TODO 移除损坏数据
			redix.HDel(ctx, "ProblemNew", id)
		}
	}
	// TODO 查看题目是否在数据库中存在
	if db.Where("id = (?)", id).First(&problem).Error != nil {
		record.Condition = "Problem Doesn't Exist"
		// TODO 将record存入redis
		v, _ := json.Marshal(record)
		redix.HSet(ctx, "RecordCompetition", fmt.Sprint(record.ID), v)
		// TODO 将record存入mysql
		db.Save(&record)
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(problem)
		redix.HSet(ctx, "Problem", id, v)
	}
leep:
	var competition model.Competition
	// TODO 先看redis中是否存在
	if ok, _ := redix.HExists(ctx, "Competition", problem.CompetitionId.String()).Result(); ok {
		cate, _ := redix.HGet(ctx, "Competition", problem.CompetitionId.String()).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			// TODO 跳过数据库搜寻competition过程
			goto leap
		} else {
			// TODO 移除损坏数据
			redix.HDel(ctx, "Competition", problem.CompetitionId.String())
		}
	}
	// TODO 查看比赛是否在数据库中存在
	if db.Where("id = (?)", problem.CompetitionId.String()).First(&competition).Error != nil {
		record.Condition = "Competition Doesn't Exist"
		// TODO 将record存入redis
		v, _ := json.Marshal(record)
		redix.HSet(ctx, "RecordCompetition", fmt.Sprint(record.ID), v)
		// TODO 将record存入mysql
		db.Save(&record)
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(competition)
		redix.HSet(ctx, "Competition", problem.CompetitionId.String(), v)
	}
leap:
	channel := make(chan map[string]string)
	go common.VjudgeMap[source].GetStatus(runid, proid, channel)
	// TODO 发布订阅用于提交列表
	recordList := vo.RecordList{
		RecordId: record.ID,
	}
	// TODO 提交订阅
	var recordCase vo.RecordCase
	recordCase.CaseId = 0
	var result map[string]string
	for result = range channel {
		// TODO 先矫正一波
		result["Result"] = util.StateCorrection(result["Result"])
		if result["Result"] != record.Condition {
			// TODO 更新状态
			record.Condition = result["Result"]
			// TODO 将recordlist打包
			v, _ := json.Marshal(recordList)
			redix.Publish(ctx, "RecordChan", v)
			// TODO 将record存入redis
			v, _ = json.Marshal(record)
			redix.HSet(ctx, "Record", fmt.Sprint(record.ID), v)
			// TODO 提交长连接
			{
				recordCase.Condition = result["Result"]
				v, _ := json.Marshal(recordCase)
				redix.Publish(ctx, "RecordChan"+record.ID.String(), v)
			}
		}
	}
	// TODO 将record存入数据库
	db.Save(&record)
	if record.Condition == "Accepted" {
		// TODO 检查是否是今日首次通过
		if db.Where("user_id = ? and condition = Accepted and to_days(created_at) = to_days(now())", record.UserId).First(&model.Record{}).Error != nil {
			Handle.Behaviors["Days"].PublishBehavior(1, record.UserId)
		}
		// TODO 检查该题目是否是首次通过
		if db.Where("user_id = ? and condition = Accepted and problem_id = ?", record.UserId, record.ProblemId).First(&model.Record{}).Error != nil {
			Handle.Behaviors["Accepts"].PublishBehavior(1, record.UserId)
			categoryMap := util.ProblemCategory(problem.ID)
			for category := range categoryMap {
				Handle.Behaviors[category].PublishBehavior(1, record.UserId)
			}
		}
	}
	// TODO 将其花费的时间与空间存入数据库
	t, err := strconv.Atoi(result["Time"])
	if err != nil {
		t = 0
	}
	m, err := strconv.Atoi(result["Memory"])
	if err != nil {
		m = 0
	}
	// TODO 将其作为第一个用例的结果存入数据库
	cas := model.CaseCondition{
		RecordId: record.ID,
		CID:      1,
		Time:     uint(t),
		Memory:   uint(m),
	}
	db.Create(&cas)

	if record.CreatedAt.After(competition.StartTime) && record.CreatedAt.Before(competition.EndTime) {
		groups, _ := redix.ZRange(ctx, "CompetitionR"+competition.ID.String(), 0, -1).Result()
		// TODO 查找组
		var group model.Group
		for i := range groups {

			// TODO 先看redis中是否存在
			if ok, _ := redix.HExists(ctx, "Group", groups[i]).Result(); ok {
				cate, _ := redix.HGet(ctx, "Group", groups[i]).Result()
				if json.Unmarshal([]byte(cate), &group) == nil {
					goto levp
				} else {
					// TODO 移除损坏数据
					redix.HDel(ctx, "Group", groups[i])
				}
			}

			// TODO 查看用户组是否在数据库中存在
			db.Where("id = (?)", groups[i]).First(&group)
			{
				// TODO 将用户组存入redis供下次使用
				v, _ := json.Marshal(group)
				redix.HSet(ctx, "Group", groups[i], v)
			}
		levp:
			if db.Where("group_id = (?) and user_id = (?)", group.ID, record.UserId).First(&model.UserList{}).Error == nil {
				break
			}
		}
		var competitionMembers []model.CompetitionMember
		// TODO 在redis中取出成员罚时具体数据
		cM, err := redix.HGet(ctx, "Competition"+competition.ID.String(), record.UserId.String()).Result()
		if err == nil {
			json.Unmarshal([]byte(cM), &competitionMembers)
		}
		// TODO 找出数组中对应的题目
		k := -1
		for i := range competitionMembers {
			if competitionMembers[i].ProblemId == record.ProblemId {
				k = i
				break
			}
		}
		if k == -1 {
			k = len(competitionMembers)
			competitionMembers = append(competitionMembers, model.CompetitionMember{
				ID:            uuid.NewV4(),
				MemberId:      group.ID,
				CompetitionId: competition.ID,
				ProblemId:     record.ProblemId,
				Pass:          0,
				Penalties:     0,
			})
		}
		// TODO 在redis中取出通过、罚时情况
		cR, err := redix.ZScore(ctx, "CompetitionR"+competition.ID.String(), group.ID.String()).Result()
		if err != nil {
			cR = 0
		}
		// TODO 先前没有通过
		if competitionMembers[k].Condition != "Accepted" {
			// TODO 记录罚时
			competitionMembers[k].Penalties += time.Time(record.CreatedAt).Sub(time.Time(competition.StartTime))
			competitionMembers[k].Condition = record.Condition
			// TODO 通过样例数量增加
			if record.Condition == "Accepted" {
				// TODO 获取用例通过分数
				score := int(problem.Score)
				// TODO 如果分数上升
				if score > 0 {
					competitionMembers[k].RecordId = record.ID
					// TODO 记入罚时
					cR -= float64(competitionMembers[k].Penalties) / 10000000000
					cR += float64(score)
					// TODO 存入redis供下次使用
					v, _ := json.Marshal(competitionMembers)
					redix.HSet(ctx, "Competition"+competition.ID.String(), group.ID.String(), v)
					redix.ZAdd(ctx, "CompetitionR"+competition.ID.String(), redis.Z{Score: cR, Member: group.ID.String()})
					// TODO 发布订阅用于滚榜
					rankList := vo.RankList{
						MemberId: group.ID,
					}
					// TODO 将ranklist打包
					v, _ = json.Marshal(rankList)
					redix.Publish(ctx, "CompetitionChan"+competition.ID.String(), v)
				}
			}
		}
	}
}

// @title    GroupClouder
// @description   组队赛提交到云端
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   void
func GroupClouder(record model.Record) {
	redix := common.GetRedisClient(0)
	db := common.GetDB()
	ctx := context.Background()
	log.Println("Single work for record")

	var problem model.ProblemNew
	// TODO 先看redis中是否存在
	id := fmt.Sprint(record.ProblemId)
	if ok, _ := redix.HExists(ctx, "ProblemNew", id).Result(); ok {
		cate, _ := redix.HGet(ctx, "ProblemNew", id).Result()
		if json.Unmarshal([]byte(cate), &problem) == nil {
			// TODO 跳过数据库搜寻problem过程
			goto leep
		} else {
			// TODO 移除损坏数据
			redix.HDel(ctx, "ProblemNew", id)
		}
	}

	// TODO 查看题目是否在数据库中存在
	if db.Where("id = (?)", id).First(&problem).Error != nil {
		record.Condition = "Problem Doesn't Exist"
		// TODO 将record存入redis
		v, _ := json.Marshal(record)
		redix.HSet(ctx, "RecordCompetition", fmt.Sprint(record.ID), v)
		// TODO 将record存入mysql
		db.Save(&record)
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(problem)
		redix.HSet(ctx, "Problem", id, v)
	}

leep:
	var competition model.Competition
	// TODO 先看redis中是否存在
	if ok, _ := redix.HExists(ctx, "Competition", problem.CompetitionId.String()).Result(); ok {
		cate, _ := redix.HGet(ctx, "Competition", problem.CompetitionId.String()).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			// TODO 跳过数据库搜寻competition过程
			goto leap
		} else {
			// TODO 移除损坏数据
			redix.HDel(ctx, "Competition", problem.CompetitionId.String())
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if db.Where("id = (?)", problem.CompetitionId.String()).First(&competition).Error != nil {
		record.Condition = "Competition Doesn't Exist"
		// TODO 将record存入redis
		v, _ := json.Marshal(record)
		redix.HSet(ctx, "RecordCompetition", fmt.Sprint(record.ID), v)
		// TODO 将record存入mysql
		db.Save(&record)
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(competition)
		redix.HSet(ctx, "Competition", problem.CompetitionId.String(), v)
	}

leap:

	// TODO 发布订阅用于提交列表
	recordList := vo.RecordList{
		RecordId: record.ID,
	}
	// TODO 确保信息进入频道
	defer func() {
		// TODO 将recordlist打包
		v, _ := json.Marshal(recordList)
		redix.Publish(ctx, "RecordCompetitionChan"+problem.CompetitionId.String(), v)
		// TODO 将record存入redis
		v, _ = json.Marshal(record)
		redix.HSet(ctx, "RecordCompetition", fmt.Sprint(record.ID), v)
		// TODO 将record存入mysql
		db.Save(&record)
	}()

	// TODO 一些准备工作
	{
		record.Condition = "Preparing"
		// TODO 将recordlist打包
		v, _ := json.Marshal(recordList)
		redix.Publish(ctx, "RecordCompetitionChan"+problem.CompetitionId.String(), v)
		// TODO 将record存入redis
		v, _ = json.Marshal(record)
		redix.HSet(ctx, "RecordCompetition", fmt.Sprint(record.ID), v)
	}
	// TODO 查看代码是否为空
	if record.Code == "" {
		record.Condition = "Code is empty"
		return
	}
	// TODO 找到提交记录后，开始判题逻辑
	if cmdI, ok := util.LanguageMap[record.Language]; ok {
		// TODO 从数据库中读出输入输出
		var cases []model.Case

		// TODO 查找用例
		if ok, _ := redix.HExists(ctx, "Case", id).Result(); ok {
			cate, _ := redix.HGet(ctx, "Case", id).Result()
			if json.Unmarshal([]byte(cate), &cases) == nil {
				// TODO 跳过数据库搜寻testInputs过程
				goto Case
			} else {
				// TODO 移除损坏数据
				redix.HDel(ctx, "Case", id)
			}
		}

		// TODO 查看题目是否在数据库中存在
		if db.Where("problem_id = (?)", id).Find(&cases).Error != nil {
			record.Condition = "Input Doesn't Exist"
			return
		}
		// TODO 将题目存入redis供下次使用
		{
			v, _ := json.Marshal(cases)
			redix.HSet(ctx, "Case", id, v)
		}
	Case:

		// TODO 记录是否通过
		flag := true

		// TODO 开始运行工作
		{
			record.Condition = "Runing"
			// TODO 将recordlist打包
			v, _ := json.Marshal(recordList)
			redix.Publish(ctx, "RecordCompetitionChan"+problem.CompetitionId.String(), v)
			// TODO 将record存入redis
			v, _ = json.Marshal(record)
			redix.HSet(ctx, "RecordCompetition", fmt.Sprint(record.ID), v)
		}

		// TODO 最终将所有用例填入cases
		caseConditions := make([]model.CaseCondition, 0)

		defer func() {
			for i := range caseConditions {
				db.Create(caseConditions[i])
			}
		}()

		for i := 0; i < len(cases); i++ {
			// TODO 将用例添加至最终数组
			cas := model.CaseCondition{
				RecordId: record.ID,
				CID:      uint(i + 1),
			}
			caseConditions = append(caseConditions, cas)

			// TODO 提交到云端
			condition, output, time, memory := TQ.CloudRun(record.Language, record.Code, cases[i].Input, problem.MemoryLimit*3, problem.TimeLimit*3)
			if condition != "ok" {
				record.Condition = condition
				flag = false
				goto final
			}

			// TODO 更新用例通过情况
			caseConditions[i].Time = uint(time)
			caseConditions[i].Memory = uint(memory)
			// TODO 超时
			if caseConditions[i].Time > problem.TimeLimit*cmdI.TimeMultiplier() {
				record.Condition = "Time Limit Exceeded"
				flag = false
				goto final
			}
			// TODO 超出内存限制
			if caseConditions[i].Memory > problem.MemoryLimit*cmdI.MemoryMultiplier() {
				record.Condition = "Memory Limit Exceeded"
				flag = false
				goto final
			}
			// TODO 答案错误
			var specalJudge model.Program
			// TODO 查看题目是否有标准程序

			// TODO 先看redis中是否存在
			if ok, _ := redix.HExists(ctx, "Program", problem.SpecialJudge.String()).Result(); ok {
				cate, _ := redix.HGet(ctx, "Program", problem.SpecialJudge.String()).Result()
				if json.Unmarshal([]byte(cate), &specalJudge) == nil {
					// TODO 跳过数据库搜寻program过程
					goto special
				} else {
					// TODO 移除损坏数据
					redix.HDel(ctx, "Program", problem.SpecialJudge.String())
				}
			}

			// TODO 查看程序是否在数据库中存在
			if db.Where("id = (?)", problem.SpecialJudge.String()).First(&specalJudge).Error != nil {
				goto outPut
			}
			// TODO 将题目存入redis供下次使用
			{
				v, _ := json.Marshal(specalJudge)
				redix.HSet(ctx, "Program", problem.SpecialJudge.String(), v)
			}
		special:
			// TODO 进行特判
			{
				if condition, output := TQ.JudgeRun(specalJudge.Language, specalJudge.Code, cases[i].Input+"\n"+output, problem.MemoryLimit*3, problem.TimeLimit*3); condition != "ok" || output != "ok" {
					record.Condition = condition
					flag = false
					goto final
				}
				goto pass
			}
		outPut:
			// TODO 正常判断
			if output != cases[i].Output {
				// TODO 去除格式后查看是否正确
				if util.RemoveWhiteSpace(output) == util.RemoveWhiteSpace(cases[i].Output) {
					record.Condition = " Presentation Error"
					flag = false
					goto final
				}
				record.Condition = "Wrong Answer"
				flag = false
				goto final
			}
		pass:
			// TODO 通过数量+1
			record.Pass++

			// TODO 数据库插入数据错误
			if db.Create(&cas).Error != nil {
				record.Condition = "System error 5"
				return
			}
		}
	final:
		// TODO 如果提交通过
		if flag {
			record.Condition = "Accepted"
			// TODO 检查是否是今日首次通过
			if db.Where("user_id = ? and condition = Accepted and to_days(created_at) = to_days(now())", record.UserId).First(&model.Record{}).Error != nil {
				Handle.Behaviors["Days"].PublishBehavior(1, record.UserId)
			}
			// TODO 检查该题目是否是首次通过
			if db.Where("user_id = ? and condition = Accepted and problem_id = ?", record.UserId, record.ProblemId).First(&model.Record{}).Error != nil {
				Handle.Behaviors["Accepts"].PublishBehavior(1, record.UserId)
				categoryMap := util.ProblemCategory(problem.ID)
				for category := range categoryMap {
					Handle.Behaviors[category].PublishBehavior(1, record.UserId)
				}
			}
		}
		// TODO 查看是否为比赛提交,且比赛已经开始
		if record.CreatedAt.After(competition.StartTime) && record.CreatedAt.Before(competition.EndTime) {
			groups, _ := redix.ZRange(ctx, "CompetitionR"+competition.ID.String(), 0, -1).Result()
			// TODO 查找组
			var group model.Group
			for i := range groups {

				// TODO 先看redis中是否存在
				if ok, _ := redix.HExists(ctx, "Group", groups[i]).Result(); ok {
					cate, _ := redix.HGet(ctx, "Group", groups[i]).Result()
					if json.Unmarshal([]byte(cate), &group) == nil {
						goto levp
					} else {
						// TODO 移除损坏数据
						redix.HDel(ctx, "Group", groups[i])
					}
				}

				// TODO 查看用户组是否在数据库中存在
				db.Where("id = (?)", groups[i]).First(&group)
				{
					// TODO 将用户组存入redis供下次使用
					v, _ := json.Marshal(group)
					redix.HSet(ctx, "Group", groups[i], v)
				}
			levp:
				if db.Where("group_id = (?) and user_id = (?)", group.ID, record.UserId).First(&model.UserList{}).Error == nil {
					break
				}
			}
			var competitionMembers []model.CompetitionMember
			// TODO 在redis中取出成员罚时具体数据
			cM, err := redix.HGet(ctx, "Competition"+competition.ID.String(), group.ID.String()).Result()
			if err == nil {
				json.Unmarshal([]byte(cM), &competitionMembers)
			}
			// TODO 找出数组中对应的题目
			k := -1
			for i := range competitionMembers {
				if competitionMembers[i].ProblemId == record.ProblemId {
					k = i
					break
				}
			}
			if k == -1 {
				k = len(competitionMembers)
				competitionMembers = append(competitionMembers, model.CompetitionMember{
					ID:            uuid.NewV4(),
					MemberId:      group.ID,
					CompetitionId: competition.ID,
					ProblemId:     record.ProblemId,
					Pass:          0,
					Penalties:     0,
				})
			}
			// TODO 在redis中取出通过、罚时情况
			cR, err := redix.ZScore(ctx, "CompetitionR"+competition.ID.String(), group.ID.String()).Result()
			if err != nil {
				cR = 0
			}
			// TODO 先前没有通过
			if competitionMembers[k].Condition != "Accepted" {
				// TODO 记录罚时
				competitionMembers[k].Penalties += time.Time(record.CreatedAt).Sub(time.Time(competition.StartTime))
				// TODO 通过样例数量增加
				if competitionMembers[k].Pass < record.Pass {
					competitionMembers[k].Condition = record.Condition
					// TODO 获取用例通过分数
					score := 0
					for i := competitionMembers[k].Pass + 1; i <= record.Pass; i++ {
						score += int(cases[i].Score)
					}
					// TODO 如果分数上升
					if score > 0 {
						competitionMembers[k].RecordId = record.ID
						// TODO 记入罚时
						cR -= float64(competitionMembers[k].Penalties) / 10000000000
						cR += float64(score)
						// TODO 存入redis供下次使用
						v, _ := json.Marshal(competitionMembers)
						redix.HSet(ctx, "Competition"+competition.ID.String(), group.ID.String(), v)
						redix.ZAdd(ctx, "CompetitionR"+competition.ID.String(), redis.Z{Score: cR, Member: group.ID.String()})
						// TODO 发布订阅用于滚榜
						rankList := vo.RankList{
							MemberId: group.ID,
						}
						// TODO 将ranklist打包
						v, _ = json.Marshal(rankList)
						redix.Publish(ctx, "CompetitionChan"+competition.ID.String(), v)
					}
				}
			}
		}
	} else {
		record.Condition = "Language Error"
	}
}

// @title    判断是否应该转发至云端
// @description   提交中间件
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   void
func SubmitCheck() bool {
	var Min float64 = 1.0
	// TODO 找出当前最小百分比
	for i := range HeartPercentageMap {
		if HeartPercentageMap[i] < Min {
			Min = HeartPercentageMap[i]
		}
	}
	// TODO 计算转发几率
	var p float64 = 0.0
	if Min > 0.5 {
		p = (Min - 0.5) * 2
	}
	// TODO 使用系统时间的不确定性来进行初始化
	rand.Seed(time.Now().Unix())
	pt := float64(rand.Intn(10000))
	// TODO 是否触发转发
	return pt < (p)*10000
}
