// @Title  ChatController
// @Description  该文件提供关于操作群聊的各种方法
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
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// IChatController			定义了群聊类接口
type IChatController interface {
	Interface.MassageInterface // 包含了信息交流相关方法
}

// ChatController			定义了群聊工具类
type ChatController struct {
	DB       *gorm.DB            // 含有一个数据库指针
	Redis    *redis.Client       // 含有一个redis指针
	UpGrader *websocket.Upgrader // 用于持久化连接
}

// @title    Send
// @description   发送一条群聊
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c ChatController) Send(ctx *gin.Context) {
	var requestChat vo.ChatRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestChat); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取指定用户组
	id := ctx.Params.ByName("id")

	var group model.Group

	// TODO 查看用户组是否存在
	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Group", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Group", id).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Group", id)
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if c.DB.Where("id = (?)", id).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		c.Redis.HSet(ctx, "Group", id, v)
	}
leep:

	// TODO 查看是否在用户组中
	if c.DB.Where("group_id = (?) and user_id = (?)", group.ID, user.ID).First(&model.UserList{}).Error != nil {
		response.Fail(ctx, nil, "不在用户组")
		return
	}

	// TODO 创建群聊信息
	chat := model.Chat{
		Content:  requestChat.Content,
		ResLong:  requestChat.ResLong,
		ResShort: requestChat.ResShort,
		GroupId:  group.ID,
		Author:   user.ID,
		ID:       uuid.NewV4(),
	}

	// TODO 将chat打包
	v, _ := json.Marshal(chat)

	// TODO 将chat存入频道
	c.Redis.Publish(ctx, "ChatChan"+group.ID.String(), v)

	// TODO 将chat存入redis数据库
	c.Redis.LPush(ctx, "ChatList"+group.ID.String(), v)

	// TODO 查找该组的成员
	var userLists []model.UserList
	// TODO 查看用户组是否存在
	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "UserList", group.ID.String()).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "UserList", group.ID.String()).Result()
		if json.Unmarshal([]byte(cate), &userLists) == nil {
			goto userlist
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "UserList", group.ID.String())
		}
	}

	// TODO 查看用户组是否在数据库中存在
	c.DB.Where("group_id = (?)", group.ID).Find(&userLists)
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(userLists)
		c.Redis.HSet(ctx, "UserList", group.ID.String(), v)
	}
userlist:

	for i := range userLists {
		// TODO 跳过自己
		if user.ID == userLists[i].UserId {
			continue
		}
		// TODO 将连接请求放入频道
		c.Redis.Publish(ctx, "ChatLinkChan"+userLists[i].UserId.String(), v)
		// TODO 将chat放入连接库
		c.Redis.HSet(ctx, "ChatLink"+userLists[i].UserId.String(), group.ID.String(), v)
	}

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    LinkList
// @description   获取多篇连接
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c ChatController) LinkList(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var chats []model.Chat

	// TODO 查找所有条目
	cats, _ := c.Redis.HGetAll(ctx, "ChatLink"+user.ID.String()).Result()

	for i := range cats {
		var chat model.Chat
		json.Unmarshal([]byte(cats[i]), &chat)
		chats = append(chats, chat)
	}

	// TODO 根据是否已读和时间排序
	sort.Sort(model.ChatSlice(chats))

	// TODO 返回数据
	response.Success(ctx, gin.H{"chats": chats}, "成功")
}

// @title    ChatList
// @description   列出聊天列表
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c ChatController) ChatList(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取指定用户组
	id := ctx.Params.ByName("id")

	var group model.Group

	// TODO 查看用户组是否存在
	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Group", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Group", id).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Group", id)
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if c.DB.Where("id = (?)", id).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		c.Redis.HSet(ctx, "Group", id, v)
	}
leep:

	// TODO 查看是否在用户组中
	if c.DB.Where("group_id = (?) and user_id = (?)", group.ID, user.ID).First(&model.UserList{}).Error != nil {
		response.Fail(ctx, nil, "不在用户组")
		return
	}

	var chats []model.Chat

	cats, _ := c.Redis.LRange(ctx, "ChatList"+group.ID.String(), 0, -1).Result()

	// TODO 整理聊天记录
	for i := range cats {
		var chat model.Chat
		json.Unmarshal([]byte(cats[i]), &chat)
		chats = append(chats, chat)
	}

	// TODO 返回数据
	response.Success(ctx, gin.H{"chats": chats}, "成功")
}

// @title    RemoveLink
// @description   移除指定连接
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c ChatController) RemoveLink(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取指定用户组
	id := ctx.Params.ByName("id")

	var group model.Group

	// TODO 查看用户组是否存在
	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Group", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Group", id).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Group", id)
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if c.DB.Where("id = (?)", id).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		c.Redis.HSet(ctx, "Group", id, v)
	}
leep:

	// TODO 删除指定条目
	c.Redis.HDel(ctx, "ChatLink"+user.ID.String(), group.ID.String()).Result()

	// TODO 返回数据
	response.Success(ctx, nil, "移除成功")
}

// @title    Receive
// @description   创建接收通信
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c ChatController) Receive(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取指定用户组
	id := ctx.Params.ByName("id")

	var group model.Group

	// TODO 查看用户组是否存在
	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Group", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Group", id).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Group", id)
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if c.DB.Where("id = (?)", id).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		c.Redis.HSet(ctx, "Group", id, v)
	}
leep:

	// TODO 查看是否在用户组中
	if c.DB.Where("group_id = (?) and user_id = (?)", group.ID, user.ID).First(&model.UserList{}).Error != nil {
		response.Fail(ctx, nil, "不在用户组")
		return
	}

	// TODO 订阅消息
	pubSub := c.Redis.Subscribe(ctx, "ChatChan"+group.ID.String())
	defer pubSub.Close()
	// TODO 获得消息管道
	ch := pubSub.Channel()

	// TODO 升级get请求为webSocket协议
	ws, err := c.UpGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	// TODO 监听消息
	for msg := range ch {
		var chat model.Chat
		json.Unmarshal([]byte(msg.Payload), &chat)
		// TODO 写入ws数据
		// TODO 断开连接
		if err := ws.WriteJSON(chat); err != nil {
			break
		}
	}

}

// @title    ReceiveLink
// @description   创建接收连接
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c ChatController) ReceiveLink(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 订阅消息
	pubSub := c.Redis.Subscribe(ctx, "ChatLinkChan"+user.ID.String())
	defer pubSub.Close()
	// TODO 获得消息管道
	ch := pubSub.Channel()

	// TODO 升级get请求为webSocket协议
	ws, err := c.UpGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	// TODO 监听消息
	for msg := range ch {
		var chat model.Chat
		json.Unmarshal([]byte(msg.Payload), &chat)
		// TODO 断开连接
		if err := ws.WriteJSON(chat); err != nil {
			break
		}
	}

}

// @title    NewChatController
// @description   新建一个INewChatController
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   IChatController		返回一个IChatController用于调用各种函数
func NewChatController() IChatController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	upGrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return ChatController{DB: db, Redis: redis, UpGrader: upGrader}
}
