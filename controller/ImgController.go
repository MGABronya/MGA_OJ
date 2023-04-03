// @Title  ImgController
// @Description  该文件提供关于操作图片的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"MGA_OJ/response"
	"log"
	"path"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// IImgController			定义了图片类接口
type IImgController interface {
	Upload(ctx *gin.Context) // 上传图片
}

// ImgController			定义了图片工具类
type ImgController struct {
}

// @title    Upload
// @description   上传一张图片
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (i ImgController) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")

	//TODO 数据验证
	if err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	extName := path.Ext(file.Filename)
	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
	}

	// TODO 格式验证
	if _, ok := allowExtMap[extName]; !ok {
		response.Fail(ctx, nil, "文件格式有误")
		return
	}

	file.Filename = uuid.NewV4().String() + extName

	// TODO 将文件存入本地
	ctx.SaveUploadedFile(file, file.Filename)

	response.Success(ctx, gin.H{"Icon": file.Filename}, "上传成功")
}

// @title    NewImgController
// @description   新建一个ImgController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   ImgController		返回一个ImgController用于调用各种函数
func NewImgController() ImgController {
	return ImgController{}
}
