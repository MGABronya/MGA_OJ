// @Title  ThreadController
// @Description  该文件提供关于操作题解的回复的各种方法
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

// IThreadController			定义了题解的回复类接口
type IThreadController interface {
	Interface.RestInterface    // 包含增删查改功能
	Interface.LikeInterface    // 包含点赞功能
	Interface.HotInterface     // 包含热度功能
	UserList(ctx *gin.Context) // 指定用户的题解回复
}

// ThreadController			定义了题解的回复工具类
type ThreadController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Create
// @description   创建一篇题解的回复
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t ThreadController) Create(ctx *gin.Context) {
	var requestThread vo.ThreadRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestThread); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 查找对应题解
	id := ctx.Params.ByName("id")

	var post model.Post

	// TODO 查看题解是否存在
	// TODO 先看redis中是否存在
	if ok, _ := t.Redis.HExists(ctx, "Post", id).Result(); ok {
		cate, _ := t.Redis.HGet(ctx, "Post", id).Result()
		if json.Unmarshal([]byte(cate), &post) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			t.Redis.HDel(ctx, "Post", id)
		}
	}

	// TODO 查看题解是否在数据库中存在
	if t.DB.Where("id = (?)", id).First(&post).Error != nil {
		response.Fail(ctx, nil, "题解不存在")
		return
	}
	{
		// TODO 将题解存入redis供下次使用
		v, _ := json.Marshal(post)
		t.Redis.HSet(ctx, "Post", id, v)
	}
leep:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 创建题解的回复
	var thread = model.Thread{
		UserId:   user.ID,
		PostId:   post.ID,
		Content:  requestThread.Content,
		ResLong:  requestThread.ResLong,
		ResShort: requestThread.ResShort,
	}

	// TODO 插入数据
	if err := t.DB.Create(&thread).Error; err != nil {
		response.Fail(ctx, nil, "题解的回复上传出错，数据库存储错误")
		return
	}

	// TODO 创建热度
	t.Redis.ZAdd(ctx, "ThreadHot"+post.ID.String(), redis.Z{Member: thread.ID.String(), Score: 100 + float64(time.Now().Unix()/86400)})

	// TODO 成功
	response.Success(ctx, gin.H{"thread": thread}, "创建成功")
}

// @title    Update
// @description   更新一篇题解的回复的内容
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t ThreadController) Update(ctx *gin.Context) {
	var requestThread vo.ThreadRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestThread); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应题解的回复
	id := ctx.Params.ByName("id")

	var Thread model.Thread

	if t.DB.Where("id = (?)", id).First(&Thread).Error != nil {
		response.Fail(ctx, nil, "题解的回复不存在")
		return
	}

	// TODO 查看是否是用户作者
	if user.ID != Thread.UserId {
		response.Fail(ctx, nil, "不是题解的回复作者，无法修改题解的回复")
		return
	}

	threadUpdate := model.Thread{
		Content:  requestThread.Content,
		ResLong:  requestThread.ResLong,
		ResShort: requestThread.ResShort,
	}

	// TODO 更新题解的回复内容
	t.DB.Where("id = (?)", id).Updates(threadUpdate)

	// TODO 移除损坏数据
	t.Redis.HDel(ctx, "Thread", id)

	// TODO 成功
	response.Success(ctx, nil, "更新成功")
}

// @title    Show
// @description   查看一篇题解的回复的内容
// @auth      MGAronya       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t ThreadController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var thread model.Thread

	// TODO 先看redis中是否存在
	if ok, _ := t.Redis.HExists(ctx, "Thread", id).Result(); ok {
		cate, _ := t.Redis.HGet(ctx, "Thread", id).Result()
		if json.Unmarshal([]byte(cate), &thread) == nil {
			response.Success(ctx, gin.H{"thread": thread}, "成功")
			return
		} else {
			// TODO 移除损坏数据
			t.Redis.HDel(ctx, "Thread", id)
		}
	}

	// TODO 查看题解的回复是否在数据库中存在
	if t.DB.Where("id = (?)", id).First(&thread).Error != nil {
		response.Fail(ctx, nil, "题解的回复不存在")
		return
	}

	response.Success(ctx, gin.H{"thread": thread}, "成功")

	// TODO 将提交存入redis供下次使用
	v, _ := json.Marshal(thread)
	t.Redis.HSet(ctx, "Thread", id, v)
}

// @title    Delete
// @description   删除一篇题解的回复
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t ThreadController) Delete(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	var thread model.Thread

	// TODO 查看题解的回复是否存在
	if t.DB.Where("id = (?)", id).First(&thread).Error != nil {
		response.Fail(ctx, nil, "题解的回复不存在")
		return
	}

	// TODO 判断当前用户是否为题解的回复的作者
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有操作题解的回复的权力
	if user.ID != thread.UserId && user.Level < 4 {
		response.Fail(ctx, nil, "题解的回复不属于您，请勿非法操作")
		return
	}

	var total int64

	// TODO 查看点赞的数量
	t.DB.Where("thread_id = (?) and like = true", id).Model(model.ThreadLike{}).Count(&total)
	t.Redis.ZIncrBy(ctx, "UserLike", -float64(total), thread.UserId.String())

	// TODO 查看点踩的数量
	t.DB.Where("thread_id = (?) and like = false", id).Model(model.ThreadLike{}).Count(&total)
	t.Redis.ZIncrBy(ctx, "UserUnLike", -float64(total), thread.UserId.String())

	// TODO 删除题解的回复
	t.DB.Delete(&thread)

	// TODO 移除损坏数据
	t.Redis.HDel(ctx, "Thread", id)

	// TODO 移除热度
	t.Redis.ZRem(ctx, "ThreadHot"+thread.PostId.String(), thread.ID.String())

	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description   获取多篇题解的回复
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t ThreadController) PageList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 分页
	var threads []model.Thread

	// TODO 查找所有分页中可见的条目
	t.DB.Where("post_id = (?)", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&threads)

	var total int64
	t.DB.Where("post_id = (?)", id).Model(model.Thread{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"threads": threads, "total": total}, "成功")
}

// @title    UserList
// @description   获取多篇指定用户的题解的回复
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t ThreadController) UserList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 分页
	var threads []model.Thread

	// TODO 查找所有分页中可见的条目
	t.DB.Where("user_id = (?)", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&threads)

	var total int64
	t.DB.Where("user_id = (?)", id).Model(model.Thread{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"threads": threads, "total": total}, "成功")
}

// @title    HotRanking
// @description   根据热度排行获取多篇讨论
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t ThreadController) HotRanking(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	// TODO 查找所有分页中可见的条目
	threads, err := t.Redis.ZRevRangeWithScores(ctx, "ThreadHot"+id, int64(pageNum-1)*int64(pageSize), int64(pageNum-1)*int64(pageSize)+int64(pageSize)-1).Result()

	if err != nil {
		response.Fail(ctx, nil, "获取失败")
	}

	for i := range threads {
		threads[i].Score -= float64(time.Now().Unix() / 86400)
	}

	// TODO 将redis中的数据取出
	total, _ := t.Redis.ZCard(ctx, "ThreadHot"+id).Result()

	// TODO 返回数据
	response.Success(ctx, gin.H{"threads": threads, "total": total}, "成功")
}

// @title    Like
// @description   点赞或点踩
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t ThreadController) Like(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	var thread model.Thread

	// TODO 先看redis中是否存在
	if ok, _ := t.Redis.HExists(ctx, "Thread", id).Result(); ok {
		cate, _ := t.Redis.HGet(ctx, "Thread", id).Result()
		if json.Unmarshal([]byte(cate), &thread) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			t.Redis.HDel(ctx, "Thread", id)
		}
	}

	// TODO 查看题解的回复是否在数据库中存在
	if t.DB.Where("id = (?)", id).First(&thread).Error != nil {
		response.Fail(ctx, nil, "题解的回复不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(thread)
		t.Redis.HSet(ctx, "Thread", id, v)
	}
leep:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var threadLike model.ThreadLike
	// TODO 如果没有点赞或者点踩
	if t.DB.Where("user_id = (?) and thread_id = (?)", user.ID, id).First(&threadLike).Error != nil {
		// TODO 插入数据
		threadLike = model.ThreadLike{
			ThreadId: thread.ID,
			UserId:   user.ID,
			Like:     like,
		}
		if err := t.DB.Create(&threadLike).Error; err != nil {
			response.Fail(ctx, nil, "点赞出错，数据库存储错误")
			return
		}
	} else {
		// TODO 热度计算
		if threadLike.Like {
			t.Redis.ZIncrBy(ctx, "ThreadHot"+thread.PostId.String(), -10.0, thread.ID.String())
			t.Redis.ZIncrBy(ctx, "UserLike", -1, thread.UserId.String())
		} else {
			t.Redis.ZIncrBy(ctx, "ThreadHot"+thread.PostId.String(), 10.0, thread.ID.String())
			t.Redis.ZIncrBy(ctx, "UserUnLike", -1, thread.UserId.String())
		}
		t.DB.Where("user_id = (?) and thread_id = (?)", user.ID, id).Model(&model.ThreadLike{}).Update("like", like)
	}

	// TODO 热度计算
	if like {
		t.Redis.ZIncrBy(ctx, "ThreadHot"+thread.PostId.String(), 10.0, thread.ID.String())
		t.Redis.ZIncrBy(ctx, "UserLike", 1, thread.UserId.String())
	} else {
		t.Redis.ZIncrBy(ctx, "ThreadHot"+thread.PostId.String(), -10.0, thread.ID.String())
		t.Redis.ZIncrBy(ctx, "UserUnLike", 1, thread.UserId.String())
	}

	Handle.Behaviors["Likes"].PublishBehavior(1, thread.UserId)
	response.Success(ctx, nil, "点赞成功")
}

// @title    CancelLike
// @description   取消点赞或者点踩
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t ThreadController) CancelLike(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var thread model.Thread

	// TODO 先看redis中是否存在
	if ok, _ := t.Redis.HExists(ctx, "Thread", id).Result(); ok {
		cate, _ := t.Redis.HGet(ctx, "Thread", id).Result()
		if json.Unmarshal([]byte(cate), &thread) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			t.Redis.HDel(ctx, "Thread", id)
		}
	}

	// TODO 查看题解的回复是否在数据库中存在
	if t.DB.Where("id = (?)", id).First(&thread).Error != nil {
		response.Fail(ctx, nil, "题解的回复不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(thread)
		t.Redis.HSet(ctx, "Thread", id, v)
	}
leep:
	// TODO 查看是否已经点赞或者点踩
	var threadLike model.ThreadLike
	if t.DB.Where("user_id = (?) and thread_id = (?)", user.ID, id).First(&threadLike).Error != nil {
		response.Fail(ctx, nil, "未点赞或点踩")
		return
	}

	// TODO 热度计算
	if threadLike.Like {
		t.Redis.ZIncrBy(ctx, "ThreadHot"+thread.PostId.String(), -10.0, thread.ID.String())
		t.Redis.ZIncrBy(ctx, "UserLike", -1, thread.UserId.String())
	} else {
		t.Redis.ZIncrBy(ctx, "ThreadHot"+thread.PostId.String(), 10.0, thread.ID.String())
		t.Redis.ZIncrBy(ctx, "UserUnLike", -1, thread.UserId.String())
	}

	// TODO 取消点赞或者点踩
	t.DB.Where("user_id = (?) and thread_id = (?)", user.ID, id).Delete(&model.ThreadLike{})

	Handle.Behaviors["Likes"].PublishBehavior(-1, thread.UserId)
	response.Success(ctx, nil, "取消成功")
}

// @title    LikeNumber
// @description   点赞或点踩的数量
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t ThreadController) LikeNumber(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	var total int64

	// TODO 查看点赞或者点踩的数量
	t.DB.Where("thread_id = (?) and `like` = ?", id, like).Model(model.ThreadLike{}).Count(&total)

	response.Success(ctx, gin.H{"total": total}, "查看成功")
}

// @title    LikeList
// @description   点赞或点踩的列表
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t ThreadController) LikeList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var threadLikes []model.ThreadLike

	var total int64

	// TODO 查看点赞或者点踩的数量
	t.DB.Where("thread_id = (?) and `like` = ?", id, like).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&threadLikes)

	t.DB.Where("thread_id = (?) and `like` = ?", id, like).Model(model.ThreadLike{}).Count(&total)

	response.Success(ctx, gin.H{"threadLikes": threadLikes, "total": total}, "查看成功")
}

// @title    LikeShow
// @description   查看用户点赞状态
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t ThreadController) LikeShow(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var threadLike model.ThreadLike

	// TODO 查看点赞状态
	if t.DB.Where("user_id = (?) and thread_id = (?)", user.ID, id).First(&threadLike).Error != nil {
		response.Success(ctx, gin.H{"like": 0}, "暂无状态")
		return
	}

	if threadLike.Like {
		response.Success(ctx, gin.H{"like": 1}, "已点赞")
	} else {
		response.Success(ctx, gin.H{"like": -1}, "已点踩")
	}

}

// @title    Likes
// @description   查看用户点赞列表
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t ThreadController) Likes(ctx *gin.Context) {

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取指定用户用户
	id := ctx.Params.ByName("id")

	// TODO 分页
	var threadLikes []model.ThreadLike

	var total int64

	// TODO 查看点赞或者点踩的数量
	t.DB.Where("user_id = (?) and `like` = ?", id, like).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&threadLikes)

	t.DB.Where("user_id = (?) and `like` = ?", id, like).Model(model.ThreadLike{}).Count(&total)

	response.Success(ctx, gin.H{"threadLikes": threadLikes, "total": total}, "查看成功")
}

// @title    NewThreadController
// @description   新建一个IThreadController
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   IThreadController		返回一个IThreadController用于调用各种函数
func NewThreadController() IThreadController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.Thread{})
	db.AutoMigrate(model.ThreadLike{})
	return ThreadController{DB: db, Redis: redis}
}
