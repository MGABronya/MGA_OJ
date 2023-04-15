// @Title  CompetitionMatchController
// @Description  该文件提供关于操作个人比赛的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"MGA_OJ/Interface"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ICompetitionMatchController			定义了匹配比赛类接口
type ICompetitionMatchController interface {
	Interface.EnterInterface // 包含报名方法
}

// CompetitionMatchController			定义了个人比赛工具类
type CompetitionMatchController struct {
	DB       *gorm.DB            // 含有一个数据库指针
	Redis    *redis.Client       // 含有一个redis指针
	UpGrader *websocket.Upgrader // 用于持久化连接
	RabbitMq *common.RabbitMQ    // 一个消息队列的指针
}

// @title    Enter
// @description   报名一篇比赛
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionMatchController) Enter(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应比赛
	id := ctx.Params.ByName("id")

	var competition model.Competition
	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Competition", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Competition", id).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			// TODO 跳过数据库搜寻competition过程
			goto leap
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Competition", id)
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if c.DB.Where("id = ?", id).First(&competition) != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(competition)
		c.Redis.HSet(ctx, "Competition", id, v)
	}

leap:

	// TODO 查看比赛是否已经进行过
	if time.Now().After(time.Time(competition.StartTime)) {
		response.Fail(ctx, nil, "比赛已经开始")
		return
	}

	var competitionRank model.CompetitionRank

	// TODO 查看是否已经报名
	if _, err := c.Redis.ZScore(ctx, "CompetitionR"+id, user.ID.String()).Result(); err == nil {
		response.Fail(ctx, nil, "已报名")
		return
	}
	if c.DB.Where("member_id = ? and competition_id = ?", user.ID, competition.ID).First(&competitionRank).Error == nil {
		response.Fail(ctx, nil, "已报名")
		// TODO 加入redis
		c.Redis.ZAdd(ctx, "CompetitionR"+id, redis.Z{Member: user.ID.String(), Score: 0})
		return
	}

	// TODO 查看比赛是否需要密码
	var passwd model.Passwd
	if c.DB.Where("id = ?", competition.PasswdId).First(&passwd).Error == nil {
		var pass model.Passwd
		// TODO 数据验证
		if err := ctx.ShouldBind(&pass); err != nil {
			log.Print(err.Error())
			response.Fail(ctx, nil, "数据验证错误")
			return
		}
		// TODO 判断密码是否正确
		if err := bcrypt.CompareHashAndPassword([]byte(passwd.Password), []byte(pass.Password)); err != nil {
			response.Fail(ctx, nil, "密码错误")
			return
		}
	}

	competitionRank.CompetitionId = competition.ID
	competitionRank.MemberId = user.ID
	competitionRank.Score = 0
	competitionRank.Penalties = 0
	// TODO 插入数据
	if err := c.DB.Create(&competitionRank).Error; err != nil {
		response.Fail(ctx, nil, "数据库存储错误")
		return
	}

	c.Redis.ZAdd(ctx, "CompetitionR"+id, redis.Z{Member: user.ID.String(), Score: 0})
	response.Success(ctx, gin.H{"competitionRank": competitionRank}, "报名成功")
	return
}

// @title    EnterCondition
// @description   查看报名状态
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionMatchController) EnterCondition(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应比赛
	id := ctx.Params.ByName("id")

	var competition model.Competition
	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Competition", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Competition", id).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			// TODO 跳过数据库搜寻competition过程
			goto leap
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Competition", id)
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if c.DB.Where("id = ?", id).First(&competition) != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(competition)
		c.Redis.HSet(ctx, "Competition", id, v)
	}

leap:

	// TODO 检查比赛类型
	if competition.Type != "Match" {
		response.Fail(ctx, nil, "不是匹配比赛，无法报名")
		return
	}

	var competitionRank model.CompetitionRank

	// TODO 查看是否已经报名
	if _, err := c.Redis.ZScore(ctx, "CompetitionR"+id, user.ID.String()).Result(); err == nil {
		response.Success(ctx, gin.H{"enter": true}, "已报名")
		return
	}
	if c.DB.Where("member_id = ? and competition_id = ?", user.ID, competition.ID).First(&competitionRank).Error == nil {
		response.Success(ctx, gin.H{"enter": true}, "已报名")
		{
			// TODO 加入redis
			c.Redis.ZAdd(ctx, "CompetitionR"+id, redis.Z{Member: user.ID.String(), Score: 0})
		}
		return
	}

	response.Success(ctx, nil, "未报名")
	return

}

// @title    CancelEnter
// @description   取消报名一篇比赛
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionMatchController) CancelEnter(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应比赛
	id := ctx.Params.ByName("id")

	var competition model.Competition
	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Competition", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Competition", id).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			// TODO 跳过数据库搜寻competition过程
			goto leap
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Competition", id)
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if c.DB.Where("id = ?", id).First(&competition) != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(competition)
		c.Redis.HSet(ctx, "Competition", id, v)
	}

leap:
	if competition.Type != "Match" {
		response.Fail(ctx, nil, "不是匹配比赛，无法报名")
		return
	}

	// TODO 查看比赛是否已经进行过
	if time.Now().After(time.Time(competition.StartTime)) {
		response.Fail(ctx, nil, "比赛已经开始")
		return
	}

	var competitionRank model.CompetitionRank

	// TODO 查看是否已经报名
	if c.DB.Where("member_id = ? and competition_id = ?", user.ID, competition.ID).First(&competitionRank).Error != nil {
		response.Fail(ctx, nil, "未报名")
		return
	}

	c.DB.Delete(&competitionRank)
	c.Redis.ZRem(ctx, "CompetitionR"+id, user.ID.String())
	response.Success(ctx, nil, "取消报名成功")
	return
}

// @title    EnterPage
// @description   查看报名列表
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionMatchController) EnterPage(ctx *gin.Context) {

	// TODO 查找对应比赛
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	var competition model.Competition
	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Competition", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Competition", id).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			// TODO 跳过数据库搜寻competition过程
			goto leap
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Competition", id)
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if c.DB.Where("id = ?", id).First(&competition) != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(competition)
		c.Redis.HSet(ctx, "Competition", id, v)
	}

leap:

	var competitionRanks []model.CompetitionRank

	// TODO 查找所有分页中可见的条目
	c.DB.Where("competition_id = ?", competition.ID).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&competitionRanks)

	var total int64
	c.DB.Where("competition_id = ?", competition.ID).Model(model.CompetitionRank{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"competitionRanks": competitionRanks, "total": total}, "成功")
}

// @title    NewCompetitionController
// @description   新建一个ICompetitionMatchController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   ICompetitionMatchController		返回一个ICompetitionMatchController用于调用各种函数
func NewCompetitionMatchController() ICompetitionMatchController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	rabbitmq := common.GetRabbitMq()
	upGrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	InitCompetition["Match"] = initMatchCompetition
	FinishCompetition["Match"] = finishMatchCompetition
	return CompetitionMatchController{DB: db, Redis: redis, UpGrader: upGrader, RabbitMq: rabbitmq}
}

func initMatchCompetition(ctx *gin.Context, redis *redis.Client, db *gorm.DB, competition model.Competition) {
	members, _ := redis.ZRange(ctx, "CompetitionR"+competition.ID.String(), 0, -1).Result()
	var users []model.User
	for i := range members {
		var user model.User

		// TODO 先看redis中是否存在
		if ok, _ := redis.HExists(ctx, "User", members[i]).Result(); ok {
			cate, _ := redis.HGet(ctx, "User", members[i]).Result()
			if json.Unmarshal([]byte(cate), &user) == nil {
				goto leap
			} else {
				// TODO 移除损坏数据
				redis.HDel(ctx, "User", members[i])
			}
		}

		// TODO 查看用户是否在数据库中存在
		if db.Where("id = ?", members[i]).First(&user).Error != nil {
			continue
		}
		{
			// TODO 将用户存入redis供下次使用
			v, _ := json.Marshal(user)
			redis.HSet(ctx, "User", members[i], v)
		}
	leap:
		users = append(users, user)
	}
	sort.Slice(users, func(i, j int) bool { // 根据Score降序
		return users[i].Score > users[j].Score
	})
	// TODO 计算需要分几组
	n := math.Ceil(float64(int(len(users))) / float64(competition.UpNum))

	index := 0

	// TODO 分配组的创建
	groups := make([]model.Group, 0)
	for i := 0; i < int(n); i++ {
		var group = model.Group{
			Title:    competition.Title + "-" + users[index].Name + "'s Group",
			Content:  fmt.Sprint("匹配组", i+1),
			LeaderId: users[index].ID,
		}
		// TODO 小组在比赛正式结束前无法更换组员
		if competition.HackTime.After(competition.EndTime) {
			group.CompetitionAt = competition.HackTime
		} else {
			group.CompetitionAt = competition.EndTime
		}
		db.Create(&group)
		groups = append(groups, group)
		var userList = model.UserList{
			GroupId: group.ID,
			UserId:  users[index].ID,
		}
		db.Create(&userList)
		index++
	}
	// TODO 进行分配
	i := int(n - 1)
	flag := true
	for index < len(users) {
		var userList = model.UserList{
			GroupId: groups[i].ID,
			UserId:  users[index].ID,
		}
		db.Create(&userList)
		if i == int(n-1) {
			flag = false
		} else if i == 0 {
			flag = true
		}
		if flag {
			i++
		} else {
			i--
		}
		index++
	}
}

func finishMatchCompetition(ctx *gin.Context, redis *redis.Client, db *gorm.DB, competition model.Competition) {
	if competition.Type != "Match" {
		log.Println("match competition's type is error!")
	} else {
		log.Println("match competition finish!", competition)
	}
	CompetitionProblemSubmit(ctx, redis, db, competition)
	CompetitionFinish(ctx, redis, db, competition)
}
