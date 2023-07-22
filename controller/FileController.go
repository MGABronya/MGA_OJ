// @Title  FileController
// @Description  该文件提供关于操作文件的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"MGA_OJ/response"
	"fmt"
	"log"
	"path"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// IFileController			定义了文件类接口
type IFileController interface {
	Upload(ctx *gin.Context)   // 上传文件
	Download(ctx *gin.Context) // 下载文件
}

// FileController			定义了文件工具类
type FileController struct {
}

// @title    Upload
// @description   上传文件
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (f FileController) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")

	//TODO 数据验证
	if err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	extName := path.Ext(file.Filename)

	file.Filename = uuid.NewV4().String() + extName

	// TODO 将文件存入本地
	ctx.SaveUploadedFile(file, file.Filename)

	response.Success(ctx, gin.H{"file": file.Filename}, "上传成功")
}

// @title    Download
// @description   下载文件
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (f FileController) Download(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")
	fileName := id

	filePath := "./" + fileName
	ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	ctx.File(filePath)
}

// @title    NewFileController
// @description   新建一个IFileController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   IFileController		返回一个IFileController用于调用各种函数
func NewFileController() IFileController {
	return FileController{}
}
