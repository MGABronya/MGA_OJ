// @Title  TopicController
// @Description  该文件提供关于操作主题的各种方法
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
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// ITopicController			定义了主题类接口
type ITopicController interface {
	Interface.RestInterface       // 包含增删查改功能
	Interface.LikeInterface       // 包含点赞功能
	Interface.CollectInterface    // 包含收藏功能
	Interface.VisitInterface      // 包含游览功能
	Interface.LabelInterface      // 包含标签功能
	Interface.SearchInterface     // 包含搜索功能
	UserList(ctx *gin.Context)    // 查看用户的主题
	ProblemList(ctx *gin.Context) // 查看主题的题目列表
}

// TopicController			定义了主题工具类
type TopicController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Create
// @description   创建一篇主题
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) Create(ctx *gin.Context) {
	var requestTopic vo.TopicRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestTopic); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 创建主题
	topic := model.Topic{
		Title:    requestTopic.Title,
		Content:  requestTopic.Content,
		Reslong:  requestTopic.Reslong,
		Resshort: requestTopic.Resshort,
		UserId:   user.ID,
	}

	// TODO 插入数据
	if err := t.DB.Create(&topic).Error; err != nil {
		response.Fail(ctx, nil, "主题上传出错，数据验证有误")
		return
	}

	// TODO 插入相关题目
	for _, v := range requestTopic.Problems {
		id := fmt.Sprint(v)
		var problem model.Problem
		// TODO 先看redis中是否存在
		if ok, _ := t.Redis.HExists(ctx, "Problem", id).Result(); ok {
			cate, _ := t.Redis.HGet(ctx, "Problem", id).Result()
			if json.Unmarshal([]byte(cate), &problem) == nil {
				// TODO 跳过数据库搜寻problem过程
				goto leep
			} else {
				// TODO 移除损坏数据
				t.Redis.HDel(ctx, "Problem", id)
			}
		}

		// TODO 查看题目是否在数据库中存在
		if t.DB.Where("id = ?", id).First(&problem).Error != nil {
			response.Fail(ctx, nil, "题目不存在")
			return
		}
		// TODO 将题目存入redis供下次使用
		{
			v, _ := json.Marshal(problem)
			t.Redis.HSet(ctx, "Problem", id, v)
		}
	leep:
		problemList := model.ProblemList{
			TopicId:   topic.ID,
			ProblemId: v,
		}
		if t.DB.Create(&problemList).Error != nil {
			response.Fail(ctx, nil, "题目上传出错，数据验证有误")
			return
		}
	}

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    Update
// @description   更新一篇主题的内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) Update(ctx *gin.Context) {
	var requestTopic vo.TopicRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestTopic); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应主题
	id := ctx.Params.ByName("id")

	var topic model.Topic

	if t.DB.Where("id = ?", id).First(&topic) != nil {
		response.Fail(ctx, nil, "主题不存在")
		return
	}

	// TODO 查看是否是用户作者
	if user.ID != topic.UserId {
		response.Fail(ctx, nil, "不是主题作者，无法修改主题")
		return
	}

	// TODO 创建主题
	topic = model.Topic{
		Title:    requestTopic.Title,
		Content:  requestTopic.Content,
		Reslong:  requestTopic.Reslong,
		Resshort: requestTopic.Resshort,
		UserId:   user.ID,
	}

	// TODO 更新主题内容
	t.DB.Where("id = ?", id).Updates(topic)

	// TODO 移除损坏数据
	t.Redis.HDel(ctx, "Topic", id)

	if len(requestTopic.Problems) != 0 {
		t.DB.Where("topic_id = ?", id).Delete(&model.ProblemList{})
		// TODO 插入相关题目
		for _, v := range requestTopic.Problems {
			id := fmt.Sprint(v)
			var problem model.Problem
			// TODO 先看redis中是否存在
			if ok, _ := t.Redis.HExists(ctx, "Problem", id).Result(); ok {
				cate, _ := t.Redis.HGet(ctx, "Problem", id).Result()
				if json.Unmarshal([]byte(cate), &problem) == nil {
					// TODO 跳过数据库搜寻problem过程
					goto leep
				} else {
					// TODO 移除损坏数据
					t.Redis.HDel(ctx, "Problem", id)
				}
			}

			// TODO 查看题目是否在数据库中存在
			if t.DB.Where("id = ?", id).First(&problem).Error != nil {
				response.Fail(ctx, nil, "题目不存在")
				return
			}
			// TODO 将题目存入redis供下次使用
			{
				v, _ := json.Marshal(problem)
				t.Redis.HSet(ctx, "Problem", id, v)
			}
		leep:
			problemList := model.ProblemList{
				TopicId:   topic.ID,
				ProblemId: v,
			}
			if t.DB.Create(&problemList).Error != nil {
				response.Fail(ctx, nil, "题目上传出错，数据验证有误")
				return
			}
		}
	}

	// TODO 成功
	response.Success(ctx, nil, "更新成功")
}

// @title    Show
// @description   查看一篇主题的内容
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var topic model.Topic

	// TODO 先看redis中是否存在
	if ok, _ := t.Redis.HExists(ctx, "Topic", id).Result(); ok {
		cate, _ := t.Redis.HGet(ctx, "Topic", id).Result()
		if json.Unmarshal([]byte(cate), &topic) == nil {
			response.Success(ctx, gin.H{"topic": topic}, "成功")
			return
		} else {
			// TODO 移除损坏数据
			t.Redis.HDel(ctx, "Topic", id)
		}
	}

	// TODO 查看主题是否在数据库中存在
	if t.DB.Where("id = ?", id).First(&topic).Error != nil {
		response.Fail(ctx, nil, "主题不存在")
		return
	}

	response.Success(ctx, gin.H{"topic": topic}, "成功")

	// TODO 将提交存入redis供下次使用
	v, _ := json.Marshal(topic)
	t.Redis.HSet(ctx, "Topic", id, v)
}

// @title    Delete
// @description   删除一篇主题
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) Delete(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	var topic model.Topic

	// TODO 查看主题是否存在
	if t.DB.Where("id = ?", id).First(&topic).Error != nil {
		response.Fail(ctx, nil, "主题不存在")
		return
	}

	// TODO 判断当前用户是否为主题的作者
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有操作主题的权力
	if user.ID != topic.UserId && user.Level < 4 {
		response.Fail(ctx, nil, "主题不属于您，请勿非法操作")
		return
	}

	// TODO 删除主题
	t.DB.Delete(&topic)

	// TODO 移除损坏数据
	t.Redis.HDel(ctx, "Topic", id)

	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description   获取多篇主题
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) PageList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var topics []model.Topic

	// TODO 查找所有分页中可见的条目
	t.DB.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&topics)

	var total int64
	t.DB.Model(model.Topic{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"topics": topics, "total": total}, "成功")
}

// @title    UserList
// @description   获取指定用户的多篇主题
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) UserList(ctx *gin.Context) {

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 分页
	var topics []model.Topic

	// TODO 查找所有分页中可见的条目
	t.DB.Where("user_id = ?", id).Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&topics)

	var total int64
	t.DB.Model(model.Topic{}).Where("user_id = ?", id).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"topics": topics, "total": total}, "成功")
}

// @title    ProblemList
// @description   获取指定主题的多篇题目
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) ProblemList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 分页
	var problemLists []model.ProblemList

	// TODO 查找所有分页中可见的条目
	t.DB.Where("topic_id = ?", id).Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&problemLists)

	var total int64
	t.DB.Model(model.Topic{}).Where("topic_id = ?", id).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"problemLists": problemLists, "total": total}, "成功")
}

// @title    Like
// @description   点赞或点踩
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) Like(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	var topic model.Topic

	// TODO 查看主题是否存在
	// TODO 先看redis中是否存在
	if ok, _ := t.Redis.HExists(ctx, "Topic", id).Result(); ok {
		cate, _ := t.Redis.HGet(ctx, "Topic", id).Result()
		if json.Unmarshal([]byte(cate), &topic) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			t.Redis.HDel(ctx, "Topic", id)
		}
	}

	// TODO 查看主题是否在数据库中存在
	if t.DB.Where("id = ?", id).First(&topic).Error != nil {
		response.Fail(ctx, nil, "主题不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(topic)
		t.Redis.HSet(ctx, "Topic", id, v)
	}
leep:
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有点赞或者点踩
	if t.DB.Where("user_id = ? and topic_id = ?", user.ID, id).Update("like", like).Error != nil {
		// TODO 插入数据
		topicLike := model.TopicLike{
			TopicId: topic.ID,
			UserId:  user.ID,
			Like:    like,
		}
		if err := t.DB.Create(&topicLike).Error; err != nil {
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
func (t TopicController) CancelLike(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 取消点赞或者点踩
	t.DB.Where("user_id = ? and topic_id = ?", user.ID, id).Delete(&model.TopicLike{})

	response.Success(ctx, nil, "取消成功")
}

// @title    LikeNumber
// @description   点赞或点踩的数量
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) LikeNumber(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	var total int64

	// TODO 查看点赞或者点踩的数量
	t.DB.Where("topic_id = ? and like = ?", id, like).Count(&total)

	response.Success(ctx, gin.H{"total": total}, "查看成功")
}

// @title    LikeList
// @description   点赞或点踩的列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) LikeList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var topicLikes []model.TopicLike

	var total int64

	// TODO 查看点赞或者点踩的数量
	t.DB.Where("topic_id = ? and like = ?", id, like).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&topicLikes).Count(&total)

	response.Success(ctx, gin.H{"topicLikes": topicLikes, "total": total}, "查看成功")
}

// @title    LikeShow
// @description   查看用户点赞状态
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) LikeShow(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var topicLike model.TopicLike

	// TODO 查看点赞状态
	if t.DB.Where("user_id = ? and topic_id = ?", user.ID, id).First(&topicLike).Error != nil {
		response.Success(ctx, gin.H{"like": 0}, "暂无状态")
		return
	}

	if topicLike.Like {
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
func (t TopicController) Likes(ctx *gin.Context) {

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 分页
	var topicLikes []model.TopicLike

	var total int64

	// TODO 查看点赞或者点踩的数量
	t.DB.Where("user_id = ? and like = ?", user.ID, like).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&topicLikes).Count(&total)

	response.Success(ctx, gin.H{"topicLikes": topicLikes, "total": total}, "查看成功")
}

// @title    Collect
// @description   收藏
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) Collect(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	var topic model.Topic

	// TODO 查看主题是否存在
	// TODO 先看redis中是否存在
	if ok, _ := t.Redis.HExists(ctx, "Topic", id).Result(); ok {
		cate, _ := t.Redis.HGet(ctx, "Topic", id).Result()
		if json.Unmarshal([]byte(cate), &topic) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			t.Redis.HDel(ctx, "Topic", id)
		}
	}

	// TODO 查看主题是否在数据库中存在
	if t.DB.Where("id = ?", id).First(&topic).Error != nil {
		response.Fail(ctx, nil, "主题不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(topic)
		t.Redis.HSet(ctx, "Topic", id, v)
	}
leep:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有收藏
	if t.DB.Where("user_id = ? and topic_id = ?", user.ID, topic.ID).First(&model.TopicCollect{}).Error != nil {
		topicCollect := model.TopicCollect{
			TopicId: topic.ID,
			UserId:  user.ID,
		}
		// TODO 插入数据
		if err := t.DB.Create(&topicCollect).Error; err != nil {
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
func (t TopicController) CancelCollect(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有收藏
	if t.DB.Where("user_id = ? and topic_id = ?", user.ID, id).First(&model.TopicCollect{}).Error != nil {
		response.Fail(ctx, nil, "未收藏")
		return
	} else {
		t.DB.Where("user_id = ? and topic_id = ?", user.ID, id).Delete(&model.TopicCollect{})
		response.Success(ctx, nil, "取消收藏成功")
		return
	}
}

// @title    CollectShow
// @description   查看收藏状态
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) CollectShow(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有收藏
	if t.DB.Where("user_id = ? and topic_id = ?", user.ID, id).First(&model.TopicCollect{}).Error != nil {
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
func (t TopicController) CollectList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var topicCollects []model.TopicCollect

	var total int64

	// TODO 查看收藏的数量
	t.DB.Where("topic_id = ?", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&topicCollects).Count(&total)

	response.Success(ctx, gin.H{"topicCollects": topicCollects, "total": total}, "查看成功")
}

// @title    CollectNumber
// @description   查看收藏用户数量
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) CollectNumber(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var total int64

	// TODO 查看收藏的数量
	t.DB.Where("topic_id = ?", id).Count(&total)

	response.Success(ctx, gin.H{"total": total}, "查看成功")
}

// @title    Collects
// @description   查看用户收藏夹
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) Collects(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var topicCollects []model.TopicCollect

	var total int64

	// TODO 查看收藏的数量
	t.DB.Where("user_id = ?", user.ID).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&topicCollects).Count(&total)

	response.Success(ctx, gin.H{"topicCollects": topicCollects, "total": total}, "查看成功")
}

// @title    Visit
// @description   游览主题
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) Visit(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	var topic model.Topic

	// TODO 查看主题是否存在
	// TODO 先看redis中是否存在
	if ok, _ := t.Redis.HExists(ctx, "Topic", id).Result(); ok {
		cate, _ := t.Redis.HGet(ctx, "Topic", id).Result()
		if json.Unmarshal([]byte(cate), &topic) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			t.Redis.HDel(ctx, "Topic", id)
		}
	}

	// TODO 查看主题是否在数据库中存在
	if t.DB.Where("id = ?", id).First(&topic).Error != nil {
		response.Fail(ctx, nil, "主题不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(topic)
		t.Redis.HSet(ctx, "Topic", id, v)
	}
leep:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	topicVisit := model.TopicVisit{
		TopicId: topic.ID,
		UserId:  user.ID,
	}

	// TODO 插入数据
	if err := t.DB.Create(&topicVisit).Error; err != nil {
		response.Fail(ctx, nil, "主题游览上传出错，数据库存储错误")
		return
	}

	response.Success(ctx, nil, "主题游览成功")
}

// @title    VisitNumber
// @description   游览主题数量
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) VisitNumber(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获得游览总数
	var total int64

	t.DB.Where("topic_id = ?", id).Count(&total)

	response.Success(ctx, gin.H{"total": total}, "请求主题游览数目成功")
}

// @title    VisitList
// @description   游览主题列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) VisitList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var topicVisits []model.TopicVisit

	var total int64

	// TODO 查看收藏的数量
	t.DB.Where("topic_id = ?", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&topicVisits).Count(&total)

	response.Success(ctx, gin.H{"topicVisits": topicVisits, "total": total}, "查看成功")
}

// @title    Visits
// @description   游览主题列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) Visits(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var topicVisits []model.TopicVisit

	var total int64

	// TODO 查看收藏的数量
	t.DB.Where("user_id = ?", user.ID).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&topicVisits).Count(&total)

	response.Success(ctx, gin.H{"topicVisits": topicVisits, "total": total}, "查看成功")
}

// @title    LabelCreate
// @description   标签创建
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) LabelCreate(ctx *gin.Context) {
	// TODO 获取指定主题
	id := ctx.Params.ByName("id")

	// TODO 获取标签
	label := ctx.Params.ByName("label")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看主题是否存在
	var topic model.Topic

	// TODO 先尝试在redis中寻找
	if ok, _ := t.Redis.HExists(ctx, "Topic", id).Result(); ok {
		art, _ := t.Redis.HGet(ctx, "Topic", id).Result()
		if json.Unmarshal([]byte(art), &topic) == nil {
			goto leep
		} else {
			// TODO 解码失败，删除字段
			t.Redis.HDel(ctx, "Topic", id)
		}
	}

	// TODO 查看主题是否在数据库中存在
	if t.DB.Where("id = ?", id).First(&topic).Error != nil {
		response.Fail(ctx, nil, "主题不存在")
		return
	}
	{
		// TODO 将主题存入redis供下次使用
		v, _ := json.Marshal(topic)
		t.Redis.HSet(ctx, "Topic", id, v)
	}
leep:

	// TODO 查看是否为主题作者
	if topic.UserId != user.ID {
		response.Fail(ctx, nil, "不是主题作者，请勿非法操作")
		return
	}

	// TODO 创建标签
	topicLabel := model.TopicLabel{
		Label:   label,
		TopicId: topic.ID,
	}

	// TODO 插入数据
	if err := t.DB.Create(&topicLabel).Error; err != nil {
		response.Fail(ctx, nil, "主题标签上传出错，数据验证有误")
		return
	}

	// TODO 解码失败，删除字段
	t.Redis.HDel(ctx, "TopicLabel", id)

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    LabelDelete
// @description   标签删除
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) LabelDelete(ctx *gin.Context) {
	// TODO 获取指定主题
	id := ctx.Params.ByName("id")

	// TODO 获取标签
	label := ctx.Params.ByName("label")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看主题是否存在
	var topic model.Topic

	// TODO 先尝试在redis中寻找
	if ok, _ := t.Redis.HExists(ctx, "Topic", id).Result(); ok {
		art, _ := t.Redis.HGet(ctx, "Topic", id).Result()
		if json.Unmarshal([]byte(art), &topic) == nil {
			goto leep
		} else {
			// TODO 解码失败，删除字段
			t.Redis.HDel(ctx, "Topic", id)
		}
	}

	// TODO 查看主题是否在数据库中存在
	if t.DB.Where("id = ?", id).First(&topic).Error != nil {
		response.Fail(ctx, nil, "主题不存在")
		return
	}
	{
		// TODO 将主题存入redis供下次使用
		v, _ := json.Marshal(topic)
		t.Redis.HSet(ctx, "Topic", id, v)
	}
leep:

	// TODO 查看是否为主题作者
	if topic.UserId != user.ID {
		response.Fail(ctx, nil, "不是主题作者，请勿非法操作")
		return
	}

	// TODO 删除主题标签
	if t.DB.Where("id = ?", label).First(&model.TopicLabel{}).Error != nil {
		response.Fail(ctx, nil, "标签不存在")
		return
	}

	t.DB.Where("id = ?", label).Delete(&model.TopicLabel{})

	// TODO 解码失败，删除字段
	t.Redis.HDel(ctx, "TopicLabel", id)

	// TODO 成功
	response.Success(ctx, nil, "删除成功")
}

// @title    LabelShow
// @description   标签查看
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) LabelShow(ctx *gin.Context) {
	// TODO 获取指定主题
	id := ctx.Params.ByName("id")

	// TODO 查找数据
	var topicLabels []model.TopicLabel
	// TODO 先尝试在redis中寻找
	if ok, _ := t.Redis.HExists(ctx, "TopicLabel", id).Result(); ok {
		art, _ := t.Redis.HGet(ctx, "TopicLabel", id).Result()
		if json.Unmarshal([]byte(art), &topicLabels) == nil {
			goto leap
		} else {
			// TODO 解码失败，删除字段
			t.Redis.HDel(ctx, "TopicLabel", id)
		}
	}

	// TODO 在数据库中查找
	t.DB.Where("topic_id = ?", id).Find(&topicLabels)
	{
		// TODO 将主题标签存入redis供下次使用
		v, _ := json.Marshal(topicLabels)
		t.Redis.HSet(ctx, "TopicLabel", id, v)
	}

leap:

	// TODO 成功
	response.Success(ctx, gin.H{"topicLabels": topicLabels}, "查看成功")
}

// @title    Search
// @description   文本搜索
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) Search(ctx *gin.Context) {
	// TODO 获取文本
	text := ctx.Params.ByName("text")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	var topics []model.Topic

	// TODO 模糊匹配
	t.DB.Where("match(title,content,res_long,res_short) against(? in boolean mode)", text+"*").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&topics)

	// TODO 查看查询总数
	var total int64
	t.DB.Where("match(title,content,res_long,res_short) against(? in boolean mode)", text+"*").Model(model.Topic{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"topics": topics, "total": total}, "成功")
}

// @title    SearchLabel
// @description   指定标签的搜索
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) SearchLabel(ctx *gin.Context) {

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
	var topicIds []struct {
		TopicId uuid.UUID `json:"topic_id"` // 题目外键
	}

	// TODO 进行标签匹配
	t.DB.Distinct("topic_id").Where("label in (?)", requestLabels.Labels).Model(model.TopicLabel{}).Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&topicIds)

	// TODO 查看查询总数
	var total int64
	t.DB.Distinct("topic_id").Where("label in (?)", requestLabels.Labels).Model(model.TopicLabel{}).Count(&total)

	// TODO 查找对应主题
	var topics []model.Topic

	t.DB.Where("id in (?)", topicIds).Find(&topics)

	// TODO 返回数据
	response.Success(ctx, gin.H{"topics": topics, "total": total}, "成功")
}

// @title    SearchWithLabel
// @description   指定标签与文本的搜索
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t TopicController) SearchWithLabel(ctx *gin.Context) {

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
	var topicIds []struct {
		TopicId uuid.UUID `json:"topic_id"` // 题目外键
	}

	// TODO 进行标签匹配
	t.DB.Distinct("topic_id").Where("label in (?)", requestLabels.Labels).Model(model.TopicLabel{}).Find(&topicIds)

	// TODO 查找对应主题
	var topics []model.Topic

	// TODO 模糊匹配
	t.DB.Where("id in (?) and match(title,content,res_long,res_short) against(? in boolean mode)", topicIds, text+"*").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&topics)

	// TODO 查看查询总数
	var total int64
	t.DB.Where("id in (?) and match(title,content,res_long,res_short) against(? in boolean mode)", topicIds, text+"*").Model(model.Topic{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"topics": topics, "total": total}, "成功")
}

// @title    NewTopicController
// @description   新建一个ITopicController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   ITopicController		返回一个ITopicController用于调用各种函数
func NewTopicController() ITopicController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.Topic{})
	db.AutoMigrate(model.TopicCollect{})
	db.AutoMigrate(model.TopicLike{})
	db.AutoMigrate(model.TopicVisit{})
	db.AutoMigrate(model.ProblemList{})
	db.AutoMigrate(model.TopicLabel{})
	return TopicController{DB: db, Redis: redis}
}
