// @Title  RealNameController
// @Description  该文件提供关于操作实名的各种方法
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
	"path"
	"strconv"

	"github.com/go-redis/redis/v9"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IRealNameController			定义了实名类接口
type IRealNameController interface {
	Interface.RestInterface       // 包含了增删查改功能
	Upload(ctx *gin.Context)      // 上传实名表单
	StudentList(ctx *gin.Context) // 查看已上传的实名表单
}

// RealNameController			定义了实名工具类
type RealNameController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Create
// @description   创建一个实名
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RealNameController) Create(ctx *gin.Context) {
	var realNameRequest vo.RealNameRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&realNameRequest); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找数据库中是否存在此人
	var student model.Student
	if r.DB.Where("student_id = ?", realNameRequest.StudentId).First(&student).Error != nil {
		response.Fail(ctx, nil, "数据库中不存在该学号")
		return
	}

	// TODO 查看名字是否对的上
	if student.Name != realNameRequest.Name {
		response.Fail(ctx, nil, "姓名与学号不匹配")
		return
	}

	// TODO 查看学号是否已经实名了
	if r.DB.Where("student_id = ?", realNameRequest.StudentId).First(&model.RealName{}).Error == nil {
		response.Fail(ctx, nil, "该学号已被实名")
		return
	}

	realName := model.RealName{
		Name:      realNameRequest.Name,
		StudentId: realNameRequest.StudentId,
		Grade:     student.Grade,
		Major:     student.Major,
		UserId:    user.ID,
	}
	// TODO 插入数据
	if err := r.DB.Create(&realName).Error; err != nil {
		response.Fail(ctx, nil, "实名上传出错，数据验证有误")
		return
	}

	// TODO 成功
	response.Success(ctx, gin.H{"realName": realName}, "创建成功")
}

// @title    Update
// @description   更新一篇实名的内容
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RealNameController) Update(ctx *gin.Context) {
	var realNameRequest vo.RealNameRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&realNameRequest); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var realName model.RealName

	if r.DB.Where("user_id = (?)", user.ID).First(&realName).Error != nil {
		response.Fail(ctx, nil, "实名不存在")
		return
	}

	// TODO 查找数据库中是否存在此人
	var student model.Student
	if r.DB.Where("student_id = ?", realNameRequest.StudentId).First(&student).Error != nil {
		response.Fail(ctx, nil, "数据库中不存在该学号")
		return
	}

	// TODO 查看名字是否对的上
	if student.Name != realNameRequest.Name {
		response.Fail(ctx, nil, "姓名与学号不匹配")
		return
	}

	// TODO 查看学号是否已经实名了
	if r.DB.Where("student_id = ?", realNameRequest.StudentId).First(&model.RealName{}).Error == nil {
		response.Fail(ctx, nil, "该学号已被实名")
		return
	}

	realNameUpdate := model.RealName{
		Name:      realNameRequest.Name,
		StudentId: realNameRequest.StudentId,
		Grade:     student.Grade,
		Major:     student.Major,
		UserId:    user.ID,
	}

	// TODO 更新实名内容
	r.DB.Model(&realName).Updates(realNameUpdate)

	// TODO 成功
	response.Success(ctx, nil, "更新成功")
	r.Redis.HDel(ctx, "RealName", user.ID.String())
}

// @title    Show
// @description   查看一篇实名的内容
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RealNameController) Show(ctx *gin.Context) {

	// TODO 查找对应实名
	id := ctx.Params.ByName("id")

	var realName model.RealName

	// TODO 先看redis中是否存在
	if ok, _ := r.Redis.HExists(ctx, "RealName", id).Result(); ok {
		cate, _ := r.Redis.HGet(ctx, "RealName", id).Result()
		if json.Unmarshal([]byte(cate), &realName) == nil {
			// TODO 跳过数据库搜寻realName过程
			goto leep
		} else {
			// TODO 移除损坏数据
			r.Redis.HDel(ctx, "RealName", id)
		}
	}

	// TODO 查看实名是否在数据库中存在
	if r.DB.Where("user_id = (?)", id).First(&realName).Error != nil {
		response.Fail(ctx, nil, "实名不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(realName)
		r.Redis.HSet(ctx, "RealName", id, v)
	}

leep:

	// TODO 成功
	response.Success(ctx, gin.H{"realName": realName}, "查看成功")
}

// @title    Delete
// @description   删除一篇实名的内容
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RealNameController) Delete(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var realName model.RealName

	if r.DB.Where("user_id = (?)", user.ID).First(&realName).Error != nil {
		response.Fail(ctx, nil, "实名不存在")
		return
	}

	// TODO 删除实名内容
	r.DB.Delete(&realName)

	r.Redis.HDel(ctx, "RealName", user.ID.String())

	// TODO 成功
	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description   查看一页实名的内容
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RealNameController) PageList(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 取出用户权限
	if user.Level < 4 {
		response.Fail(ctx, nil, "用户权限不足")
		return
	}

	var realNames []model.RealName

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	var total int64

	// TODO 查找所有分页中可见的条目
	r.DB.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&realNames)

	r.DB.Model(model.RealName{}).Count(&total)

	// TODO 成功
	response.Success(ctx, gin.H{"realNames": realNames, "total": total}, "查看成功")
}

// @title    StudentList
// @description   查看已上传的实名列表
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RealNameController) StudentList(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 取出用户权限
	if user.Level < 4 {
		response.Fail(ctx, nil, "用户权限不足")
		return
	}

	var students []model.Student

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	var total int64

	// TODO 查找所有分页中可见的条目
	r.DB.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&students)

	r.DB.Model(model.Student{}).Count(&total)

	// TODO 成功
	response.Success(ctx, gin.H{"students": students, "total": total}, "查看成功")
}

// @title    Upload
// @description   上传表单
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RealNameController) Upload(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 取出用户权限
	if user.Level < 4 {
		response.Fail(ctx, nil, "用户权限不足")
		return
	}

	file, err := ctx.FormFile("file")

	//TODO 数据验证
	if err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 验证文件格式
	extName := path.Ext(file.Filename)
	allowExtMap := map[string]bool{
		".xls":  true,
		".xlsx": true,
		".csv":  true,
	}

	// TODO 格式验证
	if _, ok := allowExtMap[extName]; !ok {
		response.Fail(ctx, nil, "文件后缀有误")
		return
	}

	// TODO 将文件存入本地
	ctx.SaveUploadedFile(file, "./students/"+file.Filename)

	// TODO 解析文件
	res, err := util.Read("./students/" + file.Filename)

	// TODO 解析有误
	if err != nil || res == nil {
		response.Fail(ctx, nil, "文件解析有误")
		return
	}

	// TODO 查看文件格式
	if res[0][0] != "StudentId" || res[0][1] != "Name" || res[0][2] != "College" || res[0][3] != "Grade" || res[0][4] != "Major" {
		response.Fail(ctx, nil, "文件格式有误")
		return
	}

	// TODO 读入文件
	for i := 1; i < len(res); i++ {
		if len(res[i]) < 5 {
			break
		}
		student := model.Student{
			StudentId: res[i][0],
			Name:      res[i][1],
			College:   res[i][2],
			Grade:     res[i][3],
			Major:     res[i][4],
		}
		r.DB.Create(&student)
	}

	// TODO 成功
	response.Success(ctx, nil, "上传成功")
}

// @title    NewRealNameController
// @description   新建一个IRealNameController
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   IRealNameController		返回一个IRealNameController用于调用各种函数
func NewRealNameController() IRealNameController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.RealName{})
	db.AutoMigrate(model.Student{})
	return RealNameController{DB: db, Redis: redis}
}
