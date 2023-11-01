// @Title  HeartController
// @Description  该文件提供关于操作heart的各种方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package controller

import (
	"MGA_OJ/Interface"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// IHeartController			定义了heart类接口
type IHeartController interface {
	Interface.HeartInterface // 包含了心跳相关功能
}

// HeartController			定义了心跳工具类
type HeartController struct {
	DB       *gorm.DB            // 含有一个数据库指针
	Redis    *redis.Client       // 含有一个redis指针
	UpGrader *websocket.Upgrader // 用于持久化连接
}

// HeartPercentageMap			心跳百分比
var HeartPercentageMap map[string]float64

// @title    Show
// @description 查看最近心跳情况
// @auth      MGAronya       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (h HeartController) Show(ctx *gin.Context) {
	id := ctx.Params.ByName("id")
	start := ctx.Params.ByName("start")
	end := ctx.Params.ByName("end")

	// TODO 找到指定时间段中的心跳
	var hearts []model.Heart
	h.DB.Where("docker_id = (?) and times_tamp >= (?) and times_tamp <= (?)", id, start, end).Find(&hearts)
	// TODO 返回数据
	response.Success(ctx, gin.H{"hearts": hearts}, "成功")
}

// @title    Publish
// @description   订阅心跳长连接
// @auth      MGAronya       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (h HeartController) Publish(ctx *gin.Context) {
	id := ctx.Params.ByName("id")

	// TODO 订阅消息
	pubSub := h.Redis.Subscribe(ctx, "heart")
	defer pubSub.Close()
	// TODO 获得消息管道
	ch := pubSub.Channel()

	// TODO 升级get请求为webSocket协议
	ws, err := h.UpGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	// TODO 监听消息
	for msg := range ch {
		var heart model.Heart
		json.Unmarshal([]byte(msg.Payload), &heart)
		if heart.DockerId == id {
			// TODO 写入ws数据
			if err := ws.WriteJSON(heart); err != nil {
				break
			}
		}
	}
}

// @title    Percentage
// @description   心跳忙碌占比
// @auth      MGAronya       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (h HeartController) Percentage(ctx *gin.Context) {
	// TODO 返回数据
	response.Success(ctx, gin.H{"HeartPercentageMap": HeartPercentageMap}, "成功")
}

// @title    NewHeartController
// @description   新建一个IHeartController
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   IHeartController		返回一个IHeartController用于调用各种函数
func NewHeartController() IHeartController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.Heart{})
	upGrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return HeartController{DB: db, Redis: redis, UpGrader: upGrader}
}

// @title    HeartCount
// @description   心跳统计
// @auth      MGAronya       2022-9-16 12:19
// @param    void
// @return   void
func HeartCount() {
	// TODO 记录心跳百分比
	go HeartPercentage()
	// TODO 记录心跳情况
	db := common.GetDB()
	redi := common.GetRedisClient(0)
	ctx := context.Background()
	// TODO 订阅消息
	pubSub := redi.Subscribe(ctx, "heart")
	defer pubSub.Close()
	// TODO 获得消息管道
	ch := pubSub.Channel()

	// TODO 监听消息
	for msg := range ch {
		var heart model.Heart
		json.Unmarshal([]byte(msg.Payload), &heart)
		db.Create(&heart)
	}
}

// @title    HeartPercentage
// @description   记录心跳百分比
// @auth      MGAronya       2022-9-16 12:19
// @param    void
// @return   void
func HeartPercentage() {
	db := common.GetDB()
	for {
		time.Sleep(10 * time.Second)
		var hearts []model.Heart
		// TODO 查找出最近10s内的心跳
		now := time.Now()
		db.Where("times_tamp between date_sub(now(),interval 10 SECOND) and now()").Order("times_tamp asc").Find(&hearts)
		var times map[string][]model.Heart = make(map[string][]model.Heart)
		var percentage map[string]float64 = make(map[string]float64)
		// TODO 按容器分类
		for i := range hearts {
			if len(times[hearts[i].DockerId]) == 0 {
				times[hearts[i].DockerId] = make([]model.Heart, 0)
			}
			times[hearts[i].DockerId] = append(times[hearts[i].DockerId], hearts[i])
		}
		// TODO 计算忙碌占比
		for id := range times {
			for i := 0; i < len(times[id]); i++ {
				if i == 0 && times[id][i].Condition == "Finish" {
					percentage[id] += float64(time.Time(times[id][i].TimesTamp).Sub(now)) + 10
				} else if times[id][i].Condition == "Running" {
					if i == len(times[id])-1 {
						percentage[id] += float64(now.Sub(time.Time(times[id][i].TimesTamp)))
					} else {
						percentage[id] += float64(time.Time(times[id][i+1].TimesTamp).Sub(time.Time(times[id][i].TimesTamp)))
					}
				}
			}
			percentage[id] /= 10
		}
		// TODO 赋值给新的百分比映射
		HeartPercentageMap = percentage
	}
}
