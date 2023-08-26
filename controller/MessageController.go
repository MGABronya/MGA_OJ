// @Title  MessageController
// @Description  该文件提供关于操作留言板的各种方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package controller

import (
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"
	"MGA_OJ/vo"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IMessageController			定义了留言类接口
type IMessageController interface {
	Create(ctx *gin.Context)   // 创建
	Delete(ctx *gin.Context)   // 删除
	PageList(ctx *gin.Context) // 列出列表
}

// MessageController			定义了留言工具类
type MessageController struct {
	DB *gorm.DB // 含有一个数据库指针
}

// @title    Create
// @description   创建一条留言
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (m MessageController) Create(ctx *gin.Context) {
	var requestMessage vo.MessageRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestMessage); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取指定用户
	id := ctx.Params.ByName("id")

	var userb model.User

	// TODO 查看用户是否存在
	if m.DB.Where("id = (?)", id).First(&userb).Error != nil {
		response.Fail(ctx, nil, "指定用户不存在")
		return
	}

	// TODO 创建留言
	message := model.Message{
		Content: requestMessage.Content,
		UserId:  userb.ID,
		Author:  user.ID,
	}

	// TODO 插入数据
	if err := m.DB.Create(&message).Error; err != nil {
		response.Fail(ctx, nil, "留言出错，数据验证有误")
		return
	}

	// TODO 成功
	response.Success(ctx, gin.H{"message": message}, "创建成功")
}

// @title    Delete
// @description   删除一条留言
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (m MessageController) Delete(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	var message model.Message

	// TODO 查看留言是否存在
	if m.DB.Where("id = (?)", id).First(&message).Error != nil {
		response.Fail(ctx, nil, "留言不存在")
		return
	}

	// TODO 判断当前用户是否为留言板的拥有者
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有操作文章的权力
	if user.ID != message.UserId && user.Level < 4 {
		response.Fail(ctx, nil, "留言不属于您，请勿非法操作")
		return
	}

	// TODO 删除文章
	m.DB.Delete(&message)

	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description   获取多篇留言
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (m MessageController) PageList(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var messages []model.Message

	// TODO 查找所有分页中可见的条目
	m.DB.Where("user_id = (?)", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&messages)

	var total int64
	m.DB.Where("user_id = (?)", id).Model(model.Article{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"messages": messages, "total": total}, "成功")
}

// @title    NewMessageController
// @description   新建一个IMessageController
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   IMessageController		返回一个IMessageController用于调用各种函数
func NewMessageController() IMessageController {
	db := common.GetDB()
	db.AutoMigrate(model.Message{})
	return MessageController{DB: db}
}
