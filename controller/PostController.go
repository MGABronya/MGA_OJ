// @Title  PostController
// @Description  该文件提供关于操作题解的各种方法
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

// IPostController			定义了题解类接口
type IPostController interface {
	Interface.RestInterface    // 包含增删查改功能
	Interface.LikeInterface    // 包含点赞功能
	Interface.CollectInterface // 包含收藏功能
	Interface.VisitInterface   // 包含游览功能
	UserList(ctx *gin.Context) // 查看指定用户的题解
}

// PostController			定义了题解工具类
type PostController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Create
// @description   创建一篇题解
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) Create(ctx *gin.Context) {
	var requestPost vo.PostRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestPost); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 查找对应题目
	id := ctx.Params.ByName("id")

	var problem model.Problem

	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "Problem", id).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Problem", id).Result()
		if json.Unmarshal([]byte(cate), &problem) == nil {
			// TODO 跳过数据库搜寻problem过程
			goto leep
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Problem", id)
		}
	}

	// TODO 查看题目是否在数据库中存在
	if p.DB.Where("id = ?", id).First(&problem).Error != nil {
		response.Fail(ctx, nil, "题目不存在")
		return
	}
	// TODO 将题目存入redis供下次使用
	{
		v, _ := json.Marshal(problem)
		p.Redis.HSet(ctx, "Problem", id, v)
	}

leep:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 创建题解
	var post = model.Post{
		UserId:    user.ID,
		ProblemId: problem.ID,
		Title:     requestPost.Title,
		Content:   requestPost.Content,
		Reslong:   requestPost.Reslong,
		Resshort:  requestPost.Resshort,
	}

	// TODO 插入数据
	if err := p.DB.Create(&post).Error; err != nil {
		response.Fail(ctx, nil, "题解上传出错，数据库存储错误")
		return
	}

	// TODO 成功
	response.Success(ctx, gin.H{"post": post}, "创建成功")
}

// @title    Update
// @description   更新一篇题解的内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) Update(ctx *gin.Context) {
	var requestPost vo.PostRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestPost); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应题解
	id := ctx.Params.ByName("id")

	var post model.Post

	if p.DB.Where("id = ?", id).First(&post) != nil {
		response.Fail(ctx, nil, "题解不存在")
		return
	}

	// TODO 查看是否是用户作者
	if user.ID != post.UserId {
		response.Fail(ctx, nil, "不是题解作者，无法修改题解")
		return
	}

	// TODO 更新题解内容
	p.DB.Where("id = ?", id).Updates(requestPost)

	// TODO 移除损坏数据
	p.Redis.HDel(ctx, "Post", id)

	// TODO 成功
	response.Success(ctx, nil, "更新成功")
}

// @title    Show
// @description   查看一篇题解的内容
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var post model.Post

	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "Post", id).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Post", id).Result()
		if json.Unmarshal([]byte(cate), &post) == nil {
			response.Success(ctx, gin.H{"post": post}, "成功")
			return
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Post", id)
		}
	}

	// TODO 查看题解是否在数据库中存在
	if p.DB.Where("id = ?", id).First(&post).Error != nil {
		response.Fail(ctx, nil, "题解不存在")
		return
	}

	response.Success(ctx, gin.H{"post": post}, "成功")

	// TODO 将题解存入redis供下次使用
	v, _ := json.Marshal(post)
	p.Redis.HSet(ctx, "Post", id, v)
}

// @title    Delete
// @description   删除一篇题解
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) Delete(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	var post model.Post

	// TODO 查看题解是否存在
	if p.DB.Where("id = ?", id).First(&post).Error != nil {
		response.Fail(ctx, nil, "题解不存在")
		return
	}

	// TODO 判断当前用户是否为题解的作者
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有操作题解的权力
	if user.ID != post.UserId && user.Level < 4 {
		response.Fail(ctx, nil, "题解不属于您，请勿非法操作")
		return
	}

	// TODO 删除题解
	p.DB.Delete(&post)

	// TODO 移除损坏数据
	p.Redis.HDel(ctx, "Post", id)

	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description   获取多篇题解
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) PageList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 分页
	var posts []model.Post

	// TODO 查找所有分页中可见的条目
	p.DB.Where("problem_id = ?", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&posts)

	var total int64
	p.DB.Where("problem_id = ?", id).Model(model.Post{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"posts": posts, "total": total}, "成功")
}

// @title    UserList
// @description   获取多篇指定用户的题解
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) UserList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 分页
	var posts []model.Post

	// TODO 查找所有分页中可见的条目
	p.DB.Where("user_id = ?", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&posts)

	var total int64
	p.DB.Where("user_id = ?", id).Model(model.Post{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"posts": posts, "total": total}, "成功")
}

// @title    Like
// @description   点赞或点踩
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) Like(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	var post model.Post

	// TODO 查看题解是否存在
	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "Post", id).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Post", id).Result()
		if json.Unmarshal([]byte(cate), &post) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Post", id)
		}
	}

	// TODO 查看题解是否在数据库中存在
	if p.DB.Where("id = ?", id).First(&post).Error != nil {
		response.Fail(ctx, nil, "题解不存在")
		return
	}
	{
		// TODO 将题解存入redis供下次使用
		v, _ := json.Marshal(post)
		p.Redis.HSet(ctx, "Post", id, v)
	}
leep:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有点赞或者点踩
	if p.DB.Where("user_id = ? and post_id = ?", user.ID, id).Update("like", like).Error != nil {
		// TODO 插入数据
		postLike := model.PostLike{
			PostId: post.ID,
			UserId: user.ID,
			Like:   like,
		}
		if err := p.DB.Create(&postLike).Error; err != nil {
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
func (p PostController) CancelLike(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 取消点赞或者点踩
	p.DB.Where("user_id = ? and post_id = ?", user.ID, id).Delete(&model.PostLike{})

	response.Success(ctx, nil, "取消成功")
}

// @title    LikeNumber
// @description   点赞或点踩的数量
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) LikeNumber(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	var total int64

	// TODO 查看点赞或者点踩的数量
	p.DB.Where("post_id = ? and like = ?", id, like).Count(&total)

	response.Success(ctx, gin.H{"total": total}, "查看成功")
}

// @title    LikeList
// @description   点赞或点踩的列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) LikeList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var postLikes []model.PostLike

	var total int64

	// TODO 查看点赞或者点踩的数量
	p.DB.Where("post_id = ? and like = ?", id, like).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&postLikes).Count(&total)

	response.Success(ctx, gin.H{"postLikes": postLikes, "total": total}, "查看成功")
}

// @title    LikeShow
// @description   查看用户点赞状态
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) LikeShow(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var postLike model.PostLike

	// TODO 查看点赞状态
	if p.DB.Where("user_id = ? and post_id = ?", user.ID, id).First(&postLike).Error != nil {
		response.Success(ctx, gin.H{"like": 0}, "暂无状态")
		return
	}

	if postLike.Like {
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
func (p PostController) Likes(ctx *gin.Context) {

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取指定用户用户
	id := ctx.Params.ByName("id")

	// TODO 分页
	var postLikes []model.PostLike

	var total int64

	// TODO 查看点赞或者点踩的数量
	p.DB.Where("user_id = ? and like = ?", id, like).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&postLikes).Count(&total)

	response.Success(ctx, gin.H{"postLikes": postLikes, "total": total}, "查看成功")
}

// @title    Collect
// @description   收藏
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) Collect(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	var post model.Post

	// TODO 查看题解是否存在
	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "Post", id).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Post", id).Result()
		if json.Unmarshal([]byte(cate), &post) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Post", id)
		}
	}

	// TODO 查看题解是否在数据库中存在
	if p.DB.Where("id = ?", id).First(&post).Error != nil {
		response.Fail(ctx, nil, "题解不存在")
		return
	}
	{
		// TODO 将题解存入redis供下次使用
		v, _ := json.Marshal(post)
		p.Redis.HSet(ctx, "Post", id, v)
	}
leep:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有收藏
	if p.DB.Where("user_id = ? and post_id = ?", user.ID, post.ID).First(&model.PostCollect{}).Error != nil {
		postCollect := model.PostCollect{
			PostId: post.ID,
			UserId: user.ID,
		}
		// TODO 插入数据
		if err := p.DB.Create(&postCollect).Error; err != nil {
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
func (p PostController) CancelCollect(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有收藏
	if p.DB.Where("user_id = ? and post_id = ?", user.ID, id).First(&model.PostCollect{}).Error != nil {
		response.Fail(ctx, nil, "未收藏")
		return
	} else {
		p.DB.Where("user_id = ? and post_id = ?", user.ID, id).Delete(&model.PostCollect{})
		response.Success(ctx, nil, "取消收藏成功")
		return
	}
}

// @title    CollectShow
// @description   查看收藏状态
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) CollectShow(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有收藏
	if p.DB.Where("user_id = ? and post_id = ?", user.ID, id).First(&model.PostCollect{}).Error != nil {
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
func (p PostController) CollectList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var postCollects []model.PostCollect

	var total int64

	// TODO 查看收藏的数量
	p.DB.Where("post_id = ?", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&postCollects).Count(&total)

	response.Success(ctx, gin.H{"postCollects": postCollects, "total": total}, "查看成功")
}

// @title    CollectNumber
// @description   查看收藏用户数量
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) CollectNumber(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	var total int64

	// TODO 查看收藏的数量
	p.DB.Where("post_id = ?", id).Count(&total)

	response.Success(ctx, gin.H{"total": total}, "查看成功")
}

// @title    Collects
// @description   查看用户收藏夹
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) Collects(ctx *gin.Context) {

	// TODO 获取指定用户用户
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var postCollects []model.PostCollect

	var total int64

	// TODO 查看收藏的数量
	p.DB.Where("user_id = ?", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&postCollects).Count(&total)

	response.Success(ctx, gin.H{"postCollects": postCollects, "total": total}, "查看成功")
}

// @title    Visit
// @description   游览题解
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) Visit(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	var post model.Post

	// TODO 查看题解是否存在
	// TODO 先看redis中是否存在
	if ok, _ := p.Redis.HExists(ctx, "Post", id).Result(); ok {
		cate, _ := p.Redis.HGet(ctx, "Post", id).Result()
		if json.Unmarshal([]byte(cate), &post) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			p.Redis.HDel(ctx, "Post", id)
		}
	}

	// TODO 查看题解是否在数据库中存在
	if p.DB.Where("id = ?", id).First(&post).Error != nil {
		response.Fail(ctx, nil, "题解不存在")
		return
	}
	{
		// TODO 将题解存入redis供下次使用
		v, _ := json.Marshal(post)
		p.Redis.HSet(ctx, "Post", id, v)
	}
leep:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	postVisit := model.PostVisit{
		PostId: post.ID,
		UserId: user.ID,
	}

	// TODO 插入数据
	if err := p.DB.Create(&postVisit).Error; err != nil {
		response.Fail(ctx, nil, "题解游览上传出错，数据库存储错误")
		return
	}

	response.Success(ctx, nil, "题解游览成功")
}

// @title    VisitNumber
// @description   游览题解数量
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) VisitNumber(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获得游览总数
	var total int64

	p.DB.Where("post_id = ?", id).Count(&total)

	response.Success(ctx, gin.H{"total": total}, "请求题解游览数目成功")
}

// @title    VisitList
// @description   游览题解列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) VisitList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var postVisits []model.PostVisit

	var total int64

	// TODO 查看收藏的数量
	p.DB.Where("post_id = ?", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&postVisits).Count(&total)

	response.Success(ctx, gin.H{"postVisits": postVisits, "total": total}, "查看成功")
}

// @title    Visits
// @description   游览题解列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) Visits(ctx *gin.Context) {
	// TODO 获取指定用户用户
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var postVisits []model.PostVisit

	var total int64

	// TODO 查看收藏的数量
	p.DB.Where("user_id = ?", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&postVisits).Count(&total)

	response.Success(ctx, gin.H{"postVisits": postVisits, "total": total}, "查看成功")
}

// @title    NewPostController
// @description   新建一个IPostController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   IPostController		返回一个IPostController用于调用各种函数
func NewPostController() IPostController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.Post{})
	db.AutoMigrate(model.PostCollect{})
	db.AutoMigrate(model.PostLike{})
	db.AutoMigrate(model.PostVisit{})
	return PostController{DB: db, Redis: redis}
}
