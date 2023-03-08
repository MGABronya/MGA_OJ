// @Title  SetController
// @Description  该文件提供关于操作表单的各种方法
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

// ISetController			定义了表单类接口
type ISetController interface {
	Interface.RestInterface      // 包含增删查改功能
	Interface.LikeInterface      // 包含点赞功能
	Interface.CollectInterface   // 包含收藏功能
	Interface.VisitInterface     // 包含游览功能
	Interface.ApplyInterface     // 包含申请接口
	Interface.LabelInterface     // 包含标签接口
	UserList(ctx *gin.Context)   // 查看指定用户的多篇表单
	TopicList(ctx *gin.Context)  // 查看指定表单的主题列表
	GroupList(ctx *gin.Context)  // 查看指定表单的用户组列表
	RankList(ctx *gin.Context)   // 查看表单内的用户排行
	RankUpdate(ctx *gin.Context) // 更新表单内用户排行
}

// SetController			定义了表单工具类
type SetController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Create
// @description   创建一篇表单
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) Create(ctx *gin.Context) {
	var requestSet vo.SetRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestSet); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 创建表单
	set := model.Set{
		Title:      requestSet.Title,
		Content:    requestSet.Content,
		Reslong:    requestSet.Reslong,
		Resshort:   requestSet.Resshort,
		UserId:     user.ID,
		AutoUpdate: requestSet.AutoUpdate,
		AutoPass:   requestSet.AutoPass,
		PassNum:    requestSet.PassNum,
		PassRe:     requestSet.PassRe,
	}

	// TODO 插入数据
	if err := s.DB.Create(&set).Error; err != nil {
		response.Fail(ctx, nil, "表单上传出错，数据验证有误")
		return
	}

	// TODO 插入相关用户组
	if user.Level >= 2 {
		for _, v := range requestSet.Groups {
			var group model.Group

			id := fmt.Sprint(v)

			// TODO 先看redis中是否存在
			if ok, _ := s.Redis.HExists(ctx, "Group", id).Result(); ok {
				cate, _ := s.Redis.HGet(ctx, "Group", id).Result()
				if json.Unmarshal([]byte(cate), &group) == nil {
					goto leep
				} else {
					// TODO 移除损坏数据
					s.Redis.HDel(ctx, "Group", id)
				}
			}

			// TODO 查看用户组是否在数据库中存在
			if s.DB.Where("id = ?", id).First(&group).Error != nil {
				response.Fail(ctx, nil, "用户组不存在")
				return
			}
			{
				// TODO 将用户组存入redis供下次使用
				v, _ := json.Marshal(group)
				s.Redis.HSet(ctx, "Group", id, v)
			}
		leep:
			if group.LeaderId != user.ID {
				response.Fail(ctx, nil, "不是该组的组长，不能进行此操作")
				return
			}
			// TODO 查看用户组是否可以合法加入表单
			if ok, err := CanAddGroup(set.ID, group.ID, set.PassNum, set.PassRe); !ok || err != nil {
				response.Fail(ctx, nil, "用户组无法合法加入表单")
				return
			}
			groupList := model.GroupList{
				SetId:   set.ID,
				GroupId: v,
			}
			if s.DB.Create(&groupList).Error != nil {
				response.Fail(ctx, nil, "用户组上传出错，数据验证有误")
				return
			}
		}
	}

	// TODO 插入相关主题
	for _, v := range requestSet.Topics {
		var topic model.Topic
		id := fmt.Sprint(v)

		// TODO 先看redis中是否存在
		if ok, _ := s.Redis.HExists(ctx, "Topic", id).Result(); ok {
			cate, _ := s.Redis.HGet(ctx, "Topic", id).Result()
			if json.Unmarshal([]byte(cate), &topic) == nil {
				goto leap
			} else {
				// TODO 移除损坏数据
				s.Redis.HDel(ctx, "Topic", id)
			}
		}

		// TODO 查看主题是否在数据库中存在
		if s.DB.Where("id = ?", id).First(&topic).Error != nil {
			response.Fail(ctx, nil, "主题不存在")
			return
		}
		{
			// TODO 将用户组存入redis供下次使用
			v, _ := json.Marshal(topic)
			s.Redis.HSet(ctx, "Topic", id, v)
		}
	leap:
		topicList := model.TopicList{
			SetId:   set.ID,
			TopicId: v,
		}
		if s.DB.Create(&topicList).Error != nil {
			response.Fail(ctx, nil, "主题上传出错，数据验证有误")
			return
		}
	}

	// TODO 尝试更新
	if updateRank(set.ID) != nil {
		response.Fail(ctx, nil, "更新出错")
		return
	}

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    Update
// @description   更新一篇表单的内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) Update(ctx *gin.Context) {
	var requestSet vo.SetRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestSet); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查找对应表单
	id := ctx.Params.ByName("id")

	var set model.Set

	if s.DB.Where("id = ?", id).First(&set) != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}

	// TODO 查看是否是用户作者
	if user.ID != set.UserId {
		response.Fail(ctx, nil, "不是表单作者，无法修改表单")
		return
	}

	// TODO 创建表单
	set = model.Set{
		Title:      requestSet.Title,
		Content:    requestSet.Content,
		Reslong:    requestSet.Reslong,
		Resshort:   requestSet.Resshort,
		AutoUpdate: requestSet.AutoUpdate,
		AutoPass:   requestSet.AutoPass,
		PassNum:    requestSet.PassNum,
		PassRe:     requestSet.PassRe,
		UserId:     user.ID,
	}

	// TODO 更新表单内容
	s.DB.Where("id = ?", id).Updates(set)

	// TODO 移除损坏数据
	s.Redis.HDel(ctx, "Set", id)

	if len(requestSet.Groups) != 0 {
		s.DB.Where("set_id = ?", id).Delete(&model.GroupList{})
		// TODO 插入相关用户组
		for _, v := range requestSet.Groups {

			var group model.Group
			id := fmt.Sprint(v)

			// TODO 先看redis中是否存在
			if ok, _ := s.Redis.HExists(ctx, "Group", id).Result(); ok {
				cate, _ := s.Redis.HGet(ctx, "Group", id).Result()
				if json.Unmarshal([]byte(cate), &group) == nil {
					goto leep
				} else {
					// TODO 移除损坏数据
					s.Redis.HDel(ctx, "Group", id)
				}
			}

			// TODO 查看用户组是否在数据库中存在
			if s.DB.Where("id = ?", id).First(&group).Error != nil {
				response.Fail(ctx, nil, "用户组不存在")
				return
			}
			{
				// TODO 将用户组存入redis供下次使用
				v, _ := json.Marshal(group)
				s.Redis.HSet(ctx, "Group", id, v)
			}
		leep:
			if group.LeaderId != user.ID {
				response.Fail(ctx, nil, "不是该组的组长，不能进行此操作")
				return
			}
			// TODO 查看用户组是否可以合法加入表单
			if ok, err := CanAddGroup(set.ID, group.ID, set.PassNum, set.PassRe); !ok || err != nil {
				response.Fail(ctx, nil, "用户组无法合法加入表单")
				return
			}

			groupList := model.GroupList{
				SetId:   set.ID,
				GroupId: v,
			}
			if s.DB.Create(&groupList).Error != nil {
				response.Fail(ctx, nil, "用户组上传出错，数据验证有误")
				return
			}
		}
	}

	if len(requestSet.Topics) != 0 {
		s.DB.Where("set_id = ?", id).Delete(&model.TopicList{})
		// TODO 插入相关主题
		var topic model.Topic
		for _, v := range requestSet.Topics {
			// TODO 先看redis中是否存在
			if ok, _ := s.Redis.HExists(ctx, "Topic", id).Result(); ok {
				cate, _ := s.Redis.HGet(ctx, "Topic", id).Result()
				if json.Unmarshal([]byte(cate), &topic) == nil {
					goto leap
				} else {
					// TODO 移除损坏数据
					s.Redis.HDel(ctx, "Topic", id)
				}
			}

			// TODO 查看主题是否在数据库中存在
			if s.DB.Where("id = ?", id).First(&topic).Error != nil {
				response.Fail(ctx, nil, "主题不存在")
				return
			}
			{
				// TODO 将用户组存入redis供下次使用
				v, _ := json.Marshal(topic)
				s.Redis.HSet(ctx, "Topic", id, v)
			}
		leap:
			topicList := model.TopicList{
				SetId:   set.ID,
				TopicId: v,
			}
			if s.DB.Create(&topicList).Error != nil {
				response.Fail(ctx, nil, "主题上传出错，数据验证有误")
				return
			}
		}
	}

	// TODO 尝试更新
	if updateRank(set.ID) != nil {
		response.Fail(ctx, nil, "更新出错")
		return
	}

	// TODO 成功
	response.Success(ctx, nil, "更新成功")
}

// @title    Show
// @description   查看一篇表单的内容
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")
	var set model.Set

	// TODO 先看redis中是否存在
	if ok, _ := s.Redis.HExists(ctx, "Set", id).Result(); ok {
		cate, _ := s.Redis.HGet(ctx, "Set", id).Result()
		if json.Unmarshal([]byte(cate), &set) == nil {
			response.Success(ctx, gin.H{"set": set}, "成功")
			return
		} else {
			// TODO 移除损坏数据
			s.Redis.HDel(ctx, "Set", id)
		}
	}

	// TODO 查看表单是否在数据库中存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}

	response.Success(ctx, gin.H{"set": set}, "成功")

	// TODO 将提交存入redis供下次使用
	v, _ := json.Marshal(set)
	s.Redis.HSet(ctx, "Set", id, v)

}

// @title    Delete
// @description   删除一篇表单
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) Delete(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	var set model.Set

	// TODO 查看表单是否存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}

	// TODO 判断当前用户是否为表单的作者
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有操作表单的权力
	if user.ID != set.UserId && user.Level < 4 {
		response.Fail(ctx, nil, "表单不属于您，请勿非法操作")
		return
	}

	// TODO 删除表单
	s.DB.Delete(&set)

	// TODO 移除损坏数据
	s.Redis.HDel(ctx, "Set", id)

	response.Success(ctx, nil, "删除成功")
}

// @title    PageList
// @description   获取多篇表单
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) PageList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var sets []model.Set

	// TODO 查找所有分页中可见的条目
	s.DB.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&sets)

	var total int64
	s.DB.Model(model.Set{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"sets": sets, "total": total}, "成功")
}

// @title    UserList
// @description   获取指定用户的多篇表单
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) UserList(ctx *gin.Context) {

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 分页
	var sets []model.Set

	// TODO 查找所有分页中可见的条目
	s.DB.Where("user_id = ?", id).Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&sets)

	var total int64
	s.DB.Model(model.Set{}).Where("user_id = ?", id).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"sets": sets, "total": total}, "成功")
}

// @title    TopicList
// @description   获取指定表单的主题列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) TopicList(ctx *gin.Context) {

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 分页
	var topicList []model.TopicList

	// TODO 查找所有分页中可见的条目
	s.DB.Where("set_id = ?", id).Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&topicList)

	var total int64
	s.DB.Model(model.Set{}).Where("user_id = ?", id).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"topicList": topicList, "total": total}, "成功")
}

// @title    GroupList
// @description   获取指定表单的用户组列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) GroupList(ctx *gin.Context) {

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 分页
	var groupList []model.GroupList

	// TODO 查找所有分页中可见的条目
	s.DB.Where("set_id = ?", id).Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&groupList)

	var total int64
	s.DB.Model(model.Set{}).Where("user_id = ?", id).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"groupList": groupList, "total": total}, "成功")
}

// @title    ProblemList
// @description   获取指定表单的多篇题目
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) ProblemList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 分页
	var problemLists []model.ProblemList

	// TODO 查找所有分页中可见的条目
	s.DB.Where("set_id = ?", id).Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&problemLists)

	var total int64
	s.DB.Model(model.Set{}).Where("set_id = ?", id).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"problemLists": problemLists, "total": total}, "成功")
}

// @title    Like
// @description   点赞或点踩
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) Like(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	var set model.Set

	// TODO 查看表单是否存在
	// TODO 先看redis中是否存在
	if ok, _ := s.Redis.HExists(ctx, "Set", id).Result(); ok {
		cate, _ := s.Redis.HGet(ctx, "Set", id).Result()
		if json.Unmarshal([]byte(cate), &set) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			s.Redis.HDel(ctx, "Set", id)
		}
	}
	// TODO 查看表单是否在数据库中存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(set)
		s.Redis.HSet(ctx, "Set", id, v)
	}
leep:
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有点赞或者点踩
	if s.DB.Where("user_id = ? and set_id = ?", user.ID, id).Update("like", like).Error != nil {
		// TODO 插入数据
		setLike := model.SetLike{
			SetId:  set.ID,
			UserId: user.ID,
			Like:   like,
		}
		if err := s.DB.Create(&setLike).Error; err != nil {
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
func (s SetController) CancelLike(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 取消点赞或者点踩
	s.DB.Where("user_id = ? and set_id = ?", user.ID, id).Delete(&model.SetLike{})

	response.Success(ctx, nil, "取消成功")
}

// @title    LikeNumber
// @description   点赞或点踩的数量
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) LikeNumber(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	var total int64

	// TODO 查看点赞或者点踩的数量
	s.DB.Where("set_id = ? and like = ?", id, like).Count(&total)

	response.Success(ctx, gin.H{"total": total}, "查看成功")
}

// @title    LikeList
// @description   点赞或点踩的列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) LikeList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var setLikes []model.SetLike

	var total int64

	// TODO 查看点赞或者点踩的数量
	s.DB.Where("set_id = ? and like = ?", id, like).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&setLikes).Count(&total)

	response.Success(ctx, gin.H{"setLikes": setLikes, "total": total}, "查看成功")
}

// @title    LikeShow
// @description   查看用户点赞状态
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) LikeShow(ctx *gin.Context) {
	// 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var setLike model.SetLike

	// TODO 查看点赞状态
	if s.DB.Where("user_id = ? and set_id = ?", user.ID, id).First(&setLike).Error != nil {
		response.Success(ctx, gin.H{"like": 0}, "暂无状态")
		return
	}

	if setLike.Like {
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
func (s SetController) Likes(ctx *gin.Context) {

	// TODO 获取like
	like, _ := strconv.ParseBool(ctx.Query("like"))

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 分页
	var setLikes []model.SetLike

	var total int64

	// TODO 查看点赞或者点踩的数量
	s.DB.Where("user_id = ? and like = ?", user.ID, like).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&setLikes).Count(&total)

	response.Success(ctx, gin.H{"setLikes": setLikes, "total": total}, "查看成功")
}

// @title    Collect
// @description   收藏
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) Collect(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	var set model.Set

	// TODO 查看表单是否存在
	// TODO 先看redis中是否存在
	if ok, _ := s.Redis.HExists(ctx, "Set", id).Result(); ok {
		cate, _ := s.Redis.HGet(ctx, "Set", id).Result()
		if json.Unmarshal([]byte(cate), &set) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			s.Redis.HDel(ctx, "Set", id)
		}
	}
	// TODO 查看表单是否在数据库中存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(set)
		s.Redis.HSet(ctx, "Set", id, v)
	}
leep:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有收藏
	if s.DB.Where("user_id = ? and set_id = ?", user.ID, set.ID).First(&model.SetCollect{}).Error != nil {
		setCollect := model.SetCollect{
			SetId:  set.ID,
			UserId: user.ID,
		}
		// TODO 插入数据
		if err := s.DB.Create(&setCollect).Error; err != nil {
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
func (s SetController) CancelCollect(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有收藏
	if s.DB.Where("user_id = ? and set_id = ?", user.ID, id).First(&model.SetCollect{}).Error != nil {
		response.Fail(ctx, nil, "未收藏")
		return
	} else {
		s.DB.Where("user_id = ? and set_id = ?", user.ID, id).Delete(&model.SetCollect{})
		response.Success(ctx, nil, "取消收藏成功")
		return
	}
}

// @title    CollectShow
// @description   查看收藏状态
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) CollectShow(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有收藏
	if s.DB.Where("user_id = ? and set_id = ?", user.ID, id).First(&model.SetCollect{}).Error != nil {
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
func (s SetController) CollectList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var setCollects []model.SetCollect

	var total int64

	// TODO 查看收藏的数量
	s.DB.Where("set_id = ?", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&setCollects).Count(&total)

	response.Success(ctx, gin.H{"setCollects": setCollects, "total": total}, "查看成功")
}

// @title    CollectNumber
// @description   查看收藏用户数量
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) CollectNumber(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	var total int64

	// TODO 查看收藏的数量
	s.DB.Where("set_id = ?", id).Count(&total)

	response.Success(ctx, gin.H{"total": total}, "查看成功")
}

// @title    Collects
// @description   查看用户收藏夹
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) Collects(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var setCollects []model.SetCollect

	var total int64

	// TODO 查看收藏的数量
	s.DB.Where("user_id = ?", user.ID).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&setCollects).Count(&total)

	response.Success(ctx, gin.H{"setCollects": setCollects, "total": total}, "查看成功")
}

// @title    Visit
// @description   游览表单
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) Visit(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	var set model.Set

	// TODO 查看表单是否存在
	// TODO 先看redis中是否存在
	if ok, _ := s.Redis.HExists(ctx, "Set", id).Result(); ok {
		cate, _ := s.Redis.HGet(ctx, "Set", id).Result()
		if json.Unmarshal([]byte(cate), &set) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			s.Redis.HDel(ctx, "Set", id)
		}
	}
	// TODO 查看表单是否在数据库中存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(set)
		s.Redis.HSet(ctx, "Set", id, v)
	}
leep:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	setVisit := model.SetVisit{
		SetId:  set.ID,
		UserId: user.ID,
	}

	// TODO 插入数据
	if err := s.DB.Create(&setVisit).Error; err != nil {
		response.Fail(ctx, nil, "表单游览上传出错，数据库存储错误")
		return
	}

	response.Success(ctx, nil, "表单游览成功")
}

// @title    VisitNumber
// @description   游览表单数量
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) VisitNumber(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获得游览总数
	var total int64

	s.DB.Where("set_id = ?", id).Count(&total)

	response.Success(ctx, gin.H{"total": total}, "请求表单游览数目成功")
}

// @title    VisitList
// @description   表单的游览表单列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) VisitList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var setVisits []model.SetVisit

	var total int64

	// TODO 查看游览的列表
	s.DB.Where("set_id = ?", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&setVisits).Count(&total)

	response.Success(ctx, gin.H{"setVisits": setVisits, "total": total}, "查看成功")
}

// @title    Visits
// @description   用户的游览表单列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) Visits(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var setVisits []model.SetVisit

	var total int64

	// TODO 查看游览的数量
	s.DB.Where("user_id = ?", user.ID).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&setVisits).Count(&total)

	response.Success(ctx, gin.H{"setVisits": setVisits, "total": total}, "查看成功")
}

// @title    RankList
// @description   游览用户排行列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) RankList(ctx *gin.Context) {

	// TODO 获取id
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var setRanks []model.SetRank

	var total int64

	// TODO 查看排行
	s.DB.Where("set_id = ?", id).Order("pass desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&setRanks).Count(&total)

	response.Success(ctx, gin.H{"setRanks": setRanks, "total": total}, "查看成功")
}

// @title    RankUpdate
// @description   更新用户排行
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) RankUpdate(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取id
	id := ctx.Params.ByName("id")

	var set model.Set

	// TODO 尝试找出set
	// TODO 查看表单是否存在
	// TODO 先看redis中是否存在
	if ok, _ := s.Redis.HExists(ctx, "Set", id).Result(); ok {
		cate, _ := s.Redis.HGet(ctx, "Set", id).Result()
		if json.Unmarshal([]byte(cate), &set) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			s.Redis.HDel(ctx, "Set", id)
		}
	}
	// TODO 查看表单是否在数据库中存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(set)
		s.Redis.HSet(ctx, "Set", id, v)
	}
leep:
	// TODO 查看是否为创建者
	if set.UserId != user.ID {
		response.Fail(ctx, nil, "不是表单创建者")
		return
	}

	// TODO 尝试更新
	if updateRank(set.ID) != nil {
		response.Fail(ctx, nil, "更新出错")
		return
	}

	response.Success(ctx, nil, "更新成功")
}

// @title    Apply
// @description   用户组申请加入某个表单
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) Apply(ctx *gin.Context) {
	var requestSetApply vo.SetApplyRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestSetApply); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取指定表单
	id := ctx.Params.ByName("id")

	var set model.Set

	// TODO 查看表单是否存在
	// TODO 先看redis中是否存在
	if ok, _ := s.Redis.HExists(ctx, "Set", id).Result(); ok {
		cate, _ := s.Redis.HGet(ctx, "Set", id).Result()
		if json.Unmarshal([]byte(cate), &set) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			s.Redis.HDel(ctx, "Set", id)
		}
	}
	// TODO 查看表单是否在数据库中存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(set)
		s.Redis.HSet(ctx, "Set", id, v)
	}
leep:

	var group model.Group
	// TODO 先看redis中是否存在
	id = fmt.Sprint(requestSetApply.GroupId)
	if ok, _ := s.Redis.HExists(ctx, "Group", id).Result(); ok {
		cate, _ := s.Redis.HGet(ctx, "Group", id).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto leap
		} else {
			// TODO 移除损坏数据
			s.Redis.HDel(ctx, "Group", id)
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if s.DB.Where("id = ?", id).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		s.Redis.HSet(ctx, "Group", id, v)
	}
leap:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看用户是否是用户组组长
	if user.ID != group.LeaderId {
		response.Fail(ctx, nil, "不是用户组组长")
		return
	}

	// TODO 查看用户组是否已经加入表单
	if s.DB.Where("set_id = ? and group_id = ?", set.ID, group.ID).First(&model.GroupList{}).Error == nil {
		response.Fail(ctx, nil, "已加入表单")
		return
	}

	// TODO 查看用户组是否被拉黑
	if s.DB.Where("set_id = ? and group_id = ?", set.ID, group.ID).First(&model.SetBlock{}).Error != nil {
		response.Fail(ctx, nil, "用户组已被拉黑")
		return
	}

	// TODO 查看用户组是否自动通过申请
	if set.AutoPass {
		if ok, err := CanAddGroup(set.ID, group.ID, set.PassNum, set.PassRe); !ok || err != nil {
			response.Fail(ctx, nil, "用户组不符合要求")
			return
		}
		groupList := model.GroupList{
			GroupId: group.ID,
			SetId:   set.ID,
		}
		// TODO 插入数据
		if err := s.DB.Create(&groupList).Error; err != nil {
			response.Fail(ctx, nil, "通过申请出错，数据验证有误")
			return
		}
		// TODO 成功
		response.Success(ctx, nil, "创建成功")
		return
	}

	var setApply model.SetApply

	// TODO 查看用户组是否已经发送过申请
	if s.DB.Where("set_id = ? and group_id = ?", set.ID, group.ID).First(&setApply).Error == nil && setApply.Condition {
		response.Fail(ctx, nil, "已发送过申请")
		return
	}

	// TODO 创建表单申请
	setApply = model.SetApply{
		SetId:     set.ID,
		GroupId:   group.ID,
		Condition: true,
		Content:   requestSetApply.Content,
		Reslong:   requestSetApply.Reslong,
		Resshort:  requestSetApply.Resshort,
	}

	// TODO 插入数据
	if err := s.DB.Create(&setApply).Error; err != nil {
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
func (s SetController) Consent(ctx *gin.Context) {

	// TODO 获取指定申请
	id := ctx.Params.ByName("id")

	var setApply model.SetApply

	// TODO 查看申请是否存在
	if s.DB.Where("id = ?", id).First(&setApply).Error != nil {
		response.Fail(ctx, nil, "申请不存在")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看该申请状态
	if !setApply.Condition {
		response.Fail(ctx, nil, "该申请已拒绝")
		return
	}

	// TODO 查看用户组是否已经加入表单
	if s.DB.Where("set_id = ? and group_id = ?", setApply.SetId, setApply.GroupId).First(&model.GroupList{}).Error == nil {
		response.Fail(ctx, nil, "已加入表单")
		return
	}

	// TODO 查看当前用户是否为表单创建者
	if s.DB.Where("id = ? and user_id = ?", setApply.SetId, user.ID).First(&model.Set{}).Error != nil {
		response.Fail(ctx, nil, "非表单创建者，无法操作")
		return
	}

	groupList := model.GroupList{
		GroupId: setApply.GroupId,
		SetId:   setApply.SetId,
	}

	// TODO 插入数据
	if err := s.DB.Create(&groupList).Error; err != nil {
		response.Fail(ctx, nil, "通过申请出错，数据验证有误")
		return
	}

	// TODO 删除申请
	s.DB.Delete(&groupList)

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    Refuse
// @description   拒绝申请
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) Refuse(ctx *gin.Context) {

	// TODO 获取指定申请
	id := ctx.Params.ByName("id")

	var setApply model.SetApply

	// TODO 查看申请是否存在
	if s.DB.Where("id = ?", id).First(&setApply).Error != nil {
		response.Fail(ctx, nil, "申请不存在")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看当前用户是否为表单创建者
	if s.DB.Where("id = ? and set_id = ?", setApply.SetId, user.ID).First(&model.Set{}).Error != nil {
		response.Fail(ctx, nil, "非表单创建者，无法操作")
		return
	}

	setApply.Condition = false

	// TODO 保存更改
	s.DB.Save(&setApply)

	// TODO 成功
	response.Success(ctx, nil, "拒绝成功")
}

// @title    Block
// @description   拉黑某用户组
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) Block(ctx *gin.Context) {

	// TODO 获取指定用户组
	group_id := ctx.Params.ByName("group")

	var group model.Group

	// TODO 查看用户组是否存在
	id := fmt.Sprint(group_id)
	if ok, _ := s.Redis.HExists(ctx, "Group", id).Result(); ok {
		cate, _ := s.Redis.HGet(ctx, "Group", id).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto leap
		} else {
			// TODO 移除损坏数据
			s.Redis.HDel(ctx, "Group", id)
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if s.DB.Where("id = ?", id).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		s.Redis.HSet(ctx, "Group", id, v)
	}
leap:

	// TODO 获取指定表单
	set_id := ctx.Params.ByName("set")

	var set model.Set

	// TODO 查看表单是否存在
	// TODO 先看redis中是否存在
	id = fmt.Sprint(set_id)
	if ok, _ := s.Redis.HExists(ctx, "Set", id).Result(); ok {
		cate, _ := s.Redis.HGet(ctx, "Set", id).Result()
		if json.Unmarshal([]byte(cate), &set) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			s.Redis.HDel(ctx, "Set", id)
		}
	}
	// TODO 查看表单是否在数据库中存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(set)
		s.Redis.HSet(ctx, "Set", id, v)
	}
leep:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看当前用户是否为表单创建者
	if set.UserId != user.ID {
		response.Fail(ctx, nil, "非表单创建者，无法操作")
		return
	}

	// TODO 查看当前用户组是否已经拉黑
	if s.DB.Where("set_id = ? and group_id = ?", set_id, group_id).First(&model.SetBlock{}).Error == nil {
		response.Fail(ctx, nil, "用户已拉黑")
		return
	}

	// TODO 将指定用户放入黑名单
	setBlock := model.SetBlock{
		SetId:   set.ID,
		GroupId: group.ID,
	}

	if s.DB.Create(&setBlock).Error != nil {
		response.Fail(ctx, nil, "黑名单入库错误")
		return
	}

	// TODO 成功
	response.Success(ctx, nil, "拉黑成功")
}

// @title    RemoveBlack
// @description   移除某用户组的黑名单
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) RemoveBlack(ctx *gin.Context) {

	// TODO 获取指定表单
	set_id := ctx.Params.ByName("set")

	var set model.Set

	// TODO 查看表单是否存在
	id := fmt.Sprint(set_id)
	if ok, _ := s.Redis.HExists(ctx, "Set", id).Result(); ok {
		cate, _ := s.Redis.HGet(ctx, "Set", id).Result()
		if json.Unmarshal([]byte(cate), &set) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			s.Redis.HDel(ctx, "Set", id)
		}
	}
	// TODO 查看表单是否在数据库中存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(set)
		s.Redis.HSet(ctx, "Set", id, v)
	}
leep:

	// TODO 获取指定用户组
	group_id := ctx.Params.ByName("group")

	var group model.Group

	// TODO 查看用户组是否存在
	id = fmt.Sprint(group_id)
	if ok, _ := s.Redis.HExists(ctx, "Group", id).Result(); ok {
		cate, _ := s.Redis.HGet(ctx, "Group", id).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto leap
		} else {
			// TODO 移除损坏数据
			s.Redis.HDel(ctx, "Group", id)
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if s.DB.Where("id = ?", id).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		s.Redis.HSet(ctx, "Group", id, v)
	}
leap:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看当前用户是否为表单创建者
	if set.UserId != user.ID {
		response.Fail(ctx, nil, "非表单创建者，无法操作")
		return
	}

	var setBlock model.SetBlock

	// TODO 查看当前用户组是否已经拉黑
	if s.DB.Where("set_id = ? and group_id = ?", set_id, group_id).First(&setBlock).Error != nil {
		response.Fail(ctx, nil, "用户组未被拉黑")
	}

	// TODO 将用户组移除黑名单
	s.DB.Delete(&setBlock)

	// TODO 成功
	response.Success(ctx, nil, "移除黑名单成功")
}

// @title    BlackList
// @description   查看黑名单
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) BlackList(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取指定表单
	id := ctx.Params.ByName("id")

	var set model.Set

	// TODO 查看表单是否存在
	if ok, _ := s.Redis.HExists(ctx, "Set", id).Result(); ok {
		cate, _ := s.Redis.HGet(ctx, "Set", id).Result()
		if json.Unmarshal([]byte(cate), &set) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			s.Redis.HDel(ctx, "Set", id)
		}
	}
	// TODO 查看表单是否在数据库中存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(set)
		s.Redis.HSet(ctx, "Set", id, v)
	}
leep:

	// TODO 查看当前用户是否为表单创建者
	if set.UserId != user.ID {
		response.Fail(ctx, nil, "非表单创建者，无法操作")
		return
	}

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var setBlocks []model.SetBlock

	var total int64

	// TODO 查看黑名单
	s.DB.Where("set_id = ?", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&setBlocks).Count(&total)

	response.Success(ctx, gin.H{"setBlocks": setBlocks, "total": total}, "查看成功")
}

// @title    ApplyingList
// @description   查看用户组申请列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) ApplyingList(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取指定用户组
	id := ctx.Params.ByName("id")

	var group model.Group

	// TODO 查看用户组是否存在
	if ok, _ := s.Redis.HExists(ctx, "Group", id).Result(); ok {
		cate, _ := s.Redis.HGet(ctx, "Group", id).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto leap
		} else {
			// TODO 移除损坏数据
			s.Redis.HDel(ctx, "Group", id)
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if s.DB.Where("id = ?", id).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		s.Redis.HSet(ctx, "Group", id, v)
	}
leap:

	// TODO 查看是否是用户组的组长
	if group.LeaderId != user.ID {
		response.Fail(ctx, nil, "非用户组组长")
		return
	}

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var setApplys []model.SetApply

	var total int64

	// TODO 查看申请的数量
	s.DB.Where("group_id = ?", group.ID).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&setApplys).Count(&total)

	response.Success(ctx, gin.H{"groupApplys": setApplys, "total": total}, "查看成功")
}

// @title    AppliedList
// @description   查看用户组申请列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) AppliedList(ctx *gin.Context) {

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取指定用户组
	id := ctx.Params.ByName("id")

	var set model.Set

	// TODO 查看表单是否存在
	if ok, _ := s.Redis.HExists(ctx, "Set", id).Result(); ok {
		cate, _ := s.Redis.HGet(ctx, "Set", id).Result()
		if json.Unmarshal([]byte(cate), &set) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			s.Redis.HDel(ctx, "Set", id)
		}
	}
	// TODO 查看表单是否在数据库中存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(set)
		s.Redis.HSet(ctx, "Set", id, v)
	}
leep:

	// TODO 查看是否是表单的创建者
	if set.UserId != user.ID {
		response.Fail(ctx, nil, "非表单创建者")
		return
	}

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var setApplys []model.SetApply

	var total int64

	// TODO 查看申请的数量
	s.DB.Where("set_id = ?", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&setApplys).Count(&total)

	response.Success(ctx, gin.H{"setApplys": setApplys, "total": total}, "查看成功")
}

// @title    Quit
// @description   退出表单
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) Quit(ctx *gin.Context) {

	// TODO 获取指定表单
	set_id := ctx.Params.ByName("set")

	var set model.Set

	// TODO 查看表单是否存在
	id := set_id
	if ok, _ := s.Redis.HExists(ctx, "Set", id).Result(); ok {
		cate, _ := s.Redis.HGet(ctx, "Set", id).Result()
		if json.Unmarshal([]byte(cate), &set) == nil {
			goto leep
		} else {
			// TODO 移除损坏数据
			s.Redis.HDel(ctx, "Set", id)
		}
	}
	// TODO 查看表单是否在数据库中存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}
	{
		// TODO 将提交存入redis供下次使用
		v, _ := json.Marshal(set)
		s.Redis.HSet(ctx, "Set", id, v)
	}
leep:

	// TODO 获取指定表单
	group_id := ctx.Params.ByName("group")

	var group model.Group

	// TODO 查看用户组是否存在
	id = group_id
	if ok, _ := s.Redis.HExists(ctx, "Group", id).Result(); ok {
		cate, _ := s.Redis.HGet(ctx, "Group", id).Result()
		if json.Unmarshal([]byte(cate), &group) == nil {
			goto leap
		} else {
			// TODO 移除损坏数据
			s.Redis.HDel(ctx, "Group", id)
		}
	}

	// TODO 查看用户组是否在数据库中存在
	if s.DB.Where("id = ?", id).First(&group).Error != nil {
		response.Fail(ctx, nil, "用户组不存在")
		return
	}
	{
		// TODO 将用户组存入redis供下次使用
		v, _ := json.Marshal(group)
		s.Redis.HSet(ctx, "Group", id, v)
	}
leap:

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否为用户组组长
	if user.ID != group.LeaderId {
		response.Fail(ctx, nil, "非用户组组长，无法操作")
		return
	}

	// TODO 查看用户组是否已经加入表单
	if s.DB.Where("set_id = ? and group_id = ?", set.ID, group.ID).First(&model.GroupList{}).Error != nil {
		response.Fail(ctx, nil, "未加入表单")
		return
	}

	// TODO 在表单中删除用户组
	s.DB.Where("set_id = ? and group_id = ?", set.ID, group.ID).Delete(&model.GroupList{})

	// TODO 成功
	response.Success(ctx, nil, "退出成功")
}

// @title    LabelCreate
// @description   标签创建
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) LabelCreate(ctx *gin.Context) {
	// TODO 获取指定表单
	id := ctx.Params.ByName("id")

	// TODO 获取标签
	label := ctx.Params.ByName("label")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看表单是否存在
	var set model.Set

	// TODO 先尝试在redis中寻找
	if ok, _ := s.Redis.HExists(ctx, "Set", id).Result(); ok {
		art, _ := s.Redis.HGet(ctx, "Set", id).Result()
		if json.Unmarshal([]byte(art), &set) == nil {
			goto leep
		} else {
			// TODO 解码失败，删除字段
			s.Redis.HDel(ctx, "Set", id)
		}
	}

	// TODO 查看表单是否在数据库中存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}
	{
		// TODO 将表单存入redis供下次使用
		v, _ := json.Marshal(set)
		s.Redis.HSet(ctx, "Set", id, v)
	}
leep:

	// TODO 查看是否为表单作者
	if set.UserId != user.ID {
		response.Fail(ctx, nil, "不是表单作者，请勿非法操作")
		return
	}

	// TODO 创建标签
	setLabel := model.SetLabel{
		Label: label,
		SetId: set.ID,
	}

	// TODO 插入数据
	if err := s.DB.Create(&setLabel).Error; err != nil {
		response.Fail(ctx, nil, "表单标签上传出错，数据验证有误")
		return
	}

	// TODO 解码失败，删除字段
	s.Redis.HDel(ctx, "SetLabel", id)

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    LabelDelete
// @description   标签删除
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) LabelDelete(ctx *gin.Context) {
	// TODO 获取指定表单
	id := ctx.Params.ByName("id")

	// TODO 获取标签
	label := ctx.Params.ByName("label")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看表单是否存在
	var set model.Set

	// TODO 先尝试在redis中寻找
	if ok, _ := s.Redis.HExists(ctx, "Set", id).Result(); ok {
		art, _ := s.Redis.HGet(ctx, "Set", id).Result()
		if json.Unmarshal([]byte(art), &set) == nil {
			goto leep
		} else {
			// TODO 解码失败，删除字段
			s.Redis.HDel(ctx, "Set", id)
		}
	}

	// TODO 查看表单是否在数据库中存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "题目不存在")
		return
	}
	{
		// TODO 将表单存入redis供下次使用
		v, _ := json.Marshal(set)
		s.Redis.HSet(ctx, "Set", id, v)
	}
leep:

	// TODO 查看是否为表单作者
	if set.UserId != user.ID {
		response.Fail(ctx, nil, "不是题目作者，请勿非法操作")
		return
	}

	// TODO 删除表单标签
	if s.DB.Where("id = ?", label).First(&model.SetLabel{}).Error != nil {
		response.Fail(ctx, nil, "标签不存在")
		return
	}

	s.DB.Where("id = ?", label).Delete(&model.SetLabel{})

	// TODO 解码失败，删除字段
	s.Redis.HDel(ctx, "SetLabel", id)

	// TODO 成功
	response.Success(ctx, nil, "删除成功")
}

// @title    LabelShow
// @description   标签查看
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) LabelShow(ctx *gin.Context) {
	// TODO 获取指定表单
	id := ctx.Params.ByName("id")

	// TODO 查找数据
	var setLabels []model.SetLabel
	// TODO 先尝试在redis中寻找
	if ok, _ := s.Redis.HExists(ctx, "SetLabel", id).Result(); ok {
		art, _ := s.Redis.HGet(ctx, "SetLabel", id).Result()
		if json.Unmarshal([]byte(art), &setLabels) == nil {
			goto leap
		} else {
			// TODO 解码失败，删除字段
			s.Redis.HDel(ctx, "SetLabel", id)
		}
	}

	// TODO 在数据库中查找
	s.DB.Where("set_id = ?", id).Find(&setLabels)
	{
		// TODO 将题目标签存入redis供下次使用
		v, _ := json.Marshal(setLabels)
		s.Redis.HSet(ctx, "SetLabel", id, v)
	}

leap:

	// TODO 成功
	response.Success(ctx, gin.H{"setLabels": setLabels}, "查看成功")
}

// @title    NewSetController
// @description   新建一个ISetController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   ISetController		返回一个ISetController用于调用各种函数
func NewSetController() ISetController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.Set{})
	db.AutoMigrate(model.SetCollect{})
	db.AutoMigrate(model.SetLike{})
	db.AutoMigrate(model.SetVisit{})
	db.AutoMigrate(model.GroupList{})
	db.AutoMigrate(model.TopicList{})
	db.AutoMigrate(model.SetRank{})
	db.AutoMigrate(model.SetApply{})
	db.AutoMigrate(model.SetBlock{})
	db.AutoMigrate(model.SetLabel{})
	return SetController{DB: db, Redis: redis}
}

// @title    updateRank
// @description   更新一个表单中的用户排行
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    set_id uuid.UUID		表示表单的id
// @return   error		返回一个error表示是否出现错误
func updateRank(set_id uuid.UUID) error {
	db := common.GetDB()
	var err error
	// TODO 删掉原先的排行
	db.Where("set_id = ?", set_id).Delete(&model.SetRank{})
	// TODO 重新建立排行
	userPass := make(map[uuid.UUID]uint, 0)
	// TODO 插入所有成员
	var groupLists []model.GroupList
	db.Where("set_id = ?", set_id).Find(&groupLists)
	for _, group := range groupLists {
		var userLists []model.UserList
		db.Where("group_id = ?", group.GroupId).Find(&userLists)
		for _, user := range userLists {
			userPass[user.UserId] = 0
		}
	}
	// 搜索所有成员的通过表单内题目数量
	var topicLists []model.TopicList
	db.Where("set_id = ?", set_id).Find(&topicLists)
	for _, topic := range topicLists {
		var problemLists []model.ProblemList
		db.Where("topic_id = ?", topic.TopicId).Find(&problemLists)
		for _, problem := range problemLists {
			for user := range userPass {
				if db.Where("user_id = ? and problem_id = ? and condition = Accepted", user, problem.ProblemId).First(&model.Record{}).Error == nil {
					userPass[user]++
				}
			}
		}
	}
	// TODO 将所有记录值放入数据库
	for user := range userPass {
		setRank := model.SetRank{
			UserId: user,
			Pass:   userPass[user],
			SetId:  set_id,
		}
		if err = db.Create(&setRank).Error; err != nil {
			return err
		}
	}
	return nil
}

// @title    CanAddGroup
// @description   更新一个表单中的用户排行
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    set_id uuid.UUID		表示表单的id
// @return   error		返回一个error表示是否出现错误
func CanAddGroup(set_id uuid.UUID, group_id uuid.UUID, PassNum uint, PassRe bool) (bool, error) {
	db := common.GetDB()

	var userLists []model.UserList
	db.Where("group_id = ?", group_id).Find(&userLists)

	// TODO 如果组员数量大于限制
	if int(PassNum) < len(userLists) {
		return false, nil
	}

	// TODO 如果没有重复限制
	if !PassRe {
		return true, nil
	}

	var groupLists []model.GroupList
	db.Where("set_id = ?", set_id).Find(&groupLists)

	userMap := make(map[uuid.UUID]bool, 0)

	// TODO 将表单内的所有用户填入map
	for _, group := range groupLists {
		var userLists []model.UserList
		db.Where("group_id = ?", group.GroupId).Find(&userLists)
		for _, user := range userLists {
			userMap[user.UserId] = true
		}
	}

	// TODO 查看成员是否重复
	for _, user := range userLists {
		if userMap[user.UserId] {
			return false, nil
		}
	}

	return true, nil
}
