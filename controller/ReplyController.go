// @Title  ReplyController
// @Description  该文件提供关于操作讨论的回复的各种方法
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

// IReplyController			定义了讨论的回复类接口
type IReplyController interface {
	Interface.RestInterface    // 包含增删查改功能
	Interface.LikeInterface    // 包含点赞功能
	Interface.HotInterface     // 包含热度功能
	UserList(ctx *gin.Context) // 查看指定用户的回复
}

// ReplyController			定义了讨论的回复工具类
type ReplyController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Create
// @description   创建一篇讨论的回复
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r ReplyController) Create(ctx *gin.Context) {
	var requestReply vo.ReplyRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestReply); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 查找对应讨论
	id := ctx.Params.ByName("id")

	var comment model.Comment

	// TODO 先看redis中是否存在
	if ok, _ := r.Redis.HExists(ctx, "Comment", id).Result(); ok {
		cate, _ := r.Redis.HGet(ctx, "Comment", id).Result()
		if json.Unmarshal([]byte(cate), &comment) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			r.Redis.HDel(ctx, "Comment", id)
		}
	}

	// TODO 查看讨论是否在数据库中存在
	if r.DB.Where("id = ?", id).First(&comment).Error != nil {
		response.Fail(ctx, nil, "讨论不存在")
		return
	}
	{
		// TODO 将讨论存入redis供下次使用
		v, _ := json.Marshal(comment)
		r.Redis.HSet(ctx, "Comment", id, v)
	}
leep:
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 创建讨论的回复
	var reply = model.Reply{
		UserId:    user.ID,
		CommentId: comment.ID,
		Content:   requestReply.Content,
		Reslong:   requestReply.Reslong,
		Resshort:  requestReply.Resshort,
	}

	// TODO 插入数据
	if err := r.DB.Create(&reply).Error; err != nil {
		response.Fail(ctx, nil, "讨论的回复上传出错，数据库存储错误")
		return
	}

	// TODO 创建热度
	r.Redis.ZAdd(ctx, "ReplyHot"+comment.ID.String(), redis.Z{Member: reply.ID.String(), Score: 100})

	// TODO 成功
	response.Success(ctx, gin.H{"reply": reply}, "创建成功")
}

// @title    Update
// @description   更新一篇讨论的回复的内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r ReplyController) Update(ctx *gin.Context) {
	var requestReply model.Reply
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestReply); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应讨论的回复
	id := ctx.Params.ByName("id")

	var reply model.Reply

	if r.DB.Where("id = ?", id).First(&reply) != nil {
		response.Fail(ctx, nil, "讨论的回复不存在")
		return
	}

	// TODO 查看是否是用户作者
	if user.ID != reply.UserId {
		response.Fail(ctx, nil, "不是讨论的回复作者，无法修改讨论的回复")
		return
	}

	// TODO 更新讨论的回复内容
	r.DB.Where("id = ?", id).Updates(requestReply)

	// TODO 移除损坏数据
	r.Redis.HDel(ctx, "Reply", id)

	// TODO 成功
	response.Success(ctx, nil, "更新成功")
}

// @title    Show
// @description   查看一篇讨论的回复的内容
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r ReplyController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var reply model.Reply

	// TODO 先看redis中是否存在
	if ok, _ := r.Redis.HExists(ctx, "Reply", id).Result(); ok {
		cate, _ := r.Redis.HGet(ctx, "Reply", id).Result()
		if json.Unmarshal([]byte(cate), &reply) == nil {
			response.Success(ctx, gin.H{"reply": reply}, "成功")
			return
		} else {
			// TODO 移除损坏数据
			r.Redis.HDel(ctx, "Reply", id)
		}
	}

	// TODO 查看讨论的回复是否在数据库中存在
	if r.DB.Where("id = ?", id).First(&reply).Error != nil {
		response.Fail(ctx, nil, "讨论的回复不存在")
		return
	}

	response.Success(ctx, gin.H{"reply": reply}, "成功")

	// TODO 将提交存入redis供下次使用
	v, _ := json.Marshal(reply)
	r.Redis.HSet(ctx, "Reply", id, v)
}

// @title    Delete
// @description   删除一篇讨论的回复
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r ReplyController) Delete(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	var reply model.Reply

	// TODO 查看讨论的回复是否存在
	if r.DB.Where("id = ?", id).First(&reply).Error != nil {
		response.Fail(ctx, nil, "讨论的回复不存在")
		return
	}

	// TODO 判断当前用户是否为讨论的回复的作者
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有操作讨论的回复的权力
	if user.ID != reply.UserId && user.Level < 4 {
		response.Fail(ctx, nil, "讨论的回复不属于您，请勿非法操作")
		return
	}

	// TODO 删除讨论的回复
	r.DB.Delete(&reply)

	// TODO 移除损坏数据
	r.Redis.HDel(ctx, "Reply", id)

	// TODO 移除热度
	r.Redis.ZRem(ctx, "ReplyHot"+reply.CommentId.String(), reply.ID.String())

	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description   获取多篇讨论的回复
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r ReplyController) PageList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 分页
	var replys []model.Reply

	// TODO 查找所有分页中可见的条目
	r.DB.Where("comment_id = ?", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&replys)

	var total int64
	r.DB.Where("comment_id = ?", id).Model(model.Reply{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"replys": replys, "total": total}, "成功")
}

// @title    UserList
// @description   获取多篇指定用户讨论的回复
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r ReplyController) UserList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 分页
	var replys []model.Reply

	// TODO 查找所有分页中可见的条目
	r.DB.Where("user_id = ?", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&replys)

	var total int64
	r.DB.Where("user_id = ?", id).Model(model.Reply{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"replys": replys, "total": total}, "成功")
}

// @title    HotRanking
// @description   根据热度排行获取多篇讨论
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r ReplyController) HotRanking(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	// TODO 查找所有分页中可见的条目
	replys, err := r.Redis.ZRevRangeWithScores(ctx, "ReplyHot"+id, int64(pageNum-1)*int64(pageSize), int64(pageNum-1)*int64(pageSize)+int64(pageSize)-1).Result()

	if err != nil {
		response.Fail(ctx, nil, "获取失败")
	}

	// TODO 将redis中的数据取出
	total, _ := r.Redis.ZCard(ctx, "ReplyHot"+id).Result()

	// TODO 返回数据
	response.Success(ctx, gin.H{"replys": replys, "total": total}, "成功")
}

// @title    Like
// @description   点赞或点踩
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r ReplyController) Like(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	var reply model.Reply

	// TODO 查看讨论是否存在
	// TODO 先看redis中是否存在
	if ok, _ := r.Redis.HExists(ctx, "Reply", id).Result(); ok {
		cate, _ := r.Redis.HGet(ctx, "Reply", id).Result()
		if json.Unmarshal([]byte(cate), &reply) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			r.Redis.HDel(ctx, "Reply", id)
		}
	}

	// TODO 查看讨论的回复是否在数据库中存在
	if r.DB.Where("id = ?", id).First(&reply).Error != nil {
		response.Fail(ctx, nil, "讨论的回复不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(reply)
		r.Redis.HSet(ctx, "Reply", id, v)
	}
leep:
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var replyLike model.ReplyLike
	// TODO 如果没有点赞或者点踩
	if r.DB.Where("user_id = ? and reply_id = ?", user.ID, id).First(&replyLike).Error != nil {
		// TODO 插入数据
		replyLike = model.ReplyLike{
			ReplyId: reply.ID,
			UserId:  user.ID,
			Like:    like,
		}
		if err := r.DB.Create(&replyLike).Error; err != nil {
			response.Fail(ctx, nil, "点赞出错，数据库存储错误")
			return
		}
		// TODO 热度计算
		if like {
			r.Redis.ZIncrBy(ctx, "ReplyHot"+reply.CommentId.String(), 10.0, reply.ID.String())
		} else {
			r.Redis.ZIncrBy(ctx, "ReplyHot"+reply.CommentId.String(), -10.0, reply.ID.String())
		}
	} else {
		// TODO 热度计算
		if replyLike.Like {
			r.Redis.ZIncrBy(ctx, "ReplyHot"+reply.CommentId.String(), -10.0, reply.ID.String())
		} else {
			r.Redis.ZIncrBy(ctx, "ReplyHot"+reply.CommentId.String(), 10.0, reply.ID.String())
		}
		if like {
			r.Redis.ZIncrBy(ctx, "ReplyHot"+reply.CommentId.String(), 10.0, reply.ID.String())
		} else {
			r.Redis.ZIncrBy(ctx, "ReplyHot"+reply.CommentId.String(), -10.0, reply.ID.String())
		}
		r.DB.Where("user_id = ? and reply_id = ?", user.ID, id).Model(&model.ReplyLike{}).Update("like", like)
	}

	response.Success(ctx, nil, "点赞成功")
}

// @title    CancelLike
// @description   取消点赞或者点踩
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r ReplyController) CancelLike(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var reply model.Reply

	// TODO 查看讨论是否存在
	// TODO 先看redis中是否存在
	if ok, _ := r.Redis.HExists(ctx, "Reply", id).Result(); ok {
		cate, _ := r.Redis.HGet(ctx, "Reply", id).Result()
		if json.Unmarshal([]byte(cate), &reply) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			r.Redis.HDel(ctx, "Reply", id)
		}
	}

	// TODO 查看讨论的回复是否在数据库中存在
	if r.DB.Where("id = ?", id).First(&reply).Error != nil {
		response.Fail(ctx, nil, "讨论的回复不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(reply)
		r.Redis.HSet(ctx, "Reply", id, v)
	}
leep:

	// TODO 查看是否已经点赞或者点踩
	var replyLike model.ReplyLike
	if r.DB.Where("user_id = ? and reply_id = ?", user.ID, id).First(&replyLike).Error != nil {
		response.Fail(ctx, nil, "未点赞或点踩")
		return
	}

	// TODO 热度计算
	if replyLike.Like {
		r.Redis.ZIncrBy(ctx, "ReplyHot"+reply.CommentId.String(), -10.0, reply.ID.String())
	} else {
		r.Redis.ZIncrBy(ctx, "ReplyHot"+reply.CommentId.String(), 10.0, reply.ID.String())
	}

	// TODO 取消点赞或者点踩
	r.DB.Where("user_id = ? and reply_id = ?", user.ID, id).Delete(&model.ReplyLike{})

	response.Success(ctx, nil, "取消成功")
}

// @title    LikeNumber
// @description   点赞或点踩的数量
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r ReplyController) LikeNumber(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	var total int64

	// TODO 查看点赞或者点踩的数量
	r.DB.Where("reply_id = ? and like = ?", id, like).Count(&total)

	response.Success(ctx, gin.H{"total": total}, "查看成功")
}

// @title    LikeList
// @description   点赞或点踩的列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r ReplyController) LikeList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var replyLikes []model.ReplyLike

	var total int64

	// TODO 查看点赞或者点踩的数量
	r.DB.Where("reply_id = ? and like = ?", id, like).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&replyLikes).Count(&total)

	response.Success(ctx, gin.H{"replyLikes": replyLikes, "total": total}, "查看成功")
}

// @title    LikeShow
// @description   查看用户点赞状态
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r ReplyController) LikeShow(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var replyLike model.ReplyLike

	// TODO 查看点赞状态
	if r.DB.Where("user_id = ? and reply_id = ?", user.ID, id).First(&replyLike).Error != nil {
		response.Success(ctx, gin.H{"like": 0}, "暂无状态")
		return
	}

	if replyLike.Like {
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
func (r ReplyController) Likes(ctx *gin.Context) {

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	/// TODO 获取指定用户用户
	id := ctx.Params.ByName("id")

	// TODO 分页
	var replyLikes []model.ReplyLike

	var total int64

	// TODO 查看点赞或者点踩的数量
	r.DB.Where("user_id = ? and like = ?", id, like).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&replyLikes).Count(&total)

	response.Success(ctx, gin.H{"replyLikes": replyLikes, "total": total}, "查看成功")
}

// @title    NewReplyController
// @description   新建一个IReplyController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   IReplyController		返回一个IReplyController用于调用各种函数
func NewReplyController() IReplyController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.Reply{})
	db.AutoMigrate(model.ReplyLike{})
	return ReplyController{DB: db, Redis: redis}
}
