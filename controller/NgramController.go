// @Title  NgramController
// @Description  该文件提供关于操作文本相似的各种方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package controller

import (
	"MGA_OJ/response"
	"MGA_OJ/util"
	"MGA_OJ/vo"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

// INgramController			定义了文本相似类接口
type INgramController interface {
	ComputeSimilarity(ctx *gin.Context) // 判断文本相似度
	JudgeSimilarity(ctx *gin.Context)   // 获得文本相似度连通块
}

// NgramController			定义了文本相似工具类
type NgramController struct {
}

// @title    ComputeSimilarity
// @description  判断文本相似度
// @auth      MGAronya       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (n NgramController) ComputeSimilarity(ctx *gin.Context) {
	var textsRequest vo.TextsRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&textsRequest); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}
	var similarity [][]float64 = make([][]float64, len(textsRequest.Texts))

	for i := range similarity {
		similarity[i] = make([]float64, len(similarity))
		for j := range similarity[i] {
			if j < i {
				similarity[i][j] = util.ComputeSimilarity(textsRequest.Texts[i], textsRequest.Texts[j], 3)
			} else if j == i {
				similarity[i][j] = 1
			} else {
				similarity[i][j] = similarity[j][i]
			}
		}
	}
	// TODO 返回数据
	response.Success(ctx, gin.H{"similarity": similarity}, "成功")
}

// @title    JudgeSimilarity
// @description  获得文本相似度连通块
// @auth      MGAronya       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (n NgramController) JudgeSimilarity(ctx *gin.Context) {
	// TODO 获取judge值
	judge, err := strconv.ParseFloat(ctx.Params.ByName("judge"), 64)
	if err != nil {
		response.Fail(ctx, nil, "judge数据错误")
		return
	}
	var arrsRequest vo.ArrsRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&arrsRequest); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 查找所有的连通块
	arrs, err := util.SimilarityJudge(arrsRequest.Arrs, judge)

	if err != nil {
		response.Fail(ctx, nil, err.Error())
		return
	}

	// TODO 返回数据
	response.Success(ctx, gin.H{"arrs": arrs}, "成功")
}

// @title    NewNgramController
// @description   新建一个INgramController
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   INgramController		返回一个INgramController用于调用各种函数
func NewNgramController() INgramController {
	return NgramController{}
}
