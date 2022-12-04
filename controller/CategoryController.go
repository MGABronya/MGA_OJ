// @Title  CategoryController
// @Description  该文件提供关于操作分类的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"MGA_OJ/Interface"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"
	"encoding/json"
	"log"
	"strconv"

	"github.com/go-redis/redis/v9"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ICategoryController			定义了分类类接口
type ICategoryController interface {
	Interface.RestInterface // 包含增删查改功能
}

// CategoryController			定义了分类工具类
type CategoryController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Create
// @description   创建一篇分类
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CategoryController) Create(ctx *gin.Context) {
	var category model.Category
	// TODO 数据验证
	if err := ctx.ShouldBind(&category); err != nil {
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

	// TODO 插入数据
	if err := c.DB.Create(&category).Error; err != nil {
		response.Fail(ctx, nil, "分类上传出错，数据库存储错误")
		return
	}

	// TODO 成功
	response.Success(ctx, gin.H{"category": category}, "创建成功")
}

// @title    Update
// @description   更新一篇分类的内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CategoryController) Update(ctx *gin.Context) {
	var category model.Category
	// TODO 数据验证
	if err := ctx.ShouldBind(&category); err != nil {
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

	// TODO 查找对应分类
	id := ctx.Params.ByName("id")

	if c.DB.Where("id = ?", id).First(&category) != nil {
		response.Fail(ctx, nil, "分类不存在")
		return
	}

	// TODO 更新分类内容
	c.DB.Where("id = ?", id).Updates(category)

	// TODO 移除损坏数据
	c.Redis.HDel(ctx, "Category", id)

	// TODO 成功
	response.Success(ctx, nil, "更新成功")
}

// @title    Show
// @description   查看一篇分类的内容
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CategoryController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var category model.Category

	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Category", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Category", id).Result()
		if json.Unmarshal([]byte(cate), &category) == nil {
			response.Success(ctx, gin.H{"category": category}, "成功")
			return
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Category", id)
		}
	}

	// TODO 查看分类是否在数据库中存在
	if c.DB.Where("id = ?", id).First(&category).Error != nil {
		response.Fail(ctx, nil, "分类不存在")
		return
	}

	response.Success(ctx, gin.H{"category": category}, "成功")

	// TODO 将分类存入redis供下次使用
	v, _ := json.Marshal(category)
	c.Redis.HSet(ctx, "Category", id, v)
}

// @title    Delete
// @description   删除一篇分类
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CategoryController) Delete(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	var category model.Category

	// TODO 查看分类是否存在
	if c.DB.Where("id = ?", id).First(&category).Error != nil {
		response.Fail(ctx, nil, "分类不存在")
		return
	}

	// TODO 判断当前用户是否为分类的作者
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有操作分类的权力
	if user.Level < 2 {
		response.Fail(ctx, nil, "用户权限不足")
		return
	}

	// TODO 删除分类
	c.DB.Delete(&category)

	// TODO 移除损坏数据
	c.Redis.HDel(ctx, "Category", id)

	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description   获取多篇分类
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CategoryController) PageList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var categorys []model.Category

	// TODO 查找所有分页中可见的条目
	c.DB.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&categorys)

	var total int64
	c.DB.Model(model.Category{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"categorys": categorys, "total": total}, "成功")
}

// @title    NewCategoryController
// @description   新建一个ICategoryController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   ICategoryController		返回一个ICategoryController用于调用各种函数
func NewCategoryController() ICategoryController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.Category{})
	return CategoryController{DB: db, Redis: redis}
}
