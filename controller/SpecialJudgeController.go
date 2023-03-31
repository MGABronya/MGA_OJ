// @Title  SpecialJudgeController
// @Description  该文件提供关于操作特判的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	TQ "MGA_OJ/Test-request"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"
	"MGA_OJ/vo"
	"log"
	"time"

	"github.com/go-redis/redis/v9"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ISpecialJudgeController			定义了特判类接口
type ISpecialJudgeController interface {
	Create(ctx *gin.Context) // 创建
	Update(ctx *gin.Context) // 更新
}

// SpecialJudgeController			定义了特判工具类
type SpecialJudgeController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Create
// @description   创建一个特判
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SpecialJudgeController) Create(ctx *gin.Context) {
	var specialJudgeRequest vo.SpecialJudgeRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&specialJudgeRequest); err != nil {
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

	// TODO 测试运行情况
	test := TQ.JudgeRun(specialJudgeRequest.Language, specialJudgeRequest.Code, specialJudgeRequest.Input, 1024*256, uint(time.Second*30))

	// TODO 查看测试情况
	if test != "ok" {
		// TODO 测试不通过
		response.Fail(ctx, gin.H{"test": test}, "测试不通过")
	}

	specialJudge := model.SpecialJudge{
		Language: specialJudgeRequest.Language,
		Code:     specialJudgeRequest.Code,
		UserId:   user.ID,
	}
	// TODO 插入数据
	if err := s.DB.Create(&specialJudge).Error; err != nil {
		response.Fail(ctx, nil, "特判上传出错，数据验证有误")
		return
	}

	// TODO 成功
	response.Success(ctx, gin.H{"specialJudge": specialJudge}, "创建成功")
}

// @title    Update
// @description   更新一篇特判的内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SpecialJudgeController) Update(ctx *gin.Context) {
	var specialJudgeRequest vo.SpecialJudgeRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&specialJudgeRequest); err != nil {
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

	var specialJudge model.SpecialJudge

	if s.DB.Where("id = ?", id).First(&specialJudge) != nil {
		response.Fail(ctx, nil, "特判不存在")
		return
	}

	// TODO 查看是否是用户作者
	if user.ID != specialJudge.UserId {
		response.Fail(ctx, nil, "不是特判作者，无法修改特判")
		return
	}

	// TODO 测试运行情况
	test := TQ.JudgeRun(specialJudgeRequest.Language, specialJudgeRequest.Code, specialJudgeRequest.Input, 1024*256, uint(time.Second*30))

	// TODO 查看测试情况
	if test != "ok" {
		// TODO 测试不通过
		response.Fail(ctx, gin.H{"test": test}, "测试不通过")
	}

	specialJudgeUpdate := model.SpecialJudge{
		Language: specialJudgeRequest.Language,
		Code:     specialJudgeRequest.Code,
	}

	// TODO 更新特判内容
	s.DB.Model(&specialJudge).Updates(specialJudgeUpdate)

	// TODO 成功
	response.Success(ctx, nil, "更新成功")
}

// @title    NewSpecialJudgeController
// @description   新建一个ISpecialJudgeController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   ISpecialJudgeController		返回一个ISpecialJudgeController用于调用各种函数
func NewSpecialJudgeController() ISpecialJudgeController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.SpecialJudge{})
	return SpecialJudgeController{DB: db, Redis: redis}
}
