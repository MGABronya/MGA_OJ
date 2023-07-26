// @Title  CommentController
// @Description  该文件提供关于操作讨论的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	Handle "MGA_OJ/Behavior"
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
	"gorm.io/gorm"
)

// ICommentController			定义了讨论类接口
type ICommentController interface {
	Interface.RestInterface    // 包含增删查改功能
	Interface.LikeInterface    // 包含点赞功能
	Interface.HotInterface     // 包含热度功能
	UserList(ctx *gin.Context) // 查看指定用户的讨论
}

// CommentController			定义了讨论工具类
type CommentController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Create
// @description   创建一篇讨论
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CommentController) Create(ctx *gin.Context) {
	var requestComment vo.CommentRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestComment); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 查找对应题目
	id := ctx.Params.ByName("id")

	var problem model.Problem

	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Problem", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Problem", id).Result()
		if json.Unmarshal([]byte(cate), &problem) == nil {
			// TODO 跳过数据库搜寻problem过程
			goto leep
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Problem", id)
		}
	}

	// TODO 查看题目是否在数据库中存在
	if c.DB.Where("id = (?)", id).First(&problem).Error != nil {
		response.Fail(ctx, nil, "题目不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(problem)
		c.Redis.HSet(ctx, "Problem", id, v)
	}

leep:
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 创建评论
	var comment = model.Comment{
		UserId:    user.ID,
		ProblemId: problem.ID,
		Content:   requestComment.Content,
		ResLong:   requestComment.ResLong,
		ResShort:  requestComment.ResShort,
	}

	// TODO 插入数据
	if err := c.DB.Create(&comment).Error; err != nil {
		response.Fail(ctx, nil, "讨论上传出错，数据库存储错误")
		return
	}
	// TODO 创建热度
	c.Redis.ZAdd(ctx, "CommentHot"+problem.ID.String(), redis.Z{Member: comment.ID.String(), Score: 100 + float64(time.Now().Unix()/86400)})

	// TODO 成功
	response.Success(ctx, gin.H{"comment": comment}, "创建成功")
}

// @title    Update
// @description   更新一篇讨论的内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CommentController) Update(ctx *gin.Context) {
	var requestComment model.Comment
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestComment); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应讨论
	id := ctx.Params.ByName("id")

	var comment model.Comment

	if c.DB.Where("id = (?)", id).First(&comment).Error != nil {
		response.Fail(ctx, nil, "讨论不存在")
		return
	}

	// TODO 查看是否是用户作者
	if user.ID != comment.UserId {
		response.Fail(ctx, nil, "不是讨论作者，无法修改讨论")
		return
	}

	// TODO 更新讨论内容
	c.DB.Where("id = (?)", id).Updates(requestComment)

	// TODO 移除损坏数据
	c.Redis.HDel(ctx, "Comment", id)

	// TODO 成功
	response.Success(ctx, nil, "更新成功")
}

// @title    Show
// @description   查看一篇讨论的内容
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CommentController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var comment model.Comment

	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Comment", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Comment", id).Result()
		if json.Unmarshal([]byte(cate), &comment) == nil {
			response.Success(ctx, gin.H{"comment": comment}, "成功")
			return
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Comment", id)
		}
	}

	// TODO 查看讨论是否在数据库中存在
	if c.DB.Where("id = (?)", id).First(&comment).Error != nil {
		response.Fail(ctx, nil, "讨论不存在")
		return
	}

	response.Success(ctx, gin.H{"comment": comment}, "成功")

	// TODO 将讨论存入redis供下次使用
	v, _ := json.Marshal(comment)
	c.Redis.HSet(ctx, "Comment", id, v)
}

// @title    Delete
// @description   删除一篇讨论
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CommentController) Delete(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	var comment model.Comment

	// TODO 查看讨论是否存在
	if c.DB.Where("id = (?)", id).First(&comment).Error != nil {
		response.Fail(ctx, nil, "讨论不存在")
		return
	}

	// TODO 判断当前用户是否为讨论的作者
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有操作讨论的权力
	if user.ID != comment.UserId && user.Level < 4 {
		response.Fail(ctx, nil, "讨论不属于您，请勿非法操作")
		return
	}

	var total int64

	// TODO 查看点赞的数量
	c.DB.Where("comment_id = (?) and like = true", id).Model(model.CommentLike{}).Count(&total)
	c.Redis.ZIncrBy(ctx, "UserLike", -float64(total), comment.UserId.String())

	// TODO 查看点踩的数量
	c.DB.Where("comment_id = (?) and like = false", id).Model(model.CommentLike{}).Count(&total)
	c.Redis.ZIncrBy(ctx, "UserUnLike", -float64(total), comment.UserId.String())

	// TODO 删除讨论
	c.DB.Delete(&comment)

	// TODO 移除损坏数据
	c.Redis.HDel(ctx, "Comment", id)

	// TODO 移除热度
	c.Redis.ZRem(ctx, "CommentHot"+comment.ProblemId.String(), comment.ID.String())

	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description   获取多篇讨论
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CommentController) PageList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 分页
	var comments []model.Comment

	// TODO 查找所有分页中可见的条目
	c.DB.Where("problem_id = (?)", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&comments)

	var total int64
	c.DB.Where("problem_id = (?)", id).Model(model.Comment{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"comments": comments, "total": total}, "成功")
}

// @title    HotRanking
// @description   根据热度排行获取多篇讨论
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CommentController) HotRanking(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	// TODO 查找所有分页中可见的条目
	comments, err := c.Redis.ZRevRangeWithScores(ctx, "CommentHot"+id, int64(pageNum-1)*int64(pageSize), int64(pageNum-1)*int64(pageSize)+int64(pageSize)-1).Result()

	if err != nil {
		response.Fail(ctx, nil, "获取失败")
	}

	for i := range comments {
		comments[i].Score -= float64(time.Now().Unix() / 86400)
	}

	// TODO 将redis中的数据取出
	total, _ := c.Redis.ZCard(ctx, "Comment"+id).Result()

	// TODO 返回数据
	response.Success(ctx, gin.H{"comments": comments, "total": total}, "成功")
}

// @title    UserList
// @description   获取指定用户的多篇讨论
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CommentController) UserList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 分页
	var comments []model.Comment

	// TODO 查找所有分页中可见的条目
	c.DB.Where("user_id = (?)", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&comments)

	var total int64
	c.DB.Where("user_id = (?)", id).Model(model.Comment{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"comments": comments, "total": total}, "成功")
}

// @title    Like
// @description   点赞或点踩
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CommentController) Like(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	var comment model.Comment

	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Comment", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Comment", id).Result()
		if json.Unmarshal([]byte(cate), &comment) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Comment", id)
		}
	}

	// TODO 查看讨论是否在数据库中存在
	if c.DB.Where("id = (?)", id).First(&comment).Error != nil {
		response.Fail(ctx, nil, "讨论不存在")
		return
	}
	{
		// TODO 将讨论存入redis供下次使用
		v, _ := json.Marshal(comment)
		c.Redis.HSet(ctx, "Comment", id, v)
	}
leep:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var commentLike model.CommentLike
	// TODO 如果没有点赞或者点踩
	if c.DB.Where("user_id = (?) and comment_id = (?)", user.ID, id).First(&commentLike).Error != nil {
		// TODO 插入数据
		commentLike := model.CommentLike{
			CommentId: comment.ID,
			UserId:    user.ID,
			Like:      like,
		}
		if err := c.DB.Create(&commentLike).Error; err != nil {
			response.Fail(ctx, nil, "点赞出错，数据库存储错误")
			return
		}
	} else {
		// TODO 热度计算
		if commentLike.Like {
			c.Redis.ZIncrBy(ctx, "CommentHot"+comment.ProblemId.String(), -10.0, comment.ID.String())
			c.Redis.ZIncrBy(ctx, "UserLike", -1, comment.UserId.String())
		} else {
			c.Redis.ZIncrBy(ctx, "CommentHot"+comment.ProblemId.String(), 10.0, comment.ID.String())
			c.Redis.ZIncrBy(ctx, "UserUnLike", -1, comment.UserId.String())
		}
		c.DB.Where("user_id = (?) and comment_id = (?)", user.ID, id).Model(&model.CommentLike{}).Update("like", like)
	}

	// TODO 热度计算
	if like {
		c.Redis.ZIncrBy(ctx, "CommentHot"+comment.ProblemId.String(), 10.0, comment.ID.String())
		c.Redis.ZIncrBy(ctx, "UserLike", 1, comment.UserId.String())
	} else {
		c.Redis.ZIncrBy(ctx, "CommentHot"+comment.ProblemId.String(), -10.0, comment.ID.String())
		c.Redis.ZIncrBy(ctx, "UserUnLike", 1, comment.UserId.String())
	}

	Handle.Behaviors["Likes"].PublishBehavior(1, comment.UserId)
	response.Success(ctx, nil, "点赞成功")
}

// @title    CancelLike
// @description   取消点赞或者点踩
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CommentController) CancelLike(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var comment model.Comment

	// TODO 先看redis中是否存在
	if ok, _ := c.Redis.HExists(ctx, "Comment", id).Result(); ok {
		cate, _ := c.Redis.HGet(ctx, "Comment", id).Result()
		if json.Unmarshal([]byte(cate), &comment) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			c.Redis.HDel(ctx, "Comment", id)
		}
	}

	// TODO 查看讨论是否在数据库中存在
	if c.DB.Where("id = (?)", id).First(&comment).Error != nil {
		response.Fail(ctx, nil, "讨论不存在")
		return
	}
	{
		// TODO 将讨论存入redis供下次使用
		v, _ := json.Marshal(comment)
		c.Redis.HSet(ctx, "Comment", id, v)
	}
leep:

	// TODO 查看是否已经点赞或者点踩
	var commentLike model.CommentLike
	if c.DB.Where("user_id = (?) and comment_id = (?)", user.ID, id).First(&commentLike).Error != nil {
		response.Fail(ctx, nil, "未点赞或点踩")
		return
	}

	// TODO 热度计算
	if commentLike.Like {
		c.Redis.ZIncrBy(ctx, "CommentHot"+comment.ProblemId.String(), -10.0, commentLike.CommentId.String())
		c.Redis.ZIncrBy(ctx, "UserLike", -1, comment.UserId.String())
	} else {
		c.Redis.ZIncrBy(ctx, "CommentHot"+comment.ProblemId.String(), 10.0, commentLike.CommentId.String())
		c.Redis.ZIncrBy(ctx, "UserUnLike", -1, comment.UserId.String())
	}

	// TODO 取消点赞或者点踩
	c.DB.Where("user_id = (?) and comment_id = (?)", user.ID, id).Delete(&model.CommentLike{})
	Handle.Behaviors["Likes"].PublishBehavior(-1, comment.UserId)
	response.Success(ctx, nil, "取消成功")
}

// @title    LikeNumber
// @description   点赞或点踩的数量
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CommentController) LikeNumber(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	var total int64

	// TODO 查看点赞或者点踩的数量
	c.DB.Where("comment_id = (?) and `like` = ?", id, like).Model(model.CommentLike{}).Count(&total)

	response.Success(ctx, gin.H{"total": total}, "查看成功")
}

// @title    LikeList
// @description   点赞或点踩的列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CommentController) LikeList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var commentLikes []model.CommentLike

	var total int64

	// TODO 查看点赞或者点踩的数量
	c.DB.Where("comment_id = (?) and `like` = ?", id, like).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&commentLikes)

	c.DB.Where("comment_id = (?) and `like` = ?", id, like).Model(model.CommentLike{}).Count(&total)

	response.Success(ctx, gin.H{"commentLikes": commentLikes, "total": total}, "查看成功")
}

// @title    LikeShow
// @description   查看用户点赞状态
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CommentController) LikeShow(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var commentLike model.CommentLike

	// TODO 查看点赞状态
	if c.DB.Where("user_id = (?) and comment_id = (?)", user.ID, id).First(&commentLike).Error != nil {
		response.Success(ctx, gin.H{"like": 0}, "暂无状态")
		return
	}

	if commentLike.Like {
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
func (c CommentController) Likes(ctx *gin.Context) {

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取指定用户用户
	id := ctx.Params.ByName("id")

	// TODO 分页
	var commentLikes []model.CommentLike

	var total int64

	// TODO 查看点赞或者点踩的数量
	c.DB.Where("user_id = (?) and `like` = ?", id, like).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&commentLikes)

	c.DB.Where("user_id = (?) and `like` = ?", id, like).Model(model.CommentLike{}).Count(&total)

	response.Success(ctx, gin.H{"commentLikes": commentLikes, "total": total}, "查看成功")
}

// @title    NewCommentController
// @description   新建一个ICommentController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   ICommentController		返回一个ICommentController用于调用各种函数
func NewCommentController() ICommentController {
	db := common.GetDB()
	db.AutoMigrate(model.Comment{})
	db.AutoMigrate(model.CommentLike{})
	return CommentController{DB: db}
}
