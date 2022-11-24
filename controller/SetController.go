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
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ISetController			定义了表单类接口
type ISetController interface {
	Interface.RestInterface     // 包含增删查改功能
	Interface.LikeInterface     // 包含点赞功能
	Interface.CollectInterface  // 包含收藏功能
	Interface.VisitInterface    // 包含游览功能
	UserList(ctx *gin.Context)  // 查看指定用户的多篇表单
	TopicList(ctx *gin.Context) // 查看指定表单的主题列表
	GroupList(ctx *gin.Context) // 查看指定表单的用户组列表
}

// SetController			定义了表单工具类
type SetController struct {
	DB *gorm.DB // 含有一个数据库指针
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
		Title:    requestSet.Title,
		Content:  requestSet.Content,
		Reslong:  requestSet.Reslong,
		Resshort: requestSet.Resshort,
		UserId:   user.ID,
	}

	// TODO 插入数据
	if err := s.DB.Create(&set).Error; err != nil {
		response.Fail(ctx, nil, "表单上传出错，数据验证有误")
		return
	}

	// TODO 插入相关主题
	for _, v := range requestSet.Topics {
		if s.DB.Where("id = ?", v).First(&model.Topic{}).Error != nil {
			response.Fail(ctx, nil, "主题不存在")
			return
		}
		topicList := model.TopicList{
			SetId:   set.ID,
			TopicId: uint(v),
		}
		if s.DB.Create(&topicList).Error != nil {
			response.Fail(ctx, nil, "主题上传出错，数据验证有误")
			return
		}
	}

	// TODO 插入相关用户组
	for _, v := range requestSet.Groups {
		var group model.Group
		if s.DB.Where("id = ?", v).First(&group).Error != nil {
			response.Fail(ctx, nil, "用户组不存在")
			return
		}
		if group.LeaderId != user.ID {
			response.Fail(ctx, nil, "不是该组的组长，不能进行此操作")
			return
		}
		groupList := model.GroupList{
			SetId:   set.ID,
			GroupId: uint(v),
		}
		if s.DB.Create(&groupList).Error != nil {
			response.Fail(ctx, nil, "用户组上传出错，数据验证有误")
			return
		}
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
		Title:    requestSet.Title,
		Content:  requestSet.Content,
		Reslong:  requestSet.Reslong,
		Resshort: requestSet.Resshort,
		UserId:   user.ID,
	}

	// TODO 更新表单内容
	s.DB.Where("id = ?", id).Updates(set)

	if len(requestSet.Topics) != 0 {
		s.DB.Where("set_id = ?", id).Delete(&model.TopicList{})
		// TODO 插入相关主题
		for _, v := range requestSet.Topics {
			if s.DB.Where("id = ?", v).First(&model.Topic{}).Error != nil {
				response.Fail(ctx, nil, "主题不存在")
				return
			}
			topicList := model.TopicList{
				SetId:   set.ID,
				TopicId: uint(v),
			}
			if s.DB.Create(&topicList).Error != nil {
				response.Fail(ctx, nil, "主题上传出错，数据验证有误")
				return
			}
		}
	}

	if len(requestSet.Groups) != 0 {
		s.DB.Where("set_id = ?", id).Delete(&model.GroupList{})
		// TODO 插入相关用户组
		for _, v := range requestSet.Groups {
			var group model.Group
			if s.DB.Where("id = ?", v).First(&group).Error != nil {
				response.Fail(ctx, nil, "用户组不存在")
				return
			}
			if group.LeaderId != user.ID {
				response.Fail(ctx, nil, "不是该组的组长，不能进行此操作")
				return
			}
			groupList := model.GroupList{
				SetId:   set.ID,
				GroupId: uint(v),
			}
			if s.DB.Create(&groupList).Error != nil {
				response.Fail(ctx, nil, "用户组上传出错，数据验证有误")
				return
			}
		}
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

	// TODO 查看表单是否在数据库中存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}

	response.Success(ctx, gin.H{"set": set}, "成功")
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
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}

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

	var set model.Set

	// TODO 查看表单是否存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}

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

	var set model.Set

	// TODO 查看表单是否存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}

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

	var set model.Set

	// TODO 查看表单是否存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}

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

	var set model.Set

	// TODO 查看讨论是否存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "讨论不存在")
		return
	}

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
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}

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

	var set model.Set

	// TODO 查看表单是否存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有收藏
	if s.DB.Where("user_id = ? and set_id = ?", user.ID, set.ID).First(&model.SetCollect{}).Error != nil {
		response.Fail(ctx, nil, "未收藏")
		return
	} else {
		s.DB.Where("user_id = ? and set_id = ?", user.ID, set.ID).Delete(&model.SetCollect{})
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

	var set model.Set

	// TODO 查看表单是否存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 如果没有收藏
	if s.DB.Where("user_id = ? and set_id = ?", user.ID, set.ID).First(&model.SetCollect{}).Error != nil {
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

	var set model.Set

	// TODO 查看题解是否存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "题解不存在")
		return
	}

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var setCollects []model.SetCollect

	var total int64

	// TODO 查看收藏的数量
	s.DB.Where("set_id = ?", set.ID).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&setCollects).Count(&total)

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

	var set model.Set

	// TODO 查看表单是否存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}

	var total int64

	// TODO 查看收藏的数量
	s.DB.Where("set_id = ?", set.ID).Count(&total)

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
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}

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

	var set model.Set

	// TODO 查看表单是否存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}

	// TODO 获得游览总数
	var total int64

	s.DB.Where("set_id = ?", set.ID).Count(&total)

	response.Success(ctx, gin.H{"total": total}, "请求表单游览数目成功")
}

// @title    VisitList
// @description   游览表单列表
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s SetController) VisitList(ctx *gin.Context) {
	// TODO 获取path中的id
	id := ctx.Params.ByName("id")

	var set model.Set

	// TODO 查看表单是否存在
	if s.DB.Where("id = ?", id).First(&set).Error != nil {
		response.Fail(ctx, nil, "表单不存在")
		return
	}

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var setVisits []model.SetVisit

	var total int64

	// TODO 查看收藏的数量
	s.DB.Where("set_id = ?", set.ID).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&setVisits).Count(&total)

	response.Success(ctx, gin.H{"setVisits": setVisits, "total": total}, "查看成功")
}

// @title    Visits
// @description   游览表单列表
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

	// TODO 查看收藏的数量
	s.DB.Where("user_id = ?", user.ID).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&setVisits).Count(&total)

	response.Success(ctx, gin.H{"setVisits": setVisits, "total": total}, "查看成功")
}

// @title    NewSetController
// @description   新建一个ISetController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   ISetController		返回一个ISetController用于调用各种函数
func NewSetController() ISetController {
	db := common.GetDB()
	db.AutoMigrate(model.Set{})
	db.AutoMigrate(model.SetCollect{})
	db.AutoMigrate(model.SetLike{})
	db.AutoMigrate(model.SetVisit{})
	db.AutoMigrate(model.GroupList{})
	db.AutoMigrate(model.TopicList{})
	return SetController{DB: db}
}
