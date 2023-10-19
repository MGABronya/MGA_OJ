// @Title  EmailController
// @Description  该文件提供关于操作邮件的各种方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package controller

import (
	"MGA_OJ/response"
	"MGA_OJ/util"
	"MGA_OJ/vo"
	"log"

	"github.com/gin-gonic/gin"
)

// IEmailController			定义了邮件类接口
type IEmailController interface {
	Send(ctx *gin.Context)    // 发送通知邮件
	Receive(ctx *gin.Context) // 发送反馈邮件
}

// EmailController			定义了邮件工具类
type EmailController struct {
}

// @title    Send
// @description   发送邮件
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (e EmailController) Send(ctx *gin.Context) {

	// TODO 查找对应文本
	var requestTest vo.TextRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestTest); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}
	// TODO 指定邮箱
	email := ctx.Params.ByName("id")

	err := util.SendEmail([]string{email}, requestTest.Text)

	// TODO 返回结果
	response.Success(ctx, nil, err.Error())

}

// @title    Receive
// @description   接收反馈邮件
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (e EmailController) Receive(ctx *gin.Context) {

	// TODO 查找对应文本
	var requestTest vo.TextRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestTest); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	err := util.SendEmail([]string{"20jzhang@stu.edu.cn"}, requestTest.Text)

	// TODO 返回结果
	response.Success(ctx, nil, err.Error())

}

// @title    NewEmailController
// @description   新建一个IEmailController
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   IEmailController		返回一个IEmailController用于调用各种函数
func NewEmailController() IEmailController {
	return EmailController{}
}
