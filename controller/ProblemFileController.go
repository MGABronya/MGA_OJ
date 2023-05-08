// @Title  ProblemFileController
// @Description  该文件提供关于操作文件题的各种方法
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
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IProblemFileController			定义了文件题类接口
type IProblemFileController interface {
	Interface.RestInterface      // 包含增删查改功能
	Submit(ctx *gin.Context)     // 提交测试
	ShowSubmit(ctx *gin.Context) // 查看提交情况
	SubmitList(ctx *gin.Context) // 查看指定问题指定用户的提交列表
}

// ProblemFileController			定义了文件题工具类
type ProblemFileController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Create
// @description   创建一篇文件题
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemFileController) Create(ctx *gin.Context) {
	var requestProblemFile vo.ProblemFileRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestProblemFile); err != nil {
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
	if p.DB.Where("id = ?", id).First(&exam).Error != nil {
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
	if p.DB.Where("id = ?", exam.GroupId).First(&group).Error != nil {
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
		response.Fail(ctx, nil, "非用户组组长，无法布置文件题")
		return
	}

	// TODO 创建文件题
	ProblemFile := model.ProblemFile{
		Description: requestProblemFile.Description,
		ResLong:     requestProblemFile.ResLong,
		ResShort:    requestProblemFile.ResShort,
		Score:       requestProblemFile.Score,
		UserId:      user.ID,
		ExamId:      exam.ID,
	}

	// TODO 插入数据
	if err := p.DB.Create(&ProblemFile).Error; err != nil {
		response.Fail(ctx, nil, "文件题上传出错，数据验证有误")
		return
	}

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    Update
// @description   更新一篇文件题的内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemFileController) Update(ctx *gin.Context) {
	var requestProblemFile vo.ProblemFileRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestProblemFile); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应文件题
	id := ctx.Params.ByName("id")

	var problemFile model.ProblemFile

	if p.DB.Where("id = ?", id).First(&problemFile) != nil {
		response.Fail(ctx, nil, "文件题不存在")
		return
	}

	// TODO 查看文件题所在测试
	var exam model.Exam
	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "Exam", problemFile.ExamId.String()).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Exam", problemFile.ExamId.String()).Result()
		if json.Unmarshal([]byte(cate), &exam) == nil {
			goto leap
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Exam", problemFile.ExamId.String())
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if p.DB.Where("id = ?", problemFile.ExamId).First(&exam).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(exam)
		p.Redis.HSet(ctx, "Exam", problemFile.ExamId.String(), v)
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
	if p.DB.Where("id = ?", exam.GroupId).First(&group).Error != nil {
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
		response.Fail(ctx, nil, "非用户组组长，无法布置文件题")
		return
	}

	// TODO 新建文件题
	ProblemFileUpdate := model.ProblemFile{
		Description: requestProblemFile.Description,
		ResLong:     requestProblemFile.ResLong,
		ResShort:    requestProblemFile.ResShort,
		Score:       requestProblemFile.Score,
	}

	// TODO 更新文件题内容
	p.DB.Model(model.ProblemFile{}).Where("id = ?", id).Updates(ProblemFileUpdate)

	// TODO 解码失败，删除字段
	p.Redis.HDel(ctx, "ProblemFile", id)

	// TODO 成功
	response.Success(ctx, nil, "更新成功")
}

// @title    Show
// @description   查看一篇文件题的内容
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemFileController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var problemFile model.ProblemFile

	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "ProblemFile", id).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "ProblemFile", id).Result()
		if json.Unmarshal([]byte(art), &problemFile) == nil {
			goto leap
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "ProblemFile", id)
		}
	}

	// TODO 查看文件题是否在数据库中存在
	if p.DB.Where("id = ?", id).First(&problemFile).Error != nil {
		response.Fail(ctx, nil, "文件题不存在")
		return
	}
	{
		// TODO 将文件题存入redis供下次使用
		v, _ := json.Marshal(problemFile)
		p.Redis.HSet(ctx, "ProblemFile", id, v)
	}
leap:
	// TODO 查看题目所属的测试
	var exam model.Exam
	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "Exam", problemFile.ExamId.String()).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "Exam", problemFile.ExamId.String()).Result()
		if json.Unmarshal([]byte(art), &exam) == nil {
			goto levp
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "Exam", problemFile.ExamId.String())
		}
	}

	// TODO 查看文件题是否在数据库中存在
	if p.DB.Where("id = ?", problemFile.ExamId).First(&exam).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}
	{
		// TODO 将文件题存入redis供下次使用
		v, _ := json.Marshal(exam)
		p.Redis.HSet(ctx, "Exam", problemFile.ExamId.String(), v)
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
	if p.DB.Where("id = ?", exam.GroupId).First(&group).Error != nil {
		response.Fail(ctx, nil, "小组不存在")
		return
	}
	{
		// TODO 将文件题存入redis供下次使用
		v, _ := json.Marshal(group)
		p.Redis.HSet(ctx, "Group", exam.GroupId.String(), v)
	}
letp:
	// TODO 查看用户是否在指定用户组
	if p.DB.Where("user_id = ? and group_id = ?", user.ID, exam.GroupId).First(&model.UserList{}).Error != nil {
		response.Fail(ctx, nil, "不在指定用户组")
		return
	}
	// TODO 查看测试是否已经开始
	if exam.StartTime.After(model.Time(time.Now())) {
		response.Fail(ctx, nil, "测试未开始")
		return
	}

	response.Success(ctx, gin.H{"problemFile": problemFile}, "成功")
}

// @title    Delete
// @description   删除一篇文件题
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemFileController) Delete(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var problemFile model.ProblemFile

	// TODO 查看文件题是否存在
	if p.DB.Where("id = ?", id).First(&problemFile).Error != nil {
		response.Fail(ctx, nil, "文件题不存在")
		return
	}

	// TODO 查看题目所属的测试
	var exam model.Exam
	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "Exam", problemFile.ExamId.String()).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "Exam", problemFile.ExamId.String()).Result()
		if json.Unmarshal([]byte(art), &exam) == nil {
			goto levp
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "Exam", problemFile.ExamId.String())
		}
	}

	// TODO 查看测试是否在数据库中存在
	if p.DB.Where("id = ?", problemFile.ExamId).First(&exam).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}
	{
		// TODO 将文件题存入redis供下次使用
		v, _ := json.Marshal(exam)
		p.Redis.HSet(ctx, "Exam", problemFile.ExamId.String(), v)
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

	// TODO 查看文件题是否在数据库中存在
	if p.DB.Where("id = ?", exam.GroupId).First(&group).Error != nil {
		response.Fail(ctx, nil, "小组不存在")
		return
	}
	{
		// TODO 将文件题存入redis供下次使用
		v, _ := json.Marshal(group)
		p.Redis.HSet(ctx, "Group", exam.GroupId.String(), v)
	}
letp:

	// TODO 查看用户是否是组长
	if group.LeaderId != user.ID {
		response.Fail(ctx, nil, "不是组长，无法删除题目")
		return
	}
	// TODO 删除文件题
	p.DB.Delete(&problemFile)

	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description   获取多篇文件题
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemFileController) PageList(ctx *gin.Context) {
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

	// TODO 查看文件题是否在数据库中存在
	if p.DB.Where("id = ?", id).First(&exam).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}
	{
		// TODO 将文件题存入redis供下次使用
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

	// TODO 查看文件题是否在数据库中存在
	if p.DB.Where("id = ?", exam.GroupId).First(&group).Error != nil {
		response.Fail(ctx, nil, "小组不存在")
		return
	}
	{
		// TODO 将文件题存入redis供下次使用
		v, _ := json.Marshal(group)
		p.Redis.HSet(ctx, "Group", exam.GroupId.String(), v)
	}
letp:
	// TODO 查看用户是否在指定用户组
	if p.DB.Where("user_id = ? and group_id = ?", user.ID, exam.GroupId).First(&model.UserList{}).Error != nil {
		response.Fail(ctx, nil, "不在指定用户组")
		return
	}
	// TODO 查看测试是否已经开始
	if exam.StartTime.After(model.Time(time.Now())) {
		response.Fail(ctx, nil, "测试未开始")
		return
	}

	// TODO 分页
	var problemFiles []model.ProblemFile

	// TODO 查找所有分页中可见的条目
	p.DB.Where("exam_id = ?", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&problemFiles)

	var total int64
	p.DB.Where("exam_id = ?", id).Model(model.ProblemFile{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"problemFiles": problemFiles, "total": total}, "成功")
}

// @title    Submit
// @description   测试提交
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemFileController) Submit(ctx *gin.Context) {
	var requestProblemAnsFile vo.ProblemFileAnsRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestProblemAnsFile); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var problemFile model.ProblemFile

	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "ProblemFile", id).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "ProblemFile", id).Result()
		if json.Unmarshal([]byte(art), &problemFile) == nil {
			goto leap
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "ProblemFile", id)
		}
	}

	// TODO 查看文件题是否在数据库中存在
	if p.DB.Where("id = ?", id).First(&problemFile).Error != nil {
		response.Fail(ctx, nil, "文件题不存在")
		return
	}
	{
		// TODO 将文件题存入redis供下次使用
		v, _ := json.Marshal(problemFile)
		p.Redis.HSet(ctx, "ProblemFile", id, v)
	}
leap:
	// TODO 查看题目所属的测试
	var exam model.Exam
	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "Exam", problemFile.ExamId.String()).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "Exam", problemFile.ExamId.String()).Result()
		if json.Unmarshal([]byte(art), &exam) == nil {
			goto levp
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "Exam", problemFile.ExamId.String())
		}
	}

	// TODO 查看文件题是否在数据库中存在
	if p.DB.Where("id = ?", problemFile.ExamId).First(&exam).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}
	{
		// TODO 将文件题存入redis供下次使用
		v, _ := json.Marshal(exam)
		p.Redis.HSet(ctx, "Exam", problemFile.ExamId.String(), v)
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
	if p.DB.Where("id = ?", exam.GroupId).First(&group).Error != nil {
		response.Fail(ctx, nil, "小组不存在")
		return
	}
	{
		// TODO 将文件题存入redis供下次使用
		v, _ := json.Marshal(group)
		p.Redis.HSet(ctx, "Group", exam.GroupId.String(), v)
	}
letp:
	// TODO 查看用户是否在指定用户组
	if p.DB.Where("user_id = ? and group_id = ?", user.ID, exam.GroupId).First(&model.UserList{}).Error != nil {
		response.Fail(ctx, nil, "不在指定用户组")
		return
	}
	// TODO 查看提交是否通过
	var score uint
	score = 0

	// TODO 查看先前是否由提交，如为第一次提交，记录分数
	if p.DB.Where("user_id = ? and problem_file_id = ?").First(&model.ProblemFileSubmit{}).Error != nil {
		var examScore model.ExamScore
		if p.DB.Where("user_id = ? and exam_id = ?", user.ID, exam.ID).First(&examScore).Error != nil {
			examScore.UserId = user.ID
			examScore.ExamId = exam.ID
			examScore.Score += score
			p.DB.Create(&examScore)
		} else {
			p.DB.Model(model.ExamScore{}).Where("user_id = ? and exam_id = ?", user.ID, exam.ID).Update("score", examScore.Score+score)
		}
	} else if exam.Type == "IO" {
		response.Fail(ctx, nil, "已经提交，不可修改")
		return
	}

	// TODO 创建提交
	problemFileSubmit := model.ProblemFileSubmit{
		UserId:        user.ID,
		ProblemFileId: problemFile.ID,
		Answer:        requestProblemAnsFile.Answer,
		Score:         score,
	}

	// TODO 插入数据
	p.DB.Create(&problemFileSubmit)

	// TODO 返回数据
	response.Success(ctx, gin.H{"problemFileSubmit": problemFileSubmit}, "成功")
}

// @title    ShowSubmit
// @description   获取提交情况
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemFileController) ShowSubmit(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var problemFileSubmit model.ProblemFileSubmit

	if p.DB.Where("id = ?", id).First(&problemFileSubmit).Error != nil {
		response.Fail(ctx, nil, "提交不存在")
		return
	}

	// TODO 查看提交所属的题目
	var problemFile model.ProblemFile

	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "ProblemFile", problemFileSubmit.ProblemFileId.String()).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "ProblemFile", problemFileSubmit.ProblemFileId.String()).Result()
		if json.Unmarshal([]byte(art), &problemFile) == nil {
			goto leap
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "ProblemFile", problemFileSubmit.ProblemFileId.String())
		}
	}

	// TODO 查看文件题是否在数据库中存在
	if p.DB.Where("id = ?", problemFileSubmit.ProblemFileId).First(&problemFile).Error != nil {
		response.Fail(ctx, nil, "文件题不存在")
		return
	}
	{
		// TODO 将文件题存入redis供下次使用
		v, _ := json.Marshal(problemFile)
		p.Redis.HSet(ctx, "ProblemFile", id, v)
	}
leap:
	// TODO 查看题目所属的测试
	var exam model.Exam
	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "Exam", problemFile.ExamId.String()).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "Exam", problemFile.ExamId.String()).Result()
		if json.Unmarshal([]byte(art), &exam) == nil {
			goto levp
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "Exam", problemFile.ExamId.String())
		}
	}

	// TODO 查看文件题是否在数据库中存在
	if p.DB.Where("id = ?", problemFile.ExamId).First(&exam).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}
	{
		// TODO 将文件题存入redis供下次使用
		v, _ := json.Marshal(exam)
		p.Redis.HSet(ctx, "Exam", problemFile.ExamId.String(), v)
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
	if p.DB.Where("id = ?", exam.GroupId).First(&group).Error != nil {
		response.Fail(ctx, nil, "小组不存在")
		return
	}
	{
		// TODO 将文件题存入redis供下次使用
		v, _ := json.Marshal(group)
		p.Redis.HSet(ctx, "Group", exam.GroupId.String(), v)
	}
letp:

	// TODO 查看用户是否是组长或提交者
	if group.LeaderId != user.ID && user.ID != problemFileSubmit.UserId {
		response.Fail(ctx, nil, "权限不足")
		return
	}
	// TODO 返回数据
	response.Success(ctx, gin.H{"problemFileSubmit": problemFileSubmit}, "成功")

}

// @title    SubmitList
// @description   获取提交情况
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemFileController) SubmitList(ctx *gin.Context) {
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
	var problemFile model.ProblemFile

	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "ProblemFile", problem_id).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "ProblemFile", problem_id).Result()
		if json.Unmarshal([]byte(art), &problemFile) == nil {
			goto leap
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "ProblemFile", problem_id)
		}
	}

	// TODO 查看文件题是否在数据库中存在
	if p.DB.Where("id = ?", problem_id).First(&problemFile).Error != nil {
		response.Fail(ctx, nil, "文件题不存在")
		return
	}
	{
		// TODO 将文件题存入redis供下次使用
		v, _ := json.Marshal(problemFile)
		p.Redis.HSet(ctx, "ProblemFile", problem_id, v)
	}
leap:
	// TODO 查看题目所属的测试
	var exam model.Exam
	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "Exam", problemFile.ExamId.String()).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "Exam", problemFile.ExamId.String()).Result()
		if json.Unmarshal([]byte(art), &exam) == nil {
			goto levp
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "Exam", problemFile.ExamId.String())
		}
	}

	// TODO 查看文件题是否在数据库中存在
	if p.DB.Where("id = ?", problemFile.ExamId).First(&exam).Error != nil {
		response.Fail(ctx, nil, "测试不存在")
		return
	}
	{
		// TODO 将文件题存入redis供下次使用
		v, _ := json.Marshal(exam)
		p.Redis.HSet(ctx, "Exam", problemFile.ExamId.String(), v)
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
	if p.DB.Where("id = ?", exam.GroupId).First(&group).Error != nil {
		response.Fail(ctx, nil, "小组不存在")
		return
	}
	{
		// TODO 将文件题存入redis供下次使用
		v, _ := json.Marshal(group)
		p.Redis.HSet(ctx, "Group", exam.GroupId.String(), v)
	}
letp:

	// TODO 查看用户是否是组长或提交者
	if group.LeaderId != user.ID && user.ID.String() != user_id {
		response.Fail(ctx, nil, "权限不足")
		return
	}
	var problemFileSubmits []model.ProblemFileSubmit
	p.DB.Where("problem_file_id = ? and user_id = ?", problemFile.ID, user_id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&problemFile)

	var total int64
	p.DB.Where("problem_file_id = ? and user_id = ?", problemFile.ID, user_id).Model(model.ProblemFileSubmit{}).Count(&total)
	// TODO 返回数据
	response.Success(ctx, gin.H{"problemFileSubmits": problemFileSubmits, "total": total}, "成功")
}

// @title    NewProblemFileController
// @description   新建一个IProblemFileController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   IProblemFileController		返回一个IProblemFileController用于调用各种函数
func NewProblemFileController() IProblemFileController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.ProblemFile{})
	db.AutoMigrate(model.ProblemFileSubmit{})
	return ProblemFileController{DB: db, Redis: redis}
}
