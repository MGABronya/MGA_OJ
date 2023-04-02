// @Title  ProblemController
// @Description  该文件提供关于操作题目的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"MGA_OJ/Interface"
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

// IProblemController			定义了题目类接口
type IProblemController interface {
	Interface.RestInterface    // 包含增删查改功能
	Interface.LikeInterface    // 包含点赞功能
	Interface.CollectInterface // 包含收藏功能
	Interface.VisitInterface   // 包含游览功能
	Interface.LabelInterface   // 包含标签功能
	Interface.SearchInterface  // 包含搜索功能
	UserList(ctx *gin.Context) // 查看指定用户上传的题目列表
	TestNum(ctx *gin.Context)  // 查看指定题目的样例数量
}

// ProblemController			定义了题目工具类
type ProblemController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Create
// @description   创建一篇题目
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) Create(ctx *gin.Context) {
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

	// TODO 取出用户权限
	if user.Level < 2 {
		response.Fail(ctx, nil, "用户权限不足")
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

	requestProblem.TimeLimit *= timeunits
	requestProblem.MemoryLimit *= memoryunits

	// TODO 查看时间限制是否合理
	if requestProblem.TimeLimit < 50 || requestProblem.TimeLimit > 60000 {
		response.Fail(ctx, nil, "时间限制不合理")
		return
	}

	// TODO 查看空间限制是否合理
	if requestProblem.MemoryLimit < 1 || requestProblem.MemoryLimit > 1024*1024*2 {
		response.Fail(ctx, nil, "空间限制不合理")
		return
	}

	// TODO 如果来源为空，为其设置默认值
	if requestProblem.Source == "" {
		requestProblem.Source = "用户" + user.Name + "上传"
	}

	// TODO 创建题目
	problem := model.Problem{
		Title:         requestProblem.Title,
		TimeLimit:     requestProblem.TimeLimit,
		MemoryLimit:   requestProblem.MemoryLimit,
		Description:   requestProblem.Description,
		Reslong:       requestProblem.Reslong,
		Resshort:      requestProblem.Resshort,
		Input:         requestProblem.Input,
		Output:        requestProblem.Output,
		SampleInput:   requestProblem.SampleInput,
		SampleOutput:  requestProblem.SampleOutput,
		Hint:          requestProblem.Hint,
		Source:        requestProblem.Source,
		UserId:        user.ID,
		CompetitionId: requestProblem.CompetitionId,
		SpecialJudge:  requestProblem.SpecialJudge,
	}

	// TODO 插入数据
	if err := p.DB.Create(&problem).Error; err != nil {
		response.Fail(ctx, nil, "题目上传出错，数据验证有误")
		return
	}

	// TODO 存储测试输入
	for i, val := range requestProblem.TestInput {
		// TODO 尝试存入数据库
		testInput := model.TestInput{
			ProblemId: problem.ID,
			Input:     val,
			Id:        uint(i + 1),
		}
		// TODO 插入数据
		if err := p.DB.Create(&testInput).Error; err != nil {
			response.Fail(ctx, nil, "题目上传出错，数据验证有误")
			return
		}
	}

	// TODO 存储测试输出
	for i, val := range requestProblem.TestOutput {
		// TODO 尝试存入数据库
		testOutput := model.TestOutput{
			ProblemId: problem.ID,
			Output:    val,
			Id:        uint(i + 1),
		}
		// TODO 插入数据
		if err := p.DB.Create(&testOutput).Error; err != nil {
			response.Fail(ctx, nil, "题目上传出错，数据验证有误")
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
func (p ProblemController) Update(ctx *gin.Context) {
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

	var problem model.Problem

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
	if requestProblem.TimeLimit != 0 && (requestProblem.TimeLimit < 50 || requestProblem.TimeLimit > 60000) {
		response.Fail(ctx, nil, "时间限制不合理")
		return
	}

	// TODO 查看空间限制是否合理
	if requestProblem.MemoryLimit != 0 && (requestProblem.MemoryLimit < 1 || requestProblem.MemoryLimit > 1024*1024*2) {
		response.Fail(ctx, nil, "空间限制不合理")
		return
	}

	// TODO 更新题目内容
	p.DB.Table("problems").Where("id = ?", id).Updates(model.Problem{
		TimeLimit:     requestProblem.TimeLimit,
		MemoryLimit:   requestProblem.MemoryLimit,
		Title:         requestProblem.Title,
		Description:   requestProblem.Description,
		Reslong:       requestProblem.Reslong,
		Resshort:      requestProblem.Resshort,
		Input:         requestProblem.Input,
		Output:        requestProblem.Output,
		SampleInput:   requestProblem.SampleInput,
		SampleOutput:  requestProblem.SampleOutput,
		Hint:          requestProblem.Hint,
		Source:        requestProblem.Source,
		CompetitionId: requestProblem.CompetitionId,
		SpecialJudge:  requestProblem.SpecialJudge,
	})

	// TODO 移除损坏数据
	p.Redis.HDel(ctx, "Problem", id)

	// TODO 查看输入测试是否变化
	if len(requestProblem.TestInput) != 0 {
		p.Redis.HDel(ctx, "Input", id)
		// TODO 清空原有的测试输入
		p.DB.Where("problem_id = ?", id).Delete(&model.TestInput{})
		// TODO 存储测试输入
		for i, val := range requestProblem.TestInput {
			// TODO 尝试存入数据库
			testInput := model.TestInput{
				ProblemId: problem.ID,
				Input:     val,
				Id:        uint(i + 1),
			}
			// TODO 插入数据
			if err := p.DB.Create(&testInput).Error; err != nil {
				response.Fail(ctx, nil, "题目上传出错，数据验证有误")
				return
			}
		}
	}

	// TODO 查看输出测试是否变化
	if len(requestProblem.TestOutput) != 0 {
		p.Redis.HDel(ctx, "Output", id)
		// TODO 清空原有的测试输出
		p.DB.Where("problem_id = ?", id).Delete(&model.TestOutput{})
		// TODO 存储测试输出
		for i, val := range requestProblem.TestOutput {
			// TODO 尝试存入数据库
			testOutput := model.TestOutput{
				ProblemId: problem.ID,
				Output:    val,
				Id:        uint(i + 1),
			}
			// TODO 插入数据
			if err := p.DB.Create(&testOutput).Error; err != nil {
				response.Fail(ctx, nil, "题目上传出错，数据验证有误")
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
func (p ProblemController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var problem model.Problem

	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "Problem", id).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Problem", id).Result()
		if json.Unmarshal([]byte(cate), &problem) == nil {
			// TODO 跳过数据库搜寻problem过程
			goto leep
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Problem", id)
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
		p.Redis.HSet(ctx, "Problem", id, v)
	}

leep:
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看problem的competition
	if problem.CompetitionId != (uuid.UUID{}) {
		var competition model.Competition
		// TODO 无法找到比赛，则返回题目
		if p.DB.Where("id = ?", problem.CompetitionId).First(&competition).Error != nil {
			response.Success(ctx, gin.H{"problem": problem}, "成功")
			return
		}
		// TODO 查看比赛是否已经结束
		if time.Now().After(time.Time(competition.EndTime)) {
			response.Success(ctx, gin.H{"problem": problem}, "成功")
			return
		}
		// TODO 查看比赛是否已经开始
		if !time.Now().After(time.Time(competition.StartTime)) {
			response.Fail(ctx, nil, "题目不存在")
			return
		}
		// TODO 查看用户是否参加了比赛
		var problemLists []model.ProblemList
		p.DB.Where("set_id = ?", competition.SetId).Find(&problemLists)
		for _, problemList := range problemLists {
			var userLists []model.UserList
			p.DB.Where("problem_id = ?", problemList.ProblemId).Find(&userLists)
			for _, userList := range userLists {
				if userList.UserId == user.ID {
					response.Success(ctx, gin.H{"problem": problem}, "成功")
					return
				}
			}
		}
		// TODO 没有参加比赛
		response.Fail(ctx, nil, "题目不存在")
		return
	}

	response.Success(ctx, gin.H{"problem": problem}, "成功")
}

// @title    TestNum
// @description   查看一篇题目的测试样例数量
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) TestNum(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var problem model.Problem

	// TODO 查看题目是否在数据库中存在
	if p.DB.Where("id = ?", id).First(&problem).Error != nil {
		response.Fail(ctx, nil, "题目不存在")
		return
	}

	var total int64
	p.DB.Where("problem_id = ?", id).Model(&model.TestInput{}).Count(&total)

	response.Success(ctx, gin.H{"total": total}, "成功")
}

// @title    Delete
// @description   删除一篇题目
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) Delete(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	var problem model.Problem

	// TODO 查看题目是否存在
	if p.DB.Where("id = ?", id).First(&problem).Error != nil {
		response.Fail(ctx, nil, "题目不存在")
		return
	}

	// TODO 判断当前用户是否为题目的作者
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有操作题目的权力
	if user.ID != problem.UserId && user.Level < 4 {
		response.Fail(ctx, nil, "题目不属于您，请勿非法操作")
		return
	}

	// TODO 删除题目
	p.DB.Delete(&problem)

	// TODO 移除损坏数据
	p.Redis.HDel(ctx, "Problem", id)

	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description   获取多篇题目
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) PageList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 尝试获取所有没有结束的比赛
	var competitions []model.Competition
	p.DB.Where("end_time > ?", time.Now()).Find(&competitions)

	// TODO 用于记录没有被用户参加的比赛
	userNotJoin := make([]uuid.UUID, 0)

	// TODO 查看哪些比赛没有被用户参加或比赛未开始
	for _, competition := range competitions {
		// TODO 查看比赛是否已经开始
		if !time.Now().After(time.Time(competition.StartTime)) {
			userNotJoin = append(userNotJoin, competition.ID)
			continue
		}
		// TODO 查看用户是否加入比赛
		ok := false
		var problemLists []model.ProblemList
		p.DB.Where("set_id = ?", competition.SetId).Find(&problemLists)
		for _, problemList := range problemLists {
			var userLists []model.UserList
			p.DB.Where("problem_id = ?", problemList.ProblemId).Find(&userLists)
			for _, userList := range userLists {
				if user.ID == userList.UserId {
					ok = true
					break
				}
			}
			if ok {
				break
			}
		}
		if !ok {
			userNotJoin = append(userNotJoin, competition.ID)
		}
	}

	// TODO 分页
	var problems []model.Problem

	// TODO 查找所有分页中可见的条目
	p.DB.Where("competition_id not in ?", userNotJoin).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&problems)

	var total int64
	p.DB.Where("competition_id not in ?", userNotJoin).Model(model.Problem{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"problems": problems, "total": total}, "成功")
}

// @title    UserList
// @description   获取指定用户的多篇题目
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) UserList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 尝试获取所有没有结束的比赛
	var competitions []model.Competition
	p.DB.Where("end_time > ?", time.Now()).Find(&competitions)

	// TODO 用于记录没有被用户参加的比赛
	userNotJoin := make([]uuid.UUID, 0)

	// TODO 查看哪些比赛没有被用户参加或比赛未开始
	for _, competition := range competitions {
		// TODO 查看比赛是否已经开始
		if !time.Now().After(time.Time(competition.StartTime)) {
			userNotJoin = append(userNotJoin, competition.ID)
			continue
		}
		// TODO 查看用户是否加入比赛
		ok := false
		var problemLists []model.ProblemList
		p.DB.Where("set_id = ?", competition.SetId).Find(&problemLists)
		for _, problemList := range problemLists {
			var userLists []model.UserList
			p.DB.Where("problem_id = ?", problemList.ProblemId).Find(&userLists)
			for _, userList := range userLists {
				if user.ID == userList.UserId {
					ok = true
					break
				}
			}
			if ok {
				break
			}
		}
		if !ok {
			userNotJoin = append(userNotJoin, competition.ID)
		}
	}

	// TODO 取出指定用户的id
	id := ctx.Params.ByName("id")

	// TODO 分页
	var problems []model.Problem

	// TODO 查找所有分页中可见的条目
	p.DB.Where("user_id = ? and competition_id not in ?", id, userNotJoin).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&problems)

	var total int64
	p.DB.Where("user_id = ? and competition_id not in ?", id, userNotJoin).Model(model.Problem{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"problems": problems, "total": total}, "成功")
}

// @title    Like
// @description   点赞或点踩
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) Like(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	var problem model.Problem

	// TODO 查看题目是否存在
	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "Problem", id).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Problem", id).Result()
		if json.Unmarshal([]byte(cate), &problem) == nil {
			// TODO 跳过数据库搜寻problem过程
			goto leep
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Problem", id)
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
		p.Redis.HSet(ctx, "Problem", id, v)
	}

leep:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有点赞或者点踩
	if p.DB.Where("user_id = ? and problem_id = ?", user.ID, id).Update("like", like).Error != nil {
		// TODO 插入数据
		problemLike := model.ProblemLike{
			ProblemId: problem.ID,
			UserId:    user.ID,
			Like:      like,
		}
		if err := p.DB.Create(&problemLike).Error; err != nil {
			response.Fail(ctx, nil, "点赞出错，数据库存储错误")
			return
		}
	}

	response.Success(ctx, nil, "点赞成功")
}

// @title    CancelLike
// @description   取消点赞或者点踩
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) CancelLike(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 取消点赞或者点踩
	p.DB.Where("user_id = ? and problem_id = ?", user.ID, id).Delete(&model.ProblemLike{})

	response.Success(ctx, nil, "取消成功")
}

// @title    LikeNumber
// @description   点赞或点踩的数量
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) LikeNumber(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	var total int64

	// TODO 查看点赞或者点踩的数量
	p.DB.Where("problem_id = ? and like = ?", id, like).Model(model.ProblemLike{}).Count(&total)

	response.Success(ctx, gin.H{"total": total}, "查看成功")
}

// @title    LikeList
// @description   点赞或点踩的列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) LikeList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var problemLikes []model.ProblemLike

	var total int64

	// TODO 查看点赞或者点踩的数量
	p.DB.Where("problem_id = ? and like = ?", id, like).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&problemLikes)

	p.DB.Where("problem_id = ? and like = ?", id, like).Model(model.ProblemLike{}).Count(&total)

	response.Success(ctx, gin.H{"problemLikes": problemLikes, "total": total}, "查看成功")
}

// @title    LikeShow
// @description   查看用户点赞状态
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) LikeShow(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var problemLike model.ProblemLike

	// TODO 查看点赞状态
	if p.DB.Where("user_id = ? and problem_id = ?", user.ID, id).First(&problemLike).Error != nil {
		response.Success(ctx, gin.H{"like": 0}, "暂无状态")
		return
	}

	if problemLike.Like {
		response.Success(ctx, gin.H{"like": 1}, "已点赞")
	} else {
		response.Success(ctx, gin.H{"like": -1}, "已点踩")
	}

}

// @title    Likes
// @description   查看用户点赞状态
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) Likes(ctx *gin.Context) {

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取指定用户用户
	id := ctx.Params.ByName("id")

	// TODO 分页
	var problemLikes []model.ProblemLike

	var total int64

	// TODO 查看点赞或者点踩的数量
	p.DB.Where("user_id = ? and like = ?", id, like).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&problemLikes)

	p.DB.Where("user_id = ? and like = ?", id, like).Model(model.ProblemLike{}).Count(&total)

	response.Success(ctx, gin.H{"problemLikes": problemLikes, "total": total}, "查看成功")
}

// @title    Collect
// @description   收藏
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) Collect(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	var problem model.Problem

	// TODO 查看题目是否存在
	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "Problem", id).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Problem", id).Result()
		if json.Unmarshal([]byte(cate), &problem) == nil {
			// TODO 跳过数据库搜寻problem过程
			goto leep
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Problem", id)
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
		p.Redis.HSet(ctx, "Problem", id, v)
	}

leep:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有收藏
	if p.DB.Where("user_id = ? and problem_id = ?", user.ID, problem.ID).First(&model.ProblemCollect{}).Error != nil {
		problemCollect := model.ProblemCollect{
			ProblemId: problem.ID,
			UserId:    user.ID,
		}
		// TODO 插入数据
		if err := p.DB.Create(&problemCollect).Error; err != nil {
			response.Fail(ctx, nil, "收藏出错，数据库存储错误")
			return
		}
	} else {
		response.Fail(ctx, nil, "已收藏")
		return
	}

	response.Success(ctx, nil, "收藏成功")
}

// @title    CancelCollect
// @description   取消收藏
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) CancelCollect(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有收藏
	if p.DB.Where("user_id = ? and problem_id = ?", user.ID, id).First(&model.ProblemCollect{}).Error != nil {
		response.Fail(ctx, nil, "未收藏")
		return
	} else {
		p.DB.Where("user_id = ? and problem_id = ?", user.ID, id).Delete(&model.ProblemCollect{})
		response.Success(ctx, nil, "取消收藏成功")
		return
	}
}

// @title    CollectShow
// @description   查看收藏状态
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) CollectShow(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有收藏
	if p.DB.Where("user_id = ? and problem_id = ?", user.ID, id).First(&model.ProblemCollect{}).Error != nil {
		response.Success(ctx, gin.H{"collect": false}, "未收藏")
		return
	} else {
		response.Success(ctx, gin.H{"collect": true}, "已收藏")
		return
	}
}

// @title    CollectList
// @description   查看收藏用户列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) CollectList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var problemCollects []model.ProblemCollect

	var total int64

	// TODO 查看收藏的数量
	p.DB.Where("problem_id = ?", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&problemCollects)

	p.DB.Where("problem_id = ?", id).Model(model.ProblemCollect{}).Count(&total)

	response.Success(ctx, gin.H{"problemCollects": problemCollects, "total": total}, "查看成功")
}

// @title    CollectNumber
// @description   查看收藏用户数量
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) CollectNumber(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	var total int64

	// TODO 查看收藏的数量
	p.DB.Where("problem_id = ?", id).Model(model.ProblemCollect{}).Count(&total)

	response.Success(ctx, gin.H{"total": total}, "查看成功")
}

// @title    Collects
// @description   查看用户收藏夹
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) Collects(ctx *gin.Context) {
	// TODO 获取指定用户用户
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var problemCollects []model.ProblemCollect

	var total int64

	// TODO 查看收藏的数量
	p.DB.Where("user_id = ?", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&problemCollects)

	p.DB.Where("user_id = ?", id).Model(model.ProblemCollect{}).Count(&total)

	response.Success(ctx, gin.H{"problemCollects": problemCollects, "total": total}, "查看成功")
}

// @title    Visit
// @description   游览题目
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) Visit(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	var problem model.Problem

	// TODO 查看题目是否存在
	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "Problem", id).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Problem", id).Result()
		if json.Unmarshal([]byte(cate), &problem) == nil {
			// TODO 跳过数据库搜寻problem过程
			goto leep
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Problem", id)
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
		p.Redis.HSet(ctx, "Problem", id, v)
	}

leep:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	problemVisit := model.ProblemVisit{
		ProblemId: problem.ID,
		UserId:    user.ID,
	}

	// TODO 插入数据
	if err := p.DB.Create(&problemVisit).Error; err != nil {
		response.Fail(ctx, nil, "题目游览上传出错，数据库存储错误")
		return
	}

	// TODO 添加入阅读库
	p.Redis.PFAdd(ctx, "ProblemVisit", id)

	response.Success(ctx, nil, "题目游览成功")
}

// @title    VisitNumber
// @description   游览题目数量
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) VisitNumber(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取阅读人数
	total, _ := p.Redis.PFCount(ctx, "ProblemVisit", id).Result()

	response.Success(ctx, gin.H{"total": total}, "请求题目游览数目成功")
}

// @title    VisitList
// @description   游览题目列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) VisitList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var problemVisits []model.ProblemVisit

	var total int64

	// TODO 查看收藏的数量
	p.DB.Where("problem_id = ?", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&problemVisits)

	p.DB.Where("problem_id = ?", id).Model(model.ProblemVisit{}).Count(&total)

	response.Success(ctx, gin.H{"problemVisits": problemVisits, "total": total}, "查看成功")
}

// @title    Visits
// @description   游览题目列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) Visits(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var problemVisits []model.ProblemVisit

	var total int64

	// TODO 查看收藏的数量
	p.DB.Where("user_id = ?", user.ID).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&problemVisits)

	p.DB.Where("user_id = ?", user.ID).Model(model.ProblemVisit{}).Count(&total)

	response.Success(ctx, gin.H{"problemVisits": problemVisits, "total": total}, "查看成功")
}

// @title    LabelCreate
// @description   标签创建
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) LabelCreate(ctx *gin.Context) {
	// TODO 获取指定题目
	id := ctx.Params.ByName("id")

	// TODO 获取标签
	label := ctx.Params.ByName("label")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看题目是否存在
	var problem model.Problem

	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "Problem", id).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "Problem", id).Result()
		if json.Unmarshal([]byte(art), &problem) == nil {
			goto leep
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "Problem", id)
		}
	}

	// TODO 查看题目是否在数据库中存在
	if p.DB.Where("id = ?", id).First(&problem).Error != nil {
		response.Fail(ctx, nil, "题目不存在")
		return
	}
	{
		// TODO 将题目存入redis供下次使用
		v, _ := json.Marshal(problem)
		p.Redis.HSet(ctx, "Problem", id, v)
	}
leep:

	// TODO 查看是否为题目作者
	if problem.UserId != user.ID {
		response.Fail(ctx, nil, "不是题目作者，请勿非法操作")
		return
	}

	// TODO 创建标签
	problemLabel := model.ProblemLabel{
		Label:     label,
		ProblemId: problem.ID,
	}

	// TODO 插入数据
	if err := p.DB.Create(&problemLabel).Error; err != nil {
		response.Fail(ctx, nil, "题目标签上传出错，数据验证有误")
		return
	}

	// TODO 解码失败，删除字段
	p.Redis.HDel(ctx, "ProblemLabel", id)

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    LabelDelete
// @description   标签删除
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) LabelDelete(ctx *gin.Context) {
	// TODO 获取指定题目
	id := ctx.Params.ByName("id")

	// TODO 获取标签
	label := ctx.Params.ByName("label")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看题目是否存在
	var problem model.Problem

	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "Problem", id).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "Problem", id).Result()
		if json.Unmarshal([]byte(art), &problem) == nil {
			goto leep
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "Problem", id)
		}
	}

	// TODO 查看题目是否在数据库中存在
	if p.DB.Where("id = ?", id).First(&problem).Error != nil {
		response.Fail(ctx, nil, "题目不存在")
		return
	}
	{
		// TODO 将题目存入redis供下次使用
		v, _ := json.Marshal(problem)
		p.Redis.HSet(ctx, "Problem", id, v)
	}
leep:

	// TODO 查看是否为题目作者
	if problem.UserId != user.ID {
		response.Fail(ctx, nil, "不是题目作者，请勿非法操作")
		return
	}

	// TODO 删除题目标签
	if p.DB.Where("id = ?", label).First(&model.ProblemLabel{}).Error != nil {
		response.Fail(ctx, nil, "标签不存在")
		return
	}

	p.DB.Where("id = ?", label).Delete(&model.ProblemLabel{})

	// TODO 解码失败，删除字段
	p.Redis.HDel(ctx, "ProblemLabel", id)

	// TODO 成功
	response.Success(ctx, nil, "删除成功")
}

// @title    LabelShow
// @description   标签查看
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) LabelShow(ctx *gin.Context) {
	// TODO 获取指定题目
	id := ctx.Params.ByName("id")

	// TODO 查找数据
	var problemLabels []model.ProblemLabel
	// TODO 先尝试在redis中寻找
	if ok, _ := p.Redis.HExists(ctx, "ProblemLabel", id).Result(); ok {
		art, _ := p.Redis.HGet(ctx, "ProblemLabel", id).Result()
		if json.Unmarshal([]byte(art), &problemLabels) == nil {
			goto leap
		} else {
			// TODO 解码失败，删除字段
			p.Redis.HDel(ctx, "ProblemLabel", id)
		}
	}

	// TODO 在数据库中查找
	p.DB.Where("problem_id = ?", id).Find(&problemLabels)
	{
		// TODO 将题目标签存入redis供下次使用
		v, _ := json.Marshal(problemLabels)
		p.Redis.HSet(ctx, "ProblemLabel", id, v)
	}

leap:

	// TODO 成功
	response.Success(ctx, gin.H{"problemLabels": problemLabels}, "查看成功")
}

// @title    Search
// @description   文本搜索
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) Search(ctx *gin.Context) {
	// TODO 获取文本
	text := ctx.Params.ByName("text")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	var problems []model.Problem

	// TODO 模糊匹配
	p.DB.Where("match(title,discription,res_long,res_short) against(? in boolean mode)", text+"*").Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&problems)

	// TODO 查看查询总数
	var total int64
	p.DB.Where("match(title,discription,res_long,res_short) against(? in boolean mode)", text+"*").Model(model.Problem{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"problems": problems, "total": total}, "成功")
}

// @title    SearchLabel
// @description   指定标签的搜索
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) SearchLabel(ctx *gin.Context) {

	var requestLabels vo.LabelsRequest

	// TODO 获取标签
	if err := ctx.ShouldBind(&requestLabels); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 通过标签寻找
	var problemIds []struct {
		ProblemId uuid.UUID `json:"problem_id"` // 题目外键
	}

	// TODO 进行标签匹配
	p.DB.Distinct("problem_id").Where("label in (?)", requestLabels.Labels).Model(model.ProblemLabel{}).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&problemIds)

	// TODO 查看查询总数
	var total int64
	p.DB.Distinct("problem_id").Where("label in (?)", requestLabels.Labels).Model(model.ProblemLabel{}).Count(&total)

	// TODO 查找对应题目
	var problems []model.Problem

	p.DB.Where("id in (?)", problemIds).Find(&problems)

	// TODO 返回数据
	response.Success(ctx, gin.H{"problems": problems, "total": total}, "成功")
}

// @title    SearchWithLabel
// @description   指定标签与文本的搜索
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p ProblemController) SearchWithLabel(ctx *gin.Context) {

	// TODO 获取文本
	text := ctx.Params.ByName("text")

	var requestLabels vo.LabelsRequest

	// TODO 获取标签
	if err := ctx.ShouldBind(&requestLabels); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 通过标签寻找
	var problemIds []struct {
		ProblemId uuid.UUID `json:"problem_id"` // 题目外键
	}

	// TODO 进行标签匹配
	p.DB.Distinct("problem_id").Where("label in (?)", requestLabels.Labels).Model(model.ProblemLabel{}).Find(&problemIds)

	// TODO 查找对应题目
	var problems []model.Problem

	// TODO 模糊匹配
	p.DB.Where("id in (?) and match(title,discription,res_long,res_short) against(? in boolean mode)", problemIds, text+"*").Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&problems)

	// TODO 查看查询总数
	var total int64
	p.DB.Where("id in (?) and match(title,discription,res_long,res_short) against(? in boolean mode)", problemIds, text+"*").Model(model.Problem{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"problems": problems, "total": total}, "成功")
}

// @title    NewProblemController
// @description   新建一个IProblemController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   IProblemController		返回一个IProblemController用于调用各种函数
func NewProblemController() IProblemController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.Problem{})
	db.AutoMigrate(model.ProblemCollect{})
	db.AutoMigrate(model.ProblemLike{})
	db.AutoMigrate(model.ProblemVisit{})
	db.AutoMigrate(model.TestInput{})
	db.AutoMigrate(model.TestOutput{})
	db.AutoMigrate(model.ProblemLabel{})
	return ProblemController{DB: db, Redis: redis}
}
