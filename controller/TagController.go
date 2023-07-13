// @Title  TagController
// @Description  该文件提供关于操作标签的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"
	"MGA_OJ/util"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ITagController			定义了自动标签类接口
type ITagController interface {
	Create(ctx *gin.Context)   // 创建自动标签
	Delete(ctx *gin.Context)   // 删除自动标签
	Show(ctx *gin.Context)     // 查看自动标签
	PageList(ctx *gin.Context) // 查看自动标签列表
	Auto(ctx *gin.Context)     // 生成自动标签
}

// TagController			定义了标签工具类
type TagController struct {
	DB *gorm.DB // 含有一个数据库指针
}

// @title    Create
// @description   创建一个标签
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TagController) Create(ctx *gin.Context) {
	// TODO 查找对应标签
	tag := ctx.Params.ByName("tag")

	T := model.Tag{
		Tag: tag,
	}
	// TODO 插入数据
	if err := t.DB.Create(&T).Error; err != nil {
		response.Fail(ctx, nil, "标签上传出错，数据验证有误")
		return
	}

	// TODO 成功
	response.Success(ctx, gin.H{"tag": T}, "创建成功")
}

// @title    Show
// @description   查看一篇标签的内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TagController) Show(ctx *gin.Context) {

	// TODO 查找对应标签
	tag := ctx.Params.ByName("tag")

	var T model.Tag

	// TODO 查看标签是否在数据库中存在
	if t.DB.Where("tag = (?)", tag).First(&T).Error != nil {
		response.Fail(ctx, nil, "标签不存在")
		return
	}

	// TODO 成功
	response.Success(ctx, gin.H{"tag": tag}, "查看成功")
}

// @title    Delete
// @description   删除一篇标签的内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TagController) Delete(ctx *gin.Context) {

	// TODO 查找对应标签
	tag := ctx.Params.ByName("tag")

	var T model.Tag

	if t.DB.Where("tag = (?)", tag).First(&T).Error != nil {
		response.Fail(ctx, nil, "标签不存在")
		return
	}

	// TODO 删除标签内容
	t.DB.Delete(&T)

	// TODO 成功
	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description   查看一页标签的内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TagController) PageList(ctx *gin.Context) {

	var tags []model.Tag

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	var total int64

	// TODO 查找所有分页中可见的条目
	t.DB.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&tags)

	t.DB.Model(model.Tag{}).Count(&total)

	// TODO 成功
	response.Success(ctx, gin.H{"tags": tags}, "查看成功")
}

// @title    Auto
// @description   生成自动标签
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TagController) Auto(ctx *gin.Context) {

	text := ctx.Query("text")

	searchResult, err := util.Search(text)

	if err != nil {
		response.Fail(ctx, nil, "搜索引擎出错："+err.Error())
		return
	}

	var tags []model.Tag
	t.DB.Find(&tags)

	var ts []string
	for _, tag := range tags {
		ts = append(ts, tag.Tag)
	}

	tagCount := util.CountTags(searchResult, ts...)
	// TODO 成功
	response.Success(ctx, gin.H{"tagCount": tagCount}, "生成成功")
}

// @title    NewTagController
// @description   新建一个ITagController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   ITagController		返回一个ITagController用于调用各种函数
func NewTagController() ITagController {
	db := common.GetDB()
	db.AutoMigrate(model.Tag{})
	return TagController{DB: db}
}
