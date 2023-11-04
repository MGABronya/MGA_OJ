// @Title  MessageController
// @Description  该文件提供关于操作留言板的各种方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package controller

import (
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

// IMessageController			定义了留言类接口
type IMessageController interface {
	Create(ctx *gin.Context)   // 创建
	Delete(ctx *gin.Context)   // 删除
	PageList(ctx *gin.Context) // 列出列表
	AICreate(ctx *gin.Context) // 创建ai模板
	AIDelete(ctx *gin.Context) // 删除ai模板
	AIShow(ctx *gin.Context)   // 查看ai模板
	AIUpdate(ctx *gin.Context) // 更新ai模板
}

// MessageController			定义了留言工具类
type MessageController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
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

	// TODO 查看是否存在自动回复留言ai
	var ai model.AI

	// TODO 先尝试在redis中寻找
	if ok, _ := m.Redis.HExists(ctx, "AI", userb.ID.String()).Result(); ok {
		art, _ := m.Redis.HGet(ctx, "AI", userb.ID.String()).Result()
		if json.Unmarshal([]byte(art), &ai) == nil {
			goto leep
		} else {
			// TODO 解码失败，删除字段
			m.Redis.HDel(ctx, "AI", userb.ID.String())
		}
	}

	// TODO 查看是否在数据库中存在
	if m.DB.Where("user_id = (?)", userb.ID.String()).First(&ai).Error != nil {
		response.Success(ctx, gin.H{"message": message}, "创建成功")
		return
	}
	{
		// TODO 将ai模板存入redis供下次使用
		v, _ := json.Marshal(ai)
		m.Redis.HSet(ctx, "AI", userb.ID.String(), v)
	}

leep:
	// TODO 设置为不回复自己时
	if userb.ID == user.ID && !ai.Reply {
		// TODO 成功
		response.Success(ctx, gin.H{"message": message}, "创建成功")
		return
	}
	// TODO 生成ai回复
	res, err := util.ChatGPT([]string{ai.Characters, message.Content}, "gpt-4")
	if err != nil {
		response.Fail(ctx, nil, "ai回复出错")
		return
	}
	// TODO 创建留言
	message = model.Message{
		Content: res.Choices[0].Message.Content,
		UserId:  userb.ID,
		Author:  userb.ID,
	}

	// TODO 插入数据
	if err := m.DB.Create(&message).Error; err != nil {
		response.Fail(ctx, nil, "ai留言出错，数据验证有误")
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

// @title    AICreate
// @description   创建AI模板
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (m MessageController) AICreate(ctx *gin.Context) {
	var requestAI vo.AIRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestAI); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 创建ai模板
	ai := model.AI{
		Characters: requestAI.Characters,
		Reply:      requestAI.Reply,
		Prologue:   requestAI.Prologue,
		UserId:     user.ID,
	}

	// TODO 插入数据
	if err := m.DB.Create(&ai).Error; err != nil {
		response.Fail(ctx, nil, "ai模板出错，数据验证有误")
		return
	}

	// TODO 创建留言
	message := model.Message{
		Content: ai.Prologue,
		UserId:  user.ID,
		Author:  user.ID,
	}

	// TODO 插入数据
	if err := m.DB.Create(&message).Error; err != nil {
		response.Fail(ctx, nil, "ai留言出错，数据验证有误")
		return
	}

	// TODO 成功
	response.Success(ctx, gin.H{"ai": ai}, "创建成功")
}

// @title    AIDelete
// @description   删除AI模板
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (m MessageController) AIDelete(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 删除AI模板
	m.DB.Where("user_id = (?)", user.ID).Delete(&model.AI{})

	// TODO 解码失败，删除字段
	m.Redis.HDel(ctx, "AI", user.ID.String())

	// TODO 成功
	response.Success(ctx, nil, "删除成功")
}

// @title    AIShow
// @description   查看AI模板
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (m MessageController) AIShow(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var ai model.AI

	// TODO 先尝试在redis中寻找
	if ok, _ := m.Redis.HExists(ctx, "AI", user.ID.String()).Result(); ok {
		art, _ := m.Redis.HGet(ctx, "AI", user.ID.String()).Result()
		if json.Unmarshal([]byte(art), &ai) == nil {
			response.Success(ctx, gin.H{"ai": ai}, "成功")
			return
		} else {
			// TODO 解码失败，删除字段
			m.Redis.HDel(ctx, "AI", user.ID.String())
		}
	}

	// TODO 查看是否在数据库中存在
	if m.DB.Where("user_id = (?)", user.ID.String()).First(&ai).Error != nil {
		response.Fail(ctx, nil, "ai模板不存在")
		return
	}

	response.Success(ctx, gin.H{"ai": ai}, "成功")

	// TODO 将ai模板存入redis供下次使用
	v, _ := json.Marshal(ai)
	m.Redis.HSet(ctx, "AI", user.ID.String(), v)
}

// @title    AIUpdate
// @description   修改AI模板
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (m MessageController) AIUpdate(ctx *gin.Context) {
	var requestAI vo.AIRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestAI); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 创建ai模板
	ai := model.AI{
		Characters: requestAI.Characters,
		Reply:      requestAI.Reply,
		Prologue:   requestAI.Prologue,
		UserId:     user.ID,
	}

	// TODO 修改数据
	m.DB.Model(model.AI{}).Where("user_id = (?)", user.ID.String()).Updates(ai)

	// TODO 成功
	response.Success(ctx, gin.H{"ai": ai}, "更新成功")

	// TODO 解码失败，删除字段
	m.Redis.HDel(ctx, "AI", user.ID.String())
}

// @title    NewMessageController
// @description   新建一个IMessageController
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   IMessageController		返回一个IMessageController用于调用各种函数
func NewMessageController() IMessageController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.Message{})
	db.AutoMigrate(model.AI{})
	return MessageController{DB: db, Redis: redis}
}
