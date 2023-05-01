// @Title  ProblemNewController
// @Description  该文件提供关于操作比赛题目的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"MGA_OJ/Interface"
	TQ "MGA_OJ/Test-request"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"
	"MGA_OJ/util"
	"MGA_OJ/vo"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// IProblemNewController			定义了比赛题目类接口
type IProblemNewController interface {
	Interface.RestInterface   // 包含了增删查改功能
	TestNum(ctx *gin.Context) // 查看指定题目的用例数量
	Quote(ctx *gin.Context)   // 引用题库
	Rematch(ctx *gin.Context) // 重现赛内题目
}

// ProblemNewController			定义了比赛题目工具类
type ProblemNewController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Create
// @description   创建一篇题目
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemNewController) Create(ctx *gin.Context) {
	var requestProblem vo.ProblemRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestProblem); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找比赛
	var competition model.Competition

	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "Competition", requestProblem.CompetitionId.String()).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Competition", requestProblem.CompetitionId.String()).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			goto leap
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Competition", requestProblem.CompetitionId.String())
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if p.DB.Where("id = ?", requestProblem.CompetitionId.String()).First(&competition).Error != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}

	{
		// TODO 将竞赛存入redis供下次使用
		v, _ := json.Marshal(competition)
		p.Redis.HSet(ctx, "Competition", requestProblem.CompetitionId.String(), v)
	}
leap:
	// TODO 查看比赛是否已经结束
	if time.Now().After(time.Time(competition.EndTime)) {
		response.Fail(ctx, nil, "比赛已经结束")
		return
	}
	// TODO 查看是否有权给比赛添加题目
	if competition.UserId != user.ID {
		if p.DB.Where("group_id = ? and user_id = ?", competition.GroupId, user.ID).First(&model.UserList{}).Error != nil {
			response.Fail(ctx, nil, "无权为比赛添加题目")
			return
		}
	}
	// TODO 尝试取出单位
	timeunits, ok := util.Units[strings.ToLower(requestProblem.TimeUnits)]

	if !ok {
		response.Fail(ctx, nil, "时间单位错误")
		return
	}

	memoryunits, ok := util.Units[strings.ToLower(requestProblem.MemoryUnits)]

	if !ok {
		response.Fail(ctx, nil, "内存单位错误")
		return
	}

	requestProblem.TimeLimit *= timeunits
	requestProblem.MemoryLimit *= memoryunits

	// TODO 查看时间限制是否合理
	if requestProblem.TimeLimit < 50 || requestProblem.TimeLimit > 10000 {
		response.Fail(ctx, nil, "时间限制不合理")
		return
	}

	// TODO 查看空间限制是否合理
	if requestProblem.MemoryLimit < 1 || requestProblem.MemoryLimit > 1024*1024*1 {
		response.Fail(ctx, nil, "空间限制不合理")
		return
	}

	// TODO 如果来源为空，为其设置默认值
	if requestProblem.Source == "" {
		requestProblem.Source = "用户" + user.Name + "上传"
	}

	// TODO 创建题目
	problem := model.ProblemNew{
		Title:         requestProblem.Title,
		TimeLimit:     requestProblem.TimeLimit,
		MemoryLimit:   requestProblem.MemoryLimit,
		Description:   requestProblem.Description,
		Reslong:       requestProblem.Reslong,
		Resshort:      requestProblem.Resshort,
		Input:         requestProblem.Input,
		Output:        requestProblem.Output,
		Hint:          requestProblem.Hint,
		Source:        requestProblem.Source,
		UserId:        user.ID,
		SpecialJudge:  requestProblem.SpecialJudge,
		Standard:      requestProblem.Standard,
		InputCheck:    requestProblem.InputCheck,
		CompetitionId: requestProblem.CompetitionId,
	}

	// TODO 如果样例输入数量与样例输出数量不对等
	if len(requestProblem.SampleInput) != len(requestProblem.SampleOutput) {
		response.Fail(ctx, nil, "题目的样例输入数量或输出数量有误")
		return
	}

	// TODO 如果用例输入数量与输出数量不对等
	if len(requestProblem.TestOutput) != len(requestProblem.TestInput) {
		response.Fail(ctx, nil, "题目的用例输入数量或输出数量有误")
		return
	}

	// TODO 分数数量与输入数量不对等
	if len(requestProblem.Scores) != len(requestProblem.TestInput) {
		response.Fail(ctx, nil, "题目的分数数量与输入数量不对等")
		return
	}

	// TODO 查看特判程序是否通过
	var program model.Program
	if p.DB.Where("id = ?", requestProblem.SpecialJudge).First(&program).Error == nil {
		for i := range requestProblem.TestInput {
			if condition, output := TQ.JudgeRun(program.Language, program.Code, requestProblem.TestInput[i]+"\n"+requestProblem.TestOutput[i], requestProblem.MemoryLimit*2, requestProblem.TimeLimit*2); condition != "ok" || output != "ok" {
				response.Fail(ctx, nil, "特判程序未通过")
				return
			}
		}
	}

	// TODO 查看标准程序是否通过
	if p.DB.Where("id = ?", requestProblem.Standard).First(&program).Error == nil {
		for i := range requestProblem.TestInput {
			if condition, output := TQ.JudgeRun(program.Language, program.Code, requestProblem.TestInput[i], requestProblem.MemoryLimit*2, requestProblem.TimeLimit*2); condition != "ok" || output != requestProblem.TestOutput[i] {
				response.Fail(ctx, nil, "标准程序未通过")
				return
			}
		}
	}

	// TODO 查看输入检查程序是否通过
	if p.DB.Where("id = ?", requestProblem.InputCheck).First(&program).Error == nil {
		for i := range requestProblem.TestInput {
			if condition, output := TQ.JudgeRun(program.Language, program.Code, requestProblem.TestInput[i], requestProblem.MemoryLimit*2, requestProblem.TimeLimit*2); condition != "ok" || output != "ok" {
				response.Fail(ctx, nil, "输入检查程序未通过")
				return
			}
		}
	}

	// TODO 插入数据
	if err := p.DB.Create(&problem).Error; err != nil {
		response.Fail(ctx, nil, "题目上传出错，数据验证有误")
		return
	}

	// TODO 存储测试样例
	for i := range requestProblem.SampleInput {
		// TODO 尝试存入数据库
		cas := model.CaseSample{
			ProblemId: problem.ID,
			Input:     requestProblem.SampleInput[i],
			Output:    requestProblem.SampleOutput[i],
			CID:       uint(i + 1),
		}
		// TODO 插入数据
		if err := p.DB.Create(&cas).Error; err != nil {
			response.Fail(ctx, nil, "题目用例上传出错，数据验证有误")
			return
		}
	}

	// TODO 存储测试用例
	for i := range requestProblem.TestInput {
		// TODO 尝试存入数据库
		cas := model.Case{
			ProblemId: problem.ID,
			Input:     requestProblem.TestInput[i],
			Output:    requestProblem.TestOutput[i],
			Score:     requestProblem.Scores[i],
			CID:       uint(i + 1),
		}
		// TODO 插入数据
		if err := p.DB.Create(&cas).Error; err != nil {
			response.Fail(ctx, nil, "题目用例上传出错，数据验证有误")
			return
		}
	}

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    Quote
// @description   引用一篇题目
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemNewController) Quote(ctx *gin.Context) {
	// 获取path中的id
	problem_id := ctx.Params.ByName("problem_id")
	competition_id := ctx.Params.ByName("competition_id")
	score, _ := strconv.Atoi(ctx.Params.ByName("score"))
	var problem model.Problem

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看题目的比赛是否开始
	var competition model.Competition

	// TODO 查看比赛是否存在
	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "Competition", competition_id).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Competition", competition_id).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			goto levp
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Competition", competition_id)
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if p.DB.Where("id = ?", competition_id).First(&competition).Error != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
	{
		// TODO 将比赛存入redis供下次使用
		v, _ := json.Marshal(competition)
		p.Redis.HSet(ctx, "Competition", competition_id, v)
	}
levp:
	// TODO 查看是否有权给比赛添加题目
	if competition.UserId != user.ID {
		if p.DB.Where("group_id = ? and user_id = ?", competition.GroupId, user.ID).First(&model.UserList{}).Error != nil {
			response.Fail(ctx, nil, "无权为比赛添加题目")
			return
		}
	}
	// TODO 查看比赛是否已经结束
	if time.Now().After(time.Time(competition.EndTime)) {
		response.Fail(ctx, nil, "比赛已经结束")
		return
	}

	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "Problem", problem_id).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Problem", problem_id).Result()
		if json.Unmarshal([]byte(cate), &problem) == nil {
			// TODO 跳过数据库搜寻problem过程
			goto leep
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Problem", problem_id)
		}
	}

	// TODO 查看题目是否在数据库中存在
	if p.DB.Where("id = ?", problem_id).First(&problem).Error != nil {
		response.Fail(ctx, nil, "题目不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(problem)
		p.Redis.HSet(ctx, "Problem", problem_id, v)
	}

leep:
	var caseSamples []model.CaseSample
	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "CaseSample", problem_id).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "CaseSample", problem_id).Result()
		if json.Unmarshal([]byte(cate), &caseSamples) == nil {
			// TODO 跳过数据库搜寻caseSample过程
			goto leap
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "CaseSample", problem_id)
		}
	}
	p.DB.Where("problem_id = ?", problem.ID).Find(&caseSamples)
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(caseSamples)
		p.Redis.HSet(ctx, "CaseSample", problem_id, v)
	}

leap:
	// TODO 从数据库中读出输入输出
	var cases []model.Case

	// TODO 查找用例
	if ok, _ := p.Redis.HExists(ctx, "Case", problem_id).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Case", problem_id).Result()
		if json.Unmarshal([]byte(cate), &cases) == nil {
			// TODO 跳过数据库搜寻testInputs过程
			goto Case
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Case", problem_id)
		}
	}

	// TODO 查看题目是否在数据库中存在
	p.DB.Where("id = ?", problem_id).Find(&cases)
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(cases)
		p.Redis.HSet(ctx, "Case", problem_id, v)
	}
Case:

	// TODO 创建题目
	problemNew := model.ProblemNew{
		Title:         problem.Title,
		TimeLimit:     problem.TimeLimit,
		MemoryLimit:   problem.MemoryLimit,
		Description:   problem.Description,
		Reslong:       problem.Reslong,
		Resshort:      problem.Resshort,
		Input:         problem.Input,
		Output:        problem.Output,
		Hint:          problem.Hint,
		Source:        problem.Source,
		UserId:        problem.UserId,
		SpecialJudge:  problem.SpecialJudge,
		Standard:      problem.Standard,
		InputCheck:    problem.InputCheck,
		CompetitionId: competition.ID,
	}

	// TODO 插入数据
	if err := p.DB.Create(&problemNew).Error; err != nil {
		response.Fail(ctx, nil, "题目上传出错，数据验证有误")
		return
	}

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
		if err := p.DB.Create(&cas).Error; err != nil {
			response.Fail(ctx, nil, "题目用例上传出错，数据验证有误")
			return
		}
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
			cas.Score = uint(score)
		}
		// TODO 插入数据
		if err := p.DB.Create(&cas).Error; err != nil {
			response.Fail(ctx, nil, "题目用例上传出错，数据验证有误")
			return
		}
	}

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    Rematch
// @description   重现赛内题
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemNewController) Rematch(ctx *gin.Context) {
	// 获取path中的id
	problem_id := ctx.Params.ByName("problem_id")
	competition_id := ctx.Params.ByName("competition_id")
	var problem model.ProblemNew
	var competition model.Competition

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "ProblemNew", problem_id).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "ProblemNew", problem_id).Result()
		if json.Unmarshal([]byte(cate), &problem) == nil {
			// TODO 跳过数据库搜寻problem过程
			goto leep
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "ProblemNew", problem_id)
		}
	}

	// TODO 查看题目是否在数据库中存在
	if p.DB.Where("id = ?", problem_id).First(&problem).Error != nil {
		response.Fail(ctx, nil, "题目不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(problem)
		p.Redis.HSet(ctx, "ProblemNew", problem_id, v)
	}

leep:

	// TODO 查看比赛是否存在
	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "Competition", problem.CompetitionId.String()).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Competition", problem.CompetitionId.String()).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			goto leap
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Competition", problem.CompetitionId.String())
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if p.DB.Where("id = ?", problem.CompetitionId.String()).First(&competition).Error != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
	{
		// TODO 将比赛存入redis供下次使用
		v, _ := json.Marshal(competition)
		p.Redis.HSet(ctx, "Competition", problem.CompetitionId.String(), v)
	}
leap:

	// TODO 查看比赛是否在结束
	if user.ID != competition.UserId && time.Now().Before(time.Time(competition.EndTime)) {
		response.Success(ctx, nil, "权限不足，请等待比赛结束")
		return
	}

	// TODO 查看比赛是否存在
	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "Competition", competition_id).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Competition", competition_id).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			goto levp
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Competition", competition_id)
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if p.DB.Where("id = ?", competition_id).First(&competition).Error != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
	{
		// TODO 将比赛存入redis供下次使用
		v, _ := json.Marshal(competition)
		p.Redis.HSet(ctx, "Competition", competition_id, v)
	}
levp:
	// TODO 查看是否有权给比赛添加题目
	if competition.UserId != user.ID {
		if p.DB.Where("group_id = ? and user_id = ?", competition.GroupId, user.ID).First(&model.UserList{}).Error != nil {
			response.Fail(ctx, nil, "无权为比赛添加题目")
			return
		}
	}
	// TODO 查看比赛是否已经结束
	if time.Now().After(time.Time(competition.EndTime)) {
		response.Fail(ctx, nil, "比赛已经结束")
		return
	}

	var caseSamples []model.CaseSample
	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "CaseSample", problem_id).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "CaseSample", problem_id).Result()
		if json.Unmarshal([]byte(cate), &caseSamples) == nil {
			// TODO 跳过数据库搜寻caseSample过程
			goto letp
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "CaseSample", problem_id)
		}
	}
	p.DB.Where("problem_id = ?", problem.ID).Find(&caseSamples)
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(caseSamples)
		p.Redis.HSet(ctx, "CaseSample", problem_id, v)
	}

letp:
	// TODO 从数据库中读出输入输出
	var cases []model.Case

	// TODO 查找用例
	if ok, _ := p.Redis.HExists(ctx, "Case", problem_id).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Case", problem_id).Result()
		if json.Unmarshal([]byte(cate), &cases) == nil {
			// TODO 跳过数据库搜寻testInputs过程
			goto Case
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Case", problem_id)
		}
	}

	// TODO 查看题目是否在数据库中存在
	p.DB.Where("id = ?", problem_id).Find(&cases)
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(cases)
		p.Redis.HSet(ctx, "Case", problem_id, v)
	}
Case:

	// TODO 创建题目
	problemNew := model.ProblemNew{
		Title:         problem.Title,
		TimeLimit:     problem.TimeLimit,
		MemoryLimit:   problem.MemoryLimit,
		Description:   problem.Description,
		Reslong:       problem.Reslong,
		Resshort:      problem.Resshort,
		Input:         problem.Input,
		Output:        problem.Output,
		Hint:          problem.Hint,
		Source:        problem.Source,
		UserId:        problem.UserId,
		SpecialJudge:  problem.SpecialJudge,
		Standard:      problem.Standard,
		InputCheck:    problem.InputCheck,
		CompetitionId: competition.ID,
	}

	// TODO 插入数据
	if err := p.DB.Create(&problemNew).Error; err != nil {
		response.Fail(ctx, nil, "题目上传出错，数据验证有误")
		return
	}

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
		if err := p.DB.Create(&cas).Error; err != nil {
			response.Fail(ctx, nil, "题目用例上传出错，数据验证有误")
			return
		}
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
		// TODO 插入数据
		if err := p.DB.Create(&cas).Error; err != nil {
			response.Fail(ctx, nil, "题目用例上传出错，数据验证有误")
			return
		}
	}

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    Update
// @description   更新一篇题目的内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemNewController) Update(ctx *gin.Context) {
	var requestProblem vo.ProblemRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestProblem); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看用户是否有权限上传题目
	if user.Level < 2 {
		response.Fail(ctx, nil, "用户权限不足")
		return
	}

	// TODO 查找对应题目
	id := ctx.Params.ByName("id")

	var problem model.ProblemNew

	if p.DB.Where("id = ?", id).First(&problem) != nil {
		response.Fail(ctx, nil, "题目不存在")
		return
	}

	// TODO 查看是否是用户作者
	if user.ID != problem.UserId {
		response.Fail(ctx, nil, "不是题目作者，无法修改题目")
		return
	}

	// TODO 尝试取出单位
	timeunits, ok := util.Units[strings.ToLower(requestProblem.TimeUnits)]

	if !ok {
		response.Fail(ctx, nil, "时间单位错误")
		return
	}

	memoryunits, ok := util.Units[strings.ToLower(requestProblem.MemoryUnits)]

	if !ok {
		response.Fail(ctx, nil, "内存单位错误")
		return
	}

	// TODO 化作默认单位
	requestProblem.TimeLimit *= timeunits
	requestProblem.MemoryLimit *= memoryunits

	requestProblem.TimeLimit *= timeunits
	requestProblem.MemoryLimit *= memoryunits

	// TODO 查看时间限制是否合理
	if requestProblem.TimeLimit != 0 && (requestProblem.TimeLimit < 50 || requestProblem.TimeLimit > 10000) {
		response.Fail(ctx, nil, "时间限制不合理")
		return
	}

	// TODO 查看空间限制是否合理
	if requestProblem.MemoryLimit != 0 && (requestProblem.MemoryLimit < 1 || requestProblem.MemoryLimit > 1024*1024*1) {
		response.Fail(ctx, nil, "空间限制不合理")
		return
	}

	// TODO 如果样例输入数量与样例输出数量不对等
	if len(requestProblem.SampleInput) != len(requestProblem.SampleOutput) {
		response.Fail(ctx, nil, "题目的样例输入数量或输出数量有误")
		return
	}

	// TODO 如果输入数量与输出数量不对等
	if len(requestProblem.TestOutput) != len(requestProblem.TestInput) {
		response.Fail(ctx, nil, "题目的输入数量或输出数量有误")
		return
	}

	// TODO 分数数量与输入数量不对等
	if len(requestProblem.Scores) != len(requestProblem.TestInput) {
		response.Fail(ctx, nil, "题目的分数数量与输入数量不对等")
		return
	}

	// TODO 查看特判程序是否通过
	var program model.Program
	if p.DB.Where("id = ?", requestProblem.SpecialJudge).First(&program).Error == nil {
		for i := range requestProblem.TestInput {
			if condition, output := TQ.JudgeRun(program.Language, program.Code, requestProblem.TestInput[i]+"\n"+requestProblem.TestOutput[i], requestProblem.MemoryLimit*2, requestProblem.TimeLimit*2); condition != "ok" || output != "ok" {
				response.Fail(ctx, nil, "特判程序未通过")
				return
			}
		}
	}

	// TODO 查看标准程序是否通过
	if p.DB.Where("id = ?", requestProblem.Standard).First(&program).Error == nil {
		for i := range requestProblem.TestInput {
			if condition, output := TQ.JudgeRun(program.Language, program.Code, requestProblem.TestInput[i], requestProblem.MemoryLimit*2, requestProblem.TimeLimit*2); condition != "ok" || output != requestProblem.TestOutput[i] {
				response.Fail(ctx, nil, "标准程序未通过")
				return
			}
		}
	}

	// TODO 查看输入检查程序是否通过
	if p.DB.Where("id = ?", requestProblem.InputCheck).First(&program).Error == nil {
		for i := range requestProblem.TestInput {
			if condition, output := TQ.JudgeRun(program.Language, program.Code, requestProblem.TestInput[i], requestProblem.MemoryLimit*2, requestProblem.TimeLimit*2); condition != "ok" || output != "ok" {
				response.Fail(ctx, nil, "输入检查程序未通过")
				return
			}
		}
	}

	// TODO 更新题目内容
	p.DB.Where("id = ?", id).Updates(model.ProblemNew{
		TimeLimit:     requestProblem.TimeLimit,
		MemoryLimit:   requestProblem.MemoryLimit,
		Title:         requestProblem.Title,
		Description:   requestProblem.Description,
		Reslong:       requestProblem.Reslong,
		Resshort:      requestProblem.Resshort,
		Input:         requestProblem.Input,
		Output:        requestProblem.Output,
		Hint:          requestProblem.Hint,
		Source:        requestProblem.Source,
		SpecialJudge:  requestProblem.SpecialJudge,
		Standard:      requestProblem.Standard,
		InputCheck:    requestProblem.InputCheck,
		CompetitionId: requestProblem.CompetitionId,
	})

	// TODO 移除损坏数据
	p.Redis.HDel(ctx, "ProblemNew", id)

	// TODO 查看输入样例是否变化
	if len(requestProblem.SampleInput) != 0 {
		p.Redis.HDel(ctx, "SampleCase", id)
		// TODO 清空原有的测试输入
		p.DB.Where("problem_id = ?", id).Delete(&model.CaseSample{})
		// TODO 存储测试输入
		for i := range requestProblem.SampleInput {
			// TODO 尝试存入数据库
			cas := model.CaseSample{
				ProblemId: problem.ID,
				Input:     requestProblem.SampleInput[i],
				Output:    requestProblem.SampleOutput[i],
				CID:       uint(i + 1),
			}
			// TODO 插入数据
			if err := p.DB.Create(&cas).Error; err != nil {
				response.Fail(ctx, nil, "题目样例上传出错，数据验证有误")
				return
			}
		}
	}

	// TODO 查看输入测试是否变化
	if len(requestProblem.TestInput) != 0 {
		p.Redis.HDel(ctx, "Case", id)
		// TODO 清空原有的测试输入
		p.DB.Where("problem_id = ?", id).Delete(&model.Case{})
		// TODO 存储测试输入
		for i := range requestProblem.TestInput {
			// TODO 尝试存入数据库
			cas := model.Case{
				ProblemId: problem.ID,
				Input:     requestProblem.TestInput[i],
				Output:    requestProblem.TestOutput[i],
				CID:       uint(i + 1),
				Score:     requestProblem.Scores[i],
			}
			// TODO 插入数据
			if err := p.DB.Create(&cas).Error; err != nil {
				response.Fail(ctx, nil, "题目用例上传出错，数据验证有误")
				return
			}
		}
	}

	// TODO 成功
	response.Success(ctx, nil, "更新成功")
}

// @title    Show
// @description   查看一篇题目的内容
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemNewController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var problem model.ProblemNew

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "ProblemNew", id).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "ProblemNew", id).Result()
		if json.Unmarshal([]byte(cate), &problem) == nil {
			// TODO 跳过数据库搜寻problem过程
			goto leep
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "ProblemNew", id)
		}
	}

	// TODO 查看题目是否在数据库中存在
	if p.DB.Where("id = ?", id).First(&problem).Error != nil {
		response.Fail(ctx, nil, "题目不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(problem)
		p.Redis.HSet(ctx, "ProblemNew", id, v)
	}

leep:
	// TODO 查看题目的比赛是否结束
	var competition model.Competition

	// TODO 查看比赛是否存在
	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "Competition", problem.CompetitionId.String()).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Competition", problem.CompetitionId.String()).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			goto leap
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Competition", problem.CompetitionId.String())
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if p.DB.Where("id = ?", problem.CompetitionId.String()).First(&competition).Error != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
	{
		// TODO 将比赛存入redis供下次使用
		v, _ := json.Marshal(competition)
		p.Redis.HSet(ctx, "Competition", problem.CompetitionId.String(), v)
	}
leap:

	// TODO 查看比赛是否在结束
	if user.ID != competition.UserId && time.Now().Before(time.Time(competition.EndTime)) {
		// TODO 查看是否报名
		var competitionRank model.CompetitionRank
		// TODO 查看是否已经报名
		// TODO 先看redis中是否存在
		if _, err := p.Redis.ZScore(ctx, "CompetitionR"+id, user.ID.String()).Result(); err != nil {
			if p.DB.Where("member_id = ? and competition_id = ?", user.ID, competition.ID).First(&competitionRank).Error != nil {
				if user.Level < 2 {
					response.Success(ctx, nil, "未报名")
					return
				}
			} else {
				// TODO 加入redis
				p.Redis.ZAdd(ctx, "CompetitionR"+id, redis.Z{Member: user.ID.String(), Score: 0})
			}
		}
	}

	var caseSamples []model.CaseSample
	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "CaseSample", id).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "CaseSample", id).Result()
		if json.Unmarshal([]byte(cate), &caseSamples) == nil {
			// TODO 跳过数据库搜寻caseSample过程
			goto levp
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "CaseSample", id)
		}
	}
	p.DB.Where("problem_id = ?", problem.ID).Find(&caseSamples)

	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(caseSamples)
		p.Redis.HSet(ctx, "CaseSample", id, v)
	}

levp:
	response.Success(ctx, gin.H{"problem": problem, "caseSamples": caseSamples}, "成功")
}

// @title    TestNum
// @description   查看一篇题目的测试样例数量
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemNewController) TestNum(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var problem model.ProblemNew

	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "ProblemNew", id).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "ProblemNew", id).Result()
		if json.Unmarshal([]byte(cate), &problem) == nil {
			// TODO 跳过数据库搜寻problem过程
			goto leep
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "ProblemNew", id)
		}
	}

	// TODO 查看题目是否在数据库中存在
	if p.DB.Where("id = ?", id).First(&problem).Error != nil {
		response.Fail(ctx, nil, "题目不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(problem)
		p.Redis.HSet(ctx, "ProblemNew", id, v)
	}

leep:

	var total int64
	p.DB.Where("problem_id = ?", id).Model(&model.Case{}).Count(&total)

	response.Success(ctx, gin.H{"total": total}, "成功")
}

// @title    Delete
// @description   删除一篇题目
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemNewController) Delete(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	var problem model.ProblemNew

	// TODO 查看题目是否存在
	if p.DB.Where("id = ?", id).First(&problem).Error != nil {
		response.Fail(ctx, nil, "题目不存在")
		return
	}

	// TODO 判断当前用户是否为比赛的作者
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看题目的比赛是否结束
	var competition model.Competition

	// TODO 查看比赛是否存在
	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "Competition", problem.CompetitionId.String()).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Competition", problem.CompetitionId.String()).Result()
		if json.Unmarshal([]byte(cate), &competition) == nil {
			goto leap
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Competition", problem.CompetitionId.String())
		}
	}

	// TODO 查看比赛是否在数据库中存在
	if p.DB.Where("id = ?", problem.CompetitionId.String()).First(&competition).Error != nil {
		response.Fail(ctx, nil, "比赛不存在")
		return
	}
	{
		// TODO 将比赛存入redis供下次使用
		v, _ := json.Marshal(competition)
		p.Redis.HSet(ctx, "Competition", problem.CompetitionId.String(), v)
	}
leap:

	// TODO 查看是否有操作题目的权力
	if user.ID != competition.UserId && user.Level < 4 {
		response.Fail(ctx, nil, "比赛不属于您，请勿非法操作")
		return
	}

	// TODO 删除题目
	p.DB.Delete(&problem)

	// TODO 移除损坏数据
	p.Redis.HDel(ctx, "ProblemNew", id)

	p.Redis.HDel(ctx, "Case", id)

	p.Redis.HDel(ctx, "CaseSample", id)

	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description   获取多篇题目
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemNewController) PageList(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var problems []model.ProblemNew

	// TODO 查找所有分页中可见的条目
	p.DB.Where("competition_id = ?", id).Order("created_at asc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&problems)

	var total int64
	p.DB.Model(model.ProblemNew{}).Count(&total)

	var problemIds []uuid.UUID

	for i := range problems {
		problemIds = append(problemIds, problems[i].ID)
	}

	// TODO 返回数据
	response.Success(ctx, gin.H{"problemIds": problemIds, "total": total}, "成功")
}

// @title    NewProblemNewController
// @description   新建一个IProblemNewController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   IProblemNewController		返回一个IProblemNewController用于调用各种函数
func NewProblemNewController() IProblemNewController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.ProblemNew{})
	return ProblemNewController{DB: db, Redis: redis}
}
