package rabbitMq

import (
	"MGA_OJ/Handle"
	"MGA_OJ/Interface"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var LanguageMap map[string]Interface.CmdInterface = map[string]Interface.CmdInterface{
	"C++":   Handle.NewCppPlusPlus(),
	"C++11": Handle.NewCppPlusPlus11(),
	"Java":  Handle.NewJava(),
}

var db *gorm.DB = common.GetDB()
var ctx context.Context = context.Background()

var redis *redis.Client = common.GetRedisClient(0)

func Test(msg []byte) {
	var record model.Record
	// TODO 先看redis中是否存在
	if ok, _ := redis.HExists(ctx, "Record", string(msg)).Result(); ok {
		cate, _ := redis.HGet(ctx, "Record", string(msg)).Result()
		// TODO 移除损坏数据
		redis.HDel(ctx, "Record", string(msg))
		if json.Unmarshal([]byte(cate), &record) == nil {
			// TODO 跳过数据库搜寻过程
			goto leep
		}
	}

	// TODO 未能找到提交记录
	if db.Where("id = ?", msg).First(&record).Error != nil {
		log.Printf("%s Record Disappear!!\n", msg)
		return
	}

leep:

	// TODO 找到提交记录后，开始判题逻辑
	if v, ok := LanguageMap[record.Language]; ok {
		handle(record, v)
	} else {
		record.Condition = "Luanguage Error"
		db.Save(&record)
	}
	// TODO 将提交存入redis供下次使用
	v, _ := json.Marshal(record)
	redis.HSet(ctx, "Record", record.ID, v)
}

func handle(record model.Record, cmdI Interface.CmdInterface) {
	// TODO 从数据库中读出输入输出
	var testInputs []model.TestInput
	var testOutputs []model.TestOutput
	var problem model.Problem
	// TODO 先看redis中是否存在
	id := fmt.Sprint(record.ProblemId)
	if ok, _ := redis.HExists(ctx, "Problem", id).Result(); ok {
		cate, _ := redis.HGet(ctx, "Problem", id).Result()
		if json.Unmarshal([]byte(cate), &problem) == nil {
			// TODO 跳过数据库搜寻problem过程
			goto leep
		} else {
			// TODO 移除损坏数据
			redis.HDel(ctx, "Problem", id)
		}
	}

	// TODO 查看题目是否在数据库中存在
	if db.Where("id = ?", id).First(&problem).Error != nil {
		record.Condition = "Problem Doesn't Exist"
		db.Save(&record)
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(problem)
		redis.HSet(ctx, "Problem", id, v)
	}

leep:
	// TODO 查看比赛是否开始
	var competition model.Competition

	// TODO 是否存在比赛
	if problem.CompetitionId != (uuid.UUID{}) {
		goto leap
	}

	// TODO 先看redis中是否存在
	if ok, _ := redis.HExists(ctx, "Competition", fmt.Sprint(problem.CompetitionId)).Result(); ok {
		cate, _ := redis.HGet(ctx, "Competition", fmt.Sprint(problem.CompetitionId)).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			goto leap
		} else {
			// TODO 移除损坏数据
			redis.HDel(ctx, "Competition", fmt.Sprint(problem.CompetitionId))
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if db.Where("id = ?", problem.CompetitionId).First(&competition).Error != nil {
		goto leap
	}

	// TODO 将竞赛存入redis供下次使用
	{
		v, _ := json.Marshal(competition)
		redis.HSet(ctx, "Competition", fmt.Sprint(problem.CompetitionId), v)
	}

	// TODO 如果比赛未开始
	if competition.StartTime.After(record.CreatedAt) {
		record.Condition = "Competition hasn't Started"
		db.Save(&record)
		return
	}
leap:

	// TODO 查找输入
	if ok, _ := redis.HExists(ctx, "Input", id).Result(); ok {
		cate, _ := redis.HGet(ctx, "Input", id).Result()
		if json.Unmarshal([]byte(cate), &testInputs) == nil {
			// TODO 跳过数据库搜寻testInputs过程
			goto Input
		} else {
			// TODO 移除损坏数据
			redis.HDel(ctx, "Input", id)
		}
	}

	// TODO 查看题目是否在数据库中存在
	if db.Where("id = ?", id).Find(&testInputs).Error != nil {
		record.Condition = "Input Doesn't Exist"
		db.Save(&record)
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(testInputs)
		redis.HSet(ctx, "Input", id, v)
	}
Input:
	// TODO 查找输入
	if ok, _ := redis.HExists(ctx, "Output", id).Result(); ok {
		cate, _ := redis.HGet(ctx, "Output", id).Result()
		if json.Unmarshal([]byte(cate), &testOutputs) == nil {
			// TODO 跳过数据库搜寻testOutputs过程
			goto Output
		} else {
			// TODO 移除损坏数据
			redis.HDel(ctx, "Output", id)
		}
	}

	// TODO 查看题目是否在数据库中存在
	if db.Where("id = ?", id).Find(&testOutputs).Error != nil {
		record.Condition = "Output Doesn't Exist"
		db.Save(&record)
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(testOutputs)
		redis.HSet(ctx, "Output", id, v)
	}
Output:
	fp, err := os.Create("user-code/" + fmt.Sprint(record.ID) + "." + cmdI.Suffix())
	// TODO 文件错误
	if err != nil {
		// TODO 创建文件失败的原因有：
		// TODO 1、路径不存在  2、权限不足  3、打开文件数量超过上限  4、磁盘空间不足等
		record.Condition = "System Error 1"
		db.Save(&record)
		return
	}

	// TODO defer延迟调用 关闭文件，释放资源
	defer fp.Close()

	// TODO 写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(fp)

	write.WriteString(record.Code)

	// TODO Flush将缓存的文件真正写入到文件中
	write.Flush()

	// TODO 编译
	cmd := cmdI.Compile("user-code/" + fmt.Sprint(record.ID))

	// TODO 系统错误
	if err := cmd.Start(); err != nil {
		record.Condition = "System Error 2"
		db.Save(&record)
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
		record.Condition = "Compile timeout"
		db.Save(&record)
		return
	case err = <-done:
	}

	// TODO 编译出错
	if err != nil {
		record.Condition = "Compile Error"
		db.Save(&record)
		return
	}

	for i := 0; i < len(testInputs); i++ {
		var bm runtime.MemStats
		runtime.ReadMemStats(&bm)

		// TODO 运行可执行文件
		cmd = cmdI.Run("user-code/" + fmt.Sprint(record.ID))

		var out, stderr bytes.Buffer
		cmd.Stderr = &stderr
		cmd.Stdout = &out
		stdinPipe, err := cmd.StdinPipe()
		// TODO 系统错误
		if err != nil {
			record.Condition = "System Error 3"
			db.Save(&record)
			return
		}
		io.WriteString(stdinPipe, testInputs[i].Input)
		now := time.Now().UnixMilli()
		// TODO 系统错误
		if err := cmd.Start(); err != nil {
			record.Condition = "System Error 4"
			db.Save(&record)
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
		cas := model.Case{
			RecordId: record.ID,
			Time:     uint(end - now - int64(cmdI.RunUpTime())),
			Memory:   uint(em.Alloc / 1024),
		}
		// TODO 超时
		if cas.Time > problem.TimeLimit*cmdI.TimeMultiplier() {
			record.Condition = "Time Limit Exceeded"
			flag = false
			goto final
		}
		// TODO 超出内存限制
		if cas.Memory > problem.MemoryLimit {
			record.Condition = "Memory Limit Exceeded"
			flag = false
			goto final
		}
		// TODO 答案错误
		if out.String() != testOutputs[i].Output {
			record.Condition = "Wrong Answer"
			flag = false
			goto final
		}
		// TODO 通过数量+1
		record.Pass++

		// TODO 数据库插入数据错误
		if db.Create(&cas).Error != nil {
			record.Condition = "System rror 5"
			db.Save(&record)
			return
		}
	}
final:
	// TODO 如果提交通过
	if flag {
		record.Condition = "Accepted"
	}
	// TODO 查看是否为比赛提交,且比赛已经开始
	if competition.ID != (uuid.UUID{}) && time.Now().After(time.Time(competition.StartTime)) {
		var TID uuid.UUID
		if competition.Type == "Single" {
			TID = record.UserId
		} else {
			// TODO 查看该用户属于哪个组
			var competition model.Competition
			db.Where("id = ?", record.CompetitionId).First(&competition)
			// TODO 如果没有参加比赛
			if db.Table("user_lists").Select("user_lists.group_id").Joins("left join group_lists on user_lists.group_id = group_lists.group_id").Where("user_id = ? and set_id = ?", record.UserId, competition.SetId).Scan(&TID).Error != nil {
				record.Condition = "Absent from the race"
				db.Save(&record)
				return
			}
		}
		var competitionMembers []model.CompetitionMember
		// TODO 在redis中取出成员罚时具体数据
		cM, err := redis.HGet(ctx, "competition"+record.CompetitionId.String(), TID.String()).Result()
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
		cR, err := redis.ZScore(ctx, "Competition"+record.CompetitionId.String(), TID.String()).Result()
		if err != nil {
			cR = 0
		}
		// TODO 如果先前没有通过该题，则记录罚时
		if competitionMembers[k].Condition != "Accepted" {
			competitionMembers[k].Penalties += time.Now().Sub(time.Time(competition.StartTime))
			cR -= float64(time.Now().Sub(time.Time(competition.StartTime))) / 10000000000
			competitionMembers[k].Condition = record.Condition
			if flag {
				cR++
			}
			// TODO 存入redis供下次使用
			v, _ := json.Marshal(competitionMembers)
			redis.HSet(ctx, "competition"+record.CompetitionId.String(), TID.String(), v)
			redis.ZAdd(ctx, "competition"+record.CompetitionId.String(), redis.Z{Score: cR, Member: TID.String()})
		}
	}
	db.Save(&record)
	return
}
