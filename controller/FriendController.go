// @Title  FriendController
// @Description  该文件提供关于操作好友的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"MGA_OJ/Interface"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"
	"MGA_OJ/vo"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IFriendController			定义了好友类接口
type IFriendController interface {
	Interface.ApplyInterface // 包含请求相关功能
	Interface.BlockInterface // 包含黑名单相关功能
}

// FriendController			定义了好友工具类
type FriendController struct {
	DB *gorm.DB // 含有一个数据库指针
}

// @title    Apply
// @description   用户申请添加某个好友
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (f FriendController) Apply(ctx *gin.Context) {
	var requestFriendApply vo.FriendApplyRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestFriendApply); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取指定用户
	id := ctx.Params.ByName("id")

	var friend model.User

	// TODO 查看好友是否存在
	if f.DB.Where("id = ?", id).First(&friend).Error != nil {
		response.Fail(ctx, nil, "指定用户不存在")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看用户是否已经加入好友
	if f.DB.Where("user_id = ? and friend_id = ?", user.ID, friend.ID).First(&model.Friend{}).Error == nil || f.DB.Where("user_id = ? and friend_id = ?", friend.ID, user.ID).First(&model.Friend{}).Error == nil {
		response.Fail(ctx, nil, "已加入好友")
		return
	}

	var friendApply model.FriendApply

	// TODO 查看是否重复申请
	if f.DB.Where("user_id = ? and friend_id = ?", user.ID, friend.ID).First(&friendApply).Error == nil && friendApply.Condition {
		response.Fail(ctx, nil, "已发送过申请")
		return
	}

	// TODO 查看是否被拉黑
	if f.DB.Where("user_id = ? and owner_id = ?", user.ID, friend.ID).First(&model.FriendBlock{}).Error == nil {
		response.Fail(ctx, nil, "已被拉黑")
		return
	}

	// TODO 创建好友申请
	friendApply = model.FriendApply{
		UserId:    user.ID,
		FriendId:  friend.ID,
		Content:   requestFriendApply.Content,
		ResLong:   requestFriendApply.ResLong,
		ResShort:  requestFriendApply.ResShort,
		Condition: true,
	}

	// TODO 插入数据
	if err := f.DB.Create(&friendApply).Error; err != nil {
		response.Fail(ctx, nil, "申请出错，数据验证有误")
		return
	}

	// TODO 成功
	response.Success(ctx, nil, "申请成功")
}

// @title    Consent
// @description   通过申请
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (f FriendController) Consent(ctx *gin.Context) {

	// TODO 获取指定申请
	id := ctx.Params.ByName("id")

	var friendApply model.FriendApply

	// TODO 查看申请是否存在
	if f.DB.Where("id = ?", id).First(&friendApply).Error != nil {
		response.Fail(ctx, nil, "申请不存在")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看用户是否已经加入好友
	if f.DB.Where("user_id = ? and friend_id = ?", friendApply.UserId, friendApply.FriendId).First(&model.Friend{}).Error == nil || f.DB.Where("user_id = ? and friend_id = ?", friendApply.FriendId, friendApply.UserId).First(&model.Friend{}).Error == nil {
		response.Fail(ctx, nil, "已加入好友")
		return
	}

	// TODO 查看当前用户是否为好友申请接收方
	if friendApply.FriendId != user.ID {
		response.Fail(ctx, nil, "非好友申请接收方，无法操作")
		return
	}

	friend := model.Friend{
		FriendId: friendApply.FriendId,
		UserId:   friendApply.UserId,
	}

	// TODO 插入数据
	if err := f.DB.Create(&friend).Error; err != nil {
		response.Fail(ctx, nil, "通过申请出错，数据验证有误")
		return
	}

	// TODO 删除申请
	f.DB.Delete(&friendApply)

	// TODO 成功
	response.Success(ctx, nil, "通过成功")
}

// @title    ApplyingList
// @description   查看用户发出申请列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (f FriendController) ApplyingList(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var friendApplys []model.FriendApply

	var total int64

	// TODO 查看申请的数量
	f.DB.Where("user_id = ?", user.ID).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&friendApplys)

	f.DB.Where("user_id = ?", user.ID).Model(model.FriendApply{}).Count(&total)

	response.Success(ctx, gin.H{"friendApplys": friendApplys, "total": total}, "查看成功")
}

// @title    AppliedList
// @description   查看用户接收申请列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (f FriendController) AppliedList(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var friendApplys []model.FriendApply

	var total int64

	// TODO 查看申请的数量
	f.DB.Where("friend_id = ?", user.ID).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&friendApplys)

	f.DB.Where("friend_id = ?", user.ID).Model(model.FriendApply{}).Count(&total)

	response.Success(ctx, gin.H{"friendApplys": friendApplys, "total": total}, "查看成功")
}

// @title    Quit
// @description   删除好友
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (f FriendController) Quit(ctx *gin.Context) {

	// TODO 获取指定好友
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看用户是否已经加入好友
	if f.DB.Where("user_id = ? and friend_id = ?", user.ID, id).First(&model.Friend{}).Error != nil && f.DB.Where("user_id = ? and friend_id = ?", id, user.ID).First(&model.Friend{}).Error != nil {
		response.Fail(ctx, nil, "未成为好友")
		return
	}

	// TODO 在用户中删除用户
	f.DB.Where("user_id = ? and friend_id = ?", user.ID, id).Delete(&model.Friend{})
	f.DB.Where("user_id = ? and friend_id = ?", id, user.ID).Delete(&model.Friend{})

	// TODO 成功
	response.Success(ctx, nil, "删除成功")
}

// @title    Refuse
// @description   拒绝申请
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (f FriendController) Refuse(ctx *gin.Context) {

	// TODO 获取指定申请
	id := ctx.Params.ByName("id")

	var friendApply model.FriendApply

	// TODO 查看申请是否存在
	if f.DB.Where("id = ?", id).First(&friendApply).Error != nil {
		response.Fail(ctx, nil, "申请不存在")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看当前用户是否为被请求方
	if friendApply.FriendId != user.ID {
		response.Fail(ctx, nil, "非被请求方，无法操作")
		return
	}

	friendApply.Condition = false

	// TODO 保存更改
	f.DB.Save(&friendApply)

	// TODO 成功
	response.Success(ctx, nil, "拒绝成功")
}

// @title    Block
// @description   拉黑某用户
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (f FriendController) Block(ctx *gin.Context) {

	// TODO 获取指定用户
	user_id := ctx.Params.ByName("id")

	var user model.User

	// TODO 查看用户是否存在
	if f.DB.Where("id = ?", user_id).First(&user).Error != nil {
		response.Fail(ctx, nil, "用户不存在")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	owner := tuser.(model.User)

	// TODO 查看当前用户是否已经拉黑
	if f.DB.Where("owner = ? and user_id = ?", owner.ID, user.ID).First(&model.FriendBlock{}).Error == nil {
		response.Fail(ctx, nil, "用户已拉黑")
	}

	// TODO 将指定用户放入黑名单
	friendBlock := model.FriendBlock{
		UserId:  user.ID,
		OwnerId: owner.ID,
	}

	if f.DB.Create(&friendBlock).Error != nil {
		response.Fail(ctx, nil, "黑名单入库错误")
		return
	}

	// TODO 成功
	response.Success(ctx, nil, "拉黑成功")
}

// @title    RemoveBlack
// @description   移除某用户的黑名单
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (f FriendController) RemoveBlack(ctx *gin.Context) {

	// TODO 获取指定用户
	id := ctx.Params.ByName("id")

	var user model.User

	// TODO 查看用户是否存在
	if f.DB.Where("id = ?", id).First(&user).Error != nil {
		response.Fail(ctx, nil, "用户不存在")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	owner := tuser.(model.User)

	var friendBlock model.FriendBlock

	// TODO 查看当前用户是否已经拉黑
	if f.DB.Where("user_id = ? and owner_id = ?", user.ID, owner.ID).First(&friendBlock).Error != nil {
		response.Fail(ctx, nil, "用户未被拉黑")
	}

	// TODO 将用户移除黑名单
	f.DB.Delete(&friendBlock)

	// TODO 成功
	response.Success(ctx, nil, "移除黑名单成功")
}

// @title    BlackList
// @description   查看黑名单
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (f FriendController) BlackList(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var friendBlocks []model.FriendBlock

	var total int64

	// TODO 查看黑名单
	f.DB.Where("owner_id = ?", user.ID).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&friendBlocks)

	f.DB.Where("owner_id = ?", user.ID).Model(model.FriendBlock{}).Count(&total)

	response.Success(ctx, gin.H{"friendBlocks": friendBlocks, "total": total}, "查看成功")
}

// @title    NewFriendController
// @description   新建一个IFriendController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   IFriendController		返回一个IFriendController用于调用各种函数
func NewFriendController() IFriendController {
	db := common.GetDB()
	db.AutoMigrate(model.Friend{})
	db.AutoMigrate(model.FriendApply{})
	db.AutoMigrate(model.FriendBlock{})
	return FriendController{DB: db}
}
