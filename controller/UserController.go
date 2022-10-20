// @Title  UserController
// @Description  该文件用于提供操作用户的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"MGA_OJ/common"
	"MGA_OJ/model"
	"MGA_OJ/response"
	"MGA_OJ/util"
	"MGA_OJ/vo"
	"log"

	"math/rand"

	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// IUserController			定义了用户类接口
type IUserController interface {
	Register(ctx *gin.Context)    // 注册
	Login(ctx *gin.Context)       // 登录
	VerifyEmail(ctx *gin.Context) // 验证码
	Security(ctx *gin.Context)    // 找回密码
	UpdatePass(ctx *gin.Context)  // 更新密码
	Info(ctx *gin.Context)        // 返回当前登录的用户
}

// UserController			定义了题目工具类
type UserController struct {
	DB *gorm.DB // 含有一个数据库指针
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
	if !util.IsEmailPass(email, requestUser.Verify) {
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
	if user.ID == 0 {
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
	if !util.IsEmailPass(requestUser.Email, requestUser.Verify) {
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
	util.SetRedisEmail(email, v)

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

// @title    NewUserController
// @description   新建一个IUserController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   IUserController		返回一个IUserController用于调用各种函数
func NewUserController() IUserController {
	db := common.GetDB()
	db.AutoMigrate(model.User{})
	return UserController{DB: db}
}
