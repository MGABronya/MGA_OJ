// @Title  PostController
// @Description  该文件提供关于操作题解的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"MGA_OJ/Interface"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"
	"MGA_OJ/vo"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IPostController			定义了题解类接口
type IPostController interface {
	Interface.RestInterface // 包含增删查改功能
}

// PostController			定义了题解工具类
type PostController struct {
	DB *gorm.DB // 含有一个数据库指针
}

// @title    Create
// @description   创建一篇题解
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) Create(ctx *gin.Context) {
	var requestPost vo.PostRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestPost); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 查找对应题目
	id := ctx.Params.ByName("id")

	var problem model.Problem

	// TODO 查看数据库中是否有该问题
	if p.DB.Where("id = ?", id).First(&problem).Error != nil {
		response.Fail(ctx, nil, "题目不存在")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 创建题解
	var post = model.Post{
		UserId:    user.ID,
		ProblemId: problem.ID,
		Title:     requestPost.Title,
		Content:   requestPost.Content,
		Reslong:   requestPost.Reslong,
		Resshort:  requestPost.Resshort,
	}

	// TODO 插入数据
	if err := p.DB.Create(&post).Error; err != nil {
		response.Fail(ctx, nil, "题解上传出错，数据库存储错误")
		return
	}

	// TODO 成功
	response.Success(ctx, gin.H{"post": post}, "创建成功")
}

// @title    Update
// @description   更新一篇题解的内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) Update(ctx *gin.Context) {
	var requestPost vo.PostRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestPost); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应题解
	id := ctx.Params.ByName("id")

	var post model.Post

	if p.DB.Where("id = ?", id).First(&post) != nil {
		response.Fail(ctx, nil, "题解不存在")
		return
	}

	// TODO 查看是否是用户作者
	if user.ID != post.UserId {
		response.Fail(ctx, nil, "不是题解作者，无法修改题解")
		return
	}

	// TODO 更新题解内容
	p.DB.Where("id = ?", id).Updates(requestPost)

	// TODO 成功
	response.Success(ctx, nil, "更新成功")
}

// @title    Show
// @description   查看一篇题解的内容
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var post model.Post

	// TODO 查看题解是否在数据库中存在
	if p.DB.Where("id = ?", id).First(&post).Error != nil {
		response.Fail(ctx, nil, "题解不存在")
		return
	}

	response.Success(ctx, gin.H{"post": post}, "成功")
}

// @title    Delete
// @description   删除一篇题解
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) Delete(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	var post model.Post

	// TODO 查看题解是否存在
	if p.DB.Where("id = ?", id).First(&post).Error != nil {
		response.Fail(ctx, nil, "题解不存在")
		return
	}

	// TODO 判断当前用户是否为题解的作者
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有操作题解的权力
	if user.ID != post.UserId && user.Level < 4 {
		response.Fail(ctx, nil, "题解不属于您，请勿非法操作")
		return
	}

	// TODO 删除题解
	p.DB.Delete(&post)

	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description   获取多篇题解
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) PageList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var posts []model.Post

	// TODO 查找所有分页中可见的条目
	p.DB.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&posts)

	var total int64
	p.DB.Model(model.Post{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"posts": posts, "total": total}, "成功")
}

// @title    NewPostController
// @description   新建一个IPostController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   IPostController		返回一个IPostController用于调用各种函数
func NewPostController() IPostController {
	db := common.GetDB()
	db.AutoMigrate(model.Post{})
	return PostController{DB: db}
}
