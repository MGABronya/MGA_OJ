package consumer

import (
	Handle "MGA_OJ/Behavior"
	"MGA_OJ/Interface"
	TQ "MGA_OJ/Test-request"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/util"
	"MGA_OJ/vo"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/go-redis/redis/v9"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// IGroup			定义了判题类接口
type IGroup interface {
	Interface.ConsumerInterface // 包含消费功能
}

// Group			定义了判断工具类
type Group struct {
	rw    *sync.RWMutex // 含有锁
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
	ctx   context.Context
}

// @title    Handle
// @description   创建一篇判断
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g Group) Handel(msg string) {
	// TODO 单核处理，上锁
	g.rw.Lock()
	// TODO 确保资源归还
	defer g.rw.Unlock()
	var record model.RecordCompetition
	// TODO 先看redis中是否存在
	if ok, _ := g.Redis.HExists(g.ctx, "RecordCompetition", msg).Result(); ok {
		cate, _ := g.Redis.HGet(g.ctx, "RecordCompetition", msg).Result()
		// TODO 移除损坏数据
		g.Redis.HDel(g.ctx, "RecordCompetition", msg)
		if json.Unmarshal([]byte(cate), &record) == nil {
			// TODO 跳过数据库搜寻过程
			goto feep
		}
	}

	// TODO 未能找到提交记录
	if g.DB.Where("id = (?)", msg).First(&record).Error != nil {
		log.Printf("%s Record Disappear!!\n", msg)
		return
	}

feep:
	var problem model.ProblemNew
	// TODO 先看redis中是否存在
	id := fmt.Sprint(record.ProblemId)
	if ok, _ := g.Redis.HExists(g.ctx, "ProblemNew", id).Result(); ok {
		cate, _ := g.Redis.HGet(g.ctx, "ProblemNew", id).Result()
		if json.Unmarshal([]byte(cate), &problem) == nil {
			// TODO 跳过数据库搜寻problem过程
			goto leep
		} else {
			// TODO 移除损坏数据
			g.Redis.HDel(g.ctx, "ProblemNew", id)
		}
	}

	// TODO 查看题目是否在数据库中存在
	if g.DB.Where("id = (?)", id).First(&problem).Error != nil {
		record.Condition = "Problem Doesn't Exist"
		// TODO 将record存入redis
		v, _ := json.Marshal(record)
		g.Redis.HSet(g.ctx, "RecordCompetition", fmt.Sprint(record.ID), v)
		// TODO 将record存入mysql
		g.DB.Save(&record)
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(problem)
		g.Redis.HSet(g.ctx, "Problem", id, v)
	}

leep:
	var competition model.Competition
	// TODO 先看redis中是否存在
	if ok, _ := g.Redis.HExists(g.ctx, "Competition", problem.CompetitionId.String()).Result(); ok {
		cate, _ := g.Redis.HGet(g.ctx, "Competition", problem.CompetitionId.String()).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			// TODO 跳过数据库搜寻competition过程
			goto leap
		} else {
			// TODO 移除损坏数据
			g.Redis.HDel(g.ctx, "Competition", problem.CompetitionId.String())
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if g.DB.Where("id = (?)", problem.CompetitionId.String()).First(&competition).Error != nil {
		record.Condition = "Competition Doesn't Exist"
		// TODO 将record存入redis
		v, _ := json.Marshal(record)
		g.Redis.HSet(g.ctx, "RecordCompetition", fmt.Sprint(record.ID), v)
		// TODO 将record存入mysql
		g.DB.Save(&record)
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(competition)
		g.Redis.HSet(g.ctx, "Competition", problem.CompetitionId.String(), v)
	}

leap:

	// TODO 发布订阅用于提交列表
	recordList := vo.RecordList{
		RecordId: record.ID,
	}
	// TODO 确保信息进入频道
	defer func() {
		// TODO 将recordlist打包
		v, _ := json.Marshal(recordList)
		g.Redis.Publish(g.ctx, "RecordCompetitionChan"+problem.CompetitionId.String(), v)
		// TODO 将record存入redis
		v, _ = json.Marshal(record)
		g.Redis.HSet(g.ctx, "RecordCompetition", fmt.Sprint(record.ID), v)
		// TODO 将record存入mysql
		g.DB.Save(&record)
		os.RemoveAll("./user-code")
		os.MkdirAll("./user-code", 0751)
		cmd1 := exec.Command("pgrep", "-u", "mgaoj")
		cmd2 := exec.Command("xargs", "kill", "-9")

		cmd2.Stdin, _ = cmd1.StdoutPipe()

		if err := cmd2.Start(); err != nil {
			return
		}

		if err := cmd1.Run(); err != nil {
			cmd2.Process.Kill()
			return
		}

		if err := cmd2.Wait(); err != nil {
			return
		}
	}()

	// TODO 一些准备工作
	{
		record.Condition = "Preparing"
		// TODO 将recordlist打包
		v, _ := json.Marshal(recordList)
		g.Redis.Publish(g.ctx, "RecordCompetitionChan"+problem.CompetitionId.String(), v)
		// TODO 将record存入redis
		v, _ = json.Marshal(record)
		g.Redis.HSet(g.ctx, "RecordCompetition", fmt.Sprint(record.ID), v)
	}
	// TODO 查看代码是否为空
	if record.Code == "" {
		record.Condition = "Code is empty"
		return
	}
	// TODO 找到提交记录后，开始判题逻辑
	if cmdI, ok := util.LanguageMap[record.Language]; ok {
		// TODO 从数据库中读出输入输出
		var cases []model.Case

		// TODO 查找用例
		if ok, _ := g.Redis.HExists(g.ctx, "Case", id).Result(); ok {
			cate, _ := g.Redis.HGet(g.ctx, "Case", id).Result()
			if json.Unmarshal([]byte(cate), &cases) == nil {
				// TODO 跳过数据库搜寻testInputs过程
				goto Case
			} else {
				// TODO 移除损坏数据
				g.Redis.HDel(g.ctx, "Case", id)
			}
		}

		// TODO 查看题目是否在数据库中存在
		if g.DB.Where("problem_id = (?)", id).Find(&cases).Error != nil {
			record.Condition = "Input Doesn't Exist"
			return
		}
		// TODO 将题目存入redis供下次使用
		{
			v, _ := json.Marshal(cases)
			g.Redis.HSet(g.ctx, "Case", id, v)
		}
	Case:

		fileId := cmdI.Name()
		fp, err := os.Create("./user-code/" + fileId + "." + cmdI.Suffix())
		// TODO 文件错误
		if err != nil {
			// TODO 创建文件失败的原因有：
			// TODO 1、路径不存在  2、权限不足  3、打开文件数量超过上限  4、磁盘空间不足等
			record.Condition = "System Error 1"
			return
		}

		// TODO 开始编译工作
		{
			record.Condition = "Compiling"
			// TODO 将recordlist打包
			v, _ := json.Marshal(recordList)
			g.Redis.Publish(g.ctx, "RecordCompetitionChan"+problem.CompetitionId.String(), v)
			// TODO 将record存入redis
			v, _ = json.Marshal(record)
			g.Redis.HSet(g.ctx, "RecordCompetition", fmt.Sprint(record.ID), v)
		}

		// TODO defer延迟调用 关闭文件，释放资源
		defer fp.Close()

		// TODO 写入文件时，使用带缓存的 *Writer
		write := bufio.NewWriter(fp)

		write.WriteString(record.Code)

		// TODO Flush将缓存的文件真正写入到文件中
		write.Flush()

		// TODO 编译
		cmd := cmdI.Compile("./user-code/", fileId)

		// TODO 系统错误
		if err := cmd.Start(); err != nil {
			record.Condition = "System Error 2"
			return
		}
		// TODO 启动routine等待结束
		done := make(chan error)
		go func() { done <- cmd.Wait() }()

		// TODO 记录是否通过
		flag := true

		// 设定超时时间，并select它
		after := time.After(time.Duration(20 * time.Second))
		select {
		// TODO 编译超时
		case <-after:
			cmd.Process.Kill()
			record.Condition = "Compile Time Out"
			return
		case err = <-done:
		}

		// TODO 编译出错
		if err != nil {
			record.Condition = "Compile Error"
			return
		}

		// TODO 获取权限
		cmd = cmdI.Chmod("./user-code/", id)

		// TODO 权限错误
		if err := cmd.Start(); err != nil {
			record.Condition = "System Error 6"
			return
		}
		// TODO 启动routine等待结束
		done = make(chan error)
		go func() { done <- cmd.Wait() }()

		// 设定超时时间，并select它
		after = time.After(time.Duration(5 * time.Second))
		select {
		// TODO 权限超时
		case <-after:
			cmd.Process.Kill()
			record.Condition = "Compile Time Out"
			return
		case err = <-done:
		}

		// TODO 编译出错
		if err != nil {
			record.Condition = "Compile Error"
			return
		}

		// TODO 开始运行工作
		{
			record.Condition = "Runing"
			// TODO 将recordlist打包
			v, _ := json.Marshal(recordList)
			g.Redis.Publish(g.ctx, "RecordCompetitionChan"+problem.CompetitionId.String(), v)
			// TODO 将record存入redis
			v, _ = json.Marshal(record)
			g.Redis.HSet(g.ctx, "RecordCompetition", fmt.Sprint(record.ID), v)
		}

		// TODO 最终将所有用例填入cases
		caseConditions := make([]model.CaseCondition, 0)

		defer func() {
			for i := range caseConditions {
				g.DB.Create(caseConditions[i])
			}
		}()

		for i := 0; i < len(cases); i++ {
			// TODO 将用例添加至最终数组
			cas := model.CaseCondition{
				RecordId: record.ID,
				CID:      uint(i + 1),
			}
			caseConditions = append(caseConditions, cas)
			var bm runtime.MemStats
			runtime.ReadMemStats(&bm)

			// TODO 通过沙箱运行可执行文件
			cmd = exec.Command("./seccomp", "-language", record.Language, "-input", cases[i].Input+"\n")

			var out, stderr bytes.Buffer
			cmd.Stderr = &stderr
			cmd.Stdout = &out

			now := time.Now().UnixMilli()
			// TODO 系统错误
			if err := cmd.Start(); err != nil {
				record.Condition = "System Error 4"
				return
			}
			// TODO 启动routine等待结束
			done = make(chan error)
			go func() { done <- cmd.Wait() }()

			// 设定超时时间，并select它
			after = time.After(time.Duration(problem.TimeLimit*cmdI.TimeMultiplier()+cmdI.RunUpTime()) * time.Millisecond)
			select {
			// TODO 运行超时
			case <-after:
				cmd.Process.Kill()
				record.Condition = "Time Limit Exceeded"
				flag = false
				goto final
			case err = <-done:
			}

			// TODO 运行时错误
			if err != nil {
				record.Condition = "Runtime Error"
				flag = false
				goto final
			}

			end := time.Now().UnixMilli()
			var em runtime.MemStats
			runtime.ReadMemStats(&em)
			// TODO 更新用例通过情况
			caseConditions[i].Time = uint(math.Max(float64(end-now-int64(cmdI.RunUpTime())), 0))
			caseConditions[i].Memory = uint(em.Alloc/1024 - bm.Alloc/1024)
			// TODO 超时
			if caseConditions[i].Time > problem.TimeLimit*cmdI.TimeMultiplier() {
				record.Condition = "Time Limit Exceeded"
				flag = false
				goto final
			}
			// TODO 超出内存限制
			if caseConditions[i].Memory > problem.MemoryLimit*cmdI.MemoryMultiplier() {
				record.Condition = "Memory Limit Exceeded"
				flag = false
				goto final
			}
			// TODO 答案错误
			var specalJudge model.Program
			// TODO 查看题目是否有标准程序

			// TODO 先看redis中是否存在
			if ok, _ := g.Redis.HExists(g.ctx, "Program", problem.SpecialJudge.String()).Result(); ok {
				cate, _ := g.Redis.HGet(g.ctx, "Program", problem.SpecialJudge.String()).Result()
				if json.Unmarshal([]byte(cate), &specalJudge) == nil {
					// TODO 跳过数据库搜寻program过程
					goto special
				} else {
					// TODO 移除损坏数据
					g.Redis.HDel(g.ctx, "Program", problem.SpecialJudge.String())
				}
			}

			// TODO 查看程序是否在数据库中存在
			if g.DB.Where("id = (?)", problem.SpecialJudge.String()).First(&specalJudge).Error != nil {
				goto outPut
			}
			// TODO 将题目存入redis供下次使用
			{
				v, _ := json.Marshal(specalJudge)
				g.Redis.HSet(g.ctx, "Program", problem.SpecialJudge.String(), v)
			}
		special:
			// TODO 进行特判
			{
				if condition, output := TQ.JudgeRun(specalJudge.Language, specalJudge.Code, cases[i].Input+"\n"+out.String(), problem.MemoryLimit*3, problem.TimeLimit*3); condition != "ok" || output != "ok" {
					record.Condition = condition
					flag = false
					goto final
				}
				goto pass
			}
		outPut:
			// TODO 正常判断
			if out.String() != cases[i].Output {
				// TODO 去除格式后查看是否正确
				if util.RemoveWhiteSpace(out.String()) == util.RemoveWhiteSpace(cases[i].Output) {
					record.Condition = " Presentation Error"
					flag = false
					goto final
				}
				record.Condition = "Wrong Answer"
				flag = false
				goto final
			}
		pass:
			// TODO 通过数量+1
			record.Pass++

			// TODO 数据库插入数据错误
			if g.DB.Create(&cas).Error != nil {
				record.Condition = "System error 5"
				return
			}
		}
	final:
		// TODO 如果提交通过
		if flag {
			record.Condition = "Accepted"
			// TODO 检查是否是今日首次通过
			if g.DB.Where("user_id = ? and condition = Accepted and to_days(created_at) = to_days(now())", record.UserId).First(&model.Record{}).Error != nil {
				Handle.Behaviors["Days"].PublishBehavior(1, record.UserId)
			}
			// TODO 检查该题目是否是首次通过
			if g.DB.Where("user_id = ? and condition = Accepted and problem_id = ?", record.UserId, record.ProblemId).First(&model.Record{}).Error != nil {
				Handle.Behaviors["Accepts"].PublishBehavior(1, record.UserId)
				categoryMap := util.ProblemCategory(problem.ID)
				for category := range categoryMap {
					Handle.Behaviors[category].PublishBehavior(1, record.UserId)
				}
			}
		}
		// TODO 查看是否为比赛提交,且比赛已经开始
		if record.CreatedAt.After(competition.StartTime) && record.CreatedAt.Before(competition.EndTime) {
			groups, _ := g.Redis.ZRange(g.ctx, "CompetitionR"+competition.ID.String(), 0, -1).Result()
			// TODO 查找组
			var group model.Group
			for i := range groups {

				// TODO 先看redis中是否存在
				if ok, _ := g.Redis.HExists(g.ctx, "Group", groups[i]).Result(); ok {
					cate, _ := g.Redis.HGet(g.ctx, "Group", groups[i]).Result()
					if json.Unmarshal([]byte(cate), &group) == nil {
						goto levp
					} else {
						// TODO 移除损坏数据
						g.Redis.HDel(g.ctx, "Group", groups[i])
					}
				}

				// TODO 查看用户组是否在数据库中存在
				g.DB.Where("id = (?)", groups[i]).First(&group)
				{
					// TODO 将用户组存入redis供下次使用
					v, _ := json.Marshal(group)
					g.Redis.HSet(g.ctx, "Group", groups[i], v)
				}
			levp:
				if g.DB.Where("group_id = (?) and user_id = (?)", group.ID, record.UserId).First(&model.UserList{}).Error == nil {
					break
				}
			}
			var competitionMembers []model.CompetitionMember
			// TODO 在redis中取出成员罚时具体数据
			cM, err := g.Redis.HGet(g.ctx, "Competition"+competition.ID.String(), group.ID.String()).Result()
			if err == nil {
				json.Unmarshal([]byte(cM), &competitionMembers)
			}
			// TODO 找出数组中对应的题目
			k := -1
			for i := range competitionMembers {
				if competitionMembers[i].ProblemId == record.ProblemId {
					k = i
					break
				}
			}
			if k == -1 {
				k = len(competitionMembers)
				competitionMembers = append(competitionMembers, model.CompetitionMember{
					ID:            uuid.NewV4(),
					MemberId:      group.ID,
					CompetitionId: competition.ID,
					ProblemId:     record.ProblemId,
					Pass:          0,
					Penalties:     0,
				})
			}
			// TODO 在redis中取出通过、罚时情况
			cR, err := g.Redis.ZScore(g.ctx, "CompetitionR"+competition.ID.String(), group.ID.String()).Result()
			if err != nil {
				cR = 0
			}
			// TODO 先前没有通过
			if competitionMembers[k].Condition != "Accepted" {
				// TODO 记录罚时
				competitionMembers[k].Penalties += time.Time(record.CreatedAt).Sub(time.Time(competition.StartTime))
				// TODO 通过样例数量增加
				if competitionMembers[k].Pass < record.Pass {
					competitionMembers[k].Condition = record.Condition
					// TODO 获取用例通过分数
					score := 0
					for i := competitionMembers[k].Pass + 1; i <= record.Pass; i++ {
						score += int(cases[i].Score)
					}
					// TODO 如果分数上升
					if score > 0 {
						// TODO 记入罚时
						cR -= float64(competitionMembers[k].Penalties) / 10000000000
						cR += float64(score)
						// TODO 存入redis供下次使用
						v, _ := json.Marshal(competitionMembers)
						g.Redis.HSet(g.ctx, "Competition"+competition.ID.String(), group.ID.String(), v)
						g.Redis.ZAdd(g.ctx, "CompetitionR"+competition.ID.String(), redis.Z{Score: cR, Member: group.ID.String()})
						// TODO 发布订阅用于滚榜
						rankList := vo.RankList{
							MemberId: group.ID,
						}
						// TODO 将ranklist打包
						v, _ = json.Marshal(rankList)
						g.Redis.Publish(g.ctx, "CompetitionChan"+competition.ID.String(), v)
					}
				}
			}
		}
	} else {
		record.Condition = "Language Error"
	}
}

// @title    NewGroup
// @description   新建一Group
// @auth      MGAronya      2022-9-16 12:23
// @param    void
// @return   IJude		返回一个IGroup用于调用各种函数
func NewGroup(rw *sync.RWMutex) IGroup {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	ctx := context.Background()
	Group := Group{rw: rw, DB: db, Redis: redis, ctx: ctx}
	Interface.ComsumerMap["Group"] = &Group
	Interface.ComsumerMap["Match"] = &Group
	return Group
}
