package timer

import (
	"MGA_OJ/common"
	"MGA_OJ/controller"
	"MGA_OJ/model"
	"MGA_OJ/util"
	"context"
	"time"

	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

// @title    StartTimer
// @description   建立一个比赛开始定时器
// @auth      MGAronya       2022-9-16 12:23
// @param    competitionId uuid.UUID	比赛id
// @return   void
func CompetitionStart() {
	ctx := context.Background()
	redis := common.GetRedisClient(0)
	db := common.GetDB()
	for competition := range controller.CompetitionChan {
		go StartTimer(ctx, redis, db, competition)
	}
}

// @title    StartTimer
// @description   建立一个比赛开始定时器
// @auth      MGAronya       2022-9-16 12:23
// @param    competitionId uuid.UUID	比赛id
// @return   void
func StartTimer(ctx context.Context, redis *redis.Client, db *gorm.DB, competition model.Competition) {

	util.TimerMap[competition.ID] = time.NewTimer(time.Until(time.Time(competition.StartTime)))
	// TODO 等待比赛开始
	<-util.TimerMap[competition.ID].C
	// TODO 比赛初始事项
	controller.InitCompetition[competition.Type](ctx, redis, db, competition)

	// TODO 创建比赛结束定时器
	util.TimerMap[competition.ID] = time.NewTimer(time.Until(time.Time(competition.EndTime)))

	// TODO 等待比赛结束
	<-util.TimerMap[competition.ID].C

	// TODO 等待hack时间结束
	if competition.HackTime.After(competition.EndTime) {
		util.TimerMap[competition.ID] = time.NewTimer(time.Until(time.Time(competition.HackTime)))
		<-util.TimerMap[competition.ID].C
	}

	controller.FinishCompetition[competition.Type](ctx, redis, db, competition)
}
