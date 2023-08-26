// @Title  Accepts
// @Description  该文件提供关于Accepts行为的各种方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package Handle

import (
	"MGA_OJ/Interface"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/vo"
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v9"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Accepts			定义了Accepts行为类
type Accepts struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
	ctx   context.Context
}

// @title    UserBehavior
// @description   查看用户做题ac次数
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a Accepts) UserBehavior(userId uuid.UUID) (float64, error) {
	var behavior model.Behavior
	if a.DB.Where("name = ? and user_id = ?", "Accepts", userId).First(&behavior).Error != nil {
		return 0, nil
	}
	return behavior.Score, nil
}

// @title    PublishBehavior
// @description   更新行为统计，并按情况通报
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a Accepts) PublishBehavior(score float64, userId uuid.UUID) error {
	var behavior model.Behavior
	// TODO 如果没有，直接创建
	if a.DB.Where("name = ? and user_id = ?", "Accepts", userId).First(&behavior).Error != nil {
		behavior.Name = "Accepts"
		behavior.UserId = userId
		behavior.Score = 0
		a.DB.Create(&behavior)
	}
	behavior.Score += score
	// TODO 更新值
	a.DB.Where("name = ? and user_id = ?", "Accepts", userId).Save(&behavior)
	var badgeBehaviors []model.BadgeBehavior
	a.DB.Where("name = ?", "Accepts").Find(&badgeBehaviors)
	// TODO 遍历映射关系，并检查是否更新
	for i := range badgeBehaviors {
		var badge model.Badge
		if err := a.DB.Where("id = ?", badgeBehaviors[i].BadgeId).First(&badge).Error; err != nil {
			log.Println(err)
			continue
		}
		var userBadge model.UserBadge
		if err := a.DB.Where("badge_id = ?", badgeBehaviors[i].BadgeId).First(&userBadge).Error; err != nil {
			log.Println(err)
			continue
		}
		score, err := EvaluateExpression(badge.Condition, userId)
		if err != nil {
			log.Println(err)
			continue
		}
		// TODO 如果需要更新最大值
		if userBadge.MaxScore < score {
			// TODO 查看是否需要发布订阅
			old, new := "", ""
			if userBadge.MaxScore >= badge.Gold {
				old = "Gold"
			} else if userBadge.MaxScore >= badge.Silver {
				old = "Silver"
			} else if userBadge.MaxScore >= badge.Copper {
				old = "Copper"
			} else if userBadge.MaxScore >= badge.Iron {
				old = "Iron"
			}

			if score >= badge.Gold {
				new = "Gold"
			} else if score >= badge.Silver {
				new = "Silver"
			} else if score >= badge.Copper {
				new = "Copper"
			} else if score >= badge.Iron {
				new = "Iron"
			}

			// TODO 如果获得了新的徽章
			if new != old {
				badgePublish := vo.BadgePublish{
					Name:  badge.Name,
					Level: new,
				}
				v, _ := json.Marshal(badgePublish)
				a.Redis.Publish(a.ctx, "BadgePublish"+userId.String(), v)
			}
			userBadge.MaxScore = score
			a.DB.Save(&userBadge)
		}
	}
	return nil
}

// @title    Description
// @description   返回行为描述
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a Accepts) Description() string {
	return "用户AC题目的总数"
}

// @title    NewAccepts
// @description   新建一个BeahviorInterface
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   BeahviorInterface 		返回一个BeahviorInterface用于调用各种函数
func NewAccepts() Interface.BehaviorInterface {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	ctx := context.Background()
	return Accepts{DB: db, Redis: redis, ctx: ctx}
}
