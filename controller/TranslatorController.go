// @Title  TranslatorController
// @Description  该文件提供关于操作标签的各种方法
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

// ITranslatorController			定义了翻译类接口
type ITranslatorController interface {
	Translate(ctx *gin.Context) // 翻译
}

// TranslatorController			定义了翻译工具类
type TranslatorController struct {
}

// @title    Translate
// @description   翻译
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TranslatorController) Translate(ctx *gin.Context) {
	// TODO 查找对应标签
	var requestTranslation vo.TranslationRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestTranslation); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	text, err := util.Translator(requestTranslation.Text)

	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}

	// TODO 成功
	response.Success(ctx, gin.H{"text": text}, "翻译成功")
}

// @title    NewTranslatorController
// @description   新建一个ITranslatorController
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   ITranslatorController		返回一个ITranslatorController用于调用各种函数
func NewTranslatorController() ITranslatorController {
	return TranslatorController{}
}
