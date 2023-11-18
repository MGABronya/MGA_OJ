// @Title  ArticleController
// @Description  该文件提供关于操作文章的各种方法
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

	"github.com/go-redis/redis/v9"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IArticleController			定义了文章类接口
type IArticleController interface {
	Interface.RestInterface        // 包含增删查改功能
	Interface.LikeInterface        // 包含点赞功能
	Interface.CollectInterface     // 包含收藏功能
	Interface.VisitInterface       // 包含游览功能
	Interface.SearchInterface      // 包含搜索功能
	Interface.HotInterface         // 包含热度功能
	UserList(ctx *gin.Context)     // 查询指定用户的文章
	CategoryList(ctx *gin.Context) // 查询某分类的文章
	Interface.LabelInterface       // 包含标签功能
}

// ArticleController			定义了文章工具类
type ArticleController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Create
// @description   创建一篇文章
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) Create(ctx *gin.Context) {
	var requestArticle vo.ArticleRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestArticle); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 创建文章
	article := model.Article{
		Title:      requestArticle.Title,
		Content:    requestArticle.Content,
		ResLong:    requestArticle.ResLong,
		ResShort:   requestArticle.ResShort,
		UserId:     user.ID,
		CategoryId: requestArticle.CategoryId,
	}

	// TODO 插入数据
	if err := a.DB.Create(&article).Error; err != nil {
		response.Fail(ctx, nil, "文章上传出错，数据验证有误")
		return
	}

	a.Redis.ZAdd(ctx, "ArticleHot", redis.Z{Member: article.ID.String(), Score: 100 + float64(time.Now().Unix()/86400)})

	// TODO 成功
	response.Success(ctx, gin.H{"article": article}, "创建成功")
}

// @title    Update
// @description   更新一篇文章的内容
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) Update(ctx *gin.Context) {
	var requestArticle vo.ArticleRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestArticle); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应文章
	id := ctx.Params.ByName("id")

	var article model.Article

	if a.DB.Where("id = (?)", id).First(&article).Error != nil {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	// TODO 查看是否是用户作者
	if user.ID != article.UserId {
		response.Fail(ctx, nil, "不是文章作者，无法修改文章")
		return
	}

	// TODO 新建文章
	articleUpdate := model.Article{
		Title:      requestArticle.Title,
		Content:    requestArticle.Content,
		ResLong:    requestArticle.ResLong,
		ResShort:   requestArticle.ResShort,
		CategoryId: requestArticle.CategoryId,
	}

	// TODO 更新文章内容
	a.DB.Model(model.Article{}).Where("id = (?)", id).Updates(articleUpdate)

	// TODO 解码失败，删除字段
	a.Redis.HDel(ctx, "Article", id)

	a.DB.Where("id = (?)", id).First(&article)

	// TODO 成功
	response.Success(ctx, gin.H{"article": article}, "更新成功")
}

// @title    Show
// @description   查看一篇文章的内容
// @auth      MGAronya       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	var article model.Article

	// TODO 先尝试在redis中寻找
	if ok, _ := a.Redis.HExists(ctx, "Article", id).Result(); ok {
		art, _ := a.Redis.HGet(ctx, "Article", id).Result()
		if json.Unmarshal([]byte(art), &article) == nil {
			response.Success(ctx, gin.H{"article": article}, "成功")
			return
		} else {
			// TODO 解码失败，删除字段
			a.Redis.HDel(ctx, "Article", id)
		}
	}

	// TODO 查看文章是否在数据库中存在
	if a.DB.Where("id = (?)", id).First(&article).Error != nil {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	response.Success(ctx, gin.H{"article": article}, "成功")

	// TODO 将文章存入redis供下次使用
	v, _ := json.Marshal(article)
	a.Redis.HSet(ctx, "Article", id, v)
}

// @title    Delete
// @description   删除一篇文章
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) Delete(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	var article model.Article

	// TODO 查看文章是否存在
	if a.DB.Where("id = (?)", id).First(&article).Error != nil {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	// TODO 判断当前用户是否为文章的作者
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有操作文章的权力
	if user.ID != article.UserId && user.Level < 4 {
		response.Fail(ctx, nil, "文章不属于您，请勿非法操作")
		return
	}

	var total int64

	// TODO 查看点赞的数量
	a.DB.Where("article_id = (?) and like = true", id).Model(model.ArticleLike{}).Count(&total)
	a.Redis.ZIncrBy(ctx, "UserLike", -float64(total), article.UserId.String())

	// TODO 查看点踩的数量
	a.DB.Where("article_id = (?) and like = false", id).Model(model.ArticleLike{}).Count(&total)
	a.Redis.ZIncrBy(ctx, "UserUnLike", -float64(total), article.UserId.String())

	// TODO 查看收藏的数量
	a.DB.Where("article_id = (?)", id).Model(model.ArticleCollect{}).Count(&total)
	a.Redis.ZIncrBy(ctx, "UserCollect", -float64(total), article.UserId.String())

	// TODO 获取阅读人数
	total, _ = a.Redis.PFCount(ctx, "ArticleVisit"+id).Result()
	a.Redis.ZIncrBy(ctx, "UserVisit", -float64(total), article.UserId.String())
	a.Redis.Del(ctx, "ArticleVisit"+id)

	// TODO 删除文章
	a.DB.Delete(&article)

	// TODO 解码失败，删除字段
	a.Redis.HDel(ctx, "Article", id)
	// TODO 删除热度
	a.Redis.ZRem(ctx, "ArticleHot", article.ID.String())

	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description   获取多篇文章
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) PageList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var articles []model.Article

	// TODO 查找所有分页中可见的条目
	a.DB.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&articles)

	var total int64
	a.DB.Model(model.Article{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"articles": articles, "total": total}, "成功")
}

// @title    UserList
// @description   获取指定用户的文章
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) UserList(ctx *gin.Context) {

	// TODO 获取指定用户用户
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var articles []model.Article

	// TODO 查找所有分页中可见的条目
	a.DB.Where("user_id = (?)", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&articles)

	var total int64
	a.DB.Where("user_id = (?)", id).Model(model.Article{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"articles": articles, "total": total}, "成功")
}

// @title    CategoryList
// @description   获取指定分类的文章
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) CategoryList(ctx *gin.Context) {

	// TODO 获取指定分类
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var articles []model.Article

	// TODO 查找所有分页中可见的条目
	a.DB.Where("category_id = (?)", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&articles)

	var total int64
	a.DB.Where("category_id = (?)", id).Model(model.Article{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"articles": articles, "total": total}, "成功")
}

// @title    HotRanking
// @description   按热度排行获取多篇文章
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) HotRanking(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 查找所有分页中可见的条目
	articles, err := a.Redis.ZRevRangeWithScores(ctx, "ArticleHot", int64(pageNum-1)*int64(pageSize), int64(pageNum-1)*int64(pageSize)+int64(pageSize)-1).Result()

	if err != nil {
		response.Fail(ctx, nil, "获取失败")
	}

	for i := range articles {
		articles[i].Score -= float64(time.Now().Unix() / 86400)
	}

	// TODO 将redis中的数据取出
	total, _ := a.Redis.ZCard(ctx, "ArticleHot").Result()

	// TODO 返回数据
	response.Success(ctx, gin.H{"articles": articles, "total": total}, "成功")
}

// @title    Like
// @description   点赞或点踩
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) Like(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	var article model.Article

	// TODO 先尝试在redis中寻找
	if ok, _ := a.Redis.HExists(ctx, "Article", id).Result(); ok {
		art, _ := a.Redis.HGet(ctx, "Article", id).Result()
		if json.Unmarshal([]byte(art), &article) == nil {
			goto leep
		} else {
			// TODO 解码失败，删除字段
			a.Redis.HDel(ctx, "Article", id)
		}
	}

	// TODO 查看文章是否在数据库中存在
	if a.DB.Where("id = (?)", id).First(&article).Error != nil {
		response.Fail(ctx, nil, "文章不存在")
		return
	}
	{
		// TODO 将文章存入redis供下次使用
		v, _ := json.Marshal(article)
		a.Redis.HSet(ctx, "Article", id, v)
	}
leep:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var articleLike model.ArticleLike
	// TODO 如果没有点赞或者点踩
	if a.DB.Where("user_id = (?) and article_id = (?)", user.ID, id).First(&articleLike).Error != nil {
		// TODO 插入数据
		articleLike = model.ArticleLike{
			ArticleId: article.ID,
			UserId:    user.ID,
			Like:      like,
		}
		if err := a.DB.Create(&articleLike).Error; err != nil {
			response.Fail(ctx, nil, "点赞出错，数据库存储错误")
			return
		}
	} else {
		// TODO 热度计算
		if articleLike.Like {
			a.Redis.ZIncrBy(ctx, "ArticleHot", -10.0, article.ID.String())
			a.Redis.ZIncrBy(ctx, "UserLike", -1, article.UserId.String())
		} else {
			a.Redis.ZIncrBy(ctx, "ArticleHot", 10.0, article.ID.String())
			a.Redis.ZIncrBy(ctx, "UserUnLike", -1, article.UserId.String())
		}
		a.DB.Where("user_id = (?) and article_id = (?)", user.ID, id).Model(&model.ArticleLike{}).Update("like", like)
	}

	// TODO 热度计算
	if like {
		a.Redis.ZIncrBy(ctx, "ArticleHot", 10.0, article.ID.String())
		a.Redis.ZIncrBy(ctx, "UserLike", 1, article.UserId.String())
	} else {
		a.Redis.ZIncrBy(ctx, "ArticleHot", -10.0, article.ID.String())
		a.Redis.ZIncrBy(ctx, "UserUnLike", 1, article.UserId.String())
	}

	Handle.Behaviors["Likes"].PublishBehavior(1, article.UserId)
	response.Success(ctx, nil, "点赞成功")
}

// @title    CancelLike
// @description   取消点赞或者点踩
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) CancelLike(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var article model.Article

	// TODO 先尝试在redis中寻找
	if ok, _ := a.Redis.HExists(ctx, "Article", id).Result(); ok {
		art, _ := a.Redis.HGet(ctx, "Article", id).Result()
		if json.Unmarshal([]byte(art), &article) == nil {
			goto leep
		} else {
			// TODO 解码失败，删除字段
			a.Redis.HDel(ctx, "Article", id)
		}
	}

	// TODO 查看文章是否在数据库中存在
	if a.DB.Where("id = (?)", id).First(&article).Error != nil {
		response.Fail(ctx, nil, "文章不存在")
		return
	}
	{
		// TODO 将文章存入redis供下次使用
		v, _ := json.Marshal(article)
		a.Redis.HSet(ctx, "Article", id, v)
	}
leep:

	// TODO 查看是否已经点赞或者点踩
	var articleLike model.ArticleLike
	if a.DB.Where("user_id = (?) and article_id = (?)", user.ID, id).First(&articleLike).Error != nil {
		response.Fail(ctx, nil, "未点赞或点踩")
		return
	}

	// TODO 热度计算
	if articleLike.Like {
		a.Redis.ZIncrBy(ctx, "ArticleHot", -10.0, article.ID.String())
		a.Redis.ZIncrBy(ctx, "UserLike", -1, article.UserId.String())
	} else {
		a.Redis.ZIncrBy(ctx, "ArticleHot", 10.0, article.ID.String())
		a.Redis.ZIncrBy(ctx, "UserUnLike", -1, article.UserId.String())
	}

	// TODO 取消点赞或者点踩
	a.DB.Where("user_id = (?) and article_id = (?)", user.ID, id).Delete(&model.ArticleLike{})

	Handle.Behaviors["Likes"].PublishBehavior(-1, article.UserId)
	response.Success(ctx, nil, "取消成功")
}

// @title    LikeNumber
// @description   点赞或点踩的数量
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) LikeNumber(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	var total int64

	// TODO 查看点赞或者点踩的数量
	a.DB.Where("article_id = (?) and `like` = ?", id, like).Model(model.ArticleLike{}).Count(&total)

	response.Success(ctx, gin.H{"total": total}, "查看成功")
}

// @title    LikeList
// @description   点赞或点踩的列表
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) LikeList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var articleLikes []model.ArticleLike

	var total int64

	// TODO 查看点赞或者点踩的数量
	a.DB.Where("article_id = (?) and `like` = ?", id, like).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&articleLikes)

	a.DB.Where("article_id = (?) and `like` = ?", id, like).Model(model.ArticleLike{}).Count(&total)
	response.Success(ctx, gin.H{"articleLikes": articleLikes, "total": total}, "查看成功")
}

// @title    LikeShow
// @description   查看用户点赞状态
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) LikeShow(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var articleLike model.ArticleLike

	// TODO 查看点赞状态
	if a.DB.Where("user_id = (?) and article_id = (?)", user.ID, id).First(&articleLike).Error != nil {
		response.Success(ctx, gin.H{"like": 0}, "暂无状态")
		return
	}

	if articleLike.Like {
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
func (a ArticleController) Likes(ctx *gin.Context) {

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取指定用户用户
	id := ctx.Params.ByName("id")

	// TODO 分页
	var articleLikes []model.ArticleLike

	var total int64

	// TODO 查看点赞或者点踩的数量
	a.DB.Where("user_id = (?) and `like` = ?", id, like).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&articleLikes)

	response.Success(ctx, gin.H{"articleLikes": articleLikes, "total": total}, "查看成功")
}

// @title    Collect
// @description   收藏
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) Collect(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	var article model.Article

	// TODO 先尝试在redis中寻找
	if ok, _ := a.Redis.HExists(ctx, "Article", id).Result(); ok {
		art, _ := a.Redis.HGet(ctx, "Article", id).Result()
		if json.Unmarshal([]byte(art), &article) == nil {
			goto leep
		} else {
			// TODO 解码失败，删除字段
			a.Redis.HDel(ctx, "Article", id)
		}
	}

	// TODO 查看文章是否在数据库中存在
	if a.DB.Where("id = (?)", id).First(&article).Error != nil {
		response.Fail(ctx, nil, "文章不存在")
		return
	}
	{
		// TODO 将文章存入redis供下次使用
		v, _ := json.Marshal(article)
		a.Redis.HSet(ctx, "Article", id, v)
	}
leep:
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有收藏
	if a.DB.Where("user_id = (?) and article_id = (?)", user.ID, article.ID).First(&model.ArticleCollect{}).Error != nil {
		articleCollect := model.ArticleCollect{
			ArticleId: article.ID,
			UserId:    user.ID,
		}
		// TODO 插入数据
		if err := a.DB.Create(&articleCollect).Error; err != nil {
			response.Fail(ctx, nil, "收藏出错，数据库存储错误")
			return
		}
	} else {
		response.Fail(ctx, nil, "已收藏")
		return
	}

	// TODO 热度计算
	a.Redis.ZIncrBy(ctx, "ArticleHot", 50.0, article.ID.String())
	a.Redis.ZIncrBy(ctx, "UserCollect", 1, article.UserId.String())

	Handle.Behaviors["Collects"].PublishBehavior(1, article.UserId)

	response.Success(ctx, nil, "收藏成功")
}

// @title    CancelCollect
// @description   取消收藏
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) CancelCollect(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	var article model.Article

	// TODO 先尝试在redis中寻找
	if ok, _ := a.Redis.HExists(ctx, "Article", id).Result(); ok {
		art, _ := a.Redis.HGet(ctx, "Article", id).Result()
		if json.Unmarshal([]byte(art), &article) == nil {
			goto leep
		} else {
			// TODO 解码失败，删除字段
			a.Redis.HDel(ctx, "Article", id)
		}
	}

	// TODO 查看文章是否在数据库中存在
	if a.DB.Where("id = (?)", id).First(&article).Error != nil {
		response.Fail(ctx, nil, "文章不存在")
		return
	}
	{
		// TODO 将文章存入redis供下次使用
		v, _ := json.Marshal(article)
		a.Redis.HSet(ctx, "Article", id, v)
	}
leep:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有收藏
	if a.DB.Where("user_id = (?) and article_id = (?)", user.ID, id).First(&model.ArticleCollect{}).Error != nil {
		response.Fail(ctx, nil, "未收藏")
	} else {
		a.DB.Where("user_id = (?) and article_id = (?)", user.ID, id).Delete(&model.ArticleCollect{})
		Handle.Behaviors["Collects"].PublishBehavior(-1, article.UserId)
		response.Success(ctx, nil, "取消收藏成功")
		// TODO 热度计算
		a.Redis.ZIncrBy(ctx, "ArticleHot", -50.0, id)
		// TODO 删除存储入库
		a.Redis.ZIncrBy(ctx, "UserCollect", -1, article.UserId.String())
	}
}

// @title    CollectShow
// @description   查看收藏状态
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) CollectShow(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有收藏
	if a.DB.Where("user_id = (?) and article_id = (?)", user.ID, id).First(&model.ArticleCollect{}).Error != nil {
		response.Success(ctx, gin.H{"collect": false}, "未收藏")
		return
	} else {
		response.Success(ctx, gin.H{"collect": true}, "已收藏")
		return
	}
}

// @title    CollectList
// @description   查看收藏用户列表
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) CollectList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var articleCollects []model.ArticleCollect

	var total int64

	// TODO 查看收藏的数量
	a.DB.Where("article_id = (?)", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&articleCollects)

	a.DB.Where("article_id = (?)", id).Model(model.ArticleCollect{}).Count(&total)

	response.Success(ctx, gin.H{"articleCollects": articleCollects, "total": total}, "查看成功")
}

// @title    CollectNumber
// @description   查看收藏用户数量
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) CollectNumber(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	var total int64

	// TODO 查看收藏的数量
	a.DB.Where("article_id = (?)", id).Model(model.ArticleCollect{}).Count(&total)

	response.Success(ctx, gin.H{"total": total}, "查看成功")
}

// @title    Collects
// @description   查看用户收藏夹
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) Collects(ctx *gin.Context) {
	// TODO 获取指定用户用户
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var articleCollects []model.ArticleCollect

	var total int64

	// TODO 查看收藏的数量
	a.DB.Where("user_id = (?)", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&articleCollects)

	a.DB.Where("user_id = (?)", id).Model(model.ArticleCollect{}).Count(&total)

	response.Success(ctx, gin.H{"articleCollects": articleCollects, "total": total}, "查看成功")
}

// @title    Visit
// @description   游览文章
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) Visit(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	var article model.Article

	// TODO 先尝试在redis中寻找
	if ok, _ := a.Redis.HExists(ctx, "Article", id).Result(); ok {
		art, _ := a.Redis.HGet(ctx, "Article", id).Result()
		if json.Unmarshal([]byte(art), &article) == nil {
			goto leep
		} else {
			// TODO 解码失败，删除字段
			a.Redis.HDel(ctx, "Article", id)
		}
	}

	// TODO 查看文章是否在数据库中存在
	if a.DB.Where("id = (?)", id).First(&article).Error != nil {
		response.Fail(ctx, nil, "文章不存在")
		return
	}
	{
		// TODO 将文章存入redis供下次使用
		v, _ := json.Marshal(article)
		a.Redis.HSet(ctx, "Article", id, v)
	}
leep:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	articleVisit := model.ArticleVisit{
		ArticleId: article.ID,
		UserId:    user.ID,
	}

	// TODO 插入数据
	if err := a.DB.Create(&articleVisit).Error; err != nil {
		response.Fail(ctx, nil, "文章游览上传出错，数据库存储错误")
		return
	}

	// TODO 获取阅读人数
	last, _ := a.Redis.PFCount(ctx, "ArticleVisit"+id).Result()

	// TODO 添加入阅读库
	a.Redis.PFAdd(ctx, "UserVisit"+article.UserId.String(), user.ID.String())

	// TODO 获取新的阅读人数
	new, _ := a.Redis.PFCount(ctx, "ArticleVisit"+id).Result()

	// TODO 如果阅读人数有增加
	if new > last {
		a.Redis.ZIncrBy(ctx, "ArticleHot", 1.0, id)
		a.Redis.ZIncrBy(ctx, "UserVisit", 1.0, article.UserId.String())
	}

	response.Success(ctx, nil, "文章游览成功")
}

// @title    VisitNumber
// @description   游览文章数量
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) VisitNumber(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取阅读人数
	total, _ := a.Redis.PFCount(ctx, "ArticleVisit"+id).Result()

	response.Success(ctx, gin.H{"total": total}, "请求文章游览数目成功")
}

// @title    VisitList
// @description   游览文章列表
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) VisitList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var articleVisits []model.ArticleVisit

	var total int64

	// TODO 查看游览的数量
	a.DB.Where("article_id = (?)", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&articleVisits).Count(&total)

	response.Success(ctx, gin.H{"articleVisits": articleVisits, "total": total}, "查看成功")
}

// @title    Visits
// @description   游览文章列表
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) Visits(ctx *gin.Context) {
	// TODO 获取指定用户用户
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var articleVisits []model.ArticleVisit

	var total int64

	// TODO 查看收藏的数量
	a.DB.Where("user_id = (?)", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&articleVisits)

	a.DB.Where("user_id = (?)", id).Model(model.ArticleVisit{}).Count(&total)

	response.Success(ctx, gin.H{"articleVisits": articleVisits, "total": total}, "查看成功")
}

// @title    LabelCreate
// @description   标签创建
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) LabelCreate(ctx *gin.Context) {
	// TODO 获取指定文章
	id := ctx.Params.ByName("id")

	// TODO 获取标签
	label := ctx.Params.ByName("label")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看文章是否存在
	var article model.Article

	// TODO 先尝试在redis中寻找
	if ok, _ := a.Redis.HExists(ctx, "Article", id).Result(); ok {
		art, _ := a.Redis.HGet(ctx, "Article", id).Result()
		if json.Unmarshal([]byte(art), &article) == nil {
			goto leep
		} else {
			// TODO 解码失败，删除字段
			a.Redis.HDel(ctx, "Article", id)
		}
	}

	// TODO 查看文章是否在数据库中存在
	if a.DB.Where("id = (?)", id).First(&article).Error != nil {
		response.Fail(ctx, nil, "文章不存在")
		return
	}
	{
		// TODO 将文章存入redis供下次使用
		v, _ := json.Marshal(article)
		a.Redis.HSet(ctx, "Article", id, v)
	}
leep:

	// TODO 查看是否为文章作者
	if article.UserId != user.ID {
		response.Fail(ctx, nil, "不是文章作者，请勿非法操作")
		return
	}

	// TODO 创建标签
	articleLabel := model.ArticleLabel{
		Label:     label,
		ArticleId: article.ID,
	}

	// TODO 插入数据
	if err := a.DB.Create(&articleLabel).Error; err != nil {
		response.Fail(ctx, nil, "文章标签上传出错，数据验证有误")
		return
	}

	// TODO 解码失败，删除字段
	a.Redis.HDel(ctx, "ArticleLabel", id)

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    LabelDelete
// @description   标签删除
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) LabelDelete(ctx *gin.Context) {
	// TODO 获取指定文章
	id := ctx.Params.ByName("id")

	// TODO 获取标签
	label := ctx.Params.ByName("label")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看文章是否存在
	var article model.Article

	// TODO 先尝试在redis中寻找
	if ok, _ := a.Redis.HExists(ctx, "Article", id).Result(); ok {
		art, _ := a.Redis.HGet(ctx, "Article", id).Result()
		if json.Unmarshal([]byte(art), &article) == nil {
			goto leep
		} else {
			// TODO 解码失败，删除字段
			a.Redis.HDel(ctx, "Article", id)
		}
	}

	// TODO 查看文章是否在数据库中存在
	if a.DB.Where("id = (?)", id).First(&article).Error != nil {
		response.Fail(ctx, nil, "文章不存在")
		return
	}
	{
		// TODO 将文章存入redis供下次使用
		v, _ := json.Marshal(article)
		a.Redis.HSet(ctx, "Article", id, v)
	}
leep:

	// TODO 查看是否为文章作者
	if article.UserId != user.ID {
		response.Fail(ctx, nil, "不是文章作者，请勿非法操作")
		return
	}

	// TODO 删除文章标签
	if a.DB.Where("id = (?)", label).First(&model.ArticleLabel{}).Error != nil {
		response.Fail(ctx, nil, "标签不存在")
		return
	}

	a.DB.Where("id = (?)", label).Delete(&model.ArticleLabel{})

	// TODO 解码失败，删除字段
	a.Redis.HDel(ctx, "ArticleLabel", id)

	// TODO 成功
	response.Success(ctx, nil, "删除成功")
}

// @title    LabelShow
// @description   标签查看
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) LabelShow(ctx *gin.Context) {
	// TODO 获取指定文章
	id := ctx.Params.ByName("id")

	// TODO 查找数据
	var articleLabels []model.ArticleLabel
	// TODO 先尝试在redis中寻找
	if ok, _ := a.Redis.HExists(ctx, "ArticleLabel", id).Result(); ok {
		art, _ := a.Redis.HGet(ctx, "ArticleLabel", id).Result()
		if json.Unmarshal([]byte(art), &articleLabels) == nil {
			goto leap
		} else {
			// TODO 解码失败，删除字段
			a.Redis.HDel(ctx, "ArticleLabel", id)
		}
	}

	// TODO 在数据库中查找
	a.DB.Where("article_id = (?)", id).Find(&articleLabels)
	{
		// TODO 将文章标签存入redis供下次使用
		v, _ := json.Marshal(articleLabels)
		a.Redis.HSet(ctx, "ArticleLabel", id, v)
	}

leap:

	// TODO 成功
	response.Success(ctx, gin.H{"articleLabels": articleLabels}, "查看成功")
}

// @title    Search
// @description   文本搜索
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) Search(ctx *gin.Context) {
	// TODO 获取文本
	text := ctx.Params.ByName("text")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	var articles []model.Article

	// TODO 模糊匹配
	a.DB.Where("match(title,content,res_long,res_short) against((?) in boolean mode)", text+"*").Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&articles)

	// TODO 查看查询总数
	var total int64
	a.DB.Where("match(title,content,res_long,res_short) against((?) in boolean mode)", text+"*").Model(model.Article{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"articles": articles, "total": total}, "成功")
}

// @title    SearchLabel
// @description   指定标签的搜索
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) SearchLabel(ctx *gin.Context) {

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

	var articleLabels []model.ArticleLabel

	// TODO 进行标签匹配
	a.DB.Distinct("article_id").Where("label in (?)", requestLabels.Labels).Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&articleLabels)

	// TODO 查看查询总数
	var total int64
	a.DB.Distinct("article_id").Where("label in (?)", requestLabels.Labels).Model(model.ArticleLabel{}).Count(&total)

	// TODO 查找对应文章
	var articles []model.Article

	// TODO 此处将所有标签对应的id加入数组
	var articleIds []string
	for i := range articleLabels {
		articleIds = append(articleIds, articleLabels[i].ArticleId.String())
	}

	a.DB.Where("id in (?)", articleIds).Find(&articles)

	// TODO 返回数据
	response.Success(ctx, gin.H{"articles": articles, "total": total}, "成功")
}

// @title    SearchWithLabel
// @description   指定标签与文本的搜索
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) SearchWithLabel(ctx *gin.Context) {

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

	var articleLabels []model.ArticleLabel

	// TODO 进行标签匹配
	a.DB.Distinct("article_id").Where("label in (?)", requestLabels.Labels).Find(&articleLabels)

	// TODO 查找对应文章
	var articles []model.Article

	// TODO 此处将所有标签对应的id加入数组
	var articleIds []string
	for i := range articleLabels {
		articleIds = append(articleIds, articleLabels[i].ArticleId.String())
	}

	// TODO 模糊匹配
	a.DB.Where("id in (?) and match(title,content,res_long,res_short) against((?) in boolean mode)", articleIds, text+"*").Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&articles)

	// TODO 查看查询总数
	var total int64
	a.DB.Where("id in (?) and match(title,content,res_long,res_short) against((?) in boolean mode)", articleIds, text+"*").Model(model.Article{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"articles": articles, "total": total}, "成功")
}

// @title    NewArticleController
// @description   新建一个IArticleController
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   IArticleController		返回一个IArticleController用于调用各种函数
func NewArticleController() IArticleController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.Article{})
	db.AutoMigrate(model.ArticleCollect{})
	db.AutoMigrate(model.ArticleLike{})
	db.AutoMigrate(model.ArticleVisit{})
	db.AutoMigrate(model.ArticleLabel{})
	return ArticleController{DB: db, Redis: redis}
}
