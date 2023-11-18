// @Title  BadgeController
// @Description  该文件提供关于操作徽章的各种方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package controller

import (
	Handle "MGA_OJ/Behavior"
	"MGA_OJ/Interface"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"
	"MGA_OJ/vo"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// IBadgeController			定义了徽章类接口
type IBadgeController interface {
	Interface.RestInterface              // 包含增删查改功能
	UserList(ctx *gin.Context)           // 查看指定用户的徽章
	UserShow(ctx *gin.Context)           // 查看指定用户的指定徽章
	EvaluateExpression(ctx *gin.Context) // 计算表达式
	BehaviorList(ctx *gin.Context)       // 查看所有变量
	BehaviorShow(ctx *gin.Context)       // 查看指定变量描述
	Publish(ctx *gin.Context)            // 徽章长连接
}

// BadgeController			定义了徽章工具类
type BadgeController struct {
	DB       *gorm.DB            // 含有一个数据库指针
	Redis    *redis.Client       // 含有一个redis指针
	UpGrader *websocket.Upgrader // 用于持久化连接
}

// @title    Create
// @description   创建一篇徽章
// @auth      MGAronya       2022-9-16 12:15
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
	if ok, err := Handle.CheckExpression([]byte(requestBadge.Condition)); !ok {
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
		File:        requestBadge.File,
	}

	// TODO 插入数据
	if err := b.DB.Create(&badge).Error; err != nil {
		response.Fail(ctx, nil, "徽章上传出错，数据库存储错误")
		return
	}

	// TODO 为勋章插入行为
	for k := range Handle.Behaviors {
		// TODO 如果字符串中包含了当前字符串数组元素
		if strings.Contains(badge.Condition, k) {
			// TODO 将该元素添加到订阅列表中
			badgeBehavior := model.BadgeBehavior{
				Name:    k,
				BadgeId: badge.ID,
			}
			b.DB.Create(&badgeBehavior)
		}
	}

	// TODO 成功
	response.Success(ctx, gin.H{"badge": badge}, "创建成功")
}

// @title    Update
// @description   更新一篇徽章的内容
// @auth      MGAronya       2022-9-16 12:15
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
	if ok, err := Handle.CheckExpression([]byte(requestBadge.Condition)); !ok {
		response.Fail(ctx, nil, err)
		return
	}

	// TODO 更新徽章内容
	b.DB.Where("id = (?)", id).Updates(requestBadge)

	// TODO 删除原先的徽章行为映射
	b.DB.Where("badge_id = ?", badge.ID).Delete(&model.BadgeBehavior{})

	// TODO 为勋章插入行为
	for k := range Handle.Behaviors {
		// TODO 如果字符串中包含了当前字符串数组元素
		if strings.Contains(badge.Condition, k) {
			// TODO 将该元素添加到订阅列表中
			badgeBehavior := model.BadgeBehavior{
				Name:    k,
				BadgeId: badge.ID,
			}
			b.DB.Create(&badgeBehavior)
		}
	}

	// TODO 移除损坏数据
	b.Redis.HDel(ctx, "Badge", id)

	b.DB.Where("id = (?)", id).First(&badge)

	// TODO 成功
	response.Success(ctx, gin.H{"badge": badge}, "更新成功")
}

// @title    Show
// @description   查看一篇徽章的内容
// @auth      MGAronya       2022-9-16 12:19
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
// @auth      MGAronya       2022-9-16 12:20
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
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (b BadgeController) PageList(ctx *gin.Context) {

	// TODO 分页
	var badges []model.Badge

	// TODO 查找所有分页中可见的条目
	b.DB.Order("created_at desc").Find(&badges)

	// TODO 返回数据
	response.Success(ctx, gin.H{"badges": badges}, "成功")
}

// @title    UserList
// @description   获取指定用户的多篇徽章
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (b BadgeController) UserList(ctx *gin.Context) {

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 分页
	var userBadges []model.UserBadge

	// TODO 查找所有分页中可见的条目
	b.DB.Where("user_id = (?)", id).Order("updated_at desc").Find(&userBadges)

	// TODO 分出金银铜铁徽章
	var goldBadges []model.UserBadge
	var silverBadges []model.UserBadge
	var copperBadges []model.UserBadge
	var ironBadges []model.UserBadge

	for _, v := range userBadges {
		var badge model.Badge
		id := v.BadgeId.String()
		// TODO 先看redis中是否存在
		if ok, _ := b.Redis.HExists(ctx, "Badge", id).Result(); ok {
			cate, _ := b.Redis.HGet(ctx, "Badge", id).Result()
			json.Unmarshal([]byte(cate), &badge)
		} else {
			// TODO 查看徽章是否在数据库中存在
			b.DB.Where("id = (?)", id).First(&badge)
			// TODO 将徽章存入redis供下次使用
			v, _ := json.Marshal(badge)
			b.Redis.HSet(ctx, "Badge", id, v)
		}
		if v.MaxScore >= badge.Gold {
			goldBadges = append(goldBadges, v)
		} else if v.MaxScore >= badge.Silver {
			silverBadges = append(silverBadges, v)
		} else if v.MaxScore >= badge.Copper {
			copperBadges = append(copperBadges, v)
		} else if v.MaxScore >= badge.Iron {
			ironBadges = append(ironBadges, v)
		}
	}

	// TODO 返回数据
	response.Success(ctx, gin.H{"goldBadges": goldBadges, "sliverBadges": silverBadges, "copperBadges": copperBadges, "ironBadges": ironBadges}, "成功")
}

// @title    UserShow
// @description   获取指定用户的指定徽章
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (b BadgeController) UserShow(ctx *gin.Context) {

	// TODO 获取path中的id
	user_id := ctx.Params.ByName("user")

	// TODO 获取path中的id
	badge_id := ctx.Params.ByName("badge")

	var badge model.UserBadge

	// TODO 查看徽章是否在数据库中存在
	if b.DB.Where("user_id = (?) and badge_id = (?)", user_id, badge_id).First(&badge).Error != nil {
		response.Fail(ctx, nil, "徽章不存在")
		return
	}

	response.Success(ctx, gin.H{"badge": badge}, "成功")
}

// @title    BehaviorList
// @description   查看变量列表
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (b BadgeController) BehaviorList(ctx *gin.Context) {
	keys := make([]string, len(Handle.Behaviors))
	j := 0
	for k := range Handle.Behaviors {
		keys[j] = k
		j++
	}

	response.Success(ctx, gin.H{"keys": keys}, "成功")
}

// @title    BehaviorShow
// @description   查看变量描述
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (b BadgeController) BehaviorShow(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	description, ok := Handle.Behaviors[id]

	if !ok {
		response.Fail(ctx, nil, "变量不存在")
		return
	}
	response.Success(ctx, gin.H{"description": description}, "成功")
}

// @title    EvaluateExpression
// @description   获取指定用户的表达式得分
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (b BadgeController) EvaluateExpression(ctx *gin.Context) {

	// TODO 获取path中的id
	user_id := ctx.Params.ByName("user")

	// TODO 获取path中的id
	expression := ctx.Params.ByName("expression")

	id, _ := uuid.FromString(user_id)

	score, err := Handle.EvaluateExpression(expression, id)

	response.Success(ctx, gin.H{"score": score, "err": err}, "成功")
}

// @title    Publish
// @description   用户连接
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (b BadgeController) Publish(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 订阅消息
	pubSub := b.Redis.Subscribe(ctx, "BadgePublish"+user.ID.String())
	defer pubSub.Close()
	// TODO 获得消息管道
	ch := pubSub.Channel()

	// TODO 升级get请求为webSocket协议
	ws, err := b.UpGrader.Upgrade(ctx.Writer, ctx.Request, nil)

	if err != nil {
		return
	}
	defer ws.Close()
	// TODO 监听消息
	for msg := range ch {
		var badgePublish vo.BadgePublish
		json.Unmarshal([]byte(msg.Payload), &badgePublish)
		// TODO 断开连接
		if err := ws.WriteJSON(badgePublish); err != nil {
			break
		}
	}
}

// @title    NewBadgeController
// @description   新建一个IBadgeController
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   IBadgeController		返回一个IBadgeController用于调用各种函数
func NewBadgeController() IBadgeController {
	db := common.GetDB()
	db.AutoMigrate(model.Badge{})
	db.AutoMigrate(model.UserBadge{})
	db.AutoMigrate(model.Behavior{})
	db.AutoMigrate(model.BadgeBehavior{})
	redis := common.GetRedisClient(0)
	upGrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	Handle.Behaviors = map[string]Interface.BehaviorInterface{
		"Accepts":               Handle.NewAccepts(),
		"Days":                  Handle.NewDays(),
		"Likes":                 Handle.NewLikes(),
		"Collects":              Handle.NewCollects(),
		"AK":                    Handle.NewAK(),
		"BasicAlgorithm":        Handle.NewBasicAlgorithm(),
		"ComputationalGeometry": Handle.NewComputationalGeometry(),
		"DataStructure":         Handle.NewDataStructure(),
		"DynamicProgramming":    Handle.NewDynamicProgramming(),
		"GraphTheory":           Handle.NewGraphTheory(),
		"NumberTheory":          Handle.NewNumberTheory(),
		"Search":                Handle.NewSearch(),
		"Explore":               Handle.NewExplore(),
		"Hack":                  Handle.NewHack(),
		"Group":                 Handle.NewGroup(),
		"Single":                Handle.NewSingle(),
		"Match":                 Handle.NewMatch(),
		"OI":                    Handle.NewOI(),
	}
	return BadgeController{DB: db, Redis: redis, UpGrader: upGrader}
}
