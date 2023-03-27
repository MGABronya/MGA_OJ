// @Title  LetterController
// @Description  该文件提供关于操作私信的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"
	"MGA_OJ/util"
	"MGA_OJ/vo"
	"encoding/json"
	"log"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// ILetterController			定义了私信类接口
type ILetterController interface {
	Send(ctx *gin.Context)        // 发送私信
	LinkList(ctx *gin.Context)    // 列出连接列表
	ChatList(ctx *gin.Context)    // 列出聊天列表
	Receive(ctx *gin.Context)     // 建立实时接收
	ReceiveLink(ctx *gin.Context) // 建立连接实时接收
	RemoveLink(ctx *gin.Context)  // 移除某个连接
	Read(ctx *gin.Context)        // 已读
}

// LetterController			定义了私信工具类
type LetterController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Send
// @description   发送一条私信
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (l LetterController) Send(ctx *gin.Context) {
	var requestLetter vo.LetterRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestLetter); err != nil {
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
	if l.DB.Where("id = ?", id).First(&userb).Error != nil {
		response.Fail(ctx, nil, "指定用户不存在")
		return
	}

	// TODO 创建留言
	letter := model.Letter{
		Content:  requestLetter.Content,
		Reslong:  requestLetter.Reslong,
		Resshort: requestLetter.Resshort,
		UserId:   userb.ID,
		Author:   user.ID,
		ID:       uuid.NewV4(),
		Read:     false,
	}

	// TODO 将letter打包
	v, _ := json.Marshal(letter)

	// TODO 将letter存入redis哈希中
	l.Redis.HSet(ctx, "Letters", letter.ID.String(), v)

	// TODO 将letter存入频道
	l.Redis.Publish(ctx, "LetterChan"+util.StringMerge(user.ID.String(), userb.ID.String()), letter.ID.String())

	// TODO 将letter存入redis数据库
	l.Redis.LPush(ctx, "LetterList"+util.StringMerge(user.ID.String(), userb.ID.String()), letter.ID.String())

	// TODO 将letter放入连接库
	l.Redis.HSet(ctx, "LetterLink"+userb.ID.String(), user.ID.String(), letter.ID.String())

	// TODO 将连接请求放入频道
	l.Redis.Publish(ctx, "LetterLinkChan"+userb.ID.String(), user.ID.String())

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    LinkList
// @description   获取多篇连接
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (l LetterController) LinkList(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var letters []model.Letter

	// TODO 查找所有条目
	lets, _ := l.Redis.HGetAll(ctx, "LetterLink"+user.ID.String()).Result()

	for i := range lets {
		var letter model.Letter
		v, _ := l.Redis.HGet(ctx, "Letters", lets[i]).Result()
		json.Unmarshal([]byte(v), &letter)
		letters = append(letters, letter)
	}

	// TODO 根据是否已读和时间排序
	sort.Sort(model.LetterSlice(letters))

	// TODO 返回数据
	response.Success(ctx, gin.H{"letters": letters}, "成功")
}

// @title    ChatList
// @description   列出聊天列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (l LetterController) ChatList(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取指定用户
	id := ctx.Params.ByName("id")

	var userb model.User

	// TODO 查看用户是否存在
	if l.DB.Where("id = ?", id).First(&userb).Error != nil {
		response.Fail(ctx, nil, "指定用户不存在")
		return
	}
	var letters []model.Letter

	lets, _ := l.Redis.LRange(ctx, "LetterList"+util.StringMerge(user.ID.String(), userb.ID.String()), 0, -1).Result()

	// TODO 整理聊天记录
	for i := range lets {
		var letter model.Letter
		v, _ := l.Redis.HGet(ctx, "Letters", lets[i]).Result()
		json.Unmarshal([]byte(v), &letter)
		letters = append(letters, letter)
	}

	// TODO 返回数据
	response.Success(ctx, gin.H{"letter": letters}, "成功")
}

// @title    RemoveLink
// @description   移除指定连接
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (l LetterController) RemoveLink(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取指定用户
	id := ctx.Params.ByName("id")

	var userb model.User

	// TODO 查看用户是否存在
	if l.DB.Where("id = ?", id).First(&userb).Error != nil {
		response.Fail(ctx, nil, "指定用户不存在")
		return
	}

	// TODO 删除指定条目
	l.Redis.HDel(ctx, "LetterLink"+user.ID.String(), userb.ID.String()).Result()

	// TODO 返回数据
	response.Success(ctx, nil, "移除成功")
}

// @title    Receive
// @description   创建接收通信
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (l LetterController) Receive(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取指定用户
	id := ctx.Params.ByName("id")

	var userb model.User

	// TODO 查看用户是否存在
	if l.DB.Where("id = ?", id).First(&userb).Error != nil {
		response.Fail(ctx, nil, "指定用户不存在")
		return
	}

	// TODO 订阅消息
	pubSub := l.Redis.Subscribe(ctx, "LetterChan"+util.StringMerge(user.ID.String(), userb.ID.String()))

	// TODO 获得消息管道
	ch := pubSub.Channel()

	// TODO 监听消息
	for msg := range ch {
		var letter model.Letter
		v, _ := l.Redis.HGet(ctx, "Letters", msg.Payload).Result()
		json.Unmarshal([]byte(v), &letter)
	}

}

// @title    ReceiveLink
// @description   创建接收连接
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (l LetterController) ReceiveLink(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 订阅消息
	pubSub := l.Redis.Subscribe(ctx, "LetterLinkChan"+user.ID.String())

	// TODO 获得消息管道
	ch := pubSub.Channel()

	// TODO 监听消息
	for msg := range ch {
		msg.Payload
	}

}

// @title    Read
// @description   已读
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (l LetterController) Read(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取指定私信
	id := ctx.Params.ByName("id")

	// TODO 将letter从redis哈希中取出
	v, _ := l.Redis.HGet(ctx, "Letters", id).Result()
	var letter model.Letter
	json.Unmarshal([]byte(v), &letter)

	// TODO 查看是否是letter的发送对象
	if letter.UserId != user.ID {
		response.Fail(ctx, nil, "不是发送对象，请勿非法操作")
	}

	// TODO 标记为已读
	letter.Read = true

	// TODO 将letter打包
	t, _ := json.Marshal(letter)

	// TODO 将letter存入redis哈希中
	l.Redis.HSet(ctx, "Letters", letter.ID.String(), t)

	// TODO 成功
	response.Success(ctx, nil, "已读成功")
}

// @title    NewLetterController
// @description   新建一个INewLetterController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   ILetterController		返回一个ILetterController用于调用各种函数
func NewLetterController() ILetterController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	return LetterController{DB: db, Redis: redis}
}
