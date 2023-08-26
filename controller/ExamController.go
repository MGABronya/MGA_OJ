// @Title  ExamController
// @Description  该文件提供关于操作测试的各种方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package controller

import (
	"MGA_OJ/Interface"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"
	"MGA_OJ/util"
	"MGA_OJ/vo"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IExamController			定义了测试类接口
type IExamController interface {
	Interface.RestInterface       // 包含增删查改功能
	ScoreShow(ctx *gin.Context)   // 分数查看
	ScoreUpdate(ctx *gin.Context) // 分数更新
	ScoreList(ctx *gin.Context)   // 分数列表
}

// ExamController			定义了测试工具类
type ExamController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Create
// @description   创建一篇测试
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (e ExamController) Create(ctx *gin.Context) {
	var requestExam vo.ExamRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestExam); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 查看是否是该组的组长
	// TODO 查找组
	var group model.Group

	// TODO 先看redis中是否存在
	if ok, _ := e.Redis.HExists(ctx, "Group", id).Result(); ok {
		cate, _ := e.Redis.HGet(ctx, "Group", id).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto levp
		} else {
			// TODO 移除损坏数据
			e.Redis.HDel(ctx, "Group", id)
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if e.DB.Where("id = (?)", id).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		e.Redis.HSet(ctx, "Group", id, v)
	}
levp:

	// TODO 查看是否是用户组组长
	if user.ID != group.LeaderId {
		response.Fail(ctx, nil, "非用户组组长，无法布置测试")
		return
	}

	// TODO 验证起始时间与终止时间是否合法
	if requestExam.StartTime.After(requestExam.EndTime) {
		response.Fail(ctx, nil, "起始时间大于了终止时间")
		return
	}
	if time.Now().After(time.Time(requestExam.StartTime)) {
		response.Fail(ctx, nil, "当前时间大于了起始时间")
		return
	}
	if time.Time(requestExam.EndTime).After(time.Now().Add(30 * 24 * time.Hour)) {
		response.Fail(ctx, nil, "终止时间不可设置为30日后")
		return
	}

	if requestExam.Type != "IOI" && requestExam.Type != "IO" {
		response.Fail(ctx, nil, "类型错误")
		return
	}

	// TODO 创建测试
	exam := model.Exam{
		Title:     requestExam.Title,
		Content:   requestExam.Content,
		ResLong:   requestExam.ResLong,
		ResShort:  requestExam.ResShort,
		GroupId:   group.ID,
		UserId:    user.ID,
		StartTime: requestExam.StartTime,
		EndTime:   requestExam.EndTime,
		Type:      requestExam.Type,
	}

	// TODO 插入数据
	if err := e.DB.Create(&exam).Error; err != nil {
		response.Fail(ctx, nil, "测试上传出错，数据验证有误")
		return
	}

	// TODO 成功
	response.Success(ctx, gin.H{"exam": exam}, "创建成功")

	// TODO 等待直至比赛结束
	ExamTimer(ctx, e.Redis, e.DB, exam)
}

// @title    Update
// @description   更新一篇测试的内容
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (e ExamController) Update(ctx *gin.Context) {
	var requestExam vo.ExamRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestExam); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应测试
	id := ctx.Params.ByName("id")

	var exam model.Exam

	if e.DB.Where("id = (?)", id).First(&exam).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}

	// TODO 查看测试所在组
	// TODO 查找组
	var group model.Group

	// TODO 先看redis中是否存在
	if ok, _ := e.Redis.HExists(ctx, "Group", exam.GroupId.String()).Result(); ok {
		cate, _ := e.Redis.HGet(ctx, "Group", exam.GroupId.String()).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto levp
		} else {
			// TODO 移除损坏数据
			e.Redis.HDel(ctx, "Group", exam.GroupId.String())
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if e.DB.Where("id = (?)", exam.GroupId).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		e.Redis.HSet(ctx, "Group", exam.GroupId.String(), v)
	}
levp:

	// TODO 查看是否是用户组组长
	if user.ID != group.LeaderId {
		response.Fail(ctx, nil, "非用户组组长，无法更新测试")
		return
	}

	// TODO 查看比赛是否已经进行过
	if time.Now().After(time.Time(exam.StartTime)) {
		response.Fail(ctx, nil, "测试已经进行过")
		return
	}
	// TODO 验证起始时间与终止时间是否合法
	if requestExam.StartTime.After(requestExam.EndTime) {
		response.Fail(ctx, nil, "起始时间大于了终止时间")
		return
	}
	if time.Now().After(time.Time(requestExam.StartTime)) {
		response.Fail(ctx, nil, "当前时间大于了起始时间")
		return
	}
	if time.Time(requestExam.EndTime).After(time.Now().Add(30 * 24 * time.Hour)) {
		response.Fail(ctx, nil, "终止时间不可设置为30日后")
		return
	}

	// TODO 新建测试
	examUpdate := model.Exam{
		Title:     requestExam.Title,
		Content:   requestExam.Content,
		ResLong:   requestExam.ResLong,
		ResShort:  requestExam.ResShort,
		StartTime: requestExam.StartTime,
		EndTime:   requestExam.EndTime,
		Type:      requestExam.Type,
	}

	// TODO 更新测试内容
	e.DB.Model(model.Exam{}).Where("id = (?)", id).Updates(examUpdate)

	// TODO 更新定时器
	util.TimerMap[exam.ID].Reset(time.Until(time.Time(examUpdate.StartTime)))

	// TODO 解码失败，删除字段
	e.Redis.HDel(ctx, "Exam", id)

	// TODO 成功
	response.Success(ctx, nil, "更新成功")
}

// @title    Show
// @description   查看一篇测试的内容
// @auth      MGAronya       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (e ExamController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var exam model.Exam

	// TODO 先尝试在redis中寻找
	if ok, _ := e.Redis.HExists(ctx, "Exam", id).Result(); ok {
		art, _ := e.Redis.HGet(ctx, "Exam", id).Result()
		if json.Unmarshal([]byte(art), &exam) == nil {
			goto leap
		} else {
			// TODO 解码失败，删除字段
			e.Redis.HDel(ctx, "Exam", id)
		}
	}

	// TODO 查看测试是否在数据库中存在
	if e.DB.Where("id = (?)", id).First(&exam).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}
	{
		// TODO 将测试存入redis供下次使用
		v, _ := json.Marshal(exam)
		e.Redis.HSet(ctx, "Exam", id, v)
	}
leap:
	// TODO 查看用户是否在指定用户组
	if e.DB.Where("user_id = (?) and group_id = (?)", user.ID, exam.GroupId).First(&model.UserList{}).Error != nil {
		response.Fail(ctx, nil, "不在指定用户组")
		return
	}

	response.Success(ctx, gin.H{"exam": exam}, "成功")
}

// @title    Delete
// @description   删除一篇测试
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (e ExamController) Delete(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	var exam model.Exam

	// TODO 查看测试是否存在
	if e.DB.Where("id = (?)", id).First(&exam).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}

	// TODO 判断当前用户是否为测试的作者
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看测试所在组
	// TODO 查找组
	var group model.Group

	// TODO 先看redis中是否存在
	if ok, _ := e.Redis.HExists(ctx, "Group", exam.GroupId.String()).Result(); ok {
		cate, _ := e.Redis.HGet(ctx, "Group", exam.GroupId.String()).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto levp
		} else {
			// TODO 移除损坏数据
			e.Redis.HDel(ctx, "Group", exam.GroupId.String())
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if e.DB.Where("id = (?)", exam.GroupId).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		e.Redis.HSet(ctx, "Group", exam.GroupId.String(), v)
	}
levp:

	// TODO 查看是否有操作测试的权力
	if user.ID != group.LeaderId && user.Level < 4 {
		response.Fail(ctx, nil, "测试不属于您，请勿非法操作")
		return
	}

	// TODO 删除测试
	e.DB.Delete(&exam)

	response.Success(ctx, nil, "删除成功")
	// TODO 解码失败，删除字段
	e.Redis.HDel(ctx, "Exam", id)
}

// @title    PageList
// @description   获取多篇测试
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (e ExamController) PageList(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看测试所在组
	// TODO 查找组
	var group model.Group

	// TODO 先看redis中是否存在
	if ok, _ := e.Redis.HExists(ctx, "Group", id).Result(); ok {
		cate, _ := e.Redis.HGet(ctx, "Group", id).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto levp
		} else {
			// TODO 移除损坏数据
			e.Redis.HDel(ctx, "Group", id)
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if e.DB.Where("id = (?)", id).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		e.Redis.HSet(ctx, "Group", id, v)
	}
levp:

	// TODO 查看是否是用户组中
	if e.DB.Where("user_id = (?) and group_id = (?)", user.ID, group.ID).First(&model.UserList{}).Error != nil {
		response.Fail(ctx, nil, "不在指定用户组中")
		return
	}

	// TODO 分页
	var exams []model.Exam

	// TODO 查找所有分页中可见的条目
	e.DB.Where("group_id = (?)", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&exams)

	var total int64
	e.DB.Where("group_id = (?)", id).Model(model.Exam{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"exams": exams, "total": total}, "成功")
}

// @title    ScoreShow
// @description   查看分数
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (e ExamController) ScoreShow(ctx *gin.Context) {
	// 获取path中的id
	user_id := ctx.Params.ByName("user_id")
	exam_id := ctx.Params.ByName("exam_id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找测试是否存在
	var exam model.Exam

	// TODO 先尝试在redis中寻找
	if ok, _ := e.Redis.HExists(ctx, "Exam", exam_id).Result(); ok {
		art, _ := e.Redis.HGet(ctx, "Exam", exam_id).Result()
		if json.Unmarshal([]byte(art), &exam) == nil {
			goto leap
		} else {
			// TODO 解码失败，删除字段
			e.Redis.HDel(ctx, "Exam", exam_id)
		}
	}

	// TODO 查看测试是否在数据库中存在
	if e.DB.Where("id = (?)", exam_id).First(&exam).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}
	{
		// TODO 将测试存入redis供下次使用
		v, _ := json.Marshal(exam)
		e.Redis.HSet(ctx, "Exam", exam_id, v)
	}
leap:

	// TODO 查看测试所在组
	// TODO 查找组
	var group model.Group

	// TODO 先看redis中是否存在
	if ok, _ := e.Redis.HExists(ctx, "Group", exam.GroupId.String()).Result(); ok {
		cate, _ := e.Redis.HGet(ctx, "Group", exam.GroupId.String()).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto levp
		} else {
			// TODO 移除损坏数据
			e.Redis.HDel(ctx, "Group", exam.GroupId.String())
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if e.DB.Where("id = (?)", exam.GroupId.String()).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		e.Redis.HSet(ctx, "Group", exam.GroupId.String(), v)
	}
levp:

	// TODO 查看是否是用户组中
	if e.DB.Where("user_id = (?) and group_id = (?)", user.ID, group.ID).First(&model.UserList{}).Error != nil {
		response.Fail(ctx, nil, "不在指定用户组中")
		return
	}

	// TODO 查看分数
	var examScore model.ExamScore
	e.DB.Where("user_id = (?) and exam_id = (?)", user_id, exam_id).First(&examScore)

	// TODO 返回数据
	response.Success(ctx, gin.H{"examScore": examScore}, "成功")
}

// @title    ScoreUpdate
// @description   修改分数
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (e ExamController) ScoreUpdate(ctx *gin.Context) {
	// 获取path中的id
	user_id := ctx.Params.ByName("user_id")
	exam_id := ctx.Params.ByName("exam_id")

	var requestScore vo.ScoreRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestScore); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找测试是否存在
	var exam model.Exam

	// TODO 先尝试在redis中寻找
	if ok, _ := e.Redis.HExists(ctx, "Exam", exam_id).Result(); ok {
		art, _ := e.Redis.HGet(ctx, "Exam", exam_id).Result()
		if json.Unmarshal([]byte(art), &exam) == nil {
			goto leap
		} else {
			// TODO 解码失败，删除字段
			e.Redis.HDel(ctx, "Exam", exam_id)
		}
	}

	// TODO 查看测试是否在数据库中存在
	if e.DB.Where("id = (?)", exam_id).First(&exam).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}
	{
		// TODO 将测试存入redis供下次使用
		v, _ := json.Marshal(exam)
		e.Redis.HSet(ctx, "Exam", exam_id, v)
	}
leap:

	// TODO 查看测试所在组
	// TODO 查找组
	var group model.Group

	// TODO 先看redis中是否存在
	if ok, _ := e.Redis.HExists(ctx, "Group", exam.GroupId.String()).Result(); ok {
		cate, _ := e.Redis.HGet(ctx, "Group", exam.GroupId.String()).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto levp
		} else {
			// TODO 移除损坏数据
			e.Redis.HDel(ctx, "Group", exam.GroupId.String())
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if e.DB.Where("id = (?)", exam.GroupId.String()).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		e.Redis.HSet(ctx, "Group", exam.GroupId.String(), v)
	}
levp:

	// TODO 查看是否是小组组长
	if user.ID != group.LeaderId {
		response.Fail(ctx, nil, "不是小组组长")
		return
	}

	e.DB.Model(model.ExamScore{}).Where("user_id = (?) and exam_id = (?)", user_id, exam_id).Update("score", requestScore.Score)

	// TODO 返回数据
	response.Success(ctx, nil, "成功")
}

// @title    ScoreList
// @description   分数列表
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (e ExamController) ScoreList(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	var requestScore vo.ScoreRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestScore); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 查找测试是否存在
	var exam model.Exam

	// TODO 先尝试在redis中寻找
	if ok, _ := e.Redis.HExists(ctx, "Exam", id).Result(); ok {
		art, _ := e.Redis.HGet(ctx, "Exam", id).Result()
		if json.Unmarshal([]byte(art), &exam) == nil {
			goto leap
		} else {
			// TODO 解码失败，删除字段
			e.Redis.HDel(ctx, "Exam", id)
		}
	}

	// TODO 查看测试是否在数据库中存在
	if e.DB.Where("id = (?)", id).First(&exam).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}
	{
		// TODO 将测试存入redis供下次使用
		v, _ := json.Marshal(exam)
		e.Redis.HSet(ctx, "Exam", id, v)
	}
leap:

	// TODO 查看测试所在组
	// TODO 查找组
	var group model.Group

	// TODO 先看redis中是否存在
	if ok, _ := e.Redis.HExists(ctx, "Group", exam.GroupId.String()).Result(); ok {
		cate, _ := e.Redis.HGet(ctx, "Group", exam.GroupId.String()).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto levp
		} else {
			// TODO 移除损坏数据
			e.Redis.HDel(ctx, "Group", exam.GroupId.String())
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if e.DB.Where("id = (?)", exam.GroupId.String()).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		e.Redis.HSet(ctx, "Group", exam.GroupId.String(), v)
	}
levp:

	// TODO 查看是否是用户组中
	if e.DB.Where("user_id = (?) and group_id = (?)", user.ID, group.ID).First(&model.UserList{}).Error != nil {
		response.Fail(ctx, nil, "不在指定用户组中")
		return
	}

	// TODO 分页
	var examScores []model.ExamScore

	// TODO 查找所有分页中可见的条目
	e.DB.Where("exam_id = (?)", id).Order("score desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&examScores)

	var total int64
	e.DB.Where("exam_id = (?)", id).Model(model.ExamScore{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"examScores": examScores, "total": total}, "成功")
}

// @title    NewExamController
// @description   新建一个IExamController
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   IExamController		返回一个IExamController用于调用各种函数
func NewExamController() IExamController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.Exam{})
	db.AutoMigrate(model.ExamScore{})
	return ExamController{DB: db, Redis: redis}
}

// @title    ExamTimer
// @description   建立一个测试开始定时器
// @auth      MGAronya       2022-9-16 12:23
// @param    competitionId uuid.UUID	比赛id
// @return   void
func ExamTimer(ctx *gin.Context, redis *redis.Client, db *gorm.DB, exam model.Exam) {

	util.TimerMap[exam.ID] = time.NewTimer(time.Until(time.Time(exam.StartTime)))
	// TODO 等待测试开始
	<-util.TimerMap[exam.ID].C
	// TODO 测试初始事项

	// TODO 创建测试结束定时器
	util.TimerMap[exam.ID] = time.NewTimer(time.Until(time.Time(exam.EndTime)))

	// TODO 等待测试结束
	<-util.TimerMap[exam.ID].C

	// TODO 分数初步统计
	var problemClozes []model.ProblemCloze
	db.Where("exam_id = (?)", exam.ID).Find(&problemClozes)

	var problemMCQss []model.ProblemMCQs
	db.Where("exam_id = (?)", exam.ID).Find(&problemMCQss)

	// TODO 使用map记录每个用户的分数
	userMap := make(map[uuid.UUID]uint)

	// TODO 计算分数
	for i := range problemClozes {
		var problemClozeSubmits []model.ProblemClozeSubmit
		db.Where("problem_cloze_id = (?)", problemClozes[i].ID).Find(&problemClozeSubmits)
		for j := range problemClozeSubmits {
			if ProblemClozeJudge(problemClozes[i].Answer, problemClozeSubmits[j].Answer) {
				userMap[problemClozeSubmits[j].UserId] += problemClozes[i].Score
			}
		}
	}
	for i := range problemMCQss {
		var problemMCQsSubmits []model.ProblemMCQsSubmit
		db.Where("problem_mcqs_id = (?)", problemMCQss[i].ID).Find(&problemMCQsSubmits)
		for j := range problemMCQsSubmits {
			if ProblemMCQsJudge(problemMCQss[i].Answer, problemMCQsSubmits[j].Answer) {
				userMap[problemMCQsSubmits[j].UserId] += problemMCQss[i].Score
			}
		}
	}
	// TODO 分数统计
	for i := range userMap {
		examScore := model.ExamScore{
			UserId: i,
			ExamId: exam.ID,
			Score:  userMap[i],
		}
		db.Create(&examScore)
	}
}
