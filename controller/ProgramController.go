// @Title  ProgramController
// @Description  该文件提供关于操作程序的各种方法
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
	"strconv"

	"github.com/go-redis/redis/v9"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IProgramController			定义了程序类接口
type IProgramController interface {
	Interface.RestInterface // 包含了增删查改功能
}

// ProgramController			定义了程序工具类
type ProgramController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Create
// @description   创建一个程序
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProgramController) Create(ctx *gin.Context) {
	var programRequest vo.ProgramRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&programRequest); err != nil {
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

	program := model.Program{
		Language: programRequest.Language,
		Code:     programRequest.Code,
		UserId:   user.ID,
	}
	// TODO 插入数据
	if err := p.DB.Create(&program).Error; err != nil {
		response.Fail(ctx, nil, "程序上传出错，数据验证有误")
		return
	}

	// TODO 成功
	response.Success(ctx, gin.H{"program": program}, "创建成功")
}

// @title    Update
// @description   更新一篇程序的内容
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProgramController) Update(ctx *gin.Context) {
	var programRequest vo.ProgramRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&programRequest); err != nil {
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

	// TODO 查找对应特判
	id := ctx.Params.ByName("id")

	var program model.Program

	if p.DB.Where("id = (?)", id).First(&program).Error != nil {
		response.Fail(ctx, nil, "程序不存在")
		return
	}

	// TODO 查看是否是用户作者
	if user.ID != program.UserId {
		response.Fail(ctx, nil, "不是程序作者，无法修改程序")
		return
	}

	programUpdate := model.Program{
		Language: programRequest.Language,
		Code:     programRequest.Code,
	}

	// TODO 更新特判内容
	p.DB.Model(&program).Updates(programUpdate)

	// TODO 成功
	response.Success(ctx, nil, "更新成功")
	p.Redis.HDel(ctx, "Program", id)
}

// @title    Show
// @description   查看一篇程序的内容
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProgramController) Show(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 取出用户权限
	if user.Level < 2 {
		response.Fail(ctx, nil, "用户权限不足")
		return
	}

	// TODO 查找对应程序
	id := ctx.Params.ByName("id")

	var program model.Program

	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "Program", id).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Program", id).Result()
		if json.Unmarshal([]byte(cate), &program) == nil {
			// TODO 跳过数据库搜寻program过程
			goto leep
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Program", id)
		}
	}

	// TODO 查看程序是否在数据库中存在
	if p.DB.Where("id = (?)", id).First(&program).Error != nil {
		response.Fail(ctx, nil, "程序不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(program)
		p.Redis.HSet(ctx, "Program", id, v)
	}

leep:

	// TODO 查看是否是用户作者
	if user.ID != program.UserId {
		response.Fail(ctx, nil, "不是程序作者，无法查看程序")
		return
	}

	// TODO 成功
	response.Success(ctx, gin.H{"program": program}, "查看成功")
}

// @title    Delete
// @description   删除一篇程序的内容
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProgramController) Delete(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 取出用户权限
	if user.Level < 2 {
		response.Fail(ctx, nil, "用户权限不足")
		return
	}

	// TODO 查找对应特判
	id := ctx.Params.ByName("id")

	var program model.Program

	if p.DB.Where("id = (?)", id).First(&program).Error != nil {
		response.Fail(ctx, nil, "程序不存在")
		return
	}

	// TODO 查看是否是用户作者
	if user.ID != program.UserId {
		response.Fail(ctx, nil, "不是程序作者，无法删除程序")
		return
	}

	// TODO 删除程序内容
	p.DB.Delete(&program)

	p.Redis.HDel(ctx, "Program", id)

	// TODO 成功
	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description   查看一页程序的内容
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProgramController) PageList(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 取出用户权限
	if user.Level < 2 {
		response.Fail(ctx, nil, "用户权限不足")
		return
	}

	var programs []model.Program

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	var total int64

	// TODO 查找所有分页中可见的条目
	p.DB.Where("user_id = (?)", user.ID).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&programs)

	p.DB.Where("user_id = (?)", user.ID).Model(model.Program{}).Count(&total)

	// TODO 成功
	response.Success(ctx, gin.H{"programs": programs}, "查看成功")
}

// @title    NewProgramController
// @description   新建一个IProgramController
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   IProgramController		返回一个IProgramController用于调用各种函数
func NewProgramController() IProgramController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.Program{})
	return ProgramController{DB: db, Redis: redis}
}
