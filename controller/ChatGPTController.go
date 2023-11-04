// @Title  ChatGPTController
// @Description  该文件提供关于操作ChatGPT的各种方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package controller

import (
	"MGA_OJ/response"
	"MGA_OJ/util"
	"MGA_OJ/vo"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// IChatGPTController			定义了ChatGPT类接口
type IChatGPTController interface {
	GenerateCode(ctx *gin.Context) // 按照注释生成代码
	GenerateNote(ctx *gin.Context) // 根据代码生成注释
	Change(ctx *gin.Context)       // 代码转换
	Opinion(ctx *gin.Context)      // 代码修改意见
}

// ChatGPTController			定义了ChatGPT工具类
type ChatGPTController struct {
}

// @title    GenerateCode
// @description   根据注释生成代码
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c ChatGPTController) GenerateCode(ctx *gin.Context) {
	// TODO 获取指定语言
	language := ctx.Params.ByName("language")

	// TODO 获取指定代码
	var requestText vo.TextRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestText); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	res, err := util.ChatGPT([]string{"请扮演一个" + language + "语言使用者", "接下来的提问，请仅给出完整代码，不要带有任何格式",
		`请根据以下代码的注释完善代码，并给出你的完整代码。
	` + requestText.Text}, "gpt-3.5-turbo")

	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}

	if len(res.Choices[0].Message.Content) >= 3 && res.Choices[0].Message.Content[0:3] == "```" {
		res.Choices[0].Message.Content = util.RemoveFirstAndLastLine(res.Choices[0].Message.Content)
	}

	// TODO 成功
	response.Success(ctx, gin.H{"res": res.Choices[0].Message.Content}, "生成成功")
}

// @title    GenerateNote
// @description   根据代码生成注释
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c ChatGPTController) GenerateNote(ctx *gin.Context) {
	// TODO 获取指定语言
	language := ctx.Params.ByName("language")

	// TODO 获取指定代码
	var requestText vo.TextRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestText); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	res, err := util.ChatGPT([]string{"请扮演一个" + language + "语言使用者", "接下来的提问，请仅给出完整代码，不要带有任何格式",
		`请根据以下代码完善注释，并给出你的完整代码。
	` + requestText.Text}, "gpt-3.5-turbo")

	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}

	if len(res.Choices[0].Message.Content) >= 3 && res.Choices[0].Message.Content[0:3] == "```" {
		res.Choices[0].Message.Content = util.RemoveFirstAndLastLine(res.Choices[0].Message.Content)
	}

	// TODO 成功
	response.Success(ctx, gin.H{"res": res.Choices[0].Message.Content}, "生成成功")
}

// @title    Change
// @description   根据代码生成代码
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c ChatGPTController) Change(ctx *gin.Context) {
	// TODO 获取指定语言
	language1 := ctx.Params.ByName("language1")
	language2 := ctx.Params.ByName("language2")

	// TODO 获取指定代码
	var requestText vo.TextRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestText); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	res, err := util.ChatGPT([]string{"接下来的提问，请仅给出完整代码，不要带有任何格式",
		`请将以下` + language1 + `代码转化成等价的` + language2 + `代码，并给出你的完整代码。
	` + requestText.Text}, "gpt-3.5-turbo")

	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}

	if len(res.Choices[0].Message.Content) >= 3 && res.Choices[0].Message.Content[0:3] == "```" {
		res.Choices[0].Message.Content = util.RemoveFirstAndLastLine(res.Choices[0].Message.Content)
	}

	// TODO 成功
	response.Success(ctx, gin.H{"res": res.Choices[0].Message.Content}, "生成成功")
}

// @title    Opinion
// @description   代码修改意见
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c ChatGPTController) Opinion(ctx *gin.Context) {
	// TODO 获取指定语言
	language := ctx.Params.ByName("language")

	// TODO 获取指定代码
	var requestProblemCode vo.ProblemCodeRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestProblemCode); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	Sample := ""
	for i := range requestProblemCode.Problem.SampleCase {
		Sample += "输入样例" + fmt.Sprint(i+1) + "：" + fmt.Sprint(requestProblemCode.Problem.SampleCase[i].Input)
		Sample += "输出样例" + fmt.Sprint(i+1) + "：" + fmt.Sprint(requestProblemCode.Problem.SampleCase[i].Output)
	}

	description, err := util.ChatGPT([]string{`在线评测中，用户需要编写代码完成指定题目的要求，用户提交代码后平台会对程序进行输入，并通过输出评测用户代码是否正确。
	用户提交的代码往往有以下运行结果：
	Wrong Answer：答案错误。
	Runtime Error：运行时错误，程序崩溃。
	Compile Error：编译错误。
	Time Limit Exceeded：运行超出时间限制。
	Memory Limit Exceeded：超出内存限制。
	Output Limit Exceeded：输出的长度超过限制。
	Presentation Error（：答案正确，但是输出格式不匹配题目要求。
	Accepted：程序通过。
	`,
		"题目如下：\n" + "描述：" + requestProblemCode.Problem.Description + "\n" +
			"输入：" + requestProblemCode.Problem.Input + "\n" +
			"输出：" + requestProblemCode.Problem.Output + "\n" +
			"时间限制：" + fmt.Sprint(requestProblemCode.Problem.TimeLimit) + requestProblemCode.Problem.TimeUnits + "\n" +
			"空间限制：" + fmt.Sprint(requestProblemCode.Problem.MemoryLimit) + requestProblemCode.Problem.MemoryUnits + "\n" +
			Sample,
		`请判断以下` + language + `代码在在线评测中的运行结果(及前述中的一种)，并通过具体的分析给出原因，仅当运行结果不为Accepted时给出详细的修改意见。
	` + requestProblemCode.Code,
	}, "gpt-4")

	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}

	res := ""

	for i := range description.Choices[0].Message.Content {
		if description.Choices[0].Message.Content[i] == 'A' {
			res = "Accepted"
			break
		} else if description.Choices[0].Message.Content[i] == 'W' {
			res = "Wrong Answer"
			break
		} else if description.Choices[0].Message.Content[i] == 'R' {
			res = "Runtime Error"
			break
		} else if description.Choices[0].Message.Content[i] == 'C' {
			res = "Compile Error"
			break
		} else if description.Choices[0].Message.Content[i] == 'T' {
			res = "Time Limit Exceeded"
			break
		} else if description.Choices[0].Message.Content[i] == 'M' {
			res = "Memory Limit Exceeded"
			break
		} else if description.Choices[0].Message.Content[i] == 'O' {
			res = "Output Limit Exceeded"
			break
		} else if description.Choices[0].Message.Content[i] == 'P' {
			res = "Presentation Error"
			break
		}
	}

	// TODO 成功
	response.Success(ctx, gin.H{"res": res, "description": description.Choices[0].Message.Content}, "返回成功")
}

// @title    NewChatGPTController
// @description   新建一个IChatGPTController
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   IChatGPTController		返回一个IChatGPTController用于调用各种函数
func NewChatGPTController() IChatGPTController {
	return ChatGPTController{}
}
