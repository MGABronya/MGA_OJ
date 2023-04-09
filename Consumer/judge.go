package consumer

import (
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
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/go-redis/redis/v9"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// IJudge			定义了判题类接口
type IJudge interface {
	Interface.ConsumerInterface // 包含消费功能
}

// Judge			定义了判断工具类
type Judge struct {
	rw    *sync.RWMutex // 含有锁
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
	ctx   context.Context
}

// @title    Handle
// @description   创建一篇判断
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (j Judge) Handel(msg []byte) {
	// TODO 单核处理，上锁
	j.rw.Lock()
	// TODO 确保资源归还
	defer j.rw.Unlock()
	var record model.Record
	// TODO 先看redis中是否存在
	if ok, _ := j.Redis.HExists(j.ctx, "Record", string(msg)).Result(); ok {
		cate, _ := j.Redis.HGet(j.ctx, "Record", string(msg)).Result()
		// TODO 移除损坏数据
		j.Redis.HDel(j.ctx, "Record", string(msg))
		if json.Unmarshal([]byte(cate), &record) == nil {
			// TODO 跳过数据库搜寻过程
			goto feep
		}
	}

	// TODO 未能找到提交记录
	if j.DB.Where("id = ?", msg).First(&record).Error != nil {
		log.Printf("%s Record Disappear!!\n", msg)
		return
	}

feep:
	// TODO 查看代码是否为空
	if record.Code == "" {
		record.Condition = "Code is empty"
		j.DB.Save(&record)
		return
	}
	// TODO 找到提交记录后，开始判题逻辑
	if cmdI, ok := util.LanguageMap[record.Language]; ok {
		// TODO 从数据库中读出输入输出
		var testInputs []model.TestInput
		var testOutputs []model.TestOutput
		var problem model.Problem
		// TODO 先看redis中是否存在
		id := fmt.Sprint(record.ProblemId)
		if ok, _ := j.Redis.HExists(j.ctx, "Problem", id).Result(); ok {
			cate, _ := j.Redis.HGet(j.ctx, "Problem", id).Result()
			if json.Unmarshal([]byte(cate), &problem) == nil {
				// TODO 跳过数据库搜寻problem过程
				goto leep
			} else {
				// TODO 移除损坏数据
				j.Redis.HDel(j.ctx, "Problem", id)
			}
		}

		// TODO 查看题目是否在数据库中存在
		if j.DB.Where("id = ?", id).First(&problem).Error != nil {
			record.Condition = "Problem Doesn't Exist"
			j.DB.Save(&record)
			return
		}
		// TODO 将题目存入redis供下次使用
		{
			v, _ := json.Marshal(problem)
			j.Redis.HSet(j.ctx, "Problem", id, v)
		}

	leep:
		// TODO 查看比赛是否开始
		var competition model.Competition

		// TODO 是否存在比赛
		if problem.CompetitionId != (uuid.UUID{}) {
			goto leap
		}

		// TODO 先看redis中是否存在
		if ok, _ := j.Redis.HExists(j.ctx, "Competition", fmt.Sprint(problem.CompetitionId)).Result(); ok {
			cate, _ := j.Redis.HGet(j.ctx, "Competition", fmt.Sprint(problem.CompetitionId)).Result()
			if json.Unmarshal([]byte(cate), &competition) == nil {
				goto leap
			} else {
				// TODO 移除损坏数据
				j.Redis.HDel(j.ctx, "Competition", fmt.Sprint(problem.CompetitionId))
			}
		}

		// TODO 查看比赛是否在数据库中存在
		if j.DB.Where("id = ?", problem.CompetitionId).First(&competition).Error != nil {
			goto leap
		}

		// TODO 将竞赛存入redis供下次使用
		{
			v, _ := json.Marshal(competition)
			j.Redis.HSet(j.ctx, "Competition", fmt.Sprint(problem.CompetitionId), v)
		}

		// TODO 如果比赛未开始
		if competition.StartTime.After(record.CreatedAt) {
			record.Condition = "Competition hasn't Started"
			j.DB.Save(&record)
			return
		}
	leap:

		// TODO 查找输入
		if ok, _ := j.Redis.HExists(j.ctx, "Input", id).Result(); ok {
			cate, _ := j.Redis.HGet(j.ctx, "Input", id).Result()
			if json.Unmarshal([]byte(cate), &testInputs) == nil {
				// TODO 跳过数据库搜寻testInputs过程
				goto Input
			} else {
				// TODO 移除损坏数据
				j.Redis.HDel(j.ctx, "Input", id)
			}
		}

		// TODO 查看题目是否在数据库中存在
		if j.DB.Where("id = ?", id).Find(&testInputs).Error != nil {
			record.Condition = "Input Doesn't Exist"
			j.DB.Save(&record)
			return
		}
		// TODO 将题目存入redis供下次使用
		{
			v, _ := json.Marshal(testInputs)
			j.Redis.HSet(j.ctx, "Input", id, v)
		}
	Input:
		// TODO 查找输入
		if ok, _ := j.Redis.HExists(j.ctx, "Output", id).Result(); ok {
			cate, _ := j.Redis.HGet(j.ctx, "Output", id).Result()
			if json.Unmarshal([]byte(cate), &testOutputs) == nil {
				// TODO 跳过数据库搜寻testOutputs过程
				goto Output
			} else {
				// TODO 移除损坏数据
				j.Redis.HDel(j.ctx, "Output", id)
			}
		}

		// TODO 查看题目是否在数据库中存在
		if j.DB.Where("id = ?", id).Find(&testOutputs).Error != nil {
			record.Condition = "Output Doesn't Exist"
			j.DB.Save(&record)
			return
		}
		// TODO 将题目存入redis供下次使用
		{
			v, _ := json.Marshal(testOutputs)
			j.Redis.HSet(j.ctx, "Output", id, v)
		}
	Output:
		fileId := cmdI.Name()
		fp, err := os.Create("user-code/" + fileId + "." + cmdI.Suffix())
		// TODO 文件错误
		if err != nil {
			// TODO 创建文件失败的原因有：
			// TODO 1、路径不存在  2、权限不足  3、打开文件数量超过上限  4、磁盘空间不足等
			record.Condition = "System Error 1"
			j.DB.Save(&record)
			goto exit
		}

		// TODO defer延迟调用 关闭文件，释放资源
		defer fp.Close()

		// TODO 写入文件时，使用带缓存的 *Writer
		write := bufio.NewWriter(fp)

		write.WriteString(record.Code)

		// TODO Flush将缓存的文件真正写入到文件中
		write.Flush()

		// TODO 编译
		cmd := cmdI.Compile("user-code/", fileId)

		// TODO 系统错误
		if err := cmd.Start(); err != nil {
			record.Condition = "System Error 2"
			j.DB.Save(&record)
			goto exit
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
			record.Condition = "Compile timeout"
			j.DB.Save(&record)
			goto exit
		case err = <-done:
		}

		// TODO 编译出错
		if err != nil {
			record.Condition = "Compile Error"
			j.DB.Save(&record)
			goto exit
		}

		for i := 0; i < len(testInputs); i++ {
			var bm runtime.MemStats
			runtime.ReadMemStats(&bm)

			// TODO 运行可执行文件
			cmd = cmdI.Run("./user-code/", fileId)

			var out, stderr bytes.Buffer
			cmd.Stderr = &stderr
			cmd.Stdout = &out
			stdinPipe, err := cmd.StdinPipe()
			// TODO 系统错误
			if err != nil {
				record.Condition = "System Error 3"
				j.DB.Save(&record)
				goto exit
			}
			io.WriteString(stdinPipe, testInputs[i].Input+"\n")
			// TODO 关闭管道制造EOF信息
			stdinPipe.Close()
			now := time.Now().UnixMilli()
			// TODO 系统错误
			if err := cmd.Start(); err != nil {
				record.Condition = "System Error 4"
				j.DB.Save(&record)
				goto exit
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
			cas := model.Case{
				RecordId: record.ID,
				Time:     uint(math.Max(float64(end-now-int64(cmdI.RunUpTime())), 0)),
				Memory:   uint(em.Alloc/1024 - bm.Alloc/1024),
			}
			// TODO 超时
			if cas.Time > problem.TimeLimit*cmdI.TimeMultiplier() {
				record.Condition = "Time Limit Exceeded"
				flag = false
				goto final
			}
			// TODO 超出内存限制
			if cas.Memory > problem.MemoryLimit*cmdI.MemoryMultiplier() {
				record.Condition = "Memory Limit Exceeded"
				flag = false
				goto final
			}
			// TODO 答案错误
			var specalJudge model.SpecialJudge
			// TODO 查看特判是否存在
			if j.DB.Where("id = ?", problem.SpecialJudge).First(&specalJudge).Error != nil {
				// TODO 进行特判
				res := TQ.JudgeRun(specalJudge.Language, specalJudge.Code, testInputs[i].Input+"\n"+out.String(), problem.MemoryLimit*5, problem.TimeLimit*5)
				if res != "ok" {
					record.Condition = res
					flag = false
					goto final
				}
			} else if out.String() != testOutputs[i].Output {
				record.Condition = "Wrong Answer"
				flag = false
				goto final
			}
			// TODO 通过数量+1
			record.Pass++

			// TODO 数据库插入数据错误
			if j.DB.Create(&cas).Error != nil {
				record.Condition = "System error 5"
				j.DB.Save(&record)
				return
			}
		}
	final:
		// TODO 如果提交通过
		if flag {
			record.Condition = "Accepted"
		}
		// TODO 查看是否为比赛提交,且比赛已经开始
		if competition.ID != (uuid.UUID{}) && record.CreatedAt.After(competition.StartTime) {
			var TID uuid.UUID
			if competition.Type == "Single" {
				TID = record.UserId
			} else {
				// TODO 查看该用户属于哪个组
				var competition model.Competition
				j.DB.Where("id = ?", record.CompetitionId).First(&competition)
				// TODO 如果没有参加比赛
				if j.DB.Table("user_lists").Select("user_lists.group_id").Joins("left join group_lists on user_lists.group_id = group_lists.group_id").Where("user_id = ? and set_id = ?", record.UserId, competition.SetId).Scan(&TID).Error != nil {
					record.Condition = "Absent from the race"
					j.DB.Save(&record)
					goto exit
				}
			}
			var competitionMembers []model.CompetitionMember
			// TODO 在redis中取出成员罚时具体数据
			cM, err := j.Redis.HGet(j.ctx, "competition"+record.CompetitionId.String(), TID.String()).Result()
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
					MemberId:      TID,
					CompetitionId: record.CompetitionId,
					ProblemId:     record.ProblemId,
				})
			}
			// TODO 在redis中取出通过、罚时情况
			cR, err := j.Redis.ZScore(j.ctx, "Competition"+record.CompetitionId.String(), TID.String()).Result()
			if err != nil {
				cR = 0
			}
			// TODO 如果先前没有通过该题，则记录罚时
			if competitionMembers[k].Condition != "Accepted" {
				competitionMembers[k].Penalties += time.Time(record.CreatedAt).Sub(time.Time(competition.StartTime))
				cR -= float64(time.Time(record.CreatedAt).Sub(time.Time(competition.StartTime))) / 10000000000
				competitionMembers[k].Condition = record.Condition
				// TODO 如果此次为第一次通过
				if flag {
					cR++
					// TODO 发布订阅用于滚榜
					rankList := vo.RankList{
						MemberId: TID,
					}
					// TODO 将ranklist打包
					v, _ := json.Marshal(rankList)
					j.Redis.Publish(j.ctx, "CompetitionChan"+competition.ID.String(), v)
				}
				// TODO 存入redis供下次使用
				v, _ := json.Marshal(competitionMembers)
				j.Redis.HSet(j.ctx, "competition"+record.CompetitionId.String(), TID.String(), v)
				j.Redis.ZAdd(j.ctx, "Competition"+record.CompetitionId.String(), redis.Z{Score: cR, Member: TID.String()})
			}
		}
		j.DB.Save(&record)
	} else {
		record.Condition = "Language Error"
		j.DB.Save(&record)
	}
exit:

	// TODO 将提交存入redis供下次使用
	v, _ := json.Marshal(record)
	j.Redis.HSet(j.ctx, "Record", record.ID, v)
}

// @title    NewJudge
// @description   新建一个Judge
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   IJudge		返回一个IJudge用于调用各种函数
func NewJudge() IJudge {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	ctx := context.Background()
	return Judge{rw: &sync.RWMutex{}, DB: db, Redis: redis, ctx: ctx}
}
