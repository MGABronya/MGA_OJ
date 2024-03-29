// @Title  CompetitionRandomGroupController
// @Description  该文件提供关于操作个人比赛的各种方法
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
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// TODO 比赛开始计时器
var competitionGroupStart time.Timer

// ICompetitionRandomGroupController			定义了个人比赛类接口
type ICompetitionRandomGroupController interface {
	Interface.EnterInterface       // 包含报名方法
	EnterPublish(ctx *gin.Context) // 实时报告情况
}

// CompetitionRandomGroupController			定义了个人比赛工具类
type CompetitionRandomGroupController struct {
	DB       *gorm.DB            // 含有一个数据库指针
	Redis    *redis.Client       // 含有一个redis指针
	UpGrader *websocket.Upgrader // 用于持久化连接
	RabbitMq *common.RabbitMQ    // 一个消息队列的指针
}

// @title    Enter
// @description   报名一篇比赛
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionRandomGroupController) Enter(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否已经报名
	if _, err := c.Redis.ZScore(ctx, "CompetitionRGroup", user.ID.String()).Result(); err == nil {
		response.Fail(ctx, nil, "已报名")
		return
	}

	ck := redis.Z{Member: user.ID.String(), Score: float64(time.Now().Unix())}

	// TODO 加入redis
	c.Redis.ZAdd(ctx, "CompetitionRGroup", ck)
	// TODO 加入频道
	v, _ := json.Marshal(ck)
	c.Redis.Publish(ctx, "CompetitionRandomGroupChan", v)
	response.Success(ctx, nil, "报名成功")
	// TODO 查看当前比赛人数是否可以直接开始比赛
	cks, _ := c.Redis.ZRangeWithScores(ctx, "CompetitionRGroup", 0, 0).Result()
	total, _ := c.Redis.ZCard(ctx, "CompetitionRGroup").Result()
	if CanStartGroupCompetition(total, cks[0].Score) {
		competitionGroupStart.Reset(0)
	} else {
		// TODO 更新比赛开始计时器
		competitionGroupStart.Reset(time.Duration(cks[0].Score+float64(10*time.Minute)) - time.Duration(time.Now().Unix()))
	}
}

// @title    EnterPublish
// @description   实时报告情况
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionRandomGroupController) EnterPublish(ctx *gin.Context) {
	// TODO 订阅消息
	pubSub := c.Redis.Subscribe(ctx, "CompetitionRandomGroupChan")
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
		var enter redis.Z
		json.Unmarshal([]byte(msg.Payload), &enter)
		// TODO 写入ws数据
		if err := ws.WriteJSON(enter); err != nil {
			break
		}

	}
}

// @title    EnterCondition
// @description   查看报名状态
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionRandomGroupController) EnterCondition(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否已经报名
	if _, err := c.Redis.ZScore(ctx, "CompetitionRGroup", user.ID.String()).Result(); err == nil {
		response.Success(ctx, gin.H{"enter": true}, "已报名")
		return
	}

	response.Success(ctx, gin.H{"enter": false}, "未报名")
	return

}

// @title    CancelEnter
// @description   取消报名一篇比赛
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionRandomGroupController) CancelEnter(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否已经报名
	if _, err := c.Redis.ZScore(ctx, "CompetitionRGroup", user.ID.String()).Result(); err != nil {
		response.Success(ctx, gin.H{"enter": true}, "未报名")
		return
	}
	c.Redis.ZRem(ctx, "CompetitionRGroup", user.ID.String())
	response.Success(ctx, nil, "取消报名成功")
	// TODO 加入频道
	ck := redis.Z{Member: user.ID.String(), Score: 0}
	v, _ := json.Marshal(ck)
	c.Redis.Publish(ctx, "CompetitionRandomGroupChan", v)
	// TODO 查看当前比赛人数是否可以直接开始比赛
	cks, _ := c.Redis.ZRangeWithScores(ctx, "CompetitionRGroup", 0, 0).Result()
	total, _ := c.Redis.ZCard(ctx, "CompetitionRGroup").Result()
	if CanStartGroupCompetition(total, cks[0].Score) {
		competitionGroupStart.Reset(0)
	} else {
		// TODO 更新比赛开始计时器
		competitionGroupStart.Reset(time.Duration(cks[0].Score+float64(10*time.Minute)) - time.Duration(time.Now().Unix()))
	}
}

// @title    EnterPage
// @description   查看报名列表
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CompetitionRandomGroupController) EnterPage(ctx *gin.Context) {

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 查找所有分页中可见的条目
	competitionRanks, _ := c.Redis.ZRevRangeWithScores(ctx, "CompetitionRGroup", int64(pageNum-1)*int64(pageSize), int64(pageNum-1)*int64(pageSize)+int64(pageSize)-1).Result()

	total, _ := c.Redis.ZCard(ctx, "CompetitionRGroup").Result()

	// TODO 返回数据
	response.Success(ctx, gin.H{"competitionRanks": competitionRanks, "total": total}, "成功")
}

// @title    NewCompetitionRandomGroupController
// @description   新建一个ICompetitionRandomGroupController
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   ICompetitionRandomGroupController		返回一个ICompetitionRandomGroupController用于调用各种函数
func NewCompetitionRandomGroupController() ICompetitionRandomGroupController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	rabbitmq := common.GetRabbitMq()
	upGrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	competitionGroupStart = *time.NewTimer(time.Hour)
	competitionGroupStart.Stop()
	return CompetitionRandomGroupController{DB: db, Redis: redis, UpGrader: upGrader, RabbitMq: rabbitmq}
}

func CanStartGroupCompetition(total int64, t float64) bool {
	t = float64(time.Now().Unix()) - t
	if total < 9 {
		return false
	}
	return total >= 48-int64(t/10)
}

func CompetitionRandomGroupGo() {
	db := common.GetDB()
	redi := common.GetRedisClient(0)
	ctx := context.Background()
	for {
		<-competitionGroupStart.C
		// TODO 取出所有报名用户
		tusers, _ := redi.ZRange(ctx, "CompetitionRGroup", 0, -1).Result()
		if len(tusers) < 9 {
			competitionGroupStart.Reset(10 * time.Minute)
			continue
		}
		competitionGroupStart.Stop()
		redi.Del(ctx, "CompetitionRGroup")
		// TODO 比赛创建
		competition := model.Competition{
			StartTime: model.Time(time.Now()),
			EndTime:   model.Time(time.Now().Add(150 * time.Minute)),
			Title:     "随机匹配组队赛",
			Content:   "随机匹配组队赛",
			HackTime:  model.Time(time.Now().Add(180 * time.Minute)),
			HackScore: 1,
			HackNum:   5,
			Type:      "Match",
			LessNum:   3,
			UpNum:     5,
		}

		// TODO 插入数据
		db.Create(&competition)
		// TODO 比赛初始事项
		// TODO 插入成员
		var users []model.User
		for i := range tusers {
			var user model.User

			// TODO 先看redis中是否存在
			if ok, _ := redi.HExists(ctx, "User", tusers[i]).Result(); ok {
				cate, _ := redi.HGet(ctx, "User", tusers[i]).Result()
				if json.Unmarshal([]byte(cate), &user) == nil {
					goto tuser
				} else {
					// TODO 移除损坏数据
					redi.HDel(ctx, "User", tusers[i])
				}
			}

			// TODO 查看用户是否在数据库中存在
			db.Where("id = (?)", tusers[i]).First(&user)
			{
				// TODO 将用户存入redis供下次使用
				v, _ := json.Marshal(user)
				redi.HSet(ctx, "User", tusers[i], v)
			}
		tuser:
			users = append(users, user)
		}
		// TODO 根据Score降序
		sort.Slice(users, func(i, j int) bool {
			return users[i].Score > users[j].Score
		})
		// TODO 计算需要分几组
		n := math.Floor(float64(int(len(users))) / float64(competition.LessNum))

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
		// TODO 插入题目
		var problems []model.Problem

		// TODO 随机获得五道题目
		db.Order("RAND()").Limit(8).Find(&problems)

		// TODO 插入题目
		for i := range problems {
			var caseSamples []model.CaseSample
			// TODO 先看redis中是否存在
			if ok, _ := redi.HExists(ctx, "CaseSample", problems[i].ID.String()).Result(); ok {
				cate, _ := redi.HGet(ctx, "CaseSample", problems[i].ID.String()).Result()
				if json.Unmarshal([]byte(cate), &caseSamples) == nil {
					// TODO 跳过数据库搜寻caseSample过程
					goto leap
				} else {
					// TODO 移除损坏数据
					redi.HDel(ctx, "CaseSample", problems[i].ID.String())
				}
			}
			db.Where("problem_id = (?)", problems[i].ID).Find(&caseSamples)
			// TODO 将题目存入redis供下次使用
			{
				v, _ := json.Marshal(caseSamples)
				redi.HSet(ctx, "CaseSample", problems[i].ID.String(), v)
			}

		leap:
			// TODO 从数据库中读出输入输出
			var cases []model.Case

			// TODO 查找用例
			if ok, _ := redi.HExists(ctx, "Case", problems[i].ID.String()).Result(); ok {
				cate, _ := redi.HGet(ctx, "Case", problems[i].ID.String()).Result()
				if json.Unmarshal([]byte(cate), &cases) == nil {
					// TODO 跳过数据库搜寻testInputs过程
					goto Case
				} else {
					// TODO 移除损坏数据
					redi.HDel(ctx, "Case", problems[i].ID.String())
				}
			}

			// TODO 查看题目是否在数据库中存在
			db.Where("id = (?)", problems[i]).Find(&cases)
			// TODO 将题目存入redis供下次使用
			{
				v, _ := json.Marshal(cases)
				redi.HSet(ctx, "Case", problems[i].ID.String(), v)
			}
		Case:

			// TODO 创建题目
			problemNew := model.ProblemNew{
				Title:         problems[i].Title,
				TimeLimit:     problems[i].TimeLimit,
				MemoryLimit:   problems[i].MemoryLimit,
				Description:   problems[i].Description,
				ResLong:       problems[i].ResLong,
				ResShort:      problems[i].ResShort,
				Input:         problems[i].Input,
				Output:        problems[i].Output,
				Hint:          problems[i].Hint,
				Source:        problems[i].Source,
				UserId:        problems[i].UserId,
				SpecialJudge:  problems[i].SpecialJudge,
				CompetitionId: competition.ID,
			}

			// TODO 插入数据
			db.Create(&problemNew)

			// TODO 存储测试样例
			for i := range caseSamples {
				// TODO 尝试存入数据库
				cas := model.CaseSample{
					ProblemId: problemNew.ID,
					Input:     caseSamples[i].Input,
					Output:    caseSamples[i].Output,
					CID:       uint(i + 1),
				}
				// TODO 插入数据
				db.Create(&cas)
			}

			// TODO 存储测试用例
			for i := range cases {
				// TODO 尝试存入数据库
				cas := model.Case{
					ProblemId: problemNew.ID,
					Input:     cases[i].Input,
					Output:    cases[i].Output,
					Score:     0,
					CID:       uint(i + 1),
				}
				if i == len(cases)-1 {
					cas.Score = 5
				}
				// TODO 插入数据
				db.Create(&cas)
			}

		}

		// TODO 等待比赛结束
		<-time.NewTimer(150 * time.Minute).C

		// TODO 等待hack时间结束
		<-time.NewTimer(30 * time.Minute).C

		// TODO 整理比赛结果
		CompetitionFinish(ctx, redi, db, competition)

	}
}
