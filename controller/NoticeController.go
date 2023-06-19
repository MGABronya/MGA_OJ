// @Title  NoticeController
// @Description  该文件提供关于操作个人比赛的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"
	"MGA_OJ/vo"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// INoticeController			定义了个人比赛类接口
type INoticeController interface {
	Create(ctx *gin.Context)   // 通告
	Publish(ctx *gin.Context)  // 订阅通告
	PageList(ctx *gin.Context) // 公告列表
	Show(ctx *gin.Context)     // 指定公告的内容
}

// NoticeController			定义了个人比赛工具类
type NoticeController struct {
	DB       *gorm.DB            // 含有一个数据库指针
	Redis    *redis.Client       // 含有一个redis指针
	UpGrader *websocket.Upgrader // 用于持久化连接
}

// @title    Create
// @description  发布通知
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (n NoticeController) Create(ctx *gin.Context) {

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找比赛
	var competition model.Competition

	// TODO 查看比赛是否存在
	// TODO 先看redis中是否存在
	if ok, _ := n.Redis.HExists(ctx, "Competition", id).Result(); ok {
		cate, _ := n.Redis.HGet(ctx, "Competition", id).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			n.Redis.HDel(ctx, "Competition", id)
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if n.DB.Where("id = (?)", id).First(&competition).Error != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
	{
		// TODO 将比赛存入redis供下次使用
		v, _ := json.Marshal(competition)
		n.Redis.HSet(ctx, "Competition", id, v)
	}
leep:

	// TODO 查看比赛是否在进行中
	if !time.Now().After(time.Time(competition.StartTime)) || !time.Now().Before(time.Time(competition.EndTime)) {
		response.Fail(ctx, nil, "比赛不在进行中")
		return
	}

	// TODO 查看是否有权给比赛添加题目
	if competition.UserId != user.ID {
		if n.DB.Where("group_id = (?) and user_id = (?)", competition.GroupId, user.ID).First(&model.UserList{}).Error != nil {
			response.Fail(ctx, nil, "无权为比赛添加公告")
			return
		}
	}

	// TODO 接收广播内容
	var notic vo.Notice
	// TODO 数据验证
	if err := ctx.ShouldBind(&notic); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}
	// TODO 存入数据库
	notice := model.Notice{
		UserId:        user.ID,
		CompetitionId: competition.ID,
		Title:         notic.Title,
		Content:       notic.Content,
		ResLong:       notic.ResLong,
		ResShort:      notic.ResShort,
	}
	n.DB.Save(&notice)
	// TODO 将notic打包
	v, _ := json.Marshal(notice)
	// TODO 推入频道
	n.Redis.Publish(ctx, "NoticeChan"+id, v)
	n.Redis.HSet(ctx, "Notice", notice.ID, v)
	response.Success(ctx, gin.H{"notice": notice}, "发布成功")
}

// @title    Publish
// @description  订阅通知
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (n NoticeController) Publish(ctx *gin.Context) {

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 查找比赛
	var competition model.Competition

	// TODO 查看比赛是否存在
	// TODO 先看redis中是否存在
	if ok, _ := n.Redis.HExists(ctx, "Competition", id).Result(); ok {
		cate, _ := n.Redis.HGet(ctx, "Competition", id).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			n.Redis.HDel(ctx, "Competition", id)
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if n.DB.Where("id = (?)", id).First(&competition).Error != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
	{
		// TODO 将比赛存入redis供下次使用
		v, _ := json.Marshal(competition)
		n.Redis.HSet(ctx, "Competition", id, v)
	}
leep:

	// TODO 查看比赛是否在进行中
	if !time.Now().After(time.Time(competition.StartTime)) || !time.Now().Before(time.Time(competition.EndTime)) {
		response.Fail(ctx, nil, "比赛不在进行中")
		return
	}

	// TODO 订阅消息
	pubSub := n.Redis.Subscribe(ctx, "NoticeChan"+id)
	defer pubSub.Close()
	// TODO 获得消息管道
	ch := pubSub.Channel()

	// TODO 升级get请求为webSocket协议
	ws, err := n.UpGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	// TODO 监听消息
	for msg := range ch {
		var notice model.Notice
		json.Unmarshal([]byte(msg.Payload), &notice)
		// TODO 写入ws数据
		if err := ws.WriteJSON(notice); err != nil {
			break
		}
	}
}

// @title    Show
// @description  查看通告
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (n NoticeController) Show(ctx *gin.Context) {

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 查找公告
	var notice model.Notice

	// TODO 查看公告是否存在
	// TODO 先看redis中是否存在
	if ok, _ := n.Redis.HExists(ctx, "Notice", id).Result(); ok {
		cate, _ := n.Redis.HGet(ctx, "Notice", id).Result()
		if json.Unmarshal([]byte(cate), &notice) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			n.Redis.HDel(ctx, "Notice", id)
		}
	}

	// TODO 查看公告是否在数据库中存在
	if n.DB.Where("id = (?)", id).First(&notice).Error != nil {
		response.Fail(ctx, nil, "公告不存在")
		return
	}
	{
		// TODO 将公告存入redis供下次使用
		v, _ := json.Marshal(notice)
		n.Redis.HSet(ctx, "Notice", id, v)
	}
leep:

	response.Success(ctx, gin.H{"notice": notice}, "成功")
}

// @title    PageList
// @description  通知列表
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (n NoticeController) PageList(ctx *gin.Context) {

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	var notices []model.Notice

	n.DB.Where("competition_id = (?)", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&notices)
	var total int64

	n.DB.Where("competition_id = (?)", id).Model(model.Notice{}).Count(&total)

	response.Success(ctx, gin.H{"notices": notices, "total": total}, "成功")
}

// @title    NewCompetitionController
// @description   新建一个INoticeController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   INoticeController		返回一个INoticeController用于调用各种函数
func NewNoticeController() INoticeController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	upGrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	db.AutoMigrate(model.Notice{})
	return NoticeController{DB: db, Redis: redis, UpGrader: upGrader}
}
