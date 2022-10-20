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
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IProblemController			定义了题目类接口
type IProblemController interface {
	Interface.RestInterface // 包含增删查改功能
}

// ProblemController			定义了题目工具类
type ProblemController struct {
	DB *gorm.DB // 含有一个数据库指针
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
		Title:        requestProblem.Title,
		TimeLimit:    requestProblem.TimeLimit,
		MemoryLimit:  requestProblem.MemoryLimit,
		Description:  requestProblem.Description,
		Reslong:      requestProblem.Reslong,
		Resshort:     requestProblem.Resshort,
		Input:        requestProblem.Input,
		Output:       requestProblem.Output,
		SampleInput:  requestProblem.SampleInput,
		SampleOutput: requestProblem.SampleOutput,
		Hint:         requestProblem.Hint,
		Source:       requestProblem.Source,
		UserId:       user.ID,
	}

	// TODO 插入数据
	if err := p.DB.Create(&problem).Error; err != nil {
		response.Fail(ctx, nil, "题目上传出错，数据验证有误")
		return
	}

	// TODO 创建题目目录
	if err := os.MkdirAll("./test_case/"+strconv.Itoa(int(problem.ID)), os.ModePerm); err != nil {
		response.Fail(ctx, nil, "题目目录创建失败")
		return
	}

	// TODO 存储测试输入
	if err := os.MkdirAll("./test_case/"+strconv.Itoa(int(problem.ID))+"/input", os.ModePerm); err != nil {
		response.Fail(ctx, nil, "题目输入测试目录创建失败")
		return
	}
	for i, val := range requestProblem.TestInput {
		file, err := os.Create("./test_case/" + strconv.Itoa(int(problem.ID)) + "/input/" + strconv.Itoa(i+1) + ".txt")
		if err != nil {
			response.Fail(ctx, nil, "题目第"+strconv.Itoa(i+1)+"个输入样例创建失败")
			return
		}
		// TODO 及时关闭file句柄
		defer file.Close()
		// TODO 写入文件时，使用带缓存的 *Writer
		write := bufio.NewWriter(file)
		write.WriteString(val)
		// TODO Flush将缓存的文件真正写入到文件中
		write.Flush()
	}

	// TODO 存储测试输出
	if err := os.MkdirAll("./test_case/"+strconv.Itoa(int(problem.ID))+"/output", os.ModePerm); err != nil {
		response.Fail(ctx, nil, "题目输入测试目录创建失败")
		return
	}
	for i, val := range requestProblem.TestOutput {
		file, err := os.Create("./test_case/" + strconv.Itoa(int(problem.ID)) + "/output/" + strconv.Itoa(i+1) + ".txt")
		if err != nil {
			response.Fail(ctx, nil, "题目第"+strconv.Itoa(i+1)+"个输出样例创建失败")
			return
		}
		// TODO 及时关闭file句柄
		defer file.Close()
		// TODO 写入文件时，使用带缓存的 *Writer
		write := bufio.NewWriter(file)
		write.WriteString(val)
		// TODO Flush将缓存的文件真正写入到文件中
		write.Flush()
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
	p.DB.Table("problems").Updates(requestProblem)

	// TODO 查看输入测试是否变化
	if len(requestProblem.TestInput) != 0 {
		// TODO 清空原有的测试输入
		os.RemoveAll("./test_case/" + strconv.Itoa(int(problem.ID)) + "/input")
		if err := os.MkdirAll("./test_case/"+strconv.Itoa(int(problem.ID))+"/input", os.ModePerm); err != nil {
			response.Fail(ctx, nil, "题目输入测试目录创建失败")
			return
		}
		for i, val := range requestProblem.TestInput {
			file, err := os.Create("./test_case/" + strconv.Itoa(int(problem.ID)) + "/input/" + strconv.Itoa(i+1) + ".txt")
			if err != nil {
				response.Fail(ctx, nil, "题目第"+strconv.Itoa(i+1)+"个输入样例创建失败")
				return
			}
			// TODO 及时关闭file句柄
			defer file.Close()
			// TODO 写入文件时，使用带缓存的 *Writer
			write := bufio.NewWriter(file)
			write.WriteString(val)
			// TODO Flush将缓存的文件真正写入到文件中
			write.Flush()
		}
	}

	// TODO 查看输出测试是否变化
	if len(requestProblem.TestOutput) != 0 {
		// TODO 清空原有的测试输入
		os.RemoveAll("./test_case/" + strconv.Itoa(int(problem.ID)) + "/output")
		if err := os.MkdirAll("./test_case/"+strconv.Itoa(int(problem.ID))+"/output", os.ModePerm); err != nil {
			response.Fail(ctx, nil, "题目输入测试目录创建失败")
			return
		}
		for i, val := range requestProblem.TestOutput {
			file, err := os.Create("./test_case/" + strconv.Itoa(int(problem.ID)) + "/output/" + strconv.Itoa(i+1) + ".txt")
			if err != nil {
				response.Fail(ctx, nil, "题目第"+strconv.Itoa(i+1)+"个输入样例创建失败")
				return
			}
			// TODO 及时关闭file句柄
			defer file.Close()
			// TODO 写入文件时，使用带缓存的 *Writer
			write := bufio.NewWriter(file)
			write.WriteString(val)
			// TODO Flush将缓存的文件真正写入到文件中
			write.Flush()
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

	// TODO 查看题目是否在数据库中存在
	if p.DB.Where("id = ?", id).First(&problem).Error != nil {
		response.Fail(ctx, nil, "题目不存在")
		return
	}

	response.Success(ctx, gin.H{"problem": problem}, "成功")
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

	// TODO 分页
	var problems []model.Problem

	// TODO 查找所有分页中可见的条目
	p.DB.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&problems)

	var total int64
	p.DB.Model(model.Problem{}).Count(&total)

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
	db.AutoMigrate(model.Problem{})
	return ProblemController{DB: db}
}
