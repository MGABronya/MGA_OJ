// @Title  RecordController
// @Description  该文件提供关于操作提交的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"MGA_OJ/Interface"
	TQ "MGA_OJ/Test-request"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"
	"MGA_OJ/vo"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
// @auth      MGAronya（张健）       2022-9-16 12:15
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
	if r.DB.Where("id = ?", requestRecord.ProblemId).First(&problem).Error != nil {
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

	// TODO 加入消息队列
	v, _ := json.Marshal(recordRaabbit)
	if err := r.Rabbitmq.PublishSimple(string(v)); err != nil {
		response.Fail(ctx, nil, "消息队列出错")
		return
	}

	// TODO 发布订阅用于提交列表
	recordList := vo.RecordList{
		RecordId: record.ID,
	}
	// TODO 将recordlist打包
	v, _ = json.Marshal(recordList)
	r.Redis.Publish(ctx, "RecordChan", v)

	// TODO 成功
	response.Success(ctx, nil, "提交成功")
}

// @title    ShowRecord
// @description   查看一篇提交的内容
// @auth      MGAronya（张健）       2022-9-16 12:19
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
	if r.DB.Where("id = ?", id).First(&record).Error != nil {
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
// @auth      MGAronya（张健）       2022-9-16 12:20
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
		db = db.Where("Language = ?", Language)
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
// @auth      MGAronya（张健）       2022-9-16 12:19
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
func (r RecordController) CaseList(ctx *gin.Context) {
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
	if r.DB.Where("id = ?", id).First(&record).Error != nil {
		response.Fail(ctx, nil, "提交不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(record)
		r.Redis.HSet(ctx, "Record", id, v)
	}
leep:
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var cases []model.CaseCondition

	var total int64

	// TODO 查找所有分页中可见的条目
	r.DB.Where("record_id = ?", record.ID).Order("id asc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&cases)

	r.DB.Where("record_id = ?", record.ID).Model(model.CaseCondition{}).Count(&total)

	response.Success(ctx, gin.H{"cases": cases}, "成功")
}

// @title    Case
// @description   查看一篇测试通过情况
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RecordController) Case(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var cas model.CaseCondition

	// TODO 查找所有分页中可见的条目
	if r.DB.Where("id = ?", id).First(&cas).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}

	response.Success(ctx, gin.H{"case": cas}, "成功")
}

// @title    Hack
// @description   Hack功能
// @auth      MGAronya（张健）       2022-9-16 12:19
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
	if r.DB.Where("id = ?", id).First(&record).Error != nil {
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
	if r.DB.Where("id = ?", record.ProblemId.String()).First(&problem).Error != nil {
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
	if r.DB.Where("id = ?", problem.InputCheck.String()).First(&inputCheckProgram).Error != nil {
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
	if r.DB.Where("id = ?", problem.Standard.String()).First(&standardProgram).Error != nil {
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
		if r.DB.Where("id = ?", problem.SpecialJudge.String()).First(&specialJudgeProgram).Error != nil {
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
// @auth      MGAronya（张健）       2022-9-16 12:23
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
