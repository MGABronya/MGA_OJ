// @Title  ProblemMCQsController
// @Description  该文件提供关于操作选择题的各种方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package controller

import (
	"MGA_OJ/Interface"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"
	"MGA_OJ/vo"
	"encoding/json"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v9"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IProblemMCQsController			定义了选择题类接口
type IProblemMCQsController interface {
	Interface.RestInterface      // 包含增删查改功能
	Submit(ctx *gin.Context)     // 提交测试
	ShowSubmit(ctx *gin.Context) // 查看提交情况
	SubmitList(ctx *gin.Context) // 查看指定问题指定用户的提交列表
}

// ProblemMCQsController			定义了选择题工具类
type ProblemMCQsController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Create
// @description   创建一篇选择题
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemMCQsController) Create(ctx *gin.Context) {
	var requestProblemMCQs vo.ProblemMCQsRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestProblemMCQs); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 查看是否是该测试所属组的组长
	// TODO 查找测试
	var exam model.Exam

	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "Exam", id).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Exam", id).Result()
		if json.Unmarshal([]byte(cate), &exam) == nil {
			goto levp
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Exam", id)
		}
	}

	// TODO 查看测试是否在数据库中存在
	if p.DB.Where("id = (?)", id).First(&exam).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(exam)
		p.Redis.HSet(ctx, "Exam", id, v)
	}
levp:
	// TODO 查找组
	var group model.Group

	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "Group", exam.GroupId.String()).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Group", exam.GroupId.String()).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto levp
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Group", exam.GroupId.String())
		}
	}

	// TODO 查看测试是否在数据库中存在
	if p.DB.Where("id = (?)", exam.GroupId).First(&group).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		p.Redis.HSet(ctx, "Group", exam.GroupId.String(), v)
	}

	// TODO 查看是否是用户组组长
	if user.ID != group.LeaderId {
		response.Fail(ctx, nil, "非用户组组长，无法布置选择题")
		return
	}

	// TODO 创建选择题
	ProblemMCQs := model.ProblemMCQs{
		Description: requestProblemMCQs.Description,
		ResLong:     requestProblemMCQs.ResLong,
		ResShort:    requestProblemMCQs.ResShort,
		Answer:      requestProblemMCQs.Answer,
		Score:       requestProblemMCQs.Score,
		UserId:      user.ID,
		ExamId:      exam.ID,
	}

	// TODO 插入数据
	if err := p.DB.Create(&ProblemMCQs).Error; err != nil {
		response.Fail(ctx, nil, "选择题上传出错，数据验证有误")
		return
	}

	// TODO 成功
	response.Success(ctx, gin.H{"problemMCQs": ProblemMCQs}, "创建成功")
}

// @title    Update
// @description   更新一篇选择题的内容
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemMCQsController) Update(ctx *gin.Context) {
	var requestProblemMCQs vo.ProblemMCQsRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestProblemMCQs); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应选择题
	id := ctx.Params.ByName("id")

	var problemMCQs model.ProblemMCQs

	if p.DB.Where("id = (?)", id).First(&problemMCQs).Error != nil {
		response.Fail(ctx, nil, "选择题不存在")
		return
	}

	// TODO 查看选择题所在测试
	var exam model.Exam
	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "Exam", problemMCQs.ExamId.String()).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Exam", problemMCQs.ExamId.String()).Result()
		if json.Unmarshal([]byte(cate), &exam) == nil {
			goto leap
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Exam", problemMCQs.ExamId.String())
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if p.DB.Where("id = (?)", problemMCQs.ExamId).First(&exam).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(exam)
		p.Redis.HSet(ctx, "Exam", problemMCQs.ExamId.String(), v)
	}
leap:
	// TODO 查找组
	var group model.Group

	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "Group", exam.GroupId.String()).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Group", exam.GroupId.String()).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto levp
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Group", exam.GroupId.String())
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if p.DB.Where("id = (?)", exam.GroupId).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		p.Redis.HSet(ctx, "Group", exam.GroupId.String(), v)
	}
levp:

	// TODO 查看是否是用户组组长
	if user.ID != group.LeaderId {
		response.Fail(ctx, nil, "非用户组组长，无法布置选择题")
		return
	}

	// TODO 新建选择题
	ProblemMCQsUpdate := model.ProblemMCQs{
		Description: requestProblemMCQs.Description,
		ResLong:     requestProblemMCQs.ResLong,
		ResShort:    requestProblemMCQs.ResShort,
		Answer:      requestProblemMCQs.Answer,
		Score:       requestProblemMCQs.Score,
	}

	// TODO 更新选择题内容
	p.DB.Model(model.ProblemMCQs{}).Where("id = (?)", id).Updates(ProblemMCQsUpdate)

	// TODO 解码失败，删除字段
	p.Redis.HDel(ctx, "ProblemMCQs", id)

	p.DB.Where("id = (?)", id).First(&problemMCQs)

	// TODO 成功
	response.Success(ctx, gin.H{"problemMCQs": problemMCQs}, "更新成功")
}

// @title    Show
// @description   查看一篇选择题的内容
// @auth      MGAronya       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemMCQsController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var problemMCQs model.ProblemMCQs

	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "ProblemMCQs", id).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "ProblemMCQs", id).Result()
		if json.Unmarshal([]byte(art), &problemMCQs) == nil {
			goto leap
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "ProblemMCQs", id)
		}
	}

	// TODO 查看选择题是否在数据库中存在
	if p.DB.Where("id = (?)", id).First(&problemMCQs).Error != nil {
		response.Fail(ctx, nil, "选择题不存在")
		return
	}
	{
		// TODO 将选择题存入redis供下次使用
		v, _ := json.Marshal(problemMCQs)
		p.Redis.HSet(ctx, "ProblemMCQs", id, v)
	}
leap:
	// TODO 查看题目所属的测试
	var exam model.Exam
	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "Exam", problemMCQs.ExamId.String()).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "Exam", problemMCQs.ExamId.String()).Result()
		if json.Unmarshal([]byte(art), &exam) == nil {
			goto levp
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "Exam", problemMCQs.ExamId.String())
		}
	}

	// TODO 查看选择题是否在数据库中存在
	if p.DB.Where("id = (?)", problemMCQs.ExamId).First(&exam).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}
	{
		// TODO 将选择题存入redis供下次使用
		v, _ := json.Marshal(exam)
		p.Redis.HSet(ctx, "Exam", problemMCQs.ExamId.String(), v)
	}
levp:
	// TODO 查看测试所在的用户组
	var group model.Group
	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "Group", exam.GroupId.String()).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "Group", exam.GroupId.String()).Result()
		if json.Unmarshal([]byte(art), &group) == nil {
			goto letp
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "Group", exam.GroupId.String())
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if p.DB.Where("id = (?)", exam.GroupId).First(&group).Error != nil {
		response.Fail(ctx, nil, "小组不存在")
		return
	}
	{
		// TODO 将选择题存入redis供下次使用
		v, _ := json.Marshal(group)
		p.Redis.HSet(ctx, "Group", exam.GroupId.String(), v)
	}
letp:
	// TODO 查看用户是否在指定用户组
	if p.DB.Where("user_id = (?) and group_id = (?)", user.ID, exam.GroupId).First(&model.UserList{}).Error != nil {
		response.Fail(ctx, nil, "不在指定用户组")
		return
	}
	// TODO 查看测试是否已经开始
	if exam.StartTime.After(model.Time(time.Now())) {
		response.Fail(ctx, nil, "测试未开始")
		return
	}

	// TODO 查看用户是否是组长
	if group.LeaderId != user.ID {
		problemMCQs.Answer = ""
	}

	response.Success(ctx, gin.H{"problemMCQs": problemMCQs}, "成功")
}

// @title    Delete
// @description   删除一篇选择题
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemMCQsController) Delete(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var problemMCQs model.ProblemMCQs

	// TODO 查看选择题是否存在
	if p.DB.Where("id = (?)", id).First(&problemMCQs).Error != nil {
		response.Fail(ctx, nil, "选择题不存在")
		return
	}

	// TODO 查看题目所属的测试
	var exam model.Exam
	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "Exam", problemMCQs.ExamId.String()).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "Exam", problemMCQs.ExamId.String()).Result()
		if json.Unmarshal([]byte(art), &exam) == nil {
			goto levp
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "Exam", problemMCQs.ExamId.String())
		}
	}

	// TODO 查看测试是否在数据库中存在
	if p.DB.Where("id = (?)", problemMCQs.ExamId).First(&exam).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}
	{
		// TODO 将选择题存入redis供下次使用
		v, _ := json.Marshal(exam)
		p.Redis.HSet(ctx, "Exam", problemMCQs.ExamId.String(), v)
	}
levp:
	// TODO 查看测试所在的用户组
	var group model.Group
	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "Group", exam.GroupId.String()).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "Group", exam.GroupId.String()).Result()
		if json.Unmarshal([]byte(art), &group) == nil {
			goto letp
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "Group", exam.GroupId.String())
		}
	}

	// TODO 查看选择题是否在数据库中存在
	if p.DB.Where("id = (?)", exam.GroupId).First(&group).Error != nil {
		response.Fail(ctx, nil, "小组不存在")
		return
	}
	{
		// TODO 将选择题存入redis供下次使用
		v, _ := json.Marshal(group)
		p.Redis.HSet(ctx, "Group", exam.GroupId.String(), v)
	}
letp:

	// TODO 查看用户是否是组长
	if group.LeaderId != user.ID {
		response.Fail(ctx, nil, "不是组长，无法删除题目")
		return
	}
	// TODO 删除选择题
	p.DB.Delete(&problemMCQs)

	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description   获取多篇选择题
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemMCQsController) PageList(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看题目所属的测试
	var exam model.Exam
	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "Exam", id).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "Exam", id).Result()
		if json.Unmarshal([]byte(art), &exam) == nil {
			goto levp
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "Exam", id)
		}
	}

	// TODO 查看选择题是否在数据库中存在
	if p.DB.Where("id = (?)", id).First(&exam).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}
	{
		// TODO 将选择题存入redis供下次使用
		v, _ := json.Marshal(exam)
		p.Redis.HSet(ctx, "Exam", id, v)
	}
levp:
	// TODO 查看测试所在的用户组
	var group model.Group
	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "Group", exam.GroupId.String()).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "Group", exam.GroupId.String()).Result()
		if json.Unmarshal([]byte(art), &group) == nil {
			goto letp
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "Group", exam.GroupId.String())
		}
	}

	// TODO 查看选择题是否在数据库中存在
	if p.DB.Where("id = (?)", exam.GroupId).First(&group).Error != nil {
		response.Fail(ctx, nil, "小组不存在")
		return
	}
	{
		// TODO 将选择题存入redis供下次使用
		v, _ := json.Marshal(group)
		p.Redis.HSet(ctx, "Group", exam.GroupId.String(), v)
	}
letp:
	// TODO 查看用户是否在指定用户组
	if p.DB.Where("user_id = (?) and group_id = (?)", user.ID, exam.GroupId).First(&model.UserList{}).Error != nil {
		response.Fail(ctx, nil, "不在指定用户组")
		return
	}
	// TODO 查看测试是否已经开始
	if exam.StartTime.After(model.Time(time.Now())) {
		response.Fail(ctx, nil, "测试未开始")
		return
	}

	// TODO 分页
	var problemMCQss []model.ProblemMCQs

	// TODO 查找所有分页中可见的条目
	p.DB.Where("exam_id = (?)", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&problemMCQss)

	var total int64
	p.DB.Where("exam_id = (?)", id).Model(model.ProblemMCQs{}).Count(&total)

	// TODO 查看用户是否是组长
	if group.LeaderId != user.ID {
		for i := range problemMCQss {
			problemMCQss[i].Answer = ""
		}
	}

	// TODO 返回数据
	response.Success(ctx, gin.H{"problemMCQss": problemMCQss, "total": total}, "成功")
}

// @title    Submit
// @description   测试提交
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemMCQsController) Submit(ctx *gin.Context) {
	var requestProblemAnsMCQs vo.ProblemMCQsAnsRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestProblemAnsMCQs); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var problemMCQs model.ProblemMCQs

	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "ProblemMCQs", id).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "ProblemMCQs", id).Result()
		if json.Unmarshal([]byte(art), &problemMCQs) == nil {
			goto leap
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "ProblemMCQs", id)
		}
	}

	// TODO 查看选择题是否在数据库中存在
	if p.DB.Where("id = (?)", id).First(&problemMCQs).Error != nil {
		response.Fail(ctx, nil, "选择题不存在")
		return
	}
	{
		// TODO 将选择题存入redis供下次使用
		v, _ := json.Marshal(problemMCQs)
		p.Redis.HSet(ctx, "ProblemMCQs", id, v)
	}
leap:
	// TODO 查看题目所属的测试
	var exam model.Exam
	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "Exam", problemMCQs.ExamId.String()).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "Exam", problemMCQs.ExamId.String()).Result()
		if json.Unmarshal([]byte(art), &exam) == nil {
			goto levp
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "Exam", problemMCQs.ExamId.String())
		}
	}

	// TODO 查看选择题是否在数据库中存在
	if p.DB.Where("id = (?)", problemMCQs.ExamId).First(&exam).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}
	{
		// TODO 将选择题存入redis供下次使用
		v, _ := json.Marshal(exam)
		p.Redis.HSet(ctx, "Exam", problemMCQs.ExamId.String(), v)
	}
levp:
	// TODO 查看测试是否已经过期
	if time.Now().After(time.Time(exam.EndTime)) {
		response.Fail(ctx, nil, "测试已过期")
		return
	}
	// TODO 查看测试所在的用户组
	var group model.Group
	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "Group", exam.GroupId.String()).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "Group", exam.GroupId.String()).Result()
		if json.Unmarshal([]byte(art), &group) == nil {
			goto letp
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "Group", exam.GroupId.String())
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if p.DB.Where("id = (?)", exam.GroupId).First(&group).Error != nil {
		response.Fail(ctx, nil, "小组不存在")
		return
	}
	{
		// TODO 将选择题存入redis供下次使用
		v, _ := json.Marshal(group)
		p.Redis.HSet(ctx, "Group", exam.GroupId.String(), v)
	}
letp:
	// TODO 查看用户是否在指定用户组
	if p.DB.Where("user_id = (?) and group_id = (?)", user.ID, exam.GroupId).First(&model.UserList{}).Error != nil {
		response.Fail(ctx, nil, "不在指定用户组")
		return
	}
	// TODO 查看提交是否通过
	var score uint
	score = 0
	if exam.Type == "IOI" && ProblemMCQsJudge(problemMCQs.Answer, requestProblemAnsMCQs.Answer) {
		score = problemMCQs.Score
	}

	// TODO 查看先前是否由提交，如为第一次提交，记录分数
	if p.DB.Where("user_id = (?) and problem_mcqs_id = (?)").First(&model.ProblemMCQsSubmit{}).Error != nil {
		var examScore model.ExamScore
		if p.DB.Where("user_id = (?) and exam_id = (?)", user.ID, exam.ID).First(&examScore).Error != nil {
			examScore.UserId = user.ID
			examScore.ExamId = exam.ID
			examScore.Score += score
			p.DB.Create(&examScore)
		} else {
			p.DB.Model(model.ExamScore{}).Where("user_id = (?) and exam_id = (?)", user.ID, exam.ID).Update("score", examScore.Score+score)
		}
	} else if exam.Type == "IO" {
		response.Fail(ctx, nil, "已经提交，不可修改")
		return
	}

	// TODO 创建提交
	problemMCQsSubmit := model.ProblemMCQsSubmit{
		UserId:        user.ID,
		ProblemMCQsId: problemMCQs.ID,
		Answer:        requestProblemAnsMCQs.Answer,
		Score:         score,
	}

	// TODO 插入数据
	p.DB.Create(&problemMCQsSubmit)

	// TODO 返回数据
	response.Success(ctx, gin.H{"problemMCQsSubmit": problemMCQsSubmit}, "成功")
}

// @title    ShowSubmit
// @description   获取提交情况
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemMCQsController) ShowSubmit(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var problemMCQsSubmit model.ProblemMCQsSubmit

	if p.DB.Where("id = (?)", id).First(&problemMCQsSubmit).Error != nil {
		response.Fail(ctx, nil, "提交不存在")
		return
	}

	// TODO 查看提交所属的题目
	var problemMCQs model.ProblemMCQs

	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "ProblemMCQs", problemMCQsSubmit.ProblemMCQsId.String()).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "ProblemMCQs", problemMCQsSubmit.ProblemMCQsId.String()).Result()
		if json.Unmarshal([]byte(art), &problemMCQs) == nil {
			goto leap
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "ProblemMCQs", problemMCQsSubmit.ProblemMCQsId.String())
		}
	}

	// TODO 查看选择题是否在数据库中存在
	if p.DB.Where("id = (?)", problemMCQsSubmit.ProblemMCQsId).First(&problemMCQs).Error != nil {
		response.Fail(ctx, nil, "选择题不存在")
		return
	}
	{
		// TODO 将选择题存入redis供下次使用
		v, _ := json.Marshal(problemMCQs)
		p.Redis.HSet(ctx, "ProblemMCQs", id, v)
	}
leap:
	// TODO 查看题目所属的测试
	var exam model.Exam
	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "Exam", problemMCQs.ExamId.String()).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "Exam", problemMCQs.ExamId.String()).Result()
		if json.Unmarshal([]byte(art), &exam) == nil {
			goto levp
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "Exam", problemMCQs.ExamId.String())
		}
	}

	// TODO 查看选择题是否在数据库中存在
	if p.DB.Where("id = (?)", problemMCQs.ExamId).First(&exam).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}
	{
		// TODO 将选择题存入redis供下次使用
		v, _ := json.Marshal(exam)
		p.Redis.HSet(ctx, "Exam", problemMCQs.ExamId.String(), v)
	}
levp:
	// TODO 查看测试所在的用户组
	var group model.Group
	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "Group", exam.GroupId.String()).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "Group", exam.GroupId.String()).Result()
		if json.Unmarshal([]byte(art), &group) == nil {
			goto letp
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "Group", exam.GroupId.String())
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if p.DB.Where("id = (?)", exam.GroupId).First(&group).Error != nil {
		response.Fail(ctx, nil, "小组不存在")
		return
	}
	{
		// TODO 将选择题存入redis供下次使用
		v, _ := json.Marshal(group)
		p.Redis.HSet(ctx, "Group", exam.GroupId.String(), v)
	}
letp:

	// TODO 查看用户是否是组长或提交者
	if group.LeaderId != user.ID && user.ID != problemMCQsSubmit.UserId {
		response.Fail(ctx, nil, "权限不足")
		return
	}
	// TODO 返回数据
	response.Success(ctx, gin.H{"problemMCQsSubmit": problemMCQsSubmit}, "成功")

}

// @title    SubmitList
// @description   获取提交情况
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemMCQsController) SubmitList(ctx *gin.Context) {
	// 获取path中的id
	user_id := ctx.Params.ByName("user_id")
	problem_id := ctx.Params.ByName("problem_id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看提交所属的题目
	var problemMCQs model.ProblemMCQs

	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "ProblemMCQs", problem_id).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "ProblemMCQs", problem_id).Result()
		if json.Unmarshal([]byte(art), &problemMCQs) == nil {
			goto leap
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "ProblemMCQs", problem_id)
		}
	}

	// TODO 查看选择题是否在数据库中存在
	if p.DB.Where("id = (?)", problem_id).First(&problemMCQs).Error != nil {
		response.Fail(ctx, nil, "选择题不存在")
		return
	}
	{
		// TODO 将选择题存入redis供下次使用
		v, _ := json.Marshal(problemMCQs)
		p.Redis.HSet(ctx, "ProblemMCQs", problem_id, v)
	}
leap:
	// TODO 查看题目所属的测试
	var exam model.Exam
	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "Exam", problemMCQs.ExamId.String()).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "Exam", problemMCQs.ExamId.String()).Result()
		if json.Unmarshal([]byte(art), &exam) == nil {
			goto levp
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "Exam", problemMCQs.ExamId.String())
		}
	}

	// TODO 查看选择题是否在数据库中存在
	if p.DB.Where("id = (?)", problemMCQs.ExamId).First(&exam).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}
	{
		// TODO 将选择题存入redis供下次使用
		v, _ := json.Marshal(exam)
		p.Redis.HSet(ctx, "Exam", problemMCQs.ExamId.String(), v)
	}
levp:
	// TODO 查看测试所在的用户组
	var group model.Group
	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "Group", exam.GroupId.String()).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "Group", exam.GroupId.String()).Result()
		if json.Unmarshal([]byte(art), &group) == nil {
			goto letp
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "Group", exam.GroupId.String())
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if p.DB.Where("id = (?)", exam.GroupId).First(&group).Error != nil {
		response.Fail(ctx, nil, "小组不存在")
		return
	}
	{
		// TODO 将选择题存入redis供下次使用
		v, _ := json.Marshal(group)
		p.Redis.HSet(ctx, "Group", exam.GroupId.String(), v)
	}
letp:

	// TODO 查看用户是否是组长或提交者
	if group.LeaderId != user.ID && user.ID.String() != user_id {
		response.Fail(ctx, nil, "权限不足")
		return
	}
	var problemMCQsSubmits []model.ProblemMCQsSubmit
	p.DB.Where("problem_mcqs_id = (?) and user_id = (?)", problemMCQs.ID, user_id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&problemMCQs)

	var total int64
	p.DB.Where("problem_mcqs_id = (?) and user_id = (?)", problemMCQs.ID, user_id).Model(model.ProblemMCQsSubmit{}).Count(&total)
	// TODO 返回数据
	response.Success(ctx, gin.H{"problemMCQsSubmits": problemMCQsSubmits, "total": total}, "成功")
}

// @title    NewProblemMCQsController
// @description   新建一个IProblemMCQsController
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   IProblemMCQsController		返回一个IProblemMCQsController用于调用各种函数
func NewProblemMCQsController() IProblemMCQsController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.ProblemMCQs{})
	db.AutoMigrate(model.ProblemMCQsSubmit{})
	return ProblemMCQsController{DB: db, Redis: redis}
}

// @title    ProblemMCQsJudge
// @description   检查提交答案是否正确
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func ProblemMCQsJudge(source string, submit string) bool {
	// TODO 转小写
	source = strings.ToLower(source)
	// TODO 进行简单排序
	b := []byte(source)
	sort.Slice(b, func(i, j int) bool {
		return b[i] < b[j]
	})
	source = string(b)

	// TODO 转小写
	submit = strings.ToLower(submit)
	// TODO 进行简单排序
	b = []byte(submit)
	sort.Slice(b, func(i, j int) bool {
		return b[i] < b[j]
	})
	submit = string(b)

	return submit == source
}
