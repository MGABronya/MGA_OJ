// @Title  NoticeBoardController
// @Description  该文件提供关于操作公告栏的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"MGA_OJ/Interface"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"
	"MGA_OJ/vo"
	"encoding/json"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

// INoticeBoardController			定义了公告栏类接口
type INoticeBoardController interface {
	Interface.RestInterface // 增删查改
}

// NoticeBoardController			定义了公告栏工具类
type NoticeBoardController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Create
// @description  发布通知
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (n NoticeBoardController) Create(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有权给公告栏添加公告
	if user.Level < 4 {
		response.Fail(ctx, nil, "无权为公告栏添加公告")
		return
	}

	// TODO 接收广播内容
	var notic vo.Notice
	// TODO 数据验证
	if err := ctx.ShouldBind(&notic); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}
	// TODO 存入数据库
	notice := model.NoticeBoard{
		UserId:   user.ID,
		Title:    notic.Title,
		Content:  notic.Content,
		ResLong:  notic.ResLong,
		ResShort: notic.ResShort,
	}
	n.DB.Save(&notice)
	// TODO 将notic打包
	v, _ := json.Marshal(notice)
	n.Redis.HSet(ctx, "NoticeBoard", notice.ID, v)
	response.Success(ctx, gin.H{"notice": notice}, "发布成功")
}

// @title    Show
// @description  查看通告
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (n NoticeBoardController) Show(ctx *gin.Context) {

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 查找公告
	var notice model.NoticeBoard

	// TODO 查看公告是否存在
	// TODO 先看redis中是否存在
	if ok, _ := n.Redis.HExists(ctx, "NoticeBoard", id).Result(); ok {
		cate, _ := n.Redis.HGet(ctx, "NoticeBoard", id).Result()
		if json.Unmarshal([]byte(cate), &notice) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			n.Redis.HDel(ctx, "NoticeBoard", id)
		}
	}

	// TODO 查看公告是否在数据库中存在
	if n.DB.Where("id = (?)", id).First(&notice).Error != nil {
		response.Fail(ctx, nil, "公告不存在")
		return
	}
	{
		// TODO 将公告存入redis供下次使用
		v, _ := json.Marshal(notice)
		n.Redis.HSet(ctx, "NoticeBoard", id, v)
	}
leep:

	response.Success(ctx, gin.H{"notice": notice}, "成功")
}

// @title    Update
// @description   更新一篇公告的内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (n NoticeBoardController) Update(ctx *gin.Context) {
	var notice vo.Notice
	// TODO 数据验证
	if err := ctx.ShouldBind(&notice); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应公告
	id := ctx.Params.ByName("id")

	var noticeBoard model.NoticeBoard

	if n.DB.Where("id = (?)", id).First(&noticeBoard).Error != nil {
		response.Fail(ctx, nil, "公告不存在")
		return
	}

	// TODO 查看是否是公告作者
	if user.ID != noticeBoard.UserId {
		response.Fail(ctx, nil, "不是公告作者，无法修改公告")
		return
	}

	// TODO 新建公告
	noticeUpdate := model.NoticeBoard{
		Title:    notice.Title,
		Content:  notice.Content,
		ResLong:  notice.ResLong,
		ResShort: notice.ResShort,
	}

	// TODO 更新文章内容
	n.DB.Model(model.NoticeBoard{}).Where("id = (?)", id).Updates(noticeUpdate)

	// TODO 解码失败，删除字段
	n.Redis.HDel(ctx, "NoticeBoard", id)

	// TODO 成功
	response.Success(ctx, nil, "更新成功")
}

// @title    Delete
// @description   删除一篇公告
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (n NoticeBoardController) Delete(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	var noticeBoard model.NoticeBoard

	// TODO 查看公告是否存在
	if n.DB.Where("id = (?)", id).First(&noticeBoard).Error != nil {
		response.Fail(ctx, nil, "公告不存在")
		return
	}

	// TODO 判断当前用户是否为公告的作者
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有操作文章的权力
	if user.ID != noticeBoard.UserId && user.Level < 4 {
		response.Fail(ctx, nil, "文章不属于您，请勿非法操作")
		return
	}

	// TODO 删除公告
	n.DB.Delete(&noticeBoard)

	// TODO 解码失败，删除字段
	n.Redis.HDel(ctx, "NoticeBoard", id)

	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description  通知列表
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (n NoticeBoardController) PageList(ctx *gin.Context) {

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	var notices []model.NoticeBoard

	n.DB.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&notices)
	var total int64

	n.DB.Model(model.Notice{}).Count(&total)

	response.Success(ctx, gin.H{"notices": notices, "total": total}, "成功")
}

// @title    NewCompetitionController
// @description   新建一个INoticeBoardController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   INoticeBoardController		返回一个INoticeBoardController用于调用各种函数
func NewNoticeBoardController() INoticeBoardController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.NoticeBoard{})
	return NoticeBoardController{DB: db, Redis: redis}
}
