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
	"log"
	"math"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/go-redis/redis/v9"
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
func (j Judge) Handel(msg string) {
	// TODO 单核处理，上锁
	j.rw.Lock()
	// TODO 确保资源归还
	defer j.rw.Unlock()
	var record model.Record
	// TODO 先看redis中是否存在
	if ok, _ := j.Redis.HExists(j.ctx, "Record", msg).Result(); ok {
		cate, _ := j.Redis.HGet(j.ctx, "Record", msg).Result()
		// TODO 移除损坏数据
		j.Redis.HDel(j.ctx, "Record", msg)
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
	// TODO 发布订阅用于提交列表
	recordList := vo.RecordList{
		RecordId: record.ID,
	}
	// TODO 确保信息进入频道
	defer func() {
		// TODO 将recordlist打包
		v, _ := json.Marshal(recordList)
		j.Redis.Publish(j.ctx, "RecordChan", v)
		// TODO 将record存入redis
		v, _ = json.Marshal(record)
		j.Redis.HSet(j.ctx, "Record", fmt.Sprint(record.ID), v)
		// TODO 将record存入mysql
		j.DB.Save(&record)
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
		j.Redis.Publish(j.ctx, "RecordChan", v)
		// TODO 将record存入redis
		v, _ = json.Marshal(record)
		j.Redis.HSet(j.ctx, "Record", fmt.Sprint(record.ID), v)
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
			return
		}
		// TODO 将题目存入redis供下次使用
		{
			v, _ := json.Marshal(problem)
			j.Redis.HSet(j.ctx, "Problem", id, v)
		}

	leep:

		// TODO 查找用例
		if ok, _ := j.Redis.HExists(j.ctx, "Case", id).Result(); ok {
			cate, _ := j.Redis.HGet(j.ctx, "Case", id).Result()
			if json.Unmarshal([]byte(cate), &cases) == nil {
				// TODO 跳过数据库搜寻testInputs过程
				goto Case
			} else {
				// TODO 移除损坏数据
				j.Redis.HDel(j.ctx, "Case", id)
			}
		}

		// TODO 查看题目是否在数据库中存在
		if j.DB.Where("problem_id = ?", id).Find(&cases).Error != nil {
			record.Condition = "Input Doesn't Exist"
			return
		}
		// TODO 将题目存入redis供下次使用
		{
			v, _ := json.Marshal(cases)
			j.Redis.HSet(j.ctx, "Case", id, v)
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
			j.Redis.Publish(j.ctx, "RecordChan", v)
			// TODO 将record存入redis
			v, _ = json.Marshal(record)
			j.Redis.HSet(j.ctx, "Record", fmt.Sprint(record.ID), v)
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
			record.Condition = "Compile timeout"
			return
		case err = <-done:
		}

		// TODO 编译出错
		if err != nil {
			record.Condition = "Compile Error"
			return
		}

		// TODO 获取权限
		cmd = cmdI.Chmod("./user-code/", fileId)

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
			record.Condition = "Compile timeout"
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
			j.Redis.Publish(j.ctx, "RecordChan", v)
			// TODO 将record存入redis
			v, _ = json.Marshal(record)
			j.Redis.HSet(j.ctx, "Record", fmt.Sprint(record.ID), v)
		}

		for i := 0; i < len(cases); i++ {
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
			cas := model.CaseCondition{
				RecordId: record.ID,
				Time:     uint(math.Max(float64(end-now-int64(cmdI.RunUpTime())), 0)),
				Memory:   uint(em.Alloc/1024 - bm.Alloc/1024),
				Input:    cases[i].Input,
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
			var specalJudge model.Program
			// TODO 查看题目是否有标准程序

			// TODO 先看redis中是否存在
			if ok, _ := j.Redis.HExists(j.ctx, "Program", problem.SpecialJudge.String()).Result(); ok {
				cate, _ := j.Redis.HGet(j.ctx, "Program", problem.SpecialJudge.String()).Result()
				if json.Unmarshal([]byte(cate), &specalJudge) == nil {
					// TODO 跳过数据库搜寻program过程
					goto special
				} else {
					// TODO 移除损坏数据
					j.Redis.HDel(j.ctx, "Program", problem.SpecialJudge.String())
				}
			}

			// TODO 查看程序是否在数据库中存在
			if j.DB.Where("id = ?", problem.SpecialJudge.String()).First(&specalJudge).Error != nil {
				goto outPut
			}
			// TODO 将题目存入redis供下次使用
			{
				v, _ := json.Marshal(specalJudge)
				j.Redis.HSet(j.ctx, "Program", problem.SpecialJudge.String(), v)
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
				record.Condition = "Wrong Answer"
				flag = false
				goto final
			}
		pass:
			// TODO 通过数量+1
			record.Pass++

			// TODO 数据库插入数据错误
			if j.DB.Create(&cas).Error != nil {
				record.Condition = "System error 5"
				return
			}
		}
	final:
		// TODO 如果提交通过
		if flag {
			record.Condition = "Accepted"
		}
	} else {
		record.Condition = "Language Error"
	}
}

// @title    NewJudge
// @description   新建一Judge
// @auth      MGAronya（张健）      2022-9-16 12:23
// @param    void
// @return   IJude		返回一个IJudge用于调用各种函数
func NewJudge(rw *sync.RWMutex) IJudge {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	ctx := context.Background()
	Normal := Judge{rw: rw, DB: db, Redis: redis, ctx: ctx}
	Interface.ComsumerMap["Normal"] = &Normal
	return Normal
}
