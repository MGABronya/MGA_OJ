// @Title  GroupController
// @Description  该文件提供关于操作用户组的各种方法
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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// IGroupController			定义了用户组类接口
type IGroupController interface {
	Interface.RestInterface      // 包含增删查改功能
	Interface.ApplyInterface     // 包含请求相关功能
	Interface.LikeInterface      // 包含点赞功能
	Interface.CollectInterface   // 包含收藏功能
	Interface.LabelInterface     // 包含标签功能
	Interface.SearchInterface    // 包含搜索功能
	LeaderList(ctx *gin.Context) // 查询指定用户领导的用户组
	MemberList(ctx *gin.Context) // 查询指定用户参加的用户组
	UserList(ctx *gin.Context)   // 查询指定用户组的用户列表
}

// GroupController			定义了用户组工具类
type GroupController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Create
// @description   创建一篇用户组
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) Create(ctx *gin.Context) {
	var requestGroup vo.GroupRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestGroup); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 创建用户组
	group := model.Group{
		Title:    requestGroup.Title,
		Content:  requestGroup.Content,
		Reslong:  requestGroup.Reslong,
		Resshort: requestGroup.Resshort,
		LeaderId: user.ID,
		Auto:     requestGroup.Auto,
	}

	// TODO 插入数据
	if err := g.DB.Create(&group).Error; err != nil {
		response.Fail(ctx, nil, "用户组上传出错，数据验证有误")
		return
	}

	// TODO 当用户权限大于2时，才能直接拉人进组
	if user.Level >= 2 {
		for _, v := range requestGroup.Users {
			userList := model.UserList{
				GroupId: group.ID,
				UserId:  v,
			}
			if g.DB.Create(&userList).Error != nil {
				response.Fail(ctx, nil, "用户上传出错，数据验证有误")
				return
			}
		}
	}

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    Update
// @description   更新一篇用户组的内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) Update(ctx *gin.Context) {
	var requestGroup vo.GroupRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestGroup); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应用户组
	id := ctx.Params.ByName("id")

	var group model.Group

	if g.DB.Where("id = ?", id).First(&group) != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}

	// TODO 查看是否是用户作者
	if user.ID != group.LeaderId {
		response.Fail(ctx, nil, "不是用户组作者，无法修改用户组")
		return
	}

	// TODO 更新用户组内容
	g.DB.Table("groups").Where("id = ?", id).Updates(requestGroup)

	// TODO 移除损坏数据
	g.Redis.HDel(ctx, "Group", id)

	if len(requestGroup.Users) != 0 {
		// TODO 查看新的用户组是否为超过最大长度

		g.DB.Where("group_id = ?", id).Delete(&model.UserList{})
		// TODO 插入相关用户
		for _, v := range requestGroup.Users {
			if ok, err := CanAddUser(v, group.ID); !ok || err != nil {
				response.Fail(ctx, nil, "用户无法添加")
				return
			}
			userList := model.UserList{
				GroupId: group.ID,
				UserId:  v,
			}
			if g.DB.Create(&userList).Error != nil {
				response.Fail(ctx, nil, "用户上传出错，数据验证有误")
				return
			}
		}
	}

	// TODO 成功
	response.Success(ctx, nil, "更新成功")
}

// @title    Show
// @description   查看一篇用户组的内容
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var group model.Group

	// TODO 先看redis中是否存在
	if ok, _ := g.Redis.HExists(ctx, "Group", id).Result(); ok {
		cate, _ := g.Redis.HGet(ctx, "Group", id).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			response.Success(ctx, gin.H{"group": group}, "成功")
			return
		} else {
			// TODO 移除损坏数据
			g.Redis.HDel(ctx, "Group", id)
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if g.DB.Where("id = ?", id).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}

	response.Success(ctx, gin.H{"group": group}, "成功")

	// TODO 将用户组存入redis供下次使用
	v, _ := json.Marshal(group)
	g.Redis.HSet(ctx, "Group", id, v)
}

// @title    Delete
// @description   删除一篇用户组
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) Delete(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	var group model.Group

	// TODO 查看用户组是否存在
	if g.DB.Where("id = ?", id).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}

	// TODO 判断当前用户是否为用户组的作者
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有操作用户组的权力
	if user.ID != group.LeaderId && user.Level < 4 {
		response.Fail(ctx, nil, "用户组不属于您，请勿非法操作")
		return
	}

	// TODO 删除用户组
	g.DB.Delete(&group)

	// TODO 移除损坏数据
	g.Redis.HDel(ctx, "Group", id)

	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description   获取多篇用户组
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) PageList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var groups []model.Group

	// TODO 查找所有分页中可见的条目
	g.DB.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&groups)

	var total int64
	g.DB.Model(model.Group{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"groups": groups, "total": total}, "成功")
}

// @title    LeaderList
// @description   获取指定用户的用户组
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) LeaderList(ctx *gin.Context) {

	// TODO 获取指定用户用户
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var groups []model.Group

	// TODO 查找所有分页中可见的条目
	g.DB.Where("leader_id = ?", id).Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&groups)

	var total int64
	g.DB.Where("leader_id = ?", id).Model(model.Group{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"groups": groups, "total": total}, "成功")
}

// @title    MemberList
// @description   获取指定用户参加的用户组
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) MemberList(ctx *gin.Context) {

	// TODO 获取指定用户用户
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var userList []model.UserList

	// TODO 查找所有分页中可见的条目
	g.DB.Where("user_id = ?", id).Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&userList)

	var total int64
	g.DB.Where("user_id = ?", id).Model(model.UserList{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"userList": userList, "total": total}, "成功")
}

// @title    UserList
// @description   获取指定用户组的用户
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) UserList(ctx *gin.Context) {

	// TODO 获取指定用户组
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var userList []model.UserList

	// TODO 查找所有分页中可见的条目
	g.DB.Where("group_id = ?", id).Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&userList)

	var total int64
	g.DB.Where("group_id = ?", id).Model(model.UserList{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"userList": userList, "total": total}, "成功")
}

// @title    Like
// @description   点赞或点踩
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) Like(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	var group model.Group

	// TODO 先看redis中是否存在
	if ok, _ := g.Redis.HExists(ctx, "Group", id).Result(); ok {
		cate, _ := g.Redis.HGet(ctx, "Group", id).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			g.Redis.HDel(ctx, "Group", id)
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if g.DB.Where("id = ?", id).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		g.Redis.HSet(ctx, "Group", id, v)
	}
leep:
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有点赞或者点踩
	if g.DB.Where("user_id = ? and group_id = ?", user.ID, id).Update("like", like).Error != nil {
		// TODO 插入数据
		groupLike := model.GroupLike{
			GroupId: group.ID,
			UserId:  user.ID,
			Like:    like,
		}
		if err := g.DB.Create(&groupLike).Error; err != nil {
			response.Fail(ctx, nil, "点赞出错，数据库存储错误")
			return
		}
	}

	response.Success(ctx, nil, "点赞成功")
}

// @title    CancelLike
// @description   取消点赞或者点踩
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) CancelLike(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 取消点赞或者点踩
	g.DB.Where("user_id = ? and group_id = ?", user.ID, id).Delete(&model.GroupLike{})

	response.Success(ctx, nil, "取消成功")
}

// @title    LikeNumber
// @description   点赞或点踩的数量
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) LikeNumber(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	var total int64

	// TODO 查看点赞或者点踩的数量
	g.DB.Where("group_id = ? and like = ?", id, like).Count(&total)

	response.Success(ctx, gin.H{"total": total}, "查看成功")
}

// @title    LikeList
// @description   点赞或点踩的列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) LikeList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var groupLikes []model.GroupLike

	var total int64

	// TODO 查看点赞或者点踩的数量
	g.DB.Where("group_id = ? and like = ?", id, like).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&groupLikes).Count(&total)

	response.Success(ctx, gin.H{"groupLikes": groupLikes, "total": total}, "查看成功")
}

// @title    LikeShow
// @description   查看用户点赞状态
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) LikeShow(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var groupLike model.GroupLike

	// TODO 查看点赞状态
	if g.DB.Where("user_id = ? and group_id = ?", user.ID, id).First(&groupLike).Error != nil {
		response.Success(ctx, gin.H{"like": 0}, "暂无状态")
		return
	}

	if groupLike.Like {
		response.Success(ctx, gin.H{"like": 1}, "已点赞")
	} else {
		response.Success(ctx, gin.H{"like": -1}, "已点踩")
	}

}

// @title    Likes
// @description   查看用户点赞状态
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) Likes(ctx *gin.Context) {

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取指定用户用户
	id := ctx.Params.ByName("id")

	// TODO 分页
	var groupLikes []model.GroupLike

	var total int64

	// TODO 查看点赞或者点踩的数量
	g.DB.Where("user_id = ? and like = ?", id, like).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&groupLikes).Count(&total)

	response.Success(ctx, gin.H{"groupLikes": groupLikes, "total": total}, "查看成功")
}

// @title    Collect
// @description   收藏
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) Collect(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	var group model.Group

	// TODO 查看用户组是否存在
	// TODO 先看redis中是否存在
	if ok, _ := g.Redis.HExists(ctx, "Group", id).Result(); ok {
		cate, _ := g.Redis.HGet(ctx, "Group", id).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			g.Redis.HDel(ctx, "Group", id)
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if g.DB.Where("id = ?", id).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		g.Redis.HSet(ctx, "Group", id, v)
	}
leep:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有收藏
	if g.DB.Where("user_id = ? and group_id = ?", user.ID, group.ID).First(&model.GroupCollect{}).Error != nil {
		groupCollect := model.GroupCollect{
			GroupId: group.ID,
			UserId:  user.ID,
		}
		// TODO 插入数据
		if err := g.DB.Create(&groupCollect).Error; err != nil {
			response.Fail(ctx, nil, "收藏出错，数据库存储错误")
			return
		}
	} else {
		response.Fail(ctx, nil, "已收藏")
		return
	}

	response.Success(ctx, nil, "收藏成功")
}

// @title    CancelCollect
// @description   取消收藏
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) CancelCollect(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有收藏
	if g.DB.Where("user_id = ? and group_id = ?", user.ID, id).First(&model.GroupCollect{}).Error != nil {
		response.Fail(ctx, nil, "未收藏")
		return
	} else {
		g.DB.Where("user_id = ? and group_id = ?", user.ID, id).Delete(&model.GroupCollect{})
		response.Success(ctx, nil, "取消收藏成功")
		return
	}
}

// @title    CollectShow
// @description   查看收藏状态
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) CollectShow(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有收藏
	if g.DB.Where("user_id = ? and group_id = ?", user.ID, id).First(&model.GroupCollect{}).Error != nil {
		response.Success(ctx, gin.H{"collect": false}, "未收藏")
		return
	} else {
		response.Success(ctx, gin.H{"collect": true}, "已收藏")
		return
	}
}

// @title    CollectList
// @description   查看收藏用户列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) CollectList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var groupCollects []model.GroupCollect

	var total int64

	// TODO 查看收藏的数量
	g.DB.Where("group_id = ?", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&groupCollects).Count(&total)

	response.Success(ctx, gin.H{"groupCollects": groupCollects, "total": total}, "查看成功")
}

// @title    CollectNumber
// @description   查看收藏用户数量
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) CollectNumber(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	var total int64

	// TODO 查看收藏的数量
	g.DB.Where("group_id = ?", id).Count(&total)

	response.Success(ctx, gin.H{"total": total}, "查看成功")
}

// @title    Collects
// @description   查看用户收藏夹
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) Collects(ctx *gin.Context) {

	// TODO 获取指定用户用户
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var groupCollects []model.GroupCollect

	var total int64

	// TODO 查看收藏的数量
	g.DB.Where("user_id = ?", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&groupCollects).Count(&total)

	response.Success(ctx, gin.H{"groupCollects": groupCollects, "total": total}, "查看成功")
}

// @title    Apply
// @description   用户申请加入某个用户组
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) Apply(ctx *gin.Context) {
	var requestGroupApply vo.GroupApplyRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestGroupApply); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取指定用户组
	id := ctx.Params.ByName("id")

	var group model.Group

	// TODO 先看redis中是否存在
	if ok, _ := g.Redis.HExists(ctx, "Group", id).Result(); ok {
		cate, _ := g.Redis.HGet(ctx, "Group", id).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			g.Redis.HDel(ctx, "Group", id)
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if g.DB.Where("id = ?", id).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		g.Redis.HSet(ctx, "Group", id, v)
	}
leep:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看用户是否已经加入用户组
	if g.DB.Where("user_id = ? and group_id = ?", user.ID, group.ID).First(&model.UserList{}).Error == nil {
		response.Fail(ctx, nil, "已加入用户组")
		return
	}

	// TODO 查看用户是否被拉黑
	if g.DB.Where("user_id = ? and group_id = ?", user.ID, group.ID).First(&model.GroupBlock{}).Error != nil {
		response.Fail(ctx, nil, "用户已被拉黑")
		return
	}

	// TODO 查看用户组是否自动通过申请
	if group.Auto {
		if ok, err := CanAddUser(user.ID, group.ID); !ok || err != nil {
			response.Fail(ctx, nil, "用户无法添加")
			return
		}
		userList := model.UserList{
			GroupId: group.ID,
			UserId:  user.ID,
		}
		// TODO 插入数据
		if err := g.DB.Create(&userList).Error; err != nil {
			response.Fail(ctx, nil, "通过申请出错，数据验证有误")
			return
		}
		// TODO 成功
		response.Success(ctx, nil, "创建成功")
		return
	}

	var groupApply model.GroupApply

	// TODO 查看用户是否已经发送过申请
	if g.DB.Where("user_id = ? and group_id = ?", user.ID, group.ID).First(&groupApply).Error == nil && groupApply.Condition {
		response.Fail(ctx, nil, "已发送过申请")
		return
	}

	// TODO 创建用户组
	groupApply = model.GroupApply{
		UserId:    user.ID,
		GroupId:   group.ID,
		Condition: true,
		Content:   requestGroupApply.Content,
		Reslong:   requestGroupApply.Reslong,
		Resshort:  requestGroupApply.Resshort,
	}

	// TODO 插入数据
	if err := g.DB.Create(&groupApply).Error; err != nil {
		response.Fail(ctx, nil, "申请出错，数据验证有误")
		return
	}

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    Consent
// @description   通过申请
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) Consent(ctx *gin.Context) {

	// TODO 获取指定申请
	id := ctx.Params.ByName("id")

	var groupApply model.GroupApply

	// TODO 查看申请是否存在
	if g.DB.Where("id = ?", id).First(&groupApply).Error != nil {
		response.Fail(ctx, nil, "申请不存在")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看该申请状态
	if !groupApply.Condition {
		response.Fail(ctx, nil, "该申请已拒绝")
		return
	}

	// TODO 查看用户是否已经加入用户组
	if g.DB.Where("user_id = ? and group_id = ?", groupApply.UserId, groupApply.GroupId).First(&model.UserList{}).Error == nil {
		response.Fail(ctx, nil, "已加入用户组")
		return
	}

	// TODO 查看当前用户是否为用户组组长
	if g.DB.Where("id = ? and leader_id = ?", groupApply.GroupId, user.ID).First(&model.Group{}).Error != nil {
		response.Fail(ctx, nil, "非用户组组长，无法操作")
		return
	}

	// TODO 查看用户是否可以添加
	if ok, err := CanAddUser(groupApply.UserId, groupApply.GroupId); !ok || err != nil {
		response.Fail(ctx, nil, "用户无法添加")
		return
	}

	userList := model.UserList{
		GroupId: groupApply.GroupId,
		UserId:  groupApply.UserId,
	}

	// TODO 插入数据
	if err := g.DB.Create(&userList).Error; err != nil {
		response.Fail(ctx, nil, "通过申请出错，数据验证有误")
		return
	}

	// TODO 删除申请
	g.DB.Delete(&groupApply)

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    Refuse
// @description   拒绝申请
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) Refuse(ctx *gin.Context) {

	// TODO 获取指定申请
	id := ctx.Params.ByName("id")

	var groupApply model.GroupApply

	// TODO 查看申请是否存在
	if g.DB.Where("id = ?", id).First(&groupApply).Error != nil {
		response.Fail(ctx, nil, "申请不存在")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看当前用户是否为用户组组长
	if g.DB.Where("id = ? and leader_id = ?", groupApply.GroupId, user.ID).First(&model.Group{}).Error != nil {
		response.Fail(ctx, nil, "非用户组组长，无法操作")
		return
	}

	groupApply.Condition = false

	// TODO 保存更改
	g.DB.Save(&groupApply)

	// TODO 成功
	response.Success(ctx, nil, "拒绝成功")
}

// @title    Block
// @description   拉黑某用户
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) Block(ctx *gin.Context) {

	// TODO 获取指定用户
	user_id := ctx.Params.ByName("user")

	var usera model.User

	// TODO 先看redis中是否存在
	if ok, _ := g.Redis.HExists(ctx, "User", user_id).Result(); ok {
		cate, _ := g.Redis.HGet(ctx, "User", user_id).Result()
		if json.Unmarshal([]byte(cate), &usera) == nil {
			goto leap
		} else {
			// TODO 移除损坏数据
			g.Redis.HDel(ctx, "User", user_id)
		}
	}

	// TODO 查看用户是否在数据库中存在
	if g.DB.Where("id = ?", user_id).First(&usera).Error != nil {
		response.Fail(ctx, nil, "用户不存在")
		return
	}
	{
		// TODO 将用户存入redis供下次使用
		v, _ := json.Marshal(usera)
		g.Redis.HSet(ctx, "User", user_id, v)
	}
leap:
	// TODO 获取指定用户组
	group_id := ctx.Params.ByName("group")

	var group model.Group

	// TODO 先看redis中是否存在
	if ok, _ := g.Redis.HExists(ctx, "Group", group_id).Result(); ok {
		cate, _ := g.Redis.HGet(ctx, "Group", group_id).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			g.Redis.HDel(ctx, "Group", group_id)
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if g.DB.Where("id = ?", group_id).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		g.Redis.HSet(ctx, "Group", group_id, v)
	}
leep:
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看当前用户是否为用户组组长
	if group.LeaderId != user.ID {
		response.Fail(ctx, nil, "非用户组组长，无法操作")
		return
	}

	// TODO 查看当前用户是否已经拉黑
	if g.DB.Where("user_id = ? and group_id = ?", user_id, group_id).First(&model.GroupBlock{}).Error == nil {
		response.Fail(ctx, nil, "用户已拉黑")
		return
	}

	// TODO 将指定用户放入黑名单
	groupBlock := model.GroupBlock{
		UserId:  usera.ID,
		GroupId: group.ID,
	}

	if g.DB.Create(&groupBlock).Error != nil {
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
func (g GroupController) RemoveBlack(ctx *gin.Context) {

	// TODO 获取指定用户
	user_id := ctx.Params.ByName("user")

	// TODO 获取指定用户组
	group_id := ctx.Params.ByName("group")

	var group model.Group

	// TODO 查看用户组是否存在
	// TODO 先看redis中是否存在
	if ok, _ := g.Redis.HExists(ctx, "Group", group_id).Result(); ok {
		cate, _ := g.Redis.HGet(ctx, "Group", group_id).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			g.Redis.HDel(ctx, "Group", group_id)
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if g.DB.Where("id = ?", group_id).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		g.Redis.HSet(ctx, "Group", group_id, v)
	}
leep:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看当前用户是否为用户组组长
	if group.LeaderId != user.ID {
		response.Fail(ctx, nil, "非用户组组长，无法操作")
		return
	}

	var groupBlock model.GroupBlock

	// TODO 查看当前用户是否已经拉黑
	if g.DB.Where("user_id = ? and group_id = ?", user_id, group_id).First(&groupBlock).Error != nil {
		response.Fail(ctx, nil, "用户未被拉黑")
	}

	// TODO 将用户移除黑名单
	g.DB.Delete(&groupBlock)

	// TODO 成功
	response.Success(ctx, nil, "移除黑名单成功")
}

// @title    BlackList
// @description   查看黑名单
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) BlackList(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取指定用户组
	id := ctx.Params.ByName("id")

	var group model.Group

	// TODO 查看用户组是否存在
	// TODO 先看redis中是否存在
	if ok, _ := g.Redis.HExists(ctx, "Group", id).Result(); ok {
		cate, _ := g.Redis.HGet(ctx, "Group", id).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			g.Redis.HDel(ctx, "Group", id)
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if g.DB.Where("id = ?", id).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		g.Redis.HSet(ctx, "Group", id, v)
	}
leep:

	// TODO 查看当前用户是否为用户组组长
	if group.LeaderId != user.ID {
		response.Fail(ctx, nil, "非用户组组长，无法操作")
		return
	}

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var groupBlocks []model.GroupBlock

	var total int64

	// TODO 查看黑名单
	g.DB.Where("group_id = ?", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&groupBlocks).Count(&total)

	response.Success(ctx, gin.H{"groupBlocks": groupBlocks, "total": total}, "查看成功")
}

// @title    ApplyingList
// @description   查看用户申请列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) ApplyingList(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var groupApplys []model.GroupApply

	var total int64

	// TODO 查看申请的数量
	g.DB.Where("user_id = ?", user.ID).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&groupApplys).Count(&total)

	response.Success(ctx, gin.H{"groupApplys": groupApplys, "total": total}, "查看成功")
}

// @title    AppliedList
// @description   查看用户申请列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) AppliedList(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 取出group
	id := ctx.Params.ByName("id")

	// TODO 查找group
	var group model.Group
	// TODO 先看redis中是否存在
	if ok, _ := g.Redis.HExists(ctx, "Group", id).Result(); ok {
		cate, _ := g.Redis.HGet(ctx, "Group", id).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			g.Redis.HDel(ctx, "Group", id)
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if g.DB.Where("id = ?", id).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		g.Redis.HSet(ctx, "Group", id, v)
	}
leep:

	// TODO 查看是否为用户组组长
	if group.LeaderId != user.ID {
		response.Fail(ctx, nil, "非用户组组长，无法操作")
		return
	}

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var groupApplys []model.GroupApply

	var total int64

	// TODO 查看申请的数量
	g.DB.Where("group_id = ?", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&groupApplys).Count(&total)

	response.Success(ctx, gin.H{"groupApplys": groupApplys, "total": total}, "查看成功")
}

// @title    Quit
// @description   退出用户组
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) Quit(ctx *gin.Context) {

	// TODO 获取指定用户组
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看用户是否已经加入用户组
	if g.DB.Where("user_id = ? and group_id = ?", user.ID, id).First(&model.UserList{}).Error != nil {
		response.Fail(ctx, nil, "未加入用户组")
		return
	}

	// TODO 在用户中删除用户
	g.DB.Where("user_id = ? and group_id = ?", user.ID, id).Delete(&model.UserList{})

	// TODO 成功
	response.Success(ctx, nil, "退出成功")
}

// @title    LabelCreate
// @description   标签创建
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) LabelCreate(ctx *gin.Context) {
	// TODO 获取指定用户组
	id := ctx.Params.ByName("id")

	// TODO 获取标签
	label := ctx.Params.ByName("label")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看用户组是否存在
	var group model.Group

	// TODO 先尝试在redis中寻找
	if ok, _ := g.Redis.HExists(ctx, "Group", id).Result(); ok {
		art, _ := g.Redis.HGet(ctx, "Group", id).Result()
		if json.Unmarshal([]byte(art), &group) == nil {
			goto leep
		} else {
			// TODO 解码失败，删除字段
			g.Redis.HDel(ctx, "Group", id)
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if g.DB.Where("id = ?", id).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		g.Redis.HSet(ctx, "Group", id, v)
	}
leep:

	// TODO 查看是否为用户组作者
	if group.LeaderId != user.ID {
		response.Fail(ctx, nil, "不是用户组作者，请勿非法操作")
		return
	}

	// TODO 创建标签
	groupLabel := model.GroupLabel{
		Label:   label,
		GroupId: group.ID,
	}

	// TODO 插入数据
	if err := g.DB.Create(&groupLabel).Error; err != nil {
		response.Fail(ctx, nil, "用户组标签上传出错，数据验证有误")
		return
	}

	// TODO 解码失败，删除字段
	g.Redis.HDel(ctx, "GroupLabel", id)

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    LabelDelete
// @description   标签删除
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) LabelDelete(ctx *gin.Context) {
	// TODO 获取指定用户组
	id := ctx.Params.ByName("id")

	// TODO 获取标签
	label := ctx.Params.ByName("label")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看用户组是否存在
	var group model.Group

	// TODO 先尝试在redis中寻找
	if ok, _ := g.Redis.HExists(ctx, "Group", id).Result(); ok {
		art, _ := g.Redis.HGet(ctx, "Group", id).Result()
		if json.Unmarshal([]byte(art), &group) == nil {
			goto leep
		} else {
			// TODO 解码失败，删除字段
			g.Redis.HDel(ctx, "Group", id)
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if g.DB.Where("id = ?", id).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		g.Redis.HSet(ctx, "Group", id, v)
	}
leep:

	// TODO 查看是否为用户组作者
	if group.LeaderId != user.ID {
		response.Fail(ctx, nil, "不是用户组作者，请勿非法操作")
		return
	}

	// TODO 删除用户组标签
	if g.DB.Where("id = ?", label).First(&model.GroupLabel{}).Error != nil {
		response.Fail(ctx, nil, "标签不存在")
		return
	}

	g.DB.Where("id = ?", label).Delete(&model.GroupLabel{})

	// TODO 解码失败，删除字段
	g.Redis.HDel(ctx, "GroupLabel", id)

	// TODO 成功
	response.Success(ctx, nil, "删除成功")
}

// @title    LabelShow
// @description   标签查看
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) LabelShow(ctx *gin.Context) {
	// TODO 获取指定用户组
	id := ctx.Params.ByName("id")

	// TODO 查找数据
	var groupLabels []model.GroupLabel
	// TODO 先尝试在redis中寻找
	if ok, _ := g.Redis.HExists(ctx, "GroupLabel", id).Result(); ok {
		art, _ := g.Redis.HGet(ctx, "GroupLabel", id).Result()
		if json.Unmarshal([]byte(art), &groupLabels) == nil {
			goto leap
		} else {
			// TODO 解码失败，删除字段
			g.Redis.HDel(ctx, "GroupLabel", id)
		}
	}

	// TODO 在数据库中查找
	g.DB.Where("group_id = ?", id).Find(&groupLabels)
	{
		// TODO 将用户组标签存入redis供下次使用
		v, _ := json.Marshal(groupLabels)
		g.Redis.HSet(ctx, "GroupLabel", id, v)
	}

leap:

	// TODO 成功
	response.Success(ctx, gin.H{"groupLabels": groupLabels}, "查看成功")
}

// @title    Search
// @description   文本搜索
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) Search(ctx *gin.Context) {
	// TODO 获取文本
	text := ctx.Params.ByName("text")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	var groups []model.Group

	// TODO 模糊匹配
	g.DB.Where("match(title,content,res_long,res_short) against(? in boolean mode)", text+"*").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&groups)

	// TODO 查看查询总数
	var total int64
	g.DB.Where("match(title,content,res_long,res_short) against(? in boolean mode)", text+"*").Model(model.Group{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"groups": groups, "total": total}, "成功")
}

// @title    SearchLabel
// @description   指定标签的搜索
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) SearchLabel(ctx *gin.Context) {

	var requestLabels vo.LabelsRequest

	// TODO 获取标签
	if err := ctx.ShouldBind(&requestLabels); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 通过标签寻找
	var groupIds []struct {
		GroupId uuid.UUID `json:"group_id"` // 用户组外键
	}

	// TODO 进行标签匹配
	g.DB.Distinct("group_id").Where("label in (?)", requestLabels.Labels).Model(model.GroupLabel{}).Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&groupIds)

	// TODO 查看查询总数
	var total int64
	g.DB.Distinct("group_id").Where("label in (?)", requestLabels.Labels).Model(model.GroupLabel{}).Count(&total)

	// TODO 查找对应用户组
	var groups []model.Group

	g.DB.Where("id in (?)", groupIds).Find(&groups)

	// TODO 返回数据
	response.Success(ctx, gin.H{"groups": groups, "total": total}, "成功")
}

// @title    SearchWithLabel
// @description   指定标签与文本的搜索
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (g GroupController) SearchWithLabel(ctx *gin.Context) {

	// TODO 获取文本
	text := ctx.Params.ByName("text")

	var requestLabels vo.LabelsRequest

	// TODO 获取标签
	if err := ctx.ShouldBind(&requestLabels); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 通过标签寻找
	var groupIds []struct {
		GroupId uuid.UUID `json:"group_id"` // 用户组外键
	}

	// TODO 进行标签匹配
	g.DB.Distinct("group_id").Where("label in (?)", requestLabels.Labels).Model(model.GroupLabel{}).Find(&groupIds)

	// TODO 查找对应用户组
	var groups []model.Group

	// TODO 模糊匹配
	g.DB.Where("id in (?) and match(title,content,res_long,res_short) against(? in boolean mode)", groupIds, text+"*").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&groups)

	// TODO 查看查询总数
	var total int64
	g.DB.Where("id in (?) and match(title,content,res_long,res_short) against(? in boolean mode)", groupIds, text+"*").Model(model.Group{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"groups": groups, "total": total}, "成功")
}

// @title    NewGroupController
// @description   新建一个IGroupController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   IGroupController		返回一个IGroupController用于调用各种函数
func NewGroupController() IGroupController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.Group{})
	db.AutoMigrate(model.GroupApply{})
	db.AutoMigrate(model.GroupBlock{})
	db.AutoMigrate(model.GroupCollect{})
	db.AutoMigrate(model.GroupLike{})
	db.AutoMigrate(model.UserList{})
	db.AutoMigrate(model.GroupLabel{})
	return GroupController{DB: db, Redis: redis}
}

// @title    CanAddUser
// @description   查看用户是否可以加入用户组
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    user_id uuid.UUID, group_id uuid.UUID	表示用户和用户组
// @return   bool, error					返回用户是否可以加入用户组
func CanAddUser(user_id uuid.UUID, group_id uuid.UUID) (bool, error) {
	// TODO 连接数据库
	db := common.GetDB()

	// TODO 获取该组参加的表单列表
	var groupLists []model.GroupList
	db.Where("group_id = ?", group_id).Find(&groupLists)

	for _, groupList := range groupLists {
		set_id := groupList.SetId
		// TODO 获取该表单下的所有组
		var groupLists []model.GroupList
		db.Where("set_id = ?", set_id).Find(&groupLists)
		// TODO 获取该表单
		var set model.Set
		if err := db.Where("id = ?", set_id).First(&set).Error; err != nil {
			return false, err
		}
		// TODO 该表单对组员有合法的人数限制
		if set.PassNum != 0 {
			var total int64
			db.Where("group_id = ?", group_id).Model(&model.UserList{}).Count(&total)
			if total >= int64(set.PassNum) {
				return false, nil
			}
		}
		// TODO 该表单禁止不同组之间成员重复
		if set.PassRe {
			for _, group := range groupLists {
				var userLists []model.UserList
				db.Where("group_id = ?", group).Find(&userLists)
				for _, user := range userLists {
					if user.UserId == user_id {
						return false, nil
					}
				}
			}
		}
		// TODO 表单的比赛不能已经开始
		var competitions []model.Competition
		db.Where("set_id = ?", set_id).Find(&competitions)
		for _, competition := range competitions {
			if time.Now().After(time.Time(competition.StartTime)) && !time.Now().After(time.Time(competition.EndTime)) {
				return false, nil
			}
		}
	}

	return true, nil

}
