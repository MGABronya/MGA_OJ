// @Title  UserController
// @Description  该文件用于提供操作用户的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
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
	"os"
	"path"

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
	Register(ctx *gin.Context)       // 注册
	Login(ctx *gin.Context)          // 登录
	VerifyEmail(ctx *gin.Context)    // 验证码
	Security(ctx *gin.Context)       // 找回密码
	UpdatePass(ctx *gin.Context)     // 更新密码
	Info(ctx *gin.Context)           // 返回当前登录的用户
	Update(ctx *gin.Context)         // 用户的信息更新
	UpdateIcon(ctx *gin.Context)     // 用户头像更新
	UpdateLevel(ctx *gin.Context)    // 修改用户的等级
	Show(ctx *gin.Context)           // 显示用户的所有信息
	AcceptNum(ctx *gin.Context)      // 显示用户ac题目数量
	AcceptRankList(ctx *gin.Context) // 显示用户ac题目的排行列表
	AcceptRank(ctx *gin.Context)     // 显示用户ac题目的排行
}

// UserController			定义了题目工具类
type UserController struct {
	DB    *gorm.DB      // 含有一个数据库指针
	Redis *redis.Client // 含有一个redis指针
}

// @title    Register
// @description   用户注册
// @auth      MGAronya（张健）       2022-9-16 12:15
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
		Name:     name,
		Email:    email,
		Password: string(hasedPassword),
		Icon:     "MGA" + strconv.Itoa(rand.Intn(9)+1) + ".jpg",
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
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) Login(ctx *gin.Context) {
	var requestUser = model.User{}
	ctx.Bind(&requestUser)
	// TODO 获取参数
	email := requestUser.Email
	password := requestUser.Password
	// TODO 数据验证
	if !util.VerifyEmailFormat(email) {
		response.Response(ctx, 201, 201, nil, "邮箱格式错误")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, 201, 201, nil, "密码不能少于6位")
		return
	}
	// TODO 判断邮箱是否存在
	var user model.User

	u.DB.Where("email = ?", email).First(&user)
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
	// TODO 返回结果
	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

// @title    Security
// @description   进行密码找回的函数
// @auth      MGAronya（张健）       2022-9-16 12:15
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
// @auth      MGAronya（张健）       2022-9-16 12:15
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
// @auth      MGAronya（张健）       2022-9-16 12:15
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
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	response.Success(ctx, gin.H{"user": vo.ToUserDto(user.(model.User))}, "查看用户成功")
}

// @title    Show
// @description   查看某个用户的信息
// @auth      MGAronya（张健）       2022-9-16 12:15
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
	if u.DB.Where("id = ?", id).First(&user).Error != nil {
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
// @auth      MGAronya（张健）       2022-9-16 12:15
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

	user.Address = requestUser.Address
	user.Blog = requestUser.Blog
	user.Name = requestUser.Name
	user.Email = requestUser.Email
	user.Sex = requestUser.Sex

	// TODO 更新信息
	u.DB.Save(&user)

	// TODO 移除损坏数据
	u.Redis.HDel(ctx, "User", fmt.Sprint(user.ID))
	response.Success(ctx, nil, "用户信息更新成功")
}

// @title    UpdateIcon
// @description   个人头像图片上传
// @auth      MGAronya（张健）       2022-9-16 12:31
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) UpdateIcon(ctx *gin.Context) {
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	file, err := ctx.FormFile("file")

	//TODO 数据验证
	if err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	extName := path.Ext(file.Filename)
	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
	}

	// TODO 格式验证
	if _, ok := allowExtMap[extName]; !ok {
		response.Fail(ctx, nil, "文件格式有误")
		return
	}

	// TODO 非默认图片，则删除原本地文件
	if !util.VerifyIconFormat(user.Icon) {
		if err := os.Remove("./Icon/" + user.Icon); err != nil {
			// TODO 如果删除失败则输出 file remove Error!
			fmt.Println("file remove Error!")
			// TODO 输出错误详细信息
			fmt.Printf("%s", err)
		} else {
			// TODO 如果删除成功则输出 file remove OK!
			fmt.Print("file remove OK!")
		}
	}
	file.Filename = user.ID.String() + extName

	// TODO 将文件存入本地
	ctx.SaveUploadedFile(file, "./Icon/"+file.Filename)

	u.DB.Where("id = ?", user.ID).Take(&user)

	user.Icon = file.Filename

	u.DB.Save(&user)

	// TODO 移除损坏数据
	u.Redis.HDel(ctx, "User", fmt.Sprint(user.ID))

	response.Success(ctx, gin.H{"Icon": user.Icon}, "更新成功")
}

// @title    UpdateLevel
// @description   修改用户的权限等级
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) UpdateLevel(ctx *gin.Context) {
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	id := ctx.Params.ByName("id")
	level, _ := strconv.Atoi(ctx.Params.ByName("level"))

	if user.Level <= level {
		response.Fail(ctx, nil, "用户权限不足")
		return
	}

	// TODO 获取对应用户
	var usera model.User
	if u.DB.Where("id = ?", id).First(&usera).Error != nil {
		response.Fail(ctx, nil, "用户不存在")
		return
	}

	usera.Level = level

	// TODO 更新信息
	u.DB.Save(&user)

	// TODO 移除损坏数据
	u.Redis.HDel(ctx, "User", id)

	response.Success(ctx, nil, "用户信息更新成功")
}

// @title    AcceptNum
// @description   查看用户ac题目的数量
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) AcceptNum(ctx *gin.Context) {

	id := ctx.Params.ByName("id")

	// TODO 获取对应用户
	var num int64
	u.DB.Table("records").Select("count(distinct problem_id)").Where("condition = Accepted and user_id = ?", id).First(&num)

	response.Success(ctx, gin.H{"num": num}, "查看ac题目数量成功")
}

// @title    AcceptRankList
// @description   查看用户ac题目的数量排行列表
// @auth      MGAronya（张健）       2022-9-16 12:15
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
	u.DB.Table("records").Select("count(distinct problem_id) as accept_num, user_id").Where("condition = Accepted").Order("accept_num desc").Group("user_id").Offset((pageNum - 1) * pageSize).Limit(pageSize).Scan(&acceptRanks).Count(&total)

	response.Success(ctx, gin.H{"acceptRanks": acceptRanks, "total": total}, "查看用户ac题目的数量排行列表成功")
}

// @title    AcceptRank
// @description   查看用户ac题目的数量排行
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (u UserController) AcceptRank(ctx *gin.Context) {

	id := ctx.Params.ByName("id")

	// TODO 获取对应用户
	var rank int64
	u.DB.Table("records").Select("rank() over(partition by condition order by count(distinct problem_id) desc)").Where("condition = Accepted and user_id = ?", id).Group("user_id").First(&rank)

	response.Success(ctx, gin.H{"rank": rank}, "查看用户ac题目的数量排行成功")
}

// @title    LabelCreate
// @description   标签创建
// @auth      MGAronya（张健）       2022-9-16 12:20
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
// @auth      MGAronya（张健）       2022-9-16 12:20
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
	if u.DB.Where("id = ?", label).First(&userLabel).Error != nil {
		response.Fail(ctx, nil, "标签不存在")
		return
	}

	if userLabel.UserId != user.ID {
		response.Fail(ctx, nil, "标签不属于你")
		return
	}

	// TODO 删除用户标签

	u.DB.Where("id = ?", label).Delete(&model.UserLabel{})

	// TODO 解码失败，删除字段
	u.Redis.HDel(ctx, "UserLabel", id)

	// TODO 成功
	response.Success(ctx, nil, "删除成功")
}

// @title    LabelShow
// @description   标签查看
// @auth      MGAronya（张健）       2022-9-16 12:20
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
	u.DB.Where("user_id = ?", id).Find(&userLabels)
	{
		// TODO 将用户标签存入redis供下次使用
		v, _ := json.Marshal(userLabels)
		u.Redis.HSet(ctx, "UserLabel", id, v)
	}

leap:

	// TODO 成功
	response.Success(ctx, gin.H{"userLabels": userLabels}, "查看成功")
}

// @title    NewUserController
// @description   新建一个IUserController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   IUserController		返回一个IUserController用于调用各种函数
func NewUserController() IUserController {
	db := common.GetDB()
	redis := common.GetRedisClient(0)
	db.AutoMigrate(model.User{})
	db.AutoMigrate(model.UserLabel{})
	return UserController{DB: db, Redis: redis}
}
