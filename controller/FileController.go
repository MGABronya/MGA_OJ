// @Title  FileController
// @Description  该文件提供关于操作文件的各种方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package controller

import (
	"MGA_OJ/model"
	"MGA_OJ/response"
	"MGA_OJ/util"
	"MGA_OJ/vo"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// IFileController			定义了文件类接口
type IFileController interface {
	Upload(ctx *gin.Context)   // 上传文件
	Download(ctx *gin.Context) // 下载文件
	Unzip(ctx *gin.Context)    // 解压文件
	ShowPath(ctx *gin.Context) // 指定所有文件下的指定文件以及分目录
	MkDir(ctx *gin.Context)    // 创建目录
	CP(ctx *gin.Context)       // 将目标文件复制到指定位置
	RM(ctx *gin.Context)       // 将目标文件删除
	Rename(ctx *gin.Context)   // 将目标文件重命名
	CPAll(ctx *gin.Context)    // 将目标目录下的所有文件复制到指定位置
	RMAll(ctx *gin.Context)    // 将目标目录下的所有文件删除
}

// FileController			定义了文件工具类
type FileController struct {
}

// @title    Upload
// @description   上传文件
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (f FileController) Upload(ctx *gin.Context) {
	// 获取path中的id
	filePath := ctx.Params.ByName("path")
	file, err := ctx.FormFile("file")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有权使用文件系统
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	//TODO 数据验证
	if err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 将文件存入本地
	ctx.SaveUploadedFile(file, "./file"+filePath+file.Filename)

	response.Success(ctx, gin.H{"file": file.Filename}, "上传成功")
}

// @title    Download
// @description   下载文件
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (f FileController) Download(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有权使用文件系统
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}
	// 获取path中的id
	id := ctx.Params.ByName("id")
	fileName := id

	filePath := "./file" + fileName
	ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	ctx.File(filePath)
}

// @title    Unzip
// @description   解压
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (f FileController) Unzip(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有权使用文件系统
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	var pairString vo.PairString
	// TODO 数据验证
	if err := ctx.ShouldBind(&pairString); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	err := util.Unzip("./file"+pairString.First, "./file"+pairString.Second)
	if err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "Failed to unzip file:"+err.Error())
		return
	}

	response.Success(ctx, nil, "解压成功")
}

// @title    ShowPath
// @description   查看路径下所有文件
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (f FileController) ShowPath(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有权使用文件系统
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	// 获取path中的id
	id := ctx.Params.ByName("id")
	filePath := "./file" + id

	// TODO 获得hour目录下的所有文件
	files, err := util.GetFiles(filePath)

	if err != nil {
		response.Fail(ctx, nil, "不存在该文件夹")
		return
	}

	response.Success(ctx, gin.H{"files": files}, "请求成功")
}

// @title    MkDir
// @description   在指定目录下创建子目录
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (f FileController) MkDir(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有权使用文件系统
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}

	// 获取path中的id
	id := ctx.Params.ByName("id")
	filePath := "./file" + id

	// TODO 获得目录下的所有文件
	err := util.Mkdir(filePath)

	if err != nil {
		response.Fail(ctx, nil, "路径不存在")
		return
	}

	response.Success(ctx, nil, "创建成功")
}

// @title    RM
// @description   删除指定文件
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (f FileController) RM(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有权使用文件系统
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}
	// 获取path中的id
	id := ctx.Params.ByName("id")
	filePath := "./file" + id
	// 调用 removeFile 函数删除文件
	err := util.RemoveFile(filePath)
	if err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "Failed to remove file:"+err.Error())
		return
	}
	response.Success(ctx, nil, "File removed successfully!")
}

// @title    CP
// @description   复制指定文件到指定位置
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (f FileController) CP(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有权使用文件系统
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}
	var pairString vo.PairString
	// TODO 数据验证
	if err := ctx.ShouldBind(&pairString); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 调用moveFile函数实现文件移动
	err := util.CopyFile("./file"+pairString.First, "./file"+pairString.Second)
	if err != nil {
		response.Fail(ctx, nil, "Failed to copy file:"+err.Error())
		return
	}

	response.Success(ctx, nil, "复制成功")
}

// @title    Rename
// @description   重命名指定文件
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (f FileController) Rename(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有权使用文件系统
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}
	var pairString vo.PairString
	// TODO 数据验证
	if err := ctx.ShouldBind(&pairString); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 调用 renameFile 函数重命名文件
	err := util.RenameFile("./file"+pairString.First, pairString.Second)
	if err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "Failed to rename file:"+err.Error())
		return
	}

	response.Success(ctx, nil, "重命名成功")
}

// @title    RMAll
// @description   删除目录
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (f FileController) RMAll(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有权使用文件系统
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}
	// 获取path中的id
	id := ctx.Params.ByName("id")
	filePath := "./file" + id

	err := util.RemoveDir(filePath)
	if err != nil {
		response.Fail(ctx, nil, "Failed to delete directory:"+err.Error())
		return
	}

	response.Success(ctx, nil, "删除成功")
}

// @title    CPAll
// @description   复制目录
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (f FileController) CPAll(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有权使用文件系统
	if user.Level < 4 {
		response.Fail(ctx, nil, "权限不足")
		return
	}
	var pairString vo.PairString
	// TODO 数据验证
	if err := ctx.ShouldBind(&pairString); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}
	// TODO 调用 copyDir 函数复制目录及其内容
	err := util.CopyDir("./file"+pairString.First, pairString.Second)
	if err != nil {
		response.Fail(ctx, nil, "Failed to copy directory:"+err.Error())
		return
	}

	response.Success(ctx, nil, "复制成功")
}

// @title    NewFileController
// @description   新建一个IFileController
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   IFileController		返回一个IFileController用于调用各种函数
func NewFileController() IFileController {
	return FileController{}
}
