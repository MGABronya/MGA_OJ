// @Title  CommentController
// @Description  该文件提供关于操作讨论的各种方法
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

// ICommentController			定义了讨论类接口
type ICommentController interface {
	Interface.RestInterface // 包含增删查改功能
}

// CommentController			定义了讨论工具类
type CommentController struct {
	DB *gorm.DB // 含有一个数据库指针
}

// @title    Create
// @description   创建一篇讨论
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CommentController) Create(ctx *gin.Context) {
	var requestComment vo.CommentRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestComment); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 查找对应题目
	id := ctx.Params.ByName("id")

	var problem model.Problem

	// TODO 查看数据库中是否有该问题
	if c.DB.Where("id = ?", id).First(&problem).Error != nil {
		response.Fail(ctx, nil, "题目不存在")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 创建评论
	var comment = model.Comment{
		UserId:    user.ID,
		ProblemId: problem.ID,
		Content:   requestComment.Content,
		Reslong:   requestComment.Reslong,
		Resshort:  requestComment.Resshort,
	}

	// TODO 插入数据
	if err := c.DB.Create(&comment).Error; err != nil {
		response.Fail(ctx, nil, "讨论上传出错，数据库存储错误")
		return
	}

	// TODO 成功
	response.Success(ctx, gin.H{"comment": comment}, "创建成功")
}

// @title    Update
// @description   更新一篇讨论的内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CommentController) Update(ctx *gin.Context) {
	var requestComment vo.CommentRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestComment); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应讨论
	id := ctx.Params.ByName("id")

	var comment model.Comment

	if c.DB.Where("id = ?", id).First(&comment) != nil {
		response.Fail(ctx, nil, "讨论不存在")
		return
	}

	// TODO 查看是否是用户作者
	if user.ID != comment.UserId {
		response.Fail(ctx, nil, "不是讨论作者，无法修改讨论")
		return
	}

	// TODO 更新讨论内容
	c.DB.Where("id = ?", id).Updates(requestComment)

	// TODO 成功
	response.Success(ctx, nil, "更新成功")
}

// @title    Show
// @description   查看一篇讨论的内容
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CommentController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var comment model.Comment

	// TODO 查看讨论是否在数据库中存在
	if c.DB.Where("id = ?", id).First(&comment).Error != nil {
		response.Fail(ctx, nil, "讨论不存在")
		return
	}

	response.Success(ctx, gin.H{"comment": comment}, "成功")
}

// @title    Delete
// @description   删除一篇讨论
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p CommentController) Delete(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	var comment model.Comment

	// TODO 查看讨论是否存在
	if p.DB.Where("id = ?", id).First(&comment).Error != nil {
		response.Fail(ctx, nil, "讨论不存在")
		return
	}

	// TODO 判断当前用户是否为讨论的作者
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有操作讨论的权力
	if user.ID != comment.UserId && user.Level < 4 {
		response.Fail(ctx, nil, "讨论不属于您，请勿非法操作")
		return
	}

	// TODO 删除讨论
	p.DB.Delete(&comment)

	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description   获取多篇讨论
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CommentController) PageList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var comments []model.Comment

	// TODO 查找所有分页中可见的条目
	c.DB.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&comments)

	var total int64
	c.DB.Model(model.Comment{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"comments": comments, "total": total}, "成功")
}

// @title    NewCommentController
// @description   新建一个ICommentController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   ICommentController		返回一个ICommentController用于调用各种函数
func NewCommentController() ICommentController {
	db := common.GetDB()
	db.AutoMigrate(model.Comment{})
	return CommentController{DB: db}
}
