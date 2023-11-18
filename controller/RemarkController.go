// @Title  RemarkController
// @Description  该文件提供关于操作文章的回复的各种方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
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

// IRemarkController			定义了文章的回复类接口
type IRemarkController interface {
	Interface.RestInterface    // 包含增删查改功能
	Interface.LikeInterface    // 包含点赞功能
	UserList(ctx *gin.Context) // 指定用户的回复
	Interface.HotInterface     // 包含热度功能
}

// RemarkController			定义了文章的回复工具类
type RemarkController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Create
// @description   创建一篇文章的回复
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RemarkController) Create(ctx *gin.Context) {
	var requestRemark vo.RemarkRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestRemark); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 查找对应文章
	id := ctx.Params.ByName("id")

	var article model.Article

	// TODO 先尝试在redis中寻找
	if ok, _ := r.Redis.HExists(ctx, "Article", id).Result(); ok {
		art, _ := r.Redis.HGet(ctx, "Article", id).Result()
		if json.Unmarshal([]byte(art), &article) == nil {
			goto leep
		} else {
			// TODO 解码失败，删除字段
			r.Redis.HDel(ctx, "Article", id)
		}
	}

	// TODO 查看文章是否在数据库中存在
	if r.DB.Where("id = (?)", id).First(&article).Error != nil {
		response.Fail(ctx, nil, "文章不存在")
		return
	}
	{
		// TODO 将文章存入redis供下次使用
		v, _ := json.Marshal(article)
		r.Redis.HSet(ctx, "Article", id, v)
	}

leep:
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 创建文章的回复
	var remark = model.Remark{
		UserId:    user.ID,
		ArticleId: article.ID,
		Content:   requestRemark.Content,
		ResLong:   requestRemark.ResLong,
		ResShort:  requestRemark.ResShort,
	}

	// TODO 插入数据
	if err := r.DB.Create(&remark).Error; err != nil {
		response.Fail(ctx, nil, "文章的回复上传出错，数据库存储错误")
		return
	}

	// TODO 创建热度
	r.Redis.ZAdd(ctx, "RemarkHot"+article.ID.String(), redis.Z{Member: remark.ID.String(), Score: 100 + float64(time.Now().Unix()/86400)})

	// TODO 成功
	response.Success(ctx, gin.H{"remark": remark}, "创建成功")
}

// @title    Update
// @description   更新一篇文章的回复的内容
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RemarkController) Update(ctx *gin.Context) {
	var requestRemark vo.RemarkRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestRemark); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应文章的回复
	id := ctx.Params.ByName("id")

	var remark model.Remark

	if r.DB.Where("id = (?)", id).First(&remark).Error != nil {
		response.Fail(ctx, nil, "文章的回复不存在")
		return
	}

	// TODO 查看是否是用户作者
	if user.ID != remark.UserId {
		response.Fail(ctx, nil, "不是文章的回复作者，无法修改文章的回复")
		return
	}

	remarkUpdate := model.Remark{
		Content:  requestRemark.Content,
		ResLong:  requestRemark.ResLong,
		ResShort: requestRemark.ResShort,
	}

	// TODO 更新文章的回复内容
	r.DB.Where("id = (?)", id).Updates(remarkUpdate)

	// TODO 移除损坏数据
	r.Redis.HDel(ctx, "Remark", id)

	r.DB.Where("id = (?)", id).First(&remark)

	// TODO 成功
	response.Success(ctx, gin.H{"remark": remark}, "更新成功")
}

// @title    Show
// @description   查看一篇文章的回复的内容
// @auth      MGAronya       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RemarkController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var remark model.Remark

	// TODO 先看redis中是否存在
	if ok, _ := r.Redis.HExists(ctx, "Remark", id).Result(); ok {
		cate, _ := r.Redis.HGet(ctx, "Remark", id).Result()
		if json.Unmarshal([]byte(cate), &remark) == nil {
			response.Success(ctx, gin.H{"remark": remark}, "成功")
			return
		} else {
			// TODO 移除损坏数据
			r.Redis.HDel(ctx, "Remark", id)
		}
	}

	// TODO 查看文章的回复是否在数据库中存在
	if r.DB.Where("id = (?)", id).First(&remark).Error != nil {
		response.Fail(ctx, nil, "文章的回复不存在")
		return
	}

	response.Success(ctx, gin.H{"remark": remark}, "成功")

	// TODO 将提交存入redis供下次使用
	v, _ := json.Marshal(remark)
	r.Redis.HSet(ctx, "Remark", id, v)
}

// @title    Delete
// @description   删除一篇文章的回复
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RemarkController) Delete(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	var remark model.Remark

	// TODO 查看文章的回复是否存在
	if r.DB.Where("id = (?)", id).First(&remark).Error != nil {
		response.Fail(ctx, nil, "文章的回复不存在")
		return
	}

	// TODO 判断当前用户是否为文章的回复的作者
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有操作文章的回复的权力
	if user.ID != remark.UserId && user.Level < 4 {
		response.Fail(ctx, nil, "文章的回复不属于您，请勿非法操作")
		return
	}

	var total int64

	// TODO 查看点赞的数量
	r.DB.Where("remark_id = (?) and like = true", id).Model(model.RemarkLike{}).Count(&total)
	r.Redis.ZIncrBy(ctx, "UserLike", -float64(total), remark.UserId.String())

	// TODO 查看点踩的数量
	r.DB.Where("remark_id = (?) and like = false", id).Model(model.RemarkLike{}).Count(&total)
	r.Redis.ZIncrBy(ctx, "UserUnLike", -float64(total), remark.UserId.String())

	// TODO 删除文章的回复
	r.DB.Delete(&remark)

	// TODO 移除损坏数据
	r.Redis.HDel(ctx, "Remark", id)

	// TODO 移除热度
	r.Redis.ZRem(ctx, "RemarkHot"+remark.ArticleId.String(), remark.ID.String())

	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description   获取多篇文章的回复
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RemarkController) PageList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 分页
	var remarks []model.Remark

	// TODO 查找所有分页中可见的条目
	r.DB.Where("article_id = (?)", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&remarks)

	var total int64
	r.DB.Where("article_id = (?)", id).Model(model.Remark{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"remarks": remarks, "total": total}, "成功")
}

// @title    UserList
// @description   获取多篇指定用户文章的回复
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RemarkController) UserList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 分页
	var remarks []model.Remark

	// TODO 查找所有分页中可见的条目
	r.DB.Where("user_id = (?)", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&remarks)

	var total int64
	r.DB.Where("user_id = (?)", id).Model(model.Remark{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"remarks": remarks, "total": total}, "成功")
}

// @title    HotRanking
// @description   根据热度排行获取多篇讨论
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RemarkController) HotRanking(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	// TODO 查找所有分页中可见的条目
	remarks, err := r.Redis.ZRevRangeWithScores(ctx, "RemarkHot"+id, int64(pageNum-1)*int64(pageSize), int64(pageNum-1)*int64(pageSize)+int64(pageSize)-1).Result()

	if err != nil {
		response.Fail(ctx, nil, "获取失败")
	}

	for i := range remarks {
		remarks[i].Score -= float64(time.Now().Unix() / 86400)
	}

	// TODO 将redis中的数据取出
	total, _ := r.Redis.ZCard(ctx, "RemarkHot"+id).Result()

	// TODO 返回数据
	response.Success(ctx, gin.H{"remarks": remarks, "total": total}, "成功")
}

// @title    Like
// @description   点赞或点踩
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RemarkController) Like(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	var remark model.Remark

	// TODO 先看redis中是否存在
	if ok, _ := r.Redis.HExists(ctx, "Remark", id).Result(); ok {
		cate, _ := r.Redis.HGet(ctx, "Remark", id).Result()
		if json.Unmarshal([]byte(cate), &remark) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			r.Redis.HDel(ctx, "Remark", id)
		}
	}

	// TODO 查看文章的回复是否在数据库中存在
	if r.DB.Where("id = (?)", id).First(&remark).Error != nil {
		response.Fail(ctx, nil, "文章的回复不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(remark)
		r.Redis.HSet(ctx, "Remark", id, v)
	}
leep:
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var remarkLike model.RemarkLike
	// TODO 如果没有点赞或者点踩
	if r.DB.Where("user_id = (?) and remark_id = (?)", user.ID, id).First(&remarkLike).Error != nil {
		// TODO 插入数据
		remarkLike = model.RemarkLike{
			RemarkId: remark.ID,
			UserId:   user.ID,
			Like:     like,
		}
		if err := r.DB.Create(&remarkLike).Error; err != nil {
			response.Fail(ctx, nil, "点赞出错，数据库存储错误")
			return
		}
	} else {
		// TODO 热度计算
		if remarkLike.Like {
			r.Redis.ZIncrBy(ctx, "RemarkHot"+remark.ArticleId.String(), -10.0, remark.ID.String())
			r.Redis.ZIncrBy(ctx, "UserLike", -1, remark.UserId.String())
		} else {
			r.Redis.ZIncrBy(ctx, "RemarkHot"+remark.ArticleId.String(), 10.0, remark.ID.String())
			r.Redis.ZIncrBy(ctx, "UserUnLike", -1, remark.UserId.String())
		}
		r.DB.Where("user_id = (?) and remark_id = (?)", user.ID, id).Model(&model.RemarkLike{}).Update("like", like)
	}

	// TODO 热度计算
	if like {
		r.Redis.ZIncrBy(ctx, "RemarkHot"+remark.ArticleId.String(), 10.0, remark.ID.String())
		r.Redis.ZIncrBy(ctx, "UserLike", 1, remark.UserId.String())
	} else {
		r.Redis.ZIncrBy(ctx, "RemarkHot"+remark.ArticleId.String(), -10.0, remark.ID.String())
		r.Redis.ZIncrBy(ctx, "UserUnLike", 1, remark.UserId.String())
	}

	Handle.Behaviors["Likes"].PublishBehavior(1, remark.UserId)
	response.Success(ctx, nil, "点赞成功")
}

// @title    CancelLike
// @description   取消点赞或者点踩
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RemarkController) CancelLike(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var remark model.Remark

	// TODO 先看redis中是否存在
	if ok, _ := r.Redis.HExists(ctx, "Remark", id).Result(); ok {
		cate, _ := r.Redis.HGet(ctx, "Remark", id).Result()
		if json.Unmarshal([]byte(cate), &remark) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			r.Redis.HDel(ctx, "Remark", id)
		}
	}

	// TODO 查看文章的回复是否在数据库中存在
	if r.DB.Where("id = (?)", id).First(&remark).Error != nil {
		response.Fail(ctx, nil, "文章的回复不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(remark)
		r.Redis.HSet(ctx, "Remark", id, v)
	}
leep:

	// TODO 查看是否已经点赞或者点踩
	var remarkLike model.RemarkLike
	if r.DB.Where("user_id = (?) and remark_id = (?)", user.ID, id).First(&remarkLike).Error != nil {
		response.Fail(ctx, nil, "未点赞或点踩")
		return
	}

	// TODO 热度计算
	if remarkLike.Like {
		r.Redis.ZIncrBy(ctx, "RemarkHot"+remark.ArticleId.String(), -10.0, remark.ID.String())
		r.Redis.ZIncrBy(ctx, "UserLike", -1, remark.UserId.String())
	} else {
		r.Redis.ZIncrBy(ctx, "RemarkHot"+remark.ArticleId.String(), 10.0, remark.ID.String())
		r.Redis.ZIncrBy(ctx, "UserUnLike", -1, remark.UserId.String())
	}

	// TODO 取消点赞或者点踩
	r.DB.Where("user_id = (?) and remark_id = (?)", user.ID, id).Delete(&model.RemarkLike{})

	Handle.Behaviors["Likes"].PublishBehavior(-1, remark.UserId)
	response.Success(ctx, nil, "取消成功")
}

// @title    LikeNumber
// @description   点赞或点踩的数量
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RemarkController) LikeNumber(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	var total int64

	// TODO 查看点赞或者点踩的数量
	r.DB.Where("remark_id = (?) and `like` = ?", id, like).Model(model.RemarkLike{}).Count(&total)

	response.Success(ctx, gin.H{"total": total}, "查看成功")
}

// @title    LikeList
// @description   点赞或点踩的列表
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RemarkController) LikeList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var remarkLikes []model.RemarkLike

	var total int64

	// TODO 查看点赞或者点踩的数量
	r.DB.Where("remark_id = (?) and `like` = ?", id, like).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&remarkLikes)

	r.DB.Where("remark_id = (?) and `like` = ?", id, like).Model(model.RemarkLike{}).Count(&total)

	response.Success(ctx, gin.H{"remarkLikes": remarkLikes, "total": total}, "查看成功")
}

// @title    LikeShow
// @description   查看用户点赞状态
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RemarkController) LikeShow(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var remarkLike model.RemarkLike

	// TODO 查看点赞状态
	if r.DB.Where("user_id = (?) and remark_id = (?)", user.ID, id).First(&remarkLike).Error != nil {
		response.Success(ctx, gin.H{"like": 0}, "暂无状态")
		return
	}

	if remarkLike.Like {
		response.Success(ctx, gin.H{"like": 1}, "已点赞")
	} else {
		response.Success(ctx, gin.H{"like": -1}, "已点踩")
	}

}

// @title    Likes
// @description   查看用户点赞状态
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (r RemarkController) Likes(ctx *gin.Context) {

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	/// TODO 获取指定用户用户
	id := ctx.Params.ByName("id")

	// TODO 分页
	var remarkLikes []model.RemarkLike

	var total int64

	// TODO 查看点赞或者点踩的数量
	r.DB.Where("user_id = (?) and `like` = ?", id, like).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&remarkLikes)

	r.DB.Where("user_id = (?) and `like` = ?", id, like).Model(model.RemarkLike{}).Count(&total)

	response.Success(ctx, gin.H{"remarkLikes": remarkLikes, "total": total}, "查看成功")
}

// @title    NewRemarkController
// @description   新建一个IRemarkController
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   IRemarkController		返回一个IRemarkController用于调用各种函数
func NewRemarkController() IRemarkController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.Remark{})
	db.AutoMigrate(model.RemarkLike{})
	return RemarkController{DB: db, Redis: redis}
}
