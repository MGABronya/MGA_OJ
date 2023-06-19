package timer

import (
	"MGA_OJ/common"
	"MGA_OJ/model"
	"context"
)

func HotStatics() {
	LikeStatics()
	UnLikeStatics()
	CollectStatics()
	VisitStatics()
}

func LikeStatics() {
	redis := common.GetRedisClient(0)
	db := common.GetDB()
	ctx := context.Background()

	members, _ := redis.ZRangeWithScores(ctx, "UserLike", 0, -1).Result()

	for i := range members {
		var user model.User
		if db.Where("id = (?)", members[i].Member).First(&user).Error != nil {
			continue
		}
		user.LikeNum += int(members[i].Score)
		db.Save(&user)
	}
	redis.Del(ctx, "UserLike")
}

func UnLikeStatics() {
	redis := common.GetRedisClient(0)
	db := common.GetDB()
	ctx := context.Background()

	members, _ := redis.ZRangeWithScores(ctx, "UserUnLike", 0, -1).Result()

	for i := range members {
		var user model.User
		if db.Where("id = (?)", members[i].Member).First(&user).Error != nil {
			continue
		}
		user.UnLikeNum += int(members[i].Score)
		db.Save(&user)
	}
	redis.Del(ctx, "UserUnLike")
}

func CollectStatics() {
	redis := common.GetRedisClient(0)
	db := common.GetDB()
	ctx := context.Background()

	members, _ := redis.ZRangeWithScores(ctx, "UserCollect", 0, -1).Result()

	for i := range members {
		var user model.User
		if db.Where("id = (?)", members[i].Member).First(&user).Error != nil {
			continue
		}
		user.CollectNum += int(members[i].Score)
		db.Save(&user)
	}
	redis.Del(ctx, "UserCollect")
}

func VisitStatics() {
	redis := common.GetRedisClient(0)
	db := common.GetDB()
	ctx := context.Background()

	members, _ := redis.ZRangeWithScores(ctx, "UserVisit", 0, -1).Result()

	for i := range members {
		var user model.User
		if db.Where("id = (?)", members[i].Member).First(&user).Error != nil {
			continue
		}
		user.VisitNum += int(members[i].Score)
		db.Save(&user)
	}
	redis.Del(ctx, "UserVisit")
}
