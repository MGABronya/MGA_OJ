// @Title  util
// @Description  收集各种需要使用的工具函数
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:47
package util

import (
	"MGA_OJ/Interface"
	Handle "MGA_OJ/Language"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/smtp"
	"os"
	"path"
	"regexp"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Units		定义了单位换算
var Units = map[string]uint{
	"mb": 1024,
	"kb": 1,
	"gb": 1024 * 1024,
	"ms": 1,
	"s":  1000,
}

var Max_run int = 4

// timerMap	    定义了当前使用的定时器
var TimerMap map[uuid.UUID]*time.Timer = make(map[uuid.UUID]*time.Timer)

// LanguageMap			定义语言字典，对应其处理方式
var LanguageMap map[string]Interface.CmdInterface = map[string]Interface.CmdInterface{
	"C++":   Handle.NewCppPlusPlus(),
	"C++11": Handle.NewCppPlusPlus11(),
	"Java":  Handle.NewJava(),
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
	// TODO 找到提交记录后，开始判题逻辑
	if cmdI, ok := LanguageMap[record.Language]; ok {
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
		// TODO 清空当前文件夹
		dir, err := ioutil.ReadDir("/user-code/")
		for _, d := range dir {
			os.RemoveAll(path.Join([]string{"user-code/", d.Name()}...))
		}
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
			io.WriteString(stdinPipe, testInputs[i].Input)
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
			if j.DB.Create(&cas).Error != nil {
				record.Condition = "System rror 5"
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
		if competition.ID != (uuid.UUID{}) && time.Now().After(time.Time(competition.StartTime)) {
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
				competitionMembers[k].Penalties += time.Now().Sub(time.Time(competition.StartTime))
				cR -= float64(time.Now().Sub(time.Time(competition.StartTime))) / 10000000000
				competitionMembers[k].Condition = record.Condition
				if flag {
					cR++
				}
				// TODO 存入redis供下次使用
				v, _ := json.Marshal(competitionMembers)
				j.Redis.HSet(j.ctx, "competition"+record.CompetitionId.String(), TID.String(), v)
				j.Redis.ZAdd(j.ctx, "competition"+record.CompetitionId.String(), redis.Z{Score: cR, Member: TID.String()})
			}
		}
		j.DB.Save(&record)
	} else {
		record.Condition = "Luanguage Error"
		j.DB.Save(&record)
	}
exit:

	// TODO 将提交存入redis供下次使用
	v, _ := json.Marshal(record)
	j.Redis.HSet(j.ctx, "Record", record.ID, v)
}

// @title    NewRecordController
// @description   新建一个IRecordController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   IRecordController		返回一个IRecordController用于调用各种函数
func NewJudge() Judge {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	ctx := context.Background()
	return Judge{rw: &sync.RWMutex{}, DB: db, Redis: redis, ctx: ctx}
}

// @title    RandomString
// @description   生成一段随机的字符串
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     n int		字符串的长度
// @return    string    一串随机的字符串
func RandomString(n int) string {
	var letters = []byte("qwertyuioplkjhgfdsazxcvbnmQWERTYUIOPLKJHGFDSAZXCVBNM")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	// TODO 不断用随机字母填充字符串
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// @title    VerifyEmailFormat
// @description   用于验证邮箱格式是否正确的工具函数
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     email string		一串字符串，表示邮箱
// @return    bool    返回是否合法
func VerifyEmailFormat(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// @title    VerifyMobileFormat
// @description   用于验证手机号格式是否正确的工具函数
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     mobileNum string		一串字符串，表示手机号
// @return    bool    返回是否合法
func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

// @title    VerifyQQFormat
// @description   用于验证QQ号格式是否正确的工具函数
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     QQNum string		一串字符串，表示QQ
// @return    bool    返回是否合法
func VerifyQQFormat(QQNum string) bool {
	regular := "[1-9][0-9]{4,10}"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(QQNum)
}

// @title    VerifyQQFormat
// @description  用于验证Icon是否为默认图片的工具函数
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     Icon string		一串字符串，表示图像名称
// @return    bool    返回是否合法
func VerifyIconFormat(Icon string) bool {
	regular := "MGA[1-9].jpg"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(Icon)
}

// @title    isEmailExist
// @description   查看email是否在数据库中存在
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func IsEmailExist(db *gorm.DB, email string) bool {
	var user model.User
	db.Where("email = ?", email).First(&user)
	return user.ID != uuid.UUID{}
}

// @title    isNameExist
// @description   查看name是否在数据库中存在
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func IsNameExist(db *gorm.DB, name string) bool {
	var user model.User
	db.Where("name = ?", name).First(&user)
	return user.ID != uuid.UUID{}
}

// @title    SendEmailValidate
// @description   发送验证邮件
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    em []string       接收一个邮箱字符串
// @return   string, error     返回验证码和error值
func SendEmailValidate(em []string) (string, error) {
	mod := `
	尊敬的%s，您好！

	您于 %s 提交的邮箱验证，本次验证码为%s，为了保证账号安全，验证码有效期为5分钟。请确认为本人操作，切勿向他人泄露，感谢您的理解与使用。
	此邮箱为系统邮箱，请勿回复。
`
	e := email.NewEmail()
	e.From = "mgAronya <2829214609@qq.com>"
	e.To = em
	// TODO 生成6位随机验证码
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vCode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	t := time.Now().Format("2006-01-02 15:04:05")
	// TODO 设置文件发送的内容
	content := fmt.Sprintf(mod, em[0], t, vCode)
	e.Text = []byte(content)
	// TODO 设置服务器相关的配置
	err := e.Send("smtp.qq.com:25", smtp.PlainAuth("", "2829214609@qq.com", "rmdtxokuuqyrdgii", "smtp.qq.com"))
	return vCode, err
}

// @title    SendEmailPass
// @description   发送密码邮件
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    em []string       接收一个邮箱字符串
// @return   string, error     返回验证码和error值
func SendEmailPass(em []string) string {
	mod := `
	尊敬的%s，您好！

	您于 %s 提交的邮箱验证，已经将密码重置为%s，为了保证账号安全。切勿向他人泄露，并尽快更改密码，感谢您的理解与使用。
	此邮箱为系统邮箱，请勿回复。
`
	e := email.NewEmail()
	e.From = "mgAronya <2829214609@qq.com>"
	e.To = em
	// TODO 生成8位随机密码
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	password := fmt.Sprintf("%08v", rnd.Int31n(100000000))
	t := time.Now().Format("2006-01-02 15:04:05")

	db := common.GetDB()

	// TODO 创建密码哈希
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "密码加密失败"
	}

	// TODO 更新密码
	err = db.Model(&model.User{}).Where("email = ?", em[0]).Updates(model.User{
		Password: string(hasedPassword),
	}).Error

	if err != nil {
		return "密码更新失败"
	}

	// TODO 设置文件发送的内容
	content := fmt.Sprintf(mod, em[0], t, password)
	e.Text = []byte(content)
	// TODO 设置服务器相关的配置
	err = e.Send("smtp.qq.com:25", smtp.PlainAuth("", "2829214609@qq.com", "rmdtxokuuqyrdgii", "smtp.qq.com"))

	if err != nil {
		return "邮件发送失败"
	}

	return "密码已重置"
}

// @title    IsEmailPass
// @description   验证邮箱是否通过
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    em []string       接收一个邮箱字符串
// @return   string, error     返回验证码和error值
func IsEmailPass(ctx *gin.Context, email string, vertify string) bool {
	client := common.GetRedisClient(0)
	V, err := client.Get(ctx, email).Result()
	if err != nil {
		return false
	}
	return V == vertify
}

// @title    SetRedisEmail
// @description   设置验证码，并令其存活五分钟
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    email string, v string       接收一个邮箱和一个验证码
// @return   void
func SetRedisEmail(ctx *gin.Context, email string, v string) {
	client := common.GetRedisClient(0)

	client.Set(ctx, email, v, 300*time.Second)
}

// @title    ScoreChange
// @description   用于计算分数变化
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    email string, v string       接收一个邮箱和一个验证码
// @return   void
func ScoreChange(fre float64, sum float64, del float64, total float64) float64 {
	return (0.07/(fre+1) + 0.04) * sum * (math.Pow(2, 10*del-0.5)) / (math.Pow(2, 10*del-0.5) + 1) * (math.Pow(2, 0.1*total-5)) / (math.Pow(2, 0.1*total-5) + 1) / total
}
