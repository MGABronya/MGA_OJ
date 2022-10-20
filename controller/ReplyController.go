// @Title  ReplyController
// @Description  该文件提供关于操作讨论的回复的各种方法
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

// IReplyController			定义了讨论的回复类接口
type IReplyController interface {
	Interface.RestInterface // 包含增删查改功能
}

// ReplyController			定义了讨论的回复工具类
type ReplyController struct {
	DB *gorm.DB // 含有一个数据库指针
}

// @title    Create
// @description   创建一篇讨论的回复
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r ReplyController) Create(ctx *gin.Context) {
	var requestReply vo.ReplyRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestReply); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 查找对应讨论
	id := ctx.Params.ByName("id")

	var comment model.Comment

	if r.DB.Where("id = ?", id).First(&comment).Error != nil {
		response.Fail(ctx, nil, "讨论不存在")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 创建讨论的回复
	var reply = model.Reply{
		UserId:    user.ID,
		CommentId: comment.ID,
		Content:   requestReply.Content,
		Reslong:   requestReply.Reslong,
		Resshort:  requestReply.Resshort,
	}

	// TODO 插入数据
	if err := r.DB.Create(&reply).Error; err != nil {
		response.Fail(ctx, nil, "讨论的回复上传出错，数据库存储错误")
		return
	}

	// TODO 成功
	response.Success(ctx, gin.H{"reply": reply}, "创建成功")
}

// @title    Update
// @description   更新一篇讨论的回复的内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r ReplyController) Update(ctx *gin.Context) {
	var requestReply vo.ReplyRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestReply); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应讨论的回复
	id := ctx.Params.ByName("id")

	var reply model.Reply

	if r.DB.Where("id = ?", id).First(&reply) != nil {
		response.Fail(ctx, nil, "讨论的回复不存在")
		return
	}

	// TODO 查看是否是用户作者
	if user.ID != reply.UserId {
		response.Fail(ctx, nil, "不是讨论的回复作者，无法修改讨论的回复")
		return
	}

	// TODO 更新讨论的回复内容
	r.DB.Where("id = ?", id).Updates(requestReply)

	// TODO 成功
	response.Success(ctx, nil, "更新成功")
}

// @title    Show
// @description   查看一篇讨论的回复的内容
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r ReplyController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var reply model.Reply

	// TODO 查看讨论的回复是否在数据库中存在
	if r.DB.Where("id = ?", id).First(&reply).Error != nil {
		response.Fail(ctx, nil, "讨论的回复不存在")
		return
	}

	response.Success(ctx, gin.H{"reply": reply}, "成功")
}

// @title    Delete
// @description   删除一篇讨论的回复
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r ReplyController) Delete(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	var reply model.Reply

	// TODO 查看讨论的回复是否存在
	if r.DB.Where("id = ?", id).First(&reply).Error != nil {
		response.Fail(ctx, nil, "讨论的回复不存在")
		return
	}

	// TODO 判断当前用户是否为讨论的回复的作者
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有操作讨论的回复的权力
	if user.ID != reply.UserId && user.Level < 4 {
		response.Fail(ctx, nil, "讨论的回复不属于您，请勿非法操作")
		return
	}

	// TODO 删除讨论的回复
	r.DB.Delete(&reply)

	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description   获取多篇讨论的回复
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r ReplyController) PageList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var replys []model.Reply

	// TODO 查找所有分页中可见的条目
	r.DB.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&replys)

	var total int64
	r.DB.Model(model.Reply{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"replys": replys, "total": total}, "成功")
}

// @title    NewReplyController
// @description   新建一个IReplyController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   IReplyController		返回一个IReplyController用于调用各种函数
func NewReplyController() IReplyController {
	db := common.GetDB()
	db.AutoMigrate(model.Reply{})
	return ReplyController{DB: db}
}
