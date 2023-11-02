// @Title  UserController
// @Description  该文件用于提供操作用户的各种方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package controller

import (
	"MGA_OJ/Interface"
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"
	"MGA_OJ/util"
	"MGA_OJ/vo"
	"encoding/json"
	"fmt"
	"log"

	"math/rand"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// IUserController			定义了用户类接口
type IUserController interface {
	Interface.LabelInterface         // 包含标签功能
	Interface.SearchInterface        // 包含搜索功能
	Register(ctx *gin.Context)       // 注册
	Login(ctx *gin.Context)          // 登录
	VerifyEmail(ctx *gin.Context)    // 验证码
	Security(ctx *gin.Context)       // 找回密码
	UpdatePass(ctx *gin.Context)     // 更新密码
	Info(ctx *gin.Context)           // 返回当前登录的用户
	Update(ctx *gin.Context)         // 用户的信息更新
	UpdateLevel(ctx *gin.Context)    // 修改用户的等级
	Show(ctx *gin.Context)           // 显示用户的所有信息
	AcceptNum(ctx *gin.Context)      // 显示用户ac题目数量
	AcceptRankList(ctx *gin.Context) // 显示用户ac题目的排行列表
	AcceptRank(ctx *gin.Context)     // 显示用户ac题目的排行
	ScoreRankList(ctx *gin.Context)  // 显示用户竞赛分排行列表
	ScoreRank(ctx *gin.Context)      // 显示用户竞赛分排行
	HotRankList(ctx *gin.Context)    // 显示用户热度排行列表
	HotRank(ctx *gin.Context)        // 显示用户热度排行
	ScoreChange(ctx *gin.Context)    // 显示用户的分数变化
	Hot(ctx *gin.Context)            // 显示用户今日热度数据
	LikeRank(ctx *gin.Context)       // 点赞榜单
	UnLikeRank(ctx *gin.Context)     // 点踩榜单
	CollectRank(ctx *gin.Context)    // 收藏榜单
	VisitRank(ctx *gin.Context)      // 游览榜单
}

// UserController			定义了用户工具类
type UserController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Register
// @description   用户注册
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) Register(ctx *gin.Context) {
	var requestUser = vo.UserRequest{}
	ctx.Bind(&requestUser)
	// TODO 获取参数
	email := requestUser.Email
	password := requestUser.Password
	name := requestUser.Name
	// TODO 数据验证
	if !util.VerifyEmailFormat(email) {
		response.Response(ctx, 201, 201, nil, "邮箱格式错误")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, 201, 201, nil, "密码不能少于6位")
		return
	}

	if len(name) == 0 {
		name = util.RandomString(10)
	}

	log.Println(name, email, password)

	// TODO 判断email是否存在
	if util.IsEmailExist(u.DB, email) {
		response.Response(ctx, 201, 201, nil, "用户已经存在")
		return
	}

	// TODO 判断email是否通过验证
	if !util.IsEmailPass(ctx, email, requestUser.Verify) {
		response.Response(ctx, 201, 201, nil, "邮箱验证码错误")
		return
	}

	if util.IsNameExist(u.DB, name) {
		response.Response(ctx, 201, 201, nil, "用户名称已被注册")
		return
	}

	// TODO 创建用户
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, 201, 201, nil, "加密错误")
		return
	}
	newUser := model.User{
		Name:       name,
		Email:      email,
		Password:   string(hasedPassword),
		Icon:       "MGA" + strconv.Itoa(rand.Intn(9)+1) + ".jpg",
		Score:      1500,
		LikeNum:    0,
		UnLikeNum:  0,
		CollectNum: 0,
		VisitNum:   0,
		Level:      0,
	}
	u.DB.Create(&newUser)

	// TODO 发放token给前端
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(ctx, 201, 201, nil, "系统异常")
		log.Printf("token generate error: %v", err)
		return
	}
	// TODO 返回结果
	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

// @title    Login
// @description   用户登录
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) Login(ctx *gin.Context) {
	var requestUser = model.User{}
	ctx.Bind(&requestUser)
	// TODO 获取参数
	email := requestUser.Email
	password := requestUser.Password

	// TODO 判断邮箱是否存在
	var user model.User

	u.DB.Where("email = (?)", email).First(&user)
	if user.ID == (uuid.UUID{}) {
		response.Response(ctx, 201, 201, nil, "用户不存在")
		return
	}
	// TODO 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Fail(ctx, nil, "密码错误")
		return
	}
	// TODO 发放token给前端
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, 201, 201, nil, "系统异常")
		log.Printf("token generate error: %v", err)
		return
	}
	log.Printf("token: %v", token)
	// TODO 返回结果
	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

// @title    Security
// @description   进行密码找回的函数
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) Security(ctx *gin.Context) {
	// TODO 数据验证
	var requestUser = vo.UserRequest{}
	ctx.Bind(&requestUser)
	if !util.VerifyEmailFormat(requestUser.Email) {
		response.Response(ctx, 201, 201, nil, "邮箱格式错误")
		return
	}
	// TODO 判断email是否存在
	if !util.IsEmailExist(u.DB, requestUser.Email) {
		response.Response(ctx, 201, 201, nil, "用户不存在")
		return
	}

	// TODO 判断email是否通过验证
	if !util.IsEmailPass(ctx, requestUser.Email, requestUser.Verify) {
		response.Response(ctx, 201, 201, nil, "邮箱验证码错误")
		return
	}

	err := util.SendEmailPass([]string{requestUser.Email})

	// TODO 返回结果
	response.Success(ctx, nil, err)
}

// @title    VerifyEmail
// @description   进行邮箱验证码发送的函数
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) VerifyEmail(ctx *gin.Context) {
	email := ctx.Params.ByName("id")
	// TODO 数据验证
	if !util.VerifyEmailFormat(email) {
		response.Response(ctx, 201, 201, nil, "邮箱格式错误")
		return
	}
	v, err := util.SendEmailValidate([]string{email})
	if err != nil {
		response.Response(ctx, 201, 201, nil, "邮箱验证码发送失败")
		return
	}
	// 验证码存入redis 并设置过期时间5分钟
	util.SetRedisEmail(ctx, email, v)

	// TODO 返回结果
	response.Success(ctx, gin.H{"email": email}, "验证码请求成功")
}

// @title    UpdatePass
// @description   进行密码修改的函数
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) UpdatePass(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	var pairString = vo.PairString{}
	ctx.Bind(&pairString)

	// TODO 获取参数
	oldPass := pairString.First
	newPass := pairString.Second

	// TODO 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPass)); err != nil {
		response.Fail(ctx, nil, "密码错误")
		return
	}

	// TODO 创建密码哈希
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)

	if err != nil {
		response.Response(ctx, 201, 201, nil, "加密错误")
		return
	}

	// TODO 更新密码
	user.Password = string(hasedPassword)

	u.DB.Save(&user)

	response.Success(ctx, nil, "密码修改成功")
}

// @title    Info
// @description   解析上下文中的token并返回user
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	response.Success(ctx, gin.H{"user": vo.ToUserDto(user.(model.User))}, "查看用户成功")
}

// @title    Show
// @description   查看某个用户的信息
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) Show(ctx *gin.Context) {
	id := ctx.Params.ByName("id")

	var user model.User

	// TODO 先看redis中是否存在
	if ok, _ := u.Redis.HExists(ctx, "User", id).Result(); ok {
		cate, _ := u.Redis.HGet(ctx, "User", id).Result()
		if json.Unmarshal([]byte(cate), &user) == nil {
			response.Success(ctx, gin.H{"user": vo.ToUserDto(user)}, "成功")
			return
		} else {
			// TODO 移除损坏数据
			u.Redis.HDel(ctx, "User", id)
		}
	}

	// TODO 查看用户是否在数据库中存在
	if u.DB.Where("id = (?)", id).First(&user).Error != nil {
		response.Fail(ctx, nil, "用户不存在")
		return
	}

	response.Success(ctx, gin.H{"user": vo.ToUserDto(user)}, "成功")

	// TODO 将用户存入redis供下次使用
	v, _ := json.Marshal(user)
	u.Redis.HSet(ctx, "User", id, v)
}

// @title    Update
// @description   修改用户的个人信息
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) Update(ctx *gin.Context) {
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取user
	var requestUser vo.UserUpdate

	// TODO 数据验证
	if err := ctx.ShouldBind(&requestUser); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 查看email是否合法
	if requestUser.Email != user.Email {
		// TODO 判断email是否存在
		if util.IsEmailExist(u.DB, requestUser.Email) {
			response.Response(ctx, 201, 201, nil, "邮箱已绑定")
			return
		}

		// TODO 判断email是否通过验证
		if !util.IsEmailPass(ctx, requestUser.Email, requestUser.Verify) {
			response.Response(ctx, 201, 201, nil, "邮箱验证码错误")
			return
		}
	}

	// TODO 更新信息
	u.DB.Where("id = (?)", user.ID).Updates(model.User{
		Address:       requestUser.Address,
		Blog:          requestUser.Blog,
		Name:          requestUser.Name,
		Email:         requestUser.Email,
		Sex:           requestUser.Sex,
		Icon:          requestUser.Icon,
		ResLong:       requestUser.ResLong,
		ResShort:      requestUser.ResShort,
		BadgeId:       requestUser.BadgeId,
		Theme:         requestUser.Theme,
		MonacoOptions: requestUser.MonacoOptions,
		Language:      requestUser.Language,
		MonacoTheme:   requestUser.MonacoTheme,
	})

	// TODO 移除损坏数据
	u.Redis.HDel(ctx, "User", fmt.Sprint(user.ID))
	response.Success(ctx, nil, "用户信息更新成功")
}

// @title    UpdateLevel
// @description   修改用户的权限等级
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) UpdateLevel(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 获取指定等级
	level, err := strconv.Atoi(ctx.Params.ByName("level"))

	if err != nil || level < 0 {
		response.Fail(ctx, nil, "权限等级有误")
		return
	}

	if level >= user.Level {
		response.Fail(ctx, nil, "权限等级大于等于了你的权限等级")
		return
	}

	// TODO 获取指定用户
	id := ctx.Params.ByName("id")

	var userb model.User

	// TODO 先看redis中是否存在
	if ok, _ := u.Redis.HExists(ctx, "User", id).Result(); ok {
		cate, _ := u.Redis.HGet(ctx, "User", id).Result()
		if json.Unmarshal([]byte(cate), &userb) == nil {
			goto leap
		} else {
			// TODO 移除损坏数据
			u.Redis.HDel(ctx, "User", id)
		}
	}

	// TODO 查看用户是否在数据库中存在
	if u.DB.Where("id = (?)", id).First(&userb).Error != nil {
		response.Fail(ctx, nil, "用户不存在")
		return
	}
	{
		// TODO 将用户存入redis供下次使用
		v, _ := json.Marshal(userb)
		u.Redis.HSet(ctx, "User", id, v)
	}
leap:

	// TODO 查看指定用户的等级
	if userb.Level >= user.Level {
		response.Fail(ctx, nil, "无法修改该用户的权限等级")
		return
	}

	// TODO 更新
	userb.Level = level
	u.DB.Save(&userb)

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    AcceptNum
// @description   查看用户ac题目的数量
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) AcceptNum(ctx *gin.Context) {

	id := ctx.Params.ByName("id")

	// TODO 获取对应用户
	var num int64
	u.DB.Table("records").Select("count(distinct problem_id)").Where("condition = Accepted and user_id = (?)", id).First(&num)

	response.Success(ctx, gin.H{"num": num}, "查看ac题目数量成功")
}

// @title    AcceptRankList
// @description   查看用户ac题目的数量排行列表
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) AcceptRankList(ctx *gin.Context) {

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 扫描结果
	type Result struct {
		AcceptNum int64 `json:"accept_num"`
		UserId    uint  `json:"user_id"`
	}
	var acceptRanks []Result
	var total int64

	// TODO 获取排行数据
	u.DB.Table("records").Select("count(distinct problem_id) as accept_num, user_id").Where("condition = Accepted").Order("accept_num desc").Group("user_id").Count(&total).Offset((pageNum - 1) * pageSize).Limit(pageSize).Scan(&acceptRanks)

	response.Success(ctx, gin.H{"acceptRanks": acceptRanks, "total": total}, "查看用户ac题目的数量排行列表成功")
}

// @title    AcceptRank
// @description   查看用户ac题目的数量排行
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) AcceptRank(ctx *gin.Context) {

	id := ctx.Params.ByName("id")

	// TODO 获取对应用户
	var rank int64
	u.DB.Table("records").Select("rank() over(partition by condition order by count(distinct problem_id) desc)").Where("condition = Accepted and user_id = (?)", id).Group("user_id").First(&rank)

	response.Success(ctx, gin.H{"rank": rank}, "查看用户ac题目的数量排行成功")
}

// @title    ScoreRankList
// @description   查看用户竞赛分数排行列表
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) ScoreRankList(ctx *gin.Context) {

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 扫描结果
	var users []model.User
	u.DB.Order("score desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&users)
	var total int64

	// TODO 获取排行数据
	u.DB.Model(model.User{}).Count(&total)

	var tusers []vo.UserDto

	for i := range users {
		tusers = append(tusers, vo.ToUserDto(users[i]))
	}

	response.Success(ctx, gin.H{"users": tusers, "total": total}, "查看用户竞赛分数排行列表成功")
}

// @title    ScoreRank
// @description   查看用户竞赛分数的数量排行
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) ScoreRank(ctx *gin.Context) {

	id := ctx.Params.ByName("id")

	// TODO 获取对应用户
	var rank int64

	u.DB.Table("users").Select("rank() over(partition by id order by score desc)").Where("id = (?)", id).Group("id").First(&rank)

	response.Success(ctx, gin.H{"rank": rank}, "查看用户竞赛分数排行成功")
}

// @title    HotRankList
// @description   查看用户热度排行列表
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) HotRankList(ctx *gin.Context) {

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 扫描结果
	var users []model.User
	u.DB.Order("(like_num + collect_num - unlike_num) desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&users)
	var total int64

	// TODO 获取排行数据
	u.DB.Model(model.User{}).Count(&total)

	var tusers []vo.UserDto

	for i := range users {
		tusers = append(tusers, vo.ToUserDto(users[i]))
	}

	response.Success(ctx, gin.H{"users": tusers, "total": total}, "查看用户竞赛分数排行列表成功")
}

// @title    HotRank
// @description   查看用户热度的排行
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) HotRank(ctx *gin.Context) {

	id := ctx.Params.ByName("id")

	// TODO 获取对应用户
	var rank int64

	u.DB.Table("users").Select("rank() over(partition by id order by (like_num + collect_num - unlike_num) desc)").Where("id = (?)", id).Group("id").First(&rank)

	response.Success(ctx, gin.H{"rank": rank}, "查看用户竞赛分数排行成功")
}

// @title    LabelCreate
// @description   标签创建
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) LabelCreate(ctx *gin.Context) {
	// TODO 获取指定用户
	id := ctx.Params.ByName("id")

	// TODO 获取标签
	label := ctx.Params.ByName("label")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 创建标签
	userLabel := model.UserLabel{
		Label:  label,
		UserId: user.ID,
	}

	// TODO 插入数据
	if err := u.DB.Create(&userLabel).Error; err != nil {
		response.Fail(ctx, nil, "用户标签上传出错，数据验证有误")
		return
	}

	// TODO 解码失败，删除字段
	u.Redis.HDel(ctx, "UserLabel", id)

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    LabelDelete
// @description   标签删除
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) LabelDelete(ctx *gin.Context) {
	// TODO 获取指定用户
	id := ctx.Params.ByName("id")

	// TODO 获取标签
	label := ctx.Params.ByName("label")

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否可以删除标签
	var userLabel model.UserLabel
	if u.DB.Where("id = (?)", label).First(&userLabel).Error != nil {
		response.Fail(ctx, nil, "标签不存在")
		return
	}

	if userLabel.UserId != user.ID {
		response.Fail(ctx, nil, "标签不属于你")
		return
	}

	// TODO 删除用户标签

	u.DB.Where("id = (?)", label).Delete(&model.UserLabel{})

	// TODO 解码失败，删除字段
	u.Redis.HDel(ctx, "UserLabel", id)

	// TODO 成功
	response.Success(ctx, nil, "删除成功")
}

// @title    LabelShow
// @description   标签查看
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) LabelShow(ctx *gin.Context) {
	// TODO 获取指定用户
	id := ctx.Params.ByName("id")

	// TODO 查找数据
	var userLabels []model.UserLabel
	// TODO 先尝试在redis中寻找
	if ok, _ := u.Redis.HExists(ctx, "UserLabel", id).Result(); ok {
		art, _ := u.Redis.HGet(ctx, "UserLabel", id).Result()
		if json.Unmarshal([]byte(art), &userLabels) == nil {
			goto leap
		} else {
			// TODO 解码失败，删除字段
			u.Redis.HDel(ctx, "UserLabel", id)
		}
	}

	// TODO 在数据库中查找
	u.DB.Where("user_id = (?)", id).Find(&userLabels)
	{
		// TODO 将用户标签存入redis供下次使用
		v, _ := json.Marshal(userLabels)
		u.Redis.HSet(ctx, "UserLabel", id, v)
	}

leap:

	// TODO 成功
	response.Success(ctx, gin.H{"userLabels": userLabels}, "查看成功")
}

// @title    Search
// @description   文本搜索
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) Search(ctx *gin.Context) {
	// TODO 获取文本
	text := ctx.Params.ByName("text")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	var users []model.User

	// TODO 模糊匹配
	u.DB.Where("match(name) against((?) in boolean mode)", text+"*").Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&users)

	// TODO 查看查询总数
	var total int64
	u.DB.Where("match(name) against((?) in boolean mode)", text+"*").Model(model.User{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"users": users, "total": total}, "成功")
}

// @title    SearchLabel
// @description   指定标签的搜索
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) SearchLabel(ctx *gin.Context) {

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
	var userLabels []model.UserLabel

	// TODO 进行标签匹配
	u.DB.Distinct("user_id").Where("label in (?)", requestLabels.Labels).Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&userLabels)

	// TODO 查看查询总数
	var total int64
	u.DB.Distinct("user_id").Where("label in (?)", requestLabels.Labels).Model(model.UserLabel{}).Count(&total)

	// TODO 查找对应用户
	var users []model.User

	var userIds []string

	for i := range userLabels {
		userIds = append(userIds, userLabels[i].UserId.String())
	}

	u.DB.Where("id in (?)", userIds).Find(&users)

	var dtoUsers []vo.UserDto

	for i := range users {
		dtoUsers = append(dtoUsers, vo.ToUserDto(users[i]))
	}

	// TODO 返回数据
	response.Success(ctx, gin.H{"users": dtoUsers, "total": total}, "成功")
}

// @title    SearchWithLabel
// @description   指定标签与文本的搜索
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) SearchWithLabel(ctx *gin.Context) {

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
	var userLabels []model.UserLabel

	// TODO 进行标签匹配
	u.DB.Distinct("user_id").Where("label in (?)", requestLabels.Labels).Find(&userLabels)

	// TODO 查找对应用户
	var users []model.User

	var userIds []string

	for i := range userLabels {
		userIds = append(userIds, userLabels[i].UserId.String())
	}

	// TODO 模糊匹配
	u.DB.Where("id in (?) and match(name) against((?) in boolean mode)", userIds, text+"*").Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&users)

	// TODO 查看查询总数
	var total int64
	u.DB.Where("id in (?) and match(name) against((?) in boolean mode)", userIds, text+"*").Model(model.User{}).Count(&total)

	var dtoUsers []vo.UserDto

	for i := range users {
		dtoUsers = append(dtoUsers, vo.ToUserDto(users[i]))
	}

	// TODO 返回数据
	response.Success(ctx, gin.H{"users": dtoUsers, "total": total}, "成功")
}

// @title    ScoreChange
// @description   指定用户的分数变化
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) ScoreChange(ctx *gin.Context) {

	// TODO 获取用户id
	id := ctx.Params.ByName("id")

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 查找对应用户分数变化
	var userScoreChanges []model.UserScoreChange

	u.DB.Where("user_id = (?)", id).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&userScoreChanges)

	// TODO 查看查询总数
	var total int64
	u.DB.Where("user_id = (?)", id).Model(model.UserScoreChange{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"userScoreChanges": userScoreChanges, "total": total}, "成功")
}

// @title    Hot
// @description   指定用户的今日热度数据
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) Hot(ctx *gin.Context) {

	// TODO 获取用户id
	id := ctx.Params.ByName("id")

	VisitNum, _ := u.Redis.ZScore(ctx, "UserVisit", id).Result()
	LikeNum, _ := u.Redis.ZScore(ctx, "UserLike", id).Result()
	UnLikeNum, _ := u.Redis.ZScore(ctx, "UserUnLike", id).Result()
	CollectNum, _ := u.Redis.ZScore(ctx, "UserCollect", id).Result()

	// TODO 返回数据
	response.Success(ctx, gin.H{"VisitNum": VisitNum, "LikeNum": LikeNum, "UnLikeNum": UnLikeNum, "CollectNum": CollectNum}, "成功")
}

// @title    LikeRank
// @description   用户的今日收到点赞排行
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) LikeRank(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	users, _ := u.Redis.ZRevRangeWithScores(ctx, "UserLike", int64(pageNum-1)*int64(pageSize), int64(pageNum-1)*int64(pageSize)+int64(pageSize)-1).Result()

	total, _ := u.Redis.ZCard(ctx, "UserLike").Result()
	// TODO 返回数据
	response.Success(ctx, gin.H{"users": users, "total": total}, "成功")
}

// @title    UnLikeRank
// @description   用户的今日收到点踩排行
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) UnLikeRank(ctx *gin.Context) {

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	users, _ := u.Redis.ZRevRangeWithScores(ctx, "UserUnLike", int64(pageNum-1)*int64(pageSize), int64(pageNum-1)*int64(pageSize)+int64(pageSize)-1).Result()

	total, _ := u.Redis.ZCard(ctx, "UserUnLike").Result()
	// TODO 返回数据
	response.Success(ctx, gin.H{"users": users, "total": total}, "成功")
}

// @title    CollectRank
// @description   用户的今日收到收藏排行
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) CollectRank(ctx *gin.Context) {

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	users, _ := u.Redis.ZRevRangeWithScores(ctx, "UserCollect", int64(pageNum-1)*int64(pageSize), int64(pageNum-1)*int64(pageSize)+int64(pageSize)-1).Result()

	total, _ := u.Redis.ZCard(ctx, "UserCollect").Result()
	// TODO 返回数据
	response.Success(ctx, gin.H{"users": users, "total": total}, "成功")
}

// @title    VisitRank
// @description   用户的今日收到游览排行
// @auth      MGAronya       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) VisitRank(ctx *gin.Context) {

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	users, _ := u.Redis.ZRevRangeWithScores(ctx, "UserVisit", int64(pageNum-1)*int64(pageSize), int64(pageNum-1)*int64(pageSize)+int64(pageSize)-1).Result()

	total, _ := u.Redis.ZCard(ctx, "UserVisit").Result()
	// TODO 返回数据
	response.Success(ctx, gin.H{"users": users, "total": total}, "成功")
}

// @title    NewUserController
// @description   新建一个IUserController
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   IUserController		返回一个IUserController用于调用各种函数
func NewUserController() IUserController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.User{})
	db.AutoMigrate(model.UserLabel{})
	db.AutoMigrate(model.UserLabel{})
	db.AutoMigrate(model.UserScoreChange{})
	return UserController{DB: db, Redis: redis}
}
