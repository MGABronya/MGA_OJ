// @Title  BadgeController
// @Description  该文件提供关于操作徽章的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
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

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

// IBadgeController			定义了徽章类接口
type IBadgeController interface {
	Interface.RestInterface    // 包含增删查改功能
	UserList(ctx *gin.Context) // 查看指定用户的徽章
}

// BadgeController			定义了徽章工具类
type BadgeController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Create
// @description   创建一篇徽章
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (b BadgeController) Create(ctx *gin.Context) {
	var requestBadge vo.BadgeRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestBadge); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 取出用户权限
	if user.Level < 4 {
		response.Fail(ctx, nil, "用户权限不足")
		return
	}

	// TODO 查看徽章参数是否合法
	if ok, err := util.CheckExpression([]byte(requestBadge.Condition)); !ok {
		response.Fail(ctx, nil, err)
		return
	}

	// TODO 创建徽章
	var badge = model.Badge{
		UserId:      user.ID,
		Name:        requestBadge.Name,
		Description: requestBadge.Description,
		ResLong:     requestBadge.ResLong,
		ResShort:    requestBadge.ResShort,
		Condition:   requestBadge.Condition,
		Iron:        requestBadge.Iron,
		Copper:      requestBadge.Copper,
		Silver:      requestBadge.Silver,
		Gold:        requestBadge.Gold,
	}

	// TODO 插入数据
	if err := b.DB.Create(&badge).Error; err != nil {
		response.Fail(ctx, nil, "徽章上传出错，数据库存储错误")
		return
	}

	// TODO 成功
	response.Success(ctx, gin.H{"badge": badge}, "创建成功")
}

// @title    Update
// @description   更新一篇徽章的内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (b BadgeController) Update(ctx *gin.Context) {
	var requestBadge model.Badge
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestBadge); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应徽章
	id := ctx.Params.ByName("id")

	var badge model.Badge

	if b.DB.Where("id = (?)", id).First(&badge).Error != nil {
		response.Fail(ctx, nil, "徽章不存在")
		return
	}

	// TODO 查看是否是用户作者
	if user.ID != badge.UserId {
		response.Fail(ctx, nil, "不是徽章作者，无法修改徽章")
		return
	}

	// TODO 查看徽章参数是否合法
	if ok, err := util.CheckExpression([]byte(requestBadge.Condition)); !ok {
		response.Fail(ctx, nil, err)
		return
	}

	// TODO 更新徽章内容
	b.DB.Where("id = (?)", id).Updates(requestBadge)

	// TODO 移除损坏数据
	b.Redis.HDel(ctx, "Badge", id)

	// TODO 成功
	response.Success(ctx, nil, "更新成功")
}

// @title    Show
// @description   查看一篇徽章的内容
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (b BadgeController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var badge model.Badge

	// TODO 先看redis中是否存在
	if ok, _ := b.Redis.HExists(ctx, "Badge", id).Result(); ok {
		cate, _ := b.Redis.HGet(ctx, "Badge", id).Result()
		if json.Unmarshal([]byte(cate), &badge) == nil {
			response.Success(ctx, gin.H{"badge": badge}, "成功")
			return
		} else {
			// TODO 移除损坏数据
			b.Redis.HDel(ctx, "Badge", id)
		}
	}

	// TODO 查看徽章是否在数据库中存在
	if b.DB.Where("id = (?)", id).First(&badge).Error != nil {
		response.Fail(ctx, nil, "徽章不存在")
		return
	}

	response.Success(ctx, gin.H{"badge": badge}, "成功")

	// TODO 将徽章存入redis供下次使用
	v, _ := json.Marshal(badge)
	b.Redis.HSet(ctx, "Badge", id, v)
}

// @title    Delete
// @description   删除一篇徽章
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (b BadgeController) Delete(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	var badge model.Badge

	// TODO 查看徽章是否存在
	if b.DB.Where("id = (?)", id).First(&badge).Error != nil {
		response.Fail(ctx, nil, "徽章不存在")
		return
	}

	// TODO 判断当前用户是否为徽章的作者
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有操作徽章的权力
	if user.ID != badge.UserId && user.Level < 4 {
		response.Fail(ctx, nil, "徽章不属于您，请勿非法操作")
		return
	}

	// TODO 删除徽章
	b.DB.Delete(&badge)

	// TODO 移除损坏数据
	b.Redis.HDel(ctx, "Badge", id)

	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description   获取多篇徽章
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (b BadgeController) PageList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 分页
	var badges []model.Badge

	// TODO 查找所有分页中可见的条目
	b.DB.Where("problem_id = (?)", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&badges)

	var total int64
	b.DB.Where("problem_id = (?)", id).Model(model.Badge{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"badges": badges, "total": total}, "成功")
}

// @title    UserList
// @description   获取指定用户的多篇徽章
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (b BadgeController) UserList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 分页
	var badges []model.Badge

	// TODO 查找所有分页中可见的条目
	b.DB.Where("user_id = (?)", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&badges)

	var total int64
	b.DB.Where("user_id = (?)", id).Model(model.Badge{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"badges": badges, "total": total}, "成功")
}

// @title    NewBadgeController
// @description   新建一个IBadgeController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   IBadgeController		返回一个IBadgeController用于调用各种函数
func NewBadgeController() IBadgeController {
	db := common.GetDB()
	db.AutoMigrate(model.Badge{})
	return BadgeController{DB: db}
}